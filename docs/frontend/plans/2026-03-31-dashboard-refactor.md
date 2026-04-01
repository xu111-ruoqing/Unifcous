# Dashboard 重构实现计划

> **For agentic workers:** Before executing this plan, follow `docs/project/agent-execution-rules.md`. All project work must use `superwork`, and frontend work must additionally use `frontend-design`.

**Goal:** 将 UniFocus Dashboard 从 antd 迁移到 shadcn/ui，实现设计文档中的深色侧边栏 + Hero Banner + 双列内容布局。

**Architecture:** 基于 Next.js 14 App Router，使用 shadcn/ui 组件库 + Tailwind CSS 构建。Dashboard 采用固定侧边栏布局，主内容区包含渐变 Hero Banner（统计卡片）和双列内容网格（机会列表 + 轨道摘要）。数据通过现有 `lib/api/` 层获取。

**Tech Stack:** Next.js 14, React 18, shadcn/ui, Tailwind CSS, TypeScript, Lucide Icons

**Design Reference:** `.superpowers/brainstorm/11941-1774938409/content/dashboard-hybrid-v1.html`
**Design Spec:** `docs/frontend/specs/2026-03-31-dashboard-refactor-design.md`

---

## File Structure

```
unifocus/web/
├── app/
│   ├── layout.tsx                          # 修改：添加字体、全局样式
│   ├── globals.css                         # 修改：Tailwind directives + CSS 变量
│   ├── dashboard/
│   │   ├── layout.tsx                      # 新建：Dashboard 布局（侧边栏 + 主区域）
│   │   ├── page.tsx                        # 新建：Dashboard 总览页（原 /dashboard 首页）
│   │   └── competitions/
│   │       └── page.tsx                    # 保留：竞赛管理页（暂不改动）
├── components/
│   ├── ui/                                 # shadcn/ui 组件（自动生成）
│   │   ├── card.tsx
│   │   ├── badge.tsx
│   │   ├── tabs.tsx
│   │   └── button.tsx
│   ├── dashboard/
│   │   ├── sidebar.tsx                     # 新建：深色侧边栏
│   │   ├── hero-banner.tsx                 # 新建：Hero Banner + 统计卡片
│   │   ├── opportunity-feed.tsx            # 新建：近期机会列表（含 Tab 筛选）
│   │   └── track-summary.tsx              # 新建：轨道摘要卡片组
├── lib/
│   └── utils.ts                            # 新建：shadcn cn() 工具函数
├── tailwind.config.ts                      # 新建：Tailwind 配置 + 自定义色彩
├── postcss.config.js                       # 新建：PostCSS 配置
├── components.json                         # 新建：shadcn/ui 配置
```

---

## Task 1: 初始化 Tailwind CSS 和 shadcn/ui

**Files:**
- Create: `unifocus/web/tailwind.config.ts`
- Create: `unifocus/web/postcss.config.js`
- Create: `unifocus/web/components.json`
- Create: `unifocus/web/lib/utils.ts`
- Modify: `unifocus/web/app/globals.css`
- Modify: `unifocus/web/app/layout.tsx`
- Modify: `unifocus/web/package.json`（通过 npm install）

- [ ] **Step 1: 安装 shadcn/ui 依赖**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
npm install tailwindcss-animate class-variance-authority clsx tailwind-merge lucide-react
```

- [ ] **Step 2: 创建 PostCSS 配置**

Create `unifocus/web/postcss.config.js`:

```js
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
```

- [ ] **Step 3: 创建 Tailwind 配置（含自定义色彩体系）**

Create `unifocus/web/tailwind.config.ts`:

```ts
import type { Config } from "tailwindcss";
import tailwindcssAnimate from "tailwindcss-animate";

const config: Config = {
  darkMode: ["class"],
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // 页面背景
        page: "#f4f5fb",
        // 侧边栏
        sidebar: {
          bg: "#0f0d2e",
          accent: "#6c5ce7",
        },
        // 品牌强调色
        accent: {
          indigo: "#6c5ce7",
          cyan: "#00cec9",
          amber: "#e8a838",
          violet: "#a855f7",
          rose: "#f43f5e",
        },
        // Banner 渐变
        banner: {
          from: "#12103a",
          via: "#2d1b69",
          to: "#6c5ce7",
        },
        // 文字色
        text: {
          primary: "#12103a",
          secondary: "#7b7a94",
          muted: "#aeadc4",
        },
        // 卡片
        card: {
          border: "#eae8f5",
        },
      },
      fontFamily: {
        display: [
          "Playfair Display",
          "Songti SC",
          "Georgia",
          "serif",
        ],
        body: ["DM Sans", "PingFang SC", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"],
      },
      borderRadius: {
        sm: "8px",
        md: "14px",
        lg: "20px",
      },
    },
  },
  plugins: [tailwindcssAnimate],
};

export default config;
```

- [ ] **Step 4: 创建 shadcn cn() 工具**

Create `unifocus/web/lib/utils.ts`:

```ts
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
```

- [ ] **Step 5: 创建 shadcn/ui 配置文件**

Create `unifocus/web/components.json`:

```json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "new-york",
  "rsc": true,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.ts",
    "css": "app/globals.css",
    "baseColor": "slate",
    "cssVariables": true,
    "prefix": ""
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils",
    "ui": "@/components/ui",
    "lib": "@/lib",
    "hooks": "@/hooks"
  }
}
```

- [ ] **Step 6: 更新 globals.css — 添加 Tailwind directives 和 CSS 变量**

Replace `unifocus/web/app/globals.css` with:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 240 10% 15%;
    --card: 0 0% 100%;
    --card-foreground: 240 10% 15%;
    --popover: 0 0% 100%;
    --popover-foreground: 240 10% 15%;
    --primary: 252 75% 60%;
    --primary-foreground: 0 0% 100%;
    --secondary: 240 5% 96%;
    --secondary-foreground: 240 10% 15%;
    --muted: 240 5% 96%;
    --muted-foreground: 240 4% 55%;
    --accent: 240 5% 96%;
    --accent-foreground: 240 10% 15%;
    --destructive: 0 84% 60%;
    --destructive-foreground: 0 0% 100%;
    --border: 240 6% 90%;
    --input: 240 6% 90%;
    --ring: 252 75% 60%;
    --radius: 0.5rem;
  }
}

@layer base {
  * {
    @apply border-border;
  }
  body {
    @apply bg-page font-body text-text-primary antialiased;
  }
}
```

- [ ] **Step 7: 更新根 layout.tsx — 添加 Google Fonts**

Replace `unifocus/web/app/layout.tsx` with:

```tsx
import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "UniFocus — Opportunity Planet",
  description: "统一展示竞赛、奖学金与实践机会",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="zh-CN">
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@600;700;800&family=DM+Sans:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap"
          rel="stylesheet"
        />
      </head>
      <body>{children}</body>
    </html>
  );
}
```

- [ ] **Step 8: 安装 shadcn/ui 组件**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
npx shadcn@latest add card badge tabs button --yes
```

- [ ] **Step 9: 验证构建**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
npx next build 2>&1 | tail -20
```

Expected: Build succeeds (可能有 lint warning，但无 error）。

- [ ] **Step 10: Commit**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
git add tailwind.config.ts postcss.config.js components.json lib/utils.ts app/globals.css app/layout.tsx package.json package-lock.json components/ui/
git commit -m "feat(web): initialize shadcn/ui with Tailwind CSS and custom color tokens"
```

---

## Task 2: 侧边栏组件

**Files:**
- Create: `unifocus/web/components/dashboard/sidebar.tsx`

- [ ] **Step 1: 创建侧边栏组件**

Create `unifocus/web/components/dashboard/sidebar.tsx`:

```tsx
"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  Trophy,
  List,
  User,
  FileUp,
  Settings,
} from "lucide-react";

const navItems = [
  { href: "/dashboard", label: "机会总览", icon: LayoutDashboard },
  { href: "/dashboard/competitions", label: "竞赛管理", icon: Trophy },
  { href: "/dashboard/opportunities", label: "机会列表", icon: List },
  { href: "/dashboard/profile", label: "用户画像", icon: User },
  { href: "/dashboard/resume", label: "简历上传", icon: FileUp },
];

const systemItems = [
  { href: "/dashboard/settings", label: "设置", icon: Settings },
];

interface StatusIndicatorProps {
  online: boolean;
  label: string;
}

function StatusIndicator({ online, label }: StatusIndicatorProps) {
  return (
    <div className="flex items-center gap-2 py-1 text-[11px] text-white/35">
      <span
        className={cn(
          "h-1.5 w-1.5 shrink-0 rounded-full",
          online
            ? "bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.5)]"
            : "bg-amber-400 shadow-[0_0_8px_rgba(251,191,36,0.4)]"
        )}
      />
      {label}
    </div>
  );
}

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="fixed inset-y-0 left-0 z-10 flex w-[240px] flex-col overflow-y-auto bg-sidebar-bg px-4 pb-5 pt-7">
      {/* Brand */}
      <div className="mb-9 flex items-center gap-3 px-2">
        <div className="flex h-9 w-9 items-center justify-center rounded-[10px] bg-gradient-to-br from-accent-indigo to-accent-cyan text-sm font-bold text-white shadow-[0_4px_16px_rgba(108,92,231,0.35)] font-display">
          U
        </div>
        <div>
          <div className="font-display text-lg font-bold tracking-tight text-white">
            UniFocus
          </div>
          <div className="font-mono text-[10px] tracking-widest text-white/35">
            opportunity planet
          </div>
        </div>
      </div>

      {/* Nav Section */}
      <div className="mb-2 px-3 text-[9px] font-semibold uppercase tracking-[0.16em] text-white/[0.22]">
        导航
      </div>
      <nav className="flex flex-col gap-0.5">
        {navItems.map((item) => {
          const isActive =
            item.href === "/dashboard"
              ? pathname === "/dashboard"
              : pathname.startsWith(item.href);
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "relative flex items-center gap-3 rounded-sm px-3 py-2.5 text-[13px] font-medium transition-all duration-200",
                isActive
                  ? "bg-accent-indigo/25 text-white"
                  : "text-white/[0.48] hover:bg-white/[0.06] hover:text-white/75"
              )}
            >
              {isActive && (
                <span className="absolute -left-4 top-1/2 h-5 w-[3px] -translate-y-1/2 rounded-r-sm bg-accent-indigo" />
              )}
              <item.icon
                className={cn("h-4 w-4", isActive ? "opacity-100" : "opacity-70")}
              />
              {item.label}
            </Link>
          );
        })}
      </nav>

      {/* System Section */}
      <div className="mb-2 mt-7 px-3 text-[9px] font-semibold uppercase tracking-[0.16em] text-white/[0.22]">
        系统
      </div>
      <nav className="flex flex-col gap-0.5">
        {systemItems.map((item) => {
          const isActive = pathname.startsWith(item.href);
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "relative flex items-center gap-3 rounded-sm px-3 py-2.5 text-[13px] font-medium transition-all duration-200",
                isActive
                  ? "bg-accent-indigo/25 text-white"
                  : "text-white/[0.48] hover:bg-white/[0.06] hover:text-white/75"
              )}
            >
              <item.icon
                className={cn("h-4 w-4", isActive ? "opacity-100" : "opacity-70")}
              />
              {item.label}
            </Link>
          );
        })}
      </nav>

      {/* Footer */}
      <div className="mt-auto border-t border-white/[0.06] px-2 pt-4">
        <StatusIndicator online={true} label="竞赛接口在线" />
        <StatusIndicator online={false} label="机会接口回落" />

        <div className="mt-3.5 flex gap-1 rounded-sm bg-white/[0.06] p-[3px]">
          <button className="flex-1 rounded-[6px] bg-gradient-to-br from-accent-indigo to-purple-500 py-[7px] text-center font-body text-[11px] font-semibold text-white shadow-[0_2px_10px_rgba(108,92,231,0.4)]">
            Dashboard
          </button>
          <button className="flex-1 rounded-[6px] bg-transparent py-[7px] text-center font-body text-[11px] font-semibold text-white/40">
            Planet
          </button>
        </div>
      </div>
    </aside>
  );
}
```

- [ ] **Step 2: Commit**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
git add components/dashboard/sidebar.tsx
git commit -m "feat(web): add dark sidebar navigation component"
```

---

## Task 3: Dashboard 布局

**Files:**
- Create: `unifocus/web/app/dashboard/layout.tsx`

- [ ] **Step 1: 创建 Dashboard 布局**

Create `unifocus/web/app/dashboard/layout.tsx`:

```tsx
import { Sidebar } from "@/components/dashboard/sidebar";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <main className="ml-[240px] flex-1 px-7 py-6">{children}</main>
    </div>
  );
}
```

- [ ] **Step 2: 验证构建**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
npx next build 2>&1 | tail -20
```

Expected: Build succeeds.

- [ ] **Step 3: Commit**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
git add app/dashboard/layout.tsx
git commit -m "feat(web): add dashboard layout with sidebar"
```

---

## Task 4: Hero Banner 组件

**Files:**
- Create: `unifocus/web/components/dashboard/hero-banner.tsx`

- [ ] **Step 1: 创建 Hero Banner 组件**

Create `unifocus/web/components/dashboard/hero-banner.tsx`:

```tsx
import { cn } from "@/lib/utils";

interface StatItem {
  label: string;
  value: number | string;
  detail: string;
  colorClass?: string;
}

interface HeroBannerProps {
  stats: StatItem[];
}

export function HeroBanner({ stats }: HeroBannerProps) {
  const today = new Date();
  const dateStr = `${today.getFullYear()}.${String(today.getMonth() + 1).padStart(2, "0")}.${String(today.getDate()).padStart(2, "0")}`;

  return (
    <section className="relative mb-5 overflow-hidden rounded-lg bg-gradient-to-br from-banner-from via-banner-via to-banner-to px-8 py-7">
      {/* Decorative circles */}
      <div className="pointer-events-none absolute -right-[60px] -top-[60px] h-[240px] w-[240px] rounded-full bg-[radial-gradient(circle,rgba(108,92,231,0.3),transparent_70%)]" />
      <div className="pointer-events-none absolute bottom-[-80px] right-20 h-[180px] w-[180px] rounded-full bg-[radial-gradient(circle,rgba(0,206,201,0.15),transparent_70%)]" />

      {/* Top row */}
      <div className="relative z-[1] flex items-start justify-between">
        <div>
          <div className="font-mono text-[10px] uppercase tracking-[0.18em] text-white/40">
            Opportunity Overview
          </div>
          <h1 className="mt-1.5 font-display text-[28px] font-bold tracking-tight text-white">
            机会总览
          </h1>
          <p className="mt-1 max-w-[360px] text-[13px] leading-relaxed text-white/50">
            统一展示竞赛、奖学金与实践机会，让每一条成长路径清晰可见。
          </p>
        </div>
        <div className="text-right font-mono text-[11px] text-white/30">
          {dateStr}
          <br />
          Spring Semester
        </div>
      </div>

      {/* Stats grid */}
      <div className="relative z-[1] mt-[22px] grid grid-cols-4 gap-3">
        {stats.map((stat) => (
          <div
            key={stat.label}
            className="rounded-md border border-white/[0.08] bg-white/[0.07] px-4 py-3.5 backdrop-blur-sm transition-all duration-200 hover:-translate-y-0.5 hover:bg-white/[0.11]"
          >
            <div className="text-[10px] font-semibold uppercase tracking-[0.1em] text-white/45">
              {stat.label}
            </div>
            <div
              className={cn(
                "mt-1 font-display text-[30px] font-extrabold leading-none",
                stat.colorClass || "text-white"
              )}
            >
              {stat.value}
            </div>
            <div className="mt-1 font-mono text-[10px] text-white/35">
              {stat.detail}
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
```

- [ ] **Step 2: Commit**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
git add components/dashboard/hero-banner.tsx
git commit -m "feat(web): add hero banner with gradient and stat cards"
```

---

## Task 5: 近期机会列表组件

**Files:**
- Create: `unifocus/web/components/dashboard/opportunity-feed.tsx`

- [ ] **Step 1: 创建机会列表组件**

Create `unifocus/web/components/dashboard/opportunity-feed.tsx`:

```tsx
"use client";

import { useState } from "react";
import { cn } from "@/lib/utils";

type OpportunityType = "competition" | "scholarship" | "internship";

interface OpportunityItem {
  id: number;
  name: string;
  meta: string;
  type: OpportunityType;
  badge: string;
  time: string;
}

interface OpportunityFeedProps {
  items: OpportunityItem[];
}

const typeConfig: Record<
  OpportunityType,
  { indicator: string; badgeBg: string; badgeText: string }
> = {
  competition: {
    indicator: "bg-accent-cyan",
    badgeBg: "bg-[#eef2ff]",
    badgeText: "text-[#4f46e5]",
  },
  scholarship: {
    indicator: "bg-accent-amber",
    badgeBg: "bg-[#fef7ec]",
    badgeText: "text-[#b45309]",
  },
  internship: {
    indicator: "bg-accent-violet",
    badgeBg: "bg-[#f3ecff]",
    badgeText: "text-[#7c3aed]",
  },
};

const tabs = [
  { key: "all", label: "全部" },
  { key: "competition", label: "竞赛" },
  { key: "scholarship", label: "奖学金" },
  { key: "internship", label: "实习" },
] as const;

type TabKey = (typeof tabs)[number]["key"];

export function OpportunityFeed({ items }: OpportunityFeedProps) {
  const [activeTab, setActiveTab] = useState<TabKey>("all");

  const filtered =
    activeTab === "all"
      ? items
      : items.filter((item) => item.type === activeTab);

  return (
    <div className="rounded-lg border border-card-border bg-white p-5 shadow-[0_1px_4px_rgba(0,0,0,0.03)]">
      {/* Header */}
      <div className="mb-4 flex items-center justify-between">
        <span className="text-[15px] font-bold tracking-tight text-text-primary">
          近期机会
        </span>
        <span className="flex cursor-pointer items-center gap-1 text-[11px] font-semibold text-accent-indigo transition-opacity hover:opacity-70">
          查看全部 →
        </span>
      </div>

      {/* Tabs */}
      <div className="mb-3.5 flex w-fit gap-0.5 rounded-sm bg-[#f4f3fa] p-[3px]">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={cn(
              "rounded-[6px] px-4 py-1.5 font-body text-xs font-semibold transition-all duration-150",
              activeTab === tab.key
                ? "bg-white text-text-primary shadow-[0_1px_3px_rgba(0,0,0,0.06)]"
                : "text-text-muted"
            )}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* List */}
      <div className="flex flex-col gap-1.5">
        {filtered.map((item) => {
          const config = typeConfig[item.type];
          return (
            <div
              key={item.id}
              className="flex cursor-pointer items-center gap-3 rounded-sm border border-transparent bg-[#fafaff] px-3.5 py-3 transition-all duration-[180ms] hover:translate-x-1 hover:border-card-border hover:bg-white hover:shadow-[0_2px_8px_rgba(108,92,231,0.06)]"
            >
              <div
                className={cn(
                  "h-8 w-1 shrink-0 rounded-full",
                  config.indicator
                )}
              />
              <div className="min-w-0 flex-1">
                <div className="truncate text-[13px] font-semibold text-text-primary">
                  {item.name}
                </div>
                <div className="mt-0.5 text-[11px] text-text-secondary">
                  {item.meta}
                </div>
              </div>
              <span
                className={cn(
                  "shrink-0 rounded-[6px] px-2.5 py-[3px] text-[10px] font-semibold",
                  config.badgeBg,
                  config.badgeText
                )}
              >
                {item.badge}
              </span>
              <span className="shrink-0 font-mono text-[11px] text-text-muted">
                {item.time}
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Commit**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
git add components/dashboard/opportunity-feed.tsx
git commit -m "feat(web): add opportunity feed component with tab filtering"
```

---

## Task 6: 轨道摘要卡片组件

**Files:**
- Create: `unifocus/web/components/dashboard/track-summary.tsx`

- [ ] **Step 1: 创建轨道摘要卡片组件**

Create `unifocus/web/components/dashboard/track-summary.tsx`:

```tsx
import { cn } from "@/lib/utils";

type TrackColor = "gold" | "violet" | "cyan";

interface TrackItem {
  name: string;
  hint: string;
}

interface TrackCardData {
  color: TrackColor;
  label: string;
  items: TrackItem[];
}

interface TrackSummaryProps {
  tracks: TrackCardData[];
}

const colorConfig: Record<
  TrackColor,
  { stripe: string; label: string }
> = {
  gold: {
    stripe: "bg-gradient-to-r from-accent-amber to-[#f6d365]",
    label: "text-accent-amber",
  },
  violet: {
    stripe: "bg-gradient-to-r from-accent-violet to-[#c084fc]",
    label: "text-accent-violet",
  },
  cyan: {
    stripe: "bg-gradient-to-r from-accent-cyan to-[#67e8f9]",
    label: "text-accent-cyan",
  },
};

export function TrackSummary({ tracks }: TrackSummaryProps) {
  return (
    <div className="flex flex-col gap-4">
      {tracks.map((track) => {
        const config = colorConfig[track.color];
        return (
          <div
            key={track.label}
            className="relative overflow-hidden rounded-lg border border-card-border bg-white px-5 py-[18px] shadow-[0_1px_4px_rgba(0,0,0,0.03)]"
          >
            {/* Top stripe */}
            <div
              className={cn(
                "absolute inset-x-0 top-0 h-[3px]",
                config.stripe
              )}
            />

            <div
              className={cn(
                "mb-3 text-[10px] font-bold uppercase tracking-[0.12em]",
                config.label
              )}
            >
              {track.label}
            </div>

            {track.items.map((item, i) => (
              <div
                key={item.name}
                className={cn(
                  "flex items-center justify-between py-[7px] text-xs text-text-primary",
                  i < track.items.length - 1 && "border-b border-[#f5f4fa]"
                )}
              >
                <span>{item.name}</span>
                <span className="font-mono text-[10px] text-text-muted">
                  {item.hint}
                </span>
              </div>
            ))}
          </div>
        );
      })}
    </div>
  );
}
```

- [ ] **Step 2: Commit**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
git add components/dashboard/track-summary.tsx
git commit -m "feat(web): add track summary cards component"
```

---

## Task 7: Dashboard 总览页 — 组装所有组件

**Files:**
- Create: `unifocus/web/app/dashboard/page.tsx`

- [ ] **Step 1: 创建 Dashboard 总览页**

Create `unifocus/web/app/dashboard/page.tsx`:

```tsx
"use client";

import { useEffect, useState } from "react";
import { HeroBanner } from "@/components/dashboard/hero-banner";
import { OpportunityFeed } from "@/components/dashboard/opportunity-feed";
import { TrackSummary } from "@/components/dashboard/track-summary";

// 静态轨道数据（后续可从 API 获取）
const trackData = [
  {
    color: "gold" as const,
    label: "奖学金轨道",
    items: [
      { name: "国家奖学金", hint: "高势能" },
      { name: "一等奖学金", hint: "绩点联动" },
      { name: "励志奖学金", hint: "成长加权" },
    ],
  },
  {
    color: "violet" as const,
    label: "实习轨道",
    items: [
      { name: "互联网大厂", hint: "高竞争" },
      { name: "产品经理", hint: "策略双栖" },
      { name: "研究助理", hint: "科研延展" },
    ],
  },
  {
    color: "cyan" as const,
    label: "竞赛热度",
    items: [
      { name: "国家级A*类", hint: "12 项" },
      { name: "国家级A类", hint: "18 项" },
      { name: "国家级B类", hint: "24 项" },
    ],
  },
];

interface Competition {
  id: number;
  name: string;
  level: string;
  typical_time_window: string;
  timeline_hint: string;
  note: string;
}

export default function DashboardPage() {
  const [competitions, setCompetitions] = useState<Competition[]>([]);

  useEffect(() => {
    let active = true;
    fetch("/api/v1/competitions?limit=500")
      .then((res) => res.json())
      .then((data) => {
        if (active && Array.isArray(data)) {
          setCompetitions(data);
        }
      })
      .catch(() => {});
    return () => {
      active = false;
    };
  }, []);

  // Derive stats from real data
  const totalCompetitions = competitions.length;
  const nationalCount = competitions.filter((c) =>
    c.level?.startsWith("国家级")
  ).length;

  const stats = [
    {
      label: "竞赛总量",
      value: totalCompetitions,
      detail: `${nationalCount} 国家级`,
    },
    {
      label: "奖学金",
      value: 5,
      detail: "荣誉轨道",
      colorClass: "text-accent-amber",
    },
    {
      label: "实习实践",
      value: 5,
      detail: "跃迁轨道",
      colorClass: "text-accent-violet",
    },
    {
      label: "机会节点",
      value: totalCompetitions + 10,
      detail: "星图映射",
      colorClass: "text-accent-cyan",
    },
  ];

  // Map competitions to opportunity feed items
  const opportunityItems = competitions.slice(0, 10).map((c) => ({
    id: c.id,
    name: c.name,
    meta: c.timeline_hint || c.note || "",
    type: "competition" as const,
    badge: c.level || "竞赛",
    time: c.typical_time_window || "",
  }));

  // Add some static scholarship/internship items for visual completeness
  const feedItems = [
    ...opportunityItems,
    {
      id: -1,
      name: "国家奖学金",
      meta: "高势能荣誉识别",
      type: "scholarship" as const,
      badge: "奖学金",
      time: "9–10月",
    },
    {
      id: -2,
      name: "互联网大厂实习",
      meta: "高竞争密度岗位",
      type: "internship" as const,
      badge: "实习",
      time: "春招",
    },
  ];

  return (
    <>
      <HeroBanner stats={stats} />
      <div className="grid grid-cols-[1.4fr_0.6fr] gap-4">
        <OpportunityFeed items={feedItems} />
        <TrackSummary tracks={trackData} />
      </div>
    </>
  );
}
```

- [ ] **Step 2: 验证构建**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
npx next build 2>&1 | tail -20
```

Expected: Build succeeds.

- [ ] **Step 3: Commit**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
git add app/dashboard/page.tsx
git commit -m "feat(web): add dashboard overview page with live competition data"
```

---

## Task 8: 视觉微调和 Grain Overlay

**Files:**
- Modify: `unifocus/web/app/globals.css`

- [ ] **Step 1: 添加 grain overlay 和微调样式**

在 `unifocus/web/app/globals.css` 末尾追加：

```css
/* Grain texture overlay */
body::after {
  content: '';
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 999;
  opacity: 0.018;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
}
```

- [ ] **Step 2: 验证 dev server 效果**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
npx next build 2>&1 | tail -10
```

Expected: Build succeeds.

- [ ] **Step 3: Commit**

```bash
cd /Users/ruoqingxu/Desktop/Unifocus-v1.2/unifocus/web
git add app/globals.css
git commit -m "style(web): add grain texture overlay for visual depth"
```

---

## Summary

| Task | Component | Description |
|------|-----------|-------------|
| 1 | 基础设施 | 初始化 Tailwind + shadcn/ui + 色彩体系 |
| 2 | Sidebar | 深色固定侧边栏 + 导航 + 状态指示 |
| 3 | Layout | Dashboard 布局（侧边栏 + 主内容区） |
| 4 | Hero Banner | 渐变 Banner + 四列统计卡片 |
| 5 | Opportunity Feed | 近期机会列表 + Tab 分类筛选 |
| 6 | Track Summary | 轨道摘要卡片组（奖学金/实习/竞赛） |
| 7 | Dashboard Page | 组装所有组件 + API 数据集成 |
| 8 | Visual Polish | Grain overlay 纹理 |
