// Planet data layer — ported and enhanced from web-vite/src/lib/opportunityData.js
import type { OpportunityScoreResult } from "@/lib/api/opportunities";
import { arrayOrEmpty } from "@/lib/utils";

export interface OrbitItem {
  title: string;
  shortName: string; // 新增：用于 3D 渲染的 Canvas 卫星文字标签
  accent: "gold" | "cyan" | "violet";
  context: string;
  window?: string;
  score?: number;
  passed?: boolean;
  gateReason?: string;
  daysUntilDeadline?: number;
  scoreComponents?: OpportunityScoreResult["components"];
  /** 0-1: engagement weight based on view_count + save_count */
  weight: number;
  /** 0-1: time proximity (1 = deadline very soon, 0 = far away or no deadline) */
  timeProximity: number;
  /** 点击跳转的落地页 URL */
  url?: string;
}

export interface OrbitGroup {
  key: "scholarships" | "competitions" | "internships";
  label: string;
  english: string;
  tone: "gold" | "cyan" | "violet";
  heroLabel: string;
  summary: string;
  items: OrbitItem[];
  count: number;
}

export interface PlanetSizes {
  scholarships: number;
  competitions: number;
  internships: number;
}

export interface OverloadWarning {
  month: string;
  count: number;
  items: string[];
}

// ─── Static reference data ───────────────────────────────────────────────────

const REFERENCE_ITEMS: Record<string, OrbitItem[]> = {
  scholarships: [
    { title: "国家奖学金", shortName: "国奖", accent: "gold", context: "高势能荣誉识别", weight: 0.9, timeProximity: 0.8 },
    { title: "一等奖学金", shortName: "一等奖", accent: "gold", context: "高绩点与成果联动", weight: 0.7, timeProximity: 0.6 },
    { title: "励志奖学金", shortName: "励志", accent: "gold", context: "成长轨迹加权", weight: 0.5, timeProximity: 0.5 },
    { title: "国际级荣誉", shortName: "国际", accent: "gold", context: "高曝光通道", weight: 0.6, timeProximity: 0.3 },
    { title: "校级荣誉", shortName: "校级", accent: "gold", context: "校内认定资源", weight: 0.3, timeProximity: 0.2 },
  ],
  competitions: [
    { title: "中国国际大学生创新大赛", shortName: "互联网+", accent: "cyan", context: "互联网+ 主轨道", weight: 0.9, timeProximity: 0.7 },
    { title: "挑战杯", shortName: "挑战杯", accent: "cyan", context: "学术科技作品", weight: 0.8, timeProximity: 0.6 },
    { title: "全国数学建模竞赛", shortName: "国赛", accent: "cyan", context: "高频能力入口", weight: 0.7, timeProximity: 0.9 },
    { title: "美国大学生数学建模竞赛", shortName: "美赛", accent: "cyan", context: "国际竞赛映射", weight: 0.6, timeProximity: 0.4 },
    { title: "创新创业", shortName: "大创", accent: "cyan", context: "产品实践联动", weight: 0.5, timeProximity: 0.3 },
  ],
  internships: [
    { title: "互联网大厂", shortName: "大厂", accent: "violet", context: "高竞争密度岗位", weight: 0.9, timeProximity: 0.8 },
    { title: "产品经理", shortName: "产品", accent: "violet", context: "策略与产品双栖", weight: 0.7, timeProximity: 0.6 },
    { title: "研究助理", shortName: "RA", accent: "violet", context: "科研路径延展", weight: 0.6, timeProximity: 0.5 },
    { title: "科研院所", shortName: "科研", accent: "violet", context: "学研协同路径", weight: 0.5, timeProximity: 0.4 },
    { title: "运营增长", shortName: "运营", accent: "violet", context: "业务落地训练", weight: 0.4, timeProximity: 0.3 },
  ],
};

const ORBIT_CONFIG: Record<string, Omit<OrbitGroup, "items" | "count">> = {
  scholarships: {
    key: "scholarships",
    label: "奖学金",
    english: "Scholarships",
    tone: "gold",
    heroLabel: "学业荣誉轨道",
    summary: "以奖学金和荣誉认定为主，强调成绩、持续性与公开证明。",
  },
  competitions: {
    key: "competitions",
    label: "学科竞赛",
    english: "Competitions",
    tone: "cyan",
    heroLabel: "竞赛成长轨道",
    summary: "由竞赛级别、赛道与时间窗口构成，是最可结构化的机会面。",
  },
  internships: {
    key: "internships",
    label: "实习实践",
    english: "Internships",
    tone: "violet",
    heroLabel: "实践跃迁轨道",
    summary: "偏向岗位和研究实践，适合将能力标签映射成职业试炼。",
  },
};

// ─── Helpers ──────────────────────────────────────────────────────────────────

function compact(value: unknown, fallback: string): string {
  const str = `${value ?? ""}`.trim();
  return str || fallback;
}

function inferGroup(item: { type?: string; title?: string }): "scholarships" | "competitions" | "internships" {
  const haystack = `${item?.type ?? ""} ${item?.title ?? ""}`.toLowerCase();
  if (haystack.includes("奖") || haystack.includes("scholar")) return "scholarships";
  if (haystack.includes("实习") || haystack.includes("intern") || haystack.includes("岗位") || haystack.includes("研究")) return "internships";
  return "competitions";
}

function computeTimeProximity(deadline?: string | null): number {
  if (!deadline) return 0.1;
  const ms = new Date(deadline).getTime() - Date.now();
  if (ms <= 0) return 0; // already passed
  const days = ms / 86_400_000;
  return Math.max(0.05, Math.min(1, 1 - days / 180));
}

function normalizeWeights(items: Array<{ rawWeight: number } & Omit<OrbitItem, "weight">>): OrbitItem[] {
  const max = Math.max(...items.map((i) => i.rawWeight), 1);
  return items.map(({ rawWeight, ...rest }) => ({
    ...rest,
    weight: max > 0 ? rawWeight / max : 0.5,
  }));
}

// 智能截取简称：优先读取 tags，若无则使用正则过滤并提取前 4 个字符
function extractShortName(title?: string, tags?: string[]): string {
  if (tags && tags.length > 0 && tags[0]) {
    return tags[0].trim().substring(0, 4);
  }
  const raw = compact(title, "未命名");
  if (raw.includes("中国国际大学生创新大赛")) return "互联网+";
  if (raw.includes("挑战杯")) return "挑战杯";
  // 移除所有括号及其内部内容
  const cleaned = raw.replace(/[(（][^)）]*[)）]/g, "").trim();
  return cleaned.substring(0, 4);
}

// ─── Public API ───────────────────────────────────────────────────────────────

export interface RawCompetition {
  id?: number;
  name?: string;
  level?: string;
  typical_time_window?: string;
  timeline_hint?: string;
  note?: string;
  official_url?: string;
}

export interface RawOpportunity {
  id?: number;
  title?: string;
  type?: string;
  organizer?: string;
  deadline?: string | null;
  start_date?: string | null;
  view_count?: number;
  save_count?: number;
  tags?: string[];
  source_url?: string;
  score_result?: OpportunityScoreResult;
}

export function buildOrbitGroups({
  competitions = [],
  opportunities = [],
}: {
  competitions?: RawCompetition[] | null;
  opportunities?: RawOpportunity[] | null;
}): OrbitGroup[] {
  const safeCompetitions = arrayOrEmpty(competitions);
  const safeOpportunities = arrayOrEmpty(opportunities);
  const byGroup: Record<string, Array<{ rawWeight: number } & Omit<OrbitItem, "weight">>> = {
    scholarships: [],
    competitions: [],
    internships: [],
  };

  safeCompetitions.slice(0, 15).forEach((c) => {
    byGroup.competitions.push({
      title: compact(c.name, "未命名竞赛"),
      shortName: extractShortName(c.name),
      accent: "cyan",
      context: compact(c.level, "国家级竞赛"),
      window: compact(c.typical_time_window, ""),
      rawWeight: 0.6,
      timeProximity: 0.4,
      url: c.official_url, 
    });
  });

  safeOpportunities.slice(0, 30).forEach((item) => {
    const group = inferGroup(item);
    const scoreResult = item.score_result;
    const score = scoreResult?.score;
    const rawWeight = score != null ? score : (item.view_count ?? 0) + (item.save_count ?? 0);
    byGroup[group].push({
      title: compact(item.title, "未命名机会"),
      shortName: extractShortName(item.title, item.tags),
      accent: ORBIT_CONFIG[group].tone,
      context: compact(item.organizer ?? item.type, "机会信号待确认"),
      window: compact(item.deadline ?? item.start_date, "开放中"),
      score,
      passed: scoreResult?.passed,
      gateReason: scoreResult?.explanation?.gate_reason,
      daysUntilDeadline: scoreResult?.explanation?.days_until_deadline,
      scoreComponents: scoreResult?.components,
      rawWeight,
      timeProximity: computeTimeProximity(item.deadline),
      url: item.source_url, 
    });
  });

  return Object.values(ORBIT_CONFIG).map((config) => {
    const raw = byGroup[config.key];
    const items: OrbitItem[] = raw.length
      ? normalizeWeights(raw).slice(0, 7)
      : REFERENCE_ITEMS[config.key];

    return { ...config, items, count: items.length };
  });
}

export function computePlanetSizes(groups: OrbitGroup[]): PlanetSizes {
  const counts = Object.fromEntries(groups.map((g) => [g.key, g.count]));
  const max = Math.max(...Object.values(counts), 1);
  const scale = (n: number) => 0.85 + ((n / max) * 0.4);
  return {
    scholarships: scale(counts.scholarships ?? 1),
    competitions: scale(counts.competitions ?? 1),
    internships: scale(counts.internships ?? 1),
  };
}

export function computeOverload(opportunities: RawOpportunity[] | null | undefined): OverloadWarning[] {
  const safeOpportunities = arrayOrEmpty(opportunities);
  const buckets: Record<string, string[]> = {};
  safeOpportunities.forEach((item) => {
    const dl = item.deadline ?? item.start_date;
    if (!dl) return;
    const d = new Date(dl);
    if (isNaN(d.getTime())) return;
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`;
    if (!buckets[key]) buckets[key] = [];
    buckets[key].push(compact(item.title, "未命名"));
  });

  return Object.entries(buckets)
    .filter(([, items]) => items.length >= 4)
    .map(([month, items]) => ({ month, count: items.length, items: items.slice(0, 5) }))
    .sort((a, b) => a.month.localeCompare(b.month));
}
