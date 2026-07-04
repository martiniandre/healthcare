import { useSearchParams } from "react-router-dom"
import { usePortalDashboardQuery } from "./queries"
import { PortalDashboardOverview } from "./PortalDashboardOverview"
import { PortalEncounters } from "./PortalEncounters"
import { PortalObservations } from "./PortalObservations"
import { PortalConditions } from "./PortalConditions"
import { PortalMedications } from "./PortalMedications"
import { PortalReports } from "./PortalReports"
import { PortalImaging } from "./PortalImaging"
import {
  LayoutDashboard,
  History,
  Heart,
  Activity,
  Pill,
  FileText,
  Image,
  Loader2,
} from "lucide-react"

const PortalTab = {
  Dashboard: "dashboard",
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
          <span className="text-sm text-gray-500 font-medium">Carregando portal...</span>
        </div>
      </div>
    )
  }

  const patientName = dashboardData.patient_info.full_name

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full">
      <div className="bg-white border border-border p-4 sm:p-6 rounded-xl">
        <div className="flex items-center gap-4">
          <div className="w-14 h-14 rounded-full bg-primary/10 flex items-center justify-center">
            <span className="text-xl font-bold text-primary">
              {patientName.split(" ").map((n: string) => n[0]).join("").slice(0, 2).toUpperCase()}
            </span>
          </div>
          <div>
            <h1 className="text-xl font-bold text-gray-900">{patientName}</h1>
            <p className="text-sm text-gray-500">Portal do Paciente</p>
          </div>
        </div>
      </div>

      <div className="flex flex-col md:flex-row gap-6 items-start">
        <div className="w-full md:w-56 shrink-0 bg-white border border-border p-4 rounded-xl flex flex-col gap-4">
          <span className="text-[10px] font-black text-gray-500 uppercase tracking-widest px-3 text-left">
            Navegação
          </span>
          <div className="flex flex-col gap-2">
            {sidebarItems.map((item) => (
              <button
                key={item.key}
                onClick={() => setActiveTab(item.key)}
                className={`w-full text-left flex items-center gap-3 px-4 py-3 rounded-lg text-xs font-extrabold transition-all duration-300 ${
                  activeTab === item.key
                    ? "bg-primary/8 text-primary"
                    : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
                }`}
              >
                {item.icon}
                {item.label}
              </button>
            ))}
          </div>
        </div>

        <div className="flex-1 flex flex-col gap-6 min-w-0 w-full">
          {activeTab === PortalTab.Dashboard && (
            <PortalDashboardOverview dashboard={dashboardData} />
          )}
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
