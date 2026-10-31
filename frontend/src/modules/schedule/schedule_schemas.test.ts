import { describe, it, expect, vi, afterEach } from 'vitest'
import { getNewAppointmentSchema, getUnavailabilitySchema } from './schedule_schemas'

afterEach(() => {
  vi.useRealTimers()
})

const mockTranslate = (key: string) => `message:${key}`

const todayIso = new Date().toISOString().slice(0, 10)
const pastIso = new Date(Date.now() - 5 * 86400000).toISOString().slice(0, 10)

const validAppointment = {
  patientFhirId: 'fhir-patient-1',
  staffId: 'fhir-staff-1',
  date: '2099-01-01',
  startTime: '09:00',
  endTime: '09:30',
  reason: 'Routine checkup',
}

describe('getNewAppointmentSchema', () => {
  it('should accept a valid 30 minute appointment', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse(validAppointment)
    expect(result.success).toBe(true)
  })

  it('should accept a valid 45 minute appointment', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      startTime: '09:30',
      endTime: '10:15',
    })
    expect(result.success).toBe(true)
  })

  it('should accept an appointment without a reason', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      reason: '',
    })
    expect(result.success).toBe(true)
  })

  it('should require a patient', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      patientFhirId: '',
    })
    expect(result.success).toBe(false)
  })

  it('should require a staff member', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      staffId: '',
    })
    expect(result.success).toBe(false)
  })

  it('should reject a past date', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      date: pastIso,
    })
    expect(result.success).toBe(false)
  })

  it('should reject an end time before the start time', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      startTime: '10:00',
      endTime: '09:00',
    })
    expect(result.success).toBe(false)
  })

  it('should reject a one hour duration', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      startTime: '09:00',
      endTime: '10:00',
    })
    expect(result.success).toBe(false)
  })

  it('should reject an arbitrary duration', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      startTime: '09:00',
      endTime: '09:20',
    })
    expect(result.success).toBe(false)
  })

  it('should accept a 15 minute aligned start time', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      startTime: '09:15',
      endTime: '09:45',
    })
    expect(result.success).toBe(true)
  })

  it('should reject an unaligned start time', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      startTime: '09:07',
      endTime: '09:45',
    })
    expect(result.success).toBe(false)
  })

  it('should reject a start time in the past when the date is today', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2099-01-01T14:32:00'))
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      date: '2099-01-01',
      startTime: '14:15',
      endTime: '14:45',
    })
    expect(result.success).toBe(false)
  })

  it('should accept a start time in the future when the date is today', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2099-01-01T14:32:00'))
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      date: '2099-01-01',
      startTime: '15:00',
      endTime: '15:30',
    })
    expect(result.success).toBe(true)
  })

  it('should ignore the past rule when the date is not today', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2099-01-01T14:32:00'))
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      date: '2099-02-01',
      startTime: '09:00',
      endTime: '09:30',
    })
    expect(result.success).toBe(true)
  })

  it('should reject a reason longer than the maximum', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      reason: 'x'.repeat(501),
    })
    expect(result.success).toBe(false)
  })
})

const validUnavailability = {
  date: todayIso,
  startTime: '09:00',
  endTime: '12:00',
  reason: 'Vacation',
}

describe('getUnavailabilitySchema', () => {
  it('should accept a valid unavailability window', () => {
    const result = getUnavailabilitySchema(mockTranslate).safeParse(validUnavailability)
    expect(result.success).toBe(true)
  })

  it('should accept an unavailability window without a reason', () => {
    const result = getUnavailabilitySchema(mockTranslate).safeParse({
      ...validUnavailability,
      reason: '',
    })
    expect(result.success).toBe(true)
  })

  it('should reject a past date', () => {
    const result = getUnavailabilitySchema(mockTranslate).safeParse({
      ...validUnavailability,
      date: pastIso,
    })
    expect(result.success).toBe(false)
  })

  it('should reject an end time before the start time', () => {
    const result = getUnavailabilitySchema(mockTranslate).safeParse({
      ...validUnavailability,
      startTime: '12:00',
      endTime: '09:00',
    })
    expect(result.success).toBe(false)
  })

  it('should require a start time', () => {
    const result = getUnavailabilitySchema(mockTranslate).safeParse({
      ...validUnavailability,
      startTime: '',
    })
    expect(result.success).toBe(false)
  })

  it('should require an end time', () => {
    const result = getUnavailabilitySchema(mockTranslate).safeParse({
      ...validUnavailability,
      endTime: '',
    })
    expect(result.success).toBe(false)
  })
})

describe('actionable validation messages', () => {
  const capturingTranslate = (key: string, options?: Record<string, string | number>) =>
    options ? `${key}|${JSON.stringify(options)}` : key

  it('should include the start time in the end-before-start message', () => {
    const result = getNewAppointmentSchema(capturingTranslate).safeParse({
      ...validAppointment,
      startTime: '09:30',
      endTime: '09:00',
    })
    expect(result.success).toBe(false)
    expect(result.error?.issues[0].message).toBe('validation.endAfterStart|{"startTime":"09:30"}')
    expect(result.error?.issues[0].path).toEqual(['endTime'])
  })

  it('should suggest valid end times when the duration is not allowed', () => {
    const result = getNewAppointmentSchema(capturingTranslate).safeParse({
      ...validAppointment,
      startTime: '09:00',
      endTime: '10:00',
    })
    expect(result.success).toBe(false)
    expect(result.error?.issues[0].message).toBe(
      'validation.slotDuration|{"startTime":"09:00","suggestions":"09:30 / 09:45"}'
    )
  })

  it('should include the start time in the unavailability end-before-start message', () => {
    const result = getUnavailabilitySchema(capturingTranslate).safeParse({
      date: '2099-01-01',
      startTime: '12:00',
      endTime: '11:00',
      reason: '',
    })
    expect(result.success).toBe(false)
    expect(result.error?.issues[0].message).toBe('validation.endAfterStart|{"startTime":"12:00"}')
  })
})
