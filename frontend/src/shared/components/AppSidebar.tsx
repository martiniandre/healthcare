import { useNavigate, useLocation } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { useAuthStore } from "../store/auth_store"
import { useLayoutStore } from "../store/layout_store"
import type { ComponentType } from "react"
import { Activity, Users, BarChart3, LogOut, X, Sparkles, History, UserRound, LayoutDashboard, CalendarClock, Stethoscope, Sun, Moon } from "lucide-react"
import { useCapabilities } from "../auth/capabilities_context"
import { isCapabilityGranted, routeRequiredCapability, type RouteAccessKey } from "../auth/route_access"
import { LanguageSwitcher } from "./LanguageSwitcher"
import { Button } from "@/shared/components/ui/button"
import { cn } from "@/shared/utils/cn"

interface NavigationItem {
  key: string
  icon: ComponentType<{ className?: string }>
  path: string
  access: RouteAccessKey
}

interface NavigationGroup {
  key: string
  items: NavigationItem[]
}

const navigationGroups: NavigationGroup[] = [
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
      { key: "portal", icon: UserRound, path: "/portal", access: "portal" },
    ],
  },
]

export const AppSidebar = () => {
  const { t } = useTranslation("sidebar")
  const navigate = useNavigate()
  const location = useLocation()
  const { email, logout } = useAuthStore()
  const { capabilities } = useCapabilities()
  const { isMobileSidebarOpen, closeMobileSidebar, theme, toggleTheme } = useLayoutStore()

  const isItemVisible = (item: NavigationItem) => {
    return isCapabilityGranted(capabilities, routeRequiredCapability[item.access])
  }

  const isItemActive = (item: NavigationItem) => {
    const isPathActive =
      location.pathname === item.path || (item.path !== "/" && location.pathname.startsWith(item.path))
    const isPatientsHomeActive =
      item.path === "/" && (location.pathname === "/" || location.pathname.startsWith("/patients"))
    return isPathActive || isPatientsHomeActive
  }

  const visibleGroups = navigationGroups
    .map((group) => ({
      ...group,
      items: group.items.filter((item) => isItemVisible(item)),
    }))
    .filter((group) => group.items.length > 0)

  const handleNavigate = (path: string) => {
    navigate(path)
    closeMobileSidebar()
  }

  const isDarkTheme = theme === "dark"

  return (
    <>
      {isMobileSidebarOpen && (
        <div
          onClick={closeMobileSidebar}
          className="fixed inset-0 z-40 bg-black/40 backdrop-blur-[1px] md:hidden transition-opacity duration-300"
        />
      )}

      <aside
        className={`w-[240px] shrink-0 h-screen fixed md:sticky top-0 left-0 bg-card border-r border-border flex flex-col z-50 transition-transform duration-300 md:translate-x-0 ${
          isMobileSidebarOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <div className="px-5 py-5 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="bg-primary/8 p-2.5 rounded-xl border border-primary/10">
              <Activity className="w-5 h-5 text-primary" />
            </div>
            <div>
              <h1 className="text-sm font-display font-bold tracking-tight text-gray-900 leading-none">
                {t("title")}
              </h1>
              <span className="text-[10px] text-muted-foreground font-medium">{t("subtitle")}</span>
            </div>
          </div>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={closeMobileSidebar}
            aria-label={t("closeMenu")}
            className="md:hidden"
          >
            <X className="size-4" />
          </Button>
        </div>

        <div className="h-px bg-border mx-4" />

        <nav className="flex-1 overflow-y-auto px-3 py-5 flex flex-col">
          {visibleGroups.map((group, groupIndex) => (
            <section
              key={group.key}
              aria-label={t(group.key)}
              className={`flex flex-col gap-1 ${
                groupIndex > 0 ? "mt-3 border-t border-border/70 pt-4" : ""
              }`}
            >
              <div className="mb-1.5 flex items-center gap-2.5 px-3">
                <span className="h-px w-2.5 bg-primary/25" aria-hidden="true" />
                <h2 className="text-[10px] font-bold uppercase tracking-[0.16em] text-gray-400">
                  {t(group.key)}
                </h2>
              </div>
              {group.items.map((item) => {
                const isActive = isItemActive(item)
                return (
                  <Button
                    key={item.path}
                    variant="ghost"
                    onClick={() => handleNavigate(item.path)}
                    aria-current={isActive ? "page" : undefined}
                    className={cn(
                      "group relative h-auto w-full justify-start gap-3 px-3 py-2 text-left text-[13px] font-medium",
                      isActive
                        ? "bg-primary/8 font-semibold text-primary hover:bg-primary/10 hover:text-primary"
                        : "text-muted-foreground hover:text-foreground"
                    )}
                  >
                    <span
                      aria-hidden="true"
                      className={`absolute left-0 top-1/2 h-5 w-[3px] -translate-y-1/2 rounded-r-full transition-colors duration-200 ${
                        isActive ? "bg-primary" : "bg-transparent group-hover:bg-primary/40"
                      }`}
                    />
                    <span
                      className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-md border transition-colors duration-200 ${
                        isActive
                          ? "border-primary/15 bg-primary/10 text-primary"
                          : "border-transparent bg-gray-100/80 text-gray-500 group-hover:bg-gray-200/70 group-hover:text-gray-700"
                      }`}
                    >
                      <item.icon className="h-4 w-4" />
                    </span>
                    <span className="truncate">{t(item.key)}</span>
                  </Button>
                )
              })}
            </section>
          ))}
        </nav>

        <div className="px-3 pt-3 border-t border-border/70 flex flex-col gap-1">
          <LanguageSwitcher sidebarLayout />
          <Button
            variant="ghost"
            onClick={toggleTheme}
            aria-pressed={isDarkTheme}
            className="h-auto w-full justify-start gap-3 px-3 py-2 text-[13px] font-medium text-muted-foreground hover:text-foreground"
          >
            <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md border border-transparent bg-gray-100/80 text-gray-500 transition-colors duration-200">
              {isDarkTheme ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
            </span>
            <span className="truncate">{t("appearance")}</span>
          </Button>
        </div>

        <div className="px-3 mt-1 pb-3">
          <Button
            variant="ghost"
            onClick={() => {
              logout()
              closeMobileSidebar()
            }}
            className="h-auto w-full justify-start gap-3 px-3 py-2.5 text-[13px] font-semibold text-destructive hover:bg-destructive/10 hover:text-destructive"
          >
            <LogOut className="size-[18px] shrink-0" />
            {t("logout")}
          </Button>
        </div>

        <div className="px-5 py-3 border-t border-border flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-1.5 h-1.5 rounded-full bg-success animate-pulse-glow" />
            <span className="text-[10px] text-muted-foreground font-medium">
              FHIR R4 · gRPC-Web
            </span>
          </div>
          <span className="text-[9px] text-gray-300 font-mono">
            {email ? email.split("@")[0] : ""}
          </span>
        </div>
      </aside>
    </>
  )
}