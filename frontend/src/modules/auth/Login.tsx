import { useTranslation } from "react-i18next"
import { Activity } from "lucide-react"
import { LoginForm } from "./components/LoginForm"

export const Login = () => {
  const { t } = useTranslation()

  return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center p-6 select-none relative overflow-hidden">
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/5 rounded-full filter blur-[120px] animate-pulse-glow" />
      <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-secondary/5 rounded-full filter blur-[120px]" />

      <div className="w-full max-w-[420px] z-10">
        <div className="flex flex-col items-center gap-2 mb-8 text-center">
          <div className="bg-primary/10 p-3.5 rounded-2xl border border-primary/20 animate-pulse-glow">
            <Activity className="w-8 h-8 text-primary" />
          </div>
          <h1 className="text-3xl font-extrabold tracking-tight text-gray-900 m-0">
            HealthCare
          </h1>
          <p className="text-sm text-muted max-w-[280px]">
            {t("auth.portalSubtitle")}
          </p>
        </div>

        <LoginForm />

        <p className="text-center text-[10px] text-gray-500 mt-8 leading-relaxed">
          {t("auth.footerNotice")}
        </p>
      </div>
    </div>
  )
}
