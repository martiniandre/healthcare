import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useAuthStore } from "../../shared/store/auth_store"
import { Card } from "../../shared/components/ui/Card"
import { Input } from "../../shared/components/ui/Input"
import { Button } from "../../shared/components/ui/Button"
import { loginFormSchema, type LoginFormData } from "./auth_schemas"
import { Activity, KeyRound, Mail, ShieldAlert } from "lucide-react"
import { authApi } from "../../shared/services/auth_api"

export const Login = () => {
  const loginToStore = useAuthStore((state) => state.login)
  const [generalError, setGeneralError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginFormSchema),
  })

  const onSubmit = async (formData: LoginFormData) => {
    setIsSubmitting(true)
    setGeneralError(null)
    try {
      const authResponseData = await authApi.login(formData.email, formData.password)
      loginToStore(
        authResponseData.userId,
        authResponseData.role,
        authResponseData.email
      )
    } catch (loginRequestError) {
      if (loginRequestError instanceof Error) {
        setGeneralError(loginRequestError.message)
      } else {
        setGeneralError("Erro ao estabelecer conexão com o servidor.")
      }
    } finally {
      setIsSubmitting(false)
    }
  }

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
            Portal de Acesso Médico Integrado ao Barramento FHIR
          </p>
        </div>

        <Card glowingType="cyan" className="p-8">
          <h2 className="text-lg font-bold text-gray-800 mb-6">Autenticação Clínica</h2>

          {generalError && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-3.5 flex gap-2.5 items-start mb-6 text-left">
              <ShieldAlert className="w-5 h-5 text-red-500 shrink-0 mt-0.5" />
              <span className="text-xs text-red-600 leading-relaxed font-medium">
                {generalError}
              </span>
            </div>
          )}

          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
            <div className="flex flex-col gap-1 text-left">
              <label className="text-xs font-semibold text-gray-600 flex items-center gap-1.5 mb-1">
                <Mail className="w-3.5 h-3.5 text-primary" />
                E-mail Profissional
              </label>
              <Input
                type="email"
                placeholder="nome.sobrenome@hospital.com"
                errorText={errors.email?.message}
                {...register("email")}
              />
            </div>

            <div className="flex flex-col gap-1 text-left">
              <label className="text-xs font-semibold text-gray-600 flex items-center gap-1.5 mb-1">
                <KeyRound className="w-3.5 h-3.5 text-primary" />
                Senha de Acesso
              </label>
              <Input
                type="password"
                placeholder="••••••••"
                errorText={errors.password?.message}
                {...register("password")}
              />
            </div>

            <Button
              type="submit"
              disabled={isSubmitting}
              className="w-full py-3.5 mt-2 text-sm font-bold tracking-wide uppercase"
            >
              {isSubmitting ? "Autenticando via gRPC..." : "Entrar no Console"}
            </Button>
          </form>
        </Card>

        <p className="text-center text-[10px] text-gray-500 mt-8 leading-relaxed">
          Uso restrito a funcionários autorizados do ecossistema hospitalar.<br />
          Todas as conexões e acessos são auditados conforme diretrizes da LGPD/HIPAA.
        </p>
      </div>
    </div>
  )
}
