import type { ComponentType } from "react"
import { Activity, Users, BarChart3, Sparkles, History, LayoutDashboard, CalendarClock, Stethoscope, Heart, Pill, FileText, Image } from "lucide-react"
import type { RouteAccessKey } from "../auth/route_access"

export interface NavigationItem {
  key: string
  icon: ComponentType<{ className?: string }>
  path: string
  tab?: string
  access: RouteAccessKey
}

export interface NavigationGroup {
  key: string
  items: NavigationItem[]
}

export const navigationGroups: NavigationGroup[] = [
  {
    key: "topics.clinical",
    items: [
      { key: "patients", icon: Users, path: "/", access: "patients" },
      { key: "dashboard", icon: LayoutDashboard, path: "/dashboard", access: "dashboard" },
      { key: "schedule", icon: CalendarClock, path: "/schedule", access: "schedule" },
    ],
  },
  {
    key: "topics.diagnostics",
    items: [
      { key: "telemetry", icon: Activity, path: "/telemetry", access: "telemetry" },
      { key: "examAnalyzer", icon: Sparkles, path: "/exam-analyzer", access: "examAnalyzer" },
    ],
  },
  {
    key: "topics.operations",
    items: [
      { key: "analytics", icon: BarChart3, path: "/analytics", access: "analytics" },
      { key: "staffManagement", icon: Stethoscope, path: "/staff", access: "staff" },
      { key: "auditLogs", icon: History, path: "/audit-logs", access: "auditLogs" },
    ],
  },
  {
    key: "topics.patient",
    items: [
      { key: "portal:tabs.dashboard", icon: LayoutDashboard, path: "/portal", tab: "dashboard", access: "portal" },
      { key: "portal:tabs.appointments", icon: CalendarClock, path: "/portal", tab: "appointments", access: "portal" },
      { key: "portal:tabs.encounters", icon: History, path: "/portal", tab: "encounters", access: "portal" },
      { key: "portal:tabs.observations", icon: Heart, path: "/portal", tab: "observations", access: "portal" },
      { key: "portal:tabs.conditions", icon: Activity, path: "/portal", tab: "conditions", access: "portal" },
      { key: "portal:tabs.medications", icon: Pill, path: "/portal", tab: "medications", access: "portal" },
      { key: "portal:tabs.reports", icon: FileText, path: "/portal", tab: "reports", access: "portal" },
      { key: "portal:tabs.imaging", icon: Image, path: "/portal", tab: "imaging", access: "portal" },
    ],
  },
]

const DEFAULT_PORTAL_TAB = "dashboard"

export const navigationItemTarget = (item: NavigationItem): string =>
  item.tab ? `${item.path}?tab=${item.tab}` : item.path

export const navigationItemId = (item: NavigationItem): string => navigationItemTarget(item)

export const isNavigationItemActive = (item: NavigationItem, pathname: string, search = ""): boolean => {
  if (item.tab) {
    const activeTab = new URLSearchParams(search).get("tab") ?? DEFAULT_PORTAL_TAB
    return pathname.startsWith(item.path) && activeTab === item.tab
  }
  const isPathActive = pathname === item.path || (item.path !== "/" && pathname.startsWith(item.path))
  const isPatientsHomeActive = item.path === "/" && (pathname === "/" || pathname.startsWith("/patients"))
  return isPathActive || isPatientsHomeActive
}

export interface ActiveNavigationMatch {
  group: NavigationGroup
  item: NavigationItem
}

export const findActiveNavigation = (pathname: string, search = ""): ActiveNavigationMatch | null => {
  for (const group of navigationGroups) {
    const matchedItem = group.items.find((item) => isNavigationItemActive(item, pathname, search))
    if (matchedItem) {
      return { group, item: matchedItem }
    }
  }
  return null
}
