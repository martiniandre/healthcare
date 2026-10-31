import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Card, CardContent, CardFooter, CardHeader } from "@/shared/components/ui/card"
import { Badge } from "@/shared/components/ui/badge"
import { StatusBadge } from "@/shared/components/ui/status-badge"
import { Button } from "@/shared/components/ui/button"
import { cn } from "@/shared/utils/cn"
import { Heart, Activity, Thermometer } from "lucide-react"
import { BedStatus, CardiacCondition } from "../../../shared/types"
import { type MutableRefObject } from "react"
import { resolveBedClinicalStatus } from "../utils/bed_helpers"

interface ActiveBed {
  id: string
  bedNumber: string
  patientName: string
  age: number
  gender: string
  bpm: number
  spo2: number
  temperature: number
  status: BedStatus
  condition: CardiacCondition
}

interface TelemetryBedMonitorProps {
  activeBed: ActiveBed | null
  activeRoomName?: string
  canvasRef: MutableRefObject<HTMLCanvasElement | null>
  updateSelectedBedCondition: (condition: CardiacCondition) => void
}

const SENSOR_LABEL_CLASS = "block text-[10px] font-bold uppercase tracking-wider text-muted-foreground"
const SENSOR_BOX_CLASS = "flex items-center justify-between rounded-xl border border-border bg-muted p-4"

export const TelemetryBedMonitor = ({
  activeBed,
  activeRoomName,
  canvasRef,
  updateSelectedBedCondition
}: TelemetryBedMonitorProps) => {
  const { t } = useTranslation("telemetry")

  if (!activeBed) {
    return (
      <Card className="flex-1 min-h-[400px] items-center justify-center gap-4 bg-muted/50 p-8 text-center">
        <div className="rounded-full border border-primary/10 bg-primary/5 p-4 text-primary">
          <Activity className="size-8" aria-hidden="true" />
        </div>
        <div className="max-w-md flex flex-col gap-1">
          <h3 className="text-base font-extrabold text-foreground">
            {t("noBedSelected")}
          </h3>
          <p className="text-xs leading-normal text-muted-foreground">
            {t("noBedSelectedDesc")}
          </p>
        </div>
      </Card>
    )
  }

  const conditionOptions: { condition: CardiacCondition; labelKey: string }[] = [
    { condition: CardiacCondition.Normal, labelKey: "simulation.normal" },
    { condition: CardiacCondition.Bradycardia, labelKey: "simulation.bradycardia" },
    { condition: CardiacCondition.Tachycardia, labelKey: "simulation.tachycardia" },
    { condition: CardiacCondition.CardiacArrest, labelKey: "simulation.cardiacArrest" },
  ]

  return (
    <Card className="gap-4 py-4">
      <CardHeader className="flex flex-wrap items-center justify-between gap-2 border-b px-4">
        <div className="text-left">
          <span className="text-xs font-bold uppercase tracking-wider text-primary">{activeRoomName}</span>
          <h3 className="text-base font-bold text-foreground">{activeBed.bedNumber} • {activeBed.patientName}</h3>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="secondary">
            {activeBed.age} {t("years")} · {activeBed.gender}
          </Badge>
          <StatusBadge
            status={resolveBedClinicalStatus(activeBed.status)}
            className={cn(activeBed.status === BedStatus.Danger && "animate-pulse-glow motion-reduce:animate-none")}
          >
            {t(`status.${resolveBedClinicalStatus(activeBed.status)}`)}
          </StatusBadge>
        </div>
      </CardHeader>

      <CardContent className="flex flex-col gap-4 px-4">
        <div className="relative overflow-hidden rounded-xl border border-border bg-slate-950 p-1">
          <canvas
            ref={canvasRef}
            width={700}
            height={260}
            className="block w-full max-w-full rounded-lg"
          />
          <div className="absolute top-4 left-4 rounded border border-emerald-500/20 bg-black/60 px-3 py-1.5 font-mono text-[10px] text-emerald-400 backdrop-blur-md">
            {t("ecgLead")}
          </div>
        </div>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div className={SENSOR_BOX_CLASS}>
            <div className="text-left">
              <span className={SENSOR_LABEL_CLASS}>{t("sensors.ecg")}</span>
              <span className="mt-1 block text-2xl font-black tabular-nums text-foreground">
                {activeBed.bpm > 0 ? `${activeBed.bpm} ` : "--- "}
                <span className="text-xs font-normal text-muted-foreground">BPM</span>
              </span>
            </div>
            <Heart className={cn("size-8 shrink-0 text-destructive", activeBed.bpm > 110 && "animate-pulse motion-reduce:animate-none")} />
          </div>

          <div className={SENSOR_BOX_CLASS}>
            <div className="text-left">
              <span className={SENSOR_LABEL_CLASS}>{t("sensors.spo2")}</span>
              <span className="mt-1 block text-2xl font-black tabular-nums text-foreground">
                {activeBed.spo2 > 0 ? `${activeBed.spo2}%` : "---%"}
              </span>
            </div>
            <Activity className="size-8 shrink-0 text-success" />
          </div>

          <div className={SENSOR_BOX_CLASS}>
            <div className="text-left">
              <span className={SENSOR_LABEL_CLASS}>{t("sensors.temp")}</span>
              <span className="mt-1 block text-2xl font-black tabular-nums text-foreground">
                {activeBed.temperature.toFixed(1)}°C
              </span>
            </div>
            <Thermometer className="size-8 shrink-0 text-primary" />
          </div>
        </div>
      </CardContent>

      <Can I={Action.Update} a={Feature.TelemetryBed}>
        <CardFooter className="flex-col items-start gap-2 border-t px-4 pt-4 text-left">
          <span className="block text-xs font-bold text-muted-foreground">{t("simulation.title")}</span>
          <div className="flex flex-wrap gap-2">
            {conditionOptions.map(({ condition, labelKey }) => {
              const isActive = activeBed.condition === condition
              const isArrest = condition === CardiacCondition.CardiacArrest
              const activeVariant = isArrest ? "destructive" : "default"
              return (
                <Button
                  key={condition}
                  variant={isActive ? activeVariant : "outline"}
                  size="sm"
                  aria-pressed={isActive}
                  onClick={() => updateSelectedBedCondition(condition)}
                  className={cn(
                    "text-[11px] font-bold",
                    isArrest && !isActive && "text-destructive hover:text-destructive"
                  )}
                >
                  {t(labelKey)}
                </Button>
              )
            })}
          </div>
        </CardFooter>
      </Can>
    </Card>
  )
}
