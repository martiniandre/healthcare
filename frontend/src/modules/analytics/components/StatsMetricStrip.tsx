import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { Activity, CalendarDays, CheckSquare, Clock, Users } from "lucide-react"
import { MetricStrip, type MetricStripItem } from "./MetricStrip"
import { resolveFhirComplianceSeverity } from "../stats_domain"
import { toneClassNames } from "../metric_tones"

interface StatsMetricStripProps {
  totalRegisteredPatients: number
  fhirComplianceRate: number
  averageServiceDurationMinutes: number
  weeklyConsultationsTotal: number
  weeklyPeakDayLabel: string | null
}

export const StatsMetricStrip = ({
  totalRegisteredPatients,
  fhirComplianceRate,
  averageServiceDurationMinutes,
  weeklyConsultationsTotal,
  weeklyPeakDayLabel,
}: StatsMetricStripProps) => {
  const { t } = useTranslation("analytics")

  const complianceSeverity = resolveFhirComplianceSeverity(fhirComplianceRate)

  const metricItems = useMemo<MetricStripItem[]>(
    () => [
      {
        label: t("metrics.activePatients"),
        value: String(totalRegisteredPatients),
        icon: Users,
        tone: "primary",
        meta: {
          text: t("metrics.registeredBase"),
          tone: "text-muted-foreground",
          icon: Users,
        },
      },
      {
        label: t("metrics.fhirCompliance"),
        value: fhirComplianceRate.toFixed(1),
        unit: "%",
        icon: CheckSquare,
        tone: complianceSeverity,
        meta: {
          text: t(`metrics.compliance.severity.${complianceSeverity}`),
          tone: toneClassNames[complianceSeverity].text,
        },
      },
      {
        label: t("metrics.avgConsultationTime"),
        value: averageServiceDurationMinutes.toFixed(1),
        unit: t("dashboard.unit.minutes"),
        icon: Clock,
        tone: "secondary",
        meta: {
          text: t("metrics.averageAcrossVisits"),
          tone: "text-muted-foreground",
          icon: Clock,
        },
      },
      {
        label: t("metrics.weeklyConsultations"),
        value: String(weeklyConsultationsTotal),
        icon: Activity,
        tone: "neutral",
        meta: {
          text:
            weeklyPeakDayLabel === null
              ? t("metrics.lastSevenDays")
              : t("metrics.peakOnDay", { day: weeklyPeakDayLabel }),
          tone: "text-muted-foreground",
          icon: CalendarDays,
        },
      },
    ],
    [
      t,
      totalRegisteredPatients,
      fhirComplianceRate,
      complianceSeverity,
      averageServiceDurationMinutes,
      weeklyConsultationsTotal,
      weeklyPeakDayLabel,
    ]
  )

  return <MetricStrip items={metricItems} columnsClassName="grid-cols-1 sm:grid-cols-2 xl:grid-cols-4" />
}
