import { Badge } from "../../../shared/components/ui/badge"
import { useNumberFormat } from "../hooks/useNumberFormat"

interface PortalVitalValueProps {
  notPerformed?: boolean
  valueQuantity: number
  valueUnit: string
  valueClassName?: string
}

export const PortalVitalValue = ({ notPerformed, valueQuantity, valueUnit, valueClassName }: PortalVitalValueProps) => {
  const formatNumber = useNumberFormat()

  if (notPerformed) {
    return (
      <Badge variant="outline" className="text-muted-foreground">
        N/A
      </Badge>
    )
  }

  return (
    <span className={valueClassName ?? "text-sm font-extrabold text-foreground"}>
      <span className="tabular-nums">{formatNumber(valueQuantity)}</span>
      <span className="ml-1 text-xs font-normal text-muted-foreground">{valueUnit}</span>
    </span>
  )
}
