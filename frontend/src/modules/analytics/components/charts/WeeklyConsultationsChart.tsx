import { useMemo, useState } from "react"
import { useTranslation } from "react-i18next"
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts"
import { AnalyticsPanel, PanelMetric } from "../AnalyticsPanel"
import { ChartTooltip } from "./ChartTooltip"
import { EmptyState } from "../../../../shared/components/ui/EmptyState"
import { toneFillVariables } from "../../metric_tones"
import { resolveDayTranslationKey, summarizeWeeklyConsultations } from "../../stats_domain"
import type { ConsultationsDayData } from "../../types"

interface WeeklyConsultationsChartProps {
  weeklyConsultations: ConsultationsDayData[]
}

interface DayTranslationInput {
  dayKey: string
  count: number
}

export const WeeklyConsultationsChart = ({
  weeklyConsultations,
}: WeeklyConsultationsChartProps) => {
  const { t } = useTranslation("analytics")
  const [activeDayKey, setActiveDayKey] = useState<string | null>(null)

  const days = useMemo<DayTranslationInput[]>(
    () =>
      weeklyConsultations.map((consultation) => ({
        dayKey: resolveDayTranslationKey(consultation.dayName),
        count: consultation.count,
      })),
    [weeklyConsultations]
  )

  const summary = useMemo(
    () => summarizeWeeklyConsultations(weeklyConsultations),
    [weeklyConsultations]
  )

  const peakCount = summary.peak

  const peakValue =
    summary.peakDayKey === null
      ? String(peakCount)
      : t("consultations.peakUnit", { count: peakCount, day: t(summary.peakDayKey) })

  return (
    <AnalyticsPanel
      title={t("consultations.title")}
      subtitle={t("consultations.subtitle")}
      trailing={
        days.length > 0 ? (
          <div className="flex flex-wrap items-center justify-end gap-1.5">
            <PanelMetric label={t("consultations.minLabel")} value={String(summary.min)} />
            <PanelMetric
              label={t("consultations.avgLabel")}
              value={String(summary.average)}
              tone="primary"
            />
            <PanelMetric label={t("consultations.peakLabel")} value={peakValue} tone="secondary" />
          </div>
        ) : null
      }
    >
      {days.length > 0 ? (
        <div className="h-[260px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={days}
              margin={{ top: 8, right: 8, bottom: 0, left: -18 }}
              onMouseLeave={() => setActiveDayKey(null)}
            >
              <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" vertical={false} />
              <XAxis
                dataKey="dayKey"
                tickFormatter={(dayKey: string) => t(dayKey)}
                axisLine={false}
                tickLine={false}
                tick={{ fontSize: 11, fontWeight: 700, fill: "var(--color-muted-foreground)" }}
                dy={4}
              />
              <YAxis
                allowDecimals={false}
                axisLine={false}
                tickLine={false}
                width={44}
                domain={[0, (maximum: number) => Math.max(Math.ceil(maximum * 1.15), 1)]}
                tick={{ fontSize: 10, fontWeight: 600, fill: "var(--color-muted-foreground)" }}
              />
              <Tooltip
                cursor={{ fill: "var(--color-gray-100)" }}
                content={({ active, payload }) => {
                  if (!active || !payload || payload.length === 0) return null
                  const day = payload[0]?.payload as DayTranslationInput | undefined
                  if (!day) return null
                  return (
                    <ChartTooltip
                      label={t(day.dayKey)}
                      value={String(day.count)}
                      unit={t("consultations.consultationsLabel")}
                      barFill={toneFillVariables.primary}
                    />
                  )
                }}
              />
              <Bar
                dataKey="count"
                radius={[6, 6, 0, 0]}
                maxBarSize={44}
                isAnimationActive={false}
              >
                {days.map((day) => (
                  <Cell
                    key={day.dayKey}
                    fill={toneFillVariables.primary}
                    fillOpacity={
                      activeDayKey === null
                        ? day.count === peakCount
                          ? 1
                          : 0.22
                        : day.dayKey === activeDayKey
                          ? 1
                          : 0.12
                    }
                    onMouseEnter={() => setActiveDayKey(day.dayKey)}
                  />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      ) : (
        <EmptyState
          title={t("empty.consultations")}
          description={t("empty.consultationsDesc")}
          className="flex-1"
        />
      )}
    </AnalyticsPanel>
  )
}
