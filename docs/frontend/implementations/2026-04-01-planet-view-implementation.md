# UniFocus 机会星球视图 — 实施计划与代码变更记录

**日期：** 2026-04-01
**分支：** dashboard planet view port + visual polish

---

## 一、需求目标

将 `web-vite/` 目录下已有的 Three.js "机会星球"移植并增强到 Next.js `web/` 应用中，功能要求：

- 星球大小/亮度反映机会数量
- 卫星轨道位置基于截止日期远近
- 技能雷达图差距诊断
- 密集截止日期过载预警
- 接入真实后端 API 数据

---

## 二、执行计划（8 Tasks）

| Task | 内容 | 状态 |
|------|------|------|
| T1 | 路由 & 页面骨架（`/dashboard/planet`） | ✅ 完成 |
| T2 | Three.js 数据层（`planet-data.ts`） | ✅ 完成 |
| T3 | Three.js 场景组件（`planet-scene.tsx`） | ✅ 完成 |
| T4 | 客户端编排组件（`planet-client.tsx`） | ✅ 完成 |
| T5 | 技能雷达 HUD 面板（`skill-radar.tsx`） | ✅ 完成 |
| T6 | 过载预警 HUD 面板（`overload-warning.tsx`） | ✅ 完成 |
| T7 | 视觉打磨（遵循 `docs/project/agent-execution-rules.md` 的前端设计协作规则） | ⚠️ 待补全 |
| T8 | Bug 修复 & 代码质量 `/simplify` | ✅ 完成 |

---

## 三、关键技术决策

- **SSR 规避**：`app/dashboard/planet/page.tsx` 用 `next/dynamic + ssr: false` 包裹 `PlanetClient`，防止 WebGL 在服务端执行崩溃
- **Three.js 生命周期**：所有 Three.js 对象在 `useEffect` 的 cleanup 函数中 dispose，包括 geometry、material、texture、renderer
- **字体方案**：Playfair Display（display）/ DM Sans（body）/ JetBrains Mono（mono）
- **HUD 动效**：CSS keyframe（`hud-enter`、`amber-pulse`、`dot-pulse`、`orbit-cw/ccw`）定义在 `globals.css`

---

## 四、新增 / 修改文件清单

### 新增文件

| 文件路径 | 说明 |
|---------|------|
| `components/dashboard/planet-data.ts` | 数据层：OrbitItem/OrbitGroup 类型、buildOrbitGroups、computePlanetSizes、computeOverload |
| `components/dashboard/planet-scene.tsx` | Three.js 场景：星球、卫星轨道、粒子、Raycaster 交互 |
| `components/dashboard/planet-client.tsx` | 客户端编排：fetch 数据、组合 HUD 面板 |
| `components/dashboard/planet-info-panel.tsx` | 悬浮节点信息 HUD 面板 |
| `components/dashboard/skill-radar.tsx` | Recharts 雷达图 HUD 面板 |
| `components/dashboard/overload-warning.tsx` | 过载预警 HUD 面板（amber 动效） |
| `app/dashboard/planet/page.tsx` | Planet 页面路由，dynamic import + 加载态 |

### 修改文件

| 文件路径 | 修改内容 |
|---------|---------|
| `app/globals.css` | 新增 keyframe 动效：`hud-enter`、`amber-pulse`、`dot-pulse`、`orbit-cw`、`orbit-ccw` |
| `tailwind.config.ts` | 新增设计系统 color tokens（sidebar、banner、text、accent 变体） + shadcn CSS 变量映射 |
| `components/dashboard/sidebar.tsx` | 底部模式切换按钮改为 Link，`/dashboard` 和 `/dashboard/planet` 高亮逻辑 |
| `app/dashboard/layout.tsx` | 确认 `ml-[240px]` 偏移配合 sidebar |

---

## 五、已修复 Bug

### CRITICAL：星球页面白屏（`planet-scene.tsx` 灯光初始化）

**位置：** `components/dashboard/planet-scene.tsx` ~ 行 211–215

**原因：** 使用 `Object.assign(new PointLight(...), { position: {x,y,z} })` 将 Three.js 的 `Vector3` 实例替换为普通对象，导致 `renderer.render()` 调用 Vector3 方法时抛出 TypeError，动画循环崩溃，画布空白。

**修复：**
```typescript
// ❌ 错误写法
Object.assign(new PointLight(0xffffff, 11, 45), { position: { x: 0, y: 0, z: 6 } })

// ✅ 正确写法
const keyLight = new PointLight(0xffffff, 11, 45);
keyLight.position.set(0, 0, 6);
```

---

## 六、`/simplify` 代码审查修复记录

### 效率修复

| 问题 | 位置 | 修复 |
|------|------|------|
| Raycaster 每帧 `.map()` 分配新数组 | `planet-scene.tsx:349` | 在 animate 循环外预计算 `satelliteMeshes` |
| opportunities 无上限迭代后再截断 | `planet-data.ts:170` | 改为 `opportunities.slice(0, 50).forEach(...)` |

### 质量修复

| 问题 | 位置 | 修复 |
|------|------|------|
| `Math.random()` 导致雷达数据不确定 | `skill-radar.tsx:31,33` | 替换为基于实际技能数据的确定性计算 |
| 卫星颜色映射使用三元嵌套 | `planet-scene.tsx:290` | 提取为 `GROUP_COLOR` Record |

### 复用修复

| 问题 | 位置 | 修复 |
|------|------|------|
| 裸 `fetch()` 绕过已有 axios 客户端 | `planet-client.tsx:35,40,45` | opportunities 和 profile 改用 `opportunitiesAPI.list()` + `profileAPI.getProfile()`，mounted flag 替换 AbortController |

---

## 七、待办事项

- [ ] **Task 7 视觉打磨**：按 `docs/project/agent-execution-rules.md` 调用 `frontend-design` skill，对面板样式、字体层次、动效细节做审查和改进
- [ ] **competitions API 模块**：`lib/api/competitions.ts` 尚未创建，当前仍用裸 `fetch()`
- [ ] **生产验证**：在真实 API 数据下验证轨道分布和星球大小缩放是否合理

---

## 八、相关参考

- 原始 vite 版实现：`web-vite/src/components/PlanetScene.jsx`
- Dashboard 重构设计稿：`docs/frontend/specs/2026-03-31-dashboard-refactor-design.md`
- shadcn 组件来源：MCP tool `mcp__shadcn__get_component`
- API 客户端：`lib/api/client.ts`（axios，含 token 注入 + 401 重定向）
