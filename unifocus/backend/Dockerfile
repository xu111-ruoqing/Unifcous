# ============================================
# 多阶段构建 Dockerfile for Go Backend
# ============================================

# -------------------- 开发环境 --------------------
FROM golang:1.25-alpine AS development

ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

# 安装必要工具
RUN apk add --no-cache git make gcc musl-dev

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 暴露端口
EXPOSE 8080

# 开发模式：使用热重载
CMD ["go", "run", "cmd/api/main.go"]

# -------------------- 构建阶段 --------------------
FROM golang:1.25-alpine AS builder

ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/api \
    ./cmd/api

# -------------------- 生产环境 --------------------
FROM alpine:latest AS production

# 安装 CA 证书和时区数据
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/api /app/api

# 复制配置文件
COPY --from=builder /app/configs /app/configs

# 创建日志目录
RUN mkdir -p /var/log/unifocus && \
    chown -R appuser:appuser /app /var/log/unifocus

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["/app/api"]
