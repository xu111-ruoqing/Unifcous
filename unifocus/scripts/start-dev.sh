#!/bin/bash

# UniFocus 开发环境启动脚本

echo "🚀 启动 UniFocus 开发环境..."

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker Desktop"
    echo "   下载地址: https://www.docker.com/products/docker-desktop"
    exit 1
fi

# 检查Docker是否运行
if ! docker info &> /dev/null; then
    echo "❌ Docker 未运行，请启动 Docker Desktop"
    exit 1
fi

# 进入项目目录
cd "$(dirname "$0")/.."

echo "📦 启动数据库和Redis服务..."
docker-compose up -d postgres redis

echo "⏳ 等待数据库就绪..."
sleep 5

# 检查数据库是否就绪
until docker exec unifocus_postgres pg_isready -U unifocus &> /dev/null; do
    echo "等待PostgreSQL启动..."
    sleep 2
done

echo "✅ 数据库已就绪！"
echo ""
echo "📊 服务状态:"
docker-compose ps postgres redis

echo ""
echo "🎯 下一步:"
echo "   1. 启动后端: cd backend && go run cmd/api/main.go"
echo "   2. 启动NLP服务: cd nlp-service && python -m app.main"
echo "   3. 启动前端: cd web && npm run dev"
echo ""
echo "💡 提示: 使用 'docker-compose logs -f' 查看日志"
echo "💡 提示: 使用 'docker-compose down' 停止服务"


