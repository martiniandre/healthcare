import { useTranslation } from "react-i18next"
import { CalendarOff, Trash2 } from "lucide-react"
import type { StaffUnavailability } from "../types"

interface UnavailabilityCardProps {
  unavailability: StaffUnavailability
  onDelete: (unavailabilityId: string) => void
}

const formatTime = (dateTimeValue: string): string => {
  return new Date(dateTimeValue).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
}

const formatDate = (dateTimeValue: string): string => {
  return new Date(dateTimeValue).toLocaleDateString()
}

export const UnavailabilityCard = ({ unavailability, onDelete }: UnavailabilityCardProps) => {
  const { t } = useTranslation("schedule")

  return (
    <div className="bg-red-50 border border-red-200 rounded-xl p-4 flex flex-col gap-2">
      <div className="flex items-center justify-between gap-2">
        <span className="flex items-center gap-1.5 text-xs font-bold text-red-700">
          <CalendarOff className="w-3.5 h-3.5" />
          {t("unavailability.badge")}
        </span>
        <span className="font-mono text-xs font-bold text-red-800">
          {formatTime(unavailability.starts_at)} — {formatTime(unavailability.ends_at)}
        </span>
      </div>
      <div className="flex flex-col gap-0.5">
        <span className="text-xs text-red-600">
          {formatDate(unavailability.starts_at)}
        </span>
        {unavailability.reason && (
          <span className="text-xs text-red-500 truncate">
            {unavailability.reason}
          </span>
        )}
      </div>
      <button
        onClick={() => onDelete(unavailability.id)}
        className="mt-1 self-start inline-flex items-center gap-1.5 text-xs font-semibold text-red-500 hover:text-red-700"
      >
        <Trash2 className="w-3.5 h-3.5" />
        {t("unavailability.delete")}
      </button>
    </div>
  )
}
