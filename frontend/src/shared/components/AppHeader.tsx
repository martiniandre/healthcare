import { useState, useEffect } from "react"
import { useAuthStore } from "../store/auth_store"
import { useLayoutStore } from "../store/layout_store"
import { Menu, Bell, ShieldCheck, WifiOff } from "lucide-react"

export const AppHeader = () => {
  const { role, email } = useAuthStore()
  const toggleMobileSidebar = useLayoutStore((state) => state.toggleMobileSidebar)
  const [isOnline, setIsOnline] = useState(navigator.onLine)

  useEffect(() => {
    const handleOnline = () => setIsOnline(true)
    const handleOffline = () => setIsOnline(false)

    window.addEventListener("online", handleOnline)
    window.addEventListener("offline", handleOffline)

    return () => {
      window.removeEventListener("online", handleOnline)
      window.removeEventListener("offline", handleOffline)
    }
  }, [])

  const translateRole = (userRole: string | null) => {
    switch (userRole) {
      case "RoleAdmin":
        return "Administrador"
      case "RoleDoctor":
        return "Médico"
      case "RoleNurse":
        return "Enfermeiro"
      case "RoleReception":
        return "Recepção"
      default:
        return "Profissional"
    }
  }

  return (
    <header className="w-full border-b border-border bg-white/80 backdrop-blur-sm sticky top-0 z-50 px-4 md:px-6 py-2.5 flex items-center justify-end gap-3">
      <button
        onClick={toggleMobileSidebar}
        className="mr-auto p-2 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-50 transition-colors md:hidden"
      >
        <Menu className="w-5 h-5" />
      </button>

      {!isOnline && (
        <div className="flex items-center gap-1.5 px-3 py-1 rounded-full bg-red-50 border border-red-200 text-red-600 text-[10px] font-bold animate-pulse select-none mr-2">
          <WifiOff className="w-3.5 h-3.5 text-red-500" />
          <span>Sem Conexão</span>
        </div>
      )}

      <button className="relative p-2 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-50 transition-colors">
        <Bell className="w-4 h-4" />
        <span className="absolute top-1.5 right-1.5 w-1.5 h-1.5 bg-primary rounded-full" />
      </button>

      <div className="h-5 w-px bg-border" />

      <div className="flex items-center gap-2.5">
        <div className="w-8 h-8 rounded-lg bg-primary/8 flex items-center justify-center text-primary text-xs font-black">
          {email ? email.charAt(0).toUpperCase() : "U"}
        </div>
        <div className="hidden sm:flex flex-col items-start">
          <span className="text-xs font-semibold text-gray-800 leading-tight">
            {email || "usuario@hospital.com"}
          </span>
          <div className="flex items-center gap-1">
            <ShieldCheck className="w-3 h-3 text-secondary" />
            <span className="text-[10px] text-muted font-medium">
              {translateRole(role)}
            </span>
          </div>
        </div>
      </div>
    </header>
  )
}
