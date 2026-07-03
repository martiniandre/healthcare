import { useTranslation } from "react-i18next"
import { Clock } from "lucide-react"

export const AuditLogsLoadingState = () => {
  const { t } = useTranslation("auditLogs")

  return (
    <div className="flex flex-col items-center justify-center py-16 gap-3">
      <Clock className="w-8 h-8 text-muted animate-spin" />
      <span className="text-sm font-medium text-muted">{t("loading")}</span>
    </div>
  )
}
