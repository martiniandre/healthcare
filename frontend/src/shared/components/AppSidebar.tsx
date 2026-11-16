import { useNavigate, useLocation } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { useAuthStore } from "../store/auth_store"
import { useLayoutStore } from "../store/layout_store"
import type { ComponentType, ReactNode } from "react"
import { Activity, Users, BarChart3, LogOut, X, Sparkles, History, UserRound, LayoutDashboard, CalendarClock, Stethoscope, Sun, Moon, Leaf } from "lucide-react"
import { useCapabilities } from "../auth/capabilities_context"
import { isCapabilityGranted, routeRequiredCapability, type RouteAccessKey } from "../auth/route_access"
import { LanguageSwitcher } from "./LanguageSwitcher"
import { Button } from "@/shared/components/ui/button"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/shared/components/ui/tooltip"
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

  const railItemClassName =
    "group relative h-auto w-full justify-start gap-3 rounded-md px-3 py-2.5 text-left text-sm font-medium md:size-11 md:justify-center md:p-0"

  const renderRailTooltip = (label: string, trigger: ReactNode) => (
    <Tooltip>
      <TooltipTrigger asChild>{trigger}</TooltipTrigger>
      <TooltipContent side="right" sideOffset={10} className="hidden font-medium md:block">
        {label}
      </TooltipContent>
    </Tooltip>
  )

  return (
    <>
      {isMobileSidebarOpen && (
        <div
          onClick={closeMobileSidebar}
          className="fixed inset-0 z-40 bg-foreground/30 backdrop-blur-[2px] md:hidden transition-opacity duration-300"
        />
      )}

      <aside
        className={cn(
          "fixed left-0 top-0 z-50 flex h-screen w-[260px] shrink-0 flex-col border-r border-border bg-card shadow-lg transition-transform duration-300",
          "md:sticky md:w-[68px] md:translate-x-0 md:shadow-none",
          isMobileSidebarOpen ? "translate-x-0" : "-translate-x-full"
        )}
      >
        <div className="flex h-full flex-col md:items-center">
          <div className="flex items-center justify-between px-5 py-5 md:justify-center md:px-0">
            <div className="flex items-center gap-3">
              <div className="flex size-10 items-center justify-center rounded-md bg-primary text-primary-foreground">
                <Leaf className="size-5" aria-hidden="true" />
              </div>
              <div className="md:sr-only">
                <h1 className="font-display text-lg font-semibold leading-none text-foreground">{t("title")}</h1>
                <span className="text-xs text-muted-foreground">{t("subtitle")}</span>
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

          <nav className="flex flex-1 flex-col gap-1 overflow-y-auto px-3 pb-4 md:items-center md:px-0">
            {visibleGroups.map((group, groupIndex) => (
              <section
                key={group.key}
                aria-label={t(group.key)}
                className={cn("flex w-full flex-col gap-1 md:items-center", groupIndex > 0 && "mt-3 pt-3")}
              >
                {groupIndex > 0 && <span aria-hidden="true" className="mb-2 hidden h-px w-6 bg-border md:block" />}
                <h2 className="px-3 pb-1 text-xs font-semibold text-muted-foreground md:sr-only">{t(group.key)}</h2>
                {group.items.map((item) => {
                  const isActive = isItemActive(item)
                  const label = t(item.key)
                  return (
                    <div key={item.path} className="w-full md:w-auto">
                      {renderRailTooltip(
                        label,
                        <Button
                          variant="ghost"
                          onClick={() => handleNavigate(item.path)}
                          aria-current={isActive ? "page" : undefined}
                          className={cn(
                            railItemClassName,
                            isActive
                              ? "bg-sage-100 text-sage-800 hover:bg-sage-100 hover:text-sage-800 dark:bg-primary/15 dark:text-primary dark:hover:bg-primary/15"
                              : "text-muted-foreground hover:bg-accent hover:text-foreground"
                          )}
                        >
                          <item.icon className="size-5" />
                          <span className="truncate md:sr-only">{label}</span>
                        </Button>
                      )}
                    </div>
                  )
                })}
              </section>
            ))}
          </nav>

          <div className="flex flex-col gap-1 border-t border-border/70 px-3 py-3 md:items-center md:border-t-0 md:px-0">
            <div className="w-full md:w-auto">
              <LanguageSwitcher sidebarLayout />
            </div>
            <div className="w-full md:w-auto">
              {renderRailTooltip(
                t("appearance"),
                <Button
                  variant="ghost"
                  onClick={toggleTheme}
                  aria-pressed={isDarkTheme}
                  className={cn(railItemClassName, "text-muted-foreground hover:bg-accent hover:text-foreground")}
                >
                  {isDarkTheme ? <Sun className="size-5" /> : <Moon className="size-5" />}
                  <span className="truncate md:sr-only">{t("appearance")}</span>
                </Button>
              )}
            </div>
            <div className="w-full md:w-auto">
              {renderRailTooltip(
                t("logout"),
                <Button
                  variant="ghost"
                  onClick={() => {
                    logout()
                    closeMobileSidebar()
                  }}
                  className={cn(railItemClassName, "text-destructive hover:bg-destructive/10 hover:text-destructive")}
                >
                  <LogOut className="size-5" />
                  <span className="truncate md:sr-only">{t("logout")}</span>
                </Button>
              )}
            </div>
            <span className="sr-only">{email ?? ""}</span>
          </div>
        </div>
      </aside>
    </>
  )
}
