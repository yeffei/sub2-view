# 智能路由与自动故障转移设计说明

日期：2026-06-30  
时区：Asia/Shanghai

## 1. 背景

网站前台已经对外承诺：

`自研智能负载均衡功能，自动故障转移与智能路由，确保高并发下的服务稳定性与低延迟。`

这句话不能只停留在页面文案或后台演示页上，必须对应一套真实可用的后端闭环能力：

- 调度能力层：请求不是随机选账号，而是按能力、负载、错误率、时延和会话粘性做选择。
- 故障转移层：上游账号异常时，系统自动切换到下一可用账号，而不是直接把失败暴露给用户。
- 健康隔离层：异常账号可以自动临时摘除，并在合适时机恢复。
- 控制面层：管理台可以看到池、成员、健康、路由结果和切换情况，并能手动干预。
- 可观测层：调度决策、失败原因、切换路径和成员状态要可追踪。

## 2. 本次范围边界

本次设计只承诺先把 `OpenAI` 主链路做成完整闭环，其他平台先复用思路，不在这一轮承诺全部落地。

本次范围内：

- 以现有 `Account` 体系为真实资源实体。
- 在 `Account` 之上增加“调度集合层”，而不是再造一套独立的上游中转系统。
- 管理台围绕“池、成员、健康、路由解释、手动摘除/恢复”建设。
- 优先复用现有 `OpenAI` 调度器、粘性会话、故障转移、临时摘除和路由解释能力。

本次范围外：

- 不重做所有平台的统一调度内核。
- 不新建一套与 `Account` 平行的账号明细系统。
- 不在这一轮承诺完整的跨平台全局池编排。

## 3. 现状判断

当前后端并不是从零开始，`OpenAI` 主链路已经有较完整的基础能力，核心实现集中在：

- `backend/internal/service/openai_gateway_service.go`
- `backend/internal/service/openai_account_scheduler.go`
- `backend/internal/service/routing_explanation.go`
- `backend/internal/service/gateway_service.go`
- `backend/internal/service/ratelimit_service.go`

已经确认存在的能力：

- 自动故障转移错误类型：`UpstreamFailoverError`
- 上游错误判定：`shouldFailoverUpstreamError`、`shouldFailoverOpenAIUpstreamResponse`
- 会话粘性：`BindStickySession`
- 账号选择：`selectAccountForModelWithExclusions`、`selectBestAccount`
- 可调度候选筛选：`listSchedulableAccounts`
- 调度打分与选择：`TopK`、weighted order、`load_skew`
- 粘性逃逸：`sticky_escape_triggered`
- 运行态观测：`RecordOpenAIAccountResult`、`SnapshotMetrics`
- 临时摘除：`temp_unschedulable`
- 路由解释：`RoutingExplanation`
- 调度开关：`openai_advanced_scheduler_enabled`

结论：底层能力已具备较强基础，但“对外可承诺”的产品闭环还未完成，主要缺的是领域模型、控制面、统一入口和验收闭环。

## 4. 核心设计原则

### 4.1 上游池不是新的上游服务

这里的“上游池”本质上是调度集合层，可以理解为“上游中转成员集合”，但不是一套新的代理服务，也不是复制一份账号系统。

真实账号仍然是 `Account`。

### 4.2 账号管理与调度管理解耦

- `Account` 负责凭据、平台、模型能力、并发、分组、限流、启停等真实资源属性。
- `Pool` 负责调度边界、策略和运维视角。
- `PoolMember` 负责把某个 `Account` 放入某个 `Pool`，并附加成员级别调度权重或状态。

这样能兼容你后续“买号做号池”的场景：号源继续走 `SUB2API` 现有账号管理，池只是把这些号组织起来参与调度。

### 4.3 先做 OpenAI 闭环，再抽象跨平台

这轮不追求一次性抽象出所有平台通用池，只先把 `OpenAI` 跑通，避免抽象过早。

## 5. 领域模型

推荐关系：

`Group -> PoolBinding -> Pool -> PoolMember -> Account`

### 5.1 Pool

表示一个可被调度的逻辑池。

建议字段：

- `id`
- `name`
- `code`
- `platform`
- `description`
- `enabled`
- `scheduler_mode`
- `policy_id`
- `default_required_capability`
- `default_required_transport`
- `sticky_enabled`
- `sticky_ttl_seconds`
- `sticky_escape_enabled`
- `sticky_escape_error_rate_threshold`
- `sticky_escape_ttft_ms_threshold`
- `load_balance_enabled`
- `failover_enabled`
- `top_k`
- `max_failover_hops`
- `wait_timeout_ms`
- `max_waiting`
- `created_at`
- `updated_at`

说明：

- `platform` 本轮固定先落 `openai`。
- `scheduler_mode` 可先保留 `basic` / `advanced` 两档，对应是否启用高级调度。
- `policy_id` 指向独立策略对象，便于未来多个池复用策略。

### 5.2 PoolMember

表示池中的一个成员，本质上引用现有 `Account`。

建议字段：

- `id`
- `pool_id`
- `account_id`
- `enabled`
- `schedulable_override`
- `manual_drained`
- `weight`
- `priority_override`
- `max_concurrency_override`
- `notes`
- `joined_at`
- `updated_at`

说明：

- `enabled` 是成员级开关。
- `schedulable_override` 用于手工强制摘除或恢复，不直接污染账号原始信息。
- `manual_drained` 用于人工下线一个成员，但保留配置。
- `weight` 和 `priority_override` 用于成员级调度偏置。

### 5.3 PoolBinding

表示某个业务分组与池的绑定关系。

建议字段：

- `id`
- `group_id`
- `pool_id`
- `platform`
- `models`
- `request_path_scope`
- `priority`
- `enabled`
- `created_at`
- `updated_at`

说明：

- 一个 `Group` 后续理论上可以绑定多个池，但本轮可先支持“每个平台一个主池”。
- `models` 可用于模型级路由覆盖，例如某些模型只走某个池。

### 5.4 PoolPolicy

表示一套可复用的调度与故障策略。

建议字段：

- `id`
- `name`
- `platform`
- `load_balance_enabled`
- `weights_json`
- `top_k`
- `sticky_enabled`
- `sticky_ttl_seconds`
- `sticky_escape_config_json`
- `failover_status_codes_json`
- `failover_message_rules_json`
- `temp_unschedulable_rules_json`
- `retry_same_account_rules_json`
- `wait_plan_json`
- `recovery_policy_json`
- `created_at`
- `updated_at`

说明：

- 当前代码里很多策略已存在于设置、运行态或硬编码逻辑中，这一层的目标是把它们显式化、可管理化。

## 6. 与现有实体的关系

### 6.1 Account

仍然是唯一真实上游资源：

- 凭据保存在 `Account`
- 能力映射保存在 `Account`
- 并发、优先级、平台类型和模型支持保存在 `Account`
- 运行态错误、限流和临时摘除仍以 `Account` 为核心对象

### 6.2 Group

`Group` 不再直接理解为“绑一堆账号”，而是“绑定某个调度池，由池再选成员”。

兼容策略：

- 旧链路仍允许从 `Group -> Account` 直接工作，避免一次性切断现有逻辑。
- 新链路优先 `Group -> PoolBinding -> Pool -> Members -> Account`。
- 过渡期可以将“分组内的可调度账号集合”视为池候选的默认来源。

## 7. 调度决策流程

以下流程以 `OpenAI` 为准。

### 7.1 请求入口

输入维度：

- `group_id`
- `requested_model`
- `required_capability`
- `required_transport`
- `require_compact`
- `session_hash`
- `previous_response_id`

### 7.2 绑定解析

1. 先解析 `Group` 对应的 `PoolBinding`。
2. 若存在平台级或模型级绑定，则取绑定池。
3. 若还未建池绑定，则回退到现有分组账号集合逻辑。

### 7.3 候选筛选

候选必须同时满足：

- `Account` 平台匹配 `openai`
- 成员在 `PoolMember` 维度启用
- 账号在 `Account` 维度可调度
- 不处于 `temp_unschedulable`
- 具备目标模型或可映射能力
- 满足 `required_capability`
- 满足 `required_transport`
- 未命中显式排除列表 `excludedIDs`
- 若开启渠道限制，也必须满足上游模型限制

### 7.4 会话粘性层

优先级顺序：

1. `previous_response_id` 粘性命中
2. `session_hash` 粘性命中
3. 负载均衡选择

粘性命中后，仍需再校验：

- 账号是否仍可调度
- 是否仍属于当前 `Group`
- 是否仍满足模型、能力、传输协议要求
- 是否已被人工摘除
- 是否命中粘性逃逸条件

### 7.5 粘性逃逸

当粘性账号虽然还活着，但已经明显变差时，应主动逃逸到重新调度：

- 错误率超过阈值
- `TTFT` 超过阈值
- 并发槽已满且等待不可接受

当前代码里已存在 `sticky_escape_triggered` 基础能力，本次要把它纳入正式产品行为，而不是仅作为隐性实现细节。

### 7.6 负载评分

在候选集中计算综合分数，建议以现有实现为基础，继续使用这些因素：

- `priority`
- `load_rate`
- `waiting_count`
- `error_rate`
- `ttft`
- `queue`
- `reset`

建议统一叫做“综合调度分”，而不是对外暴露过多实现术语。

建议默认权重方向：

- 优先低负载
- 优先低错误率
- 优先低等待
- 优先低首字延迟
- 在会话窗口即将重置时，适度提升将过期资源的利用优先级

### 7.7 TopK 候选选择

不是在全量候选中完全随机，而是：

1. 按综合分排序
2. 取前 `TopK`
3. 在 `TopK` 内做加权选择

目的：

- 避免长期固定打在单个“第一名”账号上
- 也避免高并发下退化成纯随机选号
- 在稳定性与分散度之间取得平衡

### 7.8 并发槽与等待策略

高并发下不能只看分数，还要看“现在能不能接单”。

处理顺序：

1. 尝试直接获取账号并发槽
2. 若获取成功，立即使用该账号
3. 若获取失败但可等待，则生成 `WaitPlan`
4. 若等待不值得，则切到下一候选

等待策略建议以池策略显式化：

- `wait_timeout_ms`
- `max_waiting`
- `sticky_session_wait_timeout`
- `sticky_session_max_waiting`

这样可以做到：

- 粘性账号允许短等待，减少无意义切换
- 普通负载均衡则更偏向快速选下一候选，降低整体尾延迟

### 7.9 failover 链路

当已选账号失败后：

1. 将当前账号加入 `excludedIDs`
2. 记录失败结果与运行态指标
3. 如命中临时摘除策略，则标记 `temp_unschedulable`
4. 在同池候选内重选下一账号
5. 若超过 `max_failover_hops`，再向上返回失败

### 7.10 路由解释

每次路由完成后都应生成解释对象，沿用并扩展现有 `RoutingExplanation`。

至少包含：

- 选中的账号
- 命中的层级
- 候选数
- `TopK`
- 是否 fallback
- 是否等待
- 被跳过的原因统计
- 当前请求需要的 transport / capability

## 8. 自动故障转移规则

### 8.1 应触发 failover 的错误

当前 `OpenAI` 主链路建议沿用现有规则，并明确纳入产品标准：

- `401`
- `402`
- `403`
- `429`
- `529`
- `5xx`
- 上游瞬时处理异常
- 上游 transport error
- 流式响应在未向客户端写出字节前就中断

说明：

- 这些错误说明“当前账号此刻不可继续承载该请求”，优先切到下一候选。

### 8.2 只透传、不切换的错误

默认不应 failover 的情况：

- 业务参数错误
- 请求体不合法
- 明确的用户侧输入问题
- 已经向客户端输出了不可回滚的流内容

核心原则：

- 只有当错误更可能属于“账号或上游状态问题”时才切换
- 不能把用户请求错误伪装成系统智能切换

### 8.3 同账号重试

对于少量“偶发性、同账号短暂重试可能恢复”的错误，允许先在同账号重试，再切其他账号。

例如：

- 概率性空响应
- 某些瞬时 `400` 兼容错误

这部分沿用 `RetryableOnSameAccount` 语义，但需要在池策略层显式化。

### 8.4 临时摘除

以下情况应触发 `temp_unschedulable`：

- 命中配置化 `temp_unschedulable_rules`
- 持续 `401/403/429`
- 明显的流超时或空响应
- 已知会污染后续请求成功率的账号错误

要求：

- 临时摘除必须记录到账号维度
- 原因、状态码、触发时间和恢复时间必须可见
- 摘除后同池调度应立即跳过该成员

### 8.5 恢复策略

恢复优先采用“两类恢复”：

- 自动恢复：`temp_unschedulable_until` 到期后自动回到候选集
- 手动恢复：管理员在控制台主动解除摘除

后续可加：

- 主动健康探测恢复
- 基于连续成功次数恢复权重

## 9. 管理台要求

现有 [D:\sub2api\frontend\src\views\admin\UpstreamPoolsView.vue](D:\sub2api\frontend\src\views\admin\UpstreamPoolsView.vue) 已经是一个不错的入口，但当前更偏“成员列表页”，还不够支撑完整控制面。

本轮完成后，管理台至少应支持以下视图能力。

### 9.1 池列表

字段建议：

- 池名称
- 平台
- 绑定分组数
- 成员总数
- 可调度成员数
- 临时摘除成员数
- 当前调度模式
- 近一小时 failover 次数
- 开关状态

### 9.2 成员列表

字段建议：

- 账号名称
- 账号类型
- 所属平台
- 所属池
- 当前是否可调度
- 是否手动摘除
- 当前负载
- 当前等待数
- 错误率
- `TTFT`
- 最近失败时间
- `temp_unschedulable_until`
- 最近 failover 次数

### 9.3 详情视图

详情至少包含：

- 绑定的 `Group`
- 成员健康状态
- 成员能力覆盖
- 当前策略摘要
- 最近路由解释样本
- 最近 failover 事件

### 9.4 可操作项

管理台至少要能操作：

- 开启/关闭高级调度
- 手动摘除成员
- 手动恢复成员
- 调整成员权重
- 调整池 `TopK`
- 调整粘性与逃逸阈值
- 查看路由解释

## 10. 对网站承诺的“完成定义”

只有同时满足下面几点，才算网站上的那句能力已经真正实现：

1. 当某个上游账号故障、限流、失效或超时时，请求会自动切换到同池其他可用账号。
2. 账号选择基于负载、错误率、时延、粘性和并发，而不是随机选号。
3. 异常账号会被自动临时摘除，不会持续污染后续流量。
4. 管理台能看到池、成员、健康、failover 和路由解释，并能手动干预。
5. 在高并发下，系统仍能通过并发槽与等待策略控制尾延迟，而不是简单地把请求打爆到某一个账号。

## 11. 验收标准

### 11.1 单成员故障自动切走

给定一个池中有多个 `OpenAI` 成员：

- 当当前成员返回 `429`、`5xx`、transport error 或被判定为上游瞬时故障时
- 请求应自动切到下一候选
- 用户侧不需要手工重试

### 11.2 高并发下不退化为随机选号

在高并发压测下：

- 账号分布应体现调度评分结果
- 热门账号不会被无限集中打爆
- 候选也不会退化成完全随机轮盘

### 11.3 异常账号自动退出调度

当某账号连续触发临时摘除规则时：

- 应自动标记为 `temp_unschedulable`
- 调度链路应立即跳过它
- 到期后自动恢复或允许手动恢复

### 11.4 粘性会话不是死绑

当粘性账号错误率或 `TTFT` 明显恶化时：

- 系统应触发粘性逃逸
- 后续请求应重新参与负载均衡

### 11.5 管理台可见且可操作

管理员应能在后台直接看到：

- 当前池和成员关系
- 成员是否健康
- 最近切换次数
- 最近路由解释
- 手动摘除/恢复结果

## 12. 分阶段落地建议

### Phase 1：OpenAI 闭环补齐

- 建立 `Pool / PoolMember / PoolBinding / PoolPolicy` 数据模型
- 将现有 `OpenAI` 调度与 failover 能力挂到池控制面
- 管理台补齐池、成员、健康、解释和手动操作

### Phase 2：指标与可观测增强

- 汇总 failover 事件
- 汇总成员健康趋势
- 路由解释查询化
- 最近异常原因聚合

### Phase 3：扩展到其他平台

- 抽象可复用策略层
- 将 `Gemini / Anthropic / Antigravity` 逐步接入

## 13. 当前实现建议

进入开发阶段时，建议遵循以下顺序：

1. 先补数据库与仓储层：池、成员、绑定、策略。
2. 再把 `OpenAI` 调度入口从“直接查账号”升级为“先解析池，再选成员账号”。
3. 然后补管理台真实接口，而不是继续只用 mock 数据。
4. 最后补验收测试：failover、粘性逃逸、临时摘除、并发等待、路由解释。

## 14. 本文结论

这次要实现的不是一个“上游池页面”，而是一套真实可运行的调度闭环。

针对当前项目，最合适的实现路径不是重造上游体系，而是在现有 `Account` 之上补一层 `Pool` 调度模型，把已经存在的 `OpenAI` 智能路由、自动故障转移、临时摘除和路由解释能力正式产品化。

这也是网站文案能够被真正兑现的最短路径。
