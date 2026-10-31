import { useTranslation } from "react-i18next"
import { AlertCircle, Bell, BellOff } from "lucide-react"
import { useUnreadCountQuery, useNotificationsQuery, useMarkReadMutation } from "../queries"
import { NotificationItem } from "./NotificationItem"
import { useLocale } from "../../../shared/hooks/useLocale"
import { toast } from "../../../shared/store/toast_store"
import { Badge } from "../../../shared/components/ui/badge"
import { Button } from "../../../shared/components/ui/button"
import { Popover, PopoverContent, PopoverTrigger } from "../../../shared/components/ui/popover"
import { Alert, AlertDescription } from "../../../shared/components/ui/alert"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import { Skeleton } from "../../../shared/components/ui/skeleton"
import { Spinner } from "../../../shared/components/ui/spinner"
import { Tooltip, TooltipContent, TooltipTrigger } from "../../../shared/components/ui/tooltip"

export function NotificationBell() {
  const { t } = useTranslation("notifications")
  const locale = useLocale()
  const formatCount = (value: number): string => new Intl.NumberFormat(locale).format(value)

  const { data: unreadData, isLoading: isUnreadLoading } = useUnreadCountQuery()
  const { data: notificationsData, isLoading: isListLoading, isError: hasListError, refetch: refetchList } = useNotificationsQuery(20, 0)
  const markReadMutation = useMarkReadMutation()

  const unreadCount = unreadData?.count ?? 0

  const handleMarkRead = (notificationId: string) => {
    markReadMutation.mutate(notificationId, {
      onError: () => {
        toast.error(t("markReadError"))
      },
    })
  }

  const notifications = notificationsData?.notifications ?? []

  return (
    <Popover>
      <Tooltip>
        <TooltipTrigger asChild>
          <PopoverTrigger asChild>
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="relative min-h-11 min-w-11 text-muted-foreground"
              aria-label={t("header:notificationTooltip")}
            >
              <Bell aria-hidden="true" />
              {isUnreadLoading && (
                <span className="absolute -top-0.5 -right-0.5 flex size-4 items-center justify-center">
                  <Spinner className="size-3 text-muted-foreground" />
                </span>
              )}
              {!isUnreadLoading && unreadCount > 0 && (
                <Badge
                  variant="destructive"
                  className="absolute -top-0.5 -right-0.5 size-4 justify-center p-0 text-[9px] font-bold"
                >
                  {unreadCount > 9 ? `${formatCount(9)}+` : formatCount(unreadCount)}
                </Badge>
              )}
            </Button>
          </PopoverTrigger>
        </TooltipTrigger>
        <TooltipContent>{t("header:notificationTooltip")}</TooltipContent>
      </Tooltip>

      <PopoverContent align="end" className="flex max-h-96 w-[min(20rem,calc(100vw-2rem))] flex-col gap-0 p-0">
        <div className="px-4 py-2.5 border-b border-border flex items-center justify-between shrink-0">
          <h3 className="truncate text-sm font-semibold text-foreground">{t("title")}</h3>
          {unreadCount > 0 && (
            <span className="ml-2 shrink-0 text-[10px] font-medium tabular-nums text-muted-foreground">
              {t("unreadCount", { formattedCount: formatCount(unreadCount) })}
            </span>
          )}
        </div>

        <div className="overflow-y-auto flex-1">
          {isListLoading ? (
            <div className="flex flex-col gap-3 px-4 py-4" data-testid="notifications-skeleton" aria-busy="true">
              {[0, 1, 2].map((item) => (
                <div key={item} className="flex flex-col gap-2">
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-3 w-full" />
                </div>
              ))}
            </div>
          ) : hasListError ? (
            <div className="p-3">
              <Alert variant="destructive">
                <AlertCircle aria-hidden="true" />
                <AlertDescription>
                  <p>{t("listError")}</p>
                  <Button type="button" variant="outline" className="mt-2 min-h-11 min-w-11" onClick={() => refetchList()}>
                    {t("retry")}
                  </Button>
                </AlertDescription>
              </Alert>
            </div>
          ) : notifications.length === 0 ? (
            <div className="p-3">
              <EmptyState icon={BellOff} title={t("empty")} className="py-8" />
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
      </PopoverContent>
    </Popover>
  )
}
