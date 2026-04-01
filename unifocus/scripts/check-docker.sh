#!/bin/bash

# Docker环境检查脚本

echo "🔍 检查Docker环境..."

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker未安装"
    echo "   请访问: https://www.docker.com/products/docker-desktop"
    exit 1
fi

echo "✅ Docker已安装: $(docker --version)"

# 检查Docker是否运行
if ! docker info &> /dev/null; then
    echo "❌ Docker未运行"
    echo "   请启动Docker Desktop应用"
    exit 1
fi

echo "✅ Docker正在运行"

# 检查网络连接
echo ""
echo "🌐 测试网络连接..."
if ping -c 1 registry-1.docker.io &> /dev/null; then
    echo "✅ 可以访问Docker Hub"
else
    echo "⚠️  无法访问Docker Hub，可能需要配置镜像加速器"
    echo "   参考: TROUBLESHOOTING.md"
fi

# 检查镜像是否存在
echo ""
echo "📦 检查本地镜像..."
if docker images | grep -q "postgres.*15-alpine"; then
    echo "✅ PostgreSQL镜像已存在"
else
    echo "⚠️  PostgreSQL镜像不存在，需要拉取"
fi

if docker images | grep -q "redis.*7-alpine"; then
    echo "✅ Redis镜像已存在"
else
    echo "⚠️  Redis镜像不存在，需要拉取"
fi

# 检查端口占用
echo ""
echo "🔌 检查端口占用..."
if lsof -i :5432 &> /dev/null; then
    echo "⚠️  端口5432已被占用"
    lsof -i :5432
else
    echo "✅ 端口5432可用"
fi

if lsof -i :6379 &> /dev/null; then
    echo "⚠️  端口6379已被占用"
    lsof -i :6379
else
    echo "✅ 端口6379可用"
fi

echo ""
echo "✅ 环境检查完成！"

