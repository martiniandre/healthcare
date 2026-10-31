import { useTranslation } from "react-i18next"
import { cn } from "@/shared/utils/cn"
import { buildSparklinePoints } from "../utils/bed_helpers"

interface BedTrendSparklineProperties {
  values: number[]
  metricLabel: string
  unit: string
  className?: string
}

export const BedTrendSparkline = ({ values, metricLabel, unit, className }: BedTrendSparklineProperties) => {
  const { t } = useTranslation("telemetry")

  if (values.length < 2) {
    return null
  }

  const firstValue = values[0]
  const lastValue = values[values.length - 1]

  return (
    <svg
      role="img"
      aria-label={t("trend.ariaLabel", {
        metric: metricLabel,
        count: values.length,
        first: `${firstValue}${unit}`,
        last: `${lastValue}${unit}`,
      })}
      viewBox="0 0 80 20"
      className={cn("h-5 w-20 text-primary", className)}
    >
      <polyline
        points={buildSparklinePoints(values)}
        fill="none"
        stroke="currentColor"
        strokeWidth={1.5}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}
