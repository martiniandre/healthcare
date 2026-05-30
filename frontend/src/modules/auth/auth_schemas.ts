import * as z from "zod"

export const loginFormSchema = z.object({
  email: z.string().min(1, "O e-mail é obrigatório").email("Formato de e-mail inválido").max(255),
  password: z.string().min(8, "A senha deve ter no mínimo 8 caracteres").max(128),
})

export type LoginFormData = z.infer<typeof loginFormSchema>
