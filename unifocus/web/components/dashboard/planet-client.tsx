"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import PlanetScene, { type ActiveNode } from "./planet-scene";
import { PlanetInfoPanel } from "./planet-info-panel";
import { SkillRadar, buildRadarData } from "./skill-radar";
import { OverloadWarningPanel } from "./overload-warning";
import { competitionsAPI, type Competition } from "@/lib/api/competitions";
import { opportunitiesAPI } from "@/lib/api/opportunities";
import { profileAPI, type UserProfile } from "@/lib/api/profile";
import { arrayOrEmpty } from "@/lib/utils";
import {
  buildOrbitGroups,
  computePlanetSizes,
  computeOverload,
  type RawOpportunity,
} from "./planet-data";

const DEFAULT_NODE: ActiveNode = {
  title: "UniFocus 主星",
  context: "悬浮星球卫星节点以查看机会详情。",
  groupKey: "main",
};

export default function PlanetClient() {
  const [competitions, setCompetitions] = useState<Competition[]>([]);
  const [opportunities, setOpportunities] = useState<RawOpportunity[]>([]);
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [activeNode, setActiveNode] = useState<ActiveNode>(DEFAULT_NODE);

  // Fetch data
  useEffect(() => {
    let mounted = true;

    const fetchProfileAndOpportunities = async () => {
      let resolvedProfile: UserProfile | null = null;

      try {
        resolvedProfile = await profileAPI.getProfile();
        if (mounted) {
          setProfile(resolvedProfile);
        }
      } catch {
        if (mounted) {
          setProfile(null);
        }
      }

      try {
        const scored = await opportunitiesAPI.listScored(
          { limit: 50 },
          resolvedProfile?.user_id
        );
        if (mounted) {
          setOpportunities(scored.data);
        }
      } catch {
        try {
          const fallback = await opportunitiesAPI.list({ limit: 50 });
          if (mounted) {
            setOpportunities(fallback.data);
          }
        } catch {}
      }
    };

    void fetchProfileAndOpportunities();

    competitionsAPI
      .list({ limit: 500 })
      .then((data) => { if (mounted) setCompetitions(data); })
      .catch(() => {});

    return () => { mounted = false; };
  }, []);

  const orbitGroups = useMemo(
    () => buildOrbitGroups({ competitions, opportunities }),
    [competitions, opportunities]
  );

  const planetSizes = useMemo(() => computePlanetSizes(orbitGroups), [orbitGroups]);

  const overloadWarnings = useMemo(() => computeOverload(opportunities), [opportunities]);

  const allTags = useMemo(
    () => opportunities.flatMap((o) => o.tags ?? []),
    [opportunities]
  );

  const radarData = useMemo(
    () => profile
      ? buildRadarData(arrayOrEmpty(profile.skills), arrayOrEmpty(profile.certificates), allTags)
      : null,
    [profile, allTags]
  );

  const handleActiveNode = useCallback((node: ActiveNode) => setActiveNode(node), []);

  return (
    <div className="relative h-full w-full overflow-hidden bg-[#040b18]">
      <PlanetScene
        orbitGroups={orbitGroups}
        planetSizes={planetSizes}
        onActiveNodeChange={handleActiveNode}
      />

      {/* Floating overlays */}
      <OverloadWarningPanel warnings={overloadWarnings} />
      <PlanetInfoPanel node={activeNode} />
      {radarData && <SkillRadar data={radarData} />}

      {/* Brand watermark */}
      <div className="pointer-events-none absolute right-5 top-5 text-right">
        <div className="font-display text-lg font-bold tracking-tight text-white/20">
          UniFocus
        </div>
        <div className="font-mono text-[9px] uppercase tracking-widest text-white/10">
          opportunity planet
        </div>
      </div>
    </div>
  );
}
