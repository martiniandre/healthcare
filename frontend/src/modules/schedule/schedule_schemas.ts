import * as z from "zod"
import { isTodayOrFutureDate } from "../../shared/utils/validators"

export const allowedAppointmentDurations = [30, 45] as const

export const allowedAppointmentStartMinutes = ["00", "30", "45"] as const

export const baseAppointmentSchema = z.object({
  patientFhirId: z.string(),
  staffId: z.string(),
  date: z.string(),
  startTime: z.string(),
  endTime: z.string(),
  reason: z.string().optional(),
})

export type NewAppointmentFormData = z.infer<typeof baseAppointmentSchema>

export const getNewAppointmentSchema = (translateFunction: (key: string) => string) => z.object({
  patientFhirId: z.string().min(1, translateFunction("validation.patientRequired")),
  staffId: z.string().min(1, translateFunction("validation.staffRequired")),
  date: z
    .string()
    .min(1, translateFunction("validation.dateRequired"))
    .refine(isTodayOrFutureDate, translateFunction("validation.dateFuture")),
  startTime: z.string().min(1, translateFunction("validation.startTimeRequired")),
  endTime: z.string().min(1, translateFunction("validation.endTimeRequired")),
  reason: z.string().max(500, translateFunction("validation.reasonMax")).optional(),
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
).refine(
  (formData) => {
    if (!formData.startTime) {
      return true
    }
    const minutePart = formData.startTime.split(":")[1]
    return (allowedAppointmentStartMinutes as readonly string[]).includes(minutePart)
  },
  {
    message: translateFunction("validation.startTimeAligned"),
    path: ["startTime"],
  }
).refine(
  (formData) => {
    if (!formData.startTime || !formData.endTime) {
      return true
    }
    if (formData.endTime <= formData.startTime) {
      return true
    }
    return (allowedAppointmentDurations as readonly number[]).includes(getSlotDurationMinutes(formData.startTime, formData.endTime))
  },
  {
    message: translateFunction("validation.slotDuration"),
    path: ["endTime"],
  }
)

const getSlotDurationMinutes = (startTimeValue: string, endTimeValue: string): number => {
  const [startHour, startMinute] = startTimeValue.split(":").map(Number)
  const [endHour, endMinute] = endTimeValue.split(":").map(Number)
  return endHour * 60 + endMinute - (startHour * 60 + startMinute)
}
