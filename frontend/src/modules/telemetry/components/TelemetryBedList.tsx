import { useTranslation } from "react-i18next"
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/components/ui/card"
import { StatusBadge } from "@/shared/components/ui/status-badge"
import { Button } from "@/shared/components/ui/button"
import { cn } from "@/shared/utils/cn"
import { Bell } from "lucide-react"
import { BedStatus, type CardiacCondition } from "../../../shared/types"
import { useClinicalStatusLabels } from "../../../shared/hooks/useClinicalStatusLabels"
import { resolveBedClinicalStatus, sortBedsBySeverity, type BedHistoryMap } from "../utils/bed_helpers"
import { BedTrendSparkline } from "./BedTrendSparkline"

interface Bed {
  id: string
  bedNumber: string
  patientName: string
  status: BedStatus
  condition: CardiacCondition
  bpm: number
  spo2: number
}

interface TelemetryBedListProps {
  beds: Bed[]
  selectedBedId: string | null
  setSelectedBedId: (id: string | null) => void
  history?: BedHistoryMap
}

export const TelemetryBedList = ({ beds, selectedBedId, setSelectedBedId, history = {} }: TelemetryBedListProps) => {
  const { t } = useTranslation("telemetry")
  const { cardiacConditionLabel } = useClinicalStatusLabels()
  const sortedBeds = sortBedsBySeverity(beds)

  return (
    <div className="flex flex-col gap-4 lg:col-span-1 text-left">
      <Card className="gap-4 py-4">
        <CardHeader className="border-b px-4">
          <CardTitle className="flex items-center gap-2 text-base font-extrabold">
            <Bell className="size-4 text-primary" aria-hidden="true" />
            {t("availableBeds", { count: beds.length })}
          </CardTitle>
        </CardHeader>

        <CardContent className="flex flex-col gap-3 px-4">
          {sortedBeds.map((bedItem) => {
            const isSelected = bedItem.id === selectedBedId
            const isCritical = bedItem.status === BedStatus.Danger
            const clinicalStatus = resolveBedClinicalStatus(bedItem.status)
            const bedHistory = history[bedItem.id] ?? []
            return (
              <Button
                key={bedItem.id}
                variant="outline"
                aria-pressed={isSelected}
                onClick={() => setSelectedBedId(bedItem.id)}
                className={cn(
                  "h-auto w-full flex-col items-stretch justify-start gap-0 whitespace-normal rounded-xl p-3.5 text-left font-normal",
                  isSelected && "border-primary bg-primary/5",
                  !isSelected && bedItem.status === BedStatus.Danger && "border-destructive/30 bg-destructive/5"
                )}
              >
                <div className="flex items-center justify-between">
                  <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
                    {bedItem.bedNumber}
                  </span>
                  <StatusBadge
                    status={clinicalStatus}
                    className={cn("text-[10px] font-black uppercase", isCritical && "animate-pulse-glow motion-reduce:animate-none")}
                  >
                    {t(`status.${clinicalStatus}`)}
                  </StatusBadge>
                </div>

                <h4 className="mt-1 text-sm font-bold text-foreground">
                  {bedItem.patientName}
                </h4>
                <span className="text-[11px] text-muted-foreground">{cardiacConditionLabel(bedItem.condition)}</span>

                <div className="mt-3 flex items-center justify-between text-xs">
                  <span className="text-muted-foreground">{t("sensors.ecg")}: <strong className="tabular-nums text-foreground">{bedItem.bpm > 0 ? `${bedItem.bpm} BPM` : "---"}</strong></span>
                  <span className="text-muted-foreground">{t("sensors.spo2")}: <strong className="tabular-nums text-foreground">{bedItem.spo2 > 0 ? `${bedItem.spo2}%` : "---"}</strong></span>
                </div>

                <div className="mt-2 flex items-center justify-between gap-2">
                  <BedTrendSparkline
                    values={bedHistory.map((sample) => sample.bpm)}
                    metricLabel={t("sensors.ecg")}
                    unit=" BPM"
                  />
                  <BedTrendSparkline
                    values={bedHistory.map((sample) => sample.spo2)}
                    metricLabel={t("sensors.spo2")}
                    unit="%"
                    className="text-success"
                  />
                </div>
              </Button>
            )
          })}
        </CardContent>
      </Card>
    </div>
  )
}
