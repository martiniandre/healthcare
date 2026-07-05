import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"
import { Heart, Activity, Thermometer, AlertTriangle } from "lucide-react"
import { BedStatus, CardiacCondition } from "../../../shared/types"
import { type MutableRefObject } from "react"

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

export const TelemetryBedMonitor = ({
  activeBed,
  activeRoomName,
  canvasRef,
  updateSelectedBedCondition
}: TelemetryBedMonitorProps) => {
  const { t } = useTranslation()

  if (!activeBed) {
    return (
      <Card className="flex-1 p-8 border border-border bg-gray-50/50 flex flex-col items-center justify-center text-center gap-4 min-h-[400px]">
        <div className="bg-primary/5 p-4 rounded-full border border-primary/10 text-primary">
          <Activity className="w-8 h-8 animate-pulse" />
        </div>
        <div className="max-w-md flex flex-col gap-1">
          <h3 className="text-md font-extrabold text-gray-900">
            {t("telemetry.noBedSelected")}
          </h3>
          <p className="text-xs text-gray-500 leading-normal">
            {t("telemetry.noBedSelectedDesc")}
          </p>
        </div>
      </Card>
    )
  }

  return (
    <Card className="p-4 flex flex-col gap-4">
      <div className="flex items-center justify-between border-b border-border pb-3 flex-wrap gap-2">
        <div className="text-left">
          <span className="text-xs font-bold text-primary uppercase tracking-wider">{activeRoomName}</span>
          <h3 className="text-md font-bold text-gray-900">{activeBed.bedNumber} • {activeBed.patientName}</h3>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded font-bold">
            {activeBed.age} {t("telemetry.years")} · {activeBed.gender}
          </span>
          {activeBed.status !== BedStatus.Normal && (
            <span className="text-xs bg-red-50 text-red-600 px-2.5 py-1 rounded-full font-bold flex items-center gap-1 border border-red-200">
              <AlertTriangle className="w-3.5 h-3.5" />
              {t("telemetry.criticalBadge")}
            </span>
          )}
        </div>
      </div>

      <div className="relative border border-border rounded-xl overflow-hidden bg-slate-950 p-1">
        <canvas
          ref={canvasRef}
          width={700}
          height={260}
          className="block w-full max-w-full rounded-lg"
        />
        <div className="absolute top-4 left-4 bg-black/60 backdrop-blur-md px-3 py-1.5 rounded text-[10px] text-emerald-400 font-mono border border-emerald-500/20">
          {t("telemetry.ecgLead")}
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-gray-50 border border-border p-4 rounded-xl flex items-center justify-between">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{t("telemetry.sensors.ecg")}</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">
              {activeBed.bpm > 0 ? `${activeBed.bpm} ` : "--- "}
              <span className="text-xs font-normal text-gray-400">BPM</span>
            </span>
          </div>
          <Heart className={`w-8 h-8 text-red-500 shrink-0 ${activeBed.bpm > 110 ? "animate-pulse" : ""}`} />
        </div>

        <div className="bg-gray-50 border border-border p-4 rounded-xl flex items-center justify-between">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{t("telemetry.sensors.spo2")}</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">
              {activeBed.spo2 > 0 ? `${activeBed.spo2}%` : "---%"}
            </span>
          </div>
          <Activity className="w-8 h-8 text-emerald-500 shrink-0" />
        </div>

        <div className="bg-gray-50 border border-border p-4 rounded-xl flex items-center justify-between">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{t("telemetry.sensors.temp")}</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">
              {activeBed.temperature.toFixed(1)}°C
            </span>
          </div>
          <Thermometer className="w-8 h-8 text-sky-500 shrink-0" />
        </div>
      </div>

      <Can I={Action.Update} a={Feature.TelemetryBed}>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 border-t border-border pt-4 text-left">
          <div className="flex flex-col gap-2">
            <span className="text-xs font-bold text-gray-600 block">{t("telemetry.simulation.title")}</span>
            <div className="flex gap-2 flex-wrap">
              <Button 
                variantType={activeBed.condition === CardiacCondition.Normal ? "primary" : "outline"} 
                onClick={() => updateSelectedBedCondition(CardiacCondition.Normal)}
                className="px-3 py-2 text-[11px] font-bold"
              >
                {t("telemetry.simulation.normal")}
              </Button>
              <Button 
                variantType={activeBed.condition === CardiacCondition.Bradycardia ? "primary" : "outline"} 
                onClick={() => updateSelectedBedCondition(CardiacCondition.Bradycardia)}
                className="px-3 py-2 text-[11px] font-bold"
              >
                {t("telemetry.simulation.bradycardia")}
              </Button>
            </div>
          </div>
          <div className="flex flex-col gap-2 justify-end">
            <div className="flex gap-2 flex-wrap">
              <Button 
                variantType={activeBed.condition === CardiacCondition.Tachycardia ? "primary" : "outline"} 
                onClick={() => updateSelectedBedCondition(CardiacCondition.Tachycardia)}
                className="px-3 py-2 text-[11px] font-bold"
              >
                {t("telemetry.simulation.tachycardia")}
              </Button>
              <Button 
                variantType={activeBed.condition === CardiacCondition.CardiacArrest ? "danger" : "outline"} 
                onClick={() => updateSelectedBedCondition(CardiacCondition.CardiacArrest)}
                className="px-3 py-2 text-[11px] font-bold text-red-500 hover:text-white"
              >
                {t("telemetry.simulation.cardiacArrest")}
              </Button>
            </div>
          </div>
        </div>
      </Can>
    </Card>
  )
}
