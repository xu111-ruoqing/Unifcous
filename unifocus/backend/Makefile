# UniFocus Backend Makefile

# 加载 .env 文件（如果存在）
ifneq (,$(wildcard .env))
    include .env
    export
endif

# 数据库配置（从环境变量读取，或使用默认值）
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= yufan
DB_PASSWORD ?=
DB_NAME ?= unifocus_dev
DB_SSLMODE ?= disable

# 数据库连接字符串
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: help migrate-up migrate-down migrate-force migrate-version db-reset test run

help:
	@echo "UniFocus Backend - 可用命令："
	@echo "  make migrate-up       - 执行数据库迁移（升级）"
	@echo "  make migrate-down     - 回滚最后一次迁移"
	@echo "  make migrate-version  - 查看当前迁移版本"
	@echo "  make migrate-force    - 强制设置迁移版本（需指定 VERSION=N）"
	@echo "  make db-reset         - 重置数据库（删除所有数据并重新迁移）"
	@echo "  make db-psql          - 直接用 psql 运行迁移（不需要 migrate 工具）"
	@echo "  make run              - 运行后端服务"
	@echo "  make test             - 运行测试"

# 使用 golang-migrate 工具进行迁移（需要先安装）
migrate-up:
	@echo "执行数据库迁移..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DB_URL)" up; \
	else \
		echo "错误: 未找到 migrate 工具"; \
		echo "请运行: brew install golang-migrate"; \
		echo "或使用: make db-psql（直接用 psql 运行迁移）"; \
		exit 1; \
	fi

migrate-down:
	@echo "回滚数据库迁移..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DB_URL)" down 1; \
	else \
		echo "错误: 未找到 migrate 工具"; \
		echo "请运行: brew install golang-migrate"; \
		exit 1; \
	fi

migrate-version:
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DB_URL)" version; \
	else \
		echo "错误: 未找到 migrate 工具"; \
		exit 1; \
	fi

migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "错误: 请指定版本号，例如: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DB_URL)" force $(VERSION); \
	else \
		echo "错误: 未找到 migrate 工具"; \
		exit 1; \
	fi

# 直接用 psql 运行迁移（不需要 migrate 工具）
db-psql:
	@echo "使用 psql 直接运行迁移文件..."
	@echo "数据库: $(DB_NAME), 用户: $(DB_USER)"
	@echo "执行 001_init_schema.up.sql..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f migrations/001_init_schema.up.sql
	@echo "执行 002_recognition_policy.up.sql..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f migrations/002_recognition_policy.up.sql
	@echo "✅ 迁移完成！"

# 重置数据库
db-reset:
	@echo "警告: 这将删除数据库 $(DB_NAME) 中的所有数据！"
	@read -p "确认继续？(y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@echo "删除并重建数据库..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	@echo "运行迁移..."
	@$(MAKE) db-psql

# 运行后端服务
run:
	@echo "启动后端服务..."
	go run cmd/api/main.go

# 运行测试
test:
	@echo "运行测试..."
	go test ./... -v

# 构建
build:
	@echo "构建项目..."
	go build -o bin/unifocus cmd/api/main.go

# 安装依赖
deps:
	@echo "安装 Go 依赖..."
	go mod download
	go mod tidy

# 代码格式化
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "运行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "警告: 未安装 golangci-lint，跳过..."; \
		echo "安装: brew install golangci-lint"; \
	fi
