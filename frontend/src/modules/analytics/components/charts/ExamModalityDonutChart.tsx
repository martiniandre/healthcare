import { useMemo, useState } from "react"
import { useTranslation } from "react-i18next"
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts"
import { AnalyticsPanel, PanelMetric } from "../AnalyticsPanel"
import { ChartTooltip } from "./ChartTooltip"
import { EmptyState } from "../../../../shared/components/ui/EmptyState"
import { cn } from "../../../../shared/utils/cn"
import type { ModalityData } from "../../types"

interface ExamModalityDonutChartProps {
  examModalities: ModalityData[]
}

export const ExamModalityDonutChart = ({ examModalities }: ExamModalityDonutChartProps) => {
  const { t } = useTranslation("analytics")
  const [selectedModality, setSelectedModality] = useState<string | null>(null)

  const totalStudiesCount = useMemo(
    () => examModalities.reduce((runningTotal, modality) => runningTotal + modality.count, 0),
    [examModalities]
  )

  const hasRenderableData = examModalities.length > 0 && totalStudiesCount > 0

  return (
    <AnalyticsPanel
      title={t("exams.title")}
      subtitle={t("exams.subtitle")}
      trailing={
        hasRenderableData ? (
          <PanelMetric
            label={t("exams.totalLabel")}
            value={t("exams.studiesUnit", { count: totalStudiesCount })}
          />
        ) : null
      }
    >
      {hasRenderableData ? (
        <div className="flex flex-col items-center gap-6 sm:flex-row">
          <div className="relative aspect-square w-full max-w-[200px] shrink-0">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={examModalities}
                  dataKey="count"
                  nameKey="modality"
                  innerRadius="62%"
                  outerRadius="100%"
                  paddingAngle={2}
                  stroke="var(--color-card)"
                  strokeWidth={2}
                  isAnimationActive={false}
                >
                  {examModalities.map((modality) => (
                    <Cell
                      key={modality.modality}
                      fill={modality.color}
                      opacity={
                        selectedModality === null || selectedModality === modality.modality
                          ? 1
                          : 0.25
                      }
                    />
                  ))}
                </Pie>
                <Tooltip
                  cursor={false}
                  content={({ active, payload }) => {
                    if (!active || !payload || payload.length === 0) return null
                    const modality = payload[0]?.payload as ModalityData | undefined
                    if (!modality) return null
                    return (
                      <ChartTooltip
                        label={modality.modality}
                        value={String(modality.count)}
                        unit={t("exams.examesUnit")}
                        barFill={modality.color}
                      />
                    )
                  }}
                />
              </PieChart>
            </ResponsiveContainer>

            <div className="pointer-events-none absolute inset-0 flex flex-col items-center justify-center text-center">
              <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
                {t("exams.totalLabel")}
              </span>
              <span className="text-3xl font-black leading-none tracking-tight tabular-nums text-gray-900">
                {totalStudiesCount}
              </span>
              <span className="mt-1 text-[10px] font-semibold text-muted-foreground">
                {t("exams.studiesLabel")}
              </span>
            </div>
          </div>

          <ul className="flex w-full flex-col gap-2">
            {examModalities.map((modality) => {
              const isSelected = selectedModality === modality.modality
              return (
                <li
                  key={modality.modality}
                  tabIndex={0}
                  onMouseEnter={() => setSelectedModality(modality.modality)}
                  onMouseLeave={() => setSelectedModality(null)}
                  onFocus={() => setSelectedModality(modality.modality)}
                  onBlur={() => setSelectedModality(null)}
                  className={cn(
                    "flex cursor-default items-center justify-between rounded-lg border px-2.5 py-2 transition-colors duration-200",
                    isSelected
                      ? "border-border bg-gray-50"
                      : "border-transparent bg-transparent"
                  )}
                >
                  <span className="flex min-w-0 items-center gap-2.5">
                    <span
                      className="size-3 shrink-0 rounded-full"
                      style={{ backgroundColor: modality.color }}
                      aria-hidden="true"
                    />
                    <span className="truncate text-xs font-bold text-gray-700">
                      {modality.modality}
                    </span>
                  </span>
                  <span className="flex shrink-0 items-baseline gap-2">
                    <span className="text-xs font-black tabular-nums text-gray-900">
                      {modality.count}
                    </span>
                    <span className="w-9 text-right text-[10px] font-bold tabular-nums text-muted-foreground">
                      {modality.percentage}%
                    </span>
                  </span>
                </li>
              )
            })}
          </ul>
        </div>
      ) : (
        <EmptyState title={t("empty.exams")} description={t("empty.examsDesc")} />
      )}
    </AnalyticsPanel>
  )
}
