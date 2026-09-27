import { useMemo } from "react"
import type { LucideIcon } from "lucide-react"
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
import { Card } from "../../../shared/components/ui/Card"
import { cn } from "../../../shared/utils/cn"
import {
  resolveOccupancySeverity,
  resolveWaitTimeSeverity,
  toneClassNames,
  type MetricTone,
} from "../dashboard_thresholds"

interface KpiMeta {
  text: string
  tone: string
  icon?: LucideIcon
}

interface KpiDefinition {
  label: string
  value: string
  unit?: string
  icon: LucideIcon
  tone: MetricTone
  meta: KpiMeta
}

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

  const kpiDefinitions = useMemo<KpiDefinition[]>(
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

  return (
    <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6 lg:gap-4">
      {kpiDefinitions.map((kpiDefinition) => (
        <KpiCard key={kpiDefinition.label} definition={kpiDefinition} />
      ))}
    </div>
  )
}

interface KpiCardProps {
  definition: KpiDefinition
}

const KpiCard = ({ definition }: KpiCardProps) => {
  const { label, value, unit, icon: Icon, tone, meta } = definition
  const toneClassName = toneClassNames[tone]
  const MetaIcon = meta.icon

  return (
    <Card className="relative flex flex-col gap-3 overflow-hidden p-4">
      <span
        className={cn("absolute inset-x-0 bottom-0 h-[3px]", toneClassName.bar)}
        aria-hidden="true"
      />

      <div className="flex items-start justify-between gap-2">
        <span className="text-[11px] font-bold uppercase leading-tight tracking-wider text-muted-foreground">
          {label}
        </span>
        <span
          className={cn("grid size-8 shrink-0 place-items-center rounded-lg", toneClassName.surface)}
        >
          <Icon className={cn("size-4", toneClassName.text)} aria-hidden="true" />
        </span>
      </div>

      <div className="flex items-baseline gap-1.5">
        <span className="text-3xl font-black leading-none tracking-tight tabular-nums text-gray-900">
          {value}
        </span>
        {unit ? (
          <span className="text-xs font-bold text-muted-foreground">{unit}</span>
        ) : null}
      </div>

      <div className="mt-auto flex h-4 items-center">
        <span className={cn("flex items-center gap-1 text-[11px] font-bold", meta.tone)}>
          {MetaIcon ? <MetaIcon className="size-3 shrink-0" aria-hidden="true" /> : null}
          <span className="truncate">{meta.text}</span>
        </span>
      </div>
    </Card>
  )
}
