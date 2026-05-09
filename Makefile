ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "server-go"
DOCKER_NAME = "server-go"
VERSION     ?=

include ./hack/hack-cli.mk
include ./hack/hack.mk

# 本地开发
.PHONY: dev
dev:
	@go run main.go

.PHONY: start-local-db
start-local-db:
	@go run cmd/deploy/main.go start-local-db

.PHONY: stop-local-db
stop-local-db:
	@go run cmd/deploy/main.go stop-local-db

# 构建和推送镜像
.PHONY: build.local
build.local:
	@go run cmd/deploy/main.go build local $(if $(VERSION),version=$(VERSION),)

.PHONY: build.test
build.test:
	@[ -n "$(VERSION)" ] || (echo "Error: VERSION is required for build.test" && exit 1)
	@go run cmd/deploy/main.go build test version=$(VERSION)

.PHONY: build.production
build.production:
	@[ -n "$(VERSION)" ] || (echo "Error: VERSION is required for build.production" && exit 1)
	@go run cmd/deploy/main.go build production version=$(VERSION)

.PHONY: push.local
push.local:
	@go run cmd/deploy/main.go push local $(if $(VERSION),version=$(VERSION),)

.PHONY: push.test
push.test:
	@[ -n "$(VERSION)" ] || (echo "Error: VERSION is required for push.test" && exit 1)
	@go run cmd/deploy/main.go push test version=$(VERSION)

.PHONY: push.production
push.production:
	@[ -n "$(VERSION)" ] || (echo "Error: VERSION is required for push.production" && exit 1)
	@go run cmd/deploy/main.go push production version=$(VERSION)

# 部署
.PHONY: deploy.local
deploy.local:
	@go run cmd/deploy/main.go deploy local $(if $(VERSION),version=$(VERSION),)

.PHONY: deploy.test
deploy.test:
	@go run cmd/deploy/main.go deploy test $(if $(VERSION),version=$(VERSION),)

.PHONY: deploy.production
deploy.production:
	@go run cmd/deploy/main.go deploy production $(if $(VERSION),version=$(VERSION),)

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