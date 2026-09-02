import i18next from "i18next"

const fallbackLocale = "pt-BR"

export type DateInput = string | number | Date

const resolveLocale = (explicitLocale?: string): string => {
  return explicitLocale ?? i18next.language ?? fallbackLocale
}

const dateOnlyPattern = /^\d{4}-\d{2}-\d{2}$/

const parseDateOnlyString = (value: string): Date | null => {
  const [year, month, day] = value.split("-").map(Number)
  if (month < 1 || month > 12 || day < 1 || day > 31) {
    return null
  }
  const parsedDate = new Date(year, month - 1, day)
  if (
    parsedDate.getFullYear() !== year ||
    parsedDate.getMonth() !== month - 1 ||
    parsedDate.getDate() !== day
  ) {
    return null
  }
  return parsedDate
}

const toDate = (value: DateInput): Date | null => {
  const parsedDate = typeof value === "string" && dateOnlyPattern.test(value)
    ? parseDateOnlyString(value)
    : new Date(value)
  return parsedDate && !Number.isNaN(parsedDate.getTime()) ? parsedDate : null
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

const localizedDigitsPattern = /^\d{8}$/

const isValidDateComponents = (year: number, month: number, day: number): boolean => {
  if (year < 1 || month < 1 || month > 12 || day < 1 || day > 31) {
    return false
  }
  const parsedDate = new Date(year, month - 1, day)
  return (
    parsedDate.getFullYear() === year &&
    parsedDate.getMonth() === month - 1 &&
    parsedDate.getDate() === day
  )
}

export const getDateOrder = (explicitLocale?: string): "day-first" | "month-first" => {
  const positionTokens = new Intl.DateTimeFormat(resolveLocale(explicitLocale), {
    day: "numeric",
    month: "numeric",
    year: "numeric",
  }).formatToParts(new Date(2026, 0, 15))
  let dayIndex = -1
  let monthIndex = -1
  positionTokens.forEach((part, index) => {
    if (part.type === "day") {
      dayIndex = index
    }
    if (part.type === "month") {
      monthIndex = index
    }
  })
  return dayIndex < monthIndex ? "day-first" : "month-first"
}

export const localizedDateToIso = (value: string, explicitLocale?: string): string => {
  if (dateOnlyPattern.test(value)) {
    const [year, month, day] = value.split("-").map(Number)
    return isValidDateComponents(year, month, day) ? value : ""
  }
  const digitsOnly = value.replace(/\D/g, "")
  if (!localizedDigitsPattern.test(digitsOnly)) {
    return ""
  }
  const firstPair = Number(digitsOnly.slice(0, 2))
  const secondPair = Number(digitsOnly.slice(2, 4))
  const year = Number(digitsOnly.slice(4, 8))
  const dayFirst = getDateOrder(explicitLocale) === "day-first"
  const month = dayFirst ? secondPair : firstPair
  const day = dayFirst ? firstPair : secondPair
  if (!isValidDateComponents(year, month, day)) {
    return ""
  }
  return `${year}-${String(month).padStart(2, "0")}-${String(day).padStart(2, "0")}`
}

export const isPastLocalizedDate = (value: string, explicitLocale?: string): boolean => {
  const isoDate = localizedDateToIso(value, explicitLocale)
  if (!isoDate) {
    return false
  }
  return new Date(isoDate).getTime() < Date.now()
}
