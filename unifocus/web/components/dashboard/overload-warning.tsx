"use client";

import { useState } from "react";
import { cn } from "@/lib/utils";
import type { OverloadWarning } from "./planet-data";

interface OverloadWarningPanelProps {
  warnings: OverloadWarning[];
  className?: string;
}

export function OverloadWarningPanel({ warnings, className }: OverloadWarningPanelProps) {
  const [dismissed, setDismissed] = useState(false);

  if (!warnings.length || dismissed) return null;

  return (
    <div
      className={cn("hud-enter-top amber-pulse-ring absolute left-6 top-6 w-56", className)}
      style={{
        borderRadius: "12px",
        border: "1px solid rgba(251,191,36,0.25)",
        // 保持深空暗色调的基础上，匹配主面板的背景参数
        background: "linear-gradient(135deg, rgba(12,8,0,0.82) 0%, rgba(24,14,2,0.78) 100%)",
        backdropFilter: "blur(16px) saturate(1.4)",
        // 补充缺失的泛光层，实现光学辉光对齐
        boxShadow: "0 0 24px 2px rgba(251,191,36,0.15), inset 0 1px 0 rgba(255,255,255,0.05)",
      }}
    >
      {/* Amber scan line */}
      <div
        className="rounded-t-[11px]"
        style={{ height: "2px", background: "linear-gradient(90deg, transparent, rgba(251,191,36,0.85), transparent)" }}
      />

      {/* Header */}
      <div className="flex items-start justify-between px-4 pt-3">
        <div className="flex items-start gap-2.5">
          {/* Pulsing warning icon */}
          <div className="relative mt-0.5 shrink-0">
            <div
              className="dot-pulse h-2 w-2 rounded-full"
              style={{ background: "#fbbf24", boxShadow: "0 0 8px 2px rgba(251,191,36,0.6)" }}
            />
          </div>
          <div>
            <div className="font-mono text-[8px] font-semibold tracking-[0.2em] text-amber-400">
              OVERLOAD RISK
            </div>
            <div className="font-display text-[13px] font-bold text-amber-100">过载风险</div>
          </div>
        </div>
        <button
          onClick={() => setDismissed(true)}
          className="ml-1 mt-0.5 p-0.5 text-amber-500/50 transition-colors hover:text-amber-300"
          aria-label="dismiss"
        >
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
            <path d="M2 2l8 8M10 2l-8 8" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
          </svg>
        </button>
      </div>

      {/* Warning items */}
      <div className="space-y-1.5 px-4 pb-3.5 pt-2.5">
        {warnings.slice(0, 3).map((w) => (
          <div
            key={w.month}
            className="rounded-[8px] px-3 py-2"
            style={{
              background: "rgba(251,191,36,0.08)",
              border: "1px solid rgba(251,191,36,0.15)",
            }}
          >
            <div className="flex items-center justify-between">
              <span className="font-mono text-[9.5px] font-semibold text-amber-300">{w.month}</span>
              <span
                className="rounded-full px-1.5 py-0.5 font-mono text-[8px] font-semibold"
                style={{ background: "rgba(251,191,36,0.18)", color: "#fcd34d" }}
              >
                {w.count} 项
              </span>
            </div>
            <div className="mt-1 font-body text-[10px] leading-relaxed text-amber-100/55">
              {w.items.slice(0, 2).join("、")}
              {w.count > 2 ? ` 等 ${w.count} 项` : ""}
            </div>
          </div>
        ))}

        {/* Divider + advice */}
        <div
          className="my-1"
          style={{ height: "1px", background: "rgba(251,191,36,0.15)" }}
        />
        <p className="font-mono text-[8.5px] leading-relaxed text-amber-400/45">
          机会密集，建议适当减项或错峰规划。
        </p>
      </div>
    </div>
  );
}