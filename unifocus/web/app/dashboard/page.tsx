"use client";

import { useEffect, useState } from "react";
import { HeroBanner } from "@/components/dashboard/hero-banner";
import { OpportunityFeed, type OpportunityItem } from "@/components/dashboard/opportunity-feed";
import { TrackSummary, type TrackCardData } from "@/components/dashboard/track-summary";
import { competitionsAPI, type Competition } from "@/lib/api/competitions";

const TRACK_DATA: TrackCardData[] = [
  {
    color: "gold",
    label: "奖学金轨道",
    items: [
      { name: "国家奖学金", hint: "高势能" },
      { name: "一等奖学金", hint: "绩点联动" },
      { name: "励志奖学金", hint: "成长加权" },
    ],
  },
  {
    color: "violet",
    label: "实习轨道",
    items: [
      { name: "互联网大厂", hint: "高竞争" },
      { name: "产品经理", hint: "策略双栖" },
      { name: "研究助理", hint: "科研延展" },
    ],
  },
  {
    color: "cyan",
    label: "竞赛热度",
    items: [
      { name: "国家级A*类", hint: "12 项" },
      { name: "国家级A类", hint: "18 项" },
      { name: "国家级B类", hint: "24 项" },
    ],
  },
];

const STATIC_ITEMS: OpportunityItem[] = [
  {
    id: -1,
    name: "国家奖学金",
    meta: "高势能荣誉识别",
    type: "scholarship",
    badge: "奖学金",
    time: "9–10月",
  },
  {
    id: -2,
    name: "互联网大厂实习",
    meta: "高竞争密度岗位",
    type: "internship",
    badge: "实习",
    time: "春招",
  },
];

export default function DashboardPage() {
  const [competitions, setCompetitions] = useState<Competition[]>([]);

  useEffect(() => {
    let active = true;
    competitionsAPI
      .list({ limit: 500 })
      .then((data) => {
        if (active) setCompetitions(data);
      })
      .catch(() => {});
    return () => { active = false; };
  }, []);

  const nationalCount = competitions.filter((c) => c.level?.startsWith("国家级")).length;

  const stats = [
    { label: "竞赛总量", value: competitions.length, detail: `${nationalCount} 国家级` },
    { label: "奖学金", value: 5, detail: "荣誉轨道", colorClass: "text-accent-amber" },
    { label: "实习实践", value: 5, detail: "跃迁轨道", colorClass: "text-accent-violet" },
    { label: "机会节点", value: competitions.length + 10, detail: "星图映射", colorClass: "text-accent-cyan" },
  ];

  const competitionItems: OpportunityItem[] = competitions.slice(0, 8).map((c) => ({
    id: c.id,
    name: c.name,
    meta: c.timeline_hint || c.note || "",
    type: "competition",
    badge: c.level || "竞赛",
    time: c.typical_time_window || "",
  }));

  const feedItems: OpportunityItem[] = [...competitionItems, ...STATIC_ITEMS];

  return (
    <>
      <HeroBanner stats={stats} />
      <div className="grid grid-cols-[1.4fr_0.6fr] gap-4">
        <OpportunityFeed items={feedItems} />
        <TrackSummary tracks={TRACK_DATA} />
      </div>
    </>
  );
}
