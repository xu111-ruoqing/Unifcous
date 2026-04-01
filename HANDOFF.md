# Frontend Planet Handoff

最后更新：2026-04-01
用途：这是给下一位协作者的当前批次专项交接文档。它只服务于“继续完善 `unifocus/web` 的 Planet 页面”这件事，不替代根目录 `agent.md`。

当前仓库主名称：`Unifocus-v1.2`。本文中的 `UniFocus` 或 “机会星球” 默认按产品界面命名理解。

## 1. 先读什么

开始动手前，按这个顺序读取：

1. `agent.md`
2. `HANDOFF.md`
3. `docs/project/agent-execution-rules.md`
4. `docs/project/current-status-and-direction.md`
5. `docs/frontend/issues.md`
6. `docs/frontend/implementations/2026-04-01-planet-view-implementation.md`

如果当前会话支持：

- 先调用 `superwork` 做本轮任务拆解和风险排序
- 因为本轮是前端任务，还必须调用 `frontend-design`

## 2. 当前你接手的是什么

当前主工作对象不是 `web-vite/`，而是：

- `unifocus/web/`
- 目标路由：`/dashboard/planet`

目前这条主线已经完成过一轮修复，已知这些问题曾经出现并已修掉：

- 根路由打不开目标页面
- `401` 错跳不存在的 `/login`
- 侧边栏死链接导致 404
- 多个 `next dev` 进程与混写产物导致 `.next` 损坏
- `planet` 页面对 `null` 数据调用数组方法而崩溃

所以你现在的工作重点不是“让页面先能打开”，而是在“页面已经能打开”的基础上继续把 Planet 做好。

## 3. 这轮最重要的任务

### P0. 优化 Planet 交互

这是第一优先级。

目标不是简单加一点 hover，而是让 `planet` 页面在使用上更顺畅、更有沉浸感。优先关注：

- 星球与子星的悬浮、聚焦、过渡反馈是否自然
- 信息面板出现与切换是否清晰
- 鼠标交互、镜头感、节奏感是否顺滑
- 场景在桌面端与常见笔记本尺寸下是否稳定

要求：

- 不要破坏当前可访问性
- 不要把交互做成喧宾夺主的炫技
- 优先改善真实使用体验，而不是只堆更多浮层

### P0. 提升 Planet 的视觉效果

这也是第一优先级。

当前目标是让界面更有“星球 / 宇宙 / 梦幻机会星图”的感觉，而不是继续停留在“一个可运行的 Three.js 页面”。

视觉方向请优先参考：

- `unifocus/web-vite/public/assets/opportunity-reference.png`
- `unifocus/web-vite/public/assets/plan-cover.png`
- 现有 vite 版星球实现的整体气质

期待方向：

- 更有层次的深空背景
- 更细腻的发光、雾感、粒子和轨道气氛
- 主星与子星之间更明确的主次关系
- 面板、字体、色彩、动画与“机会星球”主题一致

要求：

- 必须按 `frontend-design` 的标准做
- 不要做成通用后台卡片风
- 不要牺牲可读性去换取纯装饰效果

### P1. 子星数量要与后端更完整地匹配

这是第二优先工作，不要忽略。

当前 Planet 的子星数量偏少，和后端真实数据没有完整对齐。这个问题非常可能来自现有前端代码里的主动截断，而不只是视觉布局问题。

优先检查这些位置：

- `unifocus/web/components/dashboard/planet-client.tsx`
  - 当前请求机会列表时用了 `limit: 50`
- `unifocus/web/components/dashboard/planet-data.ts`
  - 当前对 opportunities 先 `slice(0, 50)`
  - 每个 group 最后又只 `slice(0, 7)`

本轮目标：

- 要么让子星尽可能与后端真实机会数量一致
- 要么给出明确、可解释的聚合规则，并在文档里说明为什么不是一条数据一个子星

不接受“因为页面会太挤，所以默认砍掉大部分数据但没有说明”这种状态。

### P2. 置灰的侧边栏项

这是可选项。

如果本轮时间够，且你能把其中某些页面做成真正可用的最小版本，可以推进。

如果做不到：

- 保持当前禁用状态
- 不要为了“点亮菜单”而加空页面或坏路由

## 4. 不要跑偏的边界

本轮默认不要优先做这些事：

- 不要把主线切回 `web-vite/`
- 不要优先补登录闭环
- 不要优先动 `nlp-service`
- 不要因为想做侧边栏，把主要时间从 Planet 主任务上挪走
- 不要回退已经修好的 `/`、`/dashboard`、`/dashboard/planet` 访问链路

## 5. 关键文件

你大概率会频繁改到这些文件：

- `unifocus/web/app/dashboard/planet/page.tsx`
- `unifocus/web/components/dashboard/planet-client.tsx`
- `unifocus/web/components/dashboard/planet-scene.tsx`
- `unifocus/web/components/dashboard/planet-data.ts`
- `unifocus/web/components/dashboard/planet-info-panel.tsx`
- `unifocus/web/components/dashboard/skill-radar.tsx`
- `unifocus/web/components/dashboard/overload-warning.tsx`
- `unifocus/web/components/dashboard/sidebar.tsx`
- `unifocus/web/app/globals.css`
- `unifocus/web/tailwind.config.ts`

参考资料：

- `docs/frontend/implementations/2026-04-01-planet-view-implementation.md`
- `docs/frontend/specs/2026-03-31-dashboard-refactor-design.md`
- `unifocus/web-vite/public/assets/opportunity-reference.png`
- `unifocus/web-vite/public/assets/plan-cover.png`

## 6. 开发规范

默认遵守这些规则：

- 先确认真实路由、真实数据、真实交互问题，再改代码
- 前端任务必须调用 `frontend-design`
- 本项目正式任务必须调用 `superwork`
- 先保稳定，再做重构和美化
- 所有和 Planet 数据数量相关的改动，要以“后端真实返回”为准
- 不要伪造“已完成”，要明确区分：
  - 已实现并验证
  - 已实现但未验证
  - 只是计划

如果本轮需要再次交接给下一位，请至少写清楚：

1. 当前现象
2. 已确认根因
3. 已改文件
4. 已验证结果
5. 剩余风险
6. 下一步建议

## 7. 建议的开工顺序

建议按这个顺序推进：

1. 先跑起当前 `unifocus/web`，确认 `/dashboard/planet` 现状
2. 先看交互问题，再看视觉问题
3. 再检查机会数据数量为什么和后端不一致
4. 只有前两项稳定后，才考虑置灰侧边栏

## 8. 验证清单

最少做这些验证：

- `cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web`
- 如果 `node_modules` 不存在，先执行 `npm install`
- `npm run lint`
- `npm run build`
- 访问：
  - `http://127.0.0.1:3000/`
  - `http://127.0.0.1:3000/dashboard`
  - `http://127.0.0.1:3000/dashboard/planet`

额外必须确认：

- Planet 页面不会因为空数据或异常数据崩溃
- 子星数量和后端数据的对应关系说得清楚
- 视觉增强后性能没有明显退化到不可用

## 9. 完成标准

本轮可以算“完成得不错”的标准是：

- Planet 交互比现在明显更顺滑
- Planet 视觉更接近“梦幻星球”而不是普通深色大屏
- 子星数量问题得到真正处理
- 可选侧边栏项做了就可用，做不了就保持禁用
- 没有把之前修好的访问链路重新弄坏
