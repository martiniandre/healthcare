import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { EmptyState } from "../../../shared/components/ui/EmptyState"

interface WeeklyConsultationData {
  dayName: string
  count: number
}

interface StatsConsultationsChartProps {
  consultationsWeeklyData: WeeklyConsultationData[]
  maxWeeklyConsultationCount: number
  weeklyChartSummary: { min: number; average: number; peak: number }
  hoveredBarIndex: number | null
  setHoveredBarIndex: (index: number | null) => void
}

export const StatsConsultationsChart = ({
  consultationsWeeklyData,
  maxWeeklyConsultationCount,
  weeklyChartSummary,
  hoveredBarIndex,
  setHoveredBarIndex
}: StatsConsultationsChartProps) => {
  const { t: translate } = useTranslation()

  return (
    <Card className="p-5 flex flex-col gap-5 text-left border border-border">
      <div>
        <h3 className="font-extrabold text-gray-900 text-md">{translate("analytics.consultations.title")}</h3>
        <span className="text-xs text-muted block mt-1">{translate("analytics.consultations.subtitle")}</span>
      </div>

      <div className="overflow-x-auto w-full">
        {consultationsWeeklyData?.length ? (
          <div className="flex items-end justify-between gap-2.5 h-48 border-b border-border pb-2 pt-6 px-4 min-w-[280px]">
            {consultationsWeeklyData.map((item, indexValue) => {
              const isHovered = hoveredBarIndex === indexValue
              const percentageHeight = (item.count / maxWeeklyConsultationCount) * 100
              const translatedDayName = item.dayName.startsWith("analytics.days.") 
                ? translate(item.dayName) 
                : translate(`analytics.days.${item.dayName}`)

              return (
                <div
                  key={item.dayName}
                  className="flex-1 flex flex-col items-center gap-2 group relative"
                  onMouseEnter={() => setHoveredBarIndex(indexValue)}
                  onMouseLeave={() => setHoveredBarIndex(null)}
                >
                  {isHovered && (
                    <div className="absolute -top-10 bg-gray-900 text-white text-[10px] font-bold px-2 py-1 rounded shadow-md z-10 whitespace-nowrap">
                      {item.count} {translate("analytics.consultations.consultationsLabel")}
                    </div>
                  )}

                  <div
                    className={`w-full max-w-[28px] rounded-t-md transition-all duration-300 ${
                      isHovered ? "bg-primary" : "bg-primary/20"
                    }`}
                    style={{ height: `${percentageHeight}%` }}
                  />

                  <span className="text-[10px] font-bold text-gray-500 uppercase">
                    {translatedDayName}
                  </span>
                </div>
              )
            })}
          </div>
        ) : (
          <div className="h-48 border-b border-border mb-2 pt-6">
            <EmptyState 
              title={translate("analytics.empty.consultations")} 
              description={translate("analytics.empty.consultationsDesc")} 
              className="h-full"
            />
          </div>
        )}
      </div>

      <div className="flex justify-between items-center text-xs text-gray-500 px-2.5">
        <span>{translate("analytics.consultations.minLabel")} {weeklyChartSummary.min} {translate("analytics.consultations.sunShort")}</span>
        <span>{translate("analytics.consultations.avgLabel")} {weeklyChartSummary.average} {translate("analytics.consultations.perDay")}</span>
        <span>{translate("analytics.consultations.peakLabel")} {weeklyChartSummary.peak} {translate("analytics.consultations.friShort")}</span>
      </div>
    </Card>
  )
}
