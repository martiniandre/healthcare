import { useCallback } from "react"
import { useTranslation } from "react-i18next"
import i18next from "i18next"

import type { NotificationItem } from "../types"

const resolveTemplateText = (
  translationKey: string | undefined,
  legacyText: string,
  params: Record<string, unknown> | undefined,
): string => {
  if (!translationKey || !i18next.exists(translationKey)) {
    return legacyText
  }
  return i18next.t(translationKey, { defaultValue: legacyText, ...params }) as string
}

export interface NotificationTextLabels {
  notificationTitle: (notification: NotificationItem) => string
  notificationBody: (notification: NotificationItem) => string
}

export const useNotificationText = (): NotificationTextLabels => {
  const { i18n } = useTranslation()

  const notificationTitle = useCallback(
    (notification: NotificationItem): string =>
      resolveTemplateText(notification.title_key, notification.title, notification.params),
    [i18n],
  )

  const notificationBody = useCallback(
    (notification: NotificationItem): string =>
      resolveTemplateText(notification.body_key, notification.body, notification.params),
    [i18n],
  )

  return { notificationTitle, notificationBody }
}
