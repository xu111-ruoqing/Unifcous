# 项目基线文档

最后更新：2026-04-01
维护说明：后续每次执行新的更改计划前，先阅读根目录 `agent.md`，再阅读 `docs/project/agent-execution-rules.md`、`docs/project/current-status-and-direction.md`、`docs/project/project-memory.md`、`docs/project/baseline.md`、`docs/project/structure.md`，再判断实现入口、影响范围和是否偏离当前真实结构。

## 0. 当前续接入口

后续继续本项目时，默认按这个顺序回读文档：

1. `agent.md`
   - 看完整主文档、当前协作规则、前端修复经验和文档分层。
2. `docs/project/agent-execution-rules.md`
   - 看技能调用与执行细则。
3. `docs/project/current-status-and-direction.md`
   - 看最近一次对话结论、当前启动状态、下一阶段方向。
4. `docs/project/project-memory.md`
   - 看长期决策、环境特例、边界规则和遗留项。
5. `docs/project/baseline.md`
   - 看项目真实现状、运行链路、关键偏差。
6. `docs/project/structure.md`
   - 看目录落点、入口文件和空目录边界。
7. `docs/frontend/issues.md`
   - 如果当前任务是前端，优先继续看这个问题清单和下一批任务。

## 1. 已阅读的现有文档

- `docs/setup/getting-started.md`
  - 说明本地启动顺序：先数据库，再 Go 后端，再 Next 前端。
  - 明确当前用户可见页面是 `http://localhost:3000/dashboard/competitions`。
  - 明确竞赛数据由 `unifocus/backend/internal/seed/competitions_full.json` 在后端启动时自动写入数据库。
- `docs/templates/cross-project-collaboration-template.md`
  - 这是协作规则模板，不描述本项目业务实现。
  - 可作为工作方法参考，但不能当作项目事实来源。

## 2. 项目当前真实内容

当前仓库核心是 `unifocus/`，它正在尝试组成一个“机会/竞赛信息平台”：

- `backend/`：Go + Gin API，负责认证、机会、竞赛、画像、认定规则、爬虫任务、数据库连接。
- `web/`：Next.js 14 前端，目前真正落地的页面只有竞赛管理页。
- `web-vite/`：新增并行前端实验入口，使用 `React + Vite + JavaScript + three.js` 实现“机会星球”双风格界面。
- `nlp-service/`：FastAPI 微服务，当前真正实现的是文本提取和文本向量化。
- `docker-compose.yml`：当前只编排 `postgres`、`redis`、`api`，没有把 `web` 和 `nlp-service` 放进统一编排。
- `scripts/`：本地开发辅助脚本，主要负责数据库容器启动检查。

当前真实状态不是“完整产品已经打通”，而是“后端主体已存在、旧前端只接了一个页面、新前端实验入口已补出主界面、NLP 半接入、部分目录预留未使用”。

## 3. 当前主链路

### 3.1 已打通的链路

- 浏览器访问 `web/app/dashboard/competitions/page.tsx`
- 页面通过 `fetch('/api/v1/competitions')` 请求数据
- `web/next.config.js` 将 `/api/:path*` 重写到 Go 后端 `http://localhost:8080/api/:path*`
- `backend/cmd/api/main.go` 注册 `competition` 路由
- `backend/internal/api/handlers/competition_handler.go` 处理 CRUD
- `backend/internal/repository/postgres/competition_repository.go` 直接操作 `competitions` 表

### 3.2 部分接入链路

- `backend/internal/service/profile_service.go` 会调用 NLP 客户端进行简历向量化
- 当前真实匹配的 NLP 客户端是 `backend/internal/infrastructure/client/nlp_client.go`
- 它调用 Python 服务的 `/api/v1/vectorize/text`
- Python 服务入口是 `nlp-service/app/main.py`

### 3.3 预留但未闭环链路

- 前端 `lib/api/auth.ts`、`profile.ts`、`opportunities.ts`、`metrics.ts` 已有封装
- 但前端没有对应页面真正使用这些封装
- `client.ts` 401 后跳转 `/login`，而当前前端没有 `/login` 页面
- 新增的 `web-vite/` 当前也还没有登录闭环，仍属于界面层改造阶段

## 4. 关键偏差与风险

- `unifocus/Makefile` 假定存在 `nlp-service`、`pgAdmin` 等容器，但 `unifocus/docker-compose.yml` 当前没有定义它们。
- `unifocus/scripts/check-docker.sh` 检查的镜像版本与 `docker-compose.yml` 不一致。
- `backend/internal/nlp/client/client.go` 使用的接口路径和当前 Python 服务不匹配，更像旧版本残留。
- `nlp-service/requirements.txt` 声明了很多 OCR/Spacy 依赖，但当前代码只真正实现文本提取和向量化。
- `unifocus/extension/`、`unifocus/infrastructure/`、`unifocus/tests/` 当前仍以预留结构为主，不应误判为已接入模块。
- 仓库当前是脏工作区，已有许多用户修改，后续变更必须避免覆盖现有未提交内容。
- 本机存在旧 PostgreSQL 进程监听 `127.0.0.1:5432`，会与 Docker 数据库端口发生语义冲突。
- `yufan` 在当前项目里不是单一语义：
  - `/Users/yufan/...` 在 `nlp-service/.venv` 中表示历史机器路径。
  - `DB_USER=yufan` 表示数据库连接账号。
  - 网站业务登录仍然是 `username/email/password` 这套用户体系，不等于数据库账号。

## 5. 分文件理解

以下说明以“核心源码与关键配置”为主；生成目录如 `.gocache`、`.next`、`node_modules`、`.venv` 不作为业务基线的一部分。

### 5.1 根目录与通用文档

- `竞赛final.json`
  - 顶层单独数据文件，目前未发现与 `unifocus/` 源码直接连接。
  - 后续若使用，需先确认它与 `backend/data/` 或 `internal/seed/` 的关系。
- `docs/project/current-status-and-direction.md`
  - 当前对话阶段总结文档。
  - 记录最近一次已完成事项、已知运行状态、下一阶段方向。
- `docs/frontend/issues.md`
  - 当前前端问题与已完成修复清单。
  - 当前前端工作前应优先阅读。
- `docs/setup/getting-started.md`
  - 当前唯一明确描述运行方式的项目文档。
  - 关联 `unifocus/docker-compose.yml`、`unifocus/backend/cmd/api/main.go`、`unifocus/web/package.json`。
- `docs/templates/cross-project-collaboration-template.md`
  - 协作模板，不参与运行链路。

### 5.2 `unifocus/` 根层文件

- `unifocus/Makefile`
  - 提供容器、开发、测试、格式化命令。
  - 关联 `docker-compose.yml`、`backend/`、`nlp-service/`。
  - 现状：命令集合比真实编排更大，存在与 `docker-compose.yml` 不一致的地方。
- `unifocus/docker-compose.yml`
  - 当前真实容器编排入口。
  - 只定义 `postgres`、`redis`、`api`。
  - 关联 `backend/Dockerfile` 与 `backend/migrations/init_all.sql`。
  - 当前已验证：`postgres`、`redis`、`api` 可通过 Docker 启动。
- `unifocus/LICENSE`
  - MIT 许可证，不影响代码执行。
- `unifocus/.env.example`
  - 环境变量模板。
  - 关联 `backend/internal/config/config.go` 与 `web/next.config.js`。
- `unifocus/.gitignore`
  - 忽略规则，影响仓库整洁度，不参与业务逻辑。

### 5.3 `unifocus/scripts/`

- `unifocus/scripts/start-dev.sh`
  - 启动本地开发前置依赖，主要是 `postgres` 和 `redis`。
  - 关联 `docker-compose.yml`。
  - 会提示用户手动启动 backend、nlp-service、web。
- `unifocus/scripts/start-db-only.sh`
  - 只启动数据库相关容器。
  - 适合本地调试后端与前端时使用。
- `unifocus/scripts/check-docker.sh`
  - 检查 Docker、镜像和端口。
  - 现状：检查镜像版本与真实 compose 不一致。
- `unifocus/scripts/setup-local-db.sh`
  - 描述脱离 Docker 的本地数据库安装方式。
  - 只解决 PostgreSQL/Redis，不处理 NLP 与前端依赖。

### 5.4 `unifocus/backend/` 配置与入口

- `unifocus/backend/go.mod`
  - 后端依赖声明。
  - 明确核心栈为 Gin、PostgreSQL、Redis、JWT、Zap、GoQuery。
- `unifocus/backend/go.sum`
  - Go 依赖锁定文件，不作为业务入口理解对象。
- `unifocus/backend/Dockerfile`
  - 后端容器构建方式。
  - 关联 `docker-compose.yml`。
  - 当前已修正 Go 版本到 `1.25`，以匹配 `go.mod` 的 `go >= 1.24.0` 要求。
- `unifocus/backend/Makefile`
  - 后端局部命令集合。
  - 主要辅助 Go 开发。
- `unifocus/backend/.env`
  - 后端本地环境文件，供 `config.Load()` 尝试读取。
  - 当前本地开发配置使用 `DB_USER=yufan`。
- `unifocus/backend/cmd/api/main.go`
  - Go 后端主入口。
  - 负责加载配置、初始化日志、连接 DB/Redis、执行竞赛 seed、创建 service、注册路由、启动 HTTP 服务。
  - 关联 `internal/config`、`internal/api/handlers`、`internal/service`、`internal/repository`、`internal/seed`、`pkg/jwt`、`pkg/logger`。
  - 当前已验证：容器内运行时可完成竞赛 seed，并成功监听 `8080`。
- `unifocus/backend/cmd/seed-competitions/main.go`
  - 手动触发竞赛数据 seed 的入口。
  - 关联 `internal/seed/competitions_seed.go`。
- `unifocus/backend/cmd/demo/schema_demo.go`
  - 演示/实验入口，不属于主运行链路。
- `unifocus/backend/api`
  - 当前是独立文件而非目录，未见明确业务用途，后续修改前应先核对其角色。
- `unifocus/backend/test-crawler`
  - 独立测试或脚本文件，未接入统一运行流程。

### 5.5 `unifocus/backend/configs/`

- `config.dev.yaml`
  - 本地开发默认配置。
  - `APP_ENV` 默认指向 `dev` 时使用。
  - 当前本地开发账号已同步为 `yufan`。
- `config.prod.yaml`
  - 生产环境配置模板。
- `config.docker.yaml`
  - Docker 环境配置。
  - 与 `docker-compose.yml` 中的环境变量覆盖共同作用。
- `internal/config/config.go`
  - 后端配置总入口。
  - 负责读取 `.env`、选择 YAML、环境变量覆盖、默认值和校验。
  - 关联全部运行配置来源。
- `internal/config/config_test.go`
  - 配置加载测试文件。

### 5.6 `unifocus/backend/pkg/`

- `pkg/logger/logger.go`
  - 后端日志初始化与封装。
  - 在 `cmd/api/main.go` 启动时被调用。
- `pkg/jwt/jwt.go`
  - JWT 生成、验证与刷新能力。
  - 关联 `internal/service/auth_service.go`。

### 5.7 `unifocus/backend/internal/domain/`

- `competition.go`
  - 竞赛领域对象。
  - 关联 `competition_handler.go`、`competition_repository.go`、`competitions_seed.go`。
- `opportunity.go`
  - 机会对象、过滤条件和相关请求结构。
  - 关联机会 handler/service/repository，以及 recognition 逻辑。
- `recognition.go`
  - 认定规则、认定结果、缓存结构。
  - 关联 `recognition_service.go` 与 `recognition_repository.go`。
- `crawl_task.go`
  - 爬虫任务实体与调度字段。
  - 关联 `crawler_service.go`、`crawl_task_handler.go`、`crawl_task_repository.go`。
- `user.go`
  - 用户、登录注册、用户画像等核心用户结构。
  - 关联 `auth_service.go`、`profile_service.go`、用户相关 repository。

### 5.8 `unifocus/backend/internal/api/handlers/`

- `auth_handler.go`
  - 认证接口层。
  - 关联 `auth_service.go`。
- `opportunity_handler.go`
  - 机会列表、详情、增删改接口层。
  - 关联 `opportunity_service.go`。
- `competition_handler.go`
  - 当前已被前端实际调用的竞赛 CRUD 接口层。
  - 关联 `competition_repository.go`。
- `profile_handler.go`
  - 当前用户画像与简历上传接口层。
  - 关联 `profile_service.go`。
- `metrics_handler.go`
  - 指标输出接口层。
  - 关联 `metrics.ts` 前端封装，但当前无前端页面消费。
- `crawl_task_handler.go`
  - 爬虫任务接口层。
  - 关联 `crawler_service.go`。

### 5.9 `unifocus/backend/internal/api/middleware/`

- `auth_middleware.go`
  - 解析并校验 JWT，将鉴权能力挂给受保护路由。
  - 关联 `AuthService`。
- `metrics_middleware.go`
  - 请求指标中间件。
  - 关联 `metrics_handler.go`。

### 5.10 `unifocus/backend/internal/service/`

- `auth_service.go`
  - 用户注册、登录、令牌验证和刷新。
  - 关联 `user_repository.go` 与 `pkg/jwt/jwt.go`。
- `opportunity_service.go`
  - 机会 CRUD 与认定信息补充。
  - 关联 `opportunity_repository.go` 与 `recognition_service.go`。
- `recognition_service.go`
  - 根据规则、缓存和画像计算认定结果。
  - 关联 `recognition_repository.go` 与 `domain/recognition.go`。
- `profile_service.go`
  - 用户画像更新、简历上传、向量化处理。
  - 关联 `profile_repository.go` 与 NLP 客户端。
- `crawler_service.go`
  - 爬虫任务管理与调度开关。
  - 关联 `crawler/scheduler.go` 与 `crawl_task_repository.go`。

### 5.11 `unifocus/backend/internal/repository/postgres/`

- `connection.go`
  - PostgreSQL 连接、健康检查、事务封装。
  - 是所有 Postgres repository 的基础依赖。
- `user_repository.go`
  - 用户表相关读写。
  - 关联 `auth_service.go`。
- `profile_repository.go`
  - 用户画像读写。
  - 关联 `profile_service.go`。
- `opportunity_repository.go`
  - 机会数据读写与筛选。
  - 关联 `opportunity_service.go`。
- `competition_repository.go`
  - 竞赛列表、详情、CRUD、`name_key` 归一化。
  - 关联当前唯一已接入前端业务页。
- `recognition_repository.go`
  - 认定规则与缓存读写。
  - 关联 `recognition_service.go`。
- `crawl_task_repository.go`
  - 爬虫任务存储层。
  - 关联 `crawler_service.go`。

### 5.12 `unifocus/backend/internal/repository/redis/`

- `client.go`
  - Redis 连接与健康检查。
  - 在 `cmd/api/main.go` 启动时初始化。

### 5.13 `unifocus/backend/internal/crawler/`

- `scheduler.go`
  - 爬虫调度器。
  - 当前 `cmd/api/main.go` 中以 `nil pipeline` 初始化，说明真正抓取流程仍暂停。
- `scrapers/base_scraper.go`
  - 抓取器公共抽象。
- `scrapers/static_scraper.go`
  - 静态抓取实现。
- `scrapers/cy_ncss_scraper.go`
  - 具体站点抓取实现。

### 5.14 `unifocus/backend/internal/infrastructure/` 与 `internal/nlp/`

- `internal/infrastructure/client/nlp_client.go`
  - 当前更可信的 NLP HTTP 客户端。
  - 真实匹配 Python 服务的 `/api/v1/vectorize/text`。
  - `ExtractTextFromPDF` 与 `ExtractSkills` 仍未真正接完。
- `internal/nlp/client/client.go`
  - 旧式或实验性 NLP 客户端。
  - 路径 `/vectorize`、`/extract` 与当前 Python 服务不匹配。
  - 后续修改前要先判断是否保留或废弃。

### 5.15 `unifocus/backend/internal/seed/` 与 `backend/data/`

- `internal/seed/competitions_seed.go`
  - 启动时竞赛数据 seed 逻辑。
  - 通过 `//go:embed competitions_full.json` 内嵌数据并向 `competitions` 表执行 upsert。
- `internal/seed/competitions_full.json`
  - 当前主 seed 数据源。
  - 包含竞赛名称、级别、官网、时间窗口、时间线提示、备注。
- `internal/seed/competitions_master.json`
  - 同类数据文件，字段更偏原始主表。
- `backend/data/competitions_full.json`
  - 额外保留的数据副本。
  - 当前启动链路并不直接读取它。
- `backend/data/competitions_master.json`
  - 额外保留的数据副本。
  - 当前启动链路并不直接读取它。

### 5.16 `unifocus/backend/migrations/`

- `001_init_schema.up.sql` / `.down.sql`
  - 基础表结构初始化与回滚。
- `002_recognition_policy.up.sql` / `.down.sql`
  - 认定规则相关结构。
- `003_opportunities_minimal.up.sql` / `.down.sql`
  - 机会表最小结构补充。
- `003_add_opportunity_schedule_fields.up.sql` / `.down.sql`
  - 机会时间字段扩展。
- `004_competitions_master.up.sql` / `.down.sql`
  - 竞赛主表结构补充。
- `005_competitions_name_key_dedupe.up.sql` / `.down.sql`
  - 竞赛 `name_key` 去重与约束。
- `006_competitions_timeline.up.sql` / `.down.sql`
  - 竞赛时间窗口与时间线字段扩展。
- `init_all.sql`
  - 容器初始化时一次性执行的迁移集合。
  - 被 `docker-compose.yml` 挂载到 PostgreSQL 初始化目录。
- `seed_recognition.sql`
  - 认定规则种子数据。

### 5.17 `unifocus/web/`

- `package.json`
  - 前端技术栈与脚本入口。
  - 核心依赖为 Next 14、React 18、Ant Design、Axios。
- `package-lock.json`
  - 依赖锁定文件。
- `tsconfig.json`
  - TypeScript 编译配置，开启 `strict`。
- `next.config.js`
  - 定义 `/api` 到后端的重写，是竞赛页能直接请求 `/api/v1/competitions` 的关键。
- `next-env.d.ts`
  - Next 自动生成的类型声明文件。
- `app/layout.tsx`
  - 全局根布局，当前仍是脚手架级默认布局。
- `app/globals.css`
  - 全局 reset 与基础 body 样式。
  - 仍依赖未在此处定义的 CSS 变量，样式体系尚未收口。
- `app/dashboard/competitions/page.tsx`
  - 当前唯一真正落地的前端业务页面。
  - 负责竞赛列表、搜索、筛选、新增、编辑、删除。
  - 关联后端 `competition_handler.go`。
  - 当前已修复 `useEffect` 依赖警告，`lint/build` 均已通过。
- `lib/api/client.ts`
  - Axios 统一客户端，负责 token 注入和 401 跳转。
- `lib/api/auth.ts`
  - 认证接口封装，当前无页面使用。
- `lib/api/profile.ts`
  - 用户画像接口封装，当前无页面使用。
- `lib/api/opportunities.ts`
  - 机会接口封装，当前无页面使用。
- `lib/api/metrics.ts`
  - 指标接口封装，当前无页面使用。

### 5.18 `unifocus/nlp-service/`

- `requirements.txt`
  - Python 依赖声明。
  - 当前真实已落地能力少于声明的能力。
- `Dockerfile`
  - NLP 服务容器构建方式。
  - 当前虽可独立构建，但未纳入 `docker-compose.yml`。
- `app/main.py`
  - FastAPI 入口。
  - 注册健康检查、根路径、文本提取路由、向量化路由。
  - 当前阶段不要求启动 `nlp-service`，后续前端界面完善阶段可暂时忽略。
- `app/api/routes/__init__.py`
  - 路由包标记文件。
- `app/api/routes/text_extractor.py`
  - HTML/PDF/DOCX 文本提取接口。
  - 关联 `services/text_extractor.py`。
- `app/api/routes/vectorization.py`
  - 文本向量化接口 `/api/v1/vectorize/text`。
  - 关联 `api/models/vectorization.py` 与 `services/vectorization_service.py`。
- `app/api/models/vectorization.py`
  - 向量化请求与响应模型。
- `app/services/text_extractor.py`
  - 文本提取业务逻辑，基于 BeautifulSoup、pdfminer、python-docx。
- `app/services/vectorization_service.py`
  - 懒加载 `SentenceTransformer` 的向量化服务单例。

### 5.19 当前为空或仅占位的目录

- `unifocus/extension/`
- `unifocus/infrastructure/`
- `unifocus/tests/`

这些目录当前没有承接真实文件，不应在后续讨论中被表述成“已经存在可用实现”。

## 6. 后续执行规则

后续每次开始改动前，先做这几件事：

1. 先回读根目录 `agent.md`，确认当前轮次的完整主文档与协作规则。
2. 再回读 `docs/project/agent-execution-rules.md`，确认技能调用与执行细则。
3. 回读 `docs/project/current-status-and-direction.md`，确认当前任务批次和下一阶段方向。
4. 回读 `docs/project/baseline.md` 与 `docs/project/structure.md`，确认真实运行链路、文件落点与预留模块边界。
5. 如果是前端任务，再回读 `docs/frontend/issues.md`，避免把预留前端能力误当作已接入能力。
