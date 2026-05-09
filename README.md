# 部署文档

## 目录
- [快速开始](#快速开始)
- [本地开发](#本地开发)
- [测试环境部署](#测试环境部署)
- [生产环境部署](#生产环境部署)
- [服务器部署指南](#服务器部署指南)
- [日志配置](#日志配置)
- [蓝绿部署原理](#蓝绿部署原理)
- [故障排查](#故障排查)

---

## 快速开始

### 前置要求
- Docker Desktop 已安装并运行
- Git 已安装
- Go 1.23+ 已安装
- 已登录腾讯云镜像仓库: `docker login ccr.ccs.tencentyun.com`

### 本地开发（宿主机运行）

```bash
# 启动本地数据库
go run cmd/deploy/main.go start-local

# 运行应用（Windows）
$env:GF_GCFG_FILE="config.local.yaml"
go run main.go

# 运行应用（Linux/Mac）
GF_GCFG_FILE=config.local.yaml go run main.go

# 停止数据库
go run cmd/deploy/main.go stop-local
```

### 本地蓝绿部署测试

```bash
go run cmd/deploy/main.go deploy local blue
go run cmd/deploy/main.go deploy local green
go run cmd/deploy/main.go rollback local
go run cmd/deploy/main.go status local
go run cmd/deploy/main.go cleanup local
```

访问地址：
- 应用: http://localhost:7001
- Traefik Dashboard: http://localhost:18080/dashboard/

---

## 本地开发

### 方式一：宿主机直接运行（推荐开发调试）

**使用 Go 工具（跨平台）:**
```bash
# 启动本地数据库
go run cmd/deploy/main.go start-local

# 运行应用（Windows PowerShell）
$env:GF_GCFG_FILE="config.local.yaml"
go run main.go

# 运行应用（Linux/Mac）
GF_GCFG_FILE=config.local.yaml go run main.go

# 停止数据库
go run cmd/deploy/main.go stop-local
```

**使用 Makefile（需要 make 命令）:**
```bash
make local.start  # 启动数据库
make local.run    # 运行应用
make local.stop   # 停止数据库
```

### 方式二：Docker 蓝绿部署（测试部署流程）

**使用 Go 工具（跨平台，推荐）:**
```bash
go run cmd/deploy/main.go deploy local blue    # 部署 blue
go run cmd/deploy/main.go deploy local green   # 切换到 green
go run cmd/deploy/main.go rollback local       # 回滚
go run cmd/deploy/main.go status local         # 查看状态
go run cmd/deploy/main.go cleanup local        # 清理
```

**使用 Makefile（需要 make 命令）:**
```bash
make deploy.local.blue   # 部署 blue
make deploy.local.green  # 切换到 green
make rollback.local      # 回滚
make status.local        # 查看状态
make cleanup.local       # 清理
```

访问地址：
- 应用: http://localhost:7001
- Traefik Dashboard: http://localhost:18080/dashboard/

---

## 测试环境部署

### 配置准备

1. 修改 `.env.test` 中的数据库连接信息
2. 修改 `manifest/config/config.test.yaml` 中的数据库配置

### 在服务器上部署

**使用 Go 工具（推荐）:**
```bash
# 克隆代码到服务器
git clone <your-repo-url>
cd server_go

# 配置环境
vim .env.test
vim manifest/config/config.test.yaml

# 部署
go run cmd/deploy/main.go build test
go run cmd/deploy/main.go push test
go run cmd/deploy/main.go deploy test blue
go run cmd/deploy/main.go deploy test green
go run cmd/deploy/main.go rollback test
go run cmd/deploy/main.go status test
```

**使用编译后的工具:**
```bash
# 本地编译
GOOS=linux GOARCH=amd64 go build -o deploy cmd/deploy/main.go

# 上传到服务器后运行
./deploy deploy test blue
```

**使用 Makefile:**
```bash
make deploy.test.blue   # 构建、推送并部署
make deploy.test.green  # 切换到 green
make rollback.test      # 回滚
make status.test        # 查看状态
```

---

## 生产环境部署

### 配置准备

1. 修改 `.env.production` 中的配置
2. 修改 `manifest/config/config.production.yaml` 中的数据库配置

### 在服务器上部署

**使用 Go 工具（推荐）:**
```bash
# 克隆代码到服务器
git clone <your-repo-url>
cd server_go

# 配置环境
vim .env.production
vim manifest/config/config.production.yaml

# 部署
go run cmd/deploy/main.go build production
go run cmd/deploy/main.go push production
go run cmd/deploy/main.go deploy production blue
go run cmd/deploy/main.go deploy production green
go run cmd/deploy/main.go rollback production
go run cmd/deploy/main.go status production
```

**使用编译后的工具:**
```bash
# 本地编译
GOOS=linux GOARCH=amd64 go build -o deploy cmd/deploy/main.go

# 上传到服务器后运行
./deploy deploy production blue
```

**使用 Makefile:**
```bash
make deploy.production.blue   # 构建、推送并部署
make deploy.production.green  # 切换到 green
make rollback.production      # 回滚
make status.production        # 查看状态
```

---

## 服务器部署指南

### 前置准备

#### 1. 服务器要求
- Docker 和 Docker Compose 已安装
- Git 已安装
- 已登录腾讯云镜像仓库: `docker login ccr.ccs.tencentyun.com`

#### 2. 克隆代码
```bash
git clone <your-repo-url>
cd server_go
```

#### 3. 配置环境变量

**测试环境:**
```bash
vim .env.test
vim manifest/config/config.test.yaml
```

**生产环境:**
```bash
vim .env.production
vim manifest/config/config.production.yaml
```

### 部署方式选择

#### 方式一：使用 Go 命令（推荐）

服务器上需要安装 Go 1.23+

```bash
# 测试环境部署
go run cmd/deploy/main.go build test
go run cmd/deploy/main.go push test
go run cmd/deploy/main.go deploy test blue

# 切换版本（0停机）
go run cmd/deploy/main.go deploy test green

# 回滚
go run cmd/deploy/main.go rollback test

# 查看状态
go run cmd/deploy/main.go status test

# 清理
go run cmd/deploy/main.go cleanup test
```

#### 方式二：编译部署工具

如果服务器上没有 Go 环境，可以在本地编译后上传：

```bash
# 本地编译（Linux 服务器）
GOOS=linux GOARCH=amd64 go build -o deploy cmd/deploy/main.go

# 上传到服务器
scp deploy user@server:/path/to/server_go/

# 服务器上运行
./deploy build test
./deploy push test
./deploy deploy test blue
./deploy deploy test green
./deploy rollback test
./deploy status test
```

#### 方式三：使用 Makefile

如果服务器有 make 命令：

```bash
# 测试环境
make deploy.test.blue
make deploy.test.green
make rollback.test
make status.test

# 生产环境
make deploy.production.blue
make deploy.production.green
make rollback.production
make status.production
```

### 典型部署流程

#### 首次部署

```bash
# 1. 配置环境
vim .env.test
vim manifest/config/config.test.yaml

# 2. 构建并推送镜像（可在本地或服务器执行）
go run cmd/deploy/main.go build test
go run cmd/deploy/main.go push test

# 3. 部署 blue 版本
go run cmd/deploy/main.go deploy test blue

# 4. 验证服务
curl http://localhost:7001/health/lb
```

#### 更新部署（0停机）

```bash
# 1. 拉取最新代码
git pull

# 2. 构建并推送新镜像
go run cmd/deploy/main.go build test
go run cmd/deploy/main.go push test

# 3. 部署到另一个颜色（假设当前是 blue）
go run cmd/deploy/main.go deploy test green

# 4. 验证新版本
curl http://localhost:7001/health/lb

# 5. 如果有问题，立即回滚
go run cmd/deploy/main.go rollback test
```

### CI/CD 集成

#### GitHub Actions 示例

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
            go run cmd/deploy/main.go deploy test green
```

---

## 日志配置

### 环境配置

#### 本地环境 (config.local.yaml)
- 日志级别: info
- 输出格式: 原始文本，带颜色
- 适合开发调试

#### Docker 本地测试 (config.docker.yaml)
- 日志级别: info
- 输出格式: 原始文本，带颜色
- 适合本地容器测试

#### 测试环境 (config.test.yaml)
- 日志级别: info
- 输出格式: 单行 JSON（无颜色）
- 适合腾讯云日志采集器收集

#### 生产环境 (config.production.yaml)
- 日志级别: info
- 输出格式: 单行 JSON（无颜色）
- 适合腾讯云日志采集器收集

### Docker 日志轮转

所有 Docker Compose 配置已启用日志轮转：
- 驱动: json-file
- 单文件大小: 20MB
- 保留文件数: 5个
- 总容量: 最多 100MB

配置位置：
- docker/compose/blue.yml
- docker/compose/green.yml
- docker/compose/traefik.yml

### 腾讯云日志采集

测试和生产环境的日志格式已优化为单行 JSON，方便腾讯云日志采集器收集：
- 无 ANSI 颜色代码
- JSON 格式，结构化数据
- 包含 RequestId 等上下文信息
- 每条日志一行，便于解析和查询

---

## 镜像管理

**使用 Go 工具（跨平台，推荐）:**
```bash
# 构建镜像
go run cmd/deploy/main.go build local
go run cmd/deploy/main.go build test
go run cmd/deploy/main.go build production

# 推送镜像
go run cmd/deploy/main.go push local
go run cmd/deploy/main.go push test
go run cmd/deploy/main.go push production
```

**使用 Makefile（需要 make 命令）:**
```bash
make build.local       # 构建本地镜像
make push.local        # 推送到 justsoso-local
make build.test        # 构建测试镜像
make push.test         # 推送到 justsoso-test
```

---

## 蓝绿部署原理

1. **Traefik** 作为网关，自动发现健康的容器
2. 部署新版本（blue/green）时，先启动新容器
3. 等待新容器健康检查通过（最多 60 秒）
4. 健康后，Traefik 自动将流量切换到新容器
5. 停止旧容器，完成 0 停机部署

### 健康检查

应用提供两个健康检查端点：
- `/health/ready` - 容器内部健康检查
- `/health/lb` - 负载均衡器健康检查

---

## 故障排查

### 查看日志
```bash
# 查看 blue 容器日志
docker compose -f docker/compose/blue.yml --env-file .env.local logs -f

# 查看 green 容器日志
docker compose -f docker/compose/green.yml --env-file .env.local logs -f

# 查看最近 100 行
docker compose -f docker/compose/blue.yml --env-file .env.local logs --tail=100
```

### 查看容器状态
```bash
go run cmd/deploy/main.go status test

# 或直接使用 docker
docker ps --filter "name=server-go"
```

### 访问 Traefik Dashboard
```
http://your-server:18080/dashboard/
```

### 常见问题

#### 容器无法启动
```bash
# 查看容器日志
docker compose -f docker/compose/blue.yml --env-file .env.test logs

# 检查配置文件
cat manifest/config/config.test.yaml

# 检查数据库连接
docker exec -it server-go-blue-1 sh
```

#### 健康检查失败
```bash
# 手动测试健康检查端点
curl http://localhost:7001/health/ready
curl http://localhost:7001/health/lb

# 查看应用日志
docker logs server-go-blue-1 --tail=50
```

#### 回滚失败
```bash
# 清理环境后重新部署
go run cmd/deploy/main.go cleanup test
go run cmd/deploy/main.go deploy test blue
```

---

## 安全建议

1. **不要在代码仓库中提交敏感信息**
   - `.env.*` 文件应该在 `.gitignore` 中
   - 数据库密码、API 密钥等应该通过环境变量或密钥管理系统配置

2. **限制服务器访问**
   - 使用防火墙限制端口访问
   - 只开放必要的端口（如 7001）
   - Traefik Dashboard 端口（18080）不要对外开放

3. **定期更新**
   - 定期更新 Docker 镜像
   - 定期更新依赖包
   - 定期备份数据库