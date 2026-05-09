# 部署文档

## 快速开始

### 本地开发

```bash
# 1. 启动本地数据库（MySQL + Redis）
go run cmd/deploy/main.go start-local-db

# 2. 运行应用（自动使用 config.local.yaml）
go run main.go

# 3. 停止数据库
go run cmd/deploy/main.go stop-local-db
```

访问: http://localhost:7001

### 使用 Makefile（推荐）

Makefile 提供了简化的命令封装：

```bash
make start-local-db    # 启动本地数据库（等同于 go run cmd/deploy/main.go start-local-db）
make dev               # 运行应用（等同于 go run main.go）
make stop-local-db     # 停止本地数据库（等同于 go run cmd/deploy/main.go stop-local-db）
```

> **说明**：Makefile 命令与部署脚本命令完全一致，只是省略了 `go run cmd/deploy/main.go` 前缀。

---

## 部署命令详解

### 命令格式

```bash
go run cmd/deploy/main.go <command> <env> [options]
```

**参数说明：**
- `<command>`: 命令类型（必填）
- `<env>`: 环境名称（必填，除了 start-local-db 和 stop-local-db）
- `[options]`: 可选参数，格式为 `key=value`

### 1. 构建镜像 (build)

构建 Docker 镜像并打标签。

**命令组合：**

```bash
# 构建本地环境镜像
go run cmd/deploy/main.go build local
# 说明：构建镜像并推送到 ccr.ccs.tencentyun.com/justsoso-local/server-go
# 标签：使用 git commit hash（如 abc1234）或 abc1234.dirty（有未提交修改）

# 构建测试环境镜像
go run cmd/deploy/main.go build test
# 说明：构建镜像并推送到 ccr.ccs.tencentyun.com/justsoso-test/server-go
# 标签：自动使用 git commit hash

# 构建生产环境镜像
go run cmd/deploy/main.go build production
# 说明：构建镜像并推送到 ccr.ccs.tencentyun.com/justsoso-production/server-go
# 标签：自动使用 git commit hash

# 构建测试/生产环境镜像（必须指定版本）
go run cmd/deploy/main.go build production version=v1.2.3
# 说明：构建生产环境镜像，使用指定的版本号 v1.2.3 作为标签
# 镜像：ccr.ccs.tencentyun.com/justsoso-production/server-go:v1.2.3
# 同时打 latest 标签
# 注意：test/production 环境必须指定 version 参数
```

**构建过程：**
1. 读取 git commit hash 作为默认版本号（或使用 version 参数）
2. 使用 `docker/Dockerfile` 构建镜像
3. 打两个标签：指定版本号 + latest
4. 镜像保存在本地，等待推送

### 2. 推送镜像 (push)

将本地构建的镜像推送到腾讯云镜像仓库。

**命令组合：**

```bash
# 推送本地环境镜像
go run cmd/deploy/main.go push local
# 说明：推送到 justsoso-local 命名空间
# 推送：server-go:版本号 和 server-go:latest

# 推送测试环境镜像
go run cmd/deploy/main.go push test
# 说明：推送到 justsoso-test 命名空间
# 用途：供测试服务器拉取部署

# 推送生产环境镜像
go run cmd/deploy/main.go push production
# 说明：推送到 justsoso-production 命名空间
# 用途：供生产服务器拉取部署

# 推送测试/生产环境镜像（必须指定版本）
go run cmd/deploy/main.go push production version=v1.2.3
# 说明：推送指定版本号的镜像
# 注意：test/production 环境必须指定 version 参数
# 前提：必须先用相同的 version 参数执行 build 命令
```

**推送过程：**
1. 查找本地已构建的镜像
2. 推送到对应环境的腾讯云镜像仓库
3. 同时推送版本号标签和 latest 标签

### 3. 部署 (deploy)

执行蓝绿部署，自动检测当前运行的颜色并切换到另一个颜色。

**命令组合：**

```bash
# 部署到本地环境
go run cmd/deploy/main.go deploy local
# 说明：在本地 Docker 环境执行蓝绿部署
# 流程：
#   1. 检测当前运行的颜色（blue 或 green）
#   2. 如果 blue 在运行，部署到 green；如果 green 在运行，部署到 blue
#   3. 如果都没运行（首次部署），默认部署到 blue
#   4. 启动 Traefik 网关（如果未运行）
#   5. 启动新颜色的容器
#   6. 等待健康检查通过（最多 60 秒）
#   7. 健康检查通过后，Traefik 自动切换流量到新容器
#   8. 停止旧颜色的容器
#   9. 如果健康检查失败，停止新容器，保持旧版本运行（自动回滚）

# 部署到测试环境
go run cmd/deploy/main.go deploy test
# 说明：在测试服务器执行蓝绿部署
# 前提：已执行 build test 和 push test
# 配置：使用 .env.test 和 config.test.yaml
# 数据库：连接测试环境的外部数据库

# 部署到生产环境
go run cmd/deploy/main.go deploy production
# 说明：在生产服务器执行蓝绿部署
# 前提：已执行 build production 和 push production
# 配置：使用 .env.production 和 config.production.yaml
# 数据库：连接生产环境的外部数据库
# 特点：0 停机部署，失败自动回滚

# 部署指定版本（可选）
go run cmd/deploy/main.go deploy production version=v1.2.3
# 说明：部署指定版本号的镜像到生产环境
# 前提：该版本已经 build 和 push
# 注意：version 参数可选，默认使用 latest
```

**部署流程详解：**

1. **检测阶段**
   - 检查 blue 容器是否运行
   - 检查 green 容器是否运行
   - 确定当前颜色和目标颜色

2. **准备阶段**
   - 启动 Traefik 网关（如果未运行）
   - 等待 5 秒让网关就绪

3. **部署阶段**
   - 使用 docker compose 启动目标颜色的容器
   - 容器使用对应环境的配置文件
   - 容器自动注册到 Traefik

4. **健康检查阶段**
   - 每 3 秒检查一次容器健康状态
   - 最多等待 60 秒
   - 检查 `/health/lb` 端点

5. **切换阶段**
   - 健康检查通过后，Traefik 自动将流量切到新容器
   - 停止旧颜色的容器
   - 完成 0 停机部署

6. **失败回滚**
   - 如果健康检查超时（60 秒）
   - 自动停止新容器
   - 保持旧版本继续运行
   - 输出错误日志

7. **镜像清理**
   - 部署成功后自动清理旧镜像
   - 只保留最近 10 个版本的镜像
   - 避免磁盘空间占用过多

### 4. 查看状态 (status)

查看当前机器上的容器运行状态。

**命令：**

```bash
# 查看容器状态（无需指定环境）
go run cmd/deploy/main.go status
# 显示：
#   - 运行中的容器（blue/green/traefik）
#   - 容器状态（运行时间、健康状态）
#   - 端口映射
#   - 网络信息
#   - 数据卷信息
```

> **说明**：status 命令不需要指定环境参数，它会显示当前机器上所有 server-go 相关的容器。

**输出示例：**
```
=== Environment: test ===

Running containers:
NAMES               STATUS                    PORTS
server-go-green-1   Up 2 hours (healthy)     7001/tcp
server-go-gateway   Up 2 hours (healthy)     0.0.0.0:7001->7001/tcp

Networks:
NETWORK ID     NAME                DRIVER    SCOPE
abc123         server-go-network   bridge    local

Volumes:
DRIVER    VOLUME NAME
```

### 5. 启动本地数据库 (start-local-db)

启动本地开发用的 MySQL 和 Redis 容器。

```bash
go run cmd/deploy/main.go start-local-db
# 说明：启动本地数据库服务
# 启动：
#   - MySQL 8.4 容器，端口映射到 127.0.0.1:330
#   - Redis 7.4 容器，端口映射到 127.0.0.1:637
# 配置：
#   - MySQL root 密码：root
#   - Redis 密码：root
# 数据持久化：使用 Docker volume
# 用途：本地开发时，应用连接这些数据库
```

**启动后：**
- MySQL: `mysql:root:root@tcp(127.0.0.1:330)/game_db_1`
- Redis: `127.0.0.1:637`，密码 `root`
- 数据保存在 Docker volume 中，重启不丢失

### 6. 停止本地数据库 (stop-local-db)

停止并删除本地数据库容器。

```bash
go run cmd/deploy/main.go stop-local-db
# 说明：停止本地数据库服务
# 操作：
#   - 停止 MySQL 容器
#   - 停止 Redis 容器
#   - 删除容器（数据卷保留）
# 注意：数据不会丢失，下次启动会恢复
```

---

## 完整部署流程示例

### 场景 1：本地开发

```bash
# 1. 启动数据库
go run cmd/deploy/main.go start-local-db
# 结果：MySQL 和 Redis 启动在 127.0.0.1

# 2. 运行应用
go run main.go
# 结果：应用自动使用 config.local.yaml，连接本地数据库

# 3. 停止数据库（开发结束后）
go run cmd/deploy/main.go stop-local-db
```

### 场景 2：测试环境首次部署

```bash
# 1. 配置环境（在服务器上）
vim .env.test
vim manifest/config/config.test.yaml

# 2. 构建镜像（可在本地或服务器）
go run cmd/deploy/main.go build test
# 结果：构建镜像，标签为 git commit hash

# 3. 推送镜像
go run cmd/deploy/main.go push test
# 结果：镜像推送到 justsoso-test 命名空间

# 4. 部署（在服务器上）
go run cmd/deploy/main.go deploy test
# 结果：
#   - 检测到无运行容器，部署到 blue
#   - 启动 Traefik 网关
#   - 启动 blue 容器
#   - 等待健康检查
#   - 服务可访问

# 5. 查看状态
go run cmd/deploy/main.go status test
# 结果：显示 blue 容器运行中
```

### 场景 3：测试环境更新部署

```bash
# 1. 拉取最新代码（在服务器上）
git pull

# 2. 构建新版本
go run cmd/deploy/main.go build test
# 结果：构建新的 commit hash 版本

# 3. 推送新版本
go run cmd/deploy/main.go push test

# 4. 部署新版本
go run cmd/deploy/main.go deploy test version=def5678
# 结果：
#   - 检测到 blue 在运行
#   - 部署到 green
#   - 启动 green 容器
#   - 等待健康检查
#   - 健康检查通过，流量切到 green
#   - 停止 blue 容器
#   - 完成 0 停机更新

# 5. 查看状态
go run cmd/deploy/main.go status test
# 结果：显示 green 容器运行中
```

### 场景 4：生产环境发布指定版本

```bash
# 1. 构建指定版本
go run cmd/deploy/main.go build production version=v1.2.3
# 结果：构建镜像，标签为 v1.2.3

# 2. 推送指定版本
go run cmd/deploy/main.go push production version=v1.2.3
# 结果：推送 v1.2.3 到 justsoso-production

# 3. 部署指定版本
go run cmd/deploy/main.go deploy production version=v1.2.3
# 结果：
#   - 检测当前运行颜色
#   - 部署 v1.2.3 到另一个颜色
#   - 健康检查通过后切换
#   - 0 停机发布

# 4. 查看状态
go run cmd/deploy/main.go status production
```

### 场景 5：部署失败自动回滚

```bash
# 1. 部署新版本
go run cmd/deploy/main.go deploy test
# 过程：
#   - 检测到 blue 在运行
#   - 开始部署到 green
#   - 启动 green 容器
#   - 等待健康检查...
#   - 健康检查失败（应用启动失败）
#   - 自动停止 green 容器
#   - blue 容器继续运行
#   - 输出错误信息

# 2. 查看状态
go run cmd/deploy/main.go status test
# 结果：blue 仍在运行，服务未中断

# 3. 排查问题后重新部署
go run cmd/deploy/main.go deploy test
# 结果：再次尝试部署到 green
```

---

## 环境配置

### 配置文件

| 环境 | 配置文件 | 说明 |
|------|---------|------|
| 本地开发 | `manifest/config/config.local.yaml` | 宿主机运行，连接本地 Docker 数据库 |
| Docker 本地 | `manifest/config/config.docker.yaml` | 容器内运行 |
| 测试环境 | `manifest/config/config.test.yaml` | 外部数据库 |
| 生产环境 | `manifest/config/config.production.yaml` | 外部数据库 |

### 环境变量文件

- `.env.local` - 本地环境
- `.env.test` - 测试环境
- `.env.production` - 生产环境

---

## 服务器部署

### 前置准备

1. 服务器安装 Docker 和 Docker Compose
2. 登录腾讯云镜像仓库: `docker login ccr.ccs.tencentyun.com`
3. 克隆代码: `git clone <repo-url> && cd server_go`

### 测试环境部署

```bash
# 1. 配置环境
vim .env.test
vim manifest/config/config.test.yaml

# 2. 构建并推送（可在本地或服务器执行）
go run cmd/deploy/main.go build test
go run cmd/deploy/main.go push test

# 3. 部署（自动蓝绿切换）
go run cmd/deploy/main.go deploy test

# 4. 查看状态
go run cmd/deploy/main.go status test
```

### 生产环境部署

```bash
# 1. 配置环境
vim .env.production
vim manifest/config/config.production.yaml

# 2. 构建并推送
go run cmd/deploy/main.go build production
go run cmd/deploy/main.go push production

# 3. 部署
go run cmd/deploy/main.go deploy production

# 4. 查看状态
go run cmd/deploy/main.go status production
```

### 编译部署工具（服务器无 Go 环境）

```bash
# 本地编译
GOOS=linux GOARCH=amd64 go build -o deploy cmd/deploy/main.go

# 上传到服务器
scp deploy user@server:/path/to/server_go/

# 服务器上运行
./deploy build test
./deploy push test
./deploy deploy test
```

---

## 蓝绿部署原理

1. **自动检测**：检测当前运行的颜色（blue 或 green）
2. **部署新版本**：启动另一个颜色的容器
3. **健康检查**：等待新容器健康检查通过（最多 60 秒）
4. **自动切换**：Traefik 自动将流量切换到健康的新容器
5. **停止旧版本**：停止旧颜色的容器
6. **失败回滚**：如果健康检查失败，自动停止新容器，保持旧版本运行

### 健康检查端点

- `/health/ready` - 容器内部健康检查
- `/health/lb` - 负载均衡器健康检查

---

## 日志配置

### 日志级别

所有环境统一使用 `info` 级别

### 日志格式

| 环境 | 格式 | 说明 |
|------|------|------|
| 本地开发 | 原始文本，带颜色 | 便于开发调试 |
| 测试环境 | 单行 JSON | 腾讯云日志采集 |
| 生产环境 | 单行 JSON | 腾讯云日志采集 |

### Docker 日志轮转

所有容器已配置日志轮转：
- 驱动: `json-file`
- 单文件大小: `20MB`
- 保留文件数: `5个`
- 总容量: 最多 100MB

---

## 监控和日志

### 查看容器日志

```bash
# 查看 blue 容器日志
docker compose -f docker/compose/blue.yml --env-file .env.test logs -f

# 查看 green 容器日志
docker compose -f docker/compose/green.yml --env-file .env.test logs -f

# 查看最近 100 行
docker logs server-go-blue-1 --tail=100
```

### 访问 Traefik Dashboard（仅本地环境）

本地环境默认开启 Dashboard：
```
http://localhost:18080/dashboard/
```

测试/生产环境默认关闭 Dashboard。如需开启，修改对应 `.env` 文件：
```bash
TRAEFIK_API_ENABLED=true
TRAEFIK_DASHBOARD_ENABLED=true
TRAEFIK_DASHBOARD_PORT=18080
```

---

## 故障排查

### 容器无法启动

```bash
# 查看容器日志
docker logs server-go-blue-1 --tail=50

# 检查配置文件
cat manifest/config/config.test.yaml

# 进入容器检查
docker exec -it server-go-blue-1 sh
```

### 健康检查失败

```bash
# 手动测试健康检查端点
curl http://localhost:7001/health/ready
curl http://localhost:7001/health/lb

# 查看应用日志
docker logs server-go-blue-1 --tail=50
```

### 部署失败

部署失败会自动回滚，保持旧版本运行。查看日志排查问题后重新部署即可。

---

## CI/CD 集成

### GitHub Actions 示例

```yaml
name: Deploy to Test

on:
  push:
    branches: [test]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Login to Tencent Cloud Registry
        run: echo "${{ secrets.TCR_PASSWORD }}" | docker login ccr.ccs.tencentyun.com -u ${{ secrets.TCR_USERNAME }} --password-stdin
      
      - name: Build and Push
        run: |
          go run cmd/deploy/main.go build test
          go run cmd/deploy/main.go push test
      
      - name: Deploy via SSH
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.TEST_HOST }}
          username: ${{ secrets.TEST_USER }}
          key: ${{ secrets.TEST_SSH_KEY }}
          script: |
            cd /path/to/server_go
            git pull
            go run cmd/deploy/main.go deploy test
```

---

## 安全建议

1. **不要提交敏感信息**
   - `.env.*` 文件应在 `.gitignore` 中
   - 数据库密码、API 密钥通过环境变量配置

2. **限制服务器访问**
   - 使用防火墙限制端口访问
   - 只开放必要的端口（如 7001）
   - 生产环境禁用 Traefik Dashboard（默认已关闭）

3. **定期更新**
   - 定期更新 Docker 镜像
   - 定期更新依赖包
   - 定期备份数据库