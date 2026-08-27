import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { useAuthStore } from "../../../shared/store/auth_store"
import { Card } from "../../../shared/components/ui/Card"
import { Input } from "../../../shared/components/ui/Input"
import { Button } from "../../../shared/components/ui/Button"
import { Label } from "../../../shared/components/ui/Label"
import { Alert, AlertDescription } from "../../../shared/components/ui/Alert"
import { getLoginFormSchema, type LoginFormData } from "../auth_schemas"
import { Eye, EyeOff, KeyRound, Mail, ShieldAlert } from "lucide-react"
import { useLoginMutation } from "../queries"

export const LoginForm = () => {
  const { t } = useTranslation("auth")
  const loginToStore = useAuthStore((state) => state.login)
  const [generalError, setGeneralError] = useState<string | null>(null)
  const [isPasswordVisible, setIsPasswordVisible] = useState(false)
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
        setGeneralError(t("defaultError"))
      }
    }
  }

  return (
    <Card glowingType="cyan" className="p-8">
      <h2 className="text-lg font-bold text-foreground mb-6">{t("authTitle")}</h2>

      {generalError && (
        <Alert variant="destructive" className="mb-6">
          <ShieldAlert className="h-4 w-4" />
          <AlertDescription>
            {generalError}
          </AlertDescription>
        </Alert>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5" noValidate>
        <div className="flex flex-col gap-1 text-left">
          <Label className="text-xs font-semibold text-muted-foreground flex items-center gap-1.5 mb-1">
            <Mail className="w-3.5 h-3.5 text-primary" />
            {t("emailLabel")}
          </Label>
          <Input
            type="email"
            placeholder={t("emailPlaceholder")}
            autoComplete="email"
            maxLength={255}
            errorText={errors.email?.message}
            {...register("email")}
          />
        </div>

        <div className="flex flex-col gap-1 text-left">
          <Label className="text-xs font-semibold text-muted-foreground flex items-center gap-1.5 mb-1">
            <KeyRound className="w-3.5 h-3.5 text-primary" />
            {t("passwordLabel")}
          </Label>
          <div className="relative">
            <Input
              type={isPasswordVisible ? "text" : "password"}
              placeholder={t("passwordPlaceholder")}
              autoComplete="current-password"
              errorText={errors.password?.message}
              className="pr-10"
              {...register("password")}
            />
            <Button
              type="button"
              variantType="ghost"
              size="sm"
              onClick={() => setIsPasswordVisible((visible) => !visible)}
              aria-label={isPasswordVisible ? t("hidePassword") : t("showPassword")}
              className="absolute right-1.5 top-1 h-8 w-8 p-0 shadow-none text-muted-foreground hover:text-foreground"
            >
              {isPasswordVisible ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
            </Button>
          </div>
        </div>

        <Button
          type="submit"
          disabled={loginMutation.isPending}
          isLoading={loginMutation.isPending}
          className="w-full py-3.5 mt-2 text-sm font-bold tracking-wide uppercase"
        >
          {loginMutation.isPending ? t("loadingText") : t("submitText")}
        </Button>
      </form>
    </Card>
  )
}
