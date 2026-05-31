import { useState, useMemo } from "react"
import { useTranslation } from "react-i18next"
import { Card } from "../../shared/components/ui/Card"
import { Button } from "../../shared/components/ui/Button"
import { useStatsQuery } from "./queries"
import { 
  BarChart3, 
  Users, 
  Clock, 
  CheckSquare, 
  ArrowUpRight, 
  FileSpreadsheet, 
  Activity 
} from "lucide-react"

export const Stats = () => {
  const { t: translate } = useTranslation()
  const [selectedModality, setSelectedModality] = useState<string | null>(null)
  const [hoveredBarIndex, setHoveredBarIndex] = useState<number | null>(null)

  const { data: statsData, isLoading, isError } = useStatsQuery()

  const examModalitiesWithCalculatedAngles = useMemo(() => {
    if (!statsData?.examModalitiesData) {
      return []
    }
    let cumulativePercentage = 0
    return statsData.examModalitiesData.map((item) => {
      const currentItemAngle = cumulativePercentage * 3.6
      cumulativePercentage += item.percentage
      return {
        ...item,
        dashoffset: 238.76 - (238.76 * item.percentage) / 100,
        rotationAngle: currentItemAngle,
      }
    })
  }, [statsData])

  const maxWeeklyConsultationCount = useMemo(() => {
    if (!statsData?.consultationsWeeklyData || statsData.consultationsWeeklyData.length === 0) {
      return 50
    }
    const counts = statsData.consultationsWeeklyData.map((item) => item.count)
    const peak = Math.max(...counts)
    return peak === 0 ? 50 : peak * 1.2
  }, [statsData])

  const weeklyChartSummary = useMemo(() => {
    if (!statsData?.consultationsWeeklyData || statsData.consultationsWeeklyData.length === 0) {
      return { min: 0, average: 0, peak: 0 }
    }
    const counts = statsData.consultationsWeeklyData.map((item) => item.count)
    const total = counts.reduce((sum, val) => sum + val, 0)
    const min = Math.min(...counts)
    const peak = Math.max(...counts)
    const average = Math.round(total / statsData.consultationsWeeklyData.length)
    return { min, average, peak }
  }, [statsData])


  const getTrendStyle = (pathologyCode: string): string => {
    if (pathologyCode === "E11.9") {
      return "text-red-500 font-bold"
    }
    if (pathologyCode === "J45.9") {
      return "text-emerald-600 font-bold"
    }
    return "text-gray-400 font-bold"
  }

  if (isLoading) {
    return (
      <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full animate-pulse">
        <div className="text-left">
          <div className="h-6 w-48 bg-gray-200 rounded"></div>
          <div className="h-3 w-72 bg-gray-200 rounded mt-2.5"></div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
          {[1, 2, 3, 4].map((indexValue) => (
            <Card key={indexValue} className="p-4 flex items-center justify-between border border-border h-24">
              <div className="flex-1 flex flex-col gap-2">
                <div className="h-3 w-20 bg-gray-200 rounded"></div>
                <div className="h-6 w-16 bg-gray-200 rounded"></div>
                <div className="h-3 w-28 bg-gray-200 rounded"></div>
              </div>
              <div className="w-12 h-12 bg-gray-200 rounded-xl"></div>
            </Card>
          ))}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card className="p-5 flex flex-col gap-5 border border-border h-80">
            <div className="h-4 w-32 bg-gray-200 rounded"></div>
            <div className="h-3 w-48 bg-gray-200 rounded"></div>
            <div className="flex items-center gap-6 mt-4">
              <div className="w-36 h-36 rounded-full bg-gray-200"></div>
              <div className="flex-1 flex flex-col gap-3">
                {[1, 2, 3, 4].map((indexValue) => (
                  <div key={indexValue} className="h-8 bg-gray-200 rounded"></div>
                ))}
              </div>
            </div>
          </Card>

          <Card className="p-5 flex flex-col gap-5 border border-border h-80">
            <div className="h-4 w-32 bg-gray-200 rounded"></div>
            <div className="h-3 w-48 bg-gray-200 rounded"></div>
            <div className="flex-1 flex items-end justify-between gap-3 h-40 pb-2">
              {[1, 2, 3, 4, 5, 6, 7].map((indexValue) => (
                <div key={indexValue} className="flex-1 h-32 bg-gray-200 rounded-t-md"></div>
              ))}
            </div>
          </Card>
        </div>

        <Card className="p-5 flex flex-col gap-4 border border-border">
          <div className="flex justify-between items-center pb-3">
            <div className="flex flex-col gap-2">
              <div className="h-4 w-40 bg-gray-200 rounded"></div>
              <div className="h-3 w-60 bg-gray-200 rounded"></div>
            </div>
            <div className="w-24 h-8 bg-gray-200 rounded"></div>
          </div>
          <div className="flex flex-col gap-4 mt-2">
            {[1, 2, 3, 4].map((indexValue) => (
              <div key={indexValue} className="h-10 bg-gray-200 rounded"></div>
            ))}
          </div>
        </Card>
      </div>
    )
  }

  if (isError || !statsData) {
    return (
      <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col items-center justify-center gap-4 max-w-7xl mx-auto w-full select-none">
        <div className="text-center p-8 bg-white border border-red-100 shadow-xl rounded-2xl max-w-md w-full flex flex-col items-center gap-4">
          <div className="bg-red-50 p-4 rounded-full">
            <Activity className="w-10 h-10 text-red-500 animate-bounce" />
          </div>
          <h3 className="text-lg font-black text-gray-900">{translate("stats.errorTitle") || "Erro ao carregar dados"}</h3>
          <p className="text-xs text-gray-500 leading-relaxed">
            {translate("stats.errorDescription") || "Não foi possível estabelecer conexão com o serviço de analytics FHIR."}
          </p>
          <Button 
            onClick={() => window.location.reload()} 
            className="w-full bg-red-600 hover:bg-red-700 text-white font-bold py-2 rounded-xl transition-all duration-200 mt-2"
          >
            {translate("stats.retryButton") || "Tentar Novamente"}
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full select-none">
      <div className="text-left">
        <h2 className="text-xl font-black text-gray-900 leading-none flex items-center gap-2">
          <BarChart3 className="w-5 h-5 text-primary animate-pulse-glow" />
          {translate("stats.title")}
        </h2>
        <span className="text-xs text-muted mt-1.5 block">
          {translate("stats.subtitle")}
        </span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
        <Card className="p-4 flex items-center justify-between border border-border">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{translate("stats.metrics.activePatients")}</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">{statsData.totalRegisteredPatients}</span>
            <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
              <ArrowUpRight className="w-3.5 h-3.5" />
              {translate("stats.metrics.admissionsGrowth")}
            </span>
          </div>
          <div className="bg-primary/8 p-3 rounded-xl">
            <Users className="w-6 h-6 text-primary" />
          </div>
        </Card>

        <Card className="p-4 flex items-center justify-between border border-border">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{translate("stats.metrics.fhirCompliance")}</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">{statsData.fhirComplianceRate}%</span>
            <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
              <CheckSquare className="w-3.5 h-3.5" />
              {translate("stats.metrics.fhirCompliantDesc")}
            </span>
          </div>
          <div className="bg-emerald-50 p-3 rounded-xl">
            <CheckSquare className="w-6 h-6 text-emerald-600" />
          </div>
        </Card>

        <Card className="p-4 flex items-center justify-between border border-border">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{translate("stats.metrics.avgConsultationTime")}</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">{statsData.averageServiceDurationMinutes} min</span>
            <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
              <Clock className="w-3.5 h-3.5" />
              {translate("stats.metrics.consultationTrend")}
            </span>
          </div>
          <div className="bg-purple-50 p-3 rounded-xl">
            <Clock className="w-6 h-6 text-purple-600" />
          </div>
        </Card>

        <Card className="p-4 flex items-center justify-between border border-border">
          <div className="text-left">
            <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{translate("stats.metrics.weeklyConsultations")}</span>
            <span className="text-2xl font-black text-gray-900 mt-1 block">{statsData.activeConsultationsTotal}</span>
            <span className="text-[10px] text-amber-600 font-bold flex items-center gap-1 mt-1.5">
              <Activity className="w-3.5 h-3.5" />
              {translate("stats.metrics.activeIcuBeds")}
            </span>
          </div>
          <div className="bg-amber-50 p-3 rounded-xl">
            <Activity className="w-6 h-6 text-amber-500" />
          </div>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card className="p-5 flex flex-col gap-5 text-left border border-border">
          <div>
            <h3 className="font-extrabold text-gray-900 text-md">{translate("stats.exams.title")}</h3>
            <span className="text-xs text-muted block mt-1">{translate("stats.exams.subtitle")}</span>
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
                <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider">{translate("stats.exams.totalLabel")}</span>
                <span className="text-xl font-black text-gray-900">{statsData.totalStudiesCount}</span>
                <span className="text-[9px] text-muted font-semibold mt-0.5">{translate("stats.exams.studiesLabel")}</span>
              </div>
            </div>

            <div className="flex flex-col gap-3 w-full">
              {statsData.examModalitiesData.map((item) => (
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
                    <span className="text-xs font-black text-gray-900 block">{item.count} {translate("stats.exams.examesUnit")}</span>
                    <span className="text-[10px] text-gray-500 font-semibold">{item.percentage}%</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </Card>

        <Card className="p-5 flex flex-col gap-5 text-left border border-border">
          <div>
            <h3 className="font-extrabold text-gray-900 text-md">{translate("stats.consultations.title")}</h3>
            <span className="text-xs text-muted block mt-1">{translate("stats.consultations.subtitle")}</span>
          </div>

          <div className="overflow-x-auto w-full">
            <div className="flex items-end justify-between gap-2.5 h-48 border-b border-border pb-2 pt-6 px-4 min-w-[280px]">
              {statsData.consultationsWeeklyData.map((item, indexValue) => {
                const isHovered = hoveredBarIndex === indexValue
                const percentageHeight = (item.count / maxWeeklyConsultationCount) * 100
                const translatedDayName = item.dayName.startsWith("stats.days.") 
                  ? translate(item.dayName) 
                  : translate(`stats.days.${item.dayName}`)

                return (
                  <div
                    key={item.dayName}
                    className="flex-1 flex flex-col items-center gap-2 group relative"
                    onMouseEnter={() => setHoveredBarIndex(indexValue)}
                    onMouseLeave={() => setHoveredBarIndex(null)}
                  >
                    {isHovered && (
                      <div className="absolute -top-10 bg-gray-900 text-white text-[10px] font-bold px-2 py-1 rounded shadow-md z-10 whitespace-nowrap">
                        {item.count} {translate("stats.consultations.consultationsLabel")}
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
          </div>

          <div className="flex justify-between items-center text-xs text-gray-500 px-2.5">
            <span>{translate("stats.consultations.minLabel")} {weeklyChartSummary.min} {translate("stats.consultations.sunShort")}</span>
            <span>{translate("stats.consultations.avgLabel")} {weeklyChartSummary.average} {translate("stats.consultations.perDay")}</span>
            <span>{translate("stats.consultations.peakLabel")} {weeklyChartSummary.peak} {translate("stats.consultations.friShort")}</span>
          </div>
        </Card>
      </div>

      <Card className="p-5 flex flex-col gap-4 text-left border border-border">
        <div className="flex items-center justify-between border-b border-border pb-3 flex-wrap gap-2">
          <div>
            <h3 className="font-extrabold text-gray-900 text-md">{translate("stats.epidemiology.title")}</h3>
            <span className="text-xs text-muted block mt-1">{translate("stats.epidemiology.subtitle")}</span>
          </div>
          <Button variantType="outline" className="px-3 gap-1.5 text-xs">
            <FileSpreadsheet className="w-4 h-4" />
            {translate("stats.epidemiology.exportButton")}
          </Button>
        </div>

        <div className="overflow-x-auto w-full">
          <table className="w-full text-left text-xs border-collapse min-w-[500px] md:min-w-0">
            <thead>
              <tr className="border-b border-border text-gray-500 font-bold uppercase tracking-wider">
                <th className="py-3 px-3">{translate("stats.epidemiology.table.code")}</th>
                <th className="py-3 px-3">{translate("stats.epidemiology.table.description")}</th>
                <th className="py-3 px-3">{translate("stats.epidemiology.table.category")}</th>
                <th className="py-3 px-3">{translate("stats.epidemiology.table.activeCases")}</th>
                <th className="py-3 px-3 text-right">{translate("stats.epidemiology.table.trend")}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border text-gray-700 font-medium">
              {statsData.pathologies.map((pathologyItem) => {
                const translatedDescription = pathologyItem.descriptionKey.startsWith("stats.")
                  ? translate(pathologyItem.descriptionKey)
                  : translate(`stats.pathologies.${pathologyItem.descriptionKey}`)

                const translatedCategory = pathologyItem.categoryKey.startsWith("stats.")
                  ? translate(pathologyItem.categoryKey)
                  : translate(`stats.categories.${pathologyItem.categoryKey}`)

                return (
                  <tr key={pathologyItem.code} className="hover:bg-gray-50/50">
                    <td className="py-3 px-3 font-mono font-bold text-primary">{pathologyItem.code}</td>
                    <td className="py-3 px-3">{translatedDescription}</td>
                    <td className="py-3 px-3">{translatedCategory}</td>
                    <td className="py-3 px-3 font-bold text-gray-900">{pathologyItem.activeCases}</td>
                    <td className="py-3 px-3 text-right">
                      <span className={getTrendStyle(pathologyItem.code)}>
                        {pathologyItem.trend === "stable" ? translate("stats.epidemiology.table.stable") : pathologyItem.trend}
                      </span>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  )
}
