import { useTranslation } from "react-i18next"
import { Volume2, VolumeX } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import { PageHeader } from "@/shared/components/ui/PageContainer"

interface TelemetryHeaderProps {
  isMuted: boolean
  setIsMuted: (muted: boolean) => void
}

export const TelemetryHeader = ({ isMuted, setIsMuted }: TelemetryHeaderProps) => {
  const { t } = useTranslation("telemetry")

  return (
    <PageHeader
      title={t("title")}
      actions={
        <Button variant="outline" onClick={() => setIsMuted(!isMuted)} aria-pressed={isMuted}>
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
      }
    />
  )
}
