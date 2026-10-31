import { allowedAppointmentDurations } from "./schedule_schemas"
import { todayDateString } from "../../shared/utils/validators"

const MINUTES_ON_FIFTEEN_MINUTE_GRID = [0, 15, 30, 45] as const

export interface TimeSlotOption {
  value: string
  label: string
}

const formatSlotValue = (hour: number, minute: number): string => {
  return `${hour.toString().padStart(2, "0")}:${minute.toString().padStart(2, "0")}`
}

export const getTimeSlotOptions = (): readonly TimeSlotOption[] => {
  const timeSlotOptions: TimeSlotOption[] = []
  for (let hour = 0; hour < 24; hour++) {
    for (const minute of MINUTES_ON_FIFTEEN_MINUTE_GRID) {
      timeSlotOptions.push({ value: formatSlotValue(hour, minute), label: formatSlotValue(hour, minute) })
    }
  }
  return timeSlotOptions
}

export const getAvailableStartTimeOptions = (selectedDate: string): readonly TimeSlotOption[] => {
  if (selectedDate !== todayDateString()) {
    return getTimeSlotOptions()
  }
  const currentDate = new Date()
  return getTimeSlotOptions().filter((timeSlot) => {
    const [optionHour, optionMinute] = timeSlot.value.split(":").map(Number)
    if (optionHour > currentDate.getHours()) {
      return true
    }
    if (optionHour === currentDate.getHours()) {
      return optionMinute > currentDate.getMinutes()
    }
    return false
  })
}

export const getEndTimeOptionsForStart = (startTimeValue: string): readonly TimeSlotOption[] => {
  if (!startTimeValue) {
    return []
  }
  const [startHour, startMinute] = startTimeValue.split(":").map(Number)
  const startTotalMinutes = startHour * 60 + startMinute
  const endTimeOptions: TimeSlotOption[] = []
  for (const durationMinutes of allowedAppointmentDurations) {
    const endTotalMinutes = startTotalMinutes + durationMinutes
    if (endTotalMinutes >= 24 * 60) {
      continue
    }
    const endHour = Math.floor(endTotalMinutes / 60)
    const endMinute = endTotalMinutes % 60
    endTimeOptions.push({ value: formatSlotValue(endHour, endMinute), label: formatSlotValue(endHour, endMinute) })
  }
  return endTimeOptions
}

export const isTimeOnFifteenMinuteGrid = (timeValue: string): boolean => {
  const minutePart = timeValue.split(":")[1]
  return (MINUTES_ON_FIFTEEN_MINUTE_GRID as readonly number[]).includes(Number(minutePart))
}
export const getDefaultAppointmentEndTime = (startTimeValue: string): string => {
  return getEndTimeOptionsForStart(startTimeValue)[0]?.value ?? ""
}

export const getSlotOptionsAfter = (startTimeValue: string): readonly TimeSlotOption[] => {
  if (!startTimeValue) {
    return getTimeSlotOptions()
  }
  return getTimeSlotOptions().filter((timeSlot) => timeSlot.value > startTimeValue)
}

export const getDefaultBlockedPeriodEndTime = (startTimeValue: string, durationMinutes: number): string => {
  if (!startTimeValue) {
    return ""
  }
  const [startHour, startMinute] = startTimeValue.split(":").map(Number)
  const endTotalMinutes = startHour * 60 + startMinute + durationMinutes
  if (endTotalMinutes >= 24 * 60) {
    return getSlotOptionsAfter(startTimeValue).at(-1)?.value ?? ""
  }
  return formatSlotValue(Math.floor(endTotalMinutes / 60), endTotalMinutes % 60)
}

const toTotalMinutes = (timeValue: string): number => {
  const [hour, minute] = timeValue.split(":").map(Number)
  return hour * 60 + minute
}

export const getAppointmentEndTimeForNewStart = (
  previousStartTime: string,
  previousEndTime: string,
  nextStartTime: string
): string => {
  const endTimeOptions = getEndTimeOptionsForStart(nextStartTime)
  if (previousStartTime && previousEndTime) {
    const previousDuration = toTotalMinutes(previousEndTime) - toTotalMinutes(previousStartTime)
    const keptDurationOption = endTimeOptions.find(
      (timeSlot) => toTotalMinutes(timeSlot.value) - toTotalMinutes(nextStartTime) === previousDuration
    )
    if (keptDurationOption) {
      return keptDurationOption.value
    }
  }
  return endTimeOptions[0]?.value ?? ""
}

export const getDurationInMinutes = (startTimeValue: string, endTimeValue: string): number => {
  if (!startTimeValue || !endTimeValue) {
    return 0
  }
  return toTotalMinutes(endTimeValue) - toTotalMinutes(startTimeValue)
}
