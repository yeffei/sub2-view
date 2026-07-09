# 山枢庭 / SST 生产环境稳定升级 SOP

日期：2026-07-09

## 适用范围

本文面向已经部署到服务器上的山枢庭 / SST 实例，目标是把升级流程固定成一条可重复、可回滚、不破坏现有界面和业务功能的正式链路。

适用前提：

- 线上运行的是你自己的 SST fork 产物，而不是直接跟随上游默认镜像。
- 管理员后台升级按钮已经接到自有 release 源。
- 数据库、Redis、配置文件、上传数据都独立于镜像之外持久化。

不适用场景：

- 本地开发热更新
- `source` 构建实例直接点后台按钮升级
- 未做版本管理、仍长期使用 `latest` 且无法回滚的临时部署

---

## 一句话原则

正式升级一律走：

`自己的 GitHub Release` + `版本化镜像` + `服务器 compose 切明确 tag`

管理员后台“立即更新”按钮保留，但只作为：

- 单机补丁
- 紧急热修
- release 已发布后的临时抢修手段

不能把按钮升级当成唯一升级方式，因为按钮只更新当前运行实例，不会自动更新你的镜像基线。

---

## 当前仓库已经具备的升级能力

截至 2026-07-09，仓库内已经接好这些能力：

1. 更新源可配置
   - `backend/internal/config/config.go`
   - `deploy/config.example.yaml`
   - 配置键：
     - `update.release_repo`
     - `update.github_token`

2. Docker 本地运行环境可注入更新源
   - `deploy/docker-compose.local.yml`
   - `deploy/.env.example`
   - 环境变量：
     - `UPDATE_RELEASE_REPO`
     - `UPDATE_GITHUB_TOKEN`

3. 管理员后台可显示更新源
   - 版本卡片会展示 `release_repo`

4. release workflow 会产出正式发布物
   - `.github/workflows/release.yml`
   - `.goreleaser.yaml`
   - 产物包括：
     - GitHub Release assets
     - `checksums.txt`
     - GHCR 版本化镜像与 `latest`

5. 本地 Docker 可补镜像基线
   - `deploy/build-local-release.ps1`

---

## 正式升级的标准流程

### Phase 1. 本地准备发布版本

先在自己的 SST 主仓库完成：

1. 代码合入 `main`
2. 本地最少跑固定验证
3. 确认 release 配置仍指向自己的仓库

建议最少验证：

```powershell
pnpm --dir frontend typecheck
pnpm --dir frontend build
cd backend
go test ./...
```

如果这次包含首页、登录、Dashboard、Admin 改动，还要补页面访问验收：

- `/home`
- `/login`
- `/register`
- `/dashboard`
- `/keys`
- `/usage`
- `/profile`
- `/admin/dashboard`
- `/admin/settings`

### Phase 2. 打 tag 并生成 GitHub Release

在仓库根目录执行：

```powershell
git switch main
git pull --ff-only origin main
git tag -a v0.1.147 -m "SST release v0.1.147"
git push origin v0.1.147
```

然后等待 `.github/workflows/release.yml` 完成。

必须确认三件事都存在：

1. GitHub Release 页面已生成
2. 当前平台所需资产已生成
   - 例如 `sub2api_0.1.147_linux_amd64.tar.gz`
3. `checksums.txt` 已生成

如果你线上主要跑 Docker，还要确认 GHCR 镜像已可拉取，例如：

```text
ghcr.io/<你的 GitHub owner 小写>/sub2api:0.1.147
```

注意：

- Git tag 用 `vX.Y.Z`
- GoReleaser 产出的镜像 tag 和二进制版本通常是 `X.Y.Z`

不要只看到 tag 存在就升级服务器，必须等 release 产物和镜像都 ready。

### Phase 3. 服务器切换到明确版本镜像

生产环境不要长期跟 `latest`。  
把服务器上的 `docker-compose.yml` 或部署配置固定成明确版本号：

```yaml
services:
  sub2api:
    image: ghcr.io/<your-owner>/sub2api:0.1.147
```

然后在服务器执行：

```bash
docker compose pull
docker compose up -d
docker compose ps
```

如果有反向代理或额外 sidecar，不要在同一次变更里顺手改其它东西。  
升级窗口内只做“切镜像 tag + 起容器 + 验收”。

### Phase 4. 升级后验收

先看容器和应用状态：

```bash
docker compose ps
docker compose logs --tail=100 sub2api
```

至少验这几类：

1. 页面可达
   - `/dashboard`
   - `/keys`
   - `/usage`
   - `/profile`
   - `/admin/dashboard`
   - `/admin/settings`

2. 登录和 API key 调用链正常

3. 关键业务链正常
   - 支付 / 订阅
   - 风控
   - 上游池 / 调度
   - 你当前线上最核心的 bot 调用链

4. 版本卡片状态正确
   - 当前版本正确
   - 更新源仍指向自己的 `release_repo`

### Phase 5. 留好回滚点

正式升级后，不要立刻删掉上一个可用镜像 tag。  
至少保留上一个稳定版本，例如：

- `0.1.146`
- `0.1.147`

只要 compose 仍是明确 tag，回滚就很简单：

```yaml
services:
  sub2api:
    image: ghcr.io/<your-owner>/sub2api:0.1.146
```

然后执行：

```bash
docker compose pull
docker compose up -d
```

---

## 后台按钮升级在线上应该怎么用

按钮升级现在是可用的，但在线上应该降级为“辅助能力”。

### 适合使用按钮升级的情况

- 已经发布了自有 release，但需要先快速把单台机器拉到新版本
- 某个补丁需要先热修，再补镜像基线
- 你明确知道当前实例是 `release` 构建，且预检通过

### 不适合使用按钮升级的情况

- 你准备做正式生产发版
- 你要升级多台服务器
- 你当前实例是 `source` 构建
- 你还没有对应的 GitHub Release 资产和 `checksums.txt`

### 按钮升级后的正确收口

如果你在线上先点了按钮升级，后续必须补这一步：

1. 确认当前实例已升级成功
2. 立刻补发同版本镜像基线
3. 把服务器 compose 里的镜像 tag 也切到同版本

否则以后只要执行：

```bash
docker compose up -d --force-recreate
```

就可能被旧镜像覆盖回去。

---

## 服务器部署前必须固定的配置

### 1. 固定更新源

服务器实例必须明确指向你自己的 release 仓库，例如：

```yaml
update:
  release_repo: "yeffei/sub2-view"
  github_token: ""
```

如果 release 仓库是私有仓库，或者你担心 GitHub API 限流，可配置 token。

### 2. 运行时敏感配置不入库

以下内容只能放服务器环境变量、密钥管理器或本地 `.env`：

- `UPDATE_GITHUB_TOKEN`
- 数据库密码
- Redis 密码
- 支付密钥
- 任何第三方平台 token

不要把这些值提交到 git。

### 3. 数据与配置必须外挂

至少确保这些内容不跟镜像一起丢：

- 数据库
- Redis
- `deploy/data`
- 配置文件
- 上传或缓存中需要持久化的内容

---

## 推荐的生产环境升级节奏

以后每次发版，统一按这个顺序：

1. 本地或 CI 验证通过
2. 打 `vX.Y.Z` tag
3. 等 GitHub Release 和 GHCR 镜像生成完成
4. 先在预发或单机验证
5. 服务器 compose 切到镜像 tag `X.Y.Z`
6. 做固定验收
7. 观察一段时间
8. 再决定是否让后台按钮继续暴露给运营做补丁升级

这样你的“源码版本”“release 资产”“服务器镜像基线”“后台可升级目标”会始终指向同一版本。

---

## 不建议的做法

### 1. 长期使用 `latest`

问题：

- 回滚困难
- 不知道当前实际部署的是哪一个构建
- 多台服务器容易出现版本漂移

### 2. 只点按钮，不补镜像基线

问题：

- 一旦容器重建，就可能回到旧版本

### 3. 一边升级一边改业务配置

问题：

- 出问题后很难判断是版本问题还是配置问题

### 4. source 构建实例直接依赖按钮升级

问题：

- 这本来就不属于按钮升级的设计边界

---

## 本地和线上两条线怎么分工

### 本地 / 运维演练线

适合：

- 快速验证热修
- 做单机业务态验收
- 补本地镜像基线

当前仓库可用脚本：

```powershell
deploy/build-local-release.ps1
```

这个脚本适合本地 Docker 场景，不等同于正式服务器发版。

### 线上正式发版线

适合：

- 稳定升级
- 多次可重复部署
- 可回滚
- 团队协作和长期运维

正式环境应优先依赖：

- GitHub Release
- GHCR 或你自己的镜像仓库
- 版本化 compose 配置

---

## 最小可执行清单

如果以后你只想记一版最短流程，就按这个做：

1. 合代码到 `main`
2. 跑 `pnpm --dir frontend typecheck`、`pnpm --dir frontend build`、`cd backend && go test ./...`
3. 打 tag：`vX.Y.Z`
4. 等 GitHub Release 和镜像生成完成
5. 服务器把镜像改成 `ghcr.io/<owner>/sub2api:X.Y.Z`
6. 执行 `docker compose pull && docker compose up -d`
7. 验 `/dashboard`、`/keys`、`/usage`、`/profile`、`/admin/dashboard`、`/admin/settings`
8. 出问题就把镜像 tag 切回上一个版本

---

## 相关文档

- [docs/sst-release-button-upgrade-plan.md](./sst-release-button-upgrade-plan.md)
- [docs/sst-remote-docker-upgrade-sop.md](./sst-remote-docker-upgrade-sop.md)
- [docs/sub2api-safe-upgrade-strategy.md](./sub2api-safe-upgrade-strategy.md)
- [deploy/build-local-release.ps1](../deploy/build-local-release.ps1)
- [deploy/docker-compose.local.yml](../deploy/docker-compose.local.yml)
- [.github/workflows/release.yml](../.github/workflows/release.yml)
- [.goreleaser.yaml](../.goreleaser.yaml)
