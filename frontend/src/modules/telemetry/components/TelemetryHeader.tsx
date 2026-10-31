import { useTranslation } from "react-i18next"
import { Activity, Volume2, VolumeX } from "lucide-react"
import { Button } from "@/shared/components/ui/button"

interface TelemetryHeaderProps {
  isMuted: boolean
  setIsMuted: (muted: boolean) => void
}

export const TelemetryHeader = ({ isMuted, setIsMuted }: TelemetryHeaderProps) => {
  const { t } = useTranslation("telemetry")

  return (
    <div className="flex items-center justify-between flex-wrap gap-4">
      <div className="text-left">
        <h2 className="text-xl font-display font-bold text-foreground leading-none flex items-center gap-2">
          <Activity className="size-5 text-primary" aria-hidden="true" />
          {t("title")}
        </h2>
        <span className="text-xs text-muted-foreground mt-1.5 block">
          {t("subtitle")}
        </span>
      </div>

      <div className="flex items-center gap-3">
        <Button
          variant="outline"
          onClick={() => setIsMuted(!isMuted)}
          aria-pressed={isMuted}
          className="text-xs"
        >
          {isMuted ? (
            <>
              <VolumeX className="text-destructive" />
              {t("alarmsMuted")}
            </>
          ) : (
            <>
              <Volume2 className="text-success" />
              {t("alarmsActive")}
            </>
          )}
        </Button>
      </div>
    </div>
  )
}
