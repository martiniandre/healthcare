import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { Bell } from "lucide-react"
import { BedStatus } from "../../../shared/types"

interface Bed {
  id: string
  bedNumber: string
  patientName: string
  status: BedStatus
  condition: string
  bpm: number
  spo2: number
}

interface TelemetryBedListProps {
  beds: Bed[]
  selectedBedId: string | null
  setSelectedBedId: (id: string | null) => void
}

export const TelemetryBedList = ({ beds, selectedBedId, setSelectedBedId }: TelemetryBedListProps) => {
  const { t } = useTranslation()

  return (
    <div className="flex flex-col gap-4 lg:col-span-1 text-left">
      <Card className="p-4 flex flex-col gap-4">
        <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2 border-b border-border pb-3">
          <Bell className="w-4 h-4 text-primary animate-pulse-glow" />
          {t("telemetry.availableBeds", { count: beds.length })}
        </h3>

        <div className="flex flex-col gap-3">
          {beds.map((bedItem) => {
            const isSelected = bedItem.id === selectedBedId
            return (
              <div
                key={bedItem.id}
                onClick={() => setSelectedBedId(bedItem.id)}
                className={`cursor-pointer border p-3.5 rounded-xl transition-all duration-200 ${
                  isSelected 
                    ? "bg-primary/5 border-primary" 
                    : bedItem.status === BedStatus.Danger
                      ? "bg-red-50 border-red-200 hover:border-red-300"
                      : "bg-white border-border hover:border-gray-300"
                }`}
              >
                <div className="flex items-center justify-between">
                  <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider">
                    {bedItem.bedNumber}
                  </span>
                  <span className={`text-[10px] px-2 py-0.5 rounded font-black uppercase ${
                    bedItem.status === BedStatus.Danger 
                      ? "bg-red-100 text-red-600" 
                      : bedItem.status === BedStatus.Warning
                        ? "bg-amber-100 text-amber-600"
                        : "bg-emerald-100 text-emerald-600"
                  }`}>
                    {bedItem.condition}
                  </span>
                </div>

                <h4 className="text-sm font-bold text-gray-800 mt-1">
                  {bedItem.patientName}
                </h4>

                <div className="flex items-center justify-between mt-3 text-xs">
                  <span className="text-gray-500">FC: <strong className="text-gray-800">{bedItem.bpm > 0 ? `${bedItem.bpm} BPM` : "---"}</strong></span>
                  <span className="text-gray-500">SpO₂: <strong className="text-gray-800">{bedItem.spo2 > 0 ? `${bedItem.spo2}%` : "---"}</strong></span>
                </div>
              </div>
            )
          })}
        </div>
      </Card>
    </div>
  )
}
