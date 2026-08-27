import { useState, useRef, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { Bell, Loader2 } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"
import { useUnreadCountQuery, useNotificationsQuery, useMarkReadMutation } from "../queries"
import { NotificationItem } from "./NotificationItem"
import { toast } from "../../../shared/store/toast_store"

export function NotificationBell() {
  const { t } = useTranslation("notifications")
  const [isOpen, setIsOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const { data: unreadData, isLoading: isUnreadLoading } = useUnreadCountQuery()
  const { data: notificationsData, isLoading: isListLoading, isError: hasListError } = useNotificationsQuery(20, 0)
  const markReadMutation = useMarkReadMutation()

  const unreadCount = unreadData?.count ?? 0

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [])

  const handleMarkRead = (notificationId: string) => {
    markReadMutation.mutate(notificationId, {
      onError: () => {
        toast.error(t("markReadError"))
      },
    })
  }

  const notifications = notificationsData?.notifications ?? []

  return (
    <div ref={dropdownRef} className="relative">
      <Button
        title={t("header:notificationTooltip")}
        onClick={() => setIsOpen(!isOpen)}
        variantType="ghost"
        size="sm"
        className="relative p-2 text-muted-foreground hover:text-foreground rounded-lg"
      >
        <Bell className="w-4 h-4" />
        {isUnreadLoading && (
          <span className="absolute -top-0.5 -right-0.5 w-4 h-4 flex items-center justify-center">
            <Loader2 className="w-3 h-3 animate-spin text-muted-foreground" />
          </span>
        )}
        {!isUnreadLoading && unreadCount > 0 && (
          <span className="absolute -top-0.5 -right-0.5 w-4 h-4 bg-danger text-background text-[9px] font-bold rounded-full flex items-center justify-center">
            {unreadCount > 9 ? "9+" : unreadCount}
          </span>
        )}
      </Button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 bg-surface rounded-lg shadow-xl border border-border z-50 max-h-96 flex flex-col">
          <div className="px-4 py-2.5 border-b border-border flex items-center justify-between shrink-0">
            <h3 className="text-sm font-semibold text-foreground">{t("title")}</h3>
            {unreadCount > 0 && (
              <span className="text-[10px] font-medium text-muted-foreground">
                {unreadCount} {t("unread")}
              </span>
            )}
          </div>

          <div className="overflow-y-auto flex-1">
            {isListLoading ? (
              <div className="px-4 py-8 flex justify-center">
                <Loader2 className="w-5 h-5 animate-spin text-muted-foreground" />
              </div>
            ) : hasListError ? (
              <div className="px-4 py-8 text-center text-sm text-danger">
                {t("listError")}
              </div>
            ) : notifications.length === 0 ? (
              <div className="px-4 py-8 text-center text-sm text-muted-foreground">
                {t("empty")}
              </div>
            ) : (
              notifications.map((notification) => (
                <NotificationItem
                  key={notification.id}
                  notification={notification}
                  onMarkRead={handleMarkRead}
                />
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}
