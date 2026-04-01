# Dashboard 重构设计文档

日期：2026-03-31
状态：设计中（待 shadcn MCP 配置后继续）

## 1. 目标

对 UniFocus 前端 Dashboard 进行重构，使其更整洁、大气，提升品牌辨识度和信息可读性。

## 2. 技术栈决策

- **UI 组件库**：shadcn/ui（通过 MCP 工具拉取组件）
- **前端框架**：基于现有 Next.js (`web/`) 进行重构
- **设计辅助**：执行本设计稿时，按 `docs/project/agent-execution-rules.md` 调用 `frontend-design` skill 指导视觉设计
- **字体方向**：衬线标题字体（如 Playfair Display）+ 无衬线正文（如 DM Sans）+ 等宽数据字体（如 JetBrains Mono）

## 3. 风格方向：方案 C（浅色主体 + 深色强调区）

用户从三个方案中选择了 C，核心特征：
- 浅色背景保持阅读友好
- 深色渐变区域（Banner / 侧边栏）突出品牌感
- 统计数据在深色区域更醒目
- 整洁大气但不失辨识度

## 4. 排版方案：A + B 混合（以 B 为主体）

用户选择将方案 A 和 B 混合，以 B 为主体：

### 主体结构（来自方案 B）
- **左侧固定深色侧边栏**（~240px）
  - 品牌 Logo + 名称（UniFocus / opportunity planet）
  - 导航菜单：机会总览、竞赛管理、机会列表、用户画像、简历上传、设置
  - 底部：接口状态指示（在线/回落）+ Dashboard/Planet 模式切换
  - 导航 active 态：左侧蓝色竖条 + 半透明高亮背景

### 顶部区域（来自方案 A）
- **Hero Banner**：深色渐变（indigo → violet → purple）
  - 左上：eyebrow 文字（"Opportunity Overview"）+ 大标题（"机会总览"）+ 副标题
  - 右上：日期信息
  - 下方：四列统计卡片（竞赛总量 / 奖学金 / 实习实践 / 机会节点）
  - 装饰元素：右上角半透明圆形光晕

### 内容区域
- **双列布局**：左宽右窄（约 1.4:0.6）
- **左列 - 近期机会卡片**
  - Tab 分类栏：全部 / 竞赛 / 奖学金 / 实习
  - 机会列表行：彩色竖条指示类别 + 名称 + 描述 + 级别 Badge + 时间
  - hover 效果：向右微移 + 阴影
- **右列 - 轨道摘要卡片**
  - 奖学金轨道（金色顶部条纹）
  - 实习轨道（紫色顶部条纹）
  - 竞赛热度（青色顶部条纹）
  - 每张卡片含 3 条关键项 + 简短标签

## 5. 色彩体系

- 页面背景：`#f4f5fb`
- 侧边栏背景：`#0f0d2e`
- 主强调色（Indigo）：`#6c5ce7`
- 青色（竞赛）：`#00cec9`
- 琥珀色（奖学金）：`#e8a838`
- 紫色（实习）：`#a855f7`
- Banner 渐变：`#12103a → #2d1b69 → #6c5ce7`

## 6. 视觉 Mockup 参考

完整高保真 mockup 位于：
`.superpowers/brainstorm/11941-1774938409/content/dashboard-hybrid-v1.html`

## 7. 实现范围

### 在本次重构范围内
- Dashboard 总览页的完整 UI 重构
- 使用 shadcn/ui 组件（Card, Badge, Tabs, Button 等）
- 侧边栏导航布局
- Hero Banner 统计区
- 机会列表展示
- 轨道摘要卡片

### 不在本次范围内
- Planet（three.js 星球）视图重构
- 登录页 / 登录闭环
- 其他业务页面（竞赛管理、机会详情、画像等页面内容）
- 后端 API 变更

## 8. 下一步

1. 重启 Claude Code 会话以加载 shadcn MCP
2. 使用 shadcn MCP 初始化项目并拉取所需组件
3. 基于本设计文档编写实现计划
4. 按计划实现 Dashboard 重构
