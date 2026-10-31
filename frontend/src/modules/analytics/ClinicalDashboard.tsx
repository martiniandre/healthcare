import { useTranslation } from "react-i18next";
import { Activity, RefreshCw } from "lucide-react";
import { useDashboardQuery } from "./dashboard_queries";
import { DashboardKPICards } from "./components/DashboardKPICards";
import { ConsultationsByDoctorChart } from "./components/charts/ConsultationsByDoctorChart";
import { BedOccupancyChart } from "./components/charts/BedOccupancyChart";
import { WaitTimeByDepartmentChart } from "./components/charts/WaitTimeByDepartmentChart";
import { TopDiagnosesTable } from "./components/TopDiagnosesTable";
import { StatsErrorState } from "./components/StatsErrorState";
import {
  PageContainer,
  PageTitle,
} from "../../shared/components/ui/PageContainer";
import { Card, CardContent } from "../../shared/components/ui/card";
import { Badge } from "../../shared/components/ui/badge";
import { Spinner } from "../../shared/components/ui/spinner";
import { Skeleton } from "../../shared/components/ui/skeleton";
import { Button } from "../../shared/components/ui/button";
import { cn } from "../../shared/utils/cn";

export const ClinicalDashboard = () => {
  const { t, i18n } = useTranslation("analytics");
  const {
    data: dashboardData,
    isLoading,
    isError,
    isFetching,
    dataUpdatedAt,
    refetch,
  } = useDashboardQuery();

  if (isLoading) {
    return <ClinicalDashboardSkeleton />;
  }

  if (isError || !dashboardData) {
    return (
      <StatsErrorState
        onRetry={() => refetch()}
        titleKey="dashboard.errorTitle"
        descriptionKey="dashboard.errorDescription"
        retryKey="dashboard.retryButton"
      />
    );
  }

  const lastUpdatedLabel = new Intl.DateTimeFormat(i18n.language, {
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(dataUpdatedAt));

  return (
    <PageContainer>
      <PageTitle
        icon={<Activity className="size-5 text-primary animate-pulse-glow" />}
        title={t("dashboard.title")}
        description={t("dashboard.subtitle")}
        actions={
          <div className="flex items-center gap-3">
            <Badge
              variant="outline"
              className="gap-2 rounded-lg bg-muted px-2.5 py-1.5"
            >
              <span
                className="relative flex size-1.5 shrink-0"
                aria-hidden="true"
              >
                {isFetching ? (
                  <span className="absolute inline-flex size-full animate-ping rounded-full bg-warning opacity-75" />
                ) : null}
                <span
                  className={cn(
                    "relative inline-flex size-1.5 rounded-full",
                    isFetching ? "bg-warning" : "bg-success",
                  )}
                />
              </span>
              <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
                {isFetching
                  ? t("dashboard.syncing")
                  : t("dashboard.lastUpdated")}
              </span>
              {isFetching ? null : (
                <span className="text-xs font-black tabular-nums text-foreground">
                  {lastUpdatedLabel}
                </span>
              )}
            </Badge>
            <Button
              variant="outline"
              size="icon"
              onClick={() => refetch()}
              disabled={isFetching}
              aria-label={t("dashboard.refresh")}
            >
              {isFetching ? <Spinner /> : <RefreshCw aria-hidden="true" />}
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
  );
};

const ClinicalDashboardSkeleton = () => {
  return (
    <PageContainer>
      <div>
        <Skeleton className="h-6 w-48" />
        <Skeleton className="mt-2.5 h-3 w-72" />
      </div>

      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6 lg:gap-4">
        {Array.from({ length: 6 }).map((_, indexValue) => (
          <Card key={String(indexValue)} className="gap-0 py-4">
            <CardContent className="flex flex-col gap-3 px-4">
              <div className="flex items-start justify-between gap-2">
                <Skeleton className="h-3 w-20" />
                <Skeleton className="size-8 rounded-lg" />
              </div>
              <Skeleton className="h-8 w-16" />
              <Skeleton className="h-4 w-4/5" />
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 gap-4 md:gap-6 xl:grid-cols-12">
        <Card className="gap-0 py-5 xl:col-span-7">
          <CardContent className="flex flex-col gap-5 px-5">
            <div className="flex items-start justify-between gap-4">
              <div className="flex flex-col gap-2">
                <Skeleton className="h-4 w-44" />
                <Skeleton className="h-3 w-64" />
              </div>
              <Skeleton className="h-11 w-20 rounded-lg" />
            </div>
            <div className="flex flex-col gap-6">
              {Array.from({ length: 5 }).map((_, indexValue) => (
                <div
                  key={String(indexValue)}
                  className="flex items-center gap-3"
                >
                  <Skeleton className="h-7 w-36 shrink-0" />
                  <Skeleton className="h-6 flex-1 rounded-md" />
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card className="gap-0 py-5 xl:col-span-5">
          <CardContent className="flex flex-col gap-5 px-5">
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
                <Skeleton
                  key={String(indexValue)}
                  className="h-16 w-full rounded-xl"
                />
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      <Card className="gap-0 py-5">
        <CardContent className="flex flex-col gap-5 px-5">
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
        </CardContent>
      </Card>

      <Card className="gap-0 py-5">
        <CardContent className="flex flex-col gap-4 px-5">
          <div className="flex items-start justify-between gap-4">
            <div className="flex flex-col gap-2">
              <Skeleton className="h-4 w-36" />
              <Skeleton className="h-3 w-60" />
            </div>
            <Skeleton className="h-11 w-20 rounded-lg" />
          </div>
          <div className="flex flex-col gap-3">
            {Array.from({ length: 5 }).map((_, indexValue) => (
              <Skeleton
                key={String(indexValue)}
                className="h-12 w-full rounded-lg"
              />
            ))}
          </div>
        </CardContent>
      </Card>
    </PageContainer>
  );
};
