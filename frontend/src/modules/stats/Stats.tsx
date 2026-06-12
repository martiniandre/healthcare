import { useState, useMemo } from "react"
import { useStatsQuery } from "./queries"
import { StatsHeader } from "./components/StatsHeader"
import { StatsMetricsGrid } from "./components/StatsMetricsGrid"
import { StatsExamsChart } from "./components/StatsExamsChart"
import { StatsConsultationsChart } from "./components/StatsConsultationsChart"
import { StatsEpidemiologyTable } from "./components/StatsEpidemiologyTable"
import { StatsLoadingState } from "./components/StatsLoadingState"
import { StatsErrorState } from "./components/StatsErrorState"

export const Stats = () => {
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

  if (isLoading) {
    return <StatsLoadingState />
  }

  if (isError || !statsData) {
    return <StatsErrorState />
  }

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full select-none">
      <StatsHeader />

      <StatsMetricsGrid 
        totalRegisteredPatients={statsData.totalRegisteredPatients}
        fhirComplianceRate={statsData.fhirComplianceRate}
        averageServiceDurationMinutes={statsData.averageServiceDurationMinutes}
        activeConsultationsTotal={statsData.activeConsultationsTotal}
      />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <StatsExamsChart 
          totalStudiesCount={statsData.totalStudiesCount}
          examModalitiesData={statsData.examModalitiesData}
          selectedModality={selectedModality}
          setSelectedModality={setSelectedModality}
          examModalitiesWithCalculatedAngles={examModalitiesWithCalculatedAngles}
        />

        <StatsConsultationsChart 
          consultationsWeeklyData={statsData.consultationsWeeklyData}
          maxWeeklyConsultationCount={maxWeeklyConsultationCount}
          weeklyChartSummary={weeklyChartSummary}
          hoveredBarIndex={hoveredBarIndex}
          setHoveredBarIndex={setHoveredBarIndex}
        />
      </div>

      <StatsEpidemiologyTable pathologies={statsData.pathologies} />
    </div>
  )
}
