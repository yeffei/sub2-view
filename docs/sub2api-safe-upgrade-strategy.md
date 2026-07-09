# Sub2API 安全升级策略（山枢庭 / SST）

> 如果你当前关心的是“服务器以后如何稳定升级”，优先看 [docs/sst-production-upgrade-sop.md](./sst-production-upgrade-sop.md)。
> 本文更偏“源码层面的升级边界、模块收口和长期兼容策略”。

## 目标

在继续跟进上游 `Sub2API` 版本时，同时满足三件事：

1. 升级成功，可持续吸收上游修复与新能力。
2. 不破坏当前 `山枢庭 / SST` 的前后台界面与品牌表达。
3. 不破坏已经新增的本地功能，尤其是支付、订阅、后台运维、上游池与品牌化页面。

---

## 当前事实

基于当前仓库状态，现有系统已经不是“轻量改皮肤”的原版 Sub2API，而是三层混合体：

1. 上游基础能力层
   - 认证、路由鉴权、API Key、用量、支付、订单、分组、后台管理。
2. SST 品牌与界面层
   - `/home`
   - `/dashboard`
   - `/admin/*` 多处页面的纸面/夜庭视觉
   - 全局品牌文案、图形、配色、排版
3. 本地新增功能层
   - 上游池 / 故障转移 / 调度解释
   - 订阅与支付管理扩展
   - Admin 表格和工作台体验增强
   - 其他已落地的后台运营能力

这意味着后续不能再用“直接覆盖升级”思路，否则每次都会把 UI 和新增功能重新打散。

---

## 结论

最稳的方案不是“保证一次升级完全不冲突”，而是把仓库改造成：

- `上游核心层可替换`
- `SST 视觉层可覆写`
- `本地功能层可插拔`
- `升级验收链可重复执行`

也就是说，后面升级时只允许大部分冲突集中在少数几个边界文件，而不是散落到整个前后端。

---

## 推荐架构边界

### 1. 上游核心层

这层尽量贴近原版，后续优先直接跟上游同步：

- `backend/internal/service` 中未被本地强改的通用能力
- 通用 API handler
- 通用前端基础组件
- 通用类型定义
- 路由与基础 store

要求：

- 非必要不在这层直接写 SST 视觉逻辑。
- 非必要不把本地业务判断硬塞进通用 handler。
- 如果必须修改，优先通过薄封装或扩展点完成。

### 2. SST 视觉覆写层

这层专门承接品牌化改造，避免污染上游核心：

- `frontend/src/views/*` 中明显属于品牌表达的页面
- `frontend/src/components/layout/*`
- `frontend/src/components/common/*` 里只与视觉相关的覆写
- `frontend/src/assets/brand/*`
- 主题变量、色板、排版、页面壳层

要求：

- 尽量通过 CSS token、layout shell、view-level wrapper 改视觉。
- 少改通用业务逻辑，多改视图拼装方式。
- 能放在页面级样式的，不放进通用数据逻辑。

### 3. 本地功能模块层

这层放你新增的业务能力，尽量模块化，不与上游通用流程纠缠：

- 上游池
- 调度解释
- 异常归因
- 本地支付扩展
- 订阅履约扩展
- 后台运营工作台能力

要求：

- 新增能力优先独立 service / handler / component。
- 对外只通过明确接口接入上游主链路。
- 尽量避免在一个上游核心文件里横向塞很多本地分支。

### 4. 兼容适配层

这是升级成败的关键，专门吸收“上游字段变化 / 返回结构变化 / 配置项变化”：

- 前端 API adapter
- DTO / mapper
- setting 兼容转换
- 页面 props 归一化

典型做法：

- 后端返回结构变化时，不直接让多个页面跟着散改。
- 由 adapter 先统一成 SST 当前使用的稳定格式。

---

## 必须保护的稳定边界

后续任何升级，都要把下面这些当成“不可随意破坏”的兼容面：

### 前台功能边界

- `home_content`
- `featureFlags`
- `simple mode`
- `backend mode`
- 支付开关、风控开关
- `custom_menu_items`
- `/home`
- `/login`
- `/register`
- `/dashboard`
- `/keys`
- `/usage`
- `/profile`
- `/purchase`
- `/orders`

### 后台功能边界

- `/admin/dashboard`
- `/admin/settings`
- `/admin/accounts`
- `/admin/ops`
- `/admin/subscriptions`
- `/admin/payment/*`
- `/admin/channels/upstream-pools`

### 数据与履约边界

- 支付订单状态流转
- 余额充值履约
- 订阅订单履约
- 退款回滚
- 分组与订阅关系
- 上游池绑定与调度解释

---

## 升级时的落地流程

### Phase 1. 先做“边界收口”

先不急着升版本，先把容易炸的改动从上游核心文件里往外收：

- 视觉相关，尽量收口到 view / layout / theme token。
- 本地新增能力，尽量收口到独立模块。
- 返回结构差异，收口到 adapter。

目标：

- 下次升级时，真正高风险的文件数量明显变少。

### Phase 2. 建“升级工作分支”

不要在当前主开发分支上直接升：

1. 建一个专门的 upgrade 分支。
2. 只做上游同步与冲突解决。
3. 不混入新的页面需求或产品需求。

目标：

- 把“升级问题”和“新需求问题”分开。

### Phase 3. 先合上游，再回填本地覆写

推荐顺序：

1. 先同步上游代码。
2. 让项目先恢复到“可编译、可运行”。
3. 再按模块回填 SST 视觉层和本地功能层。
4. 最后才做页面细节修正。

不要反过来做。否则你会在冲突状态里同时处理逻辑和 UI，很容易漏。

### Phase 4. 跑固定验收链

每次升级后必须跑同一套检查：

1. `pnpm --dir frontend typecheck`
2. 前台核心路由可访问
3. 后台核心路由可访问
4. 支付/订阅核心链路烟测
5. 上游池/运维核心链路烟测
6. 关键页面桌面端与移动端视觉回归

### Phase 5. 按钮式 release 自更新边界

后台版本卡片里的“立即更新”只适合 `release` 构建的二进制自更新，不等同于源码仓库升级。

当前按钮链路应遵循：

1. 先调用 `/admin/system/update/preflight` 做预检。
2. 预检通过后，才调用 `/admin/system/update`。
3. 更新成功后保留当前二进制 `.backup`，用于 `/admin/system/rollback`。
4. 更新或回滚完成后，再调用 `/admin/system/restart` 让新二进制生效。

预检至少覆盖：

- 当前构建必须是 `release`，`source` 构建必须走 worktree / git 升级流程。
- 必须有新版本。
- 必须找到当前平台匹配的 release asset。
- 下载 URL 与 checksum URL 必须来自可信 GitHub 域。
- 可执行文件路径必须可解析。
- 可执行文件所在目录必须可写。
- 回滚备份槽位必须可检查。

要求：

- `update`、`rollback`、`restart` 都必须带 `Idempotency-Key`。
- `source` 构建不得通过按钮直接升级。
- 预检失败时，不允许继续执行替换。
- 回滚属于破坏性操作，界面必须二次确认。

### Phase 6. 本地源码同步与自有 release 节奏

这一步是把“上游源码更新”和“线上按钮更新”串成一条稳定链路：

1. 本地源码先通过 `worktree / upgrade` 分支同步上游。
2. 同步后在本地或 CI 跑固定验收链。
3. 验收通过后合回自己的 SST 主分支。
4. 从自己的 SST 主分支打 `vX.Y.Z` tag，触发 GitHub Release。
5. 线上 release 部署再通过版本卡片按钮更新到这个 tag 对应的二进制。

#### 0. 当前仓库前置条件

当前 `D:\sub2api` 目录里的 `.git` 是空目录，不能直接执行：

```bash
git status
git worktree add ...
```

所以第一次必须先把本地源码恢复成一个正规 Git 仓库。推荐用“新目录重建”方式，避免误伤当前工作目录：

```powershell
cd D:\
git clone <你的 SST fork 仓库地址> sub2api-repo
cd D:\sub2api-repo
git remote add upstream https://github.com/Wei-Shaw/sub2api.git
git fetch origin --tags
git fetch upstream --tags
```

如果还没有自己的 SST fork 仓库，先在 GitHub 创建一个私有或公开仓库作为 `origin`。不要把线上实例直接指向上游原仓库，否则你自己的 SST 改造没有稳定发布源。

如果当前 `D:\sub2api` 是唯一保存了 SST 改造的目录，先不要删除或覆盖它。应当把它作为“当前源码快照”导入到 `sub2api-repo` 的一个保护分支：

```powershell
cd D:\sub2api-repo
git switch -c import/current-sst
robocopy D:\sub2api D:\sub2api-repo /E /XD .git node_modules dist .vite .tmp output /XF .env
git status --short
git add .
git commit -m "chore: import current SST source snapshot"
```

确认 `import/current-sst` 能编译、能跑关键验证后，再把它整理成你的 SST 主分支。之后所有源码升级都在这个正规 Git 仓库里做。

#### 1. 每次同步上游前

先保证自己的主分支是干净的：

```powershell
cd D:\sub2api-repo
git switch main
git status --short
git pull --ff-only origin main
git fetch upstream --tags
```

要求：

- `git status --short` 必须为空，未提交修改先提交或暂存到单独分支。
- `origin` 指自己的 SST 仓库。
- `upstream` 指原版 `https://github.com/Wei-Shaw/sub2api.git`。

#### 2. 创建独立 upgrade worktree

每次升级用当天日期建一个独立目录和分支：

```powershell
cd D:\sub2api-repo
$date = Get-Date -Format 'yyyyMMdd'
git worktree add "..\sub2api-upgrade-$date" -b "upgrade/upstream-$date" main
cd "..\sub2api-upgrade-$date"
```

然后在这个 worktree 合上游：

```powershell
git merge --no-ff upstream/main
```

如果有明确的上游 tag，也可以合指定 tag：

```powershell
git merge --no-ff v1.2.3
```

要求：

- upgrade 分支只做上游同步、冲突解决、兼容修复。
- 不在 upgrade 分支混入新页面、新产品需求、新视觉方向。
- 冲突文件优先按“上游核心层 / SST 视觉覆写层 / 本地功能模块层 / 兼容适配层”分类处理。

#### 3. 冲突处理顺序

出现冲突时先列清单：

```powershell
git diff --name-only --diff-filter=U
```

建议按这个顺序解决：

1. `go.mod`、`go.sum`、`package.json`、`pnpm-lock.yaml`
2. 后端通用 service / handler / route
3. 本地新增模块，例如 payment、subscription、upstream pool、ops
4. 前端 API adapter、store、router
5. 页面和视觉覆写
6. 文档、部署配置、release 配置

原则：

- 上游安全修复和协议兼容修复优先保留。
- SST 品牌页面、`home_content`、`custom_menu_items`、支付/订阅/上游池链路必须保留。
- 上游返回结构变化优先在 adapter 层吸收，不让多个页面分散改。

#### 4. 固定验收链

前端至少跑：

```powershell
pnpm --dir frontend install --frozen-lockfile
pnpm --dir frontend typecheck
pnpm --dir frontend exec vitest run src/api/__tests__/admin.system.spec.ts src/components/common/__tests__/VersionBadge.spec.ts
make test-frontend-critical
```

后端至少跑：

```powershell
cd backend
gofmt -w ./cmd ./internal
go test ./...
```

页面和链路至少验收：

- `/home`
- `/login`
- `/register`
- `/dashboard`
- `/keys`
- `/usage`
- `/profile`
- `/purchase`
- `/orders`
- `/admin/dashboard`
- `/admin/settings`
- `/admin/accounts`
- `/admin/ops`
- `/admin/subscriptions`
- `/admin/payment/plans`
- `/admin/payment/orders`
- `/admin/channels/upstream-pools`

业务烟测至少覆盖：

- 余额充值下单
- 订单状态更新
- 订阅订单履约
- 退款后余额或订阅回滚
- 上游池配置读取
- 调度解释展示
- 版本卡片预检接口返回正常

#### 5. 合回 SST 主分支

验收通过后，把 upgrade 分支合回自己的 SST 主分支：

```powershell
cd D:\sub2api-repo
git switch main
git merge --no-ff upgrade/upstream-YYYYMMDD
git status --short
git push origin main
```

如果 upgrade worktree 不再需要，移除它：

```powershell
git worktree remove "..\sub2api-upgrade-YYYYMMDD"
git branch -d upgrade/upstream-YYYYMMDD
```

#### 6. 打自己的 SST release

更新按钮的版本比较只按前三段数字版本比较，所以自有 release 建议使用单调递增的 `vX.Y.Z`：

```powershell
git switch main
git pull --ff-only origin main
git tag -a v1.2.4 -m "SST release v1.2.4"
git push origin v1.2.4
```

GitHub Actions 的 `Release` workflow 会：

- 写入 `backend/cmd/server/VERSION`
- 构建前端并嵌入后端
- 通过 GoReleaser 构建 `BuildType=release` 的二进制
- 生成 GitHub Release assets 和 `checksums.txt`

只有这个 release 构建出来的线上实例，才应该使用后台版本卡片按钮自更新。

#### 7. 线上按钮更新

线上环境更新顺序：

1. 在 GitHub Release 确认目标 `vX.Y.Z` 已发布，且当前平台 asset 和 `checksums.txt` 存在。
2. 打开线上后台版本卡片。
3. 点击“立即更新”。
4. 系统先跑 `/admin/system/update/preflight`。
5. 预检通过后执行 `/admin/system/update`。
6. 更新成功后点击重启。
7. 如更新后异常，使用版本卡片里的回滚按钮回到 `.backup` 二进制。

如果线上实例显示 `source` 构建，按钮会拒绝更新。这是正确行为，说明该实例不是 release 二进制部署，应该先改成 release 部署方式。

---

## 你这个仓库最该避免的做法

### 1. 直接在通用组件里混入大量品牌判断

例如把很多 SST 页面样式直接写死进全局 `DataTable`、通用业务组件、通用 store。

问题：

- 上游一改这些组件，冲突会非常大。
- 后续你不知道某个视觉问题到底是品牌覆写还是通用组件回归。

### 2. 在上游核心 handler/service 里散布本地业务分支

问题：

- 每次升级都要手动重新判断这些分支还该不该保留。
- 很容易把上游新修复覆盖掉，或者把本地逻辑漏掉。

### 3. 升级和新需求同时做

问题：

- 一旦页面坏了，不知道是上游升级造成的，还是新需求造成的。

---

## 对当前仓库的具体建议

结合当前代码，建议下一步按下面顺序做。

### A. 先把“品牌层”和“业务层”再切干净一点

重点看这些区域：

- `frontend/src/views/*`
- `frontend/src/components/layout/*`
- `frontend/src/components/common/DataTable.vue`
- `frontend/src/views/admin/*`
- `frontend/src/views/user/PaymentView.vue`
- `frontend/src/components/payment/*`

判断原则：

- 如果改的是视觉和排版，尽量留在 view 或 layout。
- 如果改的是支付/订阅行为，尽量留在 payment 模块。
- 如果只是兼容上游返回结构，做 adapter，不要散改多个页面。

### B. 订阅购买链单独收口

当前仓库已经有：

- `subscription plan`
- `order_type = subscription`
- 订阅履约 `AssignOrExtendSubscription`

但用户侧 `/purchase` 现状主要还是：

- 余额充值
- 外链购买
- 兑换码

这块后续如果要做“正式升级后的订阅购买功能”，建议单独收成一条独立链，不要继续混在余额充值主视图里。

建议拆成：

- 余额充值链
- 订阅套餐链
- 外部购买链
- 兑换码链

这样以后上游即使改支付页，也不容易把四条逻辑彼此带崩。

### C. 给关键新增能力补最小回归测试

重点不是全量测试，而是给最容易被升级打坏的点补“存在性测试”：

- 订阅订单创建
- 订阅履约成功
- 退款回滚订阅
- `/purchase` 页面能正确读取 checkout info
- `home_content` 仍可接管首页
- `custom_menu_items` 仍能注入导航
- 上游池绑定解析仍生效

---

## 建议的升级验收清单

### 编译与类型

- `pnpm --dir frontend typecheck`
- 前端关键测试至少跑支付/订阅相关用例
- 后端关键测试至少跑 payment / subscription / upstream pool 相关用例

### 前台页面

- `/home`
- `/login`
- `/register`
- `/dashboard`
- `/keys`
- `/usage`
- `/profile`
- `/purchase`
- `/orders`

### 后台页面

- `/admin/dashboard`
- `/admin/settings`
- `/admin/accounts`
- `/admin/ops`
- `/admin/subscriptions`
- `/admin/payment/plans`
- `/admin/payment/orders`
- `/admin/channels/upstream-pools`

### 关键链路

- 余额充值下单
- 订单状态更新
- 订阅订单履约
- 退款后余额或订阅回滚
- 上游池配置读取
- 调度解释展示

---

## 推荐执行方式

如果目标是“以后能持续跟上游升级”，推荐按两步走：

1. 先做一次“升级前收口”
   - 不升版本，只整理边界。
2. 再做一次“正式升级演练”
   - 用独立分支同步上游并跑完整验收。

这比现在直接硬升一个新版本更稳。

---

## 最终建议

对于 `D:\sub2api` 这条线，我建议后续采用下面的原则：

- 上游版本同步是常规动作，不再和产品需求混做。
- SST 界面改造尽量视图层化，不侵入通用核心。
- 本地新增能力尽量模块化，不散落在上游基础链路里。
- 每次升级都跑固定验收清单，而不是靠肉眼临时补查。

如果按这个方向走，后面要保证“升级成功但不破坏现有界面和新功能”，可行，而且成本会越来越低。
