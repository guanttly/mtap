# MTAP Docker镜像构建
# 核心目的：多阶段构建Go应用镜像
# 支持通过 BUILD_TARGET 参数指定编译目标cmd

# TODO: 实现多阶段构建
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
# RUN go build -o /bin/app ./cmd/${BUILD_TARGET}

FROM alpine:latest
# COPY --from=builder /bin/app /bin/app
# ENTRYPOINT ["/bin/app"]
