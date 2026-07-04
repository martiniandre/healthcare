import { useState, useMemo } from "react"
import { useTranslation } from "react-i18next"
import { useStatsQuery } from "./queries"
import { StatsHeader } from "./components/StatsHeader"
import { StatsMetricsGrid } from "./components/StatsMetricsGrid"
import { StatsExamsChart } from "./components/StatsExamsChart"
import { StatsConsultationsChart } from "./components/StatsConsultationsChart"
import { StatsEpidemiologyTable } from "./components/StatsEpidemiologyTable"
import { StatsLoadingState } from "./components/StatsLoadingState"
import { StatsErrorState } from "./components/StatsErrorState"
import { EmptyState } from "../../shared/components/ui/EmptyState"
import { BarChart3 } from "lucide-react"

export const Stats = () => {
  const { t: translate } = useTranslation()
  const [selectedModality, setSelectedModality] = useState<string | null>(null)
  const [hoveredBarIndex, setHoveredBarIndex] = useState<number | null>(null)

  const { data: analyticsData, isLoading, isError } = useStatsQuery()

  const examModalitiesWithCalculatedAngles = useMemo(() => {
    if (!analyticsData?.examModalitiesData) {
      return []
    }
    let cumulativePercentage = 0
    return analyticsData.examModalitiesData.map((item) => {
      const currentItemAngle = cumulativePercentage * 3.6
      cumulativePercentage += item.percentage
      return {
        ...item,
        dashoffset: 238.76 - (238.76 * item.percentage) / 100,
        rotationAngle: currentItemAngle,
      }
    })
  }, [analyticsData])

  const maxWeeklyConsultationCount = useMemo(() => {
    if (!analyticsData?.consultationsWeeklyData || analyticsData.consultationsWeeklyData.length === 0) {
      return 50
    }
    const counts = analyticsData.consultationsWeeklyData.map((item) => item.count)
    const peak = Math.max(...counts)
    return peak === 0 ? 50 : peak * 1.2
  }, [analyticsData])

  const weeklyChartSummary = useMemo(() => {
    if (!analyticsData?.consultationsWeeklyData || analyticsData.consultationsWeeklyData.length === 0) {
      return { min: 0, average: 0, peak: 0 }
    }
    const counts = analyticsData.consultationsWeeklyData.map((item) => item.count)
    const total = counts.reduce((sum, val) => sum + val, 0)
    const min = Math.min(...counts)
    const peak = Math.max(...counts)
    const average = Math.round(total / analyticsData.consultationsWeeklyData.length)
    return { min, average, peak }
  }, [analyticsData])

  if (isLoading) {
    return <StatsLoadingState />
  }

  if (isError || !analyticsData) {
    return <StatsErrorState />
  }

  const hasAnyData =
    analyticsData.totalRegisteredPatients > 0 ||
    (analyticsData.examModalitiesData && analyticsData.examModalitiesData.length > 0) ||
    (analyticsData.consultationsWeeklyData && analyticsData.consultationsWeeklyData.length > 0) ||
    (analyticsData.pathologies && analyticsData.pathologies.length > 0)

  if (!hasAnyData) {
    return (
      <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col items-center justify-center max-w-7xl mx-auto w-full">
        <EmptyState
          icon={BarChart3}
          title={translate("analytics.empty.title")}
          description={translate("analytics.empty.description")}
        />
      </div>
    )
  }

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full select-none">
      <StatsHeader />

      <StatsMetricsGrid 
        totalRegisteredPatients={analyticsData.totalRegisteredPatients}
        fhirComplianceRate={analyticsData.fhirComplianceRate}
        averageServiceDurationMinutes={analyticsData.averageServiceDurationMinutes}
        activeConsultationsTotal={analyticsData.activeConsultationsTotal}
      />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <StatsExamsChart 
          totalStudiesCount={analyticsData.totalStudiesCount}
          examModalitiesData={analyticsData.examModalitiesData}
          selectedModality={selectedModality}
          setSelectedModality={setSelectedModality}
          examModalitiesWithCalculatedAngles={examModalitiesWithCalculatedAngles}
        />

        <StatsConsultationsChart 
          consultationsWeeklyData={analyticsData.consultationsWeeklyData}
          maxWeeklyConsultationCount={maxWeeklyConsultationCount}
          weeklyChartSummary={weeklyChartSummary}
          hoveredBarIndex={hoveredBarIndex}
          setHoveredBarIndex={setHoveredBarIndex}
        />
      </div>

      <StatsEpidemiologyTable pathologies={analyticsData.pathologies} />
    </div>
  )
}
