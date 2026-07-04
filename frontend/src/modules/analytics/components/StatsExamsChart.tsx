import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { EmptyState } from "../../../shared/components/ui/EmptyState"

interface ExamModality {
  modality: string
  count: number
  percentage: number
  color: string
}

interface StatsExamsChartProps {
  totalStudiesCount: number
  examModalitiesData: ExamModality[]
  selectedModality: string | null
  setSelectedModality: (modality: string | null) => void
  examModalitiesWithCalculatedAngles: (ExamModality & { dashoffset: number; rotationAngle: number })[]
}

export const StatsExamsChart = ({
  totalStudiesCount,
  examModalitiesData,
  selectedModality,
  setSelectedModality,
  examModalitiesWithCalculatedAngles
}: StatsExamsChartProps) => {
  const { t: translate } = useTranslation()

  return (
    <Card className="p-5 flex flex-col gap-5 text-left border border-border">
      <div>
        <h3 className="font-extrabold text-gray-900 text-md">{translate("analytics.exams.title")}</h3>
        <span className="text-xs text-muted block mt-1">{translate("analytics.exams.subtitle")}</span>
      </div>

      <div className="flex flex-col sm:flex-row items-center justify-around gap-6 py-4">
        <div className="relative w-44 h-44 flex items-center justify-center shrink-0">
          <svg className="w-full h-full transform -rotate-90" viewBox="0 0 100 100">
            <circle cx="50" cy="50" r="38" fill="transparent" stroke="#f3f4f6" strokeWidth="8" />
            {examModalitiesWithCalculatedAngles.map((item) => (
              <circle
                key={item.modality}
                cx="50"
                cy="50"
                r="38"
                fill="transparent"
                stroke={item.color}
                strokeWidth="8.5"
                strokeDasharray="238.76"
                strokeDashoffset={item.dashoffset}
                className="transform origin-center"
                style={{
                  transform: `rotate(${item.rotationAngle}deg)`,
                  transformOrigin: "center",
                }}
              />
            ))}
          </svg>

          <div className="absolute flex flex-col items-center justify-center text-center">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider">{translate("analytics.exams.totalLabel")}</span>
            <span className="text-xl font-black text-gray-900">{totalStudiesCount}</span>
            <span className="text-[9px] text-muted font-semibold mt-0.5">{translate("analytics.exams.studiesLabel")}</span>
          </div>
        </div>

        <div className="flex flex-col gap-3 w-full">
          {examModalitiesData?.length ? (
            examModalitiesData.map((item) => (
              <div
                key={item.modality}
                onMouseEnter={() => setSelectedModality(item.modality)}
                onMouseLeave={() => setSelectedModality(null)}
                className={`flex items-center justify-between p-2.5 rounded-lg border transition-all duration-200 ${
                  selectedModality === item.modality 
                    ? "bg-gray-50 border-gray-300" 
                    : "bg-white border-transparent"
                }`}
              >
                <div className="flex items-center gap-2.5">
                  <div className="w-3.5 h-3.5 rounded-full shrink-0" style={{ backgroundColor: item.color }} />
                  <span className="text-xs font-bold text-gray-700">{item.modality}</span>
                </div>
                <div className="text-right">
                  <span className="text-xs font-black text-gray-900 block">{item.count} {translate("analytics.exams.examesUnit")}</span>
                  <span className="text-[10px] text-gray-500 font-semibold">{item.percentage}%</span>
                </div>
              </div>
            ))
          ) : (
            <EmptyState 
              title={translate("analytics.empty.exams")} 
              description={translate("analytics.empty.examsDesc")} 
              className="h-full"
            />
          )}
        </div>
      </div>
    </Card>
  )
}
