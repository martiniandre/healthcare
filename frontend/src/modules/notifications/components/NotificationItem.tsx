import { useTranslation } from "react-i18next"
import { Button } from "../../../shared/components/ui/button"
import { StatusBadge, type ClinicalStatus } from "../../../shared/components/ui/status-badge"
import { cn } from "../../../shared/utils/cn"
import { formatRelativeTime } from "../../../shared/utils/dates"
import { useLocale } from "../../../shared/hooks/useLocale"
import { useNotificationText } from "../hooks/useNotificationText"
import type { NotificationItem as NotificationItemType, NotificationPriority } from "../types"

const priorityConfig: Record<NotificationPriority, { bg: string; status: ClinicalStatus }> = {
  critical: { bg: "bg-destructive/10", status: "critical" },
  high: { bg: "bg-warning/10", status: "attention" },
  medium: { bg: "bg-primary/10", status: "info" },
  low: { bg: "bg-muted", status: "neutral" },
}

interface NotificationItemProps {
  notification: NotificationItemType
  onMarkRead: (id: string) => void
}

export function NotificationItem({ notification, onMarkRead }: NotificationItemProps) {
  const { t } = useTranslation("notifications")
  const priority: NotificationPriority = notification.priority in priorityConfig ? notification.priority : "low"
  const config = priorityConfig[priority]
  const locale = useLocale()
  const { notificationTitle, notificationBody } = useNotificationText()
  const timeAgo = formatRelativeTime(notification.created_at, locale)
  const title = notificationTitle(notification)

  return (
    <Button
      type="button"
      variant="ghost"
      onClick={() => onMarkRead(notification.id)}
      className={cn(
        "h-auto min-h-11 w-full items-start justify-start rounded-none border-b border-border px-4 py-3 text-left whitespace-normal last:border-b-0",
        !notification.is_read && config.bg,
      )}
    >
      <div className="flex w-full min-w-0 flex-col gap-1">
        <div className="flex items-start justify-between gap-2">
          <p
            title={title}
            className={cn(
              "min-w-0 flex-1 truncate text-sm",
              !notification.is_read ? "font-semibold text-foreground" : "font-medium text-muted-foreground",
            )}
          >
            {title}
          </p>
          <StatusBadge status={config.status} className="shrink-0 text-[10px]">
            {t(`priority.${priority}`)}
          </StatusBadge>
        </div>
        <p className="line-clamp-2 break-words text-xs font-normal text-muted-foreground">
          {notificationBody(notification)}
        </p>
        <p className="text-[10px] font-normal tabular-nums text-muted-foreground">{timeAgo}</p>
      </div>
    </Button>
  )
}
