ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "template-single"
DOCKER_NAME = "template-single"

include ./hack/hack-cli.mk
include ./hack/hack.mk

# 本地开发
.PHONY: dev
dev:
	@go run main.go

.PHONY: start-db
start-db:
	@go run cmd/deploy/main.go start-local-db

.PHONY: stop-db
stop-db:
	@go run cmd/deploy/main.go stop-local-db

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

# 部署
.PHONY: deploy.local
deploy.local:
	@go run cmd/deploy/main.go deploy local

.PHONY: deploy.test
deploy.test: push.test
	@go run cmd/deploy/main.go deploy test

.PHONY: deploy.production
deploy.production: push.production
	@go run cmd/deploy/main.go deploy production

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