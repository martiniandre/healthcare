import { useSearchParams } from "react-router-dom"
import { Button } from "../../shared/components/ui/Button"
import { usePortalDashboardQuery } from "./queries"
import { PortalDashboardOverview } from "./PortalDashboardOverview"
import { PortalEncounters } from "./PortalEncounters"
import { PortalObservations } from "./PortalObservations"
import { PortalConditions } from "./PortalConditions"
import { PortalMedications } from "./PortalMedications"
import { PortalReports } from "./PortalReports"
import { PortalImaging } from "./PortalImaging"
import { PortalAppointments } from "./PortalAppointments"
import {
  LayoutDashboard,
  History,
  Heart,
  Activity,
  Pill,
  FileText,
  Image,
  CalendarClock,
  Loader2,
} from "lucide-react"

const PortalTab = {
  Dashboard: "dashboard",
  Appointments: "appointments",
  Encounters: "encounters",
  Observations: "observations",
  Conditions: "conditions",
  Medications: "medications",
  Reports: "reports",
  Imaging: "imaging",
} as const

type PortalTab = (typeof PortalTab)[keyof typeof PortalTab]

const sidebarItems: { key: PortalTab; label: string; icon: React.ReactNode }[] = [
  { key: "dashboard", label: "Visão Geral", icon: <LayoutDashboard className="w-4 h-4 shrink-0" /> },
  { key: "appointments", label: "Agendamentos", icon: <CalendarClock className="w-4 h-4 shrink-0" /> },
  { key: "encounters", label: "Consultas", icon: <History className="w-4 h-4 shrink-0" /> },
  { key: "observations", label: "Sinais Vitais", icon: <Heart className="w-4 h-4 shrink-0" /> },
  { key: "conditions", label: "Condições", icon: <Activity className="w-4 h-4 shrink-0" /> },
  { key: "medications", label: "Medicamentos", icon: <Pill className="w-4 h-4 shrink-0" /> },
  { key: "reports", label: "Exames", icon: <FileText className="w-4 h-4 shrink-0" /> },
  { key: "imaging", label: "Imagens", icon: <Image className="w-4 h-4 shrink-0" /> },
]

export const PortalPage = () => {
  const [searchParameters, setSearchParameters] = useSearchParams()
  const activeTab = (searchParameters.get("tab") || PortalTab.Dashboard) as PortalTab
  const setActiveTab = (tabName: PortalTab) => {
    setSearchParameters({ tab: tabName })
  }

  const { data: dashboardData, isLoading: isDashboardLoading } = usePortalDashboardQuery()

  if (isDashboardLoading || !dashboardData) {
    return (
      <div className="flex-1 p-4 sm:p-6 md:p-8 flex items-center justify-center">
        <div className="flex flex-col items-center gap-2">
          <Loader2 className="w-8 h-8 text-primary animate-spin" />
          <span className="text-sm text-muted-foreground font-medium">Carregando portal...</span>
        </div>
      </div>
    )
  }

  const patientName = dashboardData.patient_info.full_name

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full">
      <div className="bg-surface border border-border p-4 sm:p-6 rounded-xl">
        <div className="flex items-center gap-4">
          <div className="w-14 h-14 rounded-full bg-primary-soft flex items-center justify-center">
            <span className="text-xl font-bold text-primary">
              {patientName.split(" ").map((n: string) => n[0]).join("").slice(0, 2).toUpperCase()}
            </span>
          </div>
          <div>
            <h1 className="text-xl font-display font-bold text-foreground">{patientName}</h1>
            <p className="text-sm text-muted-foreground">Portal do Paciente</p>
          </div>
        </div>
      </div>

      <div className="flex flex-col md:flex-row gap-6 items-start">
        <div className="w-full md:w-56 shrink-0 bg-surface border border-border p-4 rounded-xl flex flex-col gap-4">
          <span className="text-[10px] font-black text-muted-foreground uppercase tracking-widest px-3 text-left">
            Navegação
          </span>
          <div className="flex flex-col gap-2">
            {sidebarItems.map((item) => (
              <Button
                key={item.key}
                onClick={() => setActiveTab(item.key)}
                variantType={activeTab === item.key ? "primary" : "ghost"}
                className={`w-full justify-start gap-3 px-4 py-3 text-xs font-extrabold transition-all duration-300 h-auto shadow-none ${
                  activeTab === item.key ? "" : "text-muted-foreground hover:text-foreground"
                }`}
              >
                {item.icon}
                {item.label}
              </Button>
            ))}
          </div>
        </div>

        <div className="flex-1 flex flex-col gap-6 min-w-0 w-full">
          {activeTab === PortalTab.Dashboard && (
            <PortalDashboardOverview dashboard={dashboardData} />
          )}
          {activeTab === PortalTab.Appointments && <PortalAppointments />}
          {activeTab === PortalTab.Encounters && <PortalEncounters />}
          {activeTab === PortalTab.Observations && <PortalObservations />}
          {activeTab === PortalTab.Conditions && <PortalConditions />}
          {activeTab === PortalTab.Medications && <PortalMedications />}
          {activeTab === PortalTab.Reports && <PortalReports />}
          {activeTab === PortalTab.Imaging && <PortalImaging />}
        </div>
      </div>
    </div>
  )
}
