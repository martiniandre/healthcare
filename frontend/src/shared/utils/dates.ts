import i18next from "i18next"

const fallbackLocale = "pt-BR"

export type DateInput = string | number | Date

const resolveLocale = (explicitLocale?: string): string => {
  return explicitLocale ?? i18next.language ?? fallbackLocale
}

const toDate = (value: DateInput): Date | null => {
  const parsedDate = new Date(value)
  return Number.isNaN(parsedDate.getTime()) ? null : parsedDate
}

const formatWithIntl = (value: DateInput, options: Intl.DateTimeFormatOptions, explicitLocale?: string): string => {
  const parsedDate = toDate(value)
  if (!parsedDate) {
    return String(value)
  }
  try {
    return new Intl.DateTimeFormat(resolveLocale(explicitLocale), options).format(parsedDate)
  } catch {
    return parsedDate.toLocaleString()
  }
}

export const formatDate = (value: DateInput, explicitLocale?: string): string => {
  return formatWithIntl(value, { day: "numeric", month: "numeric", year: "numeric" }, explicitLocale)
}

export const formatLongDate = (value: DateInput, explicitLocale?: string): string => {
  return formatWithIntl(value, { day: "numeric", month: "short", year: "numeric" }, explicitLocale)
}

export const formatTime = (value: DateInput, explicitLocale?: string): string => {
  return formatWithIntl(value, { hour: "2-digit", minute: "2-digit" }, explicitLocale)
}

export const formatDateTime = (value: DateInput, explicitLocale?: string): string => {
  return formatWithIntl(
    value,
    { day: "numeric", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit" },
    explicitLocale
  )
}

export const formatRelativeTime = (value: DateInput, explicitLocale?: string): string => {
  const parsedDate = toDate(value)
  if (!parsedDate) {
    return String(value)
  }
  const activeLocale = resolveLocale(explicitLocale)
  const elapsedMinutes = Math.floor((Date.now() - parsedDate.getTime()) / 60000)

  if (elapsedMinutes < 1) {
    return i18next.t("notifications:relativeNow")
  }

  try {
    const relativeFormat = new Intl.RelativeTimeFormat(activeLocale, { numeric: "auto" })
    if (elapsedMinutes < 60) {
      return relativeFormat.format(-elapsedMinutes, "minute")
    }
    const elapsedHours = Math.floor(elapsedMinutes / 60)
    if (elapsedHours < 24) {
      return relativeFormat.format(-elapsedHours, "hour")
    }
    const elapsedDays = Math.floor(elapsedHours / 24)
    if (elapsedDays < 7) {
      return relativeFormat.format(-elapsedDays, "day")
    }
  } catch {
    return formatDate(parsedDate, activeLocale)
  }
  return formatDate(parsedDate, activeLocale)
}
