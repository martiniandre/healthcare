import { describe, it, expect } from 'vitest'
import { getNewAppointmentSchema } from './schedule_schemas'

const mockTranslate = (key: string) => `message:${key}`

const todayIso = new Date().toISOString().slice(0, 10)
const pastIso = new Date(Date.now() - 5 * 86400000).toISOString().slice(0, 10)

const validAppointment = {
  patientFhirId: 'fhir-patient-1',
  staffId: 'fhir-staff-1',
  date: todayIso,
  startTime: '09:00',
  endTime: '10:00',
  reason: 'Routine checkup',
}

describe('getNewAppointmentSchema', () => {
  it('should accept a valid appointment', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse(validAppointment)
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

  it('should reject a reason shorter than the minimum', () => {
    const result = getNewAppointmentSchema(mockTranslate).safeParse({
      ...validAppointment,
      reason: 'ab',
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
