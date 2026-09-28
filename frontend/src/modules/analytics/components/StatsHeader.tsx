import { useTranslation } from "react-i18next"
import { BarChart3 } from "lucide-react"

export const StatsHeader = () => {
  const { t } = useTranslation("analytics")

  return (
    <div className="text-left">
      <h2 className="text-xl font-display font-bold text-gray-900 leading-none flex items-center gap-2">
        <BarChart3 className="w-5 h-5 text-primary" aria-hidden="true" />
        {t("title")}
      </h2>
      <span className="text-xs text-muted mt-1.5 block">
        {t("subtitle")}
      </span>
    </div>
  )
}
