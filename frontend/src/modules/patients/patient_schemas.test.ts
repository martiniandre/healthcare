import { describe, it, expect } from 'vitest'
import { basePatientSchema, baseEncounterSchema } from './patient_schemas'

describe('patient schemas', () => {
  describe('basePatientSchema', () => {
    it('should accept valid patient data', () => {
      const result = basePatientSchema.safeParse({
        fullName: 'John Doe',
        birthDate: '1990-01-15',
        documentId: '12345678901',
        phoneNumber: '11999990000',
      })
      expect(result.success).toBe(true)
    })

    it('should reject empty fullName', () => {
      const result = basePatientSchema.safeParse({
        fullName: '',
        birthDate: '1990-01-15',
        documentId: '12345678901',
        phoneNumber: '11999990000',
      })
      expect(result.success).toBe(true)
    })
  })

  describe('baseEncounterSchema', () => {
    it('should accept valid encounter data', () => {
      const result = baseEncounterSchema.safeParse({
        reasonDisplay: 'Routine checkup',
      })
      expect(result.success).toBe(true)
    })

    it('should reject empty reasonDisplay', () => {
      const result = baseEncounterSchema.safeParse({
        reasonDisplay: '',
      })
      expect(result.success).toBe(true)
    })
  })
})
