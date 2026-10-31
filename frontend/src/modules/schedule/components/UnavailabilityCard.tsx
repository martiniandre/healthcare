import { useTranslation } from "react-i18next"
import { CalendarOff, Trash2 } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import { Card, CardContent } from "@/shared/components/ui/card"
import { formatTime, formatDate } from "../../../shared/utils/dates"
import type { StaffUnavailability } from "../types"

interface UnavailabilityCardProps {
  unavailability: StaffUnavailability
  onDelete: (unavailabilityId: string) => void
}

export const UnavailabilityCard = ({ unavailability, onDelete }: UnavailabilityCardProps) => {
  const { t } = useTranslation("schedule")

  return (
    <Card className="py-4 border-destructive/30 bg-destructive/5">
      <CardContent className="flex flex-col gap-2 px-4">
        <div className="flex items-center justify-between gap-2">
          <span className="flex items-center gap-1.5 text-xs font-bold text-destructive">
            <CalendarOff className="size-3.5" />
            {t("unavailability.badge")}
          </span>
          <span className="font-mono text-xs font-bold tabular-nums text-destructive">
            {formatTime(unavailability.starts_at)} — {formatTime(unavailability.ends_at)}
          </span>
        </div>
        <div className="flex flex-col gap-0.5">
          <span className="text-xs text-destructive">
            {formatDate(unavailability.starts_at)}
          </span>
          {unavailability.reason && (
            <span className="text-xs text-destructive/80 truncate">
              {unavailability.reason}
            </span>
          )}
        </div>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={() => onDelete(unavailability.id)}
          className="self-start text-xs font-semibold text-destructive hover:text-destructive"
        >
          <Trash2 />
          {t("unavailability.delete")}
        </Button>
      </CardContent>
    </Card>
  )
}
