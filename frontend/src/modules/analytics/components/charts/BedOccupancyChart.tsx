import { useTranslation } from "react-i18next"
import { PolarAngleAxis, RadialBar, RadialBarChart, ResponsiveContainer } from "recharts"
import { AnalyticsPanel, PanelMetric } from "../AnalyticsPanel"
import { SeverityLegend, type SeverityBand } from "../SeverityLegend"
import { resolveOccupancySeverity } from "../../dashboard_thresholds"
import { severityFillVariables, toneClassNames } from "../../metric_tones"

interface BedOccupancyChartProps {
  occupancyRate: number
  totalBeds: number
  occupiedBeds: number
}

export const BedOccupancyChart = ({
  occupancyRate,
  totalBeds,
  occupiedBeds,
}: BedOccupancyChartProps) => {
  const { t } = useTranslation("analytics")
  const clampedRate = Math.min(Math.max(occupancyRate, 0), 100)
  const severity = resolveOccupancySeverity(clampedRate)
  const availableBeds = Math.max(totalBeds - occupiedBeds, 0)

  const severityBands: SeverityBand[] = [
    {
      severity: "critical",
      label: t("dashboard.occupancy.severity.critical"),
      range: t("dashboard.occupancy.range.critical"),
    },
    {
      severity: "elevated",
      label: t("dashboard.occupancy.severity.elevated"),
      range: t("dashboard.occupancy.range.elevated"),
    },
    {
      severity: "stable",
      label: t("dashboard.occupancy.severity.stable"),
      range: t("dashboard.occupancy.range.stable"),
    },
  ]

  return (
    <AnalyticsPanel
      title={t("dashboard.occupancy.title")}
      subtitle={t("dashboard.occupancy.subtitle")}
      trailing={
        <PanelMetric
          label={t("dashboard.occupancy.bedsShort")}
          value={`${occupiedBeds}/${totalBeds}`}
        />
      }
      className="h-full"
    >
      <div className="flex flex-1 flex-col gap-5">
        <div className="relative mx-auto w-full max-w-[280px]">
          <div className="aspect-square w-full">
            <ResponsiveContainer width="100%" height="100%">
              <RadialBarChart
                data={[{ occupancy: clampedRate }]}
                startAngle={360}
                endAngle={540}
                innerRadius="78%"
                outerRadius="100%"
              >
                <PolarAngleAxis
                  type="number"
                  domain={[0, 100]}
                  angleAxisId={0}
                  tick={false}
                  axisLine={false}
                />
                <RadialBar
                  dataKey="occupancy"
                  cornerRadius={999}
                  background={{ className: "fill-gray-200" }}
                  style={{ fill: severityFillVariables[severity] }}
                  isAnimationActive={false}
                />
              </RadialBarChart>
            </ResponsiveContainer>
          </div>

          <div className="pointer-events-none absolute inset-x-0 top-[45%] flex -translate-y-1/2 flex-col items-center">
            <span
              className={`text-4xl font-black leading-none tracking-tight tabular-nums ${toneClassNames[severity].text}`}
            >
              {clampedRate.toFixed(0)}%
            </span>
            <span className="mt-1.5 text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
              {t("dashboard.occupancy.label")}
            </span>
          </div>
        </div>

        <dl className="grid grid-cols-3 divide-x divide-border rounded-xl border border-border">
          <div className="px-3 py-3 text-center">
            <dt className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
              {t("dashboard.occupancy.occupied")}
            </dt>
            <dd className="mt-1 text-xl font-black leading-none tabular-nums text-gray-900">
              {occupiedBeds}
            </dd>
          </div>
          <div className="px-3 py-3 text-center">
            <dt className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
              {t("dashboard.occupancy.total")}
            </dt>
            <dd className="mt-1 text-xl font-black leading-none tabular-nums text-gray-900">
              {totalBeds}
            </dd>
          </div>
          <div className="px-3 py-3 text-center">
            <dt className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
              {t("dashboard.occupancy.available")}
            </dt>
            <dd className="mt-1 text-xl font-black leading-none tabular-nums text-gray-900">
              {availableBeds}
            </dd>
          </div>
        </dl>

        <SeverityLegend bands={severityBands} className="mt-auto" />
      </div>
    </AnalyticsPanel>
  )
}
