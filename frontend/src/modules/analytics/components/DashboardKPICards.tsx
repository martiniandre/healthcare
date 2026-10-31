import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import {
  Activity,
  ArrowDownRight,
  ArrowUpRight,
  CalendarDays,
  Clock,
  FileText,
  Stethoscope,
  Users,
} from "lucide-react"
import { MetricStrip, type MetricStripItem } from "./MetricStrip"
import {
  resolveOccupancySeverity,
  resolveWaitTimeSeverity,
} from "../dashboard_thresholds"
import { toneClinicalStatus } from "../metric_tones"
import { useLocale } from "../../../shared/hooks/useLocale"
import { formatDecimal } from "../format_numbers"

interface DashboardKPICardsProps {
  consultationsToday: number
  consultationsTrend: string
  occupancyRate: number
  avgWaitTimeMinutes: number
  activePatients: number
  examsToday: number
  newDiagnosesToday: number
}

export const DashboardKPICards = ({
  consultationsToday,
  consultationsTrend,
  occupancyRate,
  avgWaitTimeMinutes,
  activePatients,
  examsToday,
  newDiagnosesToday,
}: DashboardKPICardsProps) => {
  const { t } = useTranslation("analytics")
  const locale = useLocale()

  const occupancySeverity = resolveOccupancySeverity(occupancyRate)
  const waitTimeSeverity = resolveWaitTimeSeverity(avgWaitTimeMinutes)

  const kpiItems = useMemo<MetricStripItem[]>(
    () => [
      {
        label: t("dashboard.kpi.consultationsToday"),
        value: formatDecimal(consultationsToday, locale),
        icon: Stethoscope,
        tone: "primary",
        meta: {
          text: consultationsTrend,
          tone: "text-gray-700",
          icon: consultationsTrend.trim().startsWith("-") ? ArrowDownRight : ArrowUpRight,
        },
      },
      {
        label: t("dashboard.kpi.occupancyRate"),
        value: formatDecimal(occupancyRate, locale, 1),
        unit: "%",
        icon: Activity,
        tone: occupancySeverity,
        meta: {
          text: t(`dashboard.occupancy.severity.${occupancySeverity}`),
          tone: "",
          status: toneClinicalStatus[occupancySeverity],
        },
      },
      {
        label: t("dashboard.kpi.avgWaitTime"),
        value: formatDecimal(avgWaitTimeMinutes, locale),
        unit: t("dashboard.unit.minutes"),
        icon: Clock,
        tone: waitTimeSeverity,
        meta: {
          text: t(`dashboard.waitTime.severity.${waitTimeSeverity}`),
          tone: "",
          status: toneClinicalStatus[waitTimeSeverity],
        },
      },
      {
        label: t("dashboard.kpi.activePatients"),
        value: formatDecimal(activePatients, locale),
        icon: Users,
        tone: "secondary",
        meta: { text: t("dashboard.kpi.inCare"), tone: "text-muted-foreground", icon: Users },
      },
      {
        label: t("dashboard.kpi.examsToday"),
        value: formatDecimal(examsToday, locale),
        icon: Activity,
        tone: "primary",
        meta: {
          text: t("dashboard.kpi.performed"),
          tone: "text-muted-foreground",
          icon: Activity,
        },
      },
      {
        label: t("dashboard.kpi.newDiagnoses"),
        value: formatDecimal(newDiagnosesToday, locale),
        icon: FileText,
        tone: "neutral",
        meta: {
          text: t("dashboard.kpi.last30Days"),
          tone: "text-muted-foreground",
          icon: CalendarDays,
        },
      },
    ],
    [
      t,
      locale,
      consultationsToday,
      consultationsTrend,
      occupancyRate,
      occupancySeverity,
      avgWaitTimeMinutes,
      waitTimeSeverity,
      activePatients,
      examsToday,
      newDiagnosesToday,
    ]
  )

  return <MetricStrip items={kpiItems} />
}
