# server_go

基于 GoFrame v2 的游戏服务端项目，提供用户登录、游戏状态、背包、棋盘、资源变更、健康检查和内部流量控制等 HTTP 接口。

## 项目定位

`server_go` 是一个面向游戏业务的后端服务，核心职责包括：

- 提供游戏客户端访问的 `/api` 接口。
- 管理用户、资源、背包、棋盘、任务等游戏数据。
- 连接 MySQL 存储业务数据。
- 连接 Redis 支撑缓存、状态和运行时能力。
- 提供健康检查与内部控制接口，支持容器化部署和蓝绿发布。

## 技术栈

- Go `1.26.0`
- GoFrame `v2.10.0`
- MySQL `8.4`
- Redis `7.4`
- Docker / Docker Compose

## 目录结构

```text
.
├── api/                    # 接口协议定义
├── cmd/                    # 辅助命令，如部署、本地压测
├── hack/                   # 构建、部署相关脚本
├── internal/
│   ├── cmd/                # HTTP 服务入口与路由注册
│   ├── controller/         # 请求处理层
│   ├── dao/                # 数据访问层
│   ├── logic/              # 业务逻辑层
│   ├── middleware/         # 签名、登录态校验等中间件
│   ├── model/              # 数据模型
│   └── service/            # 服务接口声明
├── manifest/
│   ├── config/             # 不同环境配置
│   └── docker/             # Dockerfile、Compose、数据库初始化脚本
├── resource/               # 静态资源与模板目录
├── utility/                # 通用工具
├── main.go                 # 服务启动入口
├── Makefile                # 常用开发、构建、部署命令
└── go.mod                  # Go 模块依赖
```

## 环境要求

本地开发建议准备：

- Go `1.26+`
- Docker Desktop
- Docker Compose
- Make，可选；Windows 下也可以直接执行对应的 `go run` 命令

## 快速启动

### 1. 安装依赖

```bash
go mod download
```

### 2. 启动本地 MySQL 和 Redis

如果当前环境支持 `make`：

```bash
make start-local-db
```

等价命令：

```bash
go run cmd/deploy/main.go start-local-db
```

本地开发数据库默认信息：

```text
MySQL: 127.0.0.1:330
Database: game_db_1
User: root
Password: root

Redis: 127.0.0.1:637
Password: root
DB: 0
```

初始化 SQL 位于：

```text
manifest/docker/mysql/init/01-init.sql
```

### 3. 启动服务

```bash
go run main.go
```

或：

```bash
make dev
```

服务默认监听：

```text
http://127.0.0.1:7001
```

## 配置说明

默认配置文件位于：

```text
manifest/config/config.yaml
```

常见环境配置：

```text
manifest/config/config.local.yaml       # 本地容器内运行配置
manifest/config/config.test.yaml        # 测试环境配置
manifest/config/config.production.yaml  # 生产环境配置
```

关键配置项：

| 配置项 | 说明 |
| --- | --- |
| `server.address` | HTTP 服务监听地址 |
| `server.openapiPath` | OpenAPI JSON 地址 |
| `server.swaggerPath` | Swagger UI 地址 |
| `app.keys` | 应用签名和加密密钥 |
| `database.default.link` | MySQL 连接地址 |
| `database.default.debug` | 是否开启 SQL 调试日志 |
| `redis.default.address` | Redis 地址 |
| `redis.default.pass` | Redis 密码 |
| `logger.level` | 日志级别 |

容器运行时可通过 `APP_PORT` 覆盖监听端口。

## 接口分组

### 游戏业务接口

`/api` 分组启用了签名校验和登录态校验。

| 分组 | 路径 | 说明 |
| --- | --- | --- |
| User | `/api/user/login` | 用户登录 |
| Game | `/api/game/time` | 获取服务器时间 |
| Game | `/api/game/online` | 记录在线时长 |
| Bag | `/api/bag/get_bag/{chapter}` | 获取用户背包 |
| Bag | `/api/bag/get_bag_tp/{chapter}` | 获取用户背包 TP |
| Grid | `/api/grid/get/{chapter}` | 获取棋盘数据 |
| Res | `/api/res/add_tili` | 测试增加体力 |
| Res | `/api/res/add_gold` | 测试增加金币 |
| Res | `/api/res/add_diamond` | 测试增加钻石 |

### 非登录业务接口

| 分组 | 路径 | 说明 |
| --- | --- | --- |
| Other | `/other/res_version/{key}` | 获取资源版本号 |
| Test | `/test/` | 测试接口 |
| Test | `/test/db` | 测试数据库 |

### 健康检查接口

| 路径 | 说明 |
| --- | --- |
| `/health/` | 基础健康检查 |
| `/health/ready` | 就绪检查 |
| `/health/detail` | 健康详情 |
| `/health/lb` | 负载均衡健康检查 |

### 内部控制接口

内部控制接口位于 `/internal/control`，主要用于部署、流量切换和优雅下线。

| 路径 | 说明 |
| --- | --- |
| `/internal/control/traffic-shift` | 开始流量切换 |
| `/internal/control/reject-new-requests` | 拒绝新请求 |
| `/internal/control/resume-traffic` | 恢复流量 |

## OpenAPI 与 Swagger

服务启动后可访问：

```text
http://127.0.0.1:7001/api.json
http://127.0.0.1:7001/swagger
```

## 常用命令

### 本地开发

```bash
make dev
make start-local-db
make stop-local-db
```

等价命令：

```bash
go run main.go
go run cmd/deploy/main.go start-local-db
go run cmd/deploy/main.go stop-local-db
```

### 构建镜像

```bash
make build.local
make build.test VERSION=1.0.0
make build.production VERSION=1.0.0
```

### 推送镜像

```bash
make push.local
make push.test VERSION=1.0.0
make push.production VERSION=1.0.0
```

### 部署

```bash
make deploy.local
make deploy.test VERSION=1.0.0
make deploy.production VERSION=1.0.0
```

### 查看状态

```bash
make status.local
make status.test
make status.production
```

## Docker 本地环境

本地数据库服务由以下文件维护：

```text
manifest/docker/compose/local.yml
```

包含：

- MySQL `8.4`
- Redis `7.4-alpine`
- 独立 Docker network
- MySQL 和 Redis named volume 持久化
- 首次启动自动执行初始化 SQL

停止本地数据库：

```bash
make stop-local-db
```

如果需要彻底清空本地数据，需要额外删除 Docker volume。

## 部署说明

项目内置了部署辅助命令：

```text
cmd/deploy/main.go
```

支持的环境：

```text
local
test
production
```

部署相关资源位于：

```text
manifest/docker/
manifest/docker/compose/
```

其中 `blue.yml`、`green.yml`、`traefik.yml` 用于蓝绿部署和流量入口管理。

## 开发约定

- 新接口优先在 `api/<module>/v1` 定义协议。
- 对外业务接口挂载到 `/api`，默认经过签名和登录态校验。
- 健康检查和内部控制接口不依赖登录态。
- 数据访问模型由 `internal/dao`、`internal/model` 承载。
- 业务编排放在 `internal/logic`。
- 公共能力放在 `utility`。

## 注意事项

- `app.keys` 涉及签名和加密，生产环境必须替换。
- `.env.*` 和环境配置中可能包含环境差异，请不要直接混用生产配置启动本地服务。
- `/api` 分组接口需要满足项目内签名和校验规则，直接浏览器访问可能无法通过中间件。
- 本地 MySQL 使用宿主机端口 `330`，如端口被占用，需要同步修改 Compose 和配置文件。