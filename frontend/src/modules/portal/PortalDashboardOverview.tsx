import type { PortalDashboard } from "./types"
import { StatusBadge } from "../../shared/components/ui/status-badge"
import { Card, CardContent, CardHeader, CardTitle } from "../../shared/components/ui/card"
import { History, Activity, Pill, FileText } from "lucide-react"
import { formatDate } from "../../shared/utils/dates"
import { useTranslation } from "react-i18next"
import { useLocale } from "../../shared/hooks/useLocale"
import { useClinicalStatusLabels } from "../../shared/hooks/useClinicalStatusLabels"
import { findVitalSignDisplay } from "../patients/components/vitalSignDisplay"
import { PortalVitalValue } from "./components/PortalVitalValue"
import { useNumberFormat } from "./hooks/useNumberFormat"

interface PortalDashboardOverviewProps {
  dashboard: PortalDashboard
}

export const PortalDashboardOverview = ({ dashboard }: PortalDashboardOverviewProps) => {
  const { t } = useTranslation(["portal", "patients"])
  const locale = useLocale()
  const formatNumber = useNumberFormat()
  const { encounterStatus } = useClinicalStatusLabels()
  const summaryCards = [
    {
      icon: <History className="w-5 h-5 text-blue-600" />,
      label: t("portal:summary.totalEncounters"),
      value: dashboard.upcoming_encounters.length + dashboard.recent_reports.length,
      bgColor: "bg-blue-50",
    },
    {
      icon: <Activity className="w-5 h-5 text-amber-600" />,
      label: t("portal:summary.activeConditions"),
      value: dashboard.active_conditions.length,
      bgColor: "bg-amber-50",
    },
    {
      icon: <Pill className="w-5 h-5 text-green-600" />,
      label: t("portal:summary.activeMedications"),
      value: dashboard.active_medications.length,
      bgColor: "bg-green-50",
    },
    {
      icon: <FileText className="w-5 h-5 text-purple-600" />,
      label: t("portal:summary.recentReports"),
      value: dashboard.recent_reports.length,
      bgColor: "bg-purple-50",
    },
  ]

  return (
    <div className="flex flex-col gap-4 md:gap-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-lg font-bold">
            {t("portal:greeting", { firstName: dashboard.patient_info.full_name.split(" ")[0] })}
          </CardTitle>
          <p className="text-sm text-muted-foreground">{t("portal:welcome")}</p>
        </CardHeader>
      </Card>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
        {summaryCards.map((card) => (
          <Card key={card.label} className={`${card.bgColor} border-0 py-5`}>
            <CardContent className="flex items-center gap-3 px-4 md:px-5">
              {card.icon}
              <div className="min-w-0">
                <p className="text-2xl font-bold tabular-nums text-foreground">{formatNumber(card.value)}</p>
                <p className="text-xs text-muted-foreground font-medium break-words">{card.label}</p>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {dashboard.upcoming_encounters.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-bold">{t("portal:upcomingEncounters")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {dashboard.upcoming_encounters.slice(0, 5).map((encounter) => (
              <div
                key={encounter.fhir_resource_id}
                className="flex flex-wrap items-center justify-between gap-2 p-3 bg-muted rounded-lg"
              >
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-semibold text-foreground break-words">
                    {encounter.reason_display || t("portal:fallback.encounter")}
                  </p>
                  <p className="text-xs text-muted-foreground tabular-nums">
                    {formatDate(encounter.started_at, locale)}
                  </p>
                </div>
                <StatusBadge status="info">{encounterStatus(encounter.status)}</StatusBadge>
              </div>
            ))}
          </CardContent>
        </Card>
      )}

      {dashboard.recent_observations.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-bold">{t("portal:recentVitalSigns")}</CardTitle>
          </CardHeader>
          <CardContent className="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
            {dashboard.recent_observations.slice(0, 8).map((observation) => {
              const displayMetadata = findVitalSignDisplay(observation.loinc_code)
              return (
                <div key={observation.fhir_resource_id} className="p-3 bg-muted rounded-lg">
                  <p className="text-xs text-muted-foreground font-medium break-words">
                    {displayMetadata?.labelKey ? t(`patients:${displayMetadata.labelKey}`) : observation.code_display}
                  </p>
                  <PortalVitalValue
                    notPerformed={observation.not_performed}
                    valueQuantity={observation.value_quantity}
                    valueUnit={observation.value_unit}
                    valueClassName="text-lg font-bold text-foreground"
                  />
                </div>
              )
            })}
          </CardContent>
        </Card>
      )}
    </div>
  )
}
