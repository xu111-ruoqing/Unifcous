#!/bin/bash
set -euo pipefail

# 仅启动数据库和Redis（不启动应用服务）

echo "🚀 启动数据库和Redis..."

cd "$(dirname "$0")/.."

# 检查Docker是否运行
if ! docker info &> /dev/null; then
    echo "❌ Docker未运行，请启动Docker Desktop"
    exit 1
fi

# 尝试启动PostgreSQL和Redis
echo "📥 正在拉取镜像（可能需要一些时间）..."
if ! docker-compose up -d postgres redis 2>&1; then
    echo ""
    echo "❌ 启动失败！可能的原因："
    echo "   1. 网络连接问题（无法访问Docker Hub）"
    echo "   2. DNS解析问题"
    echo ""
    echo "💡 解决方案："
    echo "   方案A: 配置Docker镜像加速器（推荐）"
    echo "   - 打开Docker Desktop -> Settings -> Docker Engine"
    echo "   - 添加镜像加速器配置（见TROUBLESHOOTING.md）"
    echo ""
    echo "   方案B: 使用本地PostgreSQL/Redis"
    echo "   - 参考TROUBLESHOOTING.md中的'本地安装'部分"
    echo ""
    exit 1
fi

echo "⏳ 等待服务就绪..."
sleep 5

# 检查服务状态
echo ""
echo "📊 服务状态:"
docker-compose ps postgres redis

echo ""
echo "✅ 数据库和Redis已启动！"
echo "   PostgreSQL: localhost:5432"
echo "   Redis: localhost:6379"
echo ""
echo "💡 查看日志: docker-compose logs -f postgres redis"
echo "💡 停止服务: docker-compose stop postgres redis"

