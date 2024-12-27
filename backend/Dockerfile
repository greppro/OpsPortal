# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装编译所需的工具和依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 启用 CGO (SQLite需要)
ENV CGO_ENABLED=1
ENV GOOS=linux

# 复制并下载依赖
COPY go.* ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN go build -o main

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache sqlite-libs ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 创建数据目录并设置权限
RUN mkdir -p /app/data && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

# 复制编译好的程序
COPY --from=builder /app/main .

# 如果有配置文件，也需要复制
# COPY --from=builder /app/config ./config

# 设置数据卷并确保权限
VOLUME ["/app/data"]
RUN chown -R appuser:appgroup /app/data

# 切换到非 root 用户
USER appuser

EXPOSE 8080

CMD ["./main"] 