import { describe, it, expect } from 'vitest'
import { getNewAppointmentSchema, getUnavailabilitySchema } from './schedule_schemas'

const mockTranslate = (key: string) => `message:${key}`

const todayIso = new Date().toISOString().slice(0, 10)
const pastIso = new Date(Date.now() - 5 * 86400000).toISOString().slice(0, 10)

const validAppointment = {
  patientFhirId: 'fhir-patient-1',
  staffId: 'fhir-staff-1',
  date: todayIso,
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

  it('should reject an unaligned start time', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      startTime: '09:15',
      endTime: '09:45',
    })
    expect(result.success).toBe(false)
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
