import { describe, it, expect } from 'vitest'
import {
  basePatientSchema,
  baseEncounterSchema,
  baseReportSchema,
  getNewEncounterSchema,
  getNewPatientSchema,
  getNewReportSchema,
  getNewMedicationSchema,
  getNewVitalSignsPanelSchema,
  vitalSignMetricDefinitions,
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

  describe('getNewVitalSignsPanelSchema', () => {
    it('should accept an entirely empty panel since every metric is optional', () => {
      const result = getNewVitalSignsPanelSchema(mockTranslate).safeParse({})
      expect(result.success).toBe(true)
    })

    it('should treat empty numeric inputs as unmeasured metrics', () => {
      const result = getNewVitalSignsPanelSchema(mockTranslate).safeParse({
        heartRate: Number.NaN,
        bodyTemperature: Number.NaN,
      })
      expect(result.success).toBe(true)
      if (result.success) {
        expect(result.data.heartRate).toBeUndefined()
        expect(result.data.bodyTemperature).toBeUndefined()
      }
    })

    it('should accept a mixed panel within the clinical ranges', () => {
      const result = getNewVitalSignsPanelSchema(mockTranslate).safeParse({
        heartRate: 72,
        bodyTemperature: 36.5,
        systolicBloodPressure: 120,
        diastolicBloodPressure: 80,
        oxygenSaturation: 98,
        respiratoryRate: 16,
        weightKg: 70.5,
        heightCm: 175,
      })
      expect(result.success).toBe(true)
    })

    it.each(vitalSignMetricDefinitions.map((metricDefinition) => [
      metricDefinition.formFieldName,
      metricDefinition.minimumValue,
      metricDefinition.maximumValue,
    ]))('should accept boundary values for %s', (formFieldName, minimumValue, maximumValue) => {
      const schema = getNewVitalSignsPanelSchema(mockTranslate)
      expect(schema.safeParse({ [formFieldName]: minimumValue }).success).toBe(true)
      expect(schema.safeParse({ [formFieldName]: maximumValue }).success).toBe(true)
    })

    it.each(vitalSignMetricDefinitions.map((metricDefinition) => [
      metricDefinition.formFieldName,
      metricDefinition.minimumValue - 1,
      metricDefinition.maximumValue + 1,
    ]))('should reject out-of-range values for %s at the offending field path', (formFieldName, belowMinimum, aboveMaximum) => {
      const schema = getNewVitalSignsPanelSchema(mockTranslate)

      const belowResult = schema.safeParse({ [formFieldName]: belowMinimum })
      expect(belowResult.success).toBe(false)
      if (!belowResult.success) {
        expect(belowResult.error.issues[0].path).toEqual([formFieldName])
        expect(belowResult.error.issues[0].message).toContain('validation.vitalSignsRange')
      }

      const aboveResult = schema.safeParse({ [formFieldName]: aboveMaximum })
      expect(aboveResult.success).toBe(false)
      if (!aboveResult.success) {
        expect(aboveResult.error.issues[0].path).toEqual([formFieldName])
      }
    })
  })
})
