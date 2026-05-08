FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.20
WORKDIR /app
ENV TZ=Asia/Shanghai
ENV GF_GCFG_FILE=config.docker.yaml
RUN apk add --no-cache tzdata
COPY --from=build /app/main .
COPY manifest/config ./manifest/config
COPY resource ./resource
EXPOSE 7001
HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=5 \
    CMD wget -q --spider http://127.0.0.1:7001/health/ready || exit 1
ENTRYPOINT ["./main"]