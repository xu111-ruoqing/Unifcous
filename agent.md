# Agent Guide

最后更新：2026-04-01
用途：这是当前仓库给 AI / agent 协作者使用的完整主文档。后续无论是新会话续接、跨智能体协作、前端修复、结构整理还是代码实现，默认先读本文件，再进入 `docs/` 中的分支文档。

当前仓库主名称：`Unifocus-v1.2`。文档中出现的 `UniFocus` 默认按产品界面或品牌文案理解；`机会寻找助手` 不再作为当前项目主名称使用。

## 1. 入口规则

- 根目录完整主文档：`agent.md`
- 当前批次专项交接：`HANDOFF.md`
  - 当前最直接的执行任务说明在这里
- 根目录兼容入口：`AGENTS.md`
  - 只用于兼容部分会自动扫描 `AGENTS.md` 的工具
  - 权威内容仍以 `agent.md` 为准
- 规则分支文档：`docs/project/agent-execution-rules.md`
  - 路径：`/Users/ruoqingxu/Desktop/Unifocus-v1.2/docs/project/agent-execution-rules.md`
- 跨项目方法论文档：`reusable-build-patterns.md`
  - 只抽象可迁移思路，不复述当前业务事实

## 2. 开工前固定阅读顺序

后续开始任何正式任务前，默认按这个顺序回读：

1. `agent.md`
2. `HANDOFF.md`
3. `docs/project/agent-execution-rules.md`
4. `docs/project/current-status-and-direction.md`
5. `docs/project/project-memory.md`
6. `docs/project/baseline.md`
7. `docs/project/structure.md`
8. 如果是前端任务，再读：
   - `docs/frontend/overview.md`
   - `docs/frontend/issues.md`
   - `docs/frontend/specs/`
   - `docs/frontend/plans/`
   - `docs/frontend/implementations/`

## 3. 必守执行规则

### 3.1 全项目规则

- 只要开始执行本项目中的正式任务，必须先调用 `superwork` agent skill，辅助完成全局理解、风险判断、优先级排序和执行顺序规划
- 如果当前会话没有 `superwork`，必须在第一次进度说明里明确说出，不能伪造“已调用”
- 重要结论不能只留在聊天里，必须回写到文档

### 3.2 前端规则

- 只要任务涉及前端页面、组件、布局、样式、动效、交互、视觉设计或前端重构，必须额外调用 `frontend-design`
- 如果当前会话没有 `frontend-design`，也必须在第一次前端进度说明里明确说出
- 前端任务默认先确认真实路由、真实接口、真实鉴权链路，再决定改代码

### 3.3 汇报与交接规则

跨轮次或跨智能体协作，默认至少交代这四件事：

1. 现在实际实现了什么
2. 当前问题为什么发生
3. 真正触发入口在哪里
4. 当前边界、残留风险和下一步建议

## 4. 当前项目真实状态

### 4.1 项目目标

当前项目的真实目标是逐步形成一个“机会 / 竞赛信息平台”，不是只维护单一竞赛页。

### 4.2 当前能力线

- 仓库根目录：`Unifocus-v1.2/`
- `unifocus/backend/`
  - 当前最完整的业务实现层
- `unifocus/web/`
  - 当前主业务前端
- `unifocus/web-vite/`
  - 并行前端实验线
- `unifocus/nlp-service/`
  - 独立 Python NLP 微服务

### 4.3 当前阶段目标

当前阶段优先级不是全栈闭环，而是：

1. 完善前端界面与页面覆盖
2. 维持主链路稳定可运行
3. 收敛主前端路线
4. 再逐步补齐登录、画像、简历、智能能力

### 4.4 当前已验证主链路

- `postgres`：Docker
- `redis`：Docker
- `api`：Docker
- `web`：本地 Next.js
- `nlp-service`：当前阶段默认不作为阻塞项

## 5. 2026-04-01 前端修复批次总结

这一批修复针对的是 `unifocus/web/` 的 Dashboard / Planet 视图在真实访问时暴露出来的一组链式问题。

### 5.1 当时的用户侧现象

- 端口打开后看不到预期前端页面
- `dashboard` 能进，`planet` 崩溃
- 出现 `Cannot find module './310.js'`
- 出现 `opportunities.slice` 对 `null` 调用的 runtime error

### 5.2 已确认的根因

- 根路由没有实际页面，用户直接打开 `/` 时看不到 Dashboard
- `401` 会被统一 API client 跳转到不存在的 `/login`
- 侧边栏里存在多个尚未实现的死链接
- 同时存在多个 `next dev` 进程，且开发 / 构建产物混写，导致 `.next` 产物损坏
- 机会数据在某些响应里返回 `null`，前端却直接按数组处理

### 5.3 已完成的修复

- 新增 `/ -> /dashboard` 重定向
- 禁用或隐藏还未实现的侧边栏死链接
- 取消 `401 -> /login` 的错误跳转
- 清理损坏的 `.next` 产物并隔离 dev / build 输出目录
- 在机会数据层与展示层都补了空数组兜底
- 修通 `dashboard` 与 `planet` 的可访问性

### 5.4 这批修复后的默认结论

- 对 `unifocus/web/` 的问题，不能只看代码，还必须同时看：
  - 文档定义的入口
  - 浏览器访问的真实路由
  - 本机正在运行的 Next 进程
  - `.next` 产物是否被混写
- 用户报告的第二个错误，通常不是推翻第一次修复，而是说明已经进入下一层问题

## 6. 前端修复协作模式

后续如果再出现“页面打不开 / 能进 Dashboard 不能进 Planet / 运行时报错 / 和文档不一致”这类问题，默认按这个顺序排查：

1. 先读文档
   - 先看 `agent.md`
   - 再看 `docs/frontend/issues.md`
   - 再看相关 `specs/`、`plans/`、`implementations/`
2. 再核对真实入口
   - 用户实际打开的是 `/`、`/dashboard` 还是 `/dashboard/planet`
   - 文档写的入口和代码里真正存在的页面是否一致
3. 再核对鉴权与路由跳转
   - 是否存在 `401`
   - 是否跳到了不存在的页面
   - 是否是导航把用户带到了未实现路由
4. 再核对运行态
   - 是否有多个 `next dev` 进程
   - 是否把 `next build` 和 `next dev` 的产物写进同一个目录
   - `.next` 是否损坏
5. 再核对数据形状
   - 不要默认接口一定返回数组
   - 对 `null`、空对象、异常字段都要做防御
6. 最后才做整理和去重
   - 先恢复可访问
   - 再收口重复请求、重复兜底和冗余代码

## 7. 跨智能体前端交接模板

后续如果一个智能体修到一半要交给另一个智能体，至少要留下这些信息：

- 用户看到的实际现象
- 已确认的根因，不要只写猜测
- 已改动的文件
- 已验证通过的访问地址或命令
- 还没解决的风险点
- 下一位智能体首先应该读哪几份文档

最低交接格式建议：

1. 现象
2. 根因
3. 已修复项
4. 验证结果
5. 剩余风险
6. 下一步

## 8. 文档分层

- `agent.md`
  - 完整主文档，给协作者先读
- `docs/project/agent-execution-rules.md`
  - 技能调用与执行规则分支
- `docs/project/current-status-and-direction.md`
  - 当前阶段续接状态
- `docs/project/project-memory.md`
  - 长期有效的项目记忆
- `docs/project/baseline.md`
  - 项目真实基线
- `docs/project/structure.md`
  - 结构、入口、文档落点
- `reusable-build-patterns.md`
  - 跨项目复用本仓库方法论时优先阅读
- `docs/frontend/overview.md`
  - 前端整体发展情况
- `docs/frontend/issues.md`
  - 前端问题边界与修复策略
- `docs/frontend/specs/`
  - 前端设计方案
- `docs/frontend/plans/`
  - 前端计划文档
- `docs/frontend/implementations/`
  - 前端实现记录

## 9. 文档维护规则

- 根目录允许保留 `agent.md`、`AGENTS.md`、`HANDOFF.md`、`reusable-build-patterns.md` 这类总入口或总方法论文档
- 其余项目 Markdown 默认放在顶层 `docs/`
- 如果目录结构变了，至少同步更新：
  - `agent.md`
  - `AGENTS.md`
  - `docs/README.md`
  - `docs/project/agent-execution-rules.md`
  - `docs/project/structure.md`
- 如果阶段目标、前端修复经验或协作规则变化了，也要同步回写 `agent.md`
