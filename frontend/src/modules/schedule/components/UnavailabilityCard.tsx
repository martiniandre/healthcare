import { useTranslation } from "react-i18next"
import { CalendarOff, Trash2 } from "lucide-react"
import { formatTime, formatDate } from "../../../shared/utils/dates"
import { Button } from "../../../shared/components/ui/Button"
import type { StaffUnavailability } from "../types"

interface UnavailabilityCardProps {
  unavailability: StaffUnavailability
  onDelete: (unavailabilityId: string) => void
}

export const UnavailabilityCard = ({ unavailability, onDelete }: UnavailabilityCardProps) => {
  const { t } = useTranslation("schedule")

  return (
    <div className="bg-danger-soft border border-danger/20 rounded-xl p-4 flex flex-col gap-2">
      <div className="flex items-center justify-between gap-2">
        <span className="flex items-center gap-1.5 text-xs font-bold text-danger">
          <CalendarOff className="w-3.5 h-3.5" />
          {t("unavailability.badge")}
        </span>
        <span className="font-mono text-xs font-bold text-danger">
          {formatTime(unavailability.starts_at)} — {formatTime(unavailability.ends_at)}
        </span>
      </div>
      <div className="flex flex-col gap-0.5">
        <span className="text-xs text-danger">
          {formatDate(unavailability.starts_at)}
        </span>
        {unavailability.reason && (
          <span className="text-xs text-danger truncate">
            {unavailability.reason}
          </span>
        )}
      </div>
      <Button
        type="button"
        variantType="danger"
        size="sm"
        onClick={() => onDelete(unavailability.id)}
        className="mt-1 self-start"
      >
        <Trash2 className="w-3.5 h-3.5" />
        {t("unavailability.delete")}
      </Button>
    </div>
  )
}
