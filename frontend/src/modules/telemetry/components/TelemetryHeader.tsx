import { useTranslation } from "react-i18next"
import { Activity, Volume2, VolumeX } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"

interface TelemetryHeaderProps {
  isMuted: boolean
  setIsMuted: (muted: boolean) => void
}

export const TelemetryHeader = ({ isMuted, setIsMuted }: TelemetryHeaderProps) => {
  const { t } = useTranslation()

  return (
    <div className="flex items-center justify-between flex-wrap gap-4">
      <div className="text-left">
        <h2 className="text-xl font-black text-gray-900 leading-none flex items-center gap-2">
          <Activity className="w-5 h-5 text-primary animate-pulse-glow" />
          {t("telemetry.title")}
        </h2>
        <span className="text-xs text-muted mt-1.5 block">
          {t("telemetry.subtitle")}
        </span>
      </div>

      <div className="flex items-center gap-3">
        <Button
          variantType="outline"
          onClick={() => setIsMuted(!isMuted)}
          className="px-3 gap-2 text-xs"
        >
          {isMuted ? (
            <>
              <VolumeX className="w-4 h-4 text-red-500" />
              {t("telemetry.alarmsMuted")}
            </>
          ) : (
            <>
              <Volume2 className="w-4 h-4 text-emerald-500 animate-pulse" />
              {t("telemetry.alarmsActive")}
            </>
          )}
        </Button>
      </div>
    </div>
  )
}
