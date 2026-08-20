import { cn } from "../../../shared/utils/cn"
import { formatRelativeTime } from "../../../shared/utils/dates"
import { useLocale } from "../../../shared/hooks/useLocale"
import type { NotificationItem as NotificationItemType, NotificationPriority } from "../types"

const priorityConfig: Record<NotificationPriority, { dot: string; bg: string }> = {
  critical: { dot: "bg-red-500", bg: "bg-red-50" },
  high: { dot: "bg-orange-500", bg: "bg-orange-50" },
  medium: { dot: "bg-blue-500", bg: "bg-blue-50" },
  low: { dot: "bg-gray-400", bg: "bg-gray-50" },
}

interface NotificationItemProps {
  notification: NotificationItemType
  onMarkRead: (id: string) => void
}

export function NotificationItem({ notification, onMarkRead }: NotificationItemProps) {
  const config = priorityConfig[notification.priority] ?? priorityConfig.low
  const locale = useLocale()
  const timeAgo = formatRelativeTime(notification.created_at, locale)

  return (
    <button
      onClick={() => onMarkRead(notification.id)}
      className={cn(
        "w-full text-left px-4 py-3 border-b border-border last:border-b-0 hover:bg-gray-50 transition-colors",
        !notification.is_read && config.bg,
      )}
    >
      <div className="flex items-start gap-3">
        <div className={cn("w-2 h-2 rounded-full mt-1.5 shrink-0", config.dot)} />
        <div className="flex-1 min-w-0">
          <p
            className={cn(
              "text-sm truncate",
              !notification.is_read ? "font-semibold text-gray-900" : "font-medium text-gray-700",
            )}
          >
            {notification.title}
          </p>
          <p className="text-xs text-gray-500 mt-0.5 line-clamp-2">{notification.body}</p>
          <p className="text-[10px] text-gray-400 mt-1">{timeAgo}</p>
        </div>
      </div>
    </button>
  )
}
