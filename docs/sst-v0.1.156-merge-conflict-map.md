# SST v0.1.156 合流冲突映射

## 范围

- 基线：`main` 的 SST 定制与进行中的性能相关提交。
- 上游：官方标签 `v0.1.156`。
- 工作树：`D:\sub2api-upgrade-v0.1.156-merge`。
- 本记录仅定义解决顺序和不可破坏边界；不代表该合并已完成或可部署。

## 当前冲突清单

`git diff --name-only --diff-filter=U` 当前返回 108 个未解决文件：42 个后端、61 个前端、5 个交付/嵌入相关文件。冲突大多发生在测试文件；所有配套生产文件及新增上游文件也必须在相同模块内一并审查，不能仅解决测试标记。

上游标签的 `backend/cmd/server/VERSION` 内容为 `0.1.155`，但标签与官方发布资产均为 `v0.1.156`。SST 合流产物使用 `0.1.156` 作为运行时版本标识，并保留该差异以便发布审计。

| 模块 | 首要责任 | 合流原则 | 解决顺序 |
| --- | --- | --- | --- |
| 迁移、Ent、数据模型 | `backend/migrations`、`backend/ent`、repository | 已发布 SST 迁移不可改名、覆盖或改 checksum；追加官方迁移并审计所有实体字段和索引 | 1 |
| 网关与协议 | OpenAI/Responses/Anthropic/WebSocket handler、service、apicompat | 以上游协议、安全、取消、超时和 failover 修复为基础；通过上游接口移植 SST 调度钩子 | 2 |
| 调度与上游池 | account、scheduler、concurrency、rate limit、token refresh | 保留 SST pool、account set、shared capacity、TTFT 保护和运行时权重；适配新 account/credential 生命周期 | 3 |
| 认证、OAuth、Grok | OpenAI OAuth、Grok import/probe、account UI | 保留 SST 既有认证和账号操作，同时接入官方 Codex identity 与 Grok 能力 | 4 |
| 支付、订阅与用量 | payment、subscription、usage billing、orders | 保留 SST 计费、支付、订单和订阅语义；不因类型重构丢失历史订单兼容 | 5 |
| 前端数据、类型、i18n、DataTable | API client、stores、types、locales、admin data surfaces | 先收敛 API 契约和类型，再解决 DataTable/筛选/用量视图；不回退 feature flags、simple/backend mode 或 custom menu | 6 |
| SST 视觉层与导航 | Home、header/sidebar、dashboard quick actions、account/keys UI | 在前端数据层稳定后保留 SST 品牌、`home_content` 接管逻辑和深链路由 | 7 |
| 交付和嵌入 | Dockerfile、version、wire、embedded web、deploy docs | 最后处理；只在测试和构建通过后更新版本化镜像流程 | 8 |

## 迁移保护规则

迁移执行器对完整文件名进行字典序排序，并以完整文件名作为 `schema_migrations.filename` 主键。因此同一数字前缀不是可接受的“覆盖”理由。当前树中以下 SST 迁移与官方迁移共享数字前缀，必须保留为独立、不可变文件并逐一审阅：

| SST 文件 | 同前缀官方文件 |
| --- | --- |
| `154_upstream_pools.sql` | `154_account_spark_shadow.sql`、`154_add_ops_system_logs_api_key_id.sql` |
| `156_upstream_account_sets.sql` | `156_content_moderation_matched_keyword.sql` |
| `158_pool_availability_snapshots.sql` | `158_add_group_peak_rate_multiplier.sql`、`158_enable_grok_media_generation_groups.sql` |
| `159_upstream_pool_runtime_weights.sql` | `159_batch_image_foundation.sql` |
| `160_upstream_account_set_shared_capacity.sql` | `160_add_user_frozen_balance.sql`、`160_batch_image_provider_refs.sql` |
| `161_upstream_capacity_snapshots.sql` | `161_batch_image_pricing_snapshot.sql` |

规则：不修改已部署迁移；新增变更必须使用新的、唯一的文件名；在连接真实数据库前先运行迁移回归测试和 schema diff 审计。此阶段不执行任何数据库迁移。

## 解决与验证门槛

每个模块完成后才进入下一个模块：清除其冲突、格式化、运行模块相关测试，再运行完整门槛：

```text
go test ./...
go vet ./...
go test -tags=embed ./internal/web
pnpm --dir frontend typecheck
```

前端构建通过后，再验证 `/home`、`/login`、`/register`、`/dashboard`、`/keys`、`/usage`、`/profile`、`/admin/dashboard`、`/admin/settings` 的桌面与移动端访问。仅当合并无冲突、迁移审计通过、镜像构建和运行检查成功后，才创建备份并切换本地 `sub2api` 容器。
