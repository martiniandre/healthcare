export const formatDecimal = (value: number, locale: string, fractionDigits = 0): string => {
  try {
    return new Intl.NumberFormat(locale, {
      minimumFractionDigits: fractionDigits,
      maximumFractionDigits: fractionDigits,
    }).format(value)
  } catch {
    return value.toFixed(fractionDigits)
  }
}

export const formatPercentage = (percentageValue: number, locale: string, fractionDigits = 0): string => {
  try {
    return new Intl.NumberFormat(locale, {
      style: "percent",
      minimumFractionDigits: fractionDigits,
      maximumFractionDigits: fractionDigits,
    }).format(percentageValue / 100)
  } catch {
    return `${percentageValue.toFixed(fractionDigits)}%`
  }
}
