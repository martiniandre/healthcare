import * as z from "zod"

export const getLoginFormSchema = (t: (key: string) => string) => z.object({
  email: z.string().min(1, t("auth.validation.emailRequired")).email(t("auth.validation.emailInvalid")).max(255),
  password: z.string().min(8, t("auth.validation.passwordMinLength")).max(128),
})

export type LoginFormData = {
  email: string
  password: string
}

