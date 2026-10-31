import i18next from "i18next"

const fallbackLocale = "pt-BR"

export const formatNumber = (
  value: number,
  options: Intl.NumberFormatOptions = { maximumFractionDigits: 2 },
  explicitLocale?: string,
): string => {
  const resolvedLocale = explicitLocale ?? i18next.language ?? fallbackLocale
  try {
    return new Intl.NumberFormat(resolvedLocale, options).format(value)
  } catch {
    return String(value)
  }
}
