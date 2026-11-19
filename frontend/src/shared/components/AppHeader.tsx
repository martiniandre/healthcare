import { Button } from "@/shared/components/ui/button"
import { useState, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { Link, useLocation } from "react-router-dom"
import { useAuthStore } from "../store/auth_store"
import { useLayoutStore } from "../store/layout_store"
import { useAuthRoleLabel } from "../hooks/useAuthRoleLabel"
import { ChevronRight, Menu, ShieldCheck, WifiOff } from "lucide-react"
import { useNotificationStream } from "../../modules/notifications/hooks/useNotificationStream"
import { NotificationBell } from "../../modules/notifications/components/NotificationBell"
import { findActiveNavigation, navigationItemTarget } from "./navigation"

interface BreadcrumbEntry {
  label: string
  path?: string
}

export const AppHeader = () => {
  const { t } = useTranslation("header")
  const { t: translateSidebar } = useTranslation("sidebar")
  const { role, email } = useAuthStore()
  const translateRole = useAuthRoleLabel()
  const toggleMobileSidebar = useLayoutStore((state) => state.toggleMobileSidebar)
  const location = useLocation()
  const [isOnline, setIsOnline] = useState(navigator.onLine)
  useNotificationStream()

  useEffect(() => {
    const handleOnline = () => setIsOnline(true)
    const handleOffline = () => setIsOnline(false)

    window.addEventListener("online", handleOnline)
    window.addEventListener("offline", handleOffline)

    return () => {
      window.removeEventListener("online", handleOnline)
      window.removeEventListener("offline", handleOffline)
    }
  }, [])

  const activeNavigation = findActiveNavigation(location.pathname, location.search)

  const breadcrumbEntries: BreadcrumbEntry[] = []
  if (activeNavigation) {
    breadcrumbEntries.push({ label: translateSidebar(activeNavigation.group.key) })
    breadcrumbEntries.push({ label: translateSidebar(activeNavigation.item.key), path: navigationItemTarget(activeNavigation.item) })
  }
  const patientPathSegments = location.pathname.match(/^\/patients\/([^/]+)(\/encounters\/[^/]+)?/)
  if (patientPathSegments) {
    breadcrumbEntries.push({ label: t("patientRecord"), path: `/patients/${patientPathSegments[1]}` })
    if (patientPathSegments[2]) {
      breadcrumbEntries.push({ label: t("encounter") })
    }
  }

  return (
    <header className="sticky top-0 z-40 flex h-14 w-full items-center gap-3 border-b border-border bg-card/95 px-4 backdrop-blur-sm md:px-6">
      <Button
        variant="ghost"
        size="icon"
        onClick={toggleMobileSidebar}
        aria-label={t("openMenu")}
        className="size-11 md:hidden"
      >
        <Menu className="size-5" />
      </Button>

      {breadcrumbEntries.length > 0 && (
        <nav aria-label={t("breadcrumbLabel")} className="hidden min-w-0 lg:block">
          <ol className="flex items-center gap-1.5 text-sm">
            {breadcrumbEntries.map((entry, entryIndex) => {
              const isLast = entryIndex === breadcrumbEntries.length - 1
              return (
                <li key={`${entry.label}-${entryIndex}`} className="flex min-w-0 items-center gap-1.5">
                  {entryIndex > 0 && <ChevronRight className="size-3.5 shrink-0 text-muted-foreground/70" aria-hidden="true" />}
                  {isLast || !entry.path ? (
                    <span
                      aria-current={isLast ? "page" : undefined}
                      className={isLast ? "truncate font-medium text-foreground" : "truncate text-muted-foreground"}
                    >
                      {entry.label}
                    </span>
                  ) : (
                    <Link to={entry.path} className="-mx-1 -my-2 truncate rounded-md px-1 py-2 text-muted-foreground transition-colors hover:text-foreground">
                      {entry.label}
                    </Link>
                  )}
                </li>
              )
            })}
          </ol>
        </nav>
      )}

      <div className="ml-auto flex items-center gap-2 md:gap-3">
        {!isOnline && (
          <div className="flex select-none items-center gap-1.5 rounded-md border border-destructive/30 bg-destructive/10 px-2.5 py-1 text-xs font-semibold text-destructive">
            <WifiOff className="size-3.5" />
            <span>{t("offlineStatus")}</span>
          </div>
        )}

        <NotificationBell />

        <div className="hidden h-6 w-px bg-border sm:block" />

        <div className="flex items-center gap-2.5">
          <div className="flex size-8 items-center justify-center rounded-full bg-sage-100 text-xs font-semibold text-sage-800 dark:bg-primary/15 dark:text-primary">
            {email ? email.charAt(0).toUpperCase() : "U"}
          </div>
          <div className="hidden flex-col items-start sm:flex">
            <span className="text-xs font-semibold leading-tight text-foreground">
              {email || t("defaultUserEmail")}
            </span>
            <div className="flex items-center gap-1">
              <ShieldCheck className="size-3 text-primary" />
              <span className="text-2xs text-muted-foreground">{translateRole(role)}</span>
            </div>
          </div>
        </div>
      </div>
    </header>
  )
}
