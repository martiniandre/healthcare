import type { z } from "zod"
import { createModuleTranslator } from "../../../../shared/i18n/i18n"
import {
  getNewAllergySchema,
  getNewConditionSchema,
  getNewEncounterSchema,
  getNewMedicationSchema,
  getNewObservationSchema,
  getNewReportSchema,
  type NewAllergyFormData,
  type NewConditionFormData,
  type NewEncounterFormData,
  type NewMedicationFormData,
  type NewObservationFormData,
  type NewReportFormData,
} from "../../patient_schemas"
import type { ClinicalFormConfig, ClinicalFormOption } from "./ClinicalFormModal"

export const LOINC_HEART_RATE = "8867-4"
export const LOINC_BODY_TEMPERATURE = "8310-5"
export const LOINC_BLOOD_PRESSURE = "85354-9"

export type SubmittedObservationFormData = NewObservationFormData & {
  codeDisplay: string
  valueUnit: string
}

interface ObservationMetricMetadata {
  display: string
  unit: string
}

const observationMetricMetadata: Record<string, ObservationMetricMetadata> = {
  [LOINC_HEART_RATE]: { display: "Frequência Cardíaca", unit: "bpm" },
  [LOINC_BODY_TEMPERATURE]: { display: "Temperatura Corporal", unit: "°C" },
  [LOINC_BLOOD_PRESSURE]: { display: "Pressão Arterial Sistólica", unit: "mmHg" },
}

const observationMetricOptions: ClinicalFormOption[] = [
  { value: LOINC_HEART_RATE, labelKey: "modals.observation.heartRate" },
  { value: LOINC_BODY_TEMPERATURE, labelKey: "modals.observation.temperature" },
  { value: LOINC_BLOOD_PRESSURE, labelKey: "modals.observation.bloodPressure" },
]

export const observationFormConfig: ClinicalFormConfig<SubmittedObservationFormData> = {
  titleKey: "modals.observation.title",
  confirmKey: "modals.observation.confirm",
  fields: [
    {
      name: "loincCode",
      labelKey: "modals.observation.selectMetric",
      placeholderKey: "modals.observation.selectMetric",
      kind: "select",
      options: observationMetricOptions,
    },
    {
      name: "valueQuantity",
      labelKey: "modals.observation.value",
      placeholderKey: "modals.observation.valuePlaceholder",
      kind: "number",
    },
  ],
  schema: getNewObservationSchema(
    createModuleTranslator("patients")
  ) as unknown as z.ZodType<SubmittedObservationFormData>,
  defaultValues: { loincCode: LOINC_HEART_RATE },
  transformOnSubmit: (formData) => {
    const metricMetadata =
      observationMetricMetadata[formData.loincCode] ?? observationMetricMetadata[LOINC_HEART_RATE]
    return {
      ...formData,
      codeDisplay: metricMetadata.display,
      valueUnit: metricMetadata.unit,
    }
  },
}

const reportExamOptions: ClinicalFormOption[] = [
  { value: "58410-2", labelKey: "modals.report.examTypes.completeBloodCount" },
  { value: "2345-7", labelKey: "modals.report.examTypes.glucose" },
  { value: "24323-8", labelKey: "modals.report.examTypes.chestXray" },
  { value: "2093-3", labelKey: "modals.report.examTypes.totalCholesterol" },
]

export const reportFormConfig: ClinicalFormConfig<NewReportFormData> = {
  titleKey: "modals.report.title",
  confirmKey: "modals.report.confirm",
  fields: [
    {
      name: "reportCode",
      labelKey: "modals.report.examType",
      placeholderKey: "modals.report.examTypePlaceholder",
      kind: "select",
      options: reportExamOptions,
    },
    {
      name: "reportDisplay",
      labelKey: "modals.report.exam",
      placeholderKey: "modals.report.examPlaceholder",
      kind: "text",
    },
    {
      name: "conclusion",
      labelKey: "modals.report.conclusion",
      placeholderKey: "modals.report.conclusionPlaceholder",
      kind: "textarea",
    },
  ],
  schema: getNewReportSchema(createModuleTranslator("patients")),
  defaultValues: { reportCode: "", reportDisplay: "", conclusion: "" },
}

export const conditionFormConfig: ClinicalFormConfig<NewConditionFormData> = {
  titleKey: "modals.condition.title",
  confirmKey: "modals.condition.confirm",
  fields: [
    {
      name: "icd10Code",
      labelKey: "modals.condition.code",
      placeholderKey: "modals.condition.codePlaceholder",
      kind: "text",
    },
    {
      name: "codeDisplay",
      labelKey: "modals.condition.display",
      placeholderKey: "modals.condition.displayPlaceholder",
      kind: "text",
    },
  ],
  schema: getNewConditionSchema(createModuleTranslator("patients")),
  transformOnSubmit: (formData) => ({
    ...formData,
    icd10Code: formData.icd10Code.toUpperCase().trim(),
  }),
  resetsAfterSubmit: true,
}

export const allergyFormConfig: ClinicalFormConfig<NewAllergyFormData> = {
  titleKey: "modals.allergy.title",
  confirmKey: "modals.allergy.confirm",
  fields: [
    {
      name: "allergenCode",
      labelKey: "modals.allergy.code",
      placeholderKey: "modals.allergy.codePlaceholder",
      kind: "text",
    },
    {
      name: "allergenDisplay",
      labelKey: "modals.allergy.display",
      placeholderKey: "modals.allergy.displayPlaceholder",
      kind: "text",
    },
    {
      name: "reaction",
      labelKey: "modals.allergy.reaction",
      placeholderKey: "modals.allergy.reactionPlaceholder",
      kind: "text",
    },
  ],
  schema: getNewAllergySchema(createModuleTranslator("patients")),
  resetsAfterSubmit: true,
}

export const medicationFormConfig: ClinicalFormConfig<NewMedicationFormData> = {
  titleKey: "modals.medication.title",
  confirmKey: "modals.medication.confirm",
  fields: [
    {
      name: "medicationDisplay",
      labelKey: "modals.medication.name",
      placeholderKey: "modals.medication.namePlaceholder",
      kind: "text",
    },
    {
      name: "dosageInstruction",
      labelKey: "modals.medication.dosage",
      placeholderKey: "modals.medication.dosagePlaceholder",
      kind: "textarea",
    },
  ],
  schema: getNewMedicationSchema(createModuleTranslator("patients")),
  resetsAfterSubmit: true,
}

export const buildEncounterFormConfig = (
  doctorOptions: ClinicalFormOption[]
): ClinicalFormConfig<NewEncounterFormData> => ({
  titleKey: "modals.encounter.title",
  confirmKey: "modals.encounter.confirm",
  fields: [
    {
      name: "reasonDisplay",
      labelKey: "modals.encounter.reason",
      placeholderKey: "modals.encounter.reasonPlaceholder",
      kind: "text",
    },
    {
      name: "practitionerId",
      labelKey: "modals.encounter.practitioner",
      placeholderKey: "modals.encounter.selectPractitioner",
      kind: "select",
      options: doctorOptions,
    },
  ],
  schema: getNewEncounterSchema(createModuleTranslator("patients")),
  defaultValues: { reasonDisplay: "", practitionerId: "" },
})
