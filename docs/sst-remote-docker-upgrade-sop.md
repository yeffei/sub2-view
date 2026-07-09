# 山枢庭 / SST 远程登录服务器执行 Docker 升级 SOP

日期：2026-07-09

## 适用场景

这份 SOP 适用于下面这种方式：

- 你在本机
- 由你本人或 Codex 从本机 `SSH` 登录服务器
- 在服务器上执行 `Docker Compose` 升级

这不是“后台按钮升级”，也不是“手工覆盖二进制”。  
它本质上仍然是正式的 Docker 版本化升级，只是执行入口从“人在服务器本机上操作”变成“从本机远程操作”。

---

## 一句话结论

以后如果你让我从本机登录服务器升级，推荐固定成这条链：

1. 本地发 `vX.Y.Z` release
2. 等镜像发布完成
3. 从本机 `SSH` 登录服务器
4. 在服务器上把 compose 镜像切到 `X.Y.Z`
5. 执行 `docker compose pull && docker compose up -d`
6. 做验收
7. 有问题立即切回旧 tag

这是可行的，而且比“后台点按钮升级”更稳。

---

## 这条方式为什么稳

因为它同时满足：

- 升级目标明确：升级到哪个镜像 tag 是确定的
- 回滚简单：切回上一个镜像 tag 即可
- 过程可审计：执行过哪些命令都能复盘
- 不会丢镜像基线：容器重建后仍然是新版本

它和“按钮升级”的最大区别是：

- 按钮升级只改当前运行实例
- 远程 Docker 升级同时改运行实例和镜像基线

---

## 前置条件

执行前必须满足这些条件：

### 1. release 已经生成完成

至少确认：

- Git tag 已推送，例如 `v0.1.147`
- GitHub Release 已生成
- 当前平台 release 资产已生成
- GHCR 镜像已生成，例如：
  - `ghcr.io/<owner>/sub2api:0.1.147`

### 2. 服务器当前就是 Docker 部署

你要能在服务器上看到：

- `docker compose.yml` 或等效 compose 文件
- 当前容器名，例如 `sub2api`
- 当前部署目录

### 3. 数据和配置独立于镜像

至少包括：

- PostgreSQL
- Redis
- `.env`
- 持久化数据卷

否则“切镜像 tag”虽然方便，但不能真正做到安全升级。

### 4. 本机具备远程登录条件

至少满足其中一种：

- 本机已有 SSH key，并已加入服务器
- 本机可用密码登录服务器
- 本机已配置 `~/.ssh/config`

如果以后要让我代你执行，这一步尤其重要，因为我只能使用你本机当前已经具备的登录条件。

---

## 推荐目录约定

建议服务器上的部署目录固定，例如：

```bash
/opt/sub2api-deploy
```

建议这里面至少有：

- `docker-compose.yml`
- `.env`
- 如有需要的反向代理配置

不要每次升级都临时去猜 compose 文件在哪。

---

## 标准升级流程

### Phase 1. 本机确认目标版本

先在本机确认本次要升到哪个版本。

例如：

- Git tag：`v0.1.147`
- Docker 镜像 tag：`0.1.147`

注意两者不要混淆：

- Git tag 通常带 `v`
- Docker 镜像 tag 通常不带 `v`

### Phase 2. 从本机 SSH 登录服务器

示例：

```bash
ssh root@your-server-ip
```

或者：

```bash
ssh your-user@your-server-ip
```

如果服务器不是 root 直登，后续命令按需加 `sudo`。

### Phase 3. 进入部署目录并做升级前检查

登录后先做最小检查：

```bash
cd /opt/sub2api-deploy
pwd
docker compose ps
docker compose images
```

还要确认当前 compose 使用的镜像 tag：

```bash
grep -n "image:" docker-compose.yml
```

如果项目使用了多个 compose 文件，比如：

- `docker-compose.yml`
- `docker-compose.prod.yml`

那就必须固定以后到底操作哪一套，不要每次临时判断。

### Phase 4. 备份当前部署文件

正式升级前，建议先备份 compose 文件和 `.env`：

```bash
cp docker-compose.yml docker-compose.yml.bak-$(date +%Y%m%d%H%M%S)
cp .env .env.bak-$(date +%Y%m%d%H%M%S)
```

如果你们把 compose 放在 git 仓库里，这一步也可以用提交记录代替，但线上文件级备份仍然建议保留。

### Phase 5. 修改镜像 tag

把 compose 中的镜像改成目标版本，例如：

```yaml
services:
  sub2api:
    image: ghcr.io/<owner>/sub2api:0.1.147
```

修改后再确认一次：

```bash
grep -n "image:" docker-compose.yml
```

原则：

- 不要继续用 `latest`
- 一次升级只切一个明确版本
- 不要在同一轮里顺手改业务配置

### Phase 6. 拉取新镜像并重建容器

执行：

```bash
docker compose pull
docker compose up -d
docker compose ps
```

如果你希望强制按新镜像重建，也可以执行：

```bash
docker compose up -d --force-recreate
```

但前提是 compose 里的镜像 tag 已经改对。

### Phase 7. 看日志和健康状态

至少检查：

```bash
docker compose logs --tail=100 sub2api
docker inspect sub2api --format '{{.State.Status}} {{if .State.Health}}{{.State.Health.Status}}{{else}}no-healthcheck{{end}}'
```

如果容器名不是 `sub2api`，替换成实际名称。

### Phase 8. 升级后验收

至少做这几类验收：

1. 页面访问
   - `/dashboard`
   - `/keys`
   - `/usage`
   - `/profile`
   - `/admin/dashboard`
   - `/admin/settings`

2. 登录与 API key 调用链

3. 关键业务链
   - 支付 / 订阅
   - 风控
   - 上游池 / 调度
   - 最核心的 bot 调用链

4. 版本确认
   - 容器内应用版本正确
   - 后台版本卡片显示正确
   - `release_repo` 仍指向自己的仓库

---

## 标准回滚流程

如果升级后发现异常，不要在服务器上临时修很多东西。  
最稳的回滚方式是：直接切回上一个镜像 tag。

### 回滚步骤

1. 把 compose 里的镜像改回旧版本，例如：

```yaml
services:
  sub2api:
    image: ghcr.io/<owner>/sub2api:0.1.146
```

2. 执行：

```bash
docker compose pull
docker compose up -d
docker compose ps
docker compose logs --tail=100 sub2api
```

3. 重新验 `/dashboard`、`/keys`、`/usage`、`/profile`、`/admin/dashboard`、`/admin/settings`

如果回滚后恢复正常，说明问题在新版本，不要继续在线上现场改代码。

---

## 如果让我来远程执行，建议怎么配合

以后如果你要让我从本机直接登录服务器操作，最稳的配合方式是：

### 你提前准备

- 本机 SSH 已可登录目标服务器
- 服务器部署目录固定
- compose 文件路径固定
- 你告诉我本次目标版本号

### 我来执行

我会按这个顺序做：

1. 登录服务器
2. 确认部署目录和当前镜像版本
3. 备份 compose / `.env`
4. 切镜像 tag
5. `docker compose pull`
6. `docker compose up -d`
7. 看日志和容器状态
8. 做你要求的最小验收
9. 如果异常，立即切回旧 tag

### 这种方式的边界

它适合：

- 单台或少量服务器
- 你本人掌控服务器
- 你希望每次升级都有人盯着验收

它不适合：

- 大规模多机自动化发布
- 需要完整 CI/CD 编排的团队化场景

---

## 不建议的操作方式

### 1. 登录服务器后直接点后台升级按钮

问题：

- 你只是换了执行入口，仍然没有解决镜像基线问题

### 2. compose 继续写 `latest`

问题：

- 无法稳定回滚
- 无法确认当前运行版本

### 3. 一次升级同时改版本和配置

问题：

- 出问题后很难定位原因

### 4. 没确认 release 镜像 ready 就上服务器升级

问题：

- 可能拉不到镜像
- 也可能拉到不完整发布链

---

## 最短执行版

如果以后只想保留一版最短流程，就按这个：

1. 本地确认目标版本 `vX.Y.Z`，对应镜像 `X.Y.Z`
2. `ssh` 登录服务器
3. `cd /opt/sub2api-deploy`
4. 备份 `docker-compose.yml` 和 `.env`
5. 把镜像改成 `ghcr.io/<owner>/sub2api:X.Y.Z`
6. 执行 `docker compose pull && docker compose up -d`
7. 看 `docker compose ps` 和 `docker compose logs --tail=100 sub2api`
8. 验 `/dashboard`、`/keys`、`/usage`、`/profile`、`/admin/dashboard`、`/admin/settings`
9. 有问题就把镜像切回上一个 tag 再 `docker compose up -d`

---

## 相关文档

- [docs/sst-production-upgrade-sop.md](./sst-production-upgrade-sop.md)
- [docs/sst-release-button-upgrade-plan.md](./sst-release-button-upgrade-plan.md)
