import { useState, useRef, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { Bell, Loader2 } from "lucide-react"
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
      <button
        title={t("header:notificationTooltip")}
        onClick={() => setIsOpen(!isOpen)}
        className="relative p-2 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-50 transition-colors"
      >
        <Bell className="w-4 h-4" />
        {isUnreadLoading && (
          <span className="absolute -top-0.5 -right-0.5 w-4 h-4 flex items-center justify-center">
            <Loader2 className="w-3 h-3 animate-spin text-gray-400" />
          </span>
        )}
        {!isUnreadLoading && unreadCount > 0 && (
          <span className="absolute -top-0.5 -right-0.5 w-4 h-4 bg-red-500 text-white text-[9px] font-bold rounded-full flex items-center justify-center">
            {unreadCount > 9 ? "9+" : unreadCount}
          </span>
        )}
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 bg-card rounded-lg shadow-xl border border-border z-50 max-h-96 flex flex-col">
          <div className="px-4 py-2.5 border-b border-border flex items-center justify-between shrink-0">
            <h3 className="text-sm font-semibold text-gray-800">{t("title")}</h3>
            {unreadCount > 0 && (
              <span className="text-[10px] font-medium text-gray-500">
                {unreadCount} {t("unread")}
              </span>
            )}
          </div>

          <div className="overflow-y-auto flex-1">
            {isListLoading ? (
              <div className="px-4 py-8 flex justify-center">
                <Loader2 className="w-5 h-5 animate-spin text-gray-300" />
              </div>
            ) : hasListError ? (
              <div className="px-4 py-8 text-center text-sm text-red-500">
                {t("listError")}
              </div>
            ) : notifications.length === 0 ? (
              <div className="px-4 py-8 text-center text-sm text-gray-400">
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
