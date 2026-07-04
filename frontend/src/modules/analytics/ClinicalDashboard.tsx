import { useDashboardQuery } from "./dashboard_queries"
import { DashboardKPICards } from "./components/DashboardKPICards"
import { DoctorConsultationsChart } from "./components/DoctorConsultationsChart"
import { OccupancyGauge } from "./components/OccupancyGauge"
import { TopDiagnosesTable } from "./components/TopDiagnosesTable"
import { WaitTimeChart } from "./components/WaitTimeChart"
import { Card } from "../../shared/components/ui/Card"
import { Skeleton } from "../../shared/components/ui/Skeleton"
import { Activity } from "lucide-react"
import { Button } from "../../shared/components/ui/Button"

export const ClinicalDashboard = () => {
  const { data: dashboardData, isLoading, isError } = useDashboardQuery()

  if (isLoading) {
    return <ClinicalDashboardSkeleton />
  }

  if (isError || !dashboardData) {
    return <ClinicalDashboardError />
  }

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full select-none">
      <div className="text-left">
        <h2 className="text-xl font-black text-gray-900 leading-none flex items-center gap-2">
          <Activity className="w-5 h-5 text-primary animate-pulse-glow" />
          Dashboard Clínico
        </h2>
        <span className="text-xs text-muted mt-1.5 block">
          Indicadores em tempo real para acompanhamento clínico
        </span>
      </div>

      <DashboardKPICards
        consultationsToday={dashboardData.consultations_today}
        consultationsTrend={dashboardData.consultations_trend}
        occupancyRate={dashboardData.occupancy_rate}
        avgWaitTimeMinutes={dashboardData.avg_wait_time_minutes}
        activePatients={dashboardData.active_patients}
        examsToday={dashboardData.exams_today}
        newDiagnosesToday={dashboardData.new_diagnoses_today}
      />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <DoctorConsultationsChart
          consultationsPerDoctor={dashboardData.consultations_per_doctor}
        />
        <OccupancyGauge
          occupancyRate={dashboardData.occupancy_rate}
          totalBeds={dashboardData.occupancy_total_beds}
          occupiedBeds={dashboardData.occupancy_occupied_beds}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <WaitTimeChart
          waitTimeByDepartment={dashboardData.wait_time_by_department}
        />
        <div>
          {/* Empty slot for future widget */}
        </div>
      </div>

      <TopDiagnosesTable
        topDiagnoses={dashboardData.top_diagnoses}
      />
    </div>
  )
}

const ClinicalDashboardSkeleton = () => {
  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full">
      <div className="text-left">
        <Skeleton className="h-6 w-48" />
        <Skeleton className="h-3 w-72 mt-2.5" />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4">
        {Array.from({ length: 6 }).map((_, indexValue) => (
          <Card key={String(indexValue)} className="p-4 flex items-center justify-between border border-border h-24">
            <div className="flex-1 flex flex-col gap-2">
              <Skeleton className="h-3 w-20" />
              <Skeleton className="h-6 w-16" />
              <Skeleton className="h-3 w-28" />
            </div>
            <Skeleton className="w-12 h-12 rounded-xl" />
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card className="p-5 flex flex-col gap-5 border border-border h-80">
          <Skeleton className="h-4 w-40" />
          <Skeleton className="h-3 w-56" />
          <div className="flex-1 flex flex-col gap-4 mt-4">
            {Array.from({ length: 4 }).map((_, indexValue) => (
              <div key={String(indexValue)} className="flex items-center gap-3">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="flex-1 h-5 rounded-full" />
                <Skeleton className="h-4 w-8" />
              </div>
            ))}
          </div>
        </Card>

        <Card className="p-5 flex flex-col gap-5 border border-border h-80">
          <Skeleton className="h-4 w-36" />
          <Skeleton className="h-3 w-52" />
          <div className="flex-1 flex items-center justify-center">
            <Skeleton className="w-40 h-28" />
          </div>
        </Card>
      </div>

      <Card className="p-5 flex flex-col gap-4 border border-border">
        <Skeleton className="h-4 w-44" />
        <Skeleton className="h-3 w-60" />
        <div className="flex flex-col gap-3 mt-2">
          {Array.from({ length: 5 }).map((_, indexValue) => (
            <Skeleton key={String(indexValue)} className="h-10 w-full" />
          ))}
        </div>
      </Card>
    </div>
  )
}

const ClinicalDashboardError = () => {
  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col items-center justify-center gap-4 max-w-7xl mx-auto w-full select-none">
      <div className="text-center p-8 bg-white border border-red-100 shadow-xl rounded-2xl max-w-md w-full flex flex-col items-center gap-4">
        <div className="bg-red-50 p-4 rounded-full">
          <Activity className="w-10 h-10 text-red-500 animate-bounce" />
        </div>
        <h3 className="text-lg font-black text-gray-900">Erro ao carregar dashboard</h3>
        <p className="text-xs text-gray-500 leading-relaxed">
          Não foi possível carregar os dados do dashboard clínico. Tente novamente mais tarde.
        </p>
        <Button
          onClick={() => window.location.reload()}
          className="w-full bg-red-600 hover:bg-red-700 text-white font-bold py-2 rounded-xl transition-all duration-200 mt-2"
        >
          Tentar Novamente
        </Button>
      </div>
    </div>
  )
}
