import { cn } from "@/lib/utils";

export type TrackColor = "gold" | "violet" | "cyan";

export interface TrackItem {
  name: string;
  hint: string;
}

export interface TrackCardData {
  color: TrackColor;
  label: string;
  items: TrackItem[];
}

const colorConfig: Record<
  TrackColor,
  { bar: string; labelColor: string; hintColor: string }
> = {
  gold: {
    bar: "bg-[#C8880A]",
    labelColor: "text-[#956500]",
    hintColor: "text-[#C8880A]",
  },
  violet: {
    bar: "bg-[#3E5672]",
    labelColor: "text-[#3E5672]",
    hintColor: "text-[#6880A0]",
  },
  cyan: {
    bar: "bg-[#3A7555]",
    labelColor: "text-[#2D6048]",
    hintColor: "text-[#3A7555]",
  },
};

export function TrackSummary({ tracks }: { tracks: TrackCardData[] }) {
  return (
    <div className="flex flex-col gap-3.5">
      {tracks.map((track) => {
        const config = colorConfig[track.color];
        return (
          <div
            key={track.label}
            className="relative overflow-hidden rounded-lg border border-card-border bg-[#FEFCF9] px-5 py-[18px] shadow-[0_1px_3px_rgba(30,18,8,0.03)]"
          >
            {/* Left accent bar */}
            <div className={cn("absolute left-0 top-0 h-full w-[3px]", config.bar)} />

            {/* Label */}
            <div
              className={cn(
                "mb-3 font-mono text-[9px] font-semibold uppercase tracking-[0.16em]",
                config.labelColor
              )}
            >
              {track.label}
            </div>

            {track.items.map((item, i) => (
              <div
                key={item.name}
                className={cn(
                  "flex items-center justify-between py-[7px] text-[12px] text-text-primary",
                  i < track.items.length - 1 && "border-b border-[#EDE8E0]"
                )}
              >
                <span className="font-medium">{item.name}</span>
                <span className={cn("font-mono text-[10px]", config.hintColor)}>
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
