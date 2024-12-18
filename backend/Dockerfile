# 构建阶段
FROM golang:1.21-alpine AS builder

# 安装编译所需的工具
RUN apk add --no-cache gcc musl-dev

# 启用 CGO
ENV CGO_ENABLED=1

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o main

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates gcc musl-dev

WORKDIR /app
COPY --from=builder /app/main .

# 如果有配置文件或其他资源文件，也需要复制
# COPY --from=builder /app/config ./config

EXPOSE 8080
CMD ["./main"] 