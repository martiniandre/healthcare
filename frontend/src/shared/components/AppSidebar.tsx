import { useNavigate, useLocation } from "react-router-dom"
import { useTranslation } from "react-i18next"
import type { ReactNode } from "react"
import { LogOut, X, Sun, Moon, Leaf, PanelLeftClose, PanelLeftOpen } from "lucide-react"
import { useAuthStore } from "../store/auth_store"
import { TABLET_MEDIA_QUERY, useMediaQuery } from "../hooks/useMediaQuery"
import { useLayoutStore } from "../store/layout_store"
import { useCapabilities } from "../auth/capabilities_context"
import { isCapabilityGranted, routeRequiredCapability } from "../auth/route_access"
import { LanguageSwitcher } from "./LanguageSwitcher"
import { navigationGroups, isNavigationItemActive, navigationItemId, navigationItemTarget, type NavigationItem } from "./navigation"
import { Button } from "@/shared/components/ui/button"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/shared/components/ui/tooltip"
import { cn } from "@/shared/utils/cn"

export const AppSidebar = () => {
  const { t } = useTranslation("sidebar")
  const navigate = useNavigate()
  const location = useLocation()
  const { logout } = useAuthStore()
  const { capabilities } = useCapabilities()
  const {
    isMobileSidebarOpen,
    closeMobileSidebar,
    theme,
    toggleTheme,
    isSidebarCollapsed: isSidebarCollapsedByUser,
    toggleSidebarCollapsed,
  } = useLayoutStore()
  const isTabletViewport = useMediaQuery(TABLET_MEDIA_QUERY)
  const isSidebarCollapsed = isSidebarCollapsedByUser || isTabletViewport

  const isItemVisible = (item: NavigationItem) =>
    isCapabilityGranted(capabilities, routeRequiredCapability[item.access])

  const visibleGroups = navigationGroups
    .map((group) => ({ ...group, items: group.items.filter((item) => isItemVisible(item)) }))
    .filter((group) => group.items.length > 0)

  const handleNavigate = (path: string) => {
    navigate(path)
    closeMobileSidebar()
  }

  const isDarkTheme = theme === "dark"
  const collapsedOnly = (className: string) => (isSidebarCollapsed ? className : "")

  const itemClassName = cn(
    "group relative h-11 w-full justify-start gap-3 rounded-md px-3 text-left text-sm font-medium md:h-9",
    collapsedOnly("md:size-10 md:justify-center md:px-0")
  )
  const labelClassName = cn("truncate", collapsedOnly("md:sr-only"))

  const withCollapsedTooltip = (label: string, trigger: ReactNode) =>
    isSidebarCollapsed ? (
      <Tooltip>
        <TooltipTrigger asChild>{trigger}</TooltipTrigger>
        <TooltipContent side="right" sideOffset={8} className="hidden font-medium md:block">
          {label}
        </TooltipContent>
      </Tooltip>
    ) : (
      trigger
    )

  return (
    <>
      {isMobileSidebarOpen && (
        <div
          onClick={closeMobileSidebar}
          className="fixed inset-0 z-40 bg-foreground/30 backdrop-blur-[1px] md:hidden transition-opacity duration-300"
        />
      )}

      <aside
        data-collapsed={isSidebarCollapsed}
        className={cn(
          "fixed left-0 top-0 z-50 flex h-screen w-[248px] shrink-0 flex-col border-r border-border bg-card shadow-lg transition-[transform,width] duration-200",
          "md:sticky md:translate-x-0 md:shadow-none",
          isSidebarCollapsed && "md:w-[68px]",
          isMobileSidebarOpen ? "translate-x-0" : "-translate-x-full"
        )}
      >
        <div className={cn("flex h-14 items-center justify-between gap-2 border-b border-border px-4", collapsedOnly("md:justify-center md:px-0"))}>
          <div className="flex min-w-0 items-center gap-2.5">
            <div className="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary text-primary-foreground">
              <Leaf className="size-4" aria-hidden="true" />
            </div>
            <div className={cn("min-w-0", collapsedOnly("md:sr-only"))}>
              <p className="truncate text-base font-semibold leading-tight text-foreground">{t("title")}</p>
              <span className="block truncate text-2xs text-muted-foreground">{t("subtitle")}</span>
            </div>
          </div>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={toggleSidebarCollapsed}
            aria-label={isSidebarCollapsed ? t("expandMenu") : t("collapseMenu")}
            aria-expanded={!isSidebarCollapsed}
            className={cn("hidden text-muted-foreground lg:inline-flex", collapsedOnly("md:hidden"))}
          >
            <PanelLeftClose />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={closeMobileSidebar}
            aria-label={t("closeMenu")}
            className="size-11 md:hidden"
          >
            <X className="size-4" />
          </Button>
        </div>

        <nav className={cn("flex flex-1 flex-col gap-5 overflow-y-auto px-3 py-4", collapsedOnly("md:items-center md:px-0"))}>
          {visibleGroups.map((group) => (
            <section
              key={group.key}
              aria-label={t(group.key)}
              className={cn("flex w-full flex-col gap-0.5", collapsedOnly("md:w-auto md:items-center"))}
            >
              <p
                aria-hidden="true"
                className={cn(
                  "mb-1 px-3 text-xs font-medium text-muted-foreground",
                  collapsedOnly("md:sr-only")
                )}
              >
                {t(group.key)}
              </p>
              {group.items.map((item) => {
                const isActive = isNavigationItemActive(item, location.pathname, location.search)
                const label = t(item.key)
                return (
                  <div key={navigationItemId(item)} className={cn("w-full", collapsedOnly("md:w-auto"))}>
                    {withCollapsedTooltip(
                      label,
                      <Button
                        variant="ghost"
                        onClick={() => handleNavigate(navigationItemTarget(item))}
                        aria-current={isActive ? "page" : undefined}
                        className={cn(
                          itemClassName,
                          isActive
                            ? "bg-sage-100 font-semibold text-sage-800 hover:bg-sage-100 hover:text-sage-800 dark:bg-primary/15 dark:text-primary dark:hover:bg-primary/15"
                            : "text-muted-foreground hover:bg-accent hover:text-foreground"
                        )}
                      >
                        <item.icon className={cn("size-[18px]", isActive ? "text-primary" : "text-muted-foreground group-hover:text-foreground")} />
                        <span className={labelClassName}>{label}</span>
                      </Button>
                    )}
                  </div>
                )
              })}
            </section>
          ))}
        </nav>

        <div className={cn("flex flex-col gap-0.5 border-t border-border px-3 py-3", collapsedOnly("md:items-center md:px-0"))}>
          <p className={cn("mb-1 px-3 text-xs font-medium text-muted-foreground", collapsedOnly("md:sr-only"))}>
            {t("preferences")}
          </p>
          <div className={cn("w-full", collapsedOnly("md:w-auto"))}>
            <LanguageSwitcher sidebarLayout compact={isSidebarCollapsed} />
          </div>
          <div className={cn("w-full", collapsedOnly("md:w-auto"))}>
            {withCollapsedTooltip(
              t("appearance"),
              <Button
                variant="ghost"
                onClick={toggleTheme}
                aria-pressed={isDarkTheme}
                className={cn(itemClassName, "text-muted-foreground hover:bg-accent hover:text-foreground")}
              >
                {isDarkTheme ? <Sun className="size-[18px]" /> : <Moon className="size-[18px]" />}
                <span className={labelClassName}>{t("appearance")}</span>
              </Button>
            )}
          </div>
          <div className={cn("w-full", collapsedOnly("md:w-auto"))}>
            {withCollapsedTooltip(
              t("logout"),
              <Button
                variant="ghost"
                onClick={() => {
                  logout()
                  closeMobileSidebar()
                }}
                className={cn(itemClassName, "text-destructive hover:bg-destructive/10 hover:text-destructive")}
              >
                <LogOut className="size-[18px]" />
                <span className={labelClassName}>{t("logout")}</span>
              </Button>
            )}
          </div>
          {isSidebarCollapsed && !isTabletViewport && (
            <div className="hidden md:block">
              {withCollapsedTooltip(
                t("expandMenu"),
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={toggleSidebarCollapsed}
                  aria-label={t("expandMenu")}
                  aria-expanded={false}
                  className="mt-1 size-10 text-muted-foreground"
                >
                  <PanelLeftOpen className="size-[18px]" />
                </Button>
              )}
            </div>
          )}
        </div>
      </aside>
    </>
  )
}
