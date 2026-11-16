import { Button } from "@/shared/components/ui/button"
import { useState, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { useAuthStore } from "../store/auth_store"
import { useLayoutStore } from "../store/layout_store"
import { useAuthRoleLabel } from "../hooks/useAuthRoleLabel"
import { Menu, ShieldCheck, WifiOff } from "lucide-react"
import { useNotificationStream } from "../../modules/notifications/hooks/useNotificationStream"
import { NotificationBell } from "../../modules/notifications/components/NotificationBell"

export const AppHeader = () => {
  const { t } = useTranslation("header")
  const { role, email } = useAuthStore()
  const translateRole = useAuthRoleLabel()
  const toggleMobileSidebar = useLayoutStore((state) => state.toggleMobileSidebar)
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

  return (
    <header className="sticky top-0 z-40 flex w-full items-center justify-end gap-3 border-b border-border bg-background/90 px-4 py-2.5 backdrop-blur-sm md:px-8">
      <Button
        variant="ghost"
        size="icon"
        onClick={toggleMobileSidebar}
        aria-label={t("openMenu")}
        className="mr-auto md:hidden"
      >
        <Menu className="size-5" />
      </Button>

      {!isOnline && (
        <div className="mr-2 flex select-none items-center gap-1.5 rounded-md border border-destructive/30 bg-destructive/10 px-3 py-1 text-xs font-semibold text-destructive animate-pulse">
          <WifiOff className="size-3.5" />
          <span>{t("offlineStatus")}</span>
        </div>
      )}

      <NotificationBell />

      <div className="flex items-center gap-2.5">
        <div className="flex size-8 items-center justify-center rounded-full bg-sage-100 text-xs font-bold text-primary dark:bg-primary/15">
          {email ? email.charAt(0).toUpperCase() : "U"}
        </div>
        <div className="hidden sm:flex flex-col items-start">
          <span className="text-xs font-semibold leading-tight text-foreground">
            {email || t("defaultUserEmail")}
          </span>
          <div className="flex items-center gap-1">
            <ShieldCheck className="size-3 text-primary" />
            <span className="text-[10px] text-muted-foreground font-medium">
              {translateRole(role)}
            </span>
          </div>
        </div>
      </div>
    </header>
  )
}
