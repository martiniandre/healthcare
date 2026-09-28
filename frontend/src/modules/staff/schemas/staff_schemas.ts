import * as z from "zod"
import { StaffRole } from "../../../shared/types"

export const staffFormSchema = z.object({
  fullName: z.string().trim().min(3, "validation.nameMinLength").max(255, "validation.nameMaxLength"),
  role: z.nativeEnum(StaffRole, {
    error: "validation.roleRequired",
  }),
  license: z
    .string()
    .trim()
    .min(1, "validation.licenseRequired")
    .max(50, "validation.licenseMaxLength")
    .refine(
      (value) => /^(CRM|COREN)(-[A-Z]{2})?[\s-]?\d{1,6}$/i.test(value),
      "validation.licenseFormat"
    ),
  email: z
    .string()
    .trim()
    .min(1, "validation.emailRequired")
    .email("validation.emailInvalid")
    .max(255, "validation.emailMaxLength"),
  departmentId: z.string().trim().min(1, "validation.departmentRequired"),
})

export type StaffFormData = z.infer<typeof staffFormSchema>
