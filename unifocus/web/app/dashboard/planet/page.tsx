"use client";

// PlanetClient is dynamically imported to avoid SSR issues with three.js
import dynamic from "next/dynamic";

const PlanetClient = dynamic(
  () => import("@/components/dashboard/planet-client"),
  { ssr: false, loading: () => <PlanetLoading /> }
);

function PlanetLoading() {
  return (
    <div className="flex h-full w-full items-center justify-center bg-[#040b18]">
      <div className="relative flex flex-col items-center gap-6">
        {/* Orbital ring spinner */}
        <div className="relative h-20 w-20">
          {/* Outer ring */}
          <div
            className="absolute inset-0 rounded-full"
            style={{
              border: "1.5px solid rgba(108,92,231,0.5)",
              animation: "orbit-cw 2.8s linear infinite",
            }}
          >
            {/* Dot on outer ring */}
            <div
              className="absolute -top-[3px] left-1/2 h-[6px] w-[6px] -translate-x-1/2 rounded-full"
              style={{ background: "#6c5ce7", boxShadow: "0 0 8px 2px rgba(108,92,231,0.7)" }}
            />
          </div>
          {/* Inner ring */}
          <div
            className="absolute inset-[14px] rounded-full"
            style={{
              border: "1.5px solid rgba(82,212,255,0.35)",
              animation: "orbit-ccw 1.9s linear infinite",
            }}
          >
            <div
              className="absolute -top-[3px] left-1/2 h-[5px] w-[5px] -translate-x-1/2 rounded-full"
              style={{ background: "#52d4ff", boxShadow: "0 0 6px 2px rgba(82,212,255,0.6)" }}
            />
          </div>
          {/* Core glyph */}
          <div className="absolute inset-0 flex items-center justify-center">
            <span
              className="font-display text-sm font-bold text-white/80"
              style={{ textShadow: "0 0 12px rgba(108,92,231,0.8)" }}
            >
              U
            </span>
          </div>
        </div>

        {/* Label */}
        <div className="text-center">
          <div className="font-mono text-[9px] uppercase tracking-[0.25em] text-white/25">
            Initializing Galaxy
          </div>
        </div>
      </div>
    </div>
  );
}

export default function PlanetPage() {
  return (
    // Escape dashboard layout padding to fill the full main area
    <div className="-mx-7 -my-6 h-[100vh]">
      <PlanetClient />
    </div>
  );
}
