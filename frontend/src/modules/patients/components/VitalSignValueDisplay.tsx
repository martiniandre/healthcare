import { Badge } from "../../../shared/components/ui/Badge"

interface VitalSignValueDisplayProps {
  notPerformed?: boolean
  valueQuantity: number
  valueUnit: string
  valueClassName?: string
}

export const VitalSignValueDisplay = ({ notPerformed, valueQuantity, valueUnit, valueClassName }: VitalSignValueDisplayProps) => {
  if (notPerformed) {
    return (
      <Badge variant="outline" className="bg-gray-100 text-gray-400 border-gray-200">
        N/A
      </Badge>
    )
  }
  return (
    <span className={valueClassName ?? "text-sm font-extrabold text-gray-800"}>
      {valueQuantity}
      <span className="text-xs text-muted font-normal ml-1">{valueUnit}</span>
    </span>
  )
}
