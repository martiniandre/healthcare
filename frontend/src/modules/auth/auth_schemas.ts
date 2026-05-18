import * as z from "zod"

export const loginFormSchema = z.object({
  email: z.string().min(1, "O e-mail é obrigatório").email("Formato de e-mail inválido"),
  password: z.string().min(6, "A senha deve ter no mínimo 6 caracteres"),
})

export type LoginFormData = z.infer<typeof loginFormSchema>
