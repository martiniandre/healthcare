import * as z from "zod"

export const baseAppointmentSchema = z.object({
  patientFhirId: z.string(),
  staffId: z.string(),
  date: z.string(),
  startTime: z.string(),
  endTime: z.string(),
  reason: z.string(),
})

export type NewAppointmentFormData = z.infer<typeof baseAppointmentSchema>

export const getNewAppointmentSchema = (translateFunction: (key: string) => string) => z.object({
  patientFhirId: z.string().min(1, translateFunction("validation.patientRequired")),
  staffId: z.string().min(1, translateFunction("validation.staffRequired")),
  date: z.string().min(1, translateFunction("validation.dateRequired")),
  startTime: z.string().min(1, translateFunction("validation.startTimeRequired")),
  endTime: z.string().min(1, translateFunction("validation.endTimeRequired")),
  reason: z.string().min(3, translateFunction("validation.reasonMin")),
}).refine(
  (formData) => {
    if (!formData.startTime || !formData.endTime) {
      return true
    }
    return formData.endTime > formData.startTime
  },
  {
    message: translateFunction("validation.endAfterStart"),
    path: ["endTime"],
  }
)
