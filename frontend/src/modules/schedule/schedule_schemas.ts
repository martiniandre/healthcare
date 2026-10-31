import * as z from "zod"
import { isTodayOrFutureDate, todayDateString } from "../../shared/utils/validators"

export const allowedAppointmentDurations = [30, 45] as const

export const allowedAppointmentStartMinutes = ["00", "15", "30", "45"] as const

export const baseAppointmentSchema = z.object({
  patientFhirId: z.string(),
  staffId: z.string(),
  date: z.string(),
  startTime: z.string(),
  endTime: z.string(),
  reason: z.string().optional(),
})

export type NewAppointmentFormData = z.infer<typeof baseAppointmentSchema>

export type ScheduleTranslate = (key: string, options?: Record<string, string | number>) => string

const getSlotDurationMinutes = (startTimeValue: string, endTimeValue: string): number => {
  const [startHour, startMinute] = startTimeValue.split(":").map(Number)
  const [endHour, endMinute] = endTimeValue.split(":").map(Number)
  return endHour * 60 + endMinute - (startHour * 60 + startMinute)
}

const addMinutesToTime = (timeValue: string, minutesToAdd: number): string => {
  const [hour, minute] = timeValue.split(":").map(Number)
  const totalMinutes = hour * 60 + minute + minutesToAdd
  const resultHour = Math.floor(totalMinutes / 60) % 24
  const resultMinute = totalMinutes % 60
  return `${String(resultHour).padStart(2, "0")}:${String(resultMinute).padStart(2, "0")}`
}

const isStartTimeInThePast = (dateValue: string, startTimeValue: string): boolean => {
  if (dateValue !== todayDateString()) {
    return false
  }
  const [startHour, startMinute] = startTimeValue.split(":").map(Number)
  const currentDate = new Date()
  if (startHour > currentDate.getHours()) {
    return false
  }
  if (startHour === currentDate.getHours()) {
    return startMinute <= currentDate.getMinutes()
  }
  return true
}

export const getNewAppointmentSchema = (translateFunction: ScheduleTranslate) => z.object({
  patientFhirId: z.string().min(1, translateFunction("validation.patientRequired")),
  staffId: z.string().min(1, translateFunction("validation.staffRequired")),
  date: z
    .string()
    .min(1, translateFunction("validation.dateRequired"))
    .refine(isTodayOrFutureDate, translateFunction("validation.dateFuture")),
  startTime: z.string().min(1, translateFunction("validation.startTimeRequired")),
  endTime: z.string().min(1, translateFunction("validation.endTimeRequired")),
  reason: z.string().max(500, translateFunction("validation.reasonMax")).optional(),
}).superRefine((formData, context) => {
  if (formData.startTime) {
    const minutePart = formData.startTime.split(":")[1]
    if (!(allowedAppointmentStartMinutes as readonly string[]).includes(minutePart)) {
      context.addIssue({
        code: "custom",
        message: translateFunction("validation.startTimeAligned"),
        path: ["startTime"],
      })
    }
  }
  if (formData.startTime && formData.endTime) {
    if (formData.endTime <= formData.startTime) {
      context.addIssue({
        code: "custom",
        message: translateFunction("validation.endAfterStart", { startTime: formData.startTime }),
        path: ["endTime"],
      })
    } else if (!(allowedAppointmentDurations as readonly number[]).includes(getSlotDurationMinutes(formData.startTime, formData.endTime))) {
      context.addIssue({
        code: "custom",
        message: translateFunction("validation.slotDuration", {
          startTime: formData.startTime,
          suggestions: allowedAppointmentDurations
            .map((durationMinutes) => addMinutesToTime(formData.startTime, durationMinutes))
            .join(" / "),
        }),
        path: ["endTime"],
      })
    }
  }
  if (formData.date && formData.startTime && isStartTimeInThePast(formData.date, formData.startTime)) {
    context.addIssue({
      code: "custom",
      message: translateFunction("validation.startTimePast"),
      path: ["startTime"],
    })
  }
})

export type UnavailabilityFormData = {
  date: string
  startTime: string
  endTime: string
  reason?: string
}

export const getUnavailabilitySchema = (translateFunction: ScheduleTranslate) => z.object({
  date: z.string().min(1, translateFunction("validation.dateRequired"))
    .refine(isTodayOrFutureDate, translateFunction("validation.dateFuture")),
  startTime: z.string().min(1, translateFunction("validation.startTimeRequired")),
  endTime: z.string().min(1, translateFunction("validation.endTimeRequired")),
  reason: z.string().max(500, translateFunction("validation.reasonMax")).optional(),
}).superRefine((formData, context) => {
  if (formData.startTime && formData.endTime && formData.endTime <= formData.startTime) {
    context.addIssue({
      code: "custom",
      message: translateFunction("validation.endAfterStart", { startTime: formData.startTime }),
      path: ["endTime"],
    })
  }
})
