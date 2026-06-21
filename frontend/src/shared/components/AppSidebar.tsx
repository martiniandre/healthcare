import { useNavigate, useLocation } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { useAuthStore } from "../store/auth_store"
import { useLayoutStore } from "../store/layout_store"
import { Activity, Users, BarChart3, Settings, LogOut, X, Sparkles, History } from "lucide-react"

const navigationItems = [
  { key: "patients", icon: Users, path: "/" },
  { key: "telemetry", icon: Activity, path: "/telemetry" },
  { key: "examAnalyzer", icon: Sparkles, path: "/exam-analyzer" },
  { key: "stats", icon: BarChart3, path: "/stats" },
  { key: "staffManagement", icon: Users, path: "/staff" },
  { key: "auditLogs", icon: History, path: "/audit-logs", adminOnly: true },
  { key: "settings", icon: Settings, path: "/settings", disabled: true },
]

export const AppSidebar = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const location = useLocation()
  const { email, logout, role } = useAuthStore()
  const { isMobileSidebarOpen, closeMobileSidebar } = useLayoutStore()

  return (
    <>
      {isMobileSidebarOpen && (
        <div
          onClick={closeMobileSidebar}
          className="fixed inset-0 z-40 bg-black/40 backdrop-blur-[1px] md:hidden transition-opacity duration-300"
        />
      )}

      <aside
        className={`w-[240px] shrink-0 h-screen fixed md:sticky top-0 left-0 bg-white border-r border-border flex flex-col z-50 transition-transform duration-300 md:translate-x-0 ${
          isMobileSidebarOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <div className="px-5 py-5 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="bg-primary/8 p-2.5 rounded-xl border border-primary/10">
              <Activity className="w-5 h-5 text-primary" />
            </div>
            <div>
              <h1 className="text-sm font-black tracking-tight text-gray-900 leading-none">
                {t("sidebar.title")}
              </h1>
              <span className="text-[10px] text-muted font-medium">{t("sidebar.subtitle")}</span>
            </div>
          </div>
          <button
            onClick={closeMobileSidebar}
            className="p-1.5 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-50 md:hidden"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        <div className="h-px bg-border mx-4" />

        <nav className="flex-1 px-3 py-5 flex flex-col gap-1">
          <span className="text-[9px] font-black text-muted/60 uppercase tracking-[0.15em] px-3 mb-3">
            {t("sidebar.menuHeader")}
          </span>
          {navigationItems
            .filter((item) => !item.adminOnly || role === "ADMIN")
            .map((item) => {
              const isActive = location.pathname === item.path ||
                (item.path !== "/" && location.pathname.startsWith(item.path))
              const isHomeActive = item.path === "/" && (
                location.pathname === "/" || location.pathname.startsWith("/patients")
              )
              const isCurrentlyActive = isActive || isHomeActive

              return (
                <button
                  key={item.path}
                  onClick={() => {
                    if (!item.disabled) {
                      navigate(item.path)
                      closeMobileSidebar()
                    }
                  }}
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
                {t(`sidebar.${item.key}`)}
                {item.disabled && (
                  <span className="ml-auto text-[8px] bg-gray-100 text-gray-400 px-1.5 py-0.5 rounded font-bold uppercase">
                    {t("sidebar.comingSoon")}
                  </span>
                )}
              </button>
            )
          })}
        </nav>

        <div className="px-3 pb-3">
          <button
            onClick={() => {
              logout()
              closeMobileSidebar()
            }}
            className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-[13px] font-semibold text-red-400 hover:text-red-600 hover:bg-red-50 transition-all duration-200"
          >
            <LogOut className="w-[18px] h-[18px] shrink-0" />
            {t("sidebar.logout")}
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
    </>
  )
}
