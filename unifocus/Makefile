.PHONY: help up down logs ps restart clean db-migrate db-reset test-api test-nlp build

# 默认目标
help: ## 显示帮助信息
	@echo "UniFocus 开发命令"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ============================================
# Docker 操作
# ============================================

up: ## 启动所有服务
	docker-compose up -d
	@echo "✅ 所有服务已启动"
	@echo "API服务: http://localhost:8080"
	@echo "NLP服务: http://localhost:8000"
	@echo "pgAdmin: http://localhost:5050"

down: ## 停止所有服务
	docker-compose down
	@echo "✅ 所有服务已停止"

logs: ## 查看所有服务日志
	docker-compose logs -f

logs-api: ## 查看API服务日志
	docker-compose logs -f api

logs-nlp: ## 查看NLP服务日志
	docker-compose logs -f nlp-service

ps: ## 查看服务状态
	docker-compose ps

restart: ## 重启所有服务
	docker-compose restart
	@echo "✅ 所有服务已重启"

restart-api: ## 重启API服务
	docker-compose restart api
	@echo "✅ API服务已重启"

restart-nlp: ## 重启NLP服务
	docker-compose restart nlp-service
	@echo "✅ NLP服务已重启"

clean: ## 清理所有容器和卷
	docker-compose down -v
	@echo "✅ 已清理所有容器和数据卷"

# ============================================
# 数据库操作
# ============================================

db-migrate: ## 执行数据库迁移
	docker exec -i unifocus_postgres psql -U unifocus -d unifocus_dev < backend/migrations/001_init_schema.up.sql
	docker exec -i unifocus_postgres psql -U unifocus -d unifocus_dev < backend/migrations/002_recognition_policy.up.sql
	@echo "✅ 数据库迁移完成"

db-reset: ## 重置数据库
	docker exec -i unifocus_postgres psql -U unifocus -d unifocus_dev < backend/migrations/002_recognition_policy.down.sql
	docker exec -i unifocus_postgres psql -U unifocus -d unifocus_dev < backend/migrations/001_init_schema.down.sql
	docker exec -i unifocus_postgres psql -U unifocus -d unifocus_dev < backend/migrations/001_init_schema.up.sql
	docker exec -i unifocus_postgres psql -U unifocus -d unifocus_dev < backend/migrations/002_recognition_policy.up.sql
	@echo "✅ 数据库已重置"

db-shell: ## 进入数据库Shell
	docker exec -it unifocus_postgres psql -U unifocus -d unifocus_dev

# ============================================
# 测试
# ============================================

test-api: ## 测试API健康检查
	@echo "测试 API 服务..."
	@curl -s http://localhost:8080/health | python3 -m json.tool || echo "❌ API服务未响应"

test-nlp: ## 测试NLP健康检查
	@echo "测试 NLP 服务..."
	@curl -s http://localhost:8000/health | python3 -m json.tool || echo "❌ NLP服务未响应"

test-all: test-api test-nlp ## 测试所有服务

# ============================================
# 开发
# ============================================

dev-backend: ## 启动后端开发环境
	cd backend && go run cmd/api/main.go

dev-nlp: ## 启动NLP开发环境
	cd nlp-service && uvicorn app.main:app --reload --port 8000

install-backend: ## 安装Go依赖
	cd backend && go mod download
	@echo "✅ Go依赖已安装"

install-nlp: ## 安装Python依赖
	cd nlp-service && pip install -r requirements.txt
	@echo "✅ Python依赖已安装"

# ============================================
# 构建
# ============================================

build: ## 构建所有服务
	docker-compose build
	@echo "✅ 所有服务已构建"

build-api: ## 构建API服务
	docker-compose build api
	@echo "✅ API服务已构建"

build-nlp: ## 构建NLP服务
	docker-compose build nlp-service
	@echo "✅ NLP服务已构建"

# ============================================
# 工具
# ============================================

fmt-go: ## 格式化Go代码
	cd backend && go fmt ./...
	@echo "✅ Go代码已格式化"

lint-go: ## Go代码静态检查
	cd backend && go vet ./...
	@echo "✅ Go代码检查完成"

fmt-python: ## 格式化Python代码
	cd nlp-service && black .
	@echo "✅ Python代码已格式化"

# ============================================
# 生产环境
# ============================================

prod-up: ## 启动生产环境
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
	@echo "✅ 生产环境已启动"

prod-down: ## 停止生产环境
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml down
	@echo "✅ 生产环境已停止"
