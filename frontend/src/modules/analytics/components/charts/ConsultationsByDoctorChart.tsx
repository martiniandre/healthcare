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
  type YAxisTickContentProps,
} from "recharts"
import { EmptyState } from "../../../../shared/components/ui/EmptyState"
import { AnalyticsPanel, PanelMetric } from "../AnalyticsPanel"
import { ChartTooltip } from "./ChartTooltip"
import type { DoctorConsultation } from "../../dashboard_types"

const BAND_HEIGHT = 44
const CHART_PADDING = 48
const MINIMUM_CHART_HEIGHT = 300

interface DoctorChartDatum {
  doctorName: string
  specialty: string
  count: number
}

interface ConsultationsByDoctorChartProps {
  consultationsPerDoctor: DoctorConsultation[]
}

export const ConsultationsByDoctorChart = ({
  consultationsPerDoctor,
}: ConsultationsByDoctorChartProps) => {
  const { t } = useTranslation("analytics")

  const chartData = useMemo<DoctorChartDatum[]>(
    () =>
      [...consultationsPerDoctor]
        .map((consultation) => ({
          doctorName: consultation.doctor_name,
          specialty: consultation.specialty,
          count: consultation.count,
        }))
        .sort((firstDoctor, secondDoctor) => secondDoctor.count - firstDoctor.count),
    [consultationsPerDoctor]
  )

  const peakCount = chartData.length > 0 ? chartData[0].count : 0
  const minimumChartHeight = Math.max(chartData.length * BAND_HEIGHT + CHART_PADDING, MINIMUM_CHART_HEIGHT)

  const renderDoctorTick = ({ x = 0, y = 0, payload }: YAxisTickContentProps) => {
    const datum = chartData[payload?.index ?? -1]

    return (
      <g transform={`translate(${x}, ${y})`}>
        <text textAnchor="end" dy={-2} className="fill-gray-900 text-[11px] font-bold">
          {datum?.doctorName}
        </text>
        <text textAnchor="end" dy={12} className="fill-gray-500 text-[10px] font-semibold">
          {datum?.specialty}
        </text>
      </g>
    )
  }

  return (
    <AnalyticsPanel
      title={t("dashboard.doctorConsultations.title")}
      subtitle={t("dashboard.doctorConsultations.subtitle")}
      trailing={
        chartData.length > 0 ? (
          <PanelMetric
            label={t("dashboard.doctorConsultations.peakLabel")}
            value={t("dashboard.doctorConsultations.peakUnit", { count: peakCount })}
            tone="primary"
          />
        ) : null
      }
      className="h-full"
    >
      {chartData.length > 0 ? (
        <div className="flex-1 min-h-0" style={{ minHeight: `${minimumChartHeight}px` }}>
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={chartData}
              layout="vertical"
              margin={{ top: 4, right: 52, bottom: 4, left: 0 }}
              barCategoryGap={12}
            >
              <CartesianGrid
                horizontal={true}
                vertical={false}
                strokeDasharray="3 4"
                className="stroke-border"
              />
              <XAxis
                type="number"
                domain={[0, peakCount]}
                tick={false}
                axisLine={false}
                tickLine={false}
              />
              <YAxis
                type="category"
                dataKey="doctorName"
                width={150}
                tick={renderDoctorTick}
                axisLine={false}
                tickLine={false}
              />
              <Tooltip
                cursor={{ className: "fill-gray-100" }}
                content={({ active, payload }) => {
                  if (!active || !payload || payload.length === 0) return null
                  const datum = payload[0].payload as DoctorChartDatum
                  return (
                    <ChartTooltip
                      label={datum.doctorName}
                      detail={datum.specialty}
                      value={String(datum.count)}
                      unit={t("dashboard.doctorConsultations.consultationsUnit")}
                      barFill="var(--color-primary)"
                    />
                  )
                }}
              />
              <Bar
                dataKey="count"
                fill="currentColor"
                className="text-primary"
                radius={[4, 6, 6, 4]}
                isAnimationActive={false}
              >
                <LabelList
                  dataKey="count"
                  position="right"
                  offset={10}
                  className="fill-gray-900 text-[11px] font-black tabular-nums"
                />
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      ) : (
        <EmptyState
          title={t("dashboard.doctorConsultations.noDataTitle")}
          description={t("dashboard.doctorConsultations.noDataDescription")}
          className="flex-1"
        />
      )}
    </AnalyticsPanel>
  )
}
