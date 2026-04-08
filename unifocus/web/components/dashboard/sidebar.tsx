"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  type LucideIcon,
  Trophy,
  List,
  User,
  FileUp,
  Settings,
} from "lucide-react";

interface SidebarItemConfig {
  href: string;
  label: string;
  icon: LucideIcon;
  disabled?: boolean;
}

const navItems = [
  { href: "/dashboard", label: "机会总览", icon: LayoutDashboard },
  { href: "/dashboard/competitions", label: "竞赛管理", icon: Trophy },
  { href: "/dashboard/opportunities", label: "机会列表", icon: List },
  { href: "/dashboard/profile", label: "用户画像", icon: User, disabled: true },
  { href: "/dashboard/resume", label: "简历上传", icon: FileUp, disabled: true },
] satisfies SidebarItemConfig[];

const systemItems = [
  { href: "/dashboard/settings", label: "设置", icon: Settings, disabled: true },
] satisfies SidebarItemConfig[];

function isItemActive(pathname: string, href: string, disabled?: boolean) {
  if (disabled) return false;
  return href === "/dashboard" ? pathname === href : pathname.startsWith(href);
}

function SidebarNavItem({
  item,
  pathname,
}: {
  item: SidebarItemConfig;
  pathname: string;
}) {
  const isActive = isItemActive(pathname, item.href, item.disabled);

  if (item.disabled) {
    return (
      <div
        aria-disabled="true"
        className="flex cursor-not-allowed items-center gap-2.5 rounded px-3 py-2.5 text-[13px] font-medium text-white/25"
      >
        <item.icon className="h-[15px] w-[15px] shrink-0 opacity-35" />
        <span className="flex-1">{item.label}</span>
        <span className="rounded px-1.5 py-0.5 font-mono text-[8px] uppercase tracking-wider text-white/20 ring-1 ring-white/[0.1]">
          待开放
        </span>
      </div>
    );
  }

  return (
    <Link
      href={item.href}
      className={cn(
        "relative flex items-center gap-2.5 rounded px-3 py-2.5 text-[13px] font-medium transition-all duration-150",
        isActive
          ? "bg-white/[0.09] text-white"
          : "text-white/42 hover:bg-white/[0.05] hover:text-white/68"
      )}
    >
      {isActive && (
        <span className="absolute left-0 top-1/2 h-[18px] w-[2px] -translate-y-1/2 rounded-r-sm bg-[#C4622D]" />
      )}
      <item.icon
        className={cn("h-[15px] w-[15px] shrink-0", isActive ? "opacity-90" : "opacity-55")}
      />
      {item.label}
    </Link>
  );
}

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="fixed inset-y-0 left-0 z-10 flex w-[240px] flex-col overflow-y-auto bg-sidebar-bg px-4 pb-5 pt-7">
      {/* Brand */}
      <div className="mb-9 flex items-center gap-3 px-2">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded bg-[#C4622D] font-display text-sm font-semibold italic text-white">
          U
        </div>
        <div>
          <div className="font-display text-[15px] font-semibold leading-tight tracking-tight text-white">
            UniFocus
          </div>
          <div className="font-mono text-[8px] uppercase tracking-[0.2em] text-white/28">
            Opportunity Record
          </div>
        </div>
      </div>

      {/* Nav section */}
      <div className="mb-1.5 px-2 font-mono text-[8px] uppercase tracking-[0.2em] text-white/22">
        Navigation
      </div>
      <nav className="flex flex-col gap-0.5">
        {navItems.map((item) => (
          <SidebarNavItem key={item.href} item={item} pathname={pathname} />
        ))}
      </nav>

      {/* System section */}
      <div className="mb-1.5 mt-6 px-2 font-mono text-[8px] uppercase tracking-[0.2em] text-white/22">
        System
      </div>
      <nav className="flex flex-col gap-0.5">
        {systemItems.map((item) => (
          <SidebarNavItem key={item.href} item={item} pathname={pathname} />
        ))}
      </nav>

      {/* Footer */}
      <div className="mt-auto">
        <div className="mb-3 border-t border-white/[0.07]" />

        {/* Status */}
        <div className="mb-4 space-y-1.5 px-2">
          {[
            { online: true, label: "竞赛接口在线" },
            { online: false, label: "机会接口回落" },
          ].map(({ online, label }) => (
            <div key={label} className="flex items-center gap-2 text-[10px] text-white/28">
              <span
                className={cn(
                  "h-[5px] w-[5px] shrink-0 rounded-full",
                  online ? "bg-[#3A7555]" : "bg-white/18"
                )}
              />
              {label}
            </div>
          ))}
        </div>

        {/* Dashboard / Planet toggle */}
        <div className="flex gap-1 rounded bg-white/[0.07] p-[3px]">
          {[
            { href: "/dashboard", label: "Dashboard" },
            { href: "/dashboard/planet", label: "Planet" },
          ].map((item) => {
            const isActive =
              item.href === "/dashboard"
                ? pathname === "/dashboard"
                : pathname.startsWith(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "flex-1 rounded-[4px] py-[7px] text-center font-mono text-[9px] uppercase tracking-wider transition-all duration-150",
                  isActive
                    ? "bg-[#C4622D] text-white shadow-sm"
                    : "text-white/32 hover:text-white/52"
                )}
              >
                {item.label}
              </Link>
            );
          })}
        </div>
      </div>
    </aside>
  );
}
