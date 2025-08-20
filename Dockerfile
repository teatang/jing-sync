
# 第一阶段：构建Go二进制文件
FROM golang:1.24-alpine AS builder

# 设置编译环境变量
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# 安装SQLite编译依赖
RUN apk add --no-cache gcc musl-dev

# 设置工作目录
WORKDIR /app

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并构建
COPY . .
RUN go build -o /bin/app ./cmd/web

# 第二阶段：运行环境
FROM alpine:3.18

# 创建数据目录并设置权限
RUN mkdir -p /data && chown nobody:nobody /data
VOLUME /data

# 安装SQLite运行时依赖
RUN apk add --no-cache sqlite

# 从构建阶段复制二进制文件
COPY --from=builder --chown=nobody:nobody /bin/app /app
COPY --from=builder /app/migrations /migrations

# 设置环境变量
ENV DB_PATH=/data/db.sqlite

# 切换到非root用户
USER nobody

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["/app"]
