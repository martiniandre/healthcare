import * as z from "zod"
import { cpfValidation, isPastDate, isValidICD10 } from "../../shared/utils/validators"

export const newPatientSchema = z.object({
  fullName: z.string().min(3, "O nome deve ter no mínimo 3 caracteres").max(255).trim(),
  birthDate: z.string().min(10, "Data de nascimento obrigatória").refine(isPastDate, "A data de nascimento deve ser no passado"),
  documentId: z.string().min(11, "CPF deve ter no mínimo 11 caracteres").refine(cpfValidation, "CPF inválido"),
  phoneNumber: z.string().regex(/^\(\d{2}\) \d{4,5}-\d{4}$/, "Formato de telefone inválido. Ex: (11) 98765-4321"),
})

export const newEncounterSchema = z.object({
  reasonDisplay: z.string().min(3, "O motivo deve ter no mínimo 3 caracteres"),
})

export const newObservationSchema = z.object({
  loincCode: z.string().min(1, "Selecione o sinal vital"),
  valueQuantity: z.number().min(0.1, "Insira um valor numérico válido"),
}).refine(
  (data) => {
    if (data.loincCode === "8867-4") {
      return data.valueQuantity >= 0 && data.valueQuantity <= 300
    }
    if (data.loincCode === "8310-5") {
      return data.valueQuantity >= 30 && data.valueQuantity <= 45
    }
    if (data.loincCode === "85354-9") {
      return data.valueQuantity >= 0 && data.valueQuantity <= 300
    }
    return true
  },
  {
    message: "Valor fora do intervalo aceitável (FC: 0-300, Temp: 30-45, PA: 0-300)",
    path: ["valueQuantity"],
  }
)

export const newReportSchema = z.object({
  reportDisplay: z.string().min(3, "O título do laudo deve ter no mínimo 3 caracteres"),
  conclusion: z.string().min(5, "A conclusão deve ter no mínimo 5 caracteres"),
})

export const newConditionSchema = z.object({
  icd10Code: z.string().min(3, "O código CID-10 deve ter no mínimo 3 caracteres").refine(isValidICD10, "Formato CID-10 inválido. Ex: I10, E11.9"),
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
