# 项目结构说明

最后更新：2026-04-01
更新说明：已将全部项目 Markdown 统一迁移到顶层 `docs/`，根目录 `agent.md` 作为完整主文档，根目录 `AGENTS.md` 作为兼容入口。当前仓库主名称为 `Unifocus-v1.2`。后续任何结构调整都应同步更新本文件和 `docs/README.md`。

## 1. 文档结构

```text
/Users/ruoqingxu/Desktop/Unifocus-v1.2
├── AGENTS.md
├── agent.md
├── HANDOFF.md
├── reusable-build-patterns.md
├── docs/
│   ├── README.md
│   ├── project/
│   │   ├── agent-execution-rules.md
│   │   ├── project-memory.md
│   │   ├── current-status-and-direction.md
│   │   ├── baseline.md
│   │   └── structure.md
│   ├── setup/
│   │   └── getting-started.md
│   ├── frontend/
│   │   ├── overview.md
│   │   ├── issues.md
│   │   ├── specs/
│   │   ├── plans/
│   │   └── implementations/
│   └── templates/
│       └── cross-project-collaboration-template.md
├── 竞赛final.json
└── unifocus/
    ├── .env.example
    ├── Makefile
    ├── docker-compose.yml
    ├── backend/
    ├── web/
    ├── web-vite/
    ├── nlp-service/
    ├── scripts/
    ├── extension/         # 预留目录
    ├── infrastructure/    # 预留目录
    └── tests/             # 预留目录
```

## 2. 核心源码结构

```text
unifocus/
├── backend/
│   ├── cmd/
│   ├── configs/
│   ├── internal/
│   ├── migrations/
│   ├── data/
│   ├── pkg/
│   ├── Dockerfile
│   └── go.mod
├── web/
│   ├── app/
│   ├── components/
│   ├── lib/
│   ├── next.config.js
│   └── package.json
├── web-vite/
│   ├── public/
│   ├── src/
│   ├── index.html
│   └── package.json
├── nlp-service/
│   ├── app/
│   ├── Dockerfile
│   └── requirements.txt
└── scripts/
```

## 3. 目录职责

- `agent.md`
  - 仓库根层完整主文档，给协作者先读
- `AGENTS.md`
  - 兼容部分会自动扫描该文件名的工具，内容应指向 `agent.md`
- `HANDOFF.md`
  - 当前专项任务交接文档，可随批次变化更新
- `reusable-build-patterns.md`
  - 跨项目复用时的方法论抽象，只讲构建思路，不讲当前业务细节
- `Unifocus-v1.2/`
  - 当前仓库根目录名称；`unifocus/` 是源码主目录，不等于完整仓库名
- `docs/`
  - 当前唯一有效的 Markdown 文档根目录
- `docs/project/`
  - 项目长期记忆、当前状态、基线和结构说明
- `docs/project/agent-execution-rules.md`
  - 当前项目对 AI / agent 协作者的统一执行规则与技能调用要求
- `docs/setup/`
  - 启动与环境说明
- `docs/frontend/`
  - 前端现状、问题、方案和实现记录
- `docs/templates/`
  - 通用协作模板
- `unifocus/backend/`
  - 当前最完整的业务实现层，是真正的 API 中心
- `unifocus/web/`
  - 当前主业务前端
- `unifocus/web-vite/`
  - 并行前端实验线
- `unifocus/nlp-service/`
  - 独立 Python NLP 微服务
- `unifocus/scripts/`
  - 本地开发辅助脚本

## 4. 当前真实入口

- Go 后端入口：`unifocus/backend/cmd/api/main.go`
- Go 手动 seed 入口：`unifocus/backend/cmd/seed-competitions/main.go`
- Next 前端布局入口：`unifocus/web/app/layout.tsx`
- Next 前端总览入口：`unifocus/web/app/dashboard/page.tsx`
- Next 前端星球入口：`unifocus/web/app/dashboard/planet/page.tsx`
- Vite 前端入口：`unifocus/web-vite/src/main.jsx`
- NLP 服务入口：`unifocus/nlp-service/app/main.py`
- 本地脚本入口：`unifocus/scripts/start-dev.sh`
- 容器编排入口：`unifocus/docker-compose.yml`

## 5. 文档放置规则

- 根目录允许保留 `agent.md`、`AGENTS.md`、`HANDOFF.md`、`reusable-build-patterns.md` 这类总入口或总方法论文档
- 其余项目 Markdown 默认放在顶层 `docs/`
- 不再向 `readme/`、`unifocus/web/docs/` 或其他子项目目录写入项目 Markdown
- 新的前端设计稿、实现记录、计划文档统一放在 `docs/frontend/`
- 结构变化后，至少同步更新：
  - `agent.md`
  - `AGENTS.md`
  - `docs/README.md`
  - `docs/project/agent-execution-rules.md`
  - `docs/project/structure.md`

## 6. 预留与生成目录提醒

以下目录存在，但当前不应误判为已接入模块：

- `unifocus/extension/`
- `unifocus/infrastructure/`
- `unifocus/tests/`

以下目录属于生成物或依赖目录，不纳入业务结构判断：

- `unifocus/backend/.gocache/`
- `unifocus/web/.next*`
- `unifocus/web/node_modules/`
- `unifocus/nlp-service/.venv/`

## 7. 更新记录

- `2026-03-25`：首次建立结构基线文档。
- `2026-03-30`：补入并行 `web-vite/` 前端结构与入口说明。
- `2026-03-31`：补入 `project-memory.md` 作为长期项目记忆文档。
- `2026-04-01`：统一全部 Markdown 到顶层 `docs/`，并在当日早些时候曾将 `agent.md` 临时收敛为导航文件。
- `2026-04-01`：新增 `agent-execution-rules.md` 作为智能体执行规则唯一权威入口。
- `2026-04-01`：将 `agent.md` 升级为完整主文档，并新增 `AGENTS.md` 兼容入口。
