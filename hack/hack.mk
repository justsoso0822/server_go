.DEFAULT_GOAL := build

# 更新 GoFrame 及其 CLI 到最新稳定版本。
.PHONY: up
up: cli.install
	@gf up -a

# 使用 hack/config.yaml 配置构建二进制文件。
.PHONY: build
build: cli.install
	@gf build -ew

# 解析 api 并生成 controller/sdk。
.PHONY: ctrl
ctrl: cli.install
	@gf gen ctrl

# 为 DAO/DO/Entity 生成 Go 文件。
.PHONY: dao
dao: cli.install
	@gf gen dao

# 解析当前项目 Go 文件并生成枚举 Go 文件。
.PHONY: enums
enums: cli.install
	@gf gen enums

# 生成 Service 的 Go 文件。
.PHONY: service
service: cli.install
	@gf gen service


# 构建 docker 镜像。
.PHONY: image
image: cli.install
	$(eval _TAG  = $(shell git rev-parse --short HEAD))
ifneq (, $(shell git status --porcelain 2>/dev/null))
	$(eval _TAG  = $(_TAG).dirty)
endif
	$(eval _TAG  = $(if ${TAG},  ${TAG}, $(_TAG)))
	$(eval _PUSH = $(if ${PUSH}, ${PUSH}, ))
	@gf docker ${_PUSH} -tn $(DOCKER_NAME):${_TAG};


# 构建 docker 镜像并自动推送到 docker 仓库。
.PHONY: image.push
image.push: cli.install
	@make image PUSH=-p;


# 将镜像和 yaml 部署到当前 kubectl 环境。
.PHONY: deploy
deploy: cli.install
	$(eval _TAG = $(if ${TAG},  ${TAG}, develop))

	@set -e; \
	mkdir -p $(ROOT_DIR)/temp/kustomize;\
	cd $(ROOT_DIR)/manifest/deploy/kustomize/overlays/${_ENV};\
	kustomize build > $(ROOT_DIR)/temp/kustomize.yaml;\
	kubectl   apply -f $(ROOT_DIR)/temp/kustomize.yaml; \
	if [ $(DEPLOY_NAME) != "" ]; then \
		kubectl patch -n $(NAMESPACE) deployment/$(DEPLOY_NAME) -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"$(shell date +%s)\"}}}}}"; \
	fi;


# 解析 protobuf 文件并生成 Go 文件。
.PHONY: pb
pb: cli.install
	@gf gen pb

# 为数据库表生成 protobuf 文件。
.PHONY: pbentity
pbentity: cli.install
	@gf gen pbentity