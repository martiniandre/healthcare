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
import { PageContainer } from "../../shared/components/ui/PageContainer"
import { BarChart3 } from "lucide-react"

export const Stats = () => {
  const { t: translate } = useTranslation("analytics")
  const [selectedModality, setSelectedModality] = useState<string | null>(null)
  const [hoveredBarIndex, setHoveredBarIndex] = useState<number | null>(null)

  const { data: analyticsData, isLoading, isError, refetch } = useStatsQuery()

  const examModalitiesWithCalculatedAngles = useMemo(() => {
    if (!analyticsData?.exam_modalities) {
      return []
    }
    let cumulativePercentage = 0
    return analyticsData.exam_modalities.map((item) => {
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
    if (!analyticsData?.weekly_consultations || analyticsData.weekly_consultations.length === 0) {
      return 50
    }
    const counts = analyticsData.weekly_consultations.map((item) => item.count)
    const peak = Math.max(...counts)
    return peak === 0 ? 50 : peak * 1.2
  }, [analyticsData])

  const weeklyChartSummary = useMemo(() => {
    if (!analyticsData?.weekly_consultations || analyticsData.weekly_consultations.length === 0) {
      return { min: 0, average: 0, peak: 0 }
    }
    const counts = analyticsData.weekly_consultations.map((item) => item.count)
    const total = counts.reduce((sum, val) => sum + val, 0)
    const min = Math.min(...counts)
    const peak = Math.max(...counts)
    const average = Math.round(total / analyticsData.weekly_consultations.length)
    return { min, average, peak }
  }, [analyticsData])

  if (isLoading) {
    return <StatsLoadingState />
  }

  if (isError || !analyticsData) {
    return <StatsErrorState onRetry={() => refetch()} />
  }

  const hasAnyData =
    analyticsData.total_patients > 0 ||
    (analyticsData.exam_modalities && analyticsData.exam_modalities.length > 0) ||
    (analyticsData.weekly_consultations && analyticsData.weekly_consultations.length > 0) ||
    (analyticsData.pathology_cases && analyticsData.pathology_cases.length > 0)

  if (!hasAnyData) {
    return (
      <PageContainer className="items-center justify-center">
        <EmptyState
          icon={BarChart3}
          title={translate("empty.title")}
          description={translate("empty.description")}
        />
      </PageContainer>
    )
  }

  const activeConsultationsTotal = analyticsData.weekly_consultations.reduce(
    (sum, item) => sum + item.count,
    0,
  )
  const totalStudiesCount = analyticsData.exam_modalities.reduce(
    (sum, item) => sum + item.count,
    0,
  )

  return (
    <PageContainer className="select-none">
      <StatsHeader />

      <StatsMetricsGrid 
        totalRegisteredPatients={analyticsData.total_patients}
        fhirComplianceRate={analyticsData.fhir_compliance_rate}
        averageServiceDurationMinutes={analyticsData.avg_service_duration_minutes}
        activeConsultationsTotal={activeConsultationsTotal}
      />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <StatsExamsChart 
          totalStudiesCount={totalStudiesCount}
          examModalitiesData={analyticsData.exam_modalities}
          selectedModality={selectedModality}
          setSelectedModality={setSelectedModality}
          examModalitiesWithCalculatedAngles={examModalitiesWithCalculatedAngles}
        />

        <StatsConsultationsChart 
          consultationsWeeklyData={analyticsData.weekly_consultations}
          maxWeeklyConsultationCount={maxWeeklyConsultationCount}
          weeklyChartSummary={weeklyChartSummary}
          hoveredBarIndex={hoveredBarIndex}
          setHoveredBarIndex={setHoveredBarIndex}
        />
      </div>

      <StatsEpidemiologyTable pathologies={analyticsData.pathology_cases} />
    </PageContainer>
  )
}
