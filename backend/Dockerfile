# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装编译所需的工具和依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 安装 swag 工具（版本与 go.mod 中的 github.com/swaggo/swag 保持一致）
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

# 启用 CGO (SQLite需要)
ENV CGO_ENABLED=1
ENV GOOS=linux

# 复制并下载依赖
COPY go.* ./
RUN go mod download

# 复制源代码
COPY . .

# 生成 Swagger 文档
RUN swag init

# 编译
RUN go build -o main

# 运行阶段
FROM alpine:latest

# 安装运行时依赖（su-exec 用于入口脚本降权）
RUN apk add --no-cache sqlite-libs ca-certificates tzdata curl su-exec

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 创建数据目录、上传目录和运行用户
RUN mkdir -p /app/data /app/uploads/logos && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

# 复制编译好的程序和入口脚本
COPY --from=builder --chown=appuser:appgroup /app/main .
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod 755 /usr/local/bin/docker-entrypoint.sh

# 入口脚本以 root 修正挂载目录属主后，用 su-exec 切换到 appuser 运行程序
EXPOSE 8080

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["./main"]
