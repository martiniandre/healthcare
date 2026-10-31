import { useTranslation } from "react-i18next"
import { StatusBadge } from "../../../shared/components/ui/status-badge"
import { cn } from "../../../shared/utils/cn"

interface VitalSignValueDisplayProps {
  notPerformed?: boolean
  valueQuantity: number
  valueUnit: string
  valueClassName?: string
}

export const VitalSignValueDisplay = ({ notPerformed, valueQuantity, valueUnit, valueClassName }: VitalSignValueDisplayProps) => {
  const { t } = useTranslation("patients")

  if (notPerformed) {
    return <StatusBadge status="neutral">{t("header.notPerformed")}</StatusBadge>
  }
  return (
    <span className={cn("tabular-nums", valueClassName ?? "text-sm font-extrabold text-gray-800")}>
      {valueQuantity}
      <span className="text-xs text-muted-foreground font-normal ml-1">{valueUnit}</span>
    </span>
  )
}
