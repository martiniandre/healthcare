import {
  Users,
  CalendarClock,
  Activity,
  LayoutDashboard,
  ScanSearch,
  BarChart3,
  IdCard,
  ShieldAlert,
  Settings,
  UserRound,
} from "lucide-react"
import type { LucideIcon } from "lucide-react"

export interface SidebarNavigationItem {
  key: string
  icon: LucideIcon
  path: string
  staffOnly?: boolean
  patientOnly?: boolean
  adminOnly?: boolean
  disabled?: boolean
}

export interface SidebarNavigationGroup {
  key: string
  items: SidebarNavigationItem[]
}

export const navigationGroups: SidebarNavigationGroup[] = [
  {
    key: "groups.care",
    items: [
      { key: "patients", icon: Users, path: "/", staffOnly: true },
      { key: "schedule", icon: CalendarClock, path: "/schedule", staffOnly: true },
      { key: "telemetry", icon: Activity, path: "/telemetry", staffOnly: true },
    ],
  },
  {
    key: "groups.intelligence",
    items: [
      { key: "dashboard", icon: LayoutDashboard, path: "/dashboard", staffOnly: true },
      { key: "examAnalyzer", icon: ScanSearch, path: "/exam-analyzer", staffOnly: true },
      { key: "analytics", icon: BarChart3, path: "/analytics", staffOnly: true },
    ],
  },
  {
    key: "groups.administration",
    items: [
      { key: "staffManagement", icon: IdCard, path: "/staff", staffOnly: true },
      { key: "auditLogs", icon: ShieldAlert, path: "/audit-logs", staffOnly: true, adminOnly: true },
      { key: "settings", icon: Settings, path: "/settings", staffOnly: true, disabled: true },
    ],
  },
  {
    key: "groups.patientAccess",
    items: [{ key: "portal", icon: UserRound, path: "/portal", patientOnly: true }],
  },
]