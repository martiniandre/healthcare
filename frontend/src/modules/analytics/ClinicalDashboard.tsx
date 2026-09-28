import { useTranslation } from "react-i18next"
import { Activity, RefreshCw, TriangleAlert } from "lucide-react"
import { useDashboardQuery } from "./dashboard_queries"
import { DashboardKPICards } from "./components/DashboardKPICards"
import { ConsultationsByDoctorChart } from "./components/charts/ConsultationsByDoctorChart"
import { BedOccupancyChart } from "./components/charts/BedOccupancyChart"
import { WaitTimeByDepartmentChart } from "./components/charts/WaitTimeByDepartmentChart"
import { TopDiagnosesTable } from "./components/TopDiagnosesTable"
import { PageContainer, PageTitle } from "../../shared/components/ui/PageContainer"
import { Card } from "../../shared/components/ui/Card"
import { Skeleton } from "../../shared/components/ui/Skeleton"
import { Button } from "../../shared/components/ui/Button"
import { cn } from "../../shared/utils/cn"

export const ClinicalDashboard = () => {
  const { t, i18n } = useTranslation("analytics")
  const {
    data: dashboardData,
    isLoading,
    isError,
    isFetching,
    dataUpdatedAt,
    refetch,
  } = useDashboardQuery()

  if (isLoading) {
    return <ClinicalDashboardSkeleton />
  }

  if (isError || !dashboardData) {
    return <ClinicalDashboardError onRetry={() => refetch()} />
  }

  const lastUpdatedLabel = new Intl.DateTimeFormat(i18n.language, {
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(dataUpdatedAt))

  return (
    <PageContainer>
      <PageTitle
        icon={<Activity className="size-5 text-primary animate-pulse-glow" />}
        title={t("dashboard.title")}
        description={t("dashboard.subtitle")}
        actions={
          <div className="flex items-center gap-3">
            <span className="flex items-center gap-2 rounded-lg border border-border bg-gray-50 px-2.5 py-1.5">
              <span className="relative flex size-1.5 shrink-0" aria-hidden="true">
                {isFetching ? (
                  <span className="absolute inline-flex size-full animate-ping rounded-full bg-warning opacity-75" />
                ) : null}
                <span
                  className={cn(
                    "relative inline-flex size-1.5 rounded-full",
                    isFetching ? "bg-warning" : "bg-success"
                  )}
                />
              </span>
              <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
                {isFetching ? t("dashboard.syncing") : t("dashboard.lastUpdated")}
              </span>
              {isFetching ? null : (
                <span className="text-xs font-black tabular-nums text-gray-900">
                  {lastUpdatedLabel}
                </span>
              )}
            </span>
            <Button
              variantType="outline"
              onClick={() => refetch()}
              isLoading={isFetching}
              aria-label={t("dashboard.refresh")}
              className="size-9 p-0"
            >
              <RefreshCw className="size-4" aria-hidden="true" />
            </Button>
          </div>
        }
      />

      <DashboardKPICards
        consultationsToday={dashboardData.consultations_today}
        consultationsTrend={dashboardData.consultations_trend}
        occupancyRate={dashboardData.occupancy_rate}
        avgWaitTimeMinutes={dashboardData.avg_wait_time_minutes}
        activePatients={dashboardData.active_patients}
        examsToday={dashboardData.exams_today}
        newDiagnosesToday={dashboardData.new_diagnoses_today}
      />

      <div className="grid grid-cols-1 gap-4 md:gap-6 xl:grid-cols-12">
        <div className="xl:col-span-7">
          <ConsultationsByDoctorChart
            consultationsPerDoctor={dashboardData.consultations_per_doctor}
          />
        </div>
        <div className="xl:col-span-5">
          <BedOccupancyChart
            occupancyRate={dashboardData.occupancy_rate}
            totalBeds={dashboardData.occupancy_total_beds}
            occupiedBeds={dashboardData.occupancy_occupied_beds}
          />
        </div>
      </div>

      <WaitTimeByDepartmentChart
        waitTimeByDepartment={dashboardData.wait_time_by_department}
      />

      <TopDiagnosesTable topDiagnoses={dashboardData.top_diagnoses} />
    </PageContainer>
  )
}

const ClinicalDashboardSkeleton = () => {
  return (
    <PageContainer>
      <div>
        <Skeleton className="h-6 w-48" />
        <Skeleton className="mt-2.5 h-3 w-72" />
      </div>

      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6 lg:gap-4">
        {Array.from({ length: 6 }).map((_, indexValue) => (
          <Card key={String(indexValue)} className="flex flex-col gap-3 p-4">
            <div className="flex items-start justify-between gap-2">
              <Skeleton className="h-3 w-20" />
              <Skeleton className="size-8 rounded-lg" />
            </div>
            <Skeleton className="h-8 w-16" />
            <Skeleton className="h-4 w-4/5" />
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 gap-4 md:gap-6 xl:grid-cols-12">
        <Card className="flex flex-col gap-5 p-5 xl:col-span-7">
          <div className="flex items-start justify-between gap-4">
            <div className="flex flex-col gap-2">
              <Skeleton className="h-4 w-44" />
              <Skeleton className="h-3 w-64" />
            </div>
            <Skeleton className="h-11 w-20 rounded-lg" />
          </div>
          <div className="flex flex-col gap-6">
            {Array.from({ length: 5 }).map((_, indexValue) => (
              <div key={String(indexValue)} className="flex items-center gap-3">
                <Skeleton className="h-7 w-36 shrink-0" />
                <Skeleton className="h-6 flex-1 rounded-md" />
              </div>
            ))}
          </div>
        </Card>

        <Card className="flex flex-col gap-5 p-5 xl:col-span-5">
          <div className="flex items-start justify-between gap-4">
            <div className="flex flex-col gap-2">
              <Skeleton className="h-4 w-36" />
              <Skeleton className="h-3 w-52" />
            </div>
            <Skeleton className="h-11 w-20 rounded-lg" />
          </div>
          <Skeleton className="mx-auto aspect-square w-full max-w-[280px] rounded-full" />
          <div className="grid grid-cols-3 gap-3">
            {Array.from({ length: 3 }).map((_, indexValue) => (
              <Skeleton key={String(indexValue)} className="h-16 w-full rounded-xl" />
            ))}
          </div>
        </Card>
      </div>

      <Card className="flex flex-col gap-5 p-5">
        <div className="flex items-start justify-between gap-4">
          <div className="flex flex-col gap-2">
            <Skeleton className="h-4 w-48" />
            <Skeleton className="h-3 w-64" />
          </div>
          <Skeleton className="h-11 w-20 rounded-lg" />
        </div>
        <div className="flex flex-col gap-5">
          {Array.from({ length: 4 }).map((_, indexValue) => (
            <div key={String(indexValue)} className="flex items-center gap-3">
              <Skeleton className="h-4 w-40 shrink-0" />
              <Skeleton className="h-5 flex-1 rounded-md" />
            </div>
          ))}
        </div>
      </Card>

      <Card className="flex flex-col gap-4 p-5">
        <div className="flex items-start justify-between gap-4">
          <div className="flex flex-col gap-2">
            <Skeleton className="h-4 w-36" />
            <Skeleton className="h-3 w-60" />
          </div>
          <Skeleton className="h-11 w-20 rounded-lg" />
        </div>
        <div className="flex flex-col gap-3">
          {Array.from({ length: 5 }).map((_, indexValue) => (
            <Skeleton key={String(indexValue)} className="h-12 w-full rounded-lg" />
          ))}
        </div>
      </Card>
    </PageContainer>
  )
}

interface ClinicalDashboardErrorProps {
  onRetry: () => void
}

const ClinicalDashboardError = ({ onRetry }: ClinicalDashboardErrorProps) => {
  const { t } = useTranslation("analytics")

  return (
    <PageContainer className="flex items-center justify-center">
      <Card className="flex w-full max-w-md flex-col items-center gap-4 border-danger/20 p-8 text-center">
        <span className="grid size-14 place-items-center rounded-full bg-danger/10">
          <TriangleAlert className="size-7 text-danger" aria-hidden="true" />
        </span>
        <h3 className="font-display text-lg font-bold text-gray-900">
          {t("dashboard.errorTitle")}
        </h3>
        <p className="text-xs leading-relaxed text-muted">{t("dashboard.errorDescription")}</p>
        <Button variantType="danger" onClick={onRetry} className="mt-2 w-full">
          {t("dashboard.retryButton")}
        </Button>
      </Card>
    </PageContainer>
  )
}
