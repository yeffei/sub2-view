# 智能路由与故障转移验收清单

日期：2026-07-02  
时区：Asia/Shanghai

## 1. 目标口径

对外这句话：

`自研负载均衡、自动故障转移与智能路由，守住高并发下的稳定与低延迟。`

要拆成两层理解：

- 功能层：系统确实会分流、切换、恢复、重试、观测。
- 保证层：系统已经通过足够的运行证据和压测，能对“高并发稳定、低延迟”做更强承诺。

当前项目已经具备较完整的功能层主干，但保证层还没有完全闭环。

## 2. 当前现状结论

### 2.1 已经具备

1. 池化调度主链路已落地
   - `Group -> PoolBinding -> Pool -> PoolMember -> Account` 已经进入 OpenAI 选路链路。
   - 分组绑定到池后，请求只从该池成员中选账号。

2. 成员级调度覆盖已接入实际调度
   - `priority_override`
   - `weight`
   - `max_concurrency_override`
   - `schedulable_override`

3. 池级策略已进入实际路由
   - `scheduler_mode`
   - `load_balance_enabled`
   - `failover_enabled`
   - `top_k`
   - `max_failover_hops`

4. 高级调度器已使用多维信号打分
   - 优先级
   - 当前负载
   - 等待队列
   - 错误率 EWMA
   - TTFT EWMA
   - 会话窗口 reset 因子

5. 自动故障转移主链路已存在
   - 上游返回可切换错误时，会抛出 `UpstreamFailoverError`
   - Handler 层会在单次请求内继续尝试后续候选
   - 受 `max_failover_hops` 控制最大切换跳数

6. 被动隔离与主动恢复都已接入
   - `rate_limit`
   - `overload`
   - `temp_unschedulable`
   - `status=error`
   - 上游池恢复探针会自动探测异常成员并在成功后恢复调度

7. 可观测链路基础已具备
   - `RoutingExplanation`
   - Scheduler metrics snapshot
   - runtime status
   - channel monitor

### 2.2 只能算“部分具备”

1. 智能低延迟选路
   - 已经有 TTFT / 错误率信号进入评分。
   - 但还缺“经过压测验证的慢线路治理规则”与稳定阈值。
   - 目前更接近“具备智能调度基础”，不是“低延迟已被证明能守住”。

2. 高并发稳定性承诺
   - 有并发、等待队列、failover、临时摘除这些机制。
   - 但还缺正式压测基线、SLO 和回归门槛。
   - 所以不能说“已经保证高并发稳定”。

3. 前台监控与真实调度的一致性
   - 已经从“展示单账号状态”往“池视角”收敛。
   - 但监控页本身不是调度真相来源，只能算辅助可视化。

### 2.3 还不能对外强承诺的点

1. 没有正式高并发压测报告。
2. 没有明确的延迟 SLO 验收阈值。
3. 没有形成慢线路自动降权的完整治理闭环。
4. 没有把“调度效果”沉淀成固定回归验证流程。

### 2.4 2026-07-02 本地运行态抽样

本地 PostgreSQL 运行库抽样显示：

1. 绑定关系
   - `Codex Plus (group_id=2) -> Codex Plus (pool_id=1)`
   - `Codex Pro (group_id=3) -> Codex Plus (pool_id=1)`

2. `group_id=2` 在最近 12 小时的命中分布
   - `account_id=3`：132 次，平均 `TTFT=16307ms`
   - `account_id=4`：105 次，平均 `TTFT=12924ms`
   - `account_id=8`：68 次，平均 `TTFT=16879ms`

3. 结论
   - 请求确实已经在池成员 `3/4/8` 中分流，不是落到池外账号。
   - 但首 token 波动仍然很大，存在 `40s~60s` 级别慢请求。
   - 因此可以证明“池化分流已生效”，还不能证明“低延迟已守住”。

## 3. 代码证据定位

### 3.1 上游池与策略接入

- `backend/internal/service/upstream_pool.go`
- `backend/internal/service/openai_gateway_service.go`
- `backend/internal/service/account.go`
- `backend/internal/repository/upstream_pool_repo.go`

### 3.2 高级调度与评分

- `backend/internal/service/openai_account_scheduler.go`

当前已能看到：

- 多层选路：`previous_response_id -> session sticky -> load balance`
- TopK 候选池
- weighted selection order
- `errorRate` / `ttft` EWMA
- sticky escape

### 3.3 自动故障转移

- `backend/internal/service/openai_gateway_service.go`
- `backend/internal/handler/failover_loop.go`

### 3.4 恢复探针

- `backend/internal/service/upstream_pool_recovery_probe_runner_service.go`
- `backend/internal/service/ratelimit_service.go`

### 3.5 路由解释与监控

- `backend/internal/service/routing_explanation.go`
- `backend/internal/service/channel_monitor_service.go`
- `backend/internal/service/channel_monitor_runner.go`

## 4. 最终验收清单

只有下面 5 组都过，才可以把那句宣传语当成“基本站得住”。

### A. 功能闭环验收

1. 分组绑定池后，请求不会再落到池外账号。
2. 成员 `priority_override` 能显著改变命中顺序。
3. 成员 `weight` 能显著影响同层候选分流比例。
4. `max_concurrency_override` 能限制单成员并发占用。
5. `schedulable_override=false` 的成员不会参与调度。

### B. 故障转移验收

1. 单个成员发生 `401/403/429/529/transport error` 时，请求能切到后续候选。
2. 切换次数受 `max_failover_hops` 约束。
3. 被判异常的成员会被摘出调度，而不是持续污染流量。
4. 恢复后的成员会重新进入候选集。

### C. 低延迟验收

1. 高级调度开启时，TTFT 更差的线路在同池竞争中会被弱化。
2. sticky account 明显变慢时，`sticky_escape` 能触发逃逸。
3. 同池多成员下，首 token 延迟不能长期集中在明显差线路上。
4. 慢线路恢复后，不会立刻重新吃满流量。

### D. 可观测验收

1. 管理台能看清池、成员、绑定、运行态。
2. 路由解释能反映本次请求走了哪一层决策。
3. 能从日志中追到：
   - 选中了谁
   - 为什么切换
   - 切到了谁
   - 是否恢复成功

### E. 高并发验收

1. 有固定压测脚本。
2. 有固定压测场景：
   - 单池双成员
   - 单池多成员
   - 单成员故障中途切换
   - 部分成员慢但不报错
3. 有明确阈值：
   - 成功率
   - P95 / P99 TTFT
   - 切换成功率
   - 平均切换耗时
4. 压测结果能重复，不是偶发一次成功。

## 5. 当前缺口清单

### P0：必须补

1. 压测基线
   - 现在缺的不是功能名词，而是证据。
   - 必须有真实并发压测和结果留档。

2. 调度效果回归验证
   - 要能证明 `priority / weight / top_k / failover` 在真实请求里生效。
   - 建议直接基于 usage logs 或 routing logs 做命中分布验证。

3. 慢线路治理规则定稿
   - 不是简单固定阈值封杀。
   - 应限定在“同池、同模型、相近负载”的对比下做相对降权。

### P1：应该补

1. 池级观测面板
   - 展示池命中分布
   - 成员近期成功率
   - 成员近期 TTFT
   - 最近切换次数

2. 调度结果留痕
   - 每次请求至少能还原：
     - 候选数量
     - TopK
     - 命中成员
     - 是否 failover
     - 是否 sticky escape

3. 恢复探针结果可视化
   - 最后探测时间
   - 最后成功时间
   - 连续失败次数

### P2：优化项

1. 池策略模板化
2. 更精细的模型级路由
3. 更完整的用户侧“池视角”监控展示

## 6. 现在可以怎么描述

当前更准确的产品口径是：

`已具备自研负载均衡、自动故障转移和基础智能路由能力；高并发稳定与低延迟保障能力正在通过压测与调度治理进一步收口。`

如果要升级成更强口径，至少要先完成：

1. 压测基线
2. 慢线路治理
3. 调度效果回归验证

## 7. 下一步建议顺序

1. 先做“调度效果验证”
   - 验证池内命中分布和 failover 是否符合配置

2. 再做“压测基线”
   - 固定场景
   - 固定指标
   - 固定结果归档

3. 最后再补“慢线路治理”
   - 这一步要建立在前两步证据之上

## 8. 相关约束文档

为避免后续调度优化把缓存率打散，缓存优先规则单独沉淀在：

- [cache-first-routing-policy.md](D:/sub2api/docs/cache-first-routing-policy.md)

慢线路治理阈值、分级动作与当前实现缺口，单独沉淀在：

- [slow-line-governance-policy.md](D:/sub2api/docs/slow-line-governance-policy.md)
