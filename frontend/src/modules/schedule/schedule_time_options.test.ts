import { describe, it, expect, vi, afterEach } from "vitest"
import {
  getTimeSlotOptions,
  getAvailableStartTimeOptions,
  getEndTimeOptionsForStart,
  isTimeOnFifteenMinuteGrid,
  getDefaultAppointmentEndTime,
  getAppointmentEndTimeForNewStart,
  getSlotOptionsAfter,
  getDefaultBlockedPeriodEndTime,
  getDurationInMinutes,
} from "./schedule_time_options"

describe("getTimeSlotOptions", () => {
  const timeSlotOptions = getTimeSlotOptions()

  it("should return exactly 96 options", () => {
    expect(timeSlotOptions).toHaveLength(96)
  })

  it("should start at 00:00 and end at 23:45", () => {
    expect(timeSlotOptions[0].value).toBe("00:00")
    expect(timeSlotOptions[timeSlotOptions.length - 1].value).toBe("23:45")
  })

  it("should return values matching the HH:MM format", () => {
    for (const timeSlot of timeSlotOptions) {
      expect(timeSlot.value).toMatch(/^\d{2}:\d{2}$/)
    }
  })

  it("should only expose minutes on the fifteen minute grid", () => {
    const allowedMinuteValues = ["00", "15", "30", "45"]
    for (const timeSlot of timeSlotOptions) {
      const minutePart = timeSlot.value.split(":")[1]
      expect(allowedMinuteValues).toContain(minutePart)
    }
  })

  it("should return strictly ascending values", () => {
    const optionValues = timeSlotOptions.map((timeSlot) => timeSlot.value)
    const sortedValues = [...optionValues].sort()
    expect(optionValues).toEqual(sortedValues)
    expect(new Set(optionValues).size).toBe(optionValues.length)
  })

  it("should use the value as the displayed label", () => {
    for (const timeSlot of timeSlotOptions) {
      expect(timeSlot.label).toBe(timeSlot.value)
    }
  })
})

describe("getAvailableStartTimeOptions", () => {
  afterEach(() => {
    vi.useRealTimers()
  })

  it("should return all slots for a future date", () => {
    expect(getAvailableStartTimeOptions("2099-01-01")).toHaveLength(96)
  })

  it("should block today's slots that are already in the past", () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date("2099-01-01T14:32:00"))
    const optionValues = getAvailableStartTimeOptions("2099-01-01").map((timeSlot) => timeSlot.value)
    expect(optionValues).toContain("14:45")
    expect(optionValues).not.toContain("14:30")
    expect(optionValues).not.toContain("09:00")
  })

  it("should keep the current aligned slot when its minute is still in the future", () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date("2099-01-01T14:28:00"))
    const optionValues = getAvailableStartTimeOptions("2099-01-01").map((timeSlot) => timeSlot.value)
    expect(optionValues).toContain("14:30")
  })
})

describe("getEndTimeOptionsForStart", () => {
  it("should return no options when no start time is selected", () => {
    expect(getEndTimeOptionsForStart("")).toEqual([])
  })

  it("should return only the 30 and 45 minute interval options", () => {
    expect(getEndTimeOptionsForStart("09:00")).toEqual([
      { value: "09:30", label: "09:30" },
      { value: "09:45", label: "09:45" },
    ])
  })

  it("should keep only the interval that fits before midnight", () => {
    expect(getEndTimeOptionsForStart("23:15")).toEqual([
      { value: "23:45", label: "23:45" },
    ])
  })

  it("should return no options when the start time is too close to midnight", () => {
    expect(getEndTimeOptionsForStart("23:45")).toEqual([])
  })

  it("should handle the hour rollover", () => {
    expect(getEndTimeOptionsForStart("08:45")).toEqual([
      { value: "09:15", label: "09:15" },
      { value: "09:30", label: "09:30" },
    ])
  })
})

describe("isTimeOnFifteenMinuteGrid", () => {
  it("should accept every quarter hour", () => {
    expect(isTimeOnFifteenMinuteGrid("09:00")).toBe(true)
    expect(isTimeOnFifteenMinuteGrid("09:15")).toBe(true)
    expect(isTimeOnFifteenMinuteGrid("09:30")).toBe(true)
    expect(isTimeOnFifteenMinuteGrid("09:45")).toBe(true)
  })

  it("should reject off grid minutes", () => {
    expect(isTimeOnFifteenMinuteGrid("09:07")).toBe(false)
    expect(isTimeOnFifteenMinuteGrid("09:59")).toBe(false)
  })
})
describe('end time helpers', () => {
  it('should default the appointment end to 30 minutes after the start', () => {
    expect(getDefaultAppointmentEndTime('09:15')).toBe('09:45')
    expect(getDefaultAppointmentEndTime('')).toBe('')
  })

  it('should keep the previous duration when the start changes', () => {
    expect(getAppointmentEndTimeForNewStart('09:00', '09:45', '10:00')).toBe('10:45')
    expect(getAppointmentEndTimeForNewStart('', '', '10:00')).toBe('10:30')
    expect(getAppointmentEndTimeForNewStart('09:00', '10:00', '10:00')).toBe('10:30')
  })

  it('should only offer blocked period end slots after the start', () => {
    const options = getSlotOptionsAfter('23:00')
    expect(options.map((option) => option.value)).toEqual(['23:15', '23:30', '23:45'])
    expect(getSlotOptionsAfter('')).toHaveLength(96)
  })

  it('should default the blocked period end and clamp at the end of the day', () => {
    expect(getDefaultBlockedPeriodEndTime('12:00', 60)).toBe('13:00')
    expect(getDefaultBlockedPeriodEndTime('23:30', 60)).toBe('23:45')
    expect(getDefaultBlockedPeriodEndTime('', 60)).toBe('')
  })

  it('should compute durations in minutes', () => {
    expect(getDurationInMinutes('09:00', '09:45')).toBe(45)
    expect(getDurationInMinutes('', '09:45')).toBe(0)
  })
})
