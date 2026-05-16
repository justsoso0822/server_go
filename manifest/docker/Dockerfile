# syntax=docker/dockerfile:1.7

# ---- 构建阶段 ----
FROM golang:1.26-alpine AS build
WORKDIR /app

ARG GOPROXY=https://proxy.golang.org,direct
ARG GOSUMDB=sum.golang.org
ARG GOPRIVATE=
ENV GOPROXY=${GOPROXY} \
    GOSUMDB=${GOSUMDB} \
    GOPRIVATE=${GOPRIVATE}

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o main .

# ---- 运行阶段 ----
# APP_PORT: 应用监听端口，构建时通过 --build-arg APP_PORT=xxx 注入
FROM alpine:3.20
ARG APP_PORT=7001
ENV APP_PORT=${APP_PORT}
WORKDIR /app
ENV TZ=Asia/Shanghai

RUN apk add --no-cache ca-certificates tzdata wget

COPY --from=build /app/main .

EXPOSE ${APP_PORT}
HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=5 \
    CMD wget -q -O /dev/null http://127.0.0.1:${APP_PORT}/health/ready || exit 1
ENTRYPOINT ["./main"]