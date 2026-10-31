import { useCallback } from "react"
import { useLocale } from "../../../shared/hooks/useLocale"
import { formatNumber } from "../../../shared/utils/numbers"

export const useNumberFormat = (): ((value: number) => string) => {
  const locale = useLocale()

  return useCallback((value: number): string => formatNumber(value, undefined, locale), [locale])
}
