# MTAP Docker镜像构建
# 多阶段构建：builder 编译 + 最终镜像只含二进制
# 构建示例：
#   docker build --build-arg BUILD_TARGET=api-server -t mtap-api .
#   docker build --build-arg BUILD_TARGET=scheduler   -t mtap-scheduler .

ARG BUILD_TARGET=api-server

# ───────────────────────── Stage 1: builder ─────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build

# 先复制 go.mod/go.sum 利用缓层
COPY src/go.mod src/go.sum ./
RUN go mod download

# 复制源码
COPY src/ .

# 编译指定目标
ARG BUILD_TARGET
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /bin/app \
    ./cmd/${BUILD_TARGET}

# ──────────────────────── Stage 2: runtime ────────────────────────
FROM alpine:3.19

RUN apk add --no-cache tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /bin/app /app/app
COPY src/configs/ /app/configs/

# 日志目录
RUN mkdir -p /var/log/mtap

ENV MTAP_CONFIG=/app/configs/config.prod.yaml

EXPOSE 8080

ENTRYPOINT ["/app/app"]
