"use client";

import { cn } from "@/lib/utils";
import type { ActiveNode } from "./planet-scene";

const GROUP_META: Record<string, {
  label: string;
  sublabel: string;
  accent: string;        // CSS color for border/glow
  textClass: string;     // Tailwind text color
  glowStyle: string;     // inline box-shadow
}> = {
  scholarships: {
    label: "SCHOLARSHIP TRACK",
    sublabel: "奖学金轨道",
    accent: "#e2c471",
    textClass: "text-amber-400",
    glowStyle: "box-shadow:0 0 18px 2px rgba(226,196,113,0.18)",
  },
  competitions: {
    label: "COMPETITION TRACK",
    sublabel: "竞赛轨道",
    accent: "#52d4ff",
    textClass: "text-cyan-400",
    glowStyle: "box-shadow:0 0 18px 2px rgba(82,212,255,0.18)",
  },
  internships: {
    label: "INTERNSHIP TRACK",
    sublabel: "实习轨道",
    accent: "#b46dff",
    textClass: "text-violet-400",
    glowStyle: "box-shadow:0 0 18px 2px rgba(180,109,255,0.18)",
  },
  main: {
    label: "UNIFOCUS CORE",
    sublabel: "主星",
    accent: "#7557ff",
    textClass: "text-indigo-300",
    glowStyle: "box-shadow:0 0 18px 2px rgba(117,87,255,0.18)",
  },
};

export function PlanetInfoPanel({ node }: { node: ActiveNode }) {
  const meta = GROUP_META[node.groupKey] ?? GROUP_META.main;

  return (
    <div
      className="hud-enter pointer-events-none absolute bottom-6 left-6 w-60"
      style={{
        borderRadius: "12px",
        border: `1px solid ${meta.accent}33`,
        background: "linear-gradient(135deg, rgba(4,11,24,0.82) 0%, rgba(8,18,40,0.78) 100%)",
        backdropFilter: "blur(16px) saturate(1.4)",
        boxShadow: `0 0 24px 2px ${meta.accent}18, inset 0 1px 0 rgba(255,255,255,0.05)`,
      }}
    >
      {/* Scan line accent at top */}
      <div
        className="rounded-t-[11px]"
        style={{
          height: "2px",
          background: `linear-gradient(90deg, transparent, ${meta.accent}cc, transparent)`,
        }}
      />

      <div className="px-4 pb-4 pt-3">
        {/* Eyebrow row */}
        <div className="mb-2.5 flex items-center gap-2">
          <span
            className="dot-pulse inline-block h-[5px] w-[5px] rounded-full"
            style={{ background: meta.accent, boxShadow: `0 0 6px 1px ${meta.accent}` }}
          />
          <span
            className={cn("font-mono text-[8px] font-semibold tracking-[0.2em]", meta.textClass)}
          >
            {meta.label}
          </span>
        </div>

        {/* Title */}
        <div className="font-display text-[15px] font-bold leading-tight text-white">
          {node.title}
        </div>

        {/* Sub-label */}
        <div className={cn("mt-0.5 font-mono text-[9px] tracking-widest", meta.textClass, "opacity-60")}>
          {meta.sublabel}
        </div>

        {/* Context */}
        {node.context && (
          <>
            <div
              className="my-3"
              style={{ height: "1px", background: `linear-gradient(90deg, ${meta.accent}44, transparent)` }}
            />
            <div className="font-body text-[11px] leading-relaxed text-white/55">
              {node.context}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
