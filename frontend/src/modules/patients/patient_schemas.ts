import * as z from "zod"

export const newPatientSchema = z.object({
  fullName: z.string().min(3, "O nome deve ter no mínimo 3 caracteres"),
  birthDate: z.string().min(10, "A data de nascimento deve estar no formato AAAA-MM-DD"),
  documentId: z.string().min(11, "O documento (CPF/RG) deve ter no mínimo 11 caracteres"),
  phoneNumber: z.string().min(10, "O telefone deve ter no mínimo 10 dígitos"),
})

export const newEncounterSchema = z.object({
  reasonDisplay: z.string().min(3, "O motivo deve ter no mínimo 3 caracteres"),
})

export const newObservationSchema = z.object({
  loincCode: z.string().min(1, "Selecione o sinal vital"),
  valueQuantity: z.number().min(0.1, "Insira um valor numérico válido"),
})

export const newReportSchema = z.object({
  reportDisplay: z.string().min(3, "O título do laudo deve ter no mínimo 3 caracteres"),
  conclusion: z.string().min(5, "A conclusão deve ter no mínimo 5 caracteres"),
})

export const newConditionSchema = z.object({
  icd10Code: z.string().min(3, "O código CID-10 deve ter no mínimo 3 caracteres"),
  codeDisplay: z.string().min(3, "A descrição do diagnóstico deve ter no mínimo 3 caracteres"),
})

export type NewPatientFormData = z.infer<typeof newPatientSchema>
export type NewEncounterFormData = z.infer<typeof newEncounterSchema>
export type NewObservationFormData = z.infer<typeof newObservationSchema>
export type NewReportFormData = z.infer<typeof newReportSchema>
export type NewConditionFormData = z.infer<typeof newConditionSchema>

export const newAllergySchema = z.object({
  allergenCode: z.string().min(3, "O código do alérgeno deve ter no mínimo 3 caracteres"),
  allergenDisplay: z.string().min(3, "A descrição do alérgeno deve ter no mínimo 3 caracteres"),
  reaction: z.string().min(3, "A reação relatada deve ter no mínimo 3 caracteres"),
})

export type NewAllergyFormData = z.infer<typeof newAllergySchema>

export const newMedicationSchema = z.object({
  medicationDisplay: z.string().min(3, "O nome da medicação deve ter no mínimo 3 caracteres"),
  dosageInstruction: z.string().min(3, "A instrução de dosagem deve ter no mínimo 3 caracteres"),
})

export type NewMedicationFormData = z.infer<typeof newMedicationSchema>
