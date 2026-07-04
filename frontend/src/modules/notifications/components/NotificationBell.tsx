import { useState, useRef, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { Bell } from "lucide-react"
import { useUnreadCountQuery, useNotificationsQuery, useMarkReadMutation } from "../queries"
import { NotificationItem } from "./NotificationItem"

export function NotificationBell() {
  const { t } = useTranslation()
  const [isOpen, setIsOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const { data: unreadData } = useUnreadCountQuery()
  const { data: notificationsData } = useNotificationsQuery(20, 0)
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
    markReadMutation.mutate(notificationId)
  }

  const notifications = notificationsData?.notifications ?? []

  return (
    <div ref={dropdownRef} className="relative">
      <button
        title={t("header.notificationTooltip")}
        onClick={() => setIsOpen(!isOpen)}
        className="relative p-2 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-50 transition-colors"
      >
        <Bell className="w-4 h-4" />
        {unreadCount > 0 && (
          <span className="absolute -top-0.5 -right-0.5 w-4 h-4 bg-red-500 text-white text-[9px] font-bold rounded-full flex items-center justify-center">
            {unreadCount > 9 ? "9+" : unreadCount}
          </span>
        )}
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 bg-white rounded-lg shadow-xl border border-border z-50 max-h-96 flex flex-col">
          <div className="px-4 py-2.5 border-b border-border flex items-center justify-between shrink-0">
            <h3 className="text-sm font-semibold text-gray-800">{t("notifications.title")}</h3>
            {unreadCount > 0 && (
              <span className="text-[10px] font-medium text-gray-500">
                {unreadCount} {t("notifications.unread")}
              </span>
            )}
          </div>

          <div className="overflow-y-auto flex-1">
            {notifications.length === 0 ? (
              <div className="px-4 py-8 text-center text-sm text-gray-400">
                {t("notifications.empty")}
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
