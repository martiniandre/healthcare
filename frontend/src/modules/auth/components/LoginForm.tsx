import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { useAuthStore } from "../../../shared/store/auth_store"
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/components/ui/card"
import { Input } from "@/shared/components/ui/input"
import { Button } from "@/shared/components/ui/button"
import { Spinner } from "@/shared/components/ui/spinner"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { Field, FieldError, FieldLabel } from "@/shared/components/ui/field"
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
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">{t("authTitle")}</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-6">
        {generalError && (
          <Alert variant="destructive">
            <ShieldAlert className="h-4 w-4" />
            <AlertDescription>{generalError}</AlertDescription>
          </Alert>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5" noValidate>
          <Field data-invalid={!!errors.email}>
            <FieldLabel htmlFor="login-email">
              <Mail className="size-3.5 text-primary" />
              {t("emailLabel")}
            </FieldLabel>
            <Input
              id="login-email"
              type="email"
              placeholder={t("emailPlaceholder")}
              autoComplete="email"
              maxLength={255}
              aria-invalid={!!errors.email}
              {...register("email")}
            />
            <FieldError>{errors.email?.message}</FieldError>
          </Field>

          <Field data-invalid={!!errors.password}>
            <FieldLabel htmlFor="login-password">
              <KeyRound className="size-3.5 text-primary" />
              {t("passwordLabel")}
            </FieldLabel>
            <div className="relative">
              <Input
                id="login-password"
                type={isPasswordVisible ? "text" : "password"}
                placeholder={t("passwordPlaceholder")}
                autoComplete="current-password"
                aria-invalid={!!errors.password}
                className="pr-10"
                {...register("password")}
              />
              <Button
                type="button"
                variant="ghost"
                size="icon-xs"
                onClick={() => setIsPasswordVisible((visible) => !visible)}
                aria-label={isPasswordVisible ? t("hidePassword") : t("showPassword")}
                className="absolute right-2 top-1.5 text-muted-foreground"
              >
                {isPasswordVisible ? <EyeOff /> : <Eye />}
              </Button>
            </div>
            <FieldError>{errors.password?.message}</FieldError>
          </Field>

          <Button
            type="submit"
            disabled={loginMutation.isPending}
            className="w-full mt-2 font-bold tracking-wide uppercase"
          >
            {loginMutation.isPending && <Spinner />}
            {loginMutation.isPending ? t("loadingText") : t("submitText")}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
