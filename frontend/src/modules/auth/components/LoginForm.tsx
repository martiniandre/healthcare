import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { useAuthStore } from "../../../shared/store/auth_store"
import { Card } from "../../../shared/components/ui/Card"
import { Input } from "../../../shared/components/ui/Input"
import { Button } from "../../../shared/components/ui/Button"
import { Alert, AlertDescription } from "../../../shared/components/ui/Alert"
import { getLoginFormSchema, type LoginFormData } from "../auth_schemas"
import { KeyRound, Mail, ShieldAlert } from "lucide-react"
import { useLoginMutation } from "../queries"

export const LoginForm = () => {
  const { t } = useTranslation()
  const loginToStore = useAuthStore((state) => state.login)
  const [generalError, setGeneralError] = useState<string | null>(null)
  const loginMutation = useLoginMutation()

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(getLoginFormSchema(t)),
  })

  const onSubmit = async (formData: LoginFormData) => {
    setGeneralError(null)
    try {
      const response = await loginMutation.mutateAsync({ email: formData.email, password: formData.password })
      loginToStore(
        response.userId,
        response.role,
        response.email,
        response.fullName,
        response.isActive,
      )
    } catch (loginRequestError) {
      if (loginRequestError instanceof Error) {
        setGeneralError(loginRequestError.message)
      } else {
        setGeneralError(t("auth.defaultError"))
      }
    }
  }

  return (
    <Card glowingType="cyan" className="p-8">
      <h2 className="text-lg font-bold text-gray-800 mb-6">{t("auth.authTitle")}</h2>

      {generalError && (
        <Alert variant="destructive" className="mb-6">
          <ShieldAlert className="h-4 w-4" />
          <AlertDescription>
            {generalError}
          </AlertDescription>
        </Alert>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
        <div className="flex flex-col gap-1 text-left">
          <label className="text-xs font-semibold text-gray-600 flex items-center gap-1.5 mb-1">
            <Mail className="w-3.5 h-3.5 text-primary" />
            {t("auth.emailLabel")}
          </label>
          <Input
            type="email"
            placeholder={t("auth.emailPlaceholder")}
            autoComplete="email"
            maxLength={255}
            errorText={errors.email?.message}
            {...register("email")}
          />
        </div>

        <div className="flex flex-col gap-1 text-left">
          <label className="text-xs font-semibold text-gray-600 flex items-center gap-1.5 mb-1">
            <KeyRound className="w-3.5 h-3.5 text-primary" />
            {t("auth.passwordLabel")}
          </label>
          <Input
            type="password"
            placeholder={t("auth.passwordPlaceholder")}
            autoComplete="current-password"
            errorText={errors.password?.message}
            {...register("password")}
          />
        </div>

        <Button
          type="submit"
          disabled={loginMutation.isPending}
          className="w-full py-3.5 mt-2 text-sm font-bold tracking-wide uppercase"
        >
          {loginMutation.isPending ? t("auth.loadingText") : t("auth.submitText")}
        </Button>
      </form>
    </Card>
  )
}
