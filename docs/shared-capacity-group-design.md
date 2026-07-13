# SST 上游共享容量组设计

状态：待架构确认，尚未进入迁移或业务代码实现。

日期：2026-07-12（Asia/Shanghai）

## 1. 目标与边界

共享容量组用于表达同一上游来源下多个账号共同受一个硬并发额度约束。例如三个 API Key 共享总并发 3000：成员可以借用其他成员的空闲份额，但所有成功获取的槽位总数不得超过 3000。

本设计不表达自动扩容，也不把整个上游池视为一个容量组。容量组边界落在 `upstream_account_sets`；同一池可以同时包含多个供应方、多个容量组及普通独立账号。

管理员可以看到精确容量、并发、等待与历史压力。用户 `/monitor` 只消费后端映射后的粗粒度状态，不返回账号名、集合名、精确额度或内部成员数。

## 2. 已确认的现有实现

- `upstream_account_sets`、`upstream_account_set_members` 和 `upstream_pool_member_sets` 已由迁移 `156_upstream_account_sets.sql` 建立。
- 集合成员在 `upstream_pool_repo.go` 中展开；管理查询保留 `source_set_id`，但运行时 `UpstreamPoolResolvedMemberConfig` 目前没有集合或容量身份。
- 同一账号同时以直接成员和集合成员出现时，查询使用 `DISTINCT ON (account_id)`，直接成员优先。因此容量身份不能依赖最终的 `source_type`，否则直接成员会绕过共享上限。
- 当前账号并发槽位为 `concurrency:account:{accountID}` ZSET，由单 key Lua 原子 acquire；Redis Cluster 兼容是现有明确约束。
- 当前最新迁移号为 `159`。实现阶段的新迁移使用下一个未占用编号；提交前必须重新检查编号冲突。

## 3. 数据模型

### 3.1 集合级配置

在 `upstream_account_sets` 增加：

```sql
shared_concurrency_limit INTEGER NULL
```

语义：

- `NULL`：普通账号集合，不参与共享容量控制。
- `> 0`：该集合同时是共享容量组，数值是上游共享硬上限。
- 不增加独立 `enabled` 字段；继续使用集合现有 `enabled`，避免出现两套启停状态。

约束：

```sql
CHECK (shared_concurrency_limit IS NULL OR shared_concurrency_limit > 0)
```

### 3.2 容量成员配置

新增表 `upstream_account_set_capacity_members`，复用账号集合作为容量组载体，但不把普通集合成员自动解释为容量成员：

```sql
CREATE TABLE upstream_account_set_capacity_members (
    account_id BIGINT PRIMARY KEY,
    set_id BIGINT NOT NULL,
    hard_concurrency_limit INTEGER NULL,
    soft_concurrency_share INTEGER NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (set_id, account_id)
        REFERENCES upstream_account_set_members(set_id, account_id)
        ON DELETE CASCADE,
    CHECK (hard_concurrency_limit IS NULL OR hard_concurrency_limit > 0),
    CHECK (soft_concurrency_share IS NULL OR soft_concurrency_share > 0)
);
```

`account_id` 作为主键，明确一个账号同一时刻最多属于一个共享容量组。账号仍可出现在多个普通账号集合中，但只有一条容量归属，避免同一请求需要同时占用多个共享组而产生优先级和死锁语义。

字段语义：

- `hard_concurrency_limit`：成员独立硬上限。达到后即使组内有余量也不能继续借用。
- `soft_concurrency_share`：调度软份额，只用于负载评分、借用统计和压力展示，不参与硬拒绝。
- `hard_concurrency_limit = NULL`：成员只受共享总上限约束。只有确认上游允许单 Key 突破原份额后才使用。
- `soft_concurrency_share = NULL`：不计算该成员的“借用”状态，仍可正常参与共享总容量。

启用校验：

1. 容量成员必须已经属于对应 `upstream_account_set_members`。
2. 集合与所有容量成员的 `platform` 必须一致。
3. `hard_concurrency_limit` 不得大于 `shared_concurrency_limit`。
4. 允许软份额之和小于或大于共享上限；它是调度目标，不是配额切片。管理员页面需要提示偏差，但不阻止保存。
5. 删除集合成员时通过复合外键级联删除容量成员；禁用集合时新请求回退到账号独立并发控制，不删除配置。

### 3.3 聪明三账号的安全初始配置

在上游尚未确认单 Key 独立硬上限前：

```text
shared_concurrency_limit = 3000
每成员 hard_concurrency_limit = 1000
每成员 soft_concurrency_share = 1000
```

这一步只建立正确的数据语义，不应直接沿用当前 `1000 / 1200 / 1000` 作为共享模型。若后续确认单 Key 可以借用空闲额度，再单独放宽或清空成员硬上限。

## 4. 运行时身份传递

扩展 `UpstreamPoolResolvedMemberConfig`：

```go
CapacityGroupID        *int64
CapacityGroupLimit     *int
CapacityMemberHardLimit *int
CapacityMemberSoftShare *int
```

运行时查询按 `account_id` 独立 `LEFT JOIN upstream_account_set_capacity_members` 与 `upstream_account_sets`，而不是从最终胜出的直接成员/集合成员分支推断容量身份。这样即使直接成员优先去重，账号仍携带唯一容量组。

账号被克隆并应用池成员覆盖时，将容量信息写入仅用于本次路由的非持久化字段，例如：

```go
type AccountCapacityScope struct {
    GroupID        int64
    GroupLimit     int
    MemberHardLimit int
    MemberSoftShare int
}
```

所有账号槽位入口统一改为接收 `AccountConcurrencyScope`，不能只修改某一个 OpenAI handler。当前 `GatewayService`、`OpenAIGatewayService`、WebSocket 和 handler fast path 都存在账号 acquire 调用，实施阶段必须以接口编译错误和调用点清单保证没有旁路。

容量组必须对成员账号的所有可调度流量生效，而不仅是“从账号集合分支展开”的流量。若某条路径无法获得容量身份，应视为实现未完成，不能以旧账号槽位静默绕过共享硬上限。

## 5. Redis key 与原子 acquire

### 5.1 Key 设计

容量组内两个 ZSET 使用相同 Redis Cluster hash tag：

```text
concurrency:capacity:{set:<setID>}:group
concurrency:capacity:{set:<setID>}:member:<accountID>
```

花括号内的 `{set:<setID>}` 保证两个 key 位于同一 hash slot，可以由一个 Lua 脚本原子读写。普通非容量组账号继续使用现有 `concurrency:account:<accountID>`。

容量组账号不再同时占用旧账号 ZSET；成员 ZSET 就是该账号在容量组内的账号级槽位。否则 Redis Cluster 下无法把不同 hash slot 的旧账号 key 与组 key 放入同一 Lua，也会产生双重释放和部分成功问题。

### 5.2 Acquire 输入与返回值

Lua 输入：

```text
KEYS[1] group ZSET
KEYS[2] member ZSET
ARGV[1] group limit
ARGV[2] member hard limit，0 表示无独立硬上限
ARGV[3] slot TTL seconds
ARGV[4] request ID
ARGV[5] member soft share，0 表示不统计借用
```

返回码：

```text
1 acquired
2 group_full
3 member_full
```

服务层将拒绝原因写入 `AcquireResult.DeniedScope`，等待队列、指标和日志不再把所有失败都归为账号满载。

### 5.3 Lua 原子步骤

单次脚本内完成：

1. 使用 Redis `TIME` 获取服务器时间。
2. 对 group/member 两个 ZSET 清理过期 request ID。
3. 检查 request ID 是否已同时存在；若存在则刷新两个 ZSET 的时间与 TTL，返回 `acquired`。
4. 若只在一个 ZSET 中存在，先删除该残缺成员，再按新请求重新判断。正常 Lua 不会产生部分写入，此分支用于兼容异常数据和人工操作。
5. 若 group `ZCARD >= group limit`，记录 `group_full`，返回 2。
6. 若成员硬上限大于 0 且 member `ZCARD >= member hard limit`，记录 `member_full`，返回 3。
7. 同时向两个 ZSET `ZADD` 同一个 request ID，并刷新 TTL。
8. 若 acquire 前成员并发已达到软份额，记录一次 `borrowed_slot`。
9. 返回 1。

因为判断和两次 `ZADD` 在同一 Lua、同一 hash slot 内完成，成功请求不可能只占账号槽位或只占共享槽位，也不会超卖共享总上限。

### 5.4 原子 release

release 使用同 hash slot 的双 key Lua：

```text
ZREM group requestID
ZREM member requestID
```

脚本幂等；任一 key 中不存在 request ID 都返回成功。释放仍使用独立的 5 秒 background context。

### 5.5 TTL 与进程清理

- 沿用当前槽位 TTL 和 request ID 进程前缀。
- `CleanupStaleProcessSlots` 增加 `concurrency:capacity:*` 扫描。
- 同一 hash tag 下可以批量清理一组容量 key；跨组仍逐组执行，不能把不同组 key 传入同一 Lua。

## 6. 等待队列与压力数据

新增组级等待 key：

```text
wait:capacity:{set:<setID>}:group
```

- `group_full` 进入组等待计数。
- `member_full` 继续计入成员等待，并同时纳入组压力统计，但不伪装成共享总容量已满。
- acquire 重试必须保留同一个容量 scope；配置缓存失效不能让等待中的请求退回旧账号槽位。

建议新增分钟快照表 `upstream_capacity_group_snapshots`：

```text
set_id
bucket_at
current_concurrency
waiting_count
capacity_limit
peak_concurrency
group_full_count
member_full_count
borrowed_slot_count
scheduler_concentration
```

5 分钟峰值和 P95 负载率从分钟快照计算；满载次数与借用次数由 Redis 分钟桶累计后落库。不要仅依赖请求日志反推并发峰值。

管理员精确指标：当前并发/共享容量、等待数、5 分钟峰值、P95 负载率、组满载次数、成员满载次数、借用次数、调度集中度、可用容量余量。

用户 `/monitor` 只返回枚举状态：

```text
ample       余量充足
observe     需要观察
tight       余量紧张
queueing    排队中
```

建议映射优先级：有等待即 `queueing`；否则 P95 或当前负载达到高阈值为 `tight`；达到观察阈值为 `observe`；其余为 `ample`。具体阈值在压力面板实现前单独确认。

## 7. 自动调权三态兼容

新字段位于现有 `policy_json.auto_weight`：

```json
{
  "auto_weight": {
    "mode": "off | observe | active",
    "enabled": false
  }
}
```

读取兼容规则：

1. `mode` 为合法值时以 `mode` 为准。
2. 没有 `mode` 时，旧 `enabled=true` 映射为 `active`。
3. 没有 `mode` 且 `enabled=false` 或缺失时映射为 `off`。

写入规则：

- `active`：写 `mode=active`、`enabled=true`。
- `observe`：写 `mode=observe`、`enabled=false`。
- `off`：写 `mode=off`、`enabled=false`。

`observe` 继续计算并持久化 `target_factor`、原因和观测时间，但调度查询的 `effective_weight` 必须保持手工权重。只有 `active` 才应用运行时 factor。管理员页面同时展示“当前生效权重”和“观察建议权重”，避免把建议误认为已生效。

## 8. 分阶段实施顺序

1. 新迁移、Repository/Service 类型、配置校验和 Redis 双作用域 Lua 单元/集成测试。
2. 统一账号 acquire 接口并接入全部 HTTP、OpenAI、WebSocket、fast path；增加无旁路测试。
3. 组级当前压力与分钟快照，先完成管理员精确面板。
4. 用户 `/monitor` 粗粒度状态。
5. 自动调权 `off / observe / active`，默认先切 `observe` 运行 24～48 小时。
6. 完整后端测试、`pnpm --dir frontend typecheck`、`pnpm --dir frontend build`、`go test -tags embed ./internal/web` 与 `127.0.0.1:18080` 路由验证。

## 9. 架构确认点

进入代码实现前只需确认两项：

1. 是否接受“账号最多属于一个共享容量组”的硬约束。
2. 聪明三个 Key 在首期是否保持每成员 `hard_concurrency_limit=1000`；若上游确认单 Key 可突破 1000，再放宽成员硬上限。

