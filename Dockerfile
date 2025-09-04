# 第一阶段：构建Node.js前端
FROM node:22-alpine as frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend .
RUN npm run build

# 第二阶段：构建Go后端（包含SQLite）
FROM golang:1.24-alpine as backend-builder
WORKDIR /go/src/app
COPY . .
# 复制前端构建产物
COPY --from=frontend-builder /app/frontend/dist ./frontend
RUN apk add --no-cache gcc musl-dev # SQLite编译依赖
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-extldflags=-static" -o /go/bin/app ./cmd/web

# 第三阶段：生成最终镜像
FROM alpine:latest
WORKDIR /app
# 安装SQLite运行时
RUN apk add --no-cache sqlite
# 复制后端可执行文件
COPY --from=backend-builder /go/bin/app .

EXPOSE 8888
VOLUME /app/data # SQLite数据存储卷
CMD ["./app"]