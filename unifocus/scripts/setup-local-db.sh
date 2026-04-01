#!/bin/bash

# 本地PostgreSQL/Redis安装指南脚本

echo "📚 本地数据库安装指南"
echo ""
echo "如果Docker镜像拉取失败，可以使用本地安装的PostgreSQL和Redis"
echo ""

# 检查操作系统
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "检测到 macOS 系统"
    echo ""
    echo "安装PostgreSQL:"
    echo "  brew install postgresql@15"
    echo "  brew services start postgresql@15"
    echo ""
    echo "安装Redis:"
    echo "  brew install redis"
    echo "  brew services start redis"
    echo ""
    echo "创建数据库:"
    echo "  createdb unifocus_dev"
    echo "  psql unifocus_dev"
    echo "  CREATE USER unifocus WITH PASSWORD 'unifocus_dev_password';"
    echo "  GRANT ALL PRIVILEGES ON DATABASE unifocus_dev TO unifocus;"
    echo "  \\q"
    echo ""
    echo "执行迁移:"
    echo "  psql -U unifocus -d unifocus_dev < backend/migrations/001_init_schema.up.sql"
    echo ""
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "检测到 Linux 系统"
    echo ""
    echo "安装PostgreSQL (Ubuntu/Debian):"
    echo "  sudo apt update"
    echo "  sudo apt install postgresql-15 postgresql-contrib"
    echo "  sudo systemctl start postgresql"
    echo ""
    echo "安装Redis:"
    echo "  sudo apt install redis-server"
    echo "  sudo systemctl start redis"
    echo ""
    echo "创建数据库:"
    echo "  sudo -u postgres psql"
    echo "  CREATE USER unifocus WITH PASSWORD 'unifocus_dev_password';"
    echo "  CREATE DATABASE unifocus_dev OWNER unifocus;"
    echo "  \\q"
    echo ""
    echo "执行迁移:"
    echo "  psql -U unifocus -d unifocus_dev < backend/migrations/001_init_schema.up.sql"
    echo ""
else
    echo "未识别的操作系统，请参考TROUBLESHOOTING.md"
fi

echo ""
echo "💡 配置完成后，确保backend/configs/config.dev.yaml中的配置正确"
echo "💡 然后运行: cd backend && go run cmd/api/main.go"


