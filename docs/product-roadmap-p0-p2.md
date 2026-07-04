# 山枢庭 SST Roadmap（2026-06-24）

## 当前结论

本轮最初规划的核心能力，已经有多项提前落地。后续继续沿本文件推进时，优先原则已经从“按表开发”切换为“避免重复实现，优先回写状态，再选择真正未完成的新项”。

当前确认：
- `P0` 用户接入体检单：已完成
- `P0` 错误复盘时间线：已完成
- `P1` 管理员异常归因工作台：已完成
- `P1` 用户侧余额续航预测：已完成
- `P1` 管理员高风险用户榜单：已完成
- `P2` 用户自助恢复向导：已完成
- `P2` 异常模式对比视图：已完成
- `P2` 品牌化周报 / 晨报：已完成首版，并已增强为支持复制与 Markdown 导出的案头简报

## 当前目标

当前更适合的目标不是重复补做上表已完成项，而是：

1. 只选择真正未完成的新能力继续推进
2. 或做小范围工程收口与文档回写，降低误判和重复劳动

原则保持不变：
- 先做首版可用版，优先复用现有前后端数据。
- 不为了首版体验强行新增复杂后端协议。
- 优先在已有页面增量实现，避免打断既有导航与深链。

## P0 / P1 / P2 优先级表

| 优先级 | 功能 | 当前状态 | 目标 | 落点 | 复用数据 / 能力 | 首版验收标准 |
| --- | --- | --- | --- | --- | --- | --- |
| P0 | 用户接入体检单 | 已完成 | 让用户在进入 `/dashboard` 后 10 秒内判断“能否直接放量、该先处理什么” | `frontend/src/views/user/DashboardView.vue` | `usage/dashboard/stats`、最近调用、平台配额、账户余额、现有跳转入口 | 首页出现体检单；至少覆盖密钥、起流、稳定性、额度 4 类检查；每项有状态和动作入口 |
| P0 | 错误复盘时间线 | 已完成 | 让用户在错误详情里看到“这次失败前后发生了什么”，区分单点失败、连续失败、恢复后失败 | `frontend/src/components/user/UserErrorDetailModal.vue` | `usage/errors/:id`、`/usage`、`/usage/errors` | 错误详情出现时间线；至少展示当前错误、相邻成功/失败记录、自动归纳结论 |
| P1 | 管理员异常归因工作台 | 已完成 | 让管理员在总览页快速判断异常主要来自用户侧、上游侧还是平台侧 | `frontend/src/views/admin/DashboardView.vue` | `admin/ops/dashboard/overview`、`admin/ops/errors`、现有日期范围 | 后台首页新增归因区；至少给出归属、来源、集中平台、处理建议 |
| P1 | 用户侧余额续航预测 | 已完成 | 把“余额不足”从静态提醒升级为“还能跑多久”的主动预警 | `/profile`、`/dashboard` | 余额、近 7~14 日消费速度 | 能展示预计续航天数，并在低续航时给出充值或降速建议 |
| P1 | 管理员高风险用户榜单 | 已完成 | 从“消费排行”升级到“异常 / 高支持成本用户榜单” | `/admin/dashboard` 或 `/admin/usage` | 用户错误率、失败量、充值/余额、请求密度 | 能看到高风险用户及主要风险标签 |
| P2 | 用户自助恢复向导 | 已完成 | 把错误分类映射成更完整的恢复脚本与复检动作 | `/usage`、错误详情 | 现有分类解释、下一步动作 | 至少覆盖常见错误的 2~3 步恢复路径 |
| P2 | 异常模式对比视图 | 已完成 | 让管理员比较不同时间窗或平台的异常结构变化 | `/admin/ops` | 错误趋势、分布、平台维度统计 | 能完成两个时间窗的核心异常占比对比 |
| P2 | 品牌化周报 / 晨报 | 已完成首版增强 | 将系统可靠性、用户增长、错误波动转成可读简报 | 管理后台导出或展示 | dashboard / ops 汇总数据 | 生成一页式摘要，适合运营或管理查看 |

## 已完成落点

- `P0` 用户接入体检单
  - `frontend/src/views/user/DashboardView.vue`
- `P0` 错误复盘时间线
  - `frontend/src/components/user/UserErrorDetailModal.vue`
- `P1` 管理员异常归因工作台
  - `frontend/src/views/admin/DashboardView.vue`
- `P1` 用户侧余额续航预测
  - `frontend/src/views/user/DashboardView.vue`
  - `frontend/src/views/user/ProfileView.vue`
- `P1` 管理员高风险用户榜单
  - `frontend/src/views/admin/DashboardView.vue`
- `P2` 用户自助恢复向导
  - `frontend/src/components/user/UserErrorDetailModal.vue`
- `P2` 异常模式对比视图
  - `frontend/src/views/admin/ops/OpsDashboard.vue`
  - `frontend/src/views/admin/ops/components/OpsAnomalyCompareCard.vue`

## 后台简报当前状态

`/admin/dashboard` 当前已经形成三层摘要能力：
- 顶部 `briefing` 晨报摘要
- 中部 `morningSheet` 一页式案头简报
- 下方 `attribution` 与 `risk` 两个执行型工作台

其中 `P2` 品牌化周报 / 晨报 当前已完成的首版增强包括：
- 新增一页式案头简报卡 `admin-morning-sheet-card`
- 支持 `复制简报`
- 支持 `导出 Markdown`
- 导出文件名为 `sst-admin-brief-<start>_<end>.md`
- Markdown 导出已补“异常归因小结”和“高风险用户明细”小节
- 复制与导出行为已有最小测试覆盖

当前判断：
- 该项不应再按“从零做晨报”重复开发
- 如后续继续增强，优先沿 `.md` 结构、正式周报版式或 `.docx` 导出继续推进

## 暂不应重复开发的项

后续若继续按本 Roadmap 推进，以下能力默认视为已完成，不再重复实现：
- `管理员高风险用户榜单`
- `用户自助恢复向导`
- `异常模式对比视图`
- `品牌化晨报 / 一页式案头简报` 的首版复制与 Markdown 导出能力

## 下一步建议

如果继续做产品能力，优先挑“真正未完成的新项”，而不是重复上表已完成内容。

如果短期没有新的产品项，更适合做这两类工作：
- 文档回写与状态收口
- 小范围体验打磨或验证回看

## 验证要求

- 前端至少做定向代码级回读。
- 尽量补最小可跑测试，优先覆盖新增交互或导出行为。
- `pnpm --dir frontend typecheck` 当前已经确认可正常退出，不是挂死；只是 `vue-tsc --noEmit` 默认静默且较慢。
- 当前环境实测：`pnpm --dir frontend typecheck` 大约需要 40 秒以上才结束。
- 排查时可使用 `pnpm --dir frontend run typecheck:diag` 查看 `vue-tsc --diagnostics` 输出，不要再把静默等待直接判定为失败。
