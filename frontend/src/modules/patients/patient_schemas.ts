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
  practitionerId: z.string().optional(),
})

export const baseObservationSchema = z.object({
  loincCode: z.string(),
  valueQuantity: z.number(),
})

export const baseReportSchema = z.object({
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
  fullName: z.string().min(3, translateFunction("validation.fullNameMin")).max(255).trim(),
  birthDate: z.string().min(10, translateFunction("validation.birthDateReq")).refine(isPastDate, translateFunction("validation.birthDatePast")),
  documentId: z.string().min(11, translateFunction("validation.documentMin")).refine(cpfValidation, translateFunction("validation.documentInvalid")),
  phoneNumber: z.string().regex(/^\(\d{2}\) \d{4,5}-\d{4}$/, translateFunction("validation.phoneFormat")),
})

export const getNewEncounterSchema = (translateFunction: (key: string) => string) => z.object({
  reasonDisplay: z.string().min(3, translateFunction("validation.reasonMin")),
  practitionerId: z.string().optional(),
})

export const getNewObservationSchema = (translateFunction: (key: string) => string) => z.object({
  loincCode: z.string().min(1, translateFunction("validation.loincReq")),
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
  reportDisplay: z.string().min(3, translateFunction("validation.reportTitleMin")),
  conclusion: z.string().min(5, translateFunction("validation.conclusionMin")),
})

export const getNewConditionSchema = (translateFunction: (key: string) => string) => z.object({
  icd10Code: z.string().min(3, translateFunction("validation.icdCodeMin")).refine(isValidICD10, translateFunction("validation.icdFormat")),
  codeDisplay: z.string().min(3, translateFunction("validation.icdDisplayMin")),
})

export const getNewAllergySchema = (translateFunction: (key: string) => string) => z.object({
  allergenCode: z.string().min(3, translateFunction("validation.allergenCodeMin")),
  allergenDisplay: z.string().min(3, translateFunction("validation.allergenDisplayMin")),
  reaction: z.string().min(3, translateFunction("validation.reactionMin")),
})

export const getNewMedicationSchema = (translateFunction: (key: string) => string) => z.object({
  medicationDisplay: z.string().min(3, translateFunction("validation.medicationMin")),
  dosageInstruction: z.string().min(3, translateFunction("validation.dosageMin")),
})
