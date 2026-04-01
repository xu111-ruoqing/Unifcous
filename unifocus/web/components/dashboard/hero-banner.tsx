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
    <section className="mb-5 rounded-lg border border-card-border bg-[#FEFCF9] px-8 py-7">
      {/* Top row */}
      <div className="flex items-start justify-between">
        {/* Title block with sienna left rule */}
        <div className="flex items-start gap-4">
          <div className="mt-1 h-10 w-[2px] shrink-0 rounded-full bg-[#C4622D]" />
          <div>
            <div className="font-display text-[11px] italic tracking-[0.14em] text-text-muted">
              Opportunity Overview
            </div>
            <h1 className="mt-0.5 font-display text-[30px] font-semibold leading-none tracking-tight text-text-primary">
              机会总览
            </h1>
            <p className="mt-1.5 max-w-[340px] text-[12px] leading-relaxed text-text-secondary">
              统一展示竞赛、奖学金与实践机会，让每一条成长路径清晰可见。
            </p>
          </div>
        </div>

        {/* Date */}
        <div className="text-right">
          <div className="font-display text-[26px] font-normal italic leading-none text-text-muted/60">
            {dateStr}
          </div>
          <div className="mt-1.5 font-mono text-[9px] uppercase tracking-[0.18em] text-text-muted/50">
            Spring Semester
          </div>
        </div>
      </div>

      {/* Divider */}
      <div className="my-6 border-t border-card-border" />

      {/* Stats row */}
      <div className="grid grid-cols-4 divide-x divide-card-border">
        {stats.map((stat, i) => (
          <div key={stat.label} className={cn("px-6", i === 0 && "pl-0", i === stats.length - 1 && "pr-0")}>
            <div
              className={cn(
                "font-display text-[50px] leading-none tracking-tight",
                stat.colorClass ?? "text-text-primary"
              )}
            >
              {stat.value}
            </div>
            <div className="mt-2.5 text-[12px] font-semibold text-text-primary">{stat.label}</div>
            <div className="mt-0.5 font-mono text-[9px] uppercase tracking-wider text-text-muted">
              {stat.detail}
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
