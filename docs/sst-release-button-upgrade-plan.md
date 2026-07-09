# 山枢庭 / SST 按钮升级切换到自有 Release 的最小方案

日期：2026-07-08

> 2026-07-09 更新：
> 本文中的最小改造已经完成。当前仓库已支持 `update.release_repo`、`update.github_token`，管理员后台版本卡片也会展示更新源。
> 如果你关心服务器长期稳定升级，请直接看 [docs/sst-production-upgrade-sop.md](./sst-production-upgrade-sop.md)。

## 结论

当前管理员后台左上角版本卡片的“立即更新”按钮，技术上已经可用。本文保留的是“当时如何把它从上游官方 release 切到 SST 自有 release”的设计收口。

所以要分开理解：

1. `现在能不能点`
可以点，前提是当前运行实例是 `release` 构建，且预检通过。

2. `现在该不该把它当正式升级主流程`
不建议。正式生产升级仍应以 `GitHub Release + 版本化镜像 + compose 切 tag` 为主。

## 当前实例实测状态

在当前本地业务态实例 `http://127.0.0.1:18080` 上，`/api/v1/admin/system/update/preflight?force=true` 返回：

- `build_type = release`
- `current_version = 0.1.137`
- `latest_version = 0.1.146`
- `has_update = true`
- `can_update = true`

这说明：

- 当前容器实例符合“按钮式二进制自更新”的运行前提。
- 按钮升级的目标取决于当前实例配置的 `release_repo`，不再固定指向上游仓库。

## 当前按钮链路现状

### 1. 前端按钮行为是正确的

当前前端不是盲点升级，而是：

1. 先调用 `/admin/system/update/preflight`
2. 预检通过后再调用 `/admin/system/update`
3. 更新完成后提示 `/admin/system/restart`
4. 支持 `/admin/system/rollback`

这条交互边界是对的，前端无需推翻重做。

### 2. 后端也确实执行了预检门

当前后端 `PerformUpdate()` 会再次构建 preflight 结果，并在 `CanUpdate=false` 时拒绝执行替换。

这意味着按钮链路本身是安全收口的，问题不在“有没有门”，而在“门后面指向哪个 release 源”。

### 3. 旧问题与当前状态

本文最初写作时，后端更新服务里存在上游仓库硬编码，因此按钮升级只能跟踪 `Wei-Shaw/sub2api`。

截至 2026-07-09，这个问题已经完成修复：

- 检查更新改为读取 `update.release_repo`
- GitHub API 请求可选携带 `update.github_token`
- `UpdateInfo` / `UpdatePreflightInfo` 会返回 `release_repo`
- 前端版本卡片会展示“更新源”

也就是说，这份文档里的“最小方案”已经不是待做，而是已落地的设计依据。

## 为什么这对 SST 不够

如果山枢庭 / SST 后续继续保留：

- 前台品牌化表达
- 导航与信息组织差异
- 后台管理改造
- 兼容层与定制链路

那么“升级按钮”应该升级到“你自己审核并发布过的 SST release”，而不是直接升级到上游官方 release。

否则会有三个长期问题：

1. 按钮升级结果不可控
上游新版本可能包含与你本地定制尚未合流的变化。

2. 线上版本来源不清晰
管理员在后台看到的是“可升级”，但实际升级目标不是你自己的发布节奏。

3. 你做的源码验收与线上按钮更新脱节
本地走的是 `origin/upstream/worktree` 合流逻辑，线上按钮走的却是另一条源。

## 最小改造目标

不改动“按钮升级的交互模型”，只把“升级源”从上游官方 release 改成可配置。

也就是说，保留这些现有能力：

- `release` 构建才允许按钮升级
- `source` 构建继续禁止按钮升级
- 必须先跑 preflight
- 必须支持 rollback / restart
- 必须保留 GitHub 下载域名白名单和 checksum 校验

只新增一个核心能力：

- 更新服务能够读取“当前实例应该跟踪哪个 GitHub release 仓库”

## 最小改造方案

### Phase 1. 把 release 仓库从硬编码改成配置项

在 `UpdateConfig` 里新增一个配置，例如：

```yaml
update:
  proxy_url: ""
  release_repo: "yeffei/sub2-view"
```

建议规则：

- 缺省值保持 `Wei-Shaw/sub2api`
- SST 环境显式配置成自己的 fork，例如 `yeffei/sub2-view`

这样改的好处是：

- 不破坏默认行为
- 上游用户仍保持原逻辑
- SST 环境只通过配置切换发布源

### Phase 2. UpdateService 持有 `releaseRepo`

当前 `UpdateService` 只持有：

- `currentVersion`
- `buildType`

建议新增：

- `releaseRepo`

然后把所有 `FetchLatestRelease(ctx, githubRepo)` 改成读取实例字段，而不是常量。

这样 `check-updates`、`preflight`、`performUpdate` 会自然统一到同一个发布源。

### Phase 3. 把 release 源暴露给前端

建议在这两个响应里新增一个只读字段：

- `UpdateInfo.release_repo`
- `UpdatePreflightInfo.release_repo`

用途不是功能必需，而是降低运营歧义：

- 管理员能看到当前检查的是哪个仓库
- 避免把 SST 环境误认成在跟官方主仓

### Phase 4. 前端版本卡片补一句来源提示

版本卡片里建议增加一句小字，例如：

- `更新源：yeffei/sub2-view`

如果是默认上游源，也可展示：

- `更新源：Wei-Shaw/sub2api`

这一步不是必须，但很值，能让“这个按钮到底升到哪”变成界面上的明确信息，而不是隐藏规则。

## 不需要改的地方

以下机制当前可以保留：

1. GitHub 下载域名校验
当前允许域名是 GitHub 官方下载域，仓库换成你的 fork 后仍然适用。

2. `release` / `source` 构建门
这是正确的边界，不应放宽。

3. preflight 结构
当前已经覆盖：
- build type
- has update
- platform asset
- checksum
- executable path
- writable directory
- backup slot

这套门继续沿用即可。

4. rollback / restart 交互
当前按钮链路已成型，不需要重做。

## 推荐的 SST 线上发布节奏

按钮升级切到自有 release 后，推荐形成固定节奏：

1. 本地 `sub2api-repo` 从 `upstream/main` 合流
2. 在 upgrade worktree 解决冲突并跑固定验收
3. 合回自己的 `main`
4. 从自己的仓库打 tag
5. 由自己的 GitHub Release 产出 Linux `amd64` 资产与 `checksums.txt`
6. 线上按钮升级只跟踪你自己的 `release_repo`

这样“源码升级”和“线上按钮升级”就会变成同一条发布链，而不是两条彼此独立的链路。

## 实施优先级

如果只做最小闭环，优先级建议如下：

1. 后端：`update.release_repo` 配置化
2. 后端：`UpdateInfo / UpdatePreflightInfo` 返回 `release_repo`
3. 前端：版本卡片显示“更新源”

做到这三步，按钮升级就已经能从“官方 release 按钮”变成“SST 自有 release 按钮”。

## 最后判断

当前结论可以固定成一句话：

> 现在管理员后台的升级按钮已经可以用，但它当前升级的是上游官方 release；如果要把它变成山枢庭 / SST 的正式升级入口，最小改造应该是把更新仓库改成可配置，并让线上实例指向你自己的 fork release。
