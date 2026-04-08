"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { ExternalLink, Pencil, Plus, RefreshCw, Search, Trash2 } from "lucide-react";
import { HeroBanner } from "@/components/dashboard/hero-banner";
import { cn, arrayOrEmpty } from "@/lib/utils";
import {
  opportunitiesAPI,
  type CreateOpportunityRequest,
  type Opportunity,
  type OpportunityScoreResult,
  type ScoredOpportunity,
} from "@/lib/api/opportunities";
import { profileAPI, type UserProfile } from "@/lib/api/profile";

type OpportunityType = "竞赛" | "奖学金" | "实习" | "科研" | "讲座" | "其他";

interface OpportunityFormState {
  title: string;
  type: OpportunityType;
  description: string;
  source_url: string;
  organizer: string;
  start_date: string;
  deadline: string;
  event_date: string;
  location: string;
  target_majors: string;
  tags: string;
}

const TYPE_OPTIONS: OpportunityType[] = ["竞赛", "奖学金", "实习", "科研", "讲座", "其他"];

const typeStyles: Record<OpportunityType, { badge: string; dot: string }> = {
  "竞赛": { badge: "bg-[#FBF0E9] text-[#C4622D] ring-[#E8C5B0]", dot: "bg-[#C4622D]" },
  "奖学金": { badge: "bg-[#FBF4E0] text-[#956500] ring-[#E8D49A]", dot: "bg-[#C8880A]" },
  "实习": { badge: "bg-[#E8F2EC] text-[#2D6048] ring-[#A8D4BA]", dot: "bg-[#3A7555]" },
  "科研": { badge: "bg-[#EAF0F7] text-[#3E5672] ring-[#B8C7DA]", dot: "bg-[#3E5672]" },
  "讲座": { badge: "bg-[#F0ECEA] text-[#6A5E53] ring-[#D8CEC1]", dot: "bg-[#8B7762]" },
  "其他": { badge: "bg-[#F0ECEA] text-[#6A5E53] ring-[#D8CEC1]", dot: "bg-[#8B7762]" },
};

const emptyForm: OpportunityFormState = {
  title: "",
  type: "竞赛",
  description: "",
  source_url: "",
  organizer: "",
  start_date: "",
  deadline: "",
  event_date: "",
  location: "",
  target_majors: "",
  tags: "",
};

function normalizeType(type?: string): OpportunityType {
  if (!type) return "其他";
  return TYPE_OPTIONS.includes(type as OpportunityType) ? (type as OpportunityType) : "其他";
}

function formatDateLabel(value?: string) {
  if (!value) return "待补充";
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return value;
  return `${parsed.getFullYear()}.${String(parsed.getMonth() + 1).padStart(2, "0")}.${String(parsed.getDate()).padStart(2, "0")}`;
}

function formatDateTimeInput(value?: string) {
  if (!value) return "";
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return "";
  const local = new Date(parsed.getTime() - parsed.getTimezoneOffset() * 60_000);
  return local.toISOString().slice(0, 16);
}

function toISOStringOrUndefined(value: string) {
  return value ? new Date(value).toISOString() : undefined;
}

function splitInput(value: string) {
  return value
    .split(/[,，]/)
    .map((item) => item.trim())
    .filter(Boolean);
}

function scoreTone(score?: number) {
  if (typeof score !== "number") return "text-text-muted";
  if (score >= 75) return "text-[#2D6048]";
  if (score >= 50) return "text-[#956500]";
  return "text-[#B84040]";
}

function buildRequest(form: OpportunityFormState): CreateOpportunityRequest {
  return {
    title: form.title.trim(),
    type: form.type,
    description: form.description.trim(),
    source_url: form.source_url.trim(),
    organizer: form.organizer.trim() || undefined,
    start_date: toISOStringOrUndefined(form.start_date),
    deadline: toISOStringOrUndefined(form.deadline),
    event_date: toISOStringOrUndefined(form.event_date),
    location: form.location.trim() || undefined,
    target_majors: splitInput(form.target_majors),
    tags: splitInput(form.tags),
  };
}

function formFromOpportunity(opportunity: Opportunity): OpportunityFormState {
  return {
    title: opportunity.title ?? "",
    type: normalizeType(opportunity.type),
    description: opportunity.description ?? "",
    source_url: opportunity.source_url ?? "",
    organizer: opportunity.organizer ?? "",
    start_date: formatDateTimeInput(opportunity.start_date),
    deadline: formatDateTimeInput(opportunity.deadline),
    event_date: formatDateTimeInput(opportunity.event_date),
    location: opportunity.location ?? "",
    target_majors: arrayOrEmpty(opportunity.target_majors).join(", "),
    tags: arrayOrEmpty(opportunity.tags).join(", "),
  };
}

function scoreSummary(score?: OpportunityScoreResult) {
  if (!score) return "待评分";
  if (!score.passed) return score.explanation?.gate_reason || "未通过门槛";
  if (score.explanation?.days_until_deadline != null) {
    return `距截止约 ${score.explanation.days_until_deadline} 天`;
  }
  return "已通过门槛";
}

export default function OpportunitiesPage() {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [items, setItems] = useState<ScoredOpportunity[]>([]);
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [selectedDetail, setSelectedDetail] = useState<Opportunity | null>(null);
  const [selectedScore, setSelectedScore] = useState<OpportunityScoreResult | undefined>(undefined);
  const [loading, setLoading] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [search, setSearch] = useState("");
  const [typeFilter, setTypeFilter] = useState<"全部" | OpportunityType>("全部");
  const [sortKey, setSortKey] = useState<"score" | "deadline" | "recent">("score");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<Opportunity | null>(null);
  const [form, setForm] = useState<OpportunityFormState>(emptyForm);
  const [submitting, setSubmitting] = useState(false);
  const [canManage, setCanManage] = useState(false);
  const [errorText, setErrorText] = useState("");

  useEffect(() => {
    if (typeof window !== "undefined") {
      setCanManage(Boolean(window.localStorage.getItem("token")));
    }
  }, []);

  const fetchList = useCallback(async () => {
    setLoading(true);
    setErrorText("");
    try {
      let resolvedProfile: UserProfile | null = null;
      try {
        resolvedProfile = await profileAPI.getProfile();
        setProfile(resolvedProfile);
      } catch {
        setProfile(null);
      }

      try {
        const scored = await opportunitiesAPI.listScored({ limit: 200 }, resolvedProfile?.user_id);
        setItems(scored.data);
        if (!selectedId && scored.data[0]?.id) setSelectedId(scored.data[0].id);
      } catch {
        const fallback = await opportunitiesAPI.list({ limit: 200 });
        setItems(fallback.data);
        if (!selectedId && fallback.data[0]?.id) setSelectedId(fallback.data[0].id);
      }
    } catch {
      setErrorText("机会数据加载失败，请检查后端服务是否已启动。");
    } finally {
      setLoading(false);
    }
  }, [selectedId]);

  useEffect(() => {
    void fetchList();
  }, [fetchList]);

  useEffect(() => {
    if (!selectedId) {
      setSelectedDetail(null);
      setSelectedScore(undefined);
      return;
    }

    let active = true;
    setDetailLoading(true);

    Promise.all([
      opportunitiesAPI.getById(selectedId),
      opportunitiesAPI.getScoreById(selectedId, profile?.user_id).catch(() => null),
    ])
      .then(([detail, scoreRes]) => {
        if (!active) return;
        setSelectedDetail(detail);
        setSelectedScore(scoreRes?.score_result);
      })
      .catch(() => {
        if (!active) return;
        setSelectedDetail(null);
        setSelectedScore(undefined);
      })
      .finally(() => {
        if (active) setDetailLoading(false);
      });

    return () => {
      active = false;
    };
  }, [selectedId, profile?.user_id]);

  const filteredItems = useMemo(() => {
    const keyword = search.trim().toLowerCase();

    const next = items.filter((item) => {
      const normalizedType = normalizeType(item.type);
      const matchesType = typeFilter === "全部" || normalizedType === typeFilter;
      const haystack = `${item.title ?? ""} ${item.organizer ?? ""} ${arrayOrEmpty(item.tags).join(" ")}`.toLowerCase();
      const matchesSearch = !keyword || haystack.includes(keyword);
      return matchesType && matchesSearch;
    });

    next.sort((a, b) => {
      if (sortKey === "score") {
        return (b.score_result?.score ?? -1) - (a.score_result?.score ?? -1);
      }
      if (sortKey === "deadline") {
        const left = a.deadline ? new Date(a.deadline).getTime() : Number.POSITIVE_INFINITY;
        const right = b.deadline ? new Date(b.deadline).getTime() : Number.POSITIVE_INFINITY;
        return left - right;
      }
      return new Date(b.created_at ?? 0).getTime() - new Date(a.created_at ?? 0).getTime();
    });

    return next;
  }, [items, search, sortKey, typeFilter]);

  useEffect(() => {
    if (!filteredItems.length) return;
    if (!selectedId || !filteredItems.some((item) => item.id === selectedId)) {
      setSelectedId(filteredItems[0].id);
    }
  }, [filteredItems, selectedId]);

  const stats = useMemo(() => {
    const passed = items.filter((item) => item.score_result?.passed).length;
    const highScore = items.filter((item) => (item.score_result?.score ?? 0) >= 75).length;
    const internships = items.filter((item) => normalizeType(item.type) === "实习").length;

    return [
      { label: "机会总量", value: items.length, detail: "已接后端真实数据" },
      { label: "高匹配机会", value: highScore, detail: "评分 ≥ 75", colorClass: "text-accent-cyan" },
      { label: "通过门槛", value: passed, detail: "资格初筛通过", colorClass: "text-accent-amber" },
      { label: "实习实践", value: internships, detail: "职业跃迁轨道", colorClass: "text-accent-violet" },
    ];
  }, [items]);

  const selectedType = normalizeType(selectedDetail?.type);
  const selectedStyle = typeStyles[selectedType];

  const openCreateModal = () => {
    setEditingItem(null);
    setForm(emptyForm);
    setIsModalOpen(true);
  };

  const openEditModal = () => {
    if (!selectedDetail) return;
    setEditingItem(selectedDetail);
    setForm(formFromOpportunity(selectedDetail));
    setIsModalOpen(true);
  };

  const closeModal = () => {
    if (submitting) return;
    setIsModalOpen(false);
  };

  const handleSubmit = async () => {
    setSubmitting(true);
    try {
      const payload = buildRequest(form);
      if (editingItem) {
        await opportunitiesAPI.update(editingItem.id, payload);
      } else {
        await opportunitiesAPI.create(payload);
      }
      setIsModalOpen(false);
      await fetchList();
    } catch {
      setErrorText(editingItem ? "机会更新失败，请确认登录状态与字段格式。" : "机会创建失败，请确认登录状态与字段格式。");
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async () => {
    if (!selectedDetail) return;
    if (!window.confirm(`确认删除「${selectedDetail.title}」吗？`)) return;
    try {
      await opportunitiesAPI.delete(selectedDetail.id);
      setSelectedId(null);
      await fetchList();
    } catch {
      setErrorText("机会删除失败，请确认登录状态。");
    }
  };

  return (
    <>
      <HeroBanner stats={stats} />

      <section className="grid grid-cols-[1.1fr_0.9fr] gap-4">
        <div className="rounded-lg border border-card-border bg-[#FEFCF9] p-5 shadow-[0_1px_3px_rgba(30,18,8,0.04)]">
          <div className="mb-4 flex items-start justify-between gap-4">
            <div>
              <div className="font-display text-[16px] font-semibold italic tracking-tight text-text-primary">
                机会列表
              </div>
              <p className="mt-1 text-[12px] leading-relaxed text-text-secondary">
                当前页已接入机会列表、批量评分、单条详情和单条评分；筛选后可直接查看每个机会的可达性结果。
              </p>
            </div>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => void fetchList()}
                className="inline-flex items-center gap-1 rounded border border-card-border bg-[#FAF7F2] px-3 py-2 text-[12px] text-text-secondary transition-colors hover:bg-white"
              >
                <RefreshCw className="h-3.5 w-3.5" />
                刷新
              </button>
              <button
                type="button"
                onClick={openCreateModal}
                disabled={!canManage}
                className={cn(
                  "inline-flex items-center gap-1 rounded px-3 py-2 text-[12px] transition-colors",
                  canManage
                    ? "bg-[#C4622D] text-white hover:bg-[#b15a2b]"
                    : "cursor-not-allowed bg-[#E7DFD3] text-text-muted"
                )}
              >
                <Plus className="h-3.5 w-3.5" />
                新增机会
              </button>
            </div>
          </div>

          <div className="mb-4 grid grid-cols-[1.2fr_0.65fr_0.65fr] gap-3">
            <label className="flex items-center gap-2 rounded border border-card-border bg-[#FAF7F2] px-3 py-2">
              <Search className="h-3.5 w-3.5 text-text-muted" />
              <input
                value={search}
                onChange={(event) => setSearch(event.target.value)}
                placeholder="搜索标题、主办方、标签"
                className="w-full bg-transparent text-[12px] text-text-primary outline-none placeholder:text-text-muted"
              />
            </label>

            <select
              value={typeFilter}
              onChange={(event) => setTypeFilter(event.target.value as "全部" | OpportunityType)}
              className="rounded border border-card-border bg-[#FAF7F2] px-3 py-2 text-[12px] text-text-primary outline-none"
            >
              <option value="全部">全部类型</option>
              {TYPE_OPTIONS.map((type) => (
                <option key={type} value={type}>
                  {type}
                </option>
              ))}
            </select>

            <select
              value={sortKey}
              onChange={(event) => setSortKey(event.target.value as "score" | "deadline" | "recent")}
              className="rounded border border-card-border bg-[#FAF7F2] px-3 py-2 text-[12px] text-text-primary outline-none"
            >
              <option value="score">按评分排序</option>
              <option value="deadline">按截止排序</option>
              <option value="recent">按创建时间</option>
            </select>
          </div>

          {!canManage && (
            <div className="mb-4 rounded border border-dashed border-card-border bg-[#FAF7F2] px-3 py-2 text-[11px] text-text-secondary">
              当前未检测到登录 token，机会的新增、编辑、删除按钮已保留但不可用；浏览和评分展示不受影响。
            </div>
          )}

          {errorText && (
            <div className="mb-4 rounded border border-[#E8C5B0] bg-[#FBF0E9] px-3 py-2 text-[11px] text-[#8A4B21]">
              {errorText}
            </div>
          )}

          <div className="grid grid-cols-[80px_1.4fr_0.8fr_0.8fr_0.65fr] border-b border-card-border px-3 pb-2 font-mono text-[9px] uppercase tracking-[0.16em] text-text-muted">
            <span>类型</span>
            <span>机会名称</span>
            <span>主办 / 标签</span>
            <span>评分</span>
            <span>截止</span>
          </div>

          <div className="mt-2 flex max-h-[620px] flex-col gap-2 overflow-y-auto pr-1">
            {loading && (
              <div className="rounded border border-card-border bg-[#FAF7F2] px-4 py-8 text-center text-[13px] text-text-secondary">
                正在拉取机会与评分数据…
              </div>
            )}

            {!loading && filteredItems.length === 0 && (
              <div className="rounded border border-card-border bg-[#FAF7F2] px-4 py-8 text-center font-display text-[15px] italic text-text-muted">
                当前筛选条件下暂无机会
              </div>
            )}

            {!loading &&
              filteredItems.map((item) => {
                const normalizedType = normalizeType(item.type);
                const typeStyle = typeStyles[normalizedType];
                const isActive = item.id === selectedId;
                const score = item.score_result?.score;
                return (
                  <button
                    key={item.id}
                    type="button"
                    onClick={() => setSelectedId(item.id)}
                    className={cn(
                      "grid grid-cols-[80px_1.4fr_0.8fr_0.8fr_0.65fr] items-start gap-3 rounded border px-3 py-3 text-left transition-all duration-150",
                      isActive
                        ? "border-[#D7B397] bg-[#FFF8F0] shadow-[0_3px_10px_rgba(196,98,45,0.08)]"
                        : "border-transparent bg-[#FAF7F2] hover:border-card-border hover:bg-white"
                    )}
                  >
                    <span
                      className={cn(
                        "inline-flex w-fit items-center rounded px-2 py-1 text-[10px] font-medium ring-1",
                        typeStyle.badge
                      )}
                    >
                      {normalizedType}
                    </span>

                    <div className="min-w-0">
                      <div className="truncate text-[13px] font-medium text-text-primary">{item.title}</div>
                      <div className="mt-1 truncate text-[11px] text-text-secondary">{item.description || item.organizer || "机会信息待补充"}</div>
                    </div>

                    <div className="text-[11px] text-text-secondary">
                      <div className="truncate">{item.organizer || "未标注主办方"}</div>
                      <div className="mt-1 truncate font-mono text-[10px] text-text-muted">
                        {arrayOrEmpty(item.tags).slice(0, 2).join(" / ") || "无标签"}
                      </div>
                    </div>

                    <div className="text-[11px]">
                      <div className={cn("font-display text-[22px] leading-none", scoreTone(score))}>
                        {typeof score === "number" ? Math.round(score) : "--"}
                      </div>
                      <div className="mt-1 text-text-secondary">{scoreSummary(item.score_result)}</div>
                    </div>

                    <div className="text-[11px] text-text-secondary">
                      {formatDateLabel(item.deadline || item.start_date)}
                    </div>
                  </button>
                );
              })}
          </div>
        </div>

        <aside className="rounded-lg border border-card-border bg-[#FEFCF9] p-5 shadow-[0_1px_3px_rgba(30,18,8,0.04)]">
          <div className="mb-4 flex items-start justify-between gap-3">
            <div>
              <div className="font-display text-[16px] font-semibold italic tracking-tight text-text-primary">
                机会详情
              </div>
              <div className="mt-1 font-mono text-[9px] uppercase tracking-[0.18em] text-text-muted">
                Opportunity Detail
              </div>
            </div>

            {selectedDetail && (
              <div className="flex gap-2">
                {selectedDetail.source_url && (
                  <a
                    href={selectedDetail.source_url}
                    target="_blank"
                    rel="noreferrer"
                    className="inline-flex items-center gap-1 rounded border border-card-border bg-[#FAF7F2] px-2.5 py-2 text-[11px] text-text-secondary transition-colors hover:bg-white"
                  >
                    <ExternalLink className="h-3.5 w-3.5" />
                    官网
                  </a>
                )}
                <button
                  type="button"
                  onClick={openEditModal}
                  disabled={!canManage || !selectedDetail}
                  className={cn(
                    "inline-flex items-center gap-1 rounded px-2.5 py-2 text-[11px] transition-colors",
                    canManage
                      ? "bg-[#FAF7F2] text-text-secondary hover:bg-white"
                      : "cursor-not-allowed bg-[#E7DFD3] text-text-muted"
                  )}
                >
                  <Pencil className="h-3.5 w-3.5" />
                  编辑
                </button>
              </div>
            )}
          </div>

          {detailLoading && (
            <div className="rounded border border-card-border bg-[#FAF7F2] px-4 py-10 text-center text-[13px] text-text-secondary">
              正在加载详情…
            </div>
          )}

          {!detailLoading && !selectedDetail && (
            <div className="rounded border border-dashed border-card-border bg-[#FAF7F2] px-4 py-10 text-center font-display text-[15px] italic text-text-muted">
              选择左侧机会以查看完整内容
            </div>
          )}

          {!detailLoading && selectedDetail && (
            <div className="space-y-5">
              <div>
                <div className="mb-2 flex items-center gap-2">
                  <span className={cn("h-2 w-2 rounded-full", selectedStyle.dot)} />
                  <span className={cn("rounded px-2 py-1 text-[10px] font-medium ring-1", selectedStyle.badge)}>
                    {selectedType}
                  </span>
                  {selectedScore && (
                    <span
                      className={cn(
                        "rounded px-2 py-1 text-[10px] font-medium ring-1",
                        selectedScore.passed
                          ? "bg-[#E8F2EC] text-[#2D6048] ring-[#A8D4BA]"
                          : "bg-[#F9E6E6] text-[#B84040] ring-[#E7B3B3]"
                      )}
                    >
                      {selectedScore.passed ? "门槛通过" : "门槛未过"}
                    </span>
                  )}
                </div>
                <h2 className="font-display text-[28px] leading-none tracking-tight text-text-primary">
                  {selectedDetail.title}
                </h2>
                <p className="mt-2 text-[13px] leading-relaxed text-text-secondary">
                  {selectedDetail.description || "当前机会暂未补充详细描述。"}
                </p>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="rounded border border-card-border bg-[#FAF7F2] px-4 py-3">
                  <div className="font-mono text-[9px] uppercase tracking-[0.18em] text-text-muted">评分</div>
                  <div className={cn("mt-2 font-display text-[34px] leading-none", scoreTone(selectedScore?.score))}>
                    {typeof selectedScore?.score === "number" ? selectedScore.score.toFixed(1) : "--"}
                  </div>
                  <div className="mt-2 text-[11px] text-text-secondary">{scoreSummary(selectedScore)}</div>
                </div>

                <div className="rounded border border-card-border bg-[#FAF7F2] px-4 py-3">
                  <div className="font-mono text-[9px] uppercase tracking-[0.18em] text-text-muted">时间窗口</div>
                  <div className="mt-2 text-[13px] font-medium text-text-primary">
                    {formatDateLabel(selectedDetail.start_date)}
                  </div>
                  <div className="mt-1 text-[11px] text-text-secondary">
                    截止：{formatDateLabel(selectedDetail.deadline)}
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                {[
                  ["主办方", selectedDetail.organizer || "未标注"],
                  ["地点", selectedDetail.location || "待确认"],
                  ["标签", arrayOrEmpty(selectedDetail.tags).join(" / ") || "无"],
                  ["面向专业", arrayOrEmpty(selectedDetail.target_majors).join(" / ") || "不限"],
                ].map(([label, value]) => (
                  <div key={label} className="rounded border border-card-border bg-[#FAF7F2] px-4 py-3">
                    <div className="font-mono text-[9px] uppercase tracking-[0.18em] text-text-muted">{label}</div>
                    <div className="mt-2 text-[12px] leading-relaxed text-text-primary">{value}</div>
                  </div>
                ))}
              </div>

              {selectedScore?.components && (
                <div className="rounded border border-card-border bg-[#FAF7F2] px-4 py-4">
                  <div className="font-display text-[15px] font-semibold italic text-text-primary">评分分量</div>
                  <div className="mt-3 space-y-2">
                    {[
                      ["时间可行性", selectedScore.components.time_feasibility],
                      ["收益", selectedScore.components.reward],
                      ["地域偏好", selectedScore.components.location_pref],
                      ["准备成本", selectedScore.components.prep_cost],
                    ].map(([label, value]) => (
                      <div key={label} className="grid grid-cols-[72px_1fr_42px] items-center gap-3 text-[11px] text-text-secondary">
                        <span>{label}</span>
                        <div className="h-1.5 overflow-hidden rounded-full bg-[#E7DFD3]">
                          <div className="h-full rounded-full bg-[#C4622D]" style={{ width: `${Math.max(5, Math.min(100, Number(value) * 100))}%` }} />
                        </div>
                        <span className="text-right font-mono">{Number(value).toFixed(2)}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {canManage && (
                <div className="flex justify-end">
                  <button
                    type="button"
                    onClick={() => void handleDelete()}
                    className="inline-flex items-center gap-1 rounded border border-[#E7B3B3] bg-[#F9E6E6] px-3 py-2 text-[12px] text-[#B84040] transition-colors hover:bg-[#F7DCDC]"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                    删除机会
                  </button>
                </div>
              )}
            </div>
          )}
        </aside>
      </section>

      {isModalOpen && (
        <div className="fixed inset-0 z-30 flex items-center justify-center bg-[#1F1B17]/45 px-6">
          <div className="w-full max-w-3xl rounded-lg border border-card-border bg-[#FEFCF9] p-6 shadow-[0_20px_60px_rgba(30,18,8,0.18)]">
            <div className="flex items-start justify-between">
              <div>
                <div className="font-display text-[24px] font-semibold leading-none tracking-tight text-text-primary">
                  {editingItem ? "编辑机会" : "新增机会"}
                </div>
                <div className="mt-1 text-[12px] text-text-secondary">
                  保持现有仪表盘风格，仅补充后端机会接口的真实操作入口。
                </div>
              </div>
              <button
                type="button"
                onClick={closeModal}
                className="rounded border border-card-border bg-[#FAF7F2] px-3 py-1.5 text-[12px] text-text-secondary"
              >
                关闭
              </button>
            </div>

            <div className="mt-5 grid grid-cols-2 gap-4">
              {[
                ["标题", "title", "text"],
                ["类型", "type", "select"],
                ["来源链接", "source_url", "text"],
                ["主办方", "organizer", "text"],
                ["开始时间", "start_date", "datetime-local"],
                ["截止时间", "deadline", "datetime-local"],
                ["活动时间", "event_date", "datetime-local"],
                ["地点", "location", "text"],
                ["面向专业", "target_majors", "text"],
                ["标签", "tags", "text"],
              ].map(([label, key, kind]) => (
                <label key={key} className={cn("flex flex-col gap-1.5", key === "source_url" && "col-span-2")}>
                  <span className="text-[12px] font-medium text-text-primary">{label}</span>
                  {kind === "select" ? (
                    <select
                      value={form.type}
                      onChange={(event) => setForm((current) => ({ ...current, type: event.target.value as OpportunityType }))}
                      className="rounded border border-card-border bg-[#FAF7F2] px-3 py-2 text-[12px] text-text-primary outline-none"
                    >
                      {TYPE_OPTIONS.map((type) => (
                        <option key={type} value={type}>
                          {type}
                        </option>
                      ))}
                    </select>
                  ) : (
                    <input
                      type={kind}
                      value={form[key as keyof OpportunityFormState]}
                      onChange={(event) => setForm((current) => ({ ...current, [key]: event.target.value }))}
                      className="rounded border border-card-border bg-[#FAF7F2] px-3 py-2 text-[12px] text-text-primary outline-none"
                    />
                  )}
                </label>
              ))}

              <label className="col-span-2 flex flex-col gap-1.5">
                <span className="text-[12px] font-medium text-text-primary">描述</span>
                <textarea
                  value={form.description}
                  onChange={(event) => setForm((current) => ({ ...current, description: event.target.value }))}
                  rows={5}
                  className="rounded border border-card-border bg-[#FAF7F2] px-3 py-2 text-[12px] text-text-primary outline-none"
                />
              </label>
            </div>

            <div className="mt-5 flex justify-end gap-2">
              <button
                type="button"
                onClick={closeModal}
                className="rounded border border-card-border bg-[#FAF7F2] px-4 py-2 text-[12px] text-text-secondary"
              >
                取消
              </button>
              <button
                type="button"
                onClick={() => void handleSubmit()}
                disabled={submitting}
                className={cn(
                  "rounded px-4 py-2 text-[12px] text-white",
                  submitting ? "bg-[#D6AA92]" : "bg-[#C4622D] hover:bg-[#b15a2b]"
                )}
              >
                {submitting ? "提交中…" : editingItem ? "保存修改" : "创建机会"}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
