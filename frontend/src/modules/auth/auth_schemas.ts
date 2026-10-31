import * as z from "zod"

export const getLoginFormSchema = (t: (key: string) => string) => z.object({
  email: z
    .string()
    .trim()
    .min(1, t("validation.emailRequired"))
    .email(t("validation.emailInvalid"))
    .max(255, t("validation.emailMaxLength")),
  password: z
    .string()
    .min(1, t("validation.passwordRequired"))
    .min(8, t("validation.passwordMinLength"))
    .max(128, t("validation.passwordMaxLength")),
})

export type LoginFormData = {
  email: string
  password: string
}
