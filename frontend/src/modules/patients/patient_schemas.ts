import * as z from "zod"
import { cpfValidation, isPastDate, isValidICD10 } from "../../shared/utils/validators"

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

export type NewVitalSignsPanelFormData = {
  heartRate?: number
  bodyTemperature?: number
  systolicBloodPressure?: number
  diastolicBloodPressure?: number
  oxygenSaturation?: number
  respiratoryRate?: number
  weightKg?: number
  heightCm?: number
}

export interface VitalSignMetricDefinition {
  formFieldName: keyof NewVitalSignsPanelFormData
  loincCode: string
  minimumValue: number
  maximumValue: number
  labelKey: string
}

export const vitalSignMetricDefinitions: VitalSignMetricDefinition[] = [
  { formFieldName: "heartRate", loincCode: "8867-4", minimumValue: 0, maximumValue: 300, labelKey: "modals.observation.panel.heartRate" },
  { formFieldName: "bodyTemperature", loincCode: "8310-5", minimumValue: 30, maximumValue: 45, labelKey: "modals.observation.panel.bodyTemperature" },
  { formFieldName: "systolicBloodPressure", loincCode: "8480-6", minimumValue: 0, maximumValue: 300, labelKey: "modals.observation.panel.systolicBloodPressure" },
  { formFieldName: "diastolicBloodPressure", loincCode: "8462-4", minimumValue: 0, maximumValue: 300, labelKey: "modals.observation.panel.diastolicBloodPressure" },
  { formFieldName: "oxygenSaturation", loincCode: "59408-5", minimumValue: 0, maximumValue: 100, labelKey: "modals.observation.panel.oxygenSaturation" },
  { formFieldName: "respiratoryRate", loincCode: "9279-1", minimumValue: 0, maximumValue: 60, labelKey: "modals.observation.panel.respiratoryRate" },
  { formFieldName: "weightKg", loincCode: "29463-7", minimumValue: 0, maximumValue: 500, labelKey: "modals.observation.panel.weightKg" },
  { formFieldName: "heightCm", loincCode: "8302-2", minimumValue: 0, maximumValue: 250, labelKey: "modals.observation.panel.heightCm" },
]

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

const buildVitalSignMetricField = (
  metricDefinition: VitalSignMetricDefinition,
  translateFunction: (key: string) => string
) => {
  const rangeErrorMessage = `${translateFunction(metricDefinition.labelKey)}: ${translateFunction("validation.vitalSignsRange")}`
  return z.preprocess(
    (rawValue) => (typeof rawValue === "number" && Number.isNaN(rawValue) ? undefined : rawValue),
    z
      .number()
      .min(metricDefinition.minimumValue, rangeErrorMessage)
      .max(metricDefinition.maximumValue, rangeErrorMessage)
      .optional()
  )
}

export const getNewVitalSignsPanelSchema = (translateFunction: (key: string) => string) =>
  z
    .object(
      Object.fromEntries(
        vitalSignMetricDefinitions.map((metricDefinition) => [
          metricDefinition.formFieldName,
          buildVitalSignMetricField(metricDefinition, translateFunction),
        ])
      )
    )
    .superRefine((values, context) => {
      const panelValues = values as unknown as NewVitalSignsPanelFormData
      const hasMeasuredMetric = vitalSignMetricDefinitions.some((metricDefinition) => {
        const metricValue = panelValues[metricDefinition.formFieldName]
        return typeof metricValue === "number" && Number.isFinite(metricValue)
      })
      if (!hasMeasuredMetric) {
        context.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["heartRate"],
          message: translateFunction("validation.vitalSignsAtLeastOne"),
        })
      }
    }) as unknown as z.ZodType<NewVitalSignsPanelFormData>

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
