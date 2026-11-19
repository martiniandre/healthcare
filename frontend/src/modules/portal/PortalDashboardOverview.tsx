import type { PortalDashboard } from "./types"
import { StatusBadge } from "../../shared/components/ui/status-badge"
import { Card, CardContent, CardHeader, CardTitle } from "../../shared/components/ui/card"
import { formatDate } from "../../shared/utils/dates"
import { useTranslation } from "react-i18next"
import { Link } from "react-router-dom"
import { useLocale } from "../../shared/hooks/useLocale"
import { useClinicalStatusLabels } from "../../shared/hooks/useClinicalStatusLabels"
import { findVitalSignDisplay } from "@/shared/clinical/vitalSignDisplay"
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
      label: t("portal:summary.upcomingEncounters"),
      value: dashboard.upcoming_encounters.length,
      target: "/portal?tab=appointments",
    },
    {
      label: t("portal:summary.activeConditions"),
      value: dashboard.active_conditions.length,
      target: "/portal?tab=conditions",
    },
    {
      label: t("portal:summary.activeMedications"),
      value: dashboard.active_medications.length,
      target: "/portal?tab=medications",
    },
    {
      label: t("portal:summary.recentReports"),
      value: dashboard.recent_reports.length,
      target: "/portal?tab=reports",
    },
  ]

  return (
    <div className="flex flex-col gap-4 md:gap-6">
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
        {summaryCards.map((card) => (
          <Link
            key={card.label}
            to={card.target}
            className="rounded-lg outline-none transition-shadow hover:shadow-md focus-visible:ring-[3px] focus-visible:ring-ring/50"
          >
          <Card className="h-full py-5">
            <CardContent className="flex items-start justify-between gap-3 px-4 md:px-5">
              <div className="min-w-0">
                <p className="break-words text-xs font-medium text-muted-foreground">{card.label}</p>
                <p className="mt-2 text-2xl font-semibold tabular-nums text-foreground">{formatNumber(card.value)}</p>
              </div>
            </CardContent>
          </Card>
          </Link>
        ))}
      </div>

      {dashboard.upcoming_encounters.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-semibold">{t("portal:upcomingEncounters")}</CardTitle>
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
            <CardTitle className="text-sm font-semibold">{t("portal:recentVitalSigns")}</CardTitle>
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
                    valueClassName="text-lg font-semibold text-foreground"
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
