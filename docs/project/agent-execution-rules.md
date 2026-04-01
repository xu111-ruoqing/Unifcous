# 智能体执行规则

最后更新：2026-04-01
用途：这是当前项目所有 AI / agent 协作者的执行规则分支文档。根目录 `agent.md` 是完整主文档；本文件负责展开技能调用与执行规则细则。

当前仓库主名称：`Unifocus-v1.2`。如果专项文档里出现 `UniFocus`，默认按产品界面或品牌文案理解。

## 1. 适用范围

本规则适用于：

- 代码实现
- 方案设计
- Bug 排查
- 文档整理
- 代码评审
- 前端设计与开发

如果其他专项文档里还保留旧的 agent 协作说明，以本文件为准。

## 2. 强制技能规则

### 2.1 全项目规则：必须调用 `superwork`

只要开始执行本项目中的任何正式任务，必须先调用 `superwork` agent skill，用它辅助完成：

- 全局任务理解
- 影响范围判断
- 风险识别
- 优先级排序
- 执行顺序规划

最低要求：

- 在开始实际改动前，就应把 `superwork` 纳入本轮协作
- 不能在没有实际调用时声称“已使用 `superwork`”
- `superwork` 产出至少应覆盖任务目标、关键风险、执行顺序三个方面

例外处理：

- 如果当前会话没有提供 `superwork`，必须在首次进度说明中明确写出这一点
- 缺少 `superwork` 时，只能按“规则未完全满足的临时例外”继续推进，不能伪造调用事实

### 2.2 前端规则：必须调用 `frontend-design`

只要任务涉及前端范围，必须额外调用 `frontend-design` skill 参与设计与开发。

前端范围包括但不限于：

- 页面
- 组件
- 布局
- 样式
- 动效
- 视觉设计
- 交互设计
- 前端信息架构
- 前端评审与重构

最低要求：

- 在前端方案确定前调用一次 `frontend-design`
- 在前端实现或审查阶段继续受 `frontend-design` 约束
- 不能在没有实际调用时声称“已使用 `frontend-design`”

例外处理：

- 如果当前会话没有提供 `frontend-design`，必须在开始前端任务时明确说明
- 缺少 `frontend-design` 时，本轮前端工作应被标记为“未满足设计协作规则的临时例外”

## 3. 标准执行顺序

后续任何 agent 在本项目内工作，默认按这个顺序执行：

1. 先读 `agent.md`
2. 再读本文件
3. 再读 `docs/project/current-status-and-direction.md`
4. 再进入 `docs/project/` 的状态、基线与结构文档
5. 如果属于前端任务，再进入 `docs/frontend/` 的相关文档

## 4. 首次说明模板

开始工作时，至少要说明这几件事：

1. 当前任务理解
2. 本轮是否已纳入 `superwork`
3. 如果是前端任务，本轮是否已纳入 `frontend-design`
4. 接下来先检查哪一层文档或代码

如果某个强制 skill 当前不可用，也必须在第一条进度说明里明确说出。

## 5. 文档回写规则

如果本轮工作改变了项目结构、协作方式或阶段性结论，至少同步更新这些文档中的相关部分：

- `agent.md`
- `docs/README.md`
- `docs/project/agent-execution-rules.md`
- `docs/project/structure.md`
- 对应专项文档

## 6. 约束优先级

如果文档之间出现冲突，按这个优先级判断：

1. `agent.md`
2. `docs/project/agent-execution-rules.md`
3. `docs/README.md`
4. `docs/project/project-memory.md`
5. 其他专项文档

## 7. 更新记录

- `2026-04-01`：首次建立统一的智能体执行规则，明确 `superwork` 为全项目强制协作技能，`frontend-design` 为前端任务强制协作技能。
