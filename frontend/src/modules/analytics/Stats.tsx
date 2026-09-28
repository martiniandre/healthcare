import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { BarChart3 } from "lucide-react"
import { useStatsQuery } from "./queries"
import { StatsHeader } from "./components/StatsHeader"
import { StatsMetricStrip } from "./components/StatsMetricStrip"
import { ExamModalityDonutChart } from "./components/charts/ExamModalityDonutChart"
import { WeeklyConsultationsChart } from "./components/charts/WeeklyConsultationsChart"
import { StatsEpidemiologyTable } from "./components/StatsEpidemiologyTable"
import { StatsLoadingState } from "./components/StatsLoadingState"
import { StatsErrorState } from "./components/StatsErrorState"
import { EmptyState } from "../../shared/components/ui/EmptyState"
import { PageContainer } from "../../shared/components/ui/PageContainer"
import { summarizeWeeklyConsultations } from "./stats_domain"

export const Stats = () => {
  const { t } = useTranslation("analytics")
  const { data: analyticsData, isLoading, isError, refetch } = useStatsQuery()

  const weeklySummary = useMemo(
    () => summarizeWeeklyConsultations(analyticsData?.weekly_consultations ?? []),
    [analyticsData]
  )

  if (isLoading) {
    return <StatsLoadingState />
  }

  if (isError || !analyticsData) {
    return <StatsErrorState onRetry={() => refetch()} />
  }

  const hasAnyData =
    analyticsData.total_patients > 0 ||
    analyticsData.exam_modalities.length > 0 ||
    analyticsData.weekly_consultations.length > 0 ||
    analyticsData.pathology_cases.length > 0

  if (!hasAnyData) {
    return (
      <PageContainer className="items-center justify-center">
        <EmptyState icon={BarChart3} title={t("empty.title")} description={t("empty.description")} />
      </PageContainer>
    )
  }

  return (
    <PageContainer>
      <StatsHeader />

      <StatsMetricStrip
        totalRegisteredPatients={analyticsData.total_patients}
        fhirComplianceRate={analyticsData.fhir_compliance_rate}
        averageServiceDurationMinutes={analyticsData.avg_service_duration_minutes}
        weeklyConsultationsTotal={weeklySummary.total}
        weeklyPeakDayLabel={
          weeklySummary.peakDayKey === null ? null : t(weeklySummary.peakDayKey)
        }
      />

      <div className="grid grid-cols-1 gap-4 md:gap-6 xl:grid-cols-12">
        <div className="xl:col-span-5">
          <ExamModalityDonutChart examModalities={analyticsData.exam_modalities} />
        </div>
        <div className="xl:col-span-7">
          <WeeklyConsultationsChart weeklyConsultations={analyticsData.weekly_consultations} />
        </div>
      </div>

      <StatsEpidemiologyTable pathologies={analyticsData.pathology_cases} />
    </PageContainer>
  )
}
