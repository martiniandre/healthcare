import { cn } from "../../../shared/utils/cn"
import { StatusBadge } from "../../../shared/components/ui/status-badge"
import { toneClinicalStatus, type Severity } from "../metric_tones"

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
          <StatusBadge status={toneClinicalStatus[band.severity]} className="text-[11px] font-bold">
            {band.label}
          </StatusBadge>
          <span className="text-[11px] tabular-nums text-muted-foreground">{band.range}</span>
        </li>
      ))}
    </ul>
  )
}
