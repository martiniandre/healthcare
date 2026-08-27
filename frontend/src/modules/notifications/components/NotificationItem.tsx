import { cn } from "../../../shared/utils/cn"
import { Button } from "../../../shared/components/ui/Button"
import { formatRelativeTime } from "../../../shared/utils/dates"
import { useLocale } from "../../../shared/hooks/useLocale"
import type { NotificationItem as NotificationItemType, NotificationPriority } from "../types"

const priorityConfig: Record<NotificationPriority, { dot: string; bg: string }> = {
  critical: { dot: "bg-danger", bg: "bg-danger-soft" },
  high: { dot: "bg-warning", bg: "bg-warning-soft" },
  medium: { dot: "bg-info", bg: "bg-info-soft" },
  low: { dot: "bg-muted-foreground", bg: "bg-muted-soft" },
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
    <Button
      onClick={() => onMarkRead(notification.id)}
      variantType="ghost"
      className={cn(
        "w-full justify-start text-left px-4 py-3 rounded-none border-b border-border last:border-b-0 hover:bg-muted-soft transition-colors h-auto",
        !notification.is_read && config.bg,
      )}
    >
      <div className="flex items-start gap-3 w-full">
        <div className={cn("w-2 h-2 rounded-full mt-1.5 shrink-0", config.dot)} />
        <div className="flex-1 min-w-0">
          <p
            className={cn(
              "text-sm truncate",
              !notification.is_read ? "font-semibold text-foreground" : "font-medium text-foreground",
            )}
          >
            {notification.title}
          </p>
          <p className="text-xs text-muted-foreground mt-0.5 line-clamp-2">{notification.body}</p>
          <p className="text-[10px] text-muted-foreground mt-1">{timeAgo}</p>
        </div>
      </div>
    </Button>
  )
}
