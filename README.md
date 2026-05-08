# server_go 项目说明

本文档面向 Go 和 GoFrame 新手，目标是让你能从零理解这个项目：怎么启动、怎么配置、怎么调用 API、每个目录做什么、每条常用命令做什么，以及一次请求在代码里怎样流转。

项目基于 GoFrame v2，提供一组游戏服务接口：登录、背包、棋盘、在线时长、资源版本、健康检查和蓝绿发布控制。

---

## 1. 你需要先知道的概念

### 1.1 Go 是什么

Go 是一种编译型语言。项目里的 `.go` 文件不会直接在线上运行，而是先通过 `go build` 编译成一个可执行文件，然后运行这个可执行文件。

常用命令：

```bash
go version
```

作用：查看本机 Go 版本。

```bash
go mod download
```

作用：根据 `go.mod` 和 `go.sum` 下载依赖。

```bash
go run .
```

作用：编译并运行当前目录下的 Go 程序。适合本地开发。

```bash
go build -o main .
```

作用：把当前项目编译成名为 `main` 的可执行文件。

```bash
go test ./...
```

作用：编译并测试当前模块下所有 Go 包。即使项目没有测试文件，这条命令也能检查所有包是否能正常编译。

```bash
go vet ./...
```

作用：做 Go 官方静态检查，能发现一些编译器不报错但可能有问题的写法。

```bash
gofmt -w .
```

作用：格式化当前目录下所有 Go 代码。Go 项目要求统一格式，通常提交前都要执行。

### 1.2 GoFrame 是什么

GoFrame 是 Go 语言 Web 框架。本项目主要使用了 GoFrame 的这些能力：

- HTTP Server：启动 Web 服务。
- Object Router：通过 `api` 里的请求结构体绑定路由。
- Middleware：处理签名、登录校验、统一响应。
- Config：读取 YAML 配置。
- Database DAO：访问 MySQL。
- Redis：做分布式锁、防重放和限流。
- CLI：通过 `gf` 命令生成 DAO、构建镜像等。

---

## 2. 项目目录结构

核心目录如下：

```text
server_go/
  api/                    # HTTP API 的请求和响应结构体
  internal/
    cmd/                  # 程序启动和路由注册
    consts/               # 常量定义
    controller/           # HTTP 控制器：接收 api.Req，调用 service
    dao/                  # GoFrame 生成的 DAO 数据访问对象
    logic/                # 业务逻辑实现
    middleware/           # HTTP 中间件
    model/                # 内部业务输入/输出模型 + 数据库模型
      do/                 # GoFrame 生成的 DO 模型，用于 DAO Data/Where
      entity/             # GoFrame 生成的实体模型，用于 Scan
      model.go            # 手写的业务层 Input/Output
    packed/               # GoFrame 资源打包入口
    service/              # 业务接口定义和注册
  utility/                # 通用工具函数
  manifest/config/        # 运行配置文件
  docker/                 # Dockerfile、docker compose、MySQL 初始化
  hack/                   # Makefile 引用的 GoFrame CLI 配置
  main.go                 # 程序入口
  go.mod                  # Go 模块和依赖声明
  Makefile                # 常用构建命令入口
```

### 2.1 `api/`

`api` 目录定义 HTTP 接口的输入输出。每个请求结构体里都有 `g.Meta`：

```go
g.Meta `path:"/user/login" method:"get,post" tags:"User" summary:"登录"`
```

含义：

- `path`：接口路径。
- `method`：允许的 HTTP 方法。
- `tags`：OpenAPI/Swagger 分组。
- `summary`：接口说明。

字段标签示例：

```go
Uid int64 `json:"uid" v:"required"`
```

含义：

- `json:"uid"`：HTTP 参数名是 `uid`。
- `v:"required"`：GoFrame 参数校验，表示必填。

### 2.2 `internal/controller/`

Controller 是 HTTP 层和业务层的连接点。

典型流程：

1. 接收 `api.XxxReq`。
2. 转成 `model.XxxInput`。
3. 调用 `service.Xxx()`。
4. 把 `model.XxxOutput` 转成 `api.XxxRes`。

例如登录接口：

```text
api/user.LoginReq
  -> controller.User.Login
  -> service.User().Login
  -> logic/user.Login
  -> dao 查询/写入数据库
  -> model.LoginOutput
  -> api/user.LoginRes
```

### 2.3 `internal/service/`

`service` 目录定义业务接口，不直接写业务逻辑。

例如 `service.IUser` 定义了用户业务应该提供的方法：

- `Login`
- `UpdateDiamond`
- `UpdateGold`
- `UpdateTili`
- `UpdateExp`
- `UpdateStar`
- `GetUser`
- `GetUserRes`

`service.User()` 返回当前注册的用户服务实现。

`logic/user` 包在 `init()` 中执行：

```go
service.RegisterUser(&sUser{})
```

这表示把 `logic/user` 里的 `sUser` 注册成用户服务实现。

### 2.4 `internal/logic/`

`logic` 目录是真正写业务逻辑的地方。

当前模块：

- `logic/user`：登录、用户资源更新、用户资源读取。
- `logic/bag`：读取用户背包。
- `logic/game`：记录在线时长。
- `logic/grid`：聚合背包、背包模板、任务数据。
- `logic/other`：资源版本查询。
- `logic/task`：任务初始化。
- `logic/lock`：Redis 分布式锁。
- `logic/gamelog`：异步写日志。

### 2.5 `internal/model/model.go`

这个文件是内部业务层的输入和输出模型，不是 HTTP API 模型。

命名规则：

- `XxxInput`：service/logic 入参。
- `XxxOutput`：service/logic 出参。

为什么不直接用 `api` 里的结构体：

- `api` 是 HTTP 协议层，包含路由、参数校验、OpenAPI 信息。
- `model` 是内部业务层，service/logic 不应该依赖 HTTP 细节。
- 以后如果业务逻辑被定时任务、后台脚本或 RPC 调用，可以继续使用 `model`，不用引入 `api`。

### 2.6 `internal/dao/`、`internal/model/entity/`、`internal/model/do/`

这些文件通常由 GoFrame CLI 根据数据库表生成。

- `dao`：表访问入口，例如 `dao.User`、`dao.UserRes`。
- `entity`：数据库实体结构体，通常用于 `Scan(&entity)`。
- `do`：DAO 操作用的数据对象，通常用于 `Where`、`Data`。

常见 DAO 调用：

```go
dao.User.Ctx(ctx).Where("uid", uid).Scan(&user)
```

逐段解释：

- `dao.User`：访问 `user` 表。
- `Ctx(ctx)`：绑定请求上下文。
- `Where("uid", uid)`：添加 SQL 条件 `uid = ?`。
- `Scan(&user)`：把查询结果扫描到结构体里。

```go
dao.UserRes.Ctx(ctx).Where("uid", uid).Data(g.Map{"gold": 100}).Update()
```

含义：更新 `user_res` 表中指定用户的金币字段。

```go
dao.UserLoginkey.Ctx(ctx).Data(g.Map{...}).Save()
```

含义：保存登录 key。`Save` 通常表示插入或更新，具体行为依赖数据库和表约束。

---

## 3. 启动流程详解

程序入口是 `main.go`。

### 3.1 `main.go`

主要做三件事：

1. 空导入 `internal/packed`。
2. 空导入所有 `logic/*` 包。
3. 运行 `cmd.Main`。

空导入写法：

```go
_ "server_go/internal/logic/user"
```

含义：不直接使用这个包里的变量或函数，但要执行这个包的 `init()`。

为什么需要空导入 `logic/user`：

- `logic/user` 的 `init()` 会调用 `service.RegisterUser(&sUser{})`。
- 如果不导入，`init()` 不执行，`service.User()` 会因为未注册而 panic。

### 3.2 `internal/cmd/cmd.go`

这里创建 HTTP Server 并注册路由。

```go
s := g.Server()
```

含义：获取 GoFrame 默认 HTTP Server。

```go
s.Group("/api", func(group *ghttp.RouterGroup) { ... })
```

含义：创建 `/api` 路由组。组内接口路径都会自动带 `/api` 前缀。

`/api` 路由组挂了三个中间件：

```go
group.Middleware(
    middleware.Sign,
    middleware.Verify,
    middleware.Response,
)
```

执行顺序：

1. `Sign`：校验请求签名。
2. `Verify`：校验登录态、时间戳、防重放。
3. `Response`：统一包装返回值。

然后绑定这些控制器：

```go
group.Bind(
    controller.User,
    controller.Game,
    controller.Bag,
    controller.Grid,
)
```

GoFrame 会扫描控制器方法和 `api` 里的 `g.Meta`，自动注册接口。

另外还有两个 `/` 路由组：

- `controller.Other`：资源版本接口，不走签名和登录校验，只走统一响应。
- `controller.Health`：健康检查和内部控制接口，不走统一响应，因为控制器自己写 JSON。

---

## 4. 配置文件

配置文件在 `manifest/config/`。

当前主要配置：

```text
manifest/config/config.docker.yaml
manifest/config/config.test.yaml
manifest/config/config.production.yaml
```

Docker 镜像默认使用：

```text
config.docker.yaml
```

因为 `docker/Dockerfile` 中设置了：

```dockerfile
ENV GF_GCFG_FILE=config.docker.yaml
```

### 4.1 server 配置

```yaml
server:
  address: ":7001"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"
  graceful: true
  gracefulTimeout: 15
```

字段解释：

- `address: ":7001"`：监听 7001 端口。
- `openapiPath`：OpenAPI JSON 地址。
- `swaggerPath`：Swagger UI 地址。
- `graceful: true`：开启优雅退出。
- `gracefulTimeout: 15`：优雅退出最多等待 15 秒。

### 4.2 app.keys

```yaml
app:
  keys:
    - "SzNg8LgHjEUgTAc4"
```

这是签名密钥列表。`middleware.Sign` 会用这些 key 计算 HMAC-SHA256 签名。

### 4.3 database

```yaml
database:
  default:
    link: "mysql:root:root@tcp(mysql:3306)/game_db_1?charset=utf8mb4&parseTime=true&loc=Local"
    debug: false
```

连接串格式：

```text
mysql:用户名:密码@tcp(主机:端口)/数据库名?参数
```

当前 Docker 配置连接到 compose 网络里的 MySQL 服务：

```text
mysql:3306
```

### 4.4 redis

```yaml
redis:
  default:
    address: "redis:6379"
    pass: "root"
    db: 0
```

字段解释：

- `address`：Redis 地址。
- `pass`：Redis 密码。
- `db`：Redis DB 编号。

Redis 在项目中用于：

- 防重放：`middleware.Verify` 使用 `replay:{uid}:{sign}`。
- 分布式锁：`logic/lock` 使用 `lock:{key}`。
- 资源版本接口限流：`logic/other` 使用 `res_version.{key}`。

---

## 5. 本地运行

### 5.1 安装 Go

安装 Go 1.23 或兼容版本后检查：

```bash
go version
```

如果能看到版本号，说明 Go 可用。

### 5.2 下载依赖

在项目根目录执行：

```bash
go mod download
```

含义：下载 `go.mod` 中声明的依赖。

### 5.3 启动 MySQL 和 Redis

本项目提供了本地 compose 文件：

```bash
docker compose -f docker/compose/local.yml up -d
```

含义：后台启动 MySQL 和 Redis。

`-f docker/compose/local.yml` 指定 compose 文件。

`up` 表示启动服务。

`-d` 表示后台运行。

默认暴露端口：

- MySQL：宿主机 `330` -> 容器 `3306`
- Redis：宿主机 `637` -> 容器 `6379`

注意：`config.docker.yaml` 里的数据库地址是 `mysql:3306`，适合容器内访问。如果你直接在宿主机用 `go run .`，需要准备一个能被 GoFrame 读取的本地配置，并把数据库地址改成宿主机能访问的地址，例如 `127.0.0.1:330`。

### 5.4 初始化数据库

compose 会挂载：

```text
docker/mysql/init:/docker-entrypoint-initdb.d:ro
```

其中 `01-init.sql` 会创建数据库：

```sql
CREATE DATABASE IF NOT EXISTS game_db_1 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

这只创建数据库，不创建业务表。业务表需要你根据项目对应的 SQL 或已有数据库准备好，否则 DAO 查询会报表不存在。

### 5.5 直接运行 Go 程序

```bash
go run .
```

含义：编译当前项目并立即运行。

成功后服务监听配置里的端口，默认是：

```text
http://127.0.0.1:7001
```

健康检查：

```bash
curl http://127.0.0.1:7001/health
```

预期返回类似：

```json
{"pid":12345,"status":"ok","timestamp":"2026/05/08 12:00:00","uptime":10}
```

### 5.6 编译后运行

Linux/macOS：

```bash
go build -o main .
./main
```

Windows PowerShell：

```powershell
go build -o main.exe .
.\main.exe
```

---

## 6. Docker 运行

### 6.1 构建镜像

```bash
docker build -f docker/Dockerfile -t server-go:local .
```

解释：

- `docker build`：构建镜像。
- `-f docker/Dockerfile`：指定 Dockerfile。
- `-t server-go:local`：镜像名和标签。
- `.`：构建上下文是项目根目录。

### 6.2 启动依赖

```bash
docker compose -f docker/compose/local.yml up -d
```

启动 MySQL 和 Redis。

### 6.3 运行应用容器

如果应用容器要访问 compose 网络里的 `mysql` 和 `redis`，需要加入同一个网络。

默认网络名是：

```text
server-go-network
```

运行命令示例：

```bash
docker run --rm --network server-go-network -p 7001:7001 server-go:local
```

解释：

- `--rm`：容器停止后自动删除。
- `--network server-go-network`：加入依赖服务所在网络。
- `-p 7001:7001`：把容器 7001 映射到宿主机 7001。
- `server-go:local`：要运行的镜像。

### 6.4 蓝绿部署 compose

项目提供：

```text
docker/compose/traefik.yml
docker/compose/blue.yml
docker/compose/green.yml
```

含义：

- `traefik.yml`：启动 Traefik 网关，监听宿主机 7001。
- `blue.yml`：启动 blue 应用实例。
- `green.yml`：启动 green 应用实例。

启动网关：

```bash
docker compose -f docker/compose/traefik.yml up -d
```

启动 blue：

```bash
docker compose -f docker/compose/blue.yml up -d --build
```

启动 green：

```bash
docker compose -f docker/compose/green.yml up -d --build
```

`--build` 表示启动前重新构建镜像。

蓝绿控制接口需要环境变量：

```text
APP_CONTROL_TOKEN
```

如果是默认值 `PLEASE_CHANGE_ME`，控制接口会返回配置错误。

---

## 7. Makefile 命令

项目根目录有 `Makefile`，它引入了：

```text
hack/hack-cli.mk
hack/hack.mk
```

这些命令依赖 GoFrame CLI：`gf`。

### 7.1 安装或检查 gf CLI

```bash
make cli.install
```

含义：检查本机是否有 `gf` 命令。如果没有，会下载并安装。

注意：`hack/hack-cli.mk` 使用 `wget`、`chmod`、shell 条件语法，Windows 原生 PowerShell 下可能不能直接运行。Windows 建议用 Git Bash、WSL 或手动安装 gf CLI。

### 7.2 构建项目

```bash
make build
```

实际执行：

```bash
gf build -ew
```

含义：使用 GoFrame CLI 编译项目，并按 `hack/config.yaml` 配置输出构建产物。

### 7.3 生成 DAO

```bash
make dao
```

实际执行：

```bash
gf gen dao
```

含义：根据 `hack/config.yaml` 中的数据库连接和表名，生成：

- `internal/dao/*`
- `internal/dao/internal/*`
- `internal/model/do/*`
- `internal/model/entity/*`

`hack/config.yaml` 中配置了这些表：

```text
user,user_res,user_item,user_loginkey,user_bag,user_bag_tp,user_data,user_task,user_online,prf_task,mem_config,sys_gm,_log_login,_log_trace,_log_msg
```

### 7.4 生成 Controller

```bash
make ctrl
```

实际执行：

```bash
gf gen ctrl
```

含义：根据 `api` 目录生成或更新控制器模板。已经手动改过 controller 时，执行前要小心查看 diff，避免覆盖业务代码。

### 7.5 生成 Service

```bash
make service
```

实际执行：

```bash
gf gen service
```

含义：根据 logic 生成 service 接口模板。当前项目已经手写 service，执行后也要检查 diff。

### 7.6 构建 Docker 镜像

```bash
make image
```

实际会：

1. 获取当前 Git commit 短 hash。
2. 如果工作区有未提交改动，tag 后追加 `.dirty`。
3. 执行 `gf docker` 构建镜像。

指定 tag：

```bash
make image TAG=1.0.0
```

构建并推送：

```bash
make image.push TAG=1.0.0
```

### 7.7 Kubernetes 部署

```bash
make deploy _ENV=develop TAG=develop
```

含义：

1. 进入 `manifest/deploy/kustomize/overlays/develop`。
2. 执行 `kustomize build` 生成 YAML。
3. 执行 `kubectl apply` 部署。
4. patch deployment，加一个时间戳 label 触发滚动更新。

依赖：

- `kubectl`
- `kustomize`
- 当前 shell 能访问目标 Kubernetes 集群

---

## 8. API 总览

### 8.1 API 前缀和中间件

以下接口在 `/api` 组内：

- `/api/user/login`
- `/api/user/add_tili`
- `/api/user/add_gold`
- `/api/user/add_diamond`
- `/api/game/online`
- `/api/game/time`
- `/api/bag/get_bag/{chapter}`
- `/api/bag/get_bag_tp/{chapter}`
- `/api/grid/get/{chapter}`

它们会经过：

1. `Sign`
2. `Verify`
3. `Response`

例外：`Verify` 会跳过 `/user/login`，但 `Sign` 不跳过，所以登录接口也需要签名。

以下接口不在 `/api` 组内：

- `/res_version/{key}`：只经过 `Response`。
- `/health/*` 和 `/_internal/control/*`：由控制器直接写 JSON，不经过统一响应。

---

## 9. 签名规则

`middleware.Sign` 的逻辑：

1. 读取所有请求参数：`r.GetMap()`。
2. 从参数 `sign`、Header `x-sign`、Header `x-signature` 中找签名。
3. 调用 `signutil.BuildParams(params)` 拼接待签名字符串。
4. 从配置 `app.keys` 读取密钥列表。
5. 对每个密钥计算 HMAC-SHA256。
6. 只要有一个计算结果等于请求签名，就通过。

### 9.1 待签名字符串

`BuildParams` 会：

1. 排除空 key。
2. 排除 key 为 `sign` 的参数。
3. 排除 nil 值。
4. 按参数名排序。
5. 拼成 `k1=v1&k2=v2`。

例如参数：

```text
uid=1001
openid=test-openid
platform=wx
version=1.0.0
login_key=abc
```

待签名字符串可能是：

```text
login_key=abc&openid=test-openid&platform=wx&uid=1001&version=1.0.0
```

然后用配置里的 secret 做 HMAC-SHA256，输出小写 hex 字符串。

### 9.2 生成签名示例

用 Node.js 生成：

```bash
node -e "const crypto=require('crypto');const s='login_key=abc&openid=test-openid&platform=wx&uid=1001&version=1.0.0';console.log(crypto.createHmac('sha256','SzNg8LgHjEUgTAc4').update(s).digest('hex'))"
```

用 Go 生成：

```go
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

func main() {
    payload := "login_key=abc&openid=test-openid&platform=wx&uid=1001&version=1.0.0"
    secret := "SzNg8LgHjEUgTAc4"
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(payload))
    fmt.Println(hex.EncodeToString(mac.Sum(nil)))
}
```

---

## 10. 登录态校验和防重放

`middleware.Verify` 做这些事：

1. 如果路径以 `/user/login` 结尾，跳过登录态校验。
2. 读取这些参数：
   - `uid`
   - `login_key`
   - `platform`
   - `version`
   - `tick`
   - `sign`
3. 如果缺参数，返回 `code=-1`。
4. 查询 `user_loginkey` 表，确认 `uid` 对应的 key 等于请求里的 `login_key`。
5. 检查 `tick` 和服务器当前时间相差不超过 1800 秒。
6. 用 Redis 检查同一个 `uid + sign` 是否重复调用。
7. 如果没有重复，写入 Redis，过期时间 300 秒。

因此，除了登录接口外，业务接口请求通常要带：

```text
uid
login_key
platform
version
tick
sign
```

---

## 11. 统一响应

`middleware.Response` 会处理 controller 返回值。

如果 controller 返回结构体，例如：

```go
return &apiGame.TimeRes{Now: gtime.TimestampMilli()}, nil
```

中间件会写 JSON，并补充 `code: 0`。

最终类似：

```json
{"code":0,"now":1778240000000}
```

如果 controller 返回 error，中间件返回：

```json
{"code":-1,"msg":"错误信息"}
```

健康检查接口没有经过这个中间件，所以它们的响应由 controller 自己决定。

---

## 12. API 详细说明

### 12.1 登录

路径：

```text
GET/POST /api/user/login
```

参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| uid | int64 | 是 | 用户 ID |
| login_key | string | 否 | 客户端登录 key，会保存到 `user_loginkey` |
| openid | string | 是 | 平台用户标识 |
| platform | string | 是 | 平台，例如 wx、ios、android |
| version | string | 是 | 客户端版本 |
| sign | string | 是 | 签名 |

示例：

```bash
curl "http://127.0.0.1:7001/api/user/login?uid=1001&login_key=abc&openid=test-openid&platform=wx&version=1.0.0&sign=你的签名"
```

业务逻辑：

1. 检查 `openid` 是否为空。
2. 查询 `user` 表是否已有该 `uid`。
3. 如果已有用户，校验 `platform/openid` 是否匹配。
4. 如果没有用户，开启事务：
   - 插入 `user`。
   - 插入初始 `user_res`，默认金币 200、钻石 100、体力 100、等级 1。
5. 异步写 `_log_login` 登录日志。
6. 保存或更新 `user_loginkey`。
7. 查询 `user_data`。
8. 查询是否 GM。
9. 查询 `user_item`。
10. 查询 `user_res`。
11. 查询 `mem_config`。
12. 返回完整登录数据。

### 12.2 增加体力

路径：

```text
GET/POST /api/user/add_tili
```

参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| uid | int64 | 是 | 用户 ID |
| login_key | string | 是 | 登录 key |
| platform | string | 是 | 平台 |
| version | string | 是 | 版本 |
| tick | int64 | 是 | 当前秒级时间戳 |
| sign | string | 是 | 签名 |

业务逻辑：调用 `UpdateTili`，固定增加 50 点体力。

### 12.3 增加金币

路径：

```text
GET/POST /api/user/add_gold
```

业务逻辑：调用 `UpdateGold`，固定增加 50 金币。

### 12.4 增加钻石

路径：

```text
GET/POST /api/user/add_diamond
```

业务逻辑：调用 `UpdateDiamond`，固定增加 50 钻石。

### 12.5 资源更新内部逻辑

`logic/user.updateResField` 用于更新金币、钻石、体力、经验、星星。

执行步骤：

1. 生成锁 key，例如 `update_gold:1001`。
2. 调用 `lock.Lock` 获取 Redis 锁。
3. 查询 `user_res`。
4. 读取旧值。
5. 计算新值：`old + cnt`。
6. 如果新值小于 0，修正为 0。
7. 更新数据库字段。
8. 更新内存中的 `entity.UserRes`。
9. 写资源流水日志 `_log_trace`。
10. 返回最新资源和实际变化值。

### 12.6 记录在线时长

路径：

```text
GET/POST /api/game/online
```

参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| uid | int64 | 是 | 用户 ID |
| seconds | int64 | 是 | 本次增加的在线秒数，不能小于 0 |
| login_key/platform/version/tick/sign | - | 是 | Verify 和 Sign 所需参数 |

业务逻辑：

1. 取当前年月日小时，格式类似 `2026-05-08, 13:00:00`。
2. 查询 `user_online` 是否已有当前用户当前小时记录。
3. 如果有，累加 `tm_online`。
4. 如果没有，插入新记录。
5. 返回当前毫秒时间戳。

### 12.7 获取服务器时间

路径：

```text
GET/POST /api/game/time
```

返回：当前毫秒时间戳。

注意：该接口在 `/api` 组内，因此仍然需要签名和登录态参数。

### 12.8 获取用户背包

路径：

```text
GET/POST /api/bag/get_bag/{chapter}
```

参数：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| chapter | int | 是 | 路径参数，章节 ID |
| uid | int64 | 是 | 用户 ID |
| login_key/platform/version/tick/sign | - | 是 | Verify 和 Sign 所需参数 |

业务逻辑：查询 `user_bag` 表：

```text
where uid = ? and chapter = ?
```

返回字段：

- `uid`
- `chapter`
- `bag`

### 12.9 获取用户背包模板

路径：

```text
GET/POST /api/bag/get_bag_tp/{chapter}
```

业务逻辑：查询 `user_bag_tp` 表。

### 12.10 获取棋盘数据

路径：

```text
GET/POST /api/grid/get/{chapter}
```

业务逻辑：`logic/grid.GetGrid` 并发做三件事：

1. 查询用户背包 `service.Bag().GetUserBag`。
2. 查询用户背包模板 `service.Bag().GetUserBagTp`。
3. 初始化或读取任务 `service.Task().InitTasks`。

它用 `sync.WaitGroup` 等待三个 goroutine 全部完成，用 `sync.Mutex` 保护共享变量。

返回：

- `bag`
- `bag_tp`
- `tasks`

### 12.11 获取资源版本

路径：

```text
GET/POST /res_version/{key}
```

该接口不在 `/api` 下，不需要登录态校验。

业务逻辑：

1. 用 Redis key `res_version.{key}` 做一小时防重复。
2. 调用 `secretutil.CheckSecret(key)` 校验 key。
3. 查询 `mem_config` 表中 `id=50` 的 `value`。
4. 返回资源版本号。

可能返回：

```json
{"code":-1036,"msg":"get_res_version: 不能重复调用"}
```

或：

```json
{"code":-1,"msg":"参数错误"}
```

或：

```json
{"code":0,"ver":"版本号"}
```

### 12.12 健康检查

这些接口不需要签名。

#### `/health/ready`

就绪检查，返回：

```json
{"ok":true}
```

#### `/health`

基础健康检查，返回：

```json
{"status":"ok","pid":123,"uptime":10,"timestamp":"2026/05/08 13:00:00"}
```

#### `/health/detail`

健康详情，额外返回：

- `draining`：是否处于流量切换状态。
- `rejecting`：是否拒绝新请求。

#### `/health/lb`

给 Traefik 使用的负载均衡健康检查。

如果正在流量切换，返回 HTTP 503：

```json
{"status":"draining"}
```

否则返回：

```json
{"status":"ok"}
```

### 12.13 内部流量控制接口

路径：

```text
/_internal/control/traffic-shift
/_internal/control/reject-new-requests
/_internal/control/resume-traffic
```

这些接口用于蓝绿发布。

必须满足：

1. 环境变量 `APP_CONTROL_TOKEN` 已配置。
2. 请求 Header `x-control-token` 等于 `APP_CONTROL_TOKEN`。
3. 请求不能带 `x-forwarded-for`，否则认为不是内部直接访问。

示例：

```bash
curl -H "x-control-token: your-token" http://127.0.0.1:7001/_internal/control/traffic-shift
```

---

## 13. 数据表和 DAO 对应关系

项目使用的主要表：

| 表名 | 作用 |
|---|---|
| user | 用户基础信息 |
| user_res | 用户资源，如金币、钻石、体力、经验 |
| user_item | 用户道具 |
| user_loginkey | 用户当前登录 key |
| user_bag | 用户背包数据 |
| user_bag_tp | 用户背包模板数据 |
| user_data | 用户扩展数据 |
| user_task | 用户任务 |
| user_online | 在线时长记录 |
| prf_task | 任务配置 |
| mem_config | 内存配置/版本配置 |
| sys_gm | GM 用户配置 |
| _log_login | 登录日志 |
| _log_trace | 资源流水日志 |
| _log_msg | 普通消息日志 |

DAO 文件命名和表名基本一致，例如：

```text
internal/dao/user.go              -> dao.User
internal/dao/user_res.go          -> dao.UserRes
internal/dao/internal/user.go     -> user 表字段定义
internal/model/entity/user.go     -> entity.User
internal/model/do/user.go         -> do.User
```

---

## 14. 核心方法讲解

### 14.1 `service.RegisterUser`

位置：`internal/service/user.go`

作用：把某个结构体注册为用户服务实现。

逻辑层调用：

```go
func init() {
    service.RegisterUser(&sUser{})
}
```

因为 `main.go` 空导入了 `logic/user`，所以程序启动时会自动执行这个注册。

### 14.2 `service.User()`

作用：获取用户服务实例。

如果没有注册，会 panic：

```go
panic("service IUser not registered")
```

这能尽早暴露忘记导入 logic 包的问题。

### 14.3 `g.DB().Transaction`

位置：`logic/user.Login`

用于数据库事务。

事务里的多条 SQL 要么全部成功，要么全部失败。

登录中新用户创建时同时插入：

1. `user`
2. `user_res`

如果第二条失败，第一条也会回滚。

### 14.4 `gctx.NeverDone(ctx)`

用于异步日志。

普通请求 `ctx` 会随着请求结束而取消。如果 goroutine 里继续用请求 ctx，可能写日志时 ctx 已经失效。

`NeverDone` 创建一个不会因为请求结束而取消的 context，适合“尽力而为”的异步日志。

### 14.5 Redis 锁 `logic/lock`

资源更新前会加锁，避免同一个用户同一种资源并发更新导致数据错乱。

锁 key 示例：

```text
lock:update_gold:1001
```

释放锁时使用 Lua 脚本：只有 token 匹配才删除锁，避免误删别的请求持有的锁。

### 14.6 `sync.WaitGroup` 和 `sync.Mutex`

位置：`logic/grid.GetGrid`

`WaitGroup` 用来等待多个 goroutine 完成。

`Mutex` 用来保护共享变量：

- `out`
- `firstErr`

如果不加锁，多个 goroutine 同时写同一个变量可能产生数据竞争。

---

## 15. 新手开发流程

### 15.1 修改一个已有接口

例如要改 `/api/game/time`：

1. 看接口定义：`api/game/game.go`。
2. 看 controller：`internal/controller/game.go`。
3. 如果只是返回字段变化，通常改 API Res 和 controller。
4. 如果涉及业务逻辑，改 `internal/logic/game/game.go`。
5. 执行：

```bash
gofmt -w api internal
go test ./...
go vet ./...
```

### 15.2 新增一个 API

以新增 `/api/game/ping` 为例：

1. 在 `api/game/game.go` 新增：
   - `PingReq`
   - `PingRes`
2. 在 `internal/controller/game.go` 给 `cGame` 新增方法：
   - `Ping(ctx, req)`
3. 如果需要业务逻辑：
   - 在 `internal/service/game.go` 的 `IGame` 增加方法。
   - 在 `internal/logic/game/game.go` 实现方法。
4. 格式化和测试。

### 15.3 新增数据库表

1. 先在 MySQL 创建表。
2. 修改 `hack/config.yaml` 的 `tables`，加入表名。
3. 执行：

```bash
make dao
```

4. 检查生成的文件：
   - `internal/dao`
   - `internal/dao/internal`
   - `internal/model/entity`
   - `internal/model/do`
5. 在 logic 中通过 `dao.表对象` 使用。

---

## 16. 常见问题

### 16.1 启动时报数据库连接失败

检查：

1. MySQL 是否启动。
2. 配置文件是否被正确读取。
3. 连接地址是否适合当前运行方式。

容器内使用：

```text
mysql:3306
```

宿主机使用：

```text
127.0.0.1:330
```

### 16.2 Redis 连接失败

检查：

1. Redis 是否启动。
2. 密码是否是 `root`。
3. 地址是否是容器内地址还是宿主机地址。

### 16.3 接口返回 `非法调用`

这是 `Sign` 中间件返回的。通常原因：

- 没有传 `sign`。
- 签名参数排序不一致。
- 使用了错误的 secret。
- 请求参数和计算签名时的参数不一致。

### 16.4 接口返回 `Verify: 参数错误`

通常是业务接口缺少：

```text
uid/login_key/platform/version/tick/sign
```

### 16.5 接口返回 `Verify: 该账号已在其他地方登陆`

说明 `user_loginkey` 表里保存的 key 和请求里的 `login_key` 不一致。

### 16.6 接口返回 `Verify: 时间校验失败`

说明请求的 `tick` 和服务器当前秒级时间戳差距超过 1800 秒。

### 16.7 接口返回 `Verify: 不能重复调用`

说明同一个 `uid + sign` 在 300 秒内重复请求。

---

## 17. 提交前检查清单

修改代码后建议执行：

```bash
gofmt -w .
go test ./...
go vet ./...
```

如果修改了依赖：

```bash
go mod tidy
```

如果修改了数据库表结构：

```bash
make dao
```

如果修改了 Docker 构建：

```bash
docker build -f docker/Dockerfile -t server-go:local .
```

---

## 18. 当前项目的几个注意点

1. `docker/mysql/init/01-init.sql` 只创建数据库，不创建所有业务表。
2. `manifest/deploy/kustomize/base/service.yaml` 中 `targetPort` 是 `8000`，而应用默认监听 `7001`；如果使用 Kubernetes 部署，需要确认这个端口是否符合实际运行配置。
3. `manifest/docker/Dockerfile` 是 GoFrame 模板风格 Dockerfile；当前主要可用 Dockerfile 是 `docker/Dockerfile`。
4. `make` 脚本更适合 Linux/macOS/WSL/Git Bash 环境，Windows PowerShell 下可能需要调整。
5. `/api/game/time` 虽然只是取时间，但因为在 `/api` 组内，仍需要签名和登录态参数。

---

## 19. 最小启动路径总结

如果你只想最快跑起来：

```bash
go mod download
docker compose -f docker/compose/local.yml up -d
go run .
```

然后访问：

```bash
curl http://127.0.0.1:7001/health
```

如果你用 Docker 跑应用：

```bash
docker compose -f docker/compose/local.yml up -d
docker build -f docker/Dockerfile -t server-go:local .
docker run --rm --network server-go-network -p 7001:7001 server-go:local
```

---

## 20. 推荐阅读顺序

如果你是 Go 新手，建议按这个顺序看代码：

1. `main.go`：看程序怎么启动。
2. `internal/cmd/cmd.go`：看路由怎么注册。
3. `api/game/game.go`：看一个简单 API 怎么定义。
4. `internal/controller/game.go`：看 controller 怎么调用 service。
5. `internal/service/game.go`：看 service 接口。
6. `internal/logic/game/game.go`：看业务逻辑。
7. `internal/middleware/sign.go`：看签名校验。
8. `internal/middleware/verify.go`：看登录态校验。
9. `internal/logic/user/user.go`：看完整业务流程。
10. `internal/dao` 和 `internal/model/entity`：看数据库访问模型。

---

## 21. 源码逐文件讲解：server / model / controller / service / logic

这一节把当前项目中手写的 server、model、controller、service、logic 代码集中贴出，并按文件解释。阅读时建议按调用链看：

```text
main.go
  -> internal/cmd/cmd.go
  -> internal/controller/*.go
  -> internal/service/*.go
  -> internal/logic/**/*.go
  -> internal/dao + internal/model/entity
```

### 21.1 Server 入口

#### `main.go`

```go
package main

import (
	_ "server_go/internal/packed"

	_ "server_go/internal/logic/bag"
	_ "server_go/internal/logic/game"
	_ "server_go/internal/logic/grid"
	_ "server_go/internal/logic/other"
	_ "server_go/internal/logic/task"
	_ "server_go/internal/logic/user"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"server_go/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
```

讲解：

- `package main` 表示这是可执行程序入口包。
- `_ "server_go/internal/packed"` 是空导入，只执行该包的 `init`，通常用于 GoFrame 资源打包。
- `_ "server_go/internal/logic/..."` 也是空导入，目的是触发各 logic 包里的 `init()`，完成 service 注册。
- `_ "github.com/gogf/gf/contrib/drivers/mysql/v2"` 注册 MySQL 驱动。
- `_ "github.com/gogf/gf/contrib/nosql/redis/v2"` 注册 Redis 支持。
- `gctx.GetInitCtx()` 创建初始化上下文。
- `cmd.Main.Run(...)` 启动 `internal/cmd/cmd.go` 里定义的命令。

#### `internal/cmd/cmd.go`

```go
package cmd

import (
	"context"

	"server_go/internal/controller"
	"server_go/internal/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// Game API routes
			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Middleware(
					middleware.Sign,
					middleware.Verify,
					middleware.Response,
				)
				group.Bind(
					controller.User,
					controller.Game,
					controller.Bag,
					controller.Grid,
				)
			})

			// Other routes (no sign/verify)
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.Response)
				group.Bind(
					controller.Other,
				)
			})

			// Health + internal control routes (no middleware)
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Bind(
					controller.Health,
				)
			})

			s.Run()
			return nil
		},
	}
)
```

讲解：

- `gcmd.Command` 是 GoFrame 命令对象。
- `Func` 是命令执行函数，程序启动后会进入这里。
- `g.Server()` 获取默认 HTTP Server。
- `/api` 路由组绑定用户、游戏、背包、棋盘接口，并使用 `Sign`、`Verify`、`Response` 三个中间件。
- `/` 路由组绑定 `Other`，用于 `/res_version/{key}`，只走统一响应，不走签名和登录校验。
- 第二个 `/` 路由组绑定健康检查和内部控制接口，不走统一响应，由 controller 手动写 JSON。
- `s.Run()` 阻塞运行 HTTP Server。

### 21.2 Model 层

#### `internal/model/model.go`

```go
package model

import (
	"server_go/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
)

// --- User ---

type LoginInput struct {
	Uid      int64  `json:"uid"`
	LoginKey string `json:"login_key"`
	Openid   string `json:"openid"`
	Platform string `json:"platform"`
	Version  string `json:"version"`
}

type LoginOutput struct {
	Uid    int64           `json:"uid"`
	Newbie int             `json:"newbie"`
	User   interface{}     `json:"user"`
	Res    *entity.UserRes `json:"res,omitempty"`
	Datas  gdb.Result      `json:"datas,omitempty"`
	Items  gdb.Result      `json:"items,omitempty"`
	Config gdb.Result      `json:"config,omitempty"`
	Gm     int             `json:"gm"`
}

type UpdateResInput struct {
	Uid    int64       `json:"uid"`
	Items  interface{} `json:"items"`
	Reason string      `json:"reason"`
}

type UpdateFieldInput struct {
	Uid    int64  `json:"uid"`
	Cnt    int64  `json:"cnt"`
	Reason string `json:"reason"`
}

type UpdateFieldOutput struct {
	Res      *entity.UserRes `json:"res"`
	AddValue int64           `json:"add_value"`
}

// --- Game ---

type OnlineInput struct {
	Uid     int64 `json:"uid"`
	Seconds int64 `json:"seconds"`
}

// --- Bag ---

type BagInput struct {
	Uid     int64 `json:"uid"`
	Chapter int   `json:"chapter"`
}

type BagOutput struct {
	Uid     int64      `json:"uid"`
	Chapter int        `json:"chapter"`
	Bag     gdb.Result `json:"bag"`
}

// --- Grid ---

type GridOutput struct {
	Bag   *BagOutput  `json:"bag,omitempty"`
	BagTp *BagOutput  `json:"bag_tp,omitempty"`
	Tasks interface{} `json:"tasks,omitempty"`
}

// --- Other ---

type ResVersionOutput struct {
	Code int    `json:"code"`
	Ver  string `json:"ver,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
```

讲解：

- `model` 包是内部业务层 DTO，不直接负责 HTTP 路由。
- `LoginInput` 是用户登录业务入参，由 controller 从 `api.LoginReq` 转换而来。
- `LoginOutput` 是用户登录业务出参，包含用户、资源、数据、道具、配置、GM 状态。
- `UpdateFieldInput` 是更新资源字段的通用入参，`Cnt` 表示变化值。
- `UpdateFieldOutput` 返回更新后的资源和实际变化值。
- `OnlineInput` 用于记录在线时长。
- `BagInput` 同时用于背包和棋盘查询，因为棋盘查询也需要 `uid + chapter`。
- `BagOutput` 中 `gdb.Result` 是 GoFrame 数据库查询结果列表。
- `GridOutput` 聚合背包、背包模板和任务数据。
- `ResVersionOutput` 是资源版本业务输出。

### 21.3 Controller 层

#### `internal/controller/bag.go`

```go
package controller

import (
	"context"

	apiBag "server_go/api/bag"
	"server_go/internal/model"
	"server_go/internal/service"
)

var Bag = &cBag{}

type cBag struct{}

func (c *cBag) GetBag(ctx context.Context, req *apiBag.GetBagReq) (res *apiBag.GetBagRes, err error) {
	out, err := service.Bag().GetUserBag(ctx, &model.BagInput{Uid: req.Uid, Chapter: req.Chapter})
	if err != nil {
		return nil, err
	}
	return (*apiBag.GetBagRes)(toBagRes(out)), nil
}

func (c *cBag) GetBagTp(ctx context.Context, req *apiBag.GetBagTpReq) (res *apiBag.GetBagTpRes, err error) {
	out, err := service.Bag().GetUserBagTp(ctx, &model.BagInput{Uid: req.Uid, Chapter: req.Chapter})
	if err != nil {
		return nil, err
	}
	return (*apiBag.GetBagTpRes)(toBagRes(out)), nil
}

func toBagRes(out *model.BagOutput) *apiBag.BagRes {
	if out == nil {
		return nil
	}
	return &apiBag.BagRes{
		Uid:     out.Uid,
		Chapter: out.Chapter,
		Bag:     out.Bag,
	}
}
```

讲解：

- `var Bag = &cBag{}` 暴露给路由绑定使用。
- `GetBag` 接收 HTTP 请求对象 `apiBag.GetBagReq`。
- controller 不直接查数据库，而是调用 `service.Bag().GetUserBag`。
- `toBagRes` 把内部业务输出 `model.BagOutput` 转成 API 响应 `apiBag.BagRes`。
- `GetBagTp` 逻辑类似，只是查询背包模板表。

#### `internal/controller/game.go`

```go
package controller

import (
	"context"

	apiGame "server_go/api/game"
	"server_go/internal/model"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

var Game = &cGame{}

type cGame struct{}

func (c *cGame) Online(ctx context.Context, req *apiGame.OnlineReq) (res *apiGame.OnlineRes, err error) {
	err = service.Game().Online(ctx, &model.OnlineInput{
		Uid: req.Uid, Seconds: req.Seconds,
	})
	if err != nil {
		return nil, err
	}
	return &apiGame.OnlineRes{Now: gtime.TimestampMilli()}, nil
}

func (c *cGame) Time(ctx context.Context, req *apiGame.TimeReq) (res *apiGame.TimeRes, err error) {
	return &apiGame.TimeRes{Now: gtime.TimestampMilli()}, nil
}
```

讲解：

- `Online` 记录用户在线时长，成功后返回服务器当前毫秒时间戳。
- `Time` 不走 service，直接返回当前毫秒时间戳。
- `gtime.TimestampMilli()` 是 GoFrame 获取当前毫秒时间戳的方法。

#### `internal/controller/grid.go`

```go
package controller

import (
	"context"

	apiGrid "server_go/api/grid"
	"server_go/internal/model"
	"server_go/internal/service"
)

var Grid = &cGrid{}

type cGrid struct{}

func (c *cGrid) GetGrid(ctx context.Context, req *apiGrid.GetGridReq) (res *apiGrid.GetGridRes, err error) {
	out, err := service.Grid().GetGrid(ctx, &model.BagInput{Uid: req.Uid, Chapter: req.Chapter})
	if err != nil {
		return nil, err
	}
	return &apiGrid.GetGridRes{
		Bag:   toBagRes(out.Bag),
		BagTp: toBagRes(out.BagTp),
		Tasks: out.Tasks,
	}, nil
}
```

讲解：

- `GetGrid` 把棋盘接口请求转换为 `model.BagInput`。
- 实际业务由 `service.Grid().GetGrid` 完成。
- 返回值复用 `toBagRes`，把 `model.BagOutput` 转成 `apiBag.BagRes`。

#### `internal/controller/other.go`

```go
package controller

import (
	"context"

	apiOther "server_go/api/other"
	"server_go/internal/service"
)

var Other = &cOther{}

type cOther struct{}

func (c *cOther) ResVersion(ctx context.Context, req *apiOther.ResVersionReq) (res *apiOther.ResVersionRes, err error) {
	out, err := service.Other().GetResVersion(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &apiOther.ResVersionRes{
		Code: out.Code,
		Ver:  out.Ver,
		Msg:  out.Msg,
	}, nil
}
```

讲解：

- `ResVersion` 处理 `/res_version/{key}`。
- 这个接口不在 `/api` 路由组内，不走签名和登录校验。
- controller 只负责转发 key 和转换响应。

#### `internal/controller/user.go`

```go
package controller

import (
	"context"

	apiUser "server_go/api/user"
	"server_go/internal/model"
	"server_go/internal/service"
)

var User = &cUser{}

type cUser struct{}

func (c *cUser) Login(ctx context.Context, req *apiUser.LoginReq) (res *apiUser.LoginRes, err error) {
	out, err := service.User().Login(ctx, &model.LoginInput{
		Uid:      req.Uid,
		LoginKey: req.LoginKey,
		Openid:   req.Openid,
		Platform: req.Platform,
		Version:  req.Version,
	})
	if err != nil {
		return nil, err
	}
	return &apiUser.LoginRes{
		Uid:    out.Uid,
		Newbie: out.Newbie,
		User:   out.User,
		Res:    out.Res,
		Datas:  out.Datas,
		Items:  out.Items,
		Config: out.Config,
		Gm:     out.Gm,
	}, nil
}

func (c *cUser) AddTili(ctx context.Context, req *apiUser.AddTiliReq) (res *apiUser.AddTiliRes, err error) {
	out, err := service.User().UpdateTili(ctx, &model.UpdateFieldInput{
		Uid: req.Uid, Cnt: 50, Reason: "测试增加体力",
	})
	if err != nil {
		return nil, err
	}
	return (*apiUser.AddTiliRes)(toUpdateFieldRes(out)), nil
}

func (c *cUser) AddGold(ctx context.Context, req *apiUser.AddGoldReq) (res *apiUser.AddGoldRes, err error) {
	out, err := service.User().UpdateGold(ctx, &model.UpdateFieldInput{
		Uid: req.Uid, Cnt: 50, Reason: "测试增加金币",
	})
	if err != nil {
		return nil, err
	}
	return (*apiUser.AddGoldRes)(toUpdateFieldRes(out)), nil
}

func (c *cUser) AddDiamond(ctx context.Context, req *apiUser.AddDiamondReq) (res *apiUser.AddDiamondRes, err error) {
	out, err := service.User().UpdateDiamond(ctx, &model.UpdateFieldInput{
		Uid: req.Uid, Cnt: 50, Reason: "测试增加钻石",
	})
	if err != nil {
		return nil, err
	}
	return (*apiUser.AddDiamondRes)(toUpdateFieldRes(out)), nil
}

func toUpdateFieldRes(out *model.UpdateFieldOutput) *apiUser.UpdateFieldRes {
	if out == nil {
		return nil
	}
	return &apiUser.UpdateFieldRes{
		Res:      out.Res,
		AddValue: out.AddValue,
	}
}
```

讲解：

- `Login` 把 HTTP 登录请求转换成内部 `model.LoginInput`。
- `AddTili`、`AddGold`、`AddDiamond` 是测试资源增加接口，固定增加 50。
- `Reason` 用于资源流水日志，方便定位资源变化来源。
- `toUpdateFieldRes` 是内部输出到 API 输出的转换函数。

#### `internal/controller/health.go`

```go
package controller

import (
	"context"
	"os"
	"sync"
	"time"

	apiHealth "server_go/api/health"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Health = &cHealth{}

type cHealth struct{}

var startTime = time.Now()

var drainState = &drainStateManager{}

type drainStateManager struct {
	mu                   sync.RWMutex
	draining             bool
	rejectingNewRequests bool
}

func (d *drainStateManager) IsTrafficShift() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.draining
}

func (d *drainStateManager) IsRejecting() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.rejectingNewRequests
}

func (d *drainStateManager) StartTrafficShift() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.draining = true
	d.rejectingNewRequests = false
}

func (d *drainStateManager) StartRejectNew() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.draining = true
	d.rejectingNewRequests = true
}

func (d *drainStateManager) Resume() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.draining = false
	d.rejectingNewRequests = false
}

func (c *cHealth) Ready(ctx context.Context, req *apiHealth.ReadyReq) (res *apiHealth.ReadyRes, err error) {
	ghttp.RequestFromCtx(ctx).Response.WriteJson(g.Map{"ok": true})
	return
}

func (c *cHealth) Health(ctx context.Context, req *apiHealth.HealthReq) (res *apiHealth.HealthRes, err error) {
	ghttp.RequestFromCtx(ctx).Response.WriteJson(g.Map{
		"status":    "ok",
		"pid":       os.Getpid(),
		"uptime":    int(time.Since(startTime).Seconds()),
		"timestamp": time.Now().Format("2006/01/02 15:04:05"),
	})
	return
}

func (c *cHealth) HealthDetail(ctx context.Context, req *apiHealth.HealthDetailReq) (res *apiHealth.HealthDetailRes, err error) {
	ghttp.RequestFromCtx(ctx).Response.WriteJson(g.Map{
		"status":    "ok",
		"pid":       os.Getpid(),
		"uptime":    int(time.Since(startTime).Seconds()),
		"timestamp": time.Now().Format("2006/01/02 15:04:05"),
		"draining":  drainState.IsTrafficShift(),
		"rejecting": drainState.IsRejecting(),
	})
	return
}

func (c *cHealth) HealthLb(ctx context.Context, req *apiHealth.HealthLbReq) (res *apiHealth.HealthLbRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if drainState.IsTrafficShift() {
		r.Response.Status = 503
		r.Response.WriteJson(g.Map{"status": "draining"})
	} else {
		r.Response.WriteJson(g.Map{"status": "ok"})
	}
	return
}

func (c *cHealth) TrafficShift(ctx context.Context, req *apiHealth.TrafficShiftReq) (res *apiHealth.TrafficShiftRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if !ensureInternalAccess(r) {
		return
	}
	drainState.StartTrafficShift()
	r.Response.WriteJson(g.Map{"ok": true, "state": "traffic-shift"})
	return
}

func (c *cHealth) RejectNew(ctx context.Context, req *apiHealth.RejectNewReq) (res *apiHealth.RejectNewRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if !ensureInternalAccess(r) {
		return
	}
	drainState.StartRejectNew()
	r.Response.WriteJson(g.Map{"ok": true, "state": "reject-new-requests"})
	return
}

func (c *cHealth) ResumeTraffic(ctx context.Context, req *apiHealth.ResumeTrafficReq) (res *apiHealth.ResumeTrafficRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if !ensureInternalAccess(r) {
		return
	}
	drainState.Resume()
	r.Response.WriteJson(g.Map{"ok": true, "state": "resume-traffic"})
	return
}

func ensureInternalAccess(r *ghttp.Request) bool {
	expected := os.Getenv("APP_CONTROL_TOKEN")
	if expected == "" || expected == "PLEASE_CHANGE_ME" {
		r.Response.Status = 500
		r.Response.WriteJson(g.Map{"ok": false, "msg": "APP_CONTROL_TOKEN not configured"})
		return false
	}
	forwarded := r.GetHeader("x-forwarded-for")
	if forwarded != "" {
		r.Response.Status = 404
		r.Response.WriteJson(g.Map{"ok": false})
		return false
	}
	received := r.GetHeader("x-control-token")
	if received != expected {
		r.Response.Status = 404
		r.Response.WriteJson(g.Map{"ok": false})
		return false
	}
	return true
}
```

讲解：

- 健康检查接口自己写响应，不依赖 `middleware.Response`。
- `startTime` 用于计算进程已运行秒数。
- `drainStateManager` 管理蓝绿发布时的流量状态。
- `sync.RWMutex` 保证并发读写 `draining/rejectingNewRequests` 安全。
- `HealthLb` 给 Traefik 调用，如果正在切流则返回 503，让网关停止转发流量。
- `ensureInternalAccess` 要求配置 `APP_CONTROL_TOKEN`，并要求请求头 `x-control-token` 匹配。

### 21.4 Service 层

#### `internal/service/bag.go`

```go
package service

import (
	"context"

	"server_go/internal/model"
)

type IBag interface {
	GetUserBag(ctx context.Context, in *model.BagInput) (*model.BagOutput, error)
	GetUserBagTp(ctx context.Context, in *model.BagInput) (*model.BagOutput, error)
}

var localBag IBag

func Bag() IBag {
	if localBag == nil {
		panic("service IBag not registered")
	}
	return localBag
}

func RegisterBag(s IBag) {
	localBag = s
}
```

讲解：

- `IBag` 定义背包服务必须实现的方法。
- `localBag` 保存真正实现。
- `RegisterBag` 由 `logic/bag` 的 `init()` 调用。
- `Bag()` 给 controller 或其他 logic 获取服务。

#### `internal/service/game.go`

```go
package service

import (
	"context"

	"server_go/internal/model"
)

type IGame interface {
	Online(ctx context.Context, in *model.OnlineInput) error
}

var localGame IGame

func Game() IGame {
	if localGame == nil {
		panic("service IGame not registered")
	}
	return localGame
}

func RegisterGame(s IGame) {
	localGame = s
}
```

讲解：

- `IGame` 当前只定义在线时长记录。
- `Game()` 是统一入口，屏蔽具体实现结构体 `sGame`。

#### `internal/service/grid.go`

```go
package service

import (
	"context"

	"server_go/internal/model"
)

type IGrid interface {
	GetGrid(ctx context.Context, in *model.BagInput) (*model.GridOutput, error)
}

var localGrid IGrid

func Grid() IGrid {
	if localGrid == nil {
		panic("service IGrid not registered")
	}
	return localGrid
}

func RegisterGrid(s IGrid) {
	localGrid = s
}
```

讲解：

- `IGrid` 定义棋盘聚合查询。
- 入参复用 `model.BagInput`，因为只需要 `uid/chapter`。

#### `internal/service/other.go`

```go
package service

import (
	"context"

	"server_go/internal/model"
)

type IOther interface {
	GetResVersion(ctx context.Context, key string) (*model.ResVersionOutput, error)
}

var localOther IOther

func Other() IOther {
	if localOther == nil {
		panic("service IOther not registered")
	}
	return localOther
}

func RegisterOther(s IOther) {
	localOther = s
}
```

讲解：

- `IOther` 当前只负责资源版本查询。
- 该服务被 `controller.Other` 调用。

#### `internal/service/task.go`

```go
package service

import (
	"context"
)

type ITask interface {
	InitTasks(ctx context.Context, uid int64) ([]map[string]interface{}, error)
}

var localTask ITask

func Task() ITask {
	if localTask == nil {
		panic("service ITask not registered")
	}
	return localTask
}

func RegisterTask(s ITask) {
	localTask = s
}
```

讲解：

- `ITask` 提供任务初始化和读取。
- `logic/grid` 会调用 `service.Task().InitTasks` 聚合任务数据。

#### `internal/service/user.go`

```go
package service

import (
	"context"

	"server_go/internal/model"
	"server_go/internal/model/entity"
)

type IUser interface {
	Login(ctx context.Context, in *model.LoginInput) (*model.LoginOutput, error)
	UpdateDiamond(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
	UpdateGold(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
	UpdateTili(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
	UpdateExp(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
	UpdateStar(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
	GetUser(ctx context.Context, uid int64) (*entity.User, error)
	GetUserRes(ctx context.Context, uid int64) (*entity.UserRes, error)
}

var localUser IUser

func User() IUser {
	if localUser == nil {
		panic("service IUser not registered")
	}
	return localUser
}

func RegisterUser(s IUser) {
	localUser = s
}
```

讲解：

- `IUser` 是用户业务接口。
- `Login` 是登录主流程。
- `UpdateDiamond/Gold/Tili/Exp/Star` 是资源更新方法。
- `GetUser/GetUserRes` 是用户基础数据和资源数据读取方法。
- `localUser == nil` 时 panic，能尽早暴露 logic 包没有注册的问题。

### 21.5 Logic 层

#### `internal/logic/bag/bag.go`

```go
package bag

import (
	"context"

	"server_go/internal/dao"
	"server_go/internal/model"
	"server_go/internal/service"
)

type sBag struct{}

func init() {
	service.RegisterBag(&sBag{})
}

func (s *sBag) GetUserBag(ctx context.Context, in *model.BagInput) (*model.BagOutput, error) {
	rows, err := dao.UserBag.Ctx(ctx).Where("uid", in.Uid).Where("chapter", in.Chapter).All()
	if err != nil {
		return nil, err
	}
	return &model.BagOutput{Uid: in.Uid, Chapter: in.Chapter, Bag: rows}, nil
}

func (s *sBag) GetUserBagTp(ctx context.Context, in *model.BagInput) (*model.BagOutput, error) {
	rows, err := dao.UserBagTp.Ctx(ctx).Where("uid", in.Uid).Where("chapter", in.Chapter).All()
	if err != nil {
		return nil, err
	}
	return &model.BagOutput{Uid: in.Uid, Chapter: in.Chapter, Bag: rows}, nil
}
```

讲解：

- `sBag` 是背包服务实现。
- `init()` 把 `sBag` 注册到 service 层。
- `GetUserBag` 查询 `user_bag`。
- `GetUserBagTp` 查询 `user_bag_tp`。
- `All()` 返回多行结果。

#### `internal/logic/game/game.go`

```go
package game

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/model"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sGame struct{}

func init() {
	service.RegisterGame(&sGame{})
}

func (s *sGame) Online(ctx context.Context, in *model.OnlineInput) error {
	now := gtime.Now()
	dayStr := fmt.Sprintf("%d-%02d-%02d, %02d:00:00", now.Year(), now.Month(), now.Day(), now.Hour())

	row, err := dao.UserOnline.Ctx(ctx).Where("uid", in.Uid).Where("day", dayStr).One()
	if err != nil {
		return err
	}

	nowTime := gtime.Now().Format("Y-m-d H:i:s")
	seconds := in.Seconds
	if !row.IsEmpty() {
		seconds += row["tm_online"].Int64()
		_, err = dao.UserOnline.Ctx(ctx).Where("uid", in.Uid).Where("day", dayStr).
			Data(g.Map{"tm_online": seconds, "tm_update": nowTime}).Update()
	} else {
		_, err = dao.UserOnline.Ctx(ctx).Data(g.Map{
			"uid": in.Uid, "day": dayStr, "tm_online": seconds, "tm_update": nowTime,
		}).Insert()
	}
	return err
}
```

讲解：

- 按“用户 + 小时”统计在线时长。
- `dayStr` 实际包含年月日和小时。
- 如果当前小时已有记录，就累加 `tm_online`。
- 如果没有记录，就插入一条。
- `g.Map` 是 GoFrame 的 `map[string]interface{}` 简写。

#### `internal/logic/grid/grid.go`

```go
package grid

import (
	"context"
	"sync"

	"server_go/internal/model"
	"server_go/internal/service"
)

type sGrid struct{}

func init() {
	service.RegisterGrid(&sGrid{})
}

func (s *sGrid) GetGrid(ctx context.Context, in *model.BagInput) (*model.GridOutput, error) {
	out := &model.GridOutput{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	wg.Add(3)

	go func() {
		defer wg.Done()
		result, err := service.Bag().GetUserBag(ctx, in)
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		out.Bag = result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		result, err := service.Bag().GetUserBagTp(ctx, in)
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		out.BagTp = result
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		tasks, err := service.Task().InitTasks(ctx, in.Uid)
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
			return
		}
		mu.Lock()
		out.Tasks = tasks
		mu.Unlock()
	}()

	wg.Wait()
	return out, firstErr
}
```

讲解：

- `GetGrid` 同时查询背包、背包模板、任务。
- 三个 goroutine 并发执行，减少总等待时间。
- `WaitGroup` 等待三个 goroutine 结束。
- `Mutex` 保护 `out` 和 `firstErr`，避免并发写冲突。
- `firstErr` 保存第一个出现的错误。

#### `internal/logic/other/other.go`

```go
package other

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/utility/secretutil"

	"github.com/gogf/gf/v2/frame/g"
)

type sOther struct{}

func init() {
	service.RegisterOther(&sOther{})
}

func (s *sOther) GetResVersion(ctx context.Context, key string) (*model.ResVersionOutput, error) {
	redis := g.Redis()
	rkey := fmt.Sprintf("res_version.%s", key)

	exists, err := redis.Do(ctx, "EXISTS", rkey)
	if err != nil {
		return nil, err
	}
	if exists.Int() > 0 {
		return &model.ResVersionOutput{Code: -1036, Msg: "get_res_version: 不能重复调用"}, nil
	}
	_, _ = redis.Do(ctx, "SET", rkey, "1", "EX", 3600)

	if !secretutil.CheckSecret(key) {
		return &model.ResVersionOutput{Code: -1, Msg: "参数错误"}, nil
	}

	ver, err := dao.MemConfig.Ctx(ctx).Where("id", 50).Value("value")
	if err != nil {
		return nil, err
	}

	return &model.ResVersionOutput{Code: 0, Ver: ver.String()}, nil
}
```

讲解：

- 用 Redis key `res_version.{key}` 做一小时防重复。
- `EXISTS` 判断是否已调用。
- `SET ... EX 3600` 写入一小时过期的标记。
- `secretutil.CheckSecret` 校验 key 合法性。
- 从 `mem_config` 表读取 `id=50` 的版本值。

#### `internal/logic/task/task.go`

```go
package task

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/logic/user"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sTask struct{}

func init() {
	service.RegisterTask(&sTask{})
}

func (s *sTask) InitTasks(ctx context.Context, uid int64) ([]map[string]interface{}, error) {
	taskConf, err := dao.UserData.Ctx(ctx).Where("uid", uid).Where("key", "task_conf").Value("value")
	if err != nil {
		return nil, err
	}

	confStr := taskConf.String()
	if confStr == "" {
		_, _ = dao.UserData.Ctx(ctx).Data(g.Map{"uid": uid, "key": "task_conf", "value": "4"}).Insert()
		confStr = "4"
	}

	serList := user.PickNumbers(confStr)
	var arr []map[string]interface{}

	for _, ser := range serList {
		task, e := getOneTask(ctx, uid, ser)
		if e != nil {
			return nil, e
		}
		if task != nil {
			task["uid"] = uid
			arr = append(arr, task)
		}
	}

	return arr, nil
}

func getOneTask(ctx context.Context, uid int64, ser int) (map[string]interface{}, error) {
	minMax, err := getTaskSerMinMax(ctx, ser)
	if err != nil {
		return nil, err
	}
	if minMax == nil {
		return nil, fmt.Errorf("用户%d的任务类型%d没有数据", uid, ser)
	}
	minId, maxId := minMax[0], minMax[1]

	row, err := dao.UserTask.Ctx(ctx).Where("uid", uid).
		WhereGTE("taskid", minId).WhereLTE("taskid", maxId).
		Limit(1).One()
	if err != nil {
		return nil, err
	}

	var taskId int
	needClear := false
	nowSec := gtime.Timestamp()

	if row.IsEmpty() {
		taskId = minId
		_, err = dao.UserTask.Ctx(ctx).Data(g.Map{
			"uid": uid, "taskid": taskId, "addtm": nowSec, "done": 0, "donetm": 0,
		}).Insert()
		if err != nil {
			return nil, err
		}
	} else if row["done"].Int() != 0 {
		taskId = row["taskid"].Int()
		if taskId >= maxId {
			v, e := dao.PrfTask.Ctx(ctx).Where("ser", ser).Where("start_loop", 1).
				OrderAsc("id").Limit(1).Value("id")
			if e != nil {
				return nil, e
			}
			taskId = v.Int()
		} else {
			v, e := dao.PrfTask.Ctx(ctx).Where("ser", ser).WhereGT("id", taskId).
				OrderAsc("id").Limit(1).Value("id")
			if e != nil {
				return nil, e
			}
			taskId = v.Int()
		}
		needClear = true
		_, err = dao.UserTask.Ctx(ctx).Data(g.Map{
			"uid": uid, "taskid": taskId, "addtm": nowSec, "done": 0, "donetm": 0,
		}).Insert()
		if err != nil {
			return nil, err
		}
	} else {
		taskId = row["taskid"].Int()
	}

	if needClear {
		_, _ = dao.UserTask.Ctx(ctx).Where("uid", uid).
			WhereGTE("taskid", minId).WhereLTE("taskid", maxId).
			Where("done", 1).Delete()
	}

	if taskId == 0 {
		return nil, nil
	}

	taskRow, err := dao.PrfTask.Ctx(ctx).Where("id", taskId).One()
	if err != nil {
		return nil, err
	}
	if taskRow.IsEmpty() {
		return nil, nil
	}
	return taskRow.Map(), nil
}

func getTaskSerMinMax(ctx context.Context, ser int) ([]int, error) {
	minVal, err := dao.PrfTask.Ctx(ctx).Where("ser", ser).Min("id")
	if err != nil {
		return nil, err
	}
	if minVal == 0 {
		return nil, nil
	}
	maxVal, err := dao.PrfTask.Ctx(ctx).Where("ser", ser).Max("id")
	if err != nil {
		return nil, err
	}
	if maxVal == 0 {
		return nil, nil
	}
	return []int{int(minVal), int(maxVal)}, nil
}
```

讲解：

- `InitTasks` 读取用户任务配置 `user_data.key=task_conf`。
- 如果没有配置，默认插入 `4`。
- `user.PickNumbers` 从配置字符串中提取任务类型编号。
- `getOneTask` 根据任务类型 `ser` 找当前任务。
- 如果用户没有该类型任务，从配置表最小任务 ID 开始。
- 如果当前任务已完成，则进入下一个任务；到最大 ID 后找 `start_loop=1` 的循环起点。
- `needClear` 为 true 时，清理该任务类型下已完成旧任务。

#### `internal/logic/lock/lock.go`

```go
package lock

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	acquireTimeoutMs = 1000
	retryBaseMs      = 20
	retryMaxMs       = 200
	lockTTLMs        = 30000
)

func Lock(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("[Lock] key is required")
	}
	redis := g.Redis()
	redisKey := "lock:" + key
	token := fmt.Sprintf("%d:%d:%d", os.Getpid(), time.Now().UnixNano(), rand.Int63())

	deadline := time.Now().Add(time.Duration(acquireTimeoutMs) * time.Millisecond)
	retryCount := 0

	for time.Now().Before(deadline) {
		ok, err := tryAcquire(ctx, redis, redisKey, token)
		if err != nil {
			return "", err
		}
		if ok {
			return token, nil
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			break
		}
		cap := math.Min(float64(retryMaxMs), float64(retryBaseMs)*math.Pow(2, float64(retryCount)))
		sleepMs := int(math.Min(float64(remaining.Milliseconds()), float64(rand.Intn(int(cap)+1))))
		retryCount++
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
	return "", nil
}

func Unlock(ctx context.Context, key, token string) error {
	if key == "" || token == "" {
		return nil
	}
	redis := g.Redis()
	redisKey := "lock:" + key
	script := `if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end`
	_, err := redis.Do(ctx, "EVAL", script, 1, redisKey, token)
	return err
}

func tryAcquire(ctx context.Context, redis *gredis.Redis, key, token string) (bool, error) {
	res, err := redis.Do(ctx, "SET", key, token, "PX", lockTTLMs, "NX")
	if err != nil {
		return false, err
	}
	return res.String() == "OK", nil
}
```

讲解：

- `Lock` 用 Redis `SET key token PX ttl NX` 获取锁。
- `NX` 表示只有 key 不存在时才设置成功。
- `PX lockTTLMs` 表示锁自动过期，当前是 30000 毫秒。
- `token` 用进程 ID、纳秒时间、随机数组成，用于标识锁持有者。
- 获取失败会短暂 sleep 后重试，最多尝试 1000 毫秒。
- `Unlock` 用 Lua 脚本保证“只有 token 匹配才删除锁”，避免误删别人持有的锁。

#### `internal/logic/gamelog/gamelog.go`

```go
package gamelog

import (
	"context"
	"math"

	"server_go/internal/dao"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// TraceRes records a resource change asynchronously.
func TraceRes(ctx context.Context, uid int64, old, now int64, resName, reason string) {
	if uid == 0 {
		return
	}
	num := now - old
	if num == 0 {
		return
	}
	if num > 0 {
		resName = "+" + resName
	} else {
		resName = "-" + resName
	}
	absNum := int64(math.Abs(float64(num)))
	bgCtx := gctx.NeverDone(ctx)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.Log().Errorf(context.Background(), "TraceRes panic: %v", r)
			}
		}()
		_, _ = dao.LogTrace.Ctx(bgCtx).Data(g.Map{
			"uid": uid, "type": resName, "num": absNum,
			"before": old, "after": now, "reason": reason,
		}).Insert()
	}()
}

// Log records a message asynchronously.
func Log(ctx context.Context, uid int64, msg string) {
	bgCtx := gctx.NeverDone(ctx)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.Log().Errorf(context.Background(), "Log panic: %v", r)
			}
		}()
		_, _ = dao.LogMsg.Ctx(bgCtx).Data(g.Map{"uid": uid, "msg": msg}).Insert()
	}()
}
```

讲解：

- `TraceRes` 异步记录资源变化流水。
- `num := now - old` 计算变化量。
- 增加资源时类型前加 `+`，减少资源时类型前加 `-`。
- `gctx.NeverDone(ctx)` 避免请求结束后异步日志 context 被取消。
- goroutine 内部使用 `recover`，避免日志写入 panic 影响主流程。
- `Log` 用于写普通消息日志 `_log_msg`。

#### `internal/logic/user/user.go`

`user.go` 是业务最重的文件，包含登录、资源更新、用户资源读取和工具函数。源码较长，下面完整贴出。

```go
package user

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"server_go/internal/consts"
	"server_go/internal/dao"
	"server_go/internal/logic/gamelog"
	"server_go/internal/logic/lock"
	"server_go/internal/model"
	"server_go/internal/model/entity"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

type sUser struct{}

func init() {
	service.RegisterUser(&sUser{})
}

func (s *sUser) Login(ctx context.Context, in *model.LoginInput) (*model.LoginOutput, error) {
	if in.Openid == "" {
		return nil, fmt.Errorf("参数错误: openid 必填")
	}

	out := &model.LoginOutput{Uid: in.Uid}

	var user *entity.User
	err := dao.User.Ctx(ctx).Where("uid", in.Uid).Scan(&user)
	if err != nil {
		return nil, err
	}

	if user != nil {
		if user.Platform != in.Platform || user.Openid != in.Openid {
			return nil, fmt.Errorf("账号信息不匹配")
		}
		out.Newbie = 0
		out.User = user
	} else {
		out.Newbie = 1
		nowDay := gtime.Now().StartOfDay().Timestamp()
		err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, e := dao.User.Ctx(ctx).TX(tx).Data(g.Map{
				"uid": in.Uid, "platform": in.Platform, "openid": in.Openid,
			}).Insert()
			if e != nil {
				return e
			}
			_, e = dao.UserRes.Ctx(ctx).TX(tx).Data(g.Map{
				"uid": in.Uid, "gold": 200, "diamond": 100, "star": 0,
				"tili": 100, "tili_time": 0, "exp": 0, "level": 1, "day_time": nowDay,
			}).Insert()
			return e
		})
		if err != nil {
			return nil, err
		}
		out.User = &entity.User{
			Uid: uint(in.Uid), Platform: in.Platform, Openid: in.Openid,
		}
	}

	// Log login (fire-and-forget with recover)
	bgCtx := gctx.NeverDone(ctx)
	go func() {
		defer func() { recover() }()
		_, _ = dao.LogLogin.Ctx(bgCtx).Data(g.Map{"uid": in.Uid, "platform": in.Platform}).Insert()
	}()

	// Upsert login key
	_, err = dao.UserLoginkey.Ctx(ctx).Data(g.Map{
		"uid": in.Uid, "key": in.LoginKey, "ver": in.Version, "time": gtime.Timestamp(),
	}).Save()
	if err != nil {
		return nil, err
	}

	out.Datas, err = dao.UserData.Ctx(ctx).Where("uid", in.Uid).All()
	if err != nil {
		return nil, err
	}

	gmVal, err := dao.SysGm.Ctx(ctx).Where("uid", in.Uid).Value("uid")
	if err != nil {
		return nil, err
	}
	if gmVal.IsEmpty() {
		out.Gm = 0
	} else {
		out.Gm = 1
	}

	out.Items, err = dao.UserItem.Ctx(ctx).Where("uid", in.Uid).All()
	if err != nil {
		return nil, err
	}

	out.Res, err = s.GetUserRes(ctx, in.Uid)
	if err != nil {
		return nil, err
	}

	out.Config, err = dao.MemConfig.Ctx(ctx).All()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *sUser) UpdateDiamond(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error) {
	return updateResField(ctx, in, "diamond", "钻石")
}

func (s *sUser) UpdateGold(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error) {
	return updateResField(ctx, in, "gold", "金币")
}

func (s *sUser) UpdateTili(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error) {
	return updateResField(ctx, in, "tili", "体力")
}

func (s *sUser) UpdateExp(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error) {
	return updateResField(ctx, in, "exp", "经验")
}

func (s *sUser) UpdateStar(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error) {
	return updateResField(ctx, in, "star", "星星")
}

func (s *sUser) GetUser(ctx context.Context, uid int64) (*entity.User, error) {
	var user *entity.User
	err := dao.User.Ctx(ctx).Where("uid", uid).Scan(&user)
	return user, err
}

func (s *sUser) GetUserRes(ctx context.Context, uid int64) (*entity.UserRes, error) {
	var res *entity.UserRes
	err := dao.UserRes.Ctx(ctx).Where("uid", uid).Scan(&res)
	if err != nil || res == nil {
		return res, err
	}
	nowDay := int(gtime.Now().StartOfDay().Timestamp())
	if res.DayTime != nowDay {
		_, _ = dao.UserRes.Ctx(ctx).Where("uid", uid).Data(g.Map{"day_conf": "", "day_time": nowDay}).Update()
		res.DayConf = ""
		res.DayTime = nowDay
	}
	return res, nil
}

func updateResField(ctx context.Context, in *model.UpdateFieldInput, field, resName string) (*model.UpdateFieldOutput, error) {
	lockKey := fmt.Sprintf("update_%s:%d", field, in.Uid)
	token, err := lock.Lock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, fmt.Errorf("系统繁忙，请稍后再试")
	}
	defer func() { _ = lock.Unlock(ctx, lockKey, token) }()

	var res *entity.UserRes
	err = dao.UserRes.Ctx(ctx).Where("uid", in.Uid).Scan(&res)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("用户资源不存在")
	}

	var oldCnt int64
	switch field {
	case "diamond":
		oldCnt = int64(res.Diamond)
	case "gold":
		oldCnt = int64(res.Gold)
	case "tili":
		oldCnt = int64(res.Tili)
	case "exp":
		oldCnt = int64(res.Exp)
	case "star":
		oldCnt = int64(res.Star)
	}

	newCnt := oldCnt + in.Cnt
	if newCnt < 0 {
		newCnt = 0
	}
	if newCnt == oldCnt {
		return &model.UpdateFieldOutput{Res: res, AddValue: 0}, nil
	}

	_, err = dao.UserRes.Ctx(ctx).Where("uid", in.Uid).Data(g.Map{field: newCnt}).Update()
	if err != nil {
		gamelog.Log(ctx, in.Uid, fmt.Sprintf("更新用户资源失败 %s %d %s %v", field, in.Cnt, in.Reason, err))
		return nil, err
	}

	// Update the struct
	switch field {
	case "diamond":
		res.Diamond = int(newCnt)
	case "gold":
		res.Gold = int(newCnt)
	case "tili":
		res.Tili = int(newCnt)
	case "exp":
		res.Exp = int(newCnt)
	case "star":
		res.Star = int(newCnt)
	}

	gamelog.TraceRes(ctx, in.Uid, oldCnt, newCnt, resName, in.Reason)
	return &model.UpdateFieldOutput{Res: res, AddValue: newCnt - oldCnt}, nil
}

// --- Utility ---

func ParseRes(items interface{}) []consts.ResItem {
	switch v := items.(type) {
	case []consts.ResItem:
		return v
	case string:
		return parseResString(v)
	default:
		return nil
	}
}

func parseResString(s string) []consts.ResItem {
	nums := PickNumbers(s)
	if len(nums) == 0 || len(nums)%3 != 0 {
		return nil
	}
	result := make([]consts.ResItem, 0, len(nums)/3)
	for i := 0; i < len(nums); i += 3 {
		result = append(result, consts.ResItem{Type: nums[i], Id: nums[i+1], Cnt: nums[i+2]})
	}
	return result
}

// PickNumbers extracts all integers from a string.
func PickNumbers(s string) []int {
	var result []int
	current := ""
	for _, c := range s {
		if (c >= '0' && c <= '9') || c == '-' || c == '+' || c == '.' {
			current += string(c)
		} else {
			if current != "" {
				if n, err := strconv.Atoi(current); err == nil {
					result = append(result, n)
				} else if f, err := strconv.ParseFloat(current, 64); err == nil {
					result = append(result, int(math.Floor(f)))
				}
				current = ""
			}
		}
	}
	if current != "" {
		if n, err := strconv.Atoi(current); err == nil {
			result = append(result, n)
		} else if f, err := strconv.ParseFloat(current, 64); err == nil {
			result = append(result, int(math.Floor(f)))
		}
	}
	return result
}

func init() {
	_ = strings.Join
}
```

讲解：

- `Login` 是用户登录主流程。
- 已存在用户时校验 `platform/openid`，防止同一 uid 被不同账号复用。
- 新用户创建时使用事务，同时插入 `user` 和 `user_res`。
- 登录日志异步写入 `_log_login`，不阻塞主流程。
- `UserLoginkey.Save()` 保存当前登录 key，供 `Verify` 中间件校验。
- `GetUserRes` 会检查跨天逻辑，如果 `DayTime` 不是今天，就重置 `day_conf` 并更新 `day_time`。
- `updateResField` 是资源更新通用方法，负责加锁、读旧值、计算新值、更新数据库、写资源流水。
- `ParseRes` 和 `PickNumbers` 是资源配置解析工具。
- 末尾 `init(){ _ = strings.Join }` 只是为了保留 `strings` 导入不报未使用；如果后续不需要 `strings`，可以清理。
