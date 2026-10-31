import { useMemo, useState, type KeyboardEvent } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { isAxiosError } from "axios"
import { useAuthStore } from "../../../shared/store/auth_store"
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/components/ui/card"
import { Input } from "@/shared/components/ui/input"
import { Button } from "@/shared/components/ui/button"
import { Spinner } from "@/shared/components/ui/spinner"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { Field, FieldDescription, FieldError, FieldLabel } from "@/shared/components/ui/field"
import { getLoginFormSchema, type LoginFormData } from "../auth_schemas"
import { Eye, EyeOff, KeyRound, Mail, ShieldAlert, TriangleAlert } from "lucide-react"
import { useLoginMutation } from "../queries"

const CREDENTIAL_ERROR_STATUSES = [400, 401, 403]

export const LoginForm = () => {
  const { t } = useTranslation("auth")
  const loginToStore = useAuthStore((state) => state.login)
  const [generalError, setGeneralError] = useState<string | null>(null)
  const [isPasswordVisible, setIsPasswordVisible] = useState(false)
  const [isCapsLockOn, setIsCapsLockOn] = useState(false)
  const loginMutation = useLoginMutation()
  const loginSchema = useMemo(() => getLoginFormSchema(t), [t])

  const {
    register,
    handleSubmit,
    setFocus,
    setValue,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    mode: "onTouched",
    defaultValues: { email: "", password: "" },
  })

  const resolveLoginErrorMessage = (loginRequestError: unknown) => {
    if (isAxiosError(loginRequestError) && loginRequestError.response) {
      return CREDENTIAL_ERROR_STATUSES.includes(loginRequestError.response.status)
        ? t("invalidCredentials")
        : t("defaultError")
    }
    return t("defaultError")
  }

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
      setGeneralError(resolveLoginErrorMessage(loginRequestError))
      setValue("password", "")
      setFocus("password")
    }
  }

  const handlePasswordKeyEvent = (event: KeyboardEvent<HTMLInputElement>) => {
    setIsCapsLockOn(event.getModifierState("CapsLock"))
  }

  const passwordRegistration = register("password")
  const isPending = loginMutation.isPending
  const passwordDescribedBy = [
    isCapsLockOn ? "login-password-caps" : null,
    errors.password ? "login-password-error" : null,
  ].filter(Boolean).join(" ") || undefined

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">{t("authTitle")}</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-6">
        {generalError && (
          <Alert variant="destructive">
            <ShieldAlert />
            <AlertDescription>{generalError}</AlertDescription>
          </Alert>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5" noValidate aria-busy={isPending}>
          <Field data-invalid={!!errors.email}>
            <FieldLabel htmlFor="login-email">
              <Mail className="size-3.5 text-primary" />
              {t("emailLabel")}
            </FieldLabel>
            <Input
              id="login-email"
              type="email"
              inputMode="email"
              placeholder={t("emailPlaceholder")}
              autoComplete="username"
              autoCapitalize="none"
              autoCorrect="off"
              spellCheck={false}
              autoFocus
              maxLength={255}
              aria-required="true"
              aria-invalid={!!errors.email}
              aria-describedby={errors.email ? "login-email-error" : undefined}
              className="h-11 md:h-9"
              {...register("email")}
            />
            <FieldError id="login-email-error">{errors.email?.message}</FieldError>
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
                autoCapitalize="none"
                spellCheck={false}
                maxLength={128}
                aria-required="true"
                aria-invalid={!!errors.password}
                aria-describedby={passwordDescribedBy}
                className="h-11 pr-28 md:h-9"
                {...passwordRegistration}
                onKeyDown={handlePasswordKeyEvent}
                onKeyUp={handlePasswordKeyEvent}
                onBlur={(event) => {
                  setIsCapsLockOn(false)
                  void passwordRegistration.onBlur(event)
                }}
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                aria-pressed={isPasswordVisible}
                onClick={() => setIsPasswordVisible((visible) => !visible)}
                className="absolute right-1 top-1/2 h-9 -translate-y-1/2 gap-1 px-2 text-xs text-muted-foreground md:h-7"
              >
                {isPasswordVisible ? <EyeOff /> : <Eye />}
                {t("togglePassword")}
              </Button>
            </div>
            {isCapsLockOn && (
              <FieldDescription id="login-password-caps" role="status" className="flex items-center gap-1.5 text-xs text-warning">
                <TriangleAlert className="size-3.5" aria-hidden="true" />
                {t("capsLockOn")}
              </FieldDescription>
            )}
            <FieldError id="login-password-error">{errors.password?.message}</FieldError>
          </Field>

          <Button
            type="submit"
            disabled={isPending}
            className="mt-2 h-11 w-full font-bold uppercase tracking-wide md:h-9"
          >
            {isPending && <Spinner />}
            {isPending ? t("loadingText") : t("submitText")}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
