import { useTranslation } from "react-i18next"
import { Leaf } from "lucide-react"
import { LoginForm } from "./components/LoginForm"

export const Login = () => {
  const { t } = useTranslation("auth")

  return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center p-6 select-none relative overflow-hidden">

      <div className="w-full max-w-[420px] z-10">
        <div className="flex flex-col items-center gap-2 mb-8 text-center">
          <div className="flex size-12 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <Leaf className="size-6" aria-hidden="true" />
          </div>
          <h1 className="m-0 font-display text-3xl font-semibold tracking-tight text-foreground">
            HealthCare
          </h1>
          <p className="text-sm text-muted-foreground max-w-[280px]">
            {t("portalSubtitle")}
          </p>
        </div>

        <LoginForm />

        <p className="text-center text-2xs text-muted-foreground mt-8 leading-relaxed">
          {t("footerNotice")}
        </p>
      </div>
    </div>
  )
}
