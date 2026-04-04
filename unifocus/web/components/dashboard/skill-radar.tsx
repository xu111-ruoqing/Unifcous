"use client";

import { useState } from "react";
import {
  RadarChart,
  Radar,
  PolarGrid,
  PolarAngleAxis,
  ResponsiveContainer,
} from "recharts";
import { cn } from "@/lib/utils";

export interface RadarDimension {
  subject: string;
  /** 0-100: user's current level */
  current: number;
  /** 0-100: opportunity requirement */
  required: number;
}

export function buildRadarData(
  skills: string[],
  certificates: Array<{ name: string; score?: number }>,
  tags: string[]
): RadarDimension[] {
  const FIXED_DIMS = ["编程开发", "数据分析", "产品设计", "研究写作", "团队协作", "国际视野"];

  return FIXED_DIMS.map((dim) => {
    const hasSkill = skills.some((s) => s.includes(dim.slice(0, 2)));
    const certBonus = certificates.filter((c) => c.name.includes(dim.slice(0, 2))).length * 10;
    const current = Math.min(100, (hasSkill ? 62 : 22) + certBonus);
    const tagHits = tags.filter((t) => t.includes(dim.slice(0, 2))).length;
    const required = Math.min(100, 50 + tagHits * 12);
    return { subject: dim, current, required };
  });
}

// Compute the largest gap for the summary badge
function maxGap(data: RadarDimension[]): { subject: string; gap: number } {
  return data.reduce(
    (max, d) => {
      const gap = d.required - d.current;
      return gap > max.gap ? { subject: d.subject, gap } : max;
    },
    { subject: "", gap: -Infinity }
  );
}

export function SkillRadar({ data, className }: { data: RadarDimension[]; className?: string }) {
  const [expanded, setExpanded] = useState(true);
  const gap = maxGap(data);

  return (
    <div
      className={cn("hud-enter-right absolute bottom-6 right-6 w-[260px]", className)}
      style={{
        borderRadius: "12px",
        border: "1px solid rgba(108,92,231,0.25)",
        // 与主面板保持一致的深空渐变和玻璃态参数
        background: "linear-gradient(135deg, rgba(4,11,24,0.82) 0%, rgba(8,18,40,0.78) 100%)",
        backdropFilter: "blur(16px) saturate(1.4)",
        boxShadow: "0 0 24px 2px rgba(108,92,231,0.18), inset 0 1px 0 rgba(255,255,255,0.05)",
      }}
    >
      {/* Top scan line */}
      <div
        className="rounded-t-[11px]"
        style={{ height: "2px", background: "linear-gradient(90deg, transparent, rgba(108,92,231,0.8), rgba(168,85,247,0.8), transparent)" }}
      />

      {/* Header */}
      <button
        onClick={() => setExpanded((v) => !v)}
        className="group flex w-full items-center justify-between px-4 py-3 text-left"
      >
        <div>
          <div className="font-mono text-[8px] font-semibold tracking-[0.2em] text-indigo-400">
            SKILLS DELTA
          </div>
          <div className="font-display text-[13px] font-bold text-white">能力缺口诊断</div>
        </div>
        <div className="flex items-center gap-2">
          {/* Gap badge */}
          {gap.gap > 0 && (
            <span
              className="rounded-md px-1.5 py-0.5 font-mono text-[8px] font-semibold"
              style={{
                background: "rgba(168,85,247,0.15)",
                color: "#c084fc",
                border: "1px solid rgba(168,85,247,0.25)",
              }}
            >
              ⬆ {gap.subject}
            </span>
          )}
          <svg
            className="text-white/30 transition-transform duration-300 group-hover:text-white/60"
            style={{ transform: expanded ? "rotate(180deg)" : "rotate(0deg)" }}
            width="14" height="14" viewBox="0 0 14 14" fill="none"
          >
            <path d="M3 5l4 4 4-4" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
        </div>
      </button>

      {/* Chart */}
      {expanded && (
        <div className="px-1 pb-3">
          <ResponsiveContainer width="100%" height={178}>
            <RadarChart data={data} margin={{ top: 8, right: 20, bottom: 4, left: 20 }}>
              {/* 增强网格线的深空质感 */}
              <PolarGrid stroke="rgba(108,92,231,0.25)" />
              <PolarAngleAxis
                dataKey="subject"
                tick={{ fill: "rgba(255,255,255,0.55)", fontSize: 8.5, fontFamily: "JetBrains Mono, monospace" }}
              />
              {/* Opportunity requirement — violet dashed outline */}
              <Radar
                name="机会要求"
                dataKey="required"
                stroke="#a855f7"
                fill="#a855f7"
                fillOpacity={0.12}
                strokeDasharray="3 2"
                strokeWidth={1.5}
                dot={false}
              />
              {/* Current ability — indigo solid fill */}
              <Radar
                name="当前能力"
                dataKey="current"
                stroke="#6c5ce7"
                fill="#6c5ce7"
                fillOpacity={0.35}
                strokeWidth={1.8}
                dot={{ r: 2.5, fill: "#6c5ce7", strokeWidth: 0 }}
              />
            </RadarChart>
          </ResponsiveContainer>

          {/* Legend row */}
          <div className="flex items-center justify-center gap-4 pb-1 pt-0.5">
            <span className="flex items-center gap-1.5 font-mono text-[8.5px] text-white/40">
              <span className="inline-block h-px w-5 border-t border-dashed border-violet-400" />
              机会要求
            </span>
            <span className="flex items-center gap-1.5 font-mono text-[8.5px] text-white/40">
              <span className="inline-block h-px w-5 bg-indigo-400" />
              当前能力
            </span>
          </div>
        </div>
      )}
    </div>
  );
}