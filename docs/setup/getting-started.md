# Unifocus-v1.2 启动文档

说明：当前仓库主名称为 `Unifocus-v1.2`；界面或设计文档中出现 `UniFocus` 时，默认按产品品牌文案理解。

## 环境要求

- macOS + OrbStack
- Go 1.21+
- Node.js 18+
- PostgreSQL（本地已安装，端口 5432）
- Redis（本地已安装，端口 6379）

---

## 第一次启动

### 1. 启动数据库容器

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus
docker-compose up -d postgres redis
```

等待约 5 秒，确认容器正常运行：

```bash
docker ps
```

### 2. 启动后端

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/backend
go run cmd/api/main.go
```

启动成功标志：
```
Database connection established: localhost:5432/unifocus_dev
Competitions seed completed
Server is running on port 8080
```

> 首次启动会自动将 84 条竞赛数据（含时间信息）写入数据库。

### 3. 启动前端

新开一个终端：

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
npm install
npm run dev
```

---

## 日常启动（第二次起）

```bash
# 终端 1：启动数据库
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus
docker-compose up -d postgres redis

# 终端 2：启动后端
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/backend
go run cmd/api/main.go

# 终端 3：启动前端
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
[ -d node_modules ] || npm install
npm run dev
```

---

## 访问地址

| 服务 | 地址 |
|------|------|
| 前端 | http://localhost:3000 |
| 竞赛管理页面 | http://localhost:3000/dashboard/competitions |
| 后端 API | http://localhost:8080 |
| 健康检查 | http://localhost:8080/health |
| 竞赛接口 | http://localhost:8080/api/v1/competitions |

---

## 竞赛数据

数据文件位置：`backend/internal/seed/competitions_full.json`

- 共 84 条竞赛记录
- 字段：名称、级别、官网链接、典型时间窗口、时间线提示、备注
- 后端启动时自动 upsert 进数据库（幂等，可重复执行）

如需更新竞赛数据，直接修改 `competitions_full.json`，重启后端即可生效。

---

## 常用命令

```bash
# 查看竞赛数据
curl -s http://localhost:8080/api/v1/competitions | python3 -m json.tool | head -30

# 查看数据库中竞赛数量
psql -h 127.0.0.1 -U unifocus -d unifocus_dev -c "SELECT COUNT(*) FROM competitions"

# 停止数据库容器
docker-compose down

# 端口被占用时 kill 进程
lsof -ti:8080 | xargs kill -9
lsof -ti:3000 | xargs kill -9
```

---

## 数据库重置（谨慎操作）

```bash
# 删除所有容器和数据卷，重新初始化
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus
docker-compose down -v
docker-compose up -d postgres redis
# 然后重启后端，会自动重新 seed 竞赛数据
```
