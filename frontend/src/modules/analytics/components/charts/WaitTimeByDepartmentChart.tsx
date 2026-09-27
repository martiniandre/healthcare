import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import {
  Bar,
  BarChart,
  CartesianGrid,
  LabelList,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
  type RenderableText,
  type YAxisTickContentProps,
} from "recharts"
import { EmptyState } from "../../../../shared/components/ui/EmptyState"
import { DashboardPanel, PanelMetric } from "../DashboardPanel"
import { ChartTooltip } from "./ChartTooltip"
import { SeverityLegend, type SeverityBand } from "../SeverityLegend"
import { resolveWaitTimeSeverity, severityFillVariables, type Severity } from "../../dashboard_thresholds"
import type { DepartmentWaitTime } from "../../dashboard_types"

const BAND_HEIGHT = 40
const CHART_PADDING = 44
const MINIMUM_CHART_HEIGHT = 260

interface WaitTimeChartDatum {
  department: string
  minutes: number
  severity: Severity
  stable: number | null
  elevated: number | null
  critical: number | null
}

interface WaitTimeByDepartmentChartProps {
  waitTimeByDepartment: DepartmentWaitTime[]
}

export const WaitTimeByDepartmentChart = ({
  waitTimeByDepartment,
}: WaitTimeByDepartmentChartProps) => {
  const { t } = useTranslation("analytics")

  const chartData = useMemo<WaitTimeChartDatum[]>(
    () =>
      [...waitTimeByDepartment]
        .map((departmentWaitTime) => {
          const severity = resolveWaitTimeSeverity(departmentWaitTime.minutes)
          return {
            department: departmentWaitTime.department,
            minutes: departmentWaitTime.minutes,
            severity,
            stable: severity === "stable" ? departmentWaitTime.minutes : null,
            elevated: severity === "elevated" ? departmentWaitTime.minutes : null,
            critical: severity === "critical" ? departmentWaitTime.minutes : null,
          }
        })
        .sort(
          (firstDepartment, secondDepartment) => secondDepartment.minutes - firstDepartment.minutes
        ),
    [waitTimeByDepartment]
  )

  const peakMinutes = chartData.length > 0 ? chartData[0].minutes : 0
  const minimumChartHeight = Math.max(
    chartData.length * BAND_HEIGHT + CHART_PADDING,
    MINIMUM_CHART_HEIGHT
  )

  const minutesUnit = t("dashboard.unit.minutes")

  const renderDepartmentTick = ({ x = 0, y = 0, payload }: YAxisTickContentProps) => {
    const datum = chartData[payload?.index ?? -1]

    return (
      <g transform={`translate(${x}, ${y})`}>
        <text textAnchor="end" dy={4} className="fill-gray-900 text-[11px] font-bold">
          {datum?.department}
        </text>
      </g>
    )
  }

  const minutesLabelList = (
    <LabelList
      dataKey="minutes"
      position="right"
      offset={10}
      formatter={(labelValue: RenderableText) =>
        `${Math.round(Number(labelValue))} ${minutesUnit}`
      }
      className="fill-gray-900 text-[11px] font-black tabular-nums"
    />
  )

  const severityBands: SeverityBand[] = [
    {
      severity: "critical",
      label: t("dashboard.waitTime.severity.critical"),
      range: t("dashboard.waitTime.range.critical"),
    },
    {
      severity: "elevated",
      label: t("dashboard.waitTime.severity.elevated"),
      range: t("dashboard.waitTime.range.elevated"),
    },
    {
      severity: "stable",
      label: t("dashboard.waitTime.severity.stable"),
      range: t("dashboard.waitTime.range.stable"),
    },
  ]

  return (
    <DashboardPanel
      title={t("dashboard.waitTime.title")}
      subtitle={t("dashboard.waitTime.subtitle")}
      trailing={
        chartData.length > 0 ? (
          <PanelMetric
            label={t("dashboard.waitTime.peakLabel")}
            value={t("dashboard.waitTime.peakUnit", { minutes: Math.round(peakMinutes) })}
          />
        ) : null
      }
    >
      {chartData.length > 0 ? (
        <div className="flex flex-col gap-4">
          <div className="min-h-0" style={{ minHeight: `${minimumChartHeight}px` }}>
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                data={chartData}
                layout="vertical"
                margin={{ top: 4, right: 76, bottom: 4, left: 0 }}
                barCategoryGap={14}
              >
                <CartesianGrid
                  horizontal={true}
                  vertical={false}
                  strokeDasharray="3 4"
                  className="stroke-border"
                />
                <XAxis
                  type="number"
                  domain={[0, peakMinutes]}
                  tick={false}
                  axisLine={false}
                  tickLine={false}
                />
                <YAxis
                  type="category"
                  dataKey="department"
                  width={172}
                  tick={renderDepartmentTick}
                  axisLine={false}
                  tickLine={false}
                />
                <Tooltip
                  cursor={{ className: "fill-gray-100" }}
                  content={({ active, payload }) => {
                    if (!active || !payload || payload.length === 0) return null
                    const datum = payload[0].payload as WaitTimeChartDatum
                    return (
                      <ChartTooltip
                        label={datum.department}
                        value={String(Math.round(datum.minutes))}
                        unit={minutesUnit}
                        barFill={severityFillVariables[datum.severity]}
                      />
                    )
                  }}
                />
                <Bar
                  dataKey="stable"
                  stackId="waitTime"
                  fill={severityFillVariables.stable}
                  radius={[4, 0, 0, 4]}
                  isAnimationActive={false}
                >
                  {minutesLabelList}
                </Bar>
                <Bar
                  dataKey="elevated"
                  stackId="waitTime"
                  fill={severityFillVariables.elevated}
                  radius={[4, 0, 0, 4]}
                  isAnimationActive={false}
                >
                  {minutesLabelList}
                </Bar>
                <Bar
                  dataKey="critical"
                  stackId="waitTime"
                  fill={severityFillVariables.critical}
                  radius={[4, 6, 6, 4]}
                  isAnimationActive={false}
                >
                  {minutesLabelList}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
          <SeverityLegend bands={severityBands} className="mt-auto" />
        </div>
      ) : (
        <EmptyState
          title={t("dashboard.waitTime.noDataTitle")}
          description={t("dashboard.waitTime.noDataDescription")}
        />
      )}
    </DashboardPanel>
  )
}
