import { describe, it, expect } from 'vitest'
import {
  basePatientSchema,
  baseEncounterSchema,
  baseReportSchema,
  getNewEncounterSchema,
  getNewPatientSchema,
  getNewReportSchema,
  getNewMedicationSchema,
} from './patient_schemas'

const mockTranslate = (key: string) => `message:${key}`

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

    it('should accept empty fullName since base schema has no constraints', () => {
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
    it('should accept valid encounter data with practitioner', () => {
      const result = baseEncounterSchema.safeParse({
        reasonDisplay: 'Routine checkup',
        practitionerId: 'fhir-practitioner-1',
      })
      expect(result.success).toBe(true)
    })

    it('should reject encounter without practitioner', () => {
      const result = baseEncounterSchema.safeParse({
        reasonDisplay: 'Routine checkup',
      })
      expect(result.success).toBe(false)
    })

    it('should accept empty reasonDisplay since base schema has no constraints', () => {
      const result = baseEncounterSchema.safeParse({
        reasonDisplay: '',
        practitionerId: 'fhir-practitioner-1',
      })
      expect(result.success).toBe(true)
    })
  })

  describe('getNewEncounterSchema', () => {
    it('should require practitionerId', () => {
      const result = getNewEncounterSchema(mockTranslate).safeParse({
        reasonDisplay: 'Routine checkup',
        practitionerId: '',
      })
      expect(result.success).toBe(false)
    })

    it('should accept valid encounter data', () => {
      const result = getNewEncounterSchema(mockTranslate).safeParse({
        reasonDisplay: 'Routine checkup',
        practitionerId: 'fhir-practitioner-1',
      })
      expect(result.success).toBe(true)
    })

    it('should reject a reason longer than the max length', () => {
      const result = getNewEncounterSchema(mockTranslate).safeParse({
        reasonDisplay: 'x'.repeat(256),
        practitionerId: 'fhir-practitioner-1',
      })
      expect(result.success).toBe(false)
    })
  })

  describe('getNewPatientSchema', () => {
    it('should reject a full name longer than the max length', () => {
      const result = getNewPatientSchema(mockTranslate).safeParse({
        fullName: 'x'.repeat(256),
        birthDate: '1990-01-15',
        documentId: '111.222.333-44',
        phoneNumber: '(11) 98765-4321',
      })
      expect(result.success).toBe(false)
    })

    it('should accept valid patient data', () => {
      const result = getNewPatientSchema(mockTranslate).safeParse({
        fullName: 'John Doe',
        birthDate: '1990-01-15',
        documentId: '111.222.333-44',
        phoneNumber: '(11) 98765-4321',
      })
      expect(result.success).toBe(true)
    })
  })

  describe('baseReportSchema', () => {
    it('should accept valid report data', () => {
      const result = baseReportSchema.safeParse({
        reportCode: '58410-2',
        reportDisplay: 'Complete Blood Count',
        conclusion: 'No anomalies found in the sample.',
      })
      expect(result.success).toBe(true)
    })

    it('should reject report without reportCode', () => {
      const result = baseReportSchema.safeParse({
        reportDisplay: 'Complete Blood Count',
        conclusion: 'No anomalies found in the sample.',
      })
      expect(result.success).toBe(false)
    })
  })

  describe('getNewReportSchema', () => {
    it('should require reportCode', () => {
      const result = getNewReportSchema(mockTranslate).safeParse({
        reportCode: '',
        reportDisplay: 'Complete Blood Count',
        conclusion: 'No anomalies found in the sample.',
      })
      expect(result.success).toBe(false)
    })

    it('should accept valid report data', () => {
      const result = getNewReportSchema(mockTranslate).safeParse({
        reportCode: '58410-2',
        reportDisplay: 'Complete Blood Count',
        conclusion: 'No anomalies found in the sample.',
      })
      expect(result.success).toBe(true)
    })

    it('should reject a conclusion longer than the max length', () => {
      const result = getNewReportSchema(mockTranslate).safeParse({
        reportCode: '58410-2',
        reportDisplay: 'Complete Blood Count',
        conclusion: 'x'.repeat(2001),
      })
      expect(result.success).toBe(false)
    })
  })

  describe('getNewMedicationSchema', () => {
    it('should accept valid medication data', () => {
      const result = getNewMedicationSchema(mockTranslate).safeParse({
        medicationDisplay: 'Paracetamol 500mg',
        dosageInstruction: 'Take one tablet every 6 hours.',
      })
      expect(result.success).toBe(true)
    })

    it('should reject a dosage instruction longer than the max length', () => {
      const result = getNewMedicationSchema(mockTranslate).safeParse({
        medicationDisplay: 'Paracetamol 500mg',
        dosageInstruction: 'x'.repeat(1001),
      })
      expect(result.success).toBe(false)
    })
  })
})
