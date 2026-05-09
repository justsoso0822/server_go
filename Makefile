ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "template-single"
DOCKER_NAME = "template-single"

include ./hack/hack-cli.mk
include ./hack/hack.mk

# 本地开发环境
.PHONY: local.start
local.start:
	@go run cmd/deploy/main.go start-local

.PHONY: local.stop
local.stop:
	@go run cmd/deploy/main.go stop-local

.PHONY: local.run
local.run:
	@GF_GCFG_FILE=config.local.yaml go run main.go

# 构建和推送镜像
.PHONY: build.local
build.local:
	@go run cmd/deploy/main.go build local

.PHONY: build.test
build.test:
	@go run cmd/deploy/main.go build test

.PHONY: build.production
build.production:
	@go run cmd/deploy/main.go build production

.PHONY: push.local
push.local: build.local
	@go run cmd/deploy/main.go push local

.PHONY: push.test
push.test: build.test
	@go run cmd/deploy/main.go push test

.PHONY: push.production
push.production: build.production
	@go run cmd/deploy/main.go push production

# 蓝绿部署
.PHONY: deploy.local.blue
deploy.local.blue:
	@go run cmd/deploy/main.go deploy local blue

.PHONY: deploy.local.green
deploy.local.green:
	@go run cmd/deploy/main.go deploy local green

.PHONY: deploy.test.blue
deploy.test.blue: push.test
	@go run cmd/deploy/main.go deploy test blue

.PHONY: deploy.test.green
deploy.test.green: push.test
	@go run cmd/deploy/main.go deploy test green

.PHONY: deploy.production.blue
deploy.production.blue: push.production
	@go run cmd/deploy/main.go deploy production blue

.PHONY: deploy.production.green
deploy.production.green: push.production
	@go run cmd/deploy/main.go deploy production green

# 回滚
.PHONY: rollback.local
rollback.local:
	@go run cmd/deploy/main.go rollback local

.PHONY: rollback.test
rollback.test:
	@go run cmd/deploy/main.go rollback test

.PHONY: rollback.production
rollback.production:
	@go run cmd/deploy/main.go rollback production

# 状态查看
.PHONY: status.local
status.local:
	@go run cmd/deploy/main.go status local

.PHONY: status.test
status.test:
	@go run cmd/deploy/main.go status test

.PHONY: status.production
status.production:
	@go run cmd/deploy/main.go status production

# 清理
.PHONY: cleanup.local
cleanup.local:
	@go run cmd/deploy/main.go cleanup local

.PHONY: cleanup.test
cleanup.test:
	@go run cmd/deploy/main.go cleanup test

.PHONY: cleanup.production
cleanup.production:
	@go run cmd/deploy/main.go cleanup production