# 当前状态与下一阶段方向

最后更新：2026-04-01
用途：这是当前对话阶段的续接文档。后续继续这个项目时，先读根目录 `agent.md`，再读本文件和相关分支文档。

当前仓库主名称：`Unifocus-v1.2`。`unifocus/` 当前应理解为源码主目录、服务命名与历史数据库前缀，而不是完整仓库名。

## 1. 本轮对话已完成的内容

### 1.1 项目基线与结构文档已建立

已建立并持续维护这些文档：

- `docs/project/project-memory.md`
- `docs/project/baseline.md`
- `docs/project/structure.md`
- `docs/frontend/issues.md`

这些文档现在承担的职责分别是：

- `project-memory.md`：长期决策、环境特例、边界规则、遗留问题
- `baseline.md`：项目真实现状、运行链路、关键偏差
- `structure.md`：目录结构、入口文件、预留目录边界
- `issues.md`：前端问题清单、已修复项、下一阶段前端方向

### 1.2 前端环境与当前页面已梳理

已确认：

- 当前 `unifocus/web/` 已存在这些真实入口：
  - `app/page.tsx`：`/ -> /dashboard`
  - `app/dashboard/page.tsx`：Dashboard 总览页
  - `app/dashboard/competitions/page.tsx`：竞赛管理页
  - `app/dashboard/planet/page.tsx`：Planet 页面
- 前端代码层面的 `lint` 和 `build` 在本轮修复后已经修通
- 已补 `.eslintrc.json`
- 已修复竞赛页的 `useEffect` 依赖警告
- 已修复 `401` 错跳不存在 `/login`、根路由缺页、侧边栏死链接和 `planet` 的空数据崩溃

补充说明：

- 当前本机缓存目录和 `node_modules` 可能已被人为清理，用于减小体积
- 这不改变代码现状，但后续复验 `web/` 前应先执行 `npm install`

### 1.3 并行 Vite 前端实验入口仍保留

本轮新增：

- `unifocus/web-vite/`

它的定位是：

- 在不打断现有 `Next.js web/` 基线的前提下，并行提供新的前端实验入口
- 技术栈为 `React + Vite + JavaScript + three.js`
- 当前已落地“机会星球”双风格界面：
  - `Dashboard`：偏深色、偏 PDF 封面语气的运营视图
  - `Planet`：偏参考图语气的沉浸式星球视图

注意：

- `web-vite/` 当前是并行入口，不等于已经正式替换 `web/`
- 当前前端主工作对象已经回到 `unifocus/web/` 的 Dashboard / Planet 链路
- `web-vite/` 当前更适合作为视觉参考和实验线，而不是默认执行入口
- 后续如果决定迁移，需要单独判断登录、路由、API client 和部署方式

### 1.4 项目主链路已验证

已确认当前可运行主链路为：

- `postgres`：Docker 容器
- `redis`：Docker 容器
- `api`：Docker 容器
- `web`：本地 Next.js 开发服务

当前不要求启动：

- `nlp-service`

### 1.5 后端启动相关结论已确认

已确认这些关键事实：

- `backend/Dockerfile` 原先 Go 版本落后，已对齐到 `1.25`
- Docker 数据库默认角色原本是 `unifocus`
- 目前为了兼容本地开发配置，数据库里已经存在 `yufan` 角色
- 但 `yufan` 不是单一语义，它在不同层里含义不同

### 1.6 主前端修复批次已完成一轮

本轮已额外确认并处理 `unifocus/web/` 的一组真实访问问题：

- 根路由没有页面，已补 `/ -> /dashboard`
- `401` 会错误跳到不存在的 `/login`，已取消该跳转
- 侧边栏存在死链接，已做禁用或收敛
- 多个 `next dev` 进程与混写产物导致 `.next` 损坏，已清理并隔离 dev / build 目录
- `planet` 页面对 `null` 机会数据直接调用数组方法，已补兜底

这意味着后续再接手 `web/` 前端问题时，不能只看代码，还必须同时检查：

- 用户打开的真实路由
- Next 进程状态
- `.next` 产物目录状态
- 接口返回数据形状

## 2. 需要持续记住的关键结论

### 2.1 `Unifocus-v1.2`、`unifocus` 与 `yufan` 不是同一层概念

当前项目里：

- `Unifocus-v1.2` 是当前仓库 / 工作区主名称
- `unifocus` 是源码主目录名、默认数据库角色、服务命名前缀
- `yufan` 在不同地方可能代表：
  - 历史机器路径中的本机用户名
  - 数据库连接账号

不要把它们直接理解成“网站业务登录用户”。

网站登录逻辑在当前代码里仍然是：

- `username`
- `email`
- `password`

这套业务用户体系。

### 2.2 当前阶段忽略 `nlp-service`

本轮已经明确：

- `nlp-service` 当前不是前端界面完善阶段的阻塞项
- 后续如果任务聚焦前端页面完善，可先无视 `nlp-service`
- 只有当任务转向简历解析、画像提取、向量化接线时，再把它拉回主任务

### 2.3 当前前端不是完整产品

当前前端状态必须始终按这个口径表述：

- 已完成：竞赛管理页
- 未完成：首页、登录页、机会页、画像页、完整导航、完整鉴权闭环

不要把 API 封装存在，表述成对应页面已经实现。

## 3. 下一阶段方向

当前已经明确的下一阶段目标是：

继续完善前端界面，当前主工作对象优先是 `unifocus/web/` 的 Planet 体验与视觉。

建议推进顺序：

1. 判断新前端路线
   - 默认继续把 `web/` 作为主工作对象
   - `web-vite/` 作为视觉参考与实验场继续保留
2. 登录闭环
   - 补真实登录页
   - 明确 token、401 跳转和前端 API client
3. 业务页面补齐
   - 机会列表页
   - 机会详情页
   - 用户画像页
   - 简历上传页
4. 请求层收口
   - 让星球视图逐步更多使用真实后端数据

## 4. 后续执行时的默认规则

后续开始任何新任务前，先做：

1. 阅读 `agent.md`
2. 阅读 `docs/project/agent-execution-rules.md`
3. 阅读 `docs/project/current-status-and-direction.md`
4. 阅读 `docs/project/project-memory.md`
5. 阅读 `docs/project/baseline.md`
6. 阅读 `docs/project/structure.md`
7. 如果任务是前端，再阅读 `docs/frontend/issues.md`

执行过程中：

- 优先沿用当前项目现有技术栈和工程结构
- 汇报时先给明确结论
- 区分“当前现状”和“规划目标”
- 重要阶段结果持续回写文档

## 5. 当前等待状态

当前文档已按 `Unifocus-v1.2` 仓库名称和现有前端入口重新同步。

下一步默认沿 `unifocus/web/` 的 Planet 主线继续推进，`web-vite/` 主要作为参考与实验线保留。
