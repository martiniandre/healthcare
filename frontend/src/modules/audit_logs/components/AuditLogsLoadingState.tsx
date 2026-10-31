import { useTranslation } from "react-i18next"
import { Spinner } from "../../../shared/components/ui/spinner"

export const AuditLogsLoadingState = () => {
  const { t } = useTranslation("auditLogs")

  return (
    <div className="flex flex-col items-center justify-center py-16 gap-3">
      <Spinner className="size-8 text-muted-foreground" />
      <span className="text-sm font-medium text-muted-foreground">{t("loading")}</span>
    </div>
  )
}
