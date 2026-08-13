import { describe, it, expect } from 'vitest'
import { cpfValidation, isPastDate, isTodayOrFutureDate, isValidICD10 } from './validators'

describe('cpfValidation', () => {
  it('should accept a valid CPF', () => {
    expect(cpfValidation('11122233344')).toBe(true)
  })

  it('should reject an invalid CPF with all same digits', () => {
    expect(cpfValidation('11111111111')).toBe(false)
  })

  it('should reject a too-short CPF', () => {
    expect(cpfValidation('123')).toBe(false)
  })

  it('should strip non-digits and validate', () => {
    expect(cpfValidation('111.222.333-44')).toBe(true)
  })
})

describe('isPastDate', () => {
  it('should return true for a date in the past', () => {
    expect(isPastDate('2020-01-01')).toBe(true)
  })

  it('should return false for an invalid date string', () => {
    expect(isPastDate('not-a-date')).toBe(false)
  })
})

describe('isTodayOrFutureDate', () => {
  const todayIso = new Date().toISOString().slice(0, 10)
  const futureIso = new Date(Date.now() + 5 * 86400000).toISOString().slice(0, 10)
  const pastIso = new Date(Date.now() - 5 * 86400000).toISOString().slice(0, 10)

  it('should return true for today', () => {
    expect(isTodayOrFutureDate(todayIso)).toBe(true)
  })

  it('should return true for a future date', () => {
    expect(isTodayOrFutureDate(futureIso)).toBe(true)
  })

  it('should return false for a past date', () => {
    expect(isTodayOrFutureDate(pastIso)).toBe(false)
  })

  it('should return false for an invalid date string', () => {
    expect(isTodayOrFutureDate('not-a-date')).toBe(false)
  })
})

describe('isValidICD10', () => {
  it('should accept a valid ICD-10 code', () => {
    expect(isValidICD10('I10')).toBe(true)
  })

  it('should accept a valid ICD-10 code with subcategory', () => {
    expect(isValidICD10('J45.9')).toBe(true)
  })

  it('should reject an invalid code', () => {
    expect(isValidICD10('abc')).toBe(false)
  })
})
