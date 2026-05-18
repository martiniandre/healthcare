import { useNavigate, useLocation } from "react-router-dom"
import { useAuthStore } from "../store/auth_store"
import { Activity, Users, Image as ImageIcon, BarChart3, Settings, LogOut } from "lucide-react"

const navigationItems = [
  { label: "Pacientes", icon: Users, path: "/" },
  { label: "Telemetria UTI", icon: Activity, path: "/telemetry" },
  { label: "PACS Viewer", icon: ImageIcon, path: "/imaging" },
  { label: "Estatísticas", icon: BarChart3, path: "/stats" },
  { label: "Gestão de Equipes", icon: Users, path: "/staff" },
  { label: "Configurações", icon: Settings, path: "/settings", disabled: true },
]

export const AppSidebar = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { email, logout } = useAuthStore()

  return (
    <aside className="w-[240px] shrink-0 h-screen sticky top-0 bg-white border-r border-border flex flex-col">
      <div className="px-5 py-5 flex items-center gap-3">
        <div className="bg-primary/8 p-2.5 rounded-xl border border-primary/10">
          <Activity className="w-5 h-5 text-primary" />
        </div>
        <div>
          <h1 className="text-sm font-black tracking-tight text-gray-900 leading-none">
            HealthCare
          </h1>
          <span className="text-[10px] text-muted font-medium">Console Clínico v1.0</span>
        </div>
      </div>

      <div className="h-px bg-border mx-4" />

      <nav className="flex-1 px-3 py-5 flex flex-col gap-1">
        <span className="text-[9px] font-black text-muted/60 uppercase tracking-[0.15em] px-3 mb-3">
          Menu Principal
        </span>
        {navigationItems.map((item) => {
          const isActive = location.pathname === item.path ||
            (item.path !== "/" && location.pathname.startsWith(item.path))
          const isHomeActive = item.path === "/" && (
            location.pathname === "/" || location.pathname.startsWith("/patients")
          )
          const isCurrentlyActive = isActive || isHomeActive

          return (
            <button
              key={item.path}
              onClick={() => !item.disabled && navigate(item.path)}
              disabled={item.disabled}
              className={`w-full text-left flex items-center gap-3 px-3 py-2.5 rounded-lg text-[13px] font-semibold transition-all duration-200 ${
                item.disabled
                  ? "text-gray-300 cursor-not-allowed"
                  : isCurrentlyActive
                    ? "bg-primary/8 text-primary"
                    : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
              }`}
            >
              <item.icon className="w-[18px] h-[18px] shrink-0" />
              {item.label}
              {item.disabled && (
                <span className="ml-auto text-[8px] bg-gray-100 text-gray-400 px-1.5 py-0.5 rounded font-bold uppercase">
                  Em breve
                </span>
              )}
            </button>
          )
        })}
      </nav>

      <div className="px-3 pb-3">
        <button
          onClick={logout}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-[13px] font-semibold text-red-400 hover:text-red-600 hover:bg-red-50 transition-all duration-200"
        >
          <LogOut className="w-[18px] h-[18px] shrink-0" />
          Sair
        </button>
      </div>

      <div className="px-5 py-3 border-t border-border flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="w-1.5 h-1.5 rounded-full bg-success animate-pulse-glow" />
          <span className="text-[10px] text-muted font-medium">
            FHIR R4 · gRPC-Web
          </span>
        </div>
        <span className="text-[9px] text-gray-300 font-mono">
          {email ? email.split("@")[0] : ""}
        </span>
      </div>
    </aside>
  )
}
