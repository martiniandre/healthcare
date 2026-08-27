import { useState, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { useAuthStore } from "../store/auth_store"
import { useLayoutStore } from "../store/layout_store"
import { Menu, ShieldCheck, WifiOff } from "lucide-react"
import { LanguageSwitcher } from "./LanguageSwitcher"
import { NotificationBell } from "../../modules/notifications/components/NotificationBell"
import { Button } from "./ui/Button"
import { Badge } from "./ui/Badge"

export const AppHeader = () => {
  const { t } = useTranslation("header")
  const { role, email } = useAuthStore()
  const toggleMobileSidebar = useLayoutStore((state) => state.toggleMobileSidebar)
  const [isOnline, setIsOnline] = useState(navigator.onLine)

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

  const translateRole = (userRole: string | null) => {
    if (!userRole) {
      return t("roles.RoleDefault")
    }
    return t(`roles.${userRole}`, { defaultValue: t("roles.RoleDefault") })
  }

  return (
    <header className="w-full border-b border-border bg-white/80 backdrop-blur-sm sticky top-0 z-50 px-4 md:px-6 py-2.5 flex items-center justify-end gap-3">
      <Button
        onClick={toggleMobileSidebar}
        aria-label={t("openMenu")}
        size="sm"
        variantType="ghost"
        className="mr-auto p-2 md:hidden"
      >
        <Menu className="w-5 h-5" />
      </Button>

      {!isOnline && (
        <Badge variant="destructive" className="mr-2 gap-1.5 animate-pulse">
          <WifiOff className="w-3.5 h-3.5" />
          {t("offlineStatus")}
        </Badge>
      )}

      <LanguageSwitcher />

      <NotificationBell />

      <div className="h-5 w-px bg-border" />

      <div className="flex items-center gap-2.5">
        <div className="w-8 h-8 rounded-lg bg-primary/8 flex items-center justify-center text-primary text-xs font-black">
          {email ? email.charAt(0).toUpperCase() : "U"}
        </div>
        <div className="hidden sm:flex flex-col items-start">
          <span className="text-xs font-semibold text-foreground leading-tight">
            {email || t("defaultUserEmail")}
          </span>
          <div className="flex items-center gap-1">
            <ShieldCheck className="w-3 h-3 text-secondary" />
            <span className="text-[10px] text-muted font-medium">
              {translateRole(role)}
            </span>
          </div>
        </div>
      </div>
    </header>
  )
}
