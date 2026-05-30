import * as z from "zod"
import { StaffRole } from "../../../shared/types"

export const staffFormSchema = z.object({
  fullName: z.string().min(3, "O nome deve ter no mínimo 3 caracteres").max(255),
  role: z.nativeEnum(StaffRole, {
    message: "Selecione uma categoria válida",
  }),
  license: z
    .string()
    .optional()
    .refine(
      (value) => {
        if (!value || value.trim() === "") {
          return true
        }
        return /^(CRM|COREN)(-[A-Z]{2})?[\s\-]?\d{1,6}$/i.test(value.trim())
      },
      "Formato inválido. Ex: CRM-SP 12345"
    ),
  email: z.string().email("E-mail inválido").max(255),
  department: z.string().max(100).optional(),
})

export type StaffFormData = z.infer<typeof staffFormSchema>
