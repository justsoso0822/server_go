# ---- 构建阶段 ----
FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w" -o main .

# ---- 运行阶段 ----
# APP_PORT: 应用监听端口，构建时通过 --build-arg APP_PORT=xxx 注入
FROM alpine:3.20
ARG APP_PORT=7001
ENV APP_PORT=${APP_PORT}
WORKDIR /app
ENV TZ=Asia/Shanghai
COPY --from=build /app/main .
EXPOSE ${APP_PORT}
HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=5 \
    CMD wget -q --spider http://127.0.0.1:${APP_PORT}/health/ready || exit 1
ENTRYPOINT ["./main"]