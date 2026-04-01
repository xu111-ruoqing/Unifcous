# 项目文档索引

最后更新：2026-04-01
用途：这是当前仓库 Markdown 文档的统一索引。根目录 `agent.md` 是给智能体协作者使用的完整主文档；根目录 `AGENTS.md` 仅作兼容入口；`HANDOFF.md` 与 `reusable-build-patterns.md` 作为根层专项文档保留；其余项目文档默认保存在 `docs/` 下。

当前仓库主名称：`Unifocus-v1.2`。如果部分前端文档仍出现 `UniFocus`，默认理解为产品界面品牌，而不是仓库名称。

## 1. 默认阅读顺序

### 1.1 所有任务通用

1. `agent.md`
2. `docs/project/agent-execution-rules.md`
3. `docs/project/current-status-and-direction.md`
4. `docs/project/project-memory.md`
5. `docs/project/baseline.md`
6. `docs/project/structure.md`
7. `docs/README.md`

### 1.2 前端任务加读

1. `docs/frontend/overview.md`
2. `docs/frontend/issues.md`
3. `docs/frontend/specs/`
4. `docs/frontend/plans/`
5. `docs/frontend/implementations/`

### 1.3 启动与环境任务加读

1. `docs/setup/getting-started.md`

### 1.4 协作方式参考

1. `docs/templates/cross-project-collaboration-template.md`

### 1.5 跨项目方法论参考

1. `reusable-build-patterns.md`

## 2. 当前文档结构

```text
docs/
├── README.md
├── project/
│   ├── agent-execution-rules.md
│   ├── project-memory.md
│   ├── current-status-and-direction.md
│   ├── baseline.md
│   └── structure.md
├── setup/
│   └── getting-started.md
├── frontend/
│   ├── overview.md
│   ├── issues.md
│   ├── specs/
│   ├── plans/
│   └── implementations/
└── templates/
    └── cross-project-collaboration-template.md
```

## 3. 归档规则

- 根目录完整主文档固定为 `agent.md`
- 根目录 `AGENTS.md` 只作为兼容入口，不承载完整项目事实
- 根目录 `HANDOFF.md` 作为当前专项交接文档
- 根目录 `reusable-build-patterns.md` 作为跨项目可复用方法论文档
- 智能体协作与技能调用规则放 `docs/project/agent-execution-rules.md`
- 项目长期状态、边界、结构类文档放 `docs/project/`
- 跨项目可复用方法论文档放根目录 `reusable-build-patterns.md`，避免和当前项目事实混写
- 启动、部署、环境说明放 `docs/setup/`
- 前端现状、方案、实现记录放 `docs/frontend/`
- 通用模板放 `docs/templates/`
- 新增 Markdown 时，不再写入 `readme/`、`unifocus/web/docs/` 或其他子项目私有 docs 目录

## 4. 维护规则

- 如果目录结构变了，至少同步更新 `agent.md`、`AGENTS.md`、`docs/README.md`、`docs/project/agent-execution-rules.md`、`docs/project/structure.md`
- 如果项目现状或阶段目标变了，至少同步更新 `docs/project/project-memory.md` 和相关专项文档
- 如果新增某个前端阶段性实现记录，优先放到 `docs/frontend/implementations/`
