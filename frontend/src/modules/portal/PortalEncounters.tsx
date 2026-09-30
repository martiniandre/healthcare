import { useTranslation } from "react-i18next"
import { usePortalEncountersQuery } from "./queries"
import { Card } from "../../shared/components/ui/Card"
import { Loader2 } from "lucide-react"
import { formatDate } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { useClinicalStatusLabels } from "../../shared/hooks/useClinicalStatusLabels"

export const PortalEncounters = () => {
  const { t } = useTranslation("portal")
  const { data: encounters, isLoading } = usePortalEncountersQuery()
  const locale = useLocale()
  const { encounterStatus } = useClinicalStatusLabels()

  if (isLoading) {
    return (
      <Card className="flex items-center justify-center min-h-[300px]">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
      </Card>
    )
  }

  if (!encounters || encounters.length === 0) {
    return (
      <Card className="py-16 text-center">
        <p className="text-sm text-gray-500">{t("emptyEncounters")}</p>
      </Card>
    )
  }

  return (
    <div className="bg-card border border-border rounded-xl overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-gray-50 border-b border-border">
              <th className="text-left p-4 text-xs font-bold text-gray-500 uppercase tracking-wider">{t("table.date")}</th>
              <th className="text-left p-4 text-xs font-bold text-gray-500 uppercase tracking-wider">{t("table.reason")}</th>
              <th className="text-left p-4 text-xs font-bold text-gray-500 uppercase tracking-wider">{t("table.status")}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {encounters.map((encounter) => (
              <tr key={encounter.fhir_resource_id} className="hover:bg-gray-50">
                <td className="p-4 text-gray-900 font-medium whitespace-nowrap">
                  {formatDate(encounter.started_at, locale)}
                </td>
                <td className="p-4 text-gray-700">{encounter.reason_display || t("fallback.encounter")}</td>
                <td className="p-4">
                  <span className="text-xs font-bold px-2 py-1 rounded-full bg-blue-100 text-blue-700">
                    {encounterStatus(encounter.status)}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
