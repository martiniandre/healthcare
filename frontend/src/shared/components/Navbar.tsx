import { useAuthStore } from "../store/auth_store"
import { Button } from "./ui/Button"
import { Activity, LogOut, ShieldCheck, User } from "lucide-react"

export const Navbar = () => {
  const { role, email, logout } = useAuthStore()

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
    <header className="w-full border-b border-border bg-slate-950/80 backdrop-blur-md sticky top-0 z-50 px-6 py-4 flex items-center justify-between">
      <div className="flex items-center gap-3">
        <div className="bg-primary/10 p-2 rounded-lg border border-primary/20 animate-pulse-glow">
          <Activity className="w-6 h-6 text-primary" />
        </div>
        <div>
          <h1 className="text-lg font-bold tracking-tight text-white m-0 leading-none">
            HealthCare
          </h1>
          <span className="text-xs text-muted">Console Clínico & PACS</span>
        </div>
      </div>

      <div className="flex items-center gap-4">
        <div className="flex items-center gap-2.5 bg-slate-900/60 border border-border px-3 py-1.5 rounded-lg">
          <div className="bg-primary/10 p-1.5 rounded-full">
            <User className="w-4 h-4 text-primary" />
          </div>
          <div className="flex flex-col items-start">
            <span className="text-xs font-semibold text-gray-200">
              {email || "usuario@hospital.com"}
            </span>
            <div className="flex items-center gap-1 mt-0.5">
              <ShieldCheck className="w-3.5 h-3.5 text-secondary" />
              <span className="text-[10px] text-gray-400 font-bold uppercase tracking-wider">
                {translateRole(role)}
              </span>
            </div>
          </div>
        </div>

        <Button variantType="danger" onClick={logout} className="px-3 py-2 text-xs">
          <LogOut className="w-4 h-4" />
          Sair
        </Button>
      </div>
    </header>
  )
}
