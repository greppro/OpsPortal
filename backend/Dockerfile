# 构建阶段
FROM golang:1.19-alpine as builder

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 生产阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/main .
RUN mkdir -p data

EXPOSE 8080
CMD ["./main"] 