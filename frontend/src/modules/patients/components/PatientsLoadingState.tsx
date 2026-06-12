import { useTranslation } from "react-i18next"
import { Clock } from "lucide-react"

export const PatientsLoadingState = () => {
  const { t } = useTranslation()

  return (
    <div className="flex-1 flex items-center justify-center py-20">
      <div className="flex items-center gap-3 text-muted">
        <Clock className="w-5 h-5 animate-spin" />
        <span className="text-sm font-medium">{t("patients.loading")}</span>
      </div>
    </div>
  )
}
