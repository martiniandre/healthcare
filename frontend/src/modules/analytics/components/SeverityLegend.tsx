import { cn } from "../../../shared/utils/cn"
import { toneClassNames, type Severity } from "../dashboard_thresholds"

export interface SeverityBand {
  severity: Severity
  label: string
  range: string
}

interface SeverityLegendProps {
  bands: SeverityBand[]
  className?: string
}

export const SeverityLegend = ({ bands, className }: SeverityLegendProps) => {
  return (
    <ul className={cn("flex flex-wrap items-center gap-x-4 gap-y-1.5", className)}>
      {bands.map((band) => (
        <li key={band.severity} className="flex items-center gap-1.5">
          <span
            className={cn("size-2 shrink-0 rounded-full", toneClassNames[band.severity].bar)}
            aria-hidden="true"
          />
          <span className="text-[11px] font-bold text-gray-700">{band.label}</span>
          <span className="text-[11px] tabular-nums text-muted-foreground">{band.range}</span>
        </li>
      ))}
    </ul>
  )
}
