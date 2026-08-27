import { useMemo } from "react"
import { useNavigate, useLocation } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { LogOut, X } from "lucide-react"
import { useAuthStore } from "../store/auth_store"
import { useLayoutStore } from "../store/layout_store"
import { getVisibleNavigationGroups, isNavigationItemActive } from "../navigation/navigationFilter"

export const AppSidebar = () => {
  const { t } = useTranslation("sidebar")
  const navigate = useNavigate()
  const location = useLocation()
  const { email, logout, role } = useAuthStore()
  const { isMobileSidebarOpen, closeMobileSidebar } = useLayoutStore()

  const visibleNavigationGroups = useMemo(() => getVisibleNavigationGroups(role), [role])

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
            <div className="relative w-9 h-9 rounded-xl bg-gradient-to-br from-primary to-accent flex items-center justify-center overflow-hidden">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2.2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className="w-5 h-5 text-white animate-ecg-trace"
              >
                <polyline
                  points="1,12 5,12 7.5,6 10.5,18 13,12 16,12 23,12"
                  style={{ strokeDasharray: "20 6" }}
                />
              </svg>
            </div>
            <div>
              <h1 className="text-sm font-display font-bold tracking-tight text-gray-900 leading-none">
                {t("title")}
              </h1>
              <span className="text-[10px] text-muted font-medium">{t("subtitle")}</span>
            </div>
          </div>
          <button
            onClick={closeMobileSidebar}
            aria-label={t("closeMenu")}
            className="p-1.5 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-50 md:hidden"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        <div className="h-px bg-border mx-4" />

        <nav className="flex-1 min-h-0 overflow-y-auto px-3 py-4 flex flex-col gap-5">
          {visibleNavigationGroups.map((navigationGroup, groupIndex) => (
            <section key={navigationGroup.key} className="flex flex-col gap-1">
              <div className="px-3 mb-1.5 flex items-center gap-2">
                <span className="font-mono text-[10px] font-semibold text-primary/70 tabular-nums">
                  {String(groupIndex + 1).padStart(2, "0")}
                </span>
                <span className="text-[9px] font-black uppercase tracking-[0.18em] text-gray-400">
                  {t(navigationGroup.key)}
                </span>
                <span className="flex-1 h-px bg-gradient-to-r from-border to-transparent" />
              </div>

              {navigationGroup.items.map((navigationItem) => {
                const isCurrentlyActive = isNavigationItemActive(navigationItem, location.pathname)
                const isInteractive = !navigationItem.disabled

                return (
                  <button
                    key={navigationItem.path}
                    onClick={() => {
                      if (isInteractive) {
                        navigate(navigationItem.path)
                        closeMobileSidebar()
                      }
                    }}
                    disabled={!isInteractive}
                    aria-current={isInteractive && isCurrentlyActive ? "page" : undefined}
                    className={`group relative w-full text-left flex items-center gap-3 px-3 py-2 rounded-lg text-[13px] transition-all duration-200 ${
                      isInteractive
                        ? isCurrentlyActive
                          ? "bg-primary/10 text-primary font-semibold"
                          : "text-gray-500 hover:text-gray-900 hover:bg-gray-100/80"
                        : "text-gray-300 cursor-not-allowed"
                    }`}
                  >
                    {isInteractive && isCurrentlyActive && (
                      <span className="absolute left-0 top-1/2 -translate-y-1/2 w-[3px] h-4 rounded-full bg-primary" />
                    )}
                    <navigationItem.icon
                      className={`w-[18px] h-[18px] shrink-0 transition-colors duration-200 ${
                        isInteractive
                          ? isCurrentlyActive
                            ? "text-primary"
                            : "text-gray-400 group-hover:text-gray-600"
                          : "text-gray-300"
                      }`}
                    />
                    <span className="truncate">{t(navigationItem.key)}</span>
                    {navigationItem.disabled && (
                      <span className="ml-auto shrink-0 text-[8px] bg-gray-100 text-gray-400 px-1.5 py-0.5 rounded font-bold uppercase">
                        {t("comingSoon")}
                      </span>
                    )}
                    {navigationItem.adminOnly && !navigationItem.disabled && (
                      <span className="ml-auto shrink-0 text-[8px] bg-secondary/10 text-secondary px-1.5 py-0.5 rounded font-bold uppercase tracking-wider">
                        {t("restricted")}
                      </span>
                    )}
                  </button>
                )
              })}
            </section>
          ))}
        </nav>

        <div className="px-3 pb-3 pt-1">
          <button
            onClick={() => {
              logout()
              closeMobileSidebar()
            }}
            className="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-[13px] font-semibold text-red-400 hover:text-red-600 hover:bg-red-50 transition-all duration-200"
          >
            <LogOut className="w-[18px] h-[18px] shrink-0" />
            {t("logout")}
          </button>
        </div>

        <div className="px-5 py-3 border-t border-border flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-1.5 h-1.5 rounded-full bg-success animate-pulse-glow" />
            <span className="text-[10px] text-muted font-medium">
              FHIR R4 · gRPC-Web
            </span>
          </div>
          <span className="text-[9px] text-gray-300 font-mono truncate">
            {email ? email.split("@")[0] : ""}
          </span>
        </div>
      </aside>
    </>
  )
}