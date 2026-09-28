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
import { toneClassNames } from "../metric_tones"

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

  const occupancySeverity = resolveOccupancySeverity(occupancyRate)
  const waitTimeSeverity = resolveWaitTimeSeverity(avgWaitTimeMinutes)

  const kpiItems = useMemo<MetricStripItem[]>(
    () => [
      {
        label: t("dashboard.kpi.consultationsToday"),
        value: String(consultationsToday),
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
        value: occupancyRate.toFixed(1),
        unit: "%",
        icon: Activity,
        tone: occupancySeverity,
        meta: {
          text: t(`dashboard.occupancy.severity.${occupancySeverity}`),
          tone: toneClassNames[occupancySeverity].text,
        },
      },
      {
        label: t("dashboard.kpi.avgWaitTime"),
        value: avgWaitTimeMinutes.toFixed(0),
        unit: t("dashboard.unit.minutes"),
        icon: Clock,
        tone: waitTimeSeverity,
        meta: {
          text: t(`dashboard.waitTime.severity.${waitTimeSeverity}`),
          tone: toneClassNames[waitTimeSeverity].text,
        },
      },
      {
        label: t("dashboard.kpi.activePatients"),
        value: String(activePatients),
        icon: Users,
        tone: "secondary",
        meta: { text: t("dashboard.kpi.inCare"), tone: "text-muted-foreground", icon: Users },
      },
      {
        label: t("dashboard.kpi.examsToday"),
        value: String(examsToday),
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
        value: String(newDiagnosesToday),
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
