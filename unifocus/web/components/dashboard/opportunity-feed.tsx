"use client";

import { useState } from "react";
import { cn } from "@/lib/utils";

export type OpportunityType = "competition" | "scholarship" | "internship";

export interface OpportunityItem {
  id: number;
  name: string;
  meta: string;
  type: OpportunityType;
  badge: string;
  time: string;
}

const typeConfig: Record<
  OpportunityType,
  { bar: string; badgeBg: string; badgeText: string; badgeRing: string }
> = {
  competition: {
    bar: "bg-[#C4622D]",
    badgeBg: "bg-[#FBF0E9]",
    badgeText: "text-[#C4622D]",
    badgeRing: "ring-[#E8C5B0]",
  },
  scholarship: {
    bar: "bg-[#C8880A]",
    badgeBg: "bg-[#FBF4E0]",
    badgeText: "text-[#956500]",
    badgeRing: "ring-[#E8D49A]",
  },
  internship: {
    bar: "bg-[#3A7555]",
    badgeBg: "bg-[#E8F2EC]",
    badgeText: "text-[#2D6048]",
    badgeRing: "ring-[#A8D4BA]",
  },
};

const tabs = [
  { key: "all", label: "全部" },
  { key: "competition", label: "竞赛" },
  { key: "scholarship", label: "奖学金" },
  { key: "internship", label: "实习" },
] as const;

type TabKey = (typeof tabs)[number]["key"];

export function OpportunityFeed({ items }: { items: OpportunityItem[] }) {
  const [activeTab, setActiveTab] = useState<TabKey>("all");

  const filtered =
    activeTab === "all" ? items : items.filter((item) => item.type === activeTab);

  return (
    <div className="rounded-lg border border-card-border bg-[#FEFCF9] p-5 shadow-[0_1px_3px_rgba(30,18,8,0.04)]">
      {/* Header */}
      <div className="mb-4 flex items-center justify-between">
        <span className="font-display text-[16px] font-semibold italic tracking-tight text-text-primary">
          近期机会
        </span>
        <span className="flex cursor-pointer items-center gap-1 font-mono text-[10px] uppercase tracking-wider text-[#C4622D] transition-opacity hover:opacity-60">
          查看全部 →
        </span>
      </div>

      {/* Tab bar */}
      <div className="mb-4 flex w-fit gap-0.5 rounded bg-[#EDE8DE] p-[3px]">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={cn(
              "rounded-[4px] px-4 py-1.5 text-xs font-medium transition-all duration-150",
              activeTab === tab.key
                ? "bg-white text-text-primary shadow-[0_1px_2px_rgba(30,18,8,0.08)]"
                : "text-text-muted hover:text-text-secondary"
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
              className="flex cursor-pointer items-center gap-3 rounded border border-transparent bg-[#FAF7F2] px-3.5 py-3 transition-all duration-[160ms] hover:translate-x-[3px] hover:border-card-border hover:bg-[#FEFCF9] hover:shadow-[0_2px_8px_rgba(30,18,8,0.05)]"
            >
              <div className={cn("h-7 w-[2px] shrink-0 rounded-full", config.bar)} />
              <div className="min-w-0 flex-1">
                <div className="truncate text-[13px] font-medium text-text-primary">
                  {item.name}
                </div>
                <div className="mt-0.5 text-[11px] text-text-muted">{item.meta}</div>
              </div>
              <span
                className={cn(
                  "shrink-0 rounded px-2 py-[3px] text-[10px] font-medium ring-1",
                  config.badgeBg,
                  config.badgeText,
                  config.badgeRing
                )}
              >
                {item.badge}
              </span>
              <span className="shrink-0 font-mono text-[10px] text-text-muted">
                {item.time}
              </span>
            </div>
          );
        })}
        {filtered.length === 0 && (
          <div className="py-8 text-center font-display text-[14px] italic text-text-muted">
            暂无数据
          </div>
        )}
      </div>
    </div>
  );
}
