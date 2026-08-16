import * as z from "zod"
import { cpfValidation, isPastDate, isValidICD10 } from "../../shared/utils/validators"
import { LoincCode } from "../../shared/types"

export const basePatientSchema = z.object({
  fullName: z.string(),
  birthDate: z.string(),
  documentId: z.string(),
  phoneNumber: z.string(),
})

export const baseEncounterSchema = z.object({
  reasonDisplay: z.string(),
  reasonCode: z.string().optional(),
  practitionerId: z.string(),
})

export const baseObservationSchema = z.object({
  loincCode: z.string(),
  valueQuantity: z.number(),
})

export const baseReportSchema = z.object({
  reportCode: z.string(),
  reportDisplay: z.string(),
  conclusion: z.string(),
})

export const baseConditionSchema = z.object({
  icd10Code: z.string(),
  codeDisplay: z.string(),
})

export const baseAllergySchema = z.object({
  allergenCode: z.string(),
  allergenDisplay: z.string(),
  reaction: z.string(),
})

export const baseMedicationSchema = z.object({
  medicationDisplay: z.string(),
  dosageInstruction: z.string(),
})

export type NewPatientFormData = z.infer<typeof basePatientSchema>
export type NewEncounterFormData = z.infer<typeof baseEncounterSchema>
export type NewObservationFormData = z.infer<typeof baseObservationSchema>
export type NewReportFormData = z.infer<typeof baseReportSchema>
export type NewConditionFormData = z.infer<typeof baseConditionSchema>
export type NewAllergyFormData = z.infer<typeof baseAllergySchema>
export type NewMedicationFormData = z.infer<typeof baseMedicationSchema>

export const getNewPatientSchema = (translateFunction: (key: string) => string) => z.object({
  fullName: z.string().min(3, translateFunction("validation.fullNameMin")).max(255, translateFunction("validation.maxLength")).trim(),
  birthDate: z.string().min(10, translateFunction("validation.birthDateReq")).refine(isPastDate, translateFunction("validation.birthDatePast")),
  documentId: z.string().min(11, translateFunction("validation.documentMin")).refine(cpfValidation, translateFunction("validation.documentInvalid")),
  phoneNumber: z.string().regex(/^\(\d{2}\) \d{4,5}-\d{4}$/, translateFunction("validation.phoneFormat")),
})

export const getNewEncounterSchema = (translateFunction: (key: string) => string) => z.object({
  reasonDisplay: z.string().min(3, translateFunction("validation.reasonMin")).max(255, translateFunction("validation.maxLength")),
  reasonCode: z.string().optional(),
  practitionerId: z.string().min(1, translateFunction("validation.practitionerReq")),
})

export const getNewObservationSchema = (translateFunction: (key: string) => string) => z.object({
  loincCode: z.string().min(1, translateFunction("validation.loincReq")).max(10, translateFunction("validation.maxLength")),
  valueQuantity: z.number().min(0.1, translateFunction("validation.valueReq")),
}).refine(
  (data) => {
    if (data.loincCode === LoincCode.HeartRate) {
      return data.valueQuantity >= 0 && data.valueQuantity <= 300
    }
    if (data.loincCode === LoincCode.BodyTemperature) {
      return data.valueQuantity >= 30 && data.valueQuantity <= 45
    }
    if (data.loincCode === LoincCode.BloodPressure) {
      return data.valueQuantity >= 0 && data.valueQuantity <= 300
    }
    return true
  },
  {
    message: translateFunction("validation.rangeError"),
    path: ["valueQuantity"],
  }
)

export const getNewReportSchema = (translateFunction: (key: string) => string) => z.object({
  reportCode: z.string().min(1, translateFunction("validation.reportCodeReq")).max(10, translateFunction("validation.maxLength")),
  reportDisplay: z.string().min(3, translateFunction("validation.reportTitleMin")).max(255, translateFunction("validation.maxLength")),
  conclusion: z.string().min(5, translateFunction("validation.conclusionMin")).max(2000, translateFunction("validation.maxLength")),
})

export const getNewConditionSchema = (translateFunction: (key: string) => string) => z.object({
  icd10Code: z.string().min(3, translateFunction("validation.icdCodeMin")).max(10, translateFunction("validation.maxLength")).refine(isValidICD10, translateFunction("validation.icdFormat")),
  codeDisplay: z.string().min(3, translateFunction("validation.icdDisplayMin")).max(255, translateFunction("validation.maxLength")),
})

export const getNewAllergySchema = (translateFunction: (key: string) => string) => z.object({
  allergenCode: z.string().min(3, translateFunction("validation.allergenCodeMin")).max(50, translateFunction("validation.maxLength")),
  allergenDisplay: z.string().min(3, translateFunction("validation.allergenDisplayMin")).max(255, translateFunction("validation.maxLength")),
  reaction: z.string().min(3, translateFunction("validation.reactionMin")).max(500, translateFunction("validation.maxLength")),
})

export const getNewMedicationSchema = (translateFunction: (key: string) => string) => z.object({
  medicationDisplay: z.string().min(3, translateFunction("validation.medicationMin")).max(255, translateFunction("validation.maxLength")),
  dosageInstruction: z.string().min(3, translateFunction("validation.dosageMin")).max(1000, translateFunction("validation.maxLength")),
})
