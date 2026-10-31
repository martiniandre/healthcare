import type { z } from "zod"
import { createModuleTranslator } from "../../../../shared/i18n/i18n"
import {
  getNewAllergySchema,
  getNewConditionSchema,
  getNewEncounterSchema,
  getNewMedicationSchema,
  getNewReportSchema,
  getNewVitalSignsPanelSchema,
  vitalSignMetricDefinitions,
  type NewAllergyFormData,
  type NewConditionFormData,
  type NewEncounterFormData,
  type NewMedicationFormData,
  type NewReportFormData,
  type NewVitalSignsPanelFormData,
} from "../../patient_schemas"
import type { ClinicalFormConfig, ClinicalFormOption } from "./ClinicalFormModal"

interface VitalSignFieldLayout {
  formFieldName: keyof NewVitalSignsPanelFormData
  sectionKey: string
  inputMode: "numeric" | "decimal"
}

const cardiorespiratorySectionKey = "modals.observation.sections.cardiorespiratory"
const bodySectionKey = "modals.observation.sections.body"

const vitalSignFieldLayouts: VitalSignFieldLayout[] = [
  { formFieldName: "heartRate", sectionKey: cardiorespiratorySectionKey, inputMode: "numeric" },
  { formFieldName: "oxygenSaturation", sectionKey: cardiorespiratorySectionKey, inputMode: "numeric" },
  { formFieldName: "systolicBloodPressure", sectionKey: cardiorespiratorySectionKey, inputMode: "numeric" },
  { formFieldName: "diastolicBloodPressure", sectionKey: cardiorespiratorySectionKey, inputMode: "numeric" },
  { formFieldName: "respiratoryRate", sectionKey: cardiorespiratorySectionKey, inputMode: "numeric" },
  { formFieldName: "bodyTemperature", sectionKey: bodySectionKey, inputMode: "decimal" },
  { formFieldName: "weightKg", sectionKey: bodySectionKey, inputMode: "decimal" },
  { formFieldName: "heightCm", sectionKey: bodySectionKey, inputMode: "numeric" },
]

const vitalSignFields = vitalSignFieldLayouts.map((fieldLayout) => {
  const metricDefinition = vitalSignMetricDefinitions.find(
    (definition) => definition.formFieldName === fieldLayout.formFieldName
  ) as (typeof vitalSignMetricDefinitions)[number]
  return {
    name: metricDefinition.formFieldName,
    labelKey: metricDefinition.labelKey,
    kind: "number" as const,
    required: false,
    span: "half" as const,
    sectionKey: fieldLayout.sectionKey,
    inputMode: fieldLayout.inputMode,
    range: { minimum: metricDefinition.minimumValue, maximum: metricDefinition.maximumValue },
  }
})

export const vitalSignsPanelFormConfig: ClinicalFormConfig<NewVitalSignsPanelFormData> = {
  titleKey: "modals.observation.title",
  descriptionKey: "modals.observation.description",
  confirmKey: "modals.observation.confirm",
  fields: vitalSignFields,
  schema: getNewVitalSignsPanelSchema(
    createModuleTranslator("patients")
  ) as unknown as z.ZodType<NewVitalSignsPanelFormData>,
}

const reportExamOptions: ClinicalFormOption[] = [
  { value: "58410-2", labelKey: "modals.report.examTypes.completeBloodCount" },
  { value: "2345-7", labelKey: "modals.report.examTypes.glucose" },
  { value: "24323-8", labelKey: "modals.report.examTypes.chestXray" },
  { value: "2093-3", labelKey: "modals.report.examTypes.totalCholesterol" },
]

export const reportFormConfig: ClinicalFormConfig<NewReportFormData> = {
  titleKey: "modals.report.title",
  descriptionKey: "modals.report.description",
  confirmKey: "modals.report.confirm",
  fields: [
    {
      name: "reportCode",
      labelKey: "modals.report.examType",
      placeholderKey: "modals.report.examTypePlaceholder",
      descriptionKey: "modals.report.examTypeHint",
      kind: "select",
      options: reportExamOptions,
      autofillTarget: "reportDisplay",
    },
    {
      name: "reportDisplay",
      labelKey: "modals.report.exam",
      placeholderKey: "modals.report.examPlaceholder",
      descriptionKey: "modals.report.examHint",
      kind: "text",
      maxLength: 255,
    },
    {
      name: "conclusion",
      labelKey: "modals.report.conclusion",
      placeholderKey: "modals.report.conclusionPlaceholder",
      descriptionKey: "modals.report.conclusionHint",
      kind: "textarea",
      maxLength: 2000,
    },
  ],
  schema: getNewReportSchema(createModuleTranslator("patients")),
  defaultValues: { reportCode: "", reportDisplay: "", conclusion: "" },
}

export const conditionFormConfig: ClinicalFormConfig<NewConditionFormData> = {
  titleKey: "modals.condition.title",
  descriptionKey: "modals.condition.description",
  confirmKey: "modals.condition.confirm",
  fields: [
    {
      name: "icd10Code",
      labelKey: "modals.condition.code",
      placeholderKey: "modals.condition.codePlaceholder",
      descriptionKey: "modals.condition.codeHint",
      kind: "text",
      autoCapitalize: "characters",
      spellCheck: false,
      maxLength: 10,
      span: "half",
    },
    {
      name: "codeDisplay",
      labelKey: "modals.condition.display",
      placeholderKey: "modals.condition.displayPlaceholder",
      descriptionKey: "modals.condition.displayHint",
      kind: "text",
      maxLength: 255,
    },
  ],
  schema: getNewConditionSchema(createModuleTranslator("patients")),
  defaultValues: { icd10Code: "", codeDisplay: "" },
  transformOnSubmit: (formData) => ({
    ...formData,
    icd10Code: formData.icd10Code.toUpperCase().trim(),
  }),
}

export const allergyFormConfig: ClinicalFormConfig<NewAllergyFormData> = {
  titleKey: "modals.allergy.title",
  descriptionKey: "modals.allergy.description",
  confirmKey: "modals.allergy.confirm",
  fields: [
    {
      name: "allergenDisplay",
      labelKey: "modals.allergy.display",
      placeholderKey: "modals.allergy.displayPlaceholder",
      descriptionKey: "modals.allergy.displayHint",
      kind: "text",
      maxLength: 255,
    },
    {
      name: "allergenCode",
      labelKey: "modals.allergy.code",
      placeholderKey: "modals.allergy.codePlaceholder",
      descriptionKey: "modals.allergy.codeHint",
      kind: "text",
      autoCapitalize: "none",
      spellCheck: false,
      inputMode: "text",
      maxLength: 50,
      span: "half",
    },
    {
      name: "reaction",
      labelKey: "modals.allergy.reaction",
      placeholderKey: "modals.allergy.reactionPlaceholder",
      descriptionKey: "modals.allergy.reactionHint",
      kind: "text",
      maxLength: 500,
    },
  ],
  schema: getNewAllergySchema(createModuleTranslator("patients")),
  defaultValues: { allergenCode: "", allergenDisplay: "", reaction: "" },
}

export const medicationFormConfig: ClinicalFormConfig<NewMedicationFormData> = {
  titleKey: "modals.medication.title",
  confirmKey: "modals.medication.confirm",
  fields: [
    {
      name: "medicationDisplay",
      labelKey: "modals.medication.name",
      placeholderKey: "modals.medication.namePlaceholder",
      descriptionKey: "modals.medication.nameHint",
      kind: "text",
      spellCheck: false,
      maxLength: 255,
    },
    {
      name: "dosageInstruction",
      labelKey: "modals.medication.dosage",
      placeholderKey: "modals.medication.dosagePlaceholder",
      descriptionKey: "modals.medication.dosageHint",
      kind: "textarea",
      maxLength: 1000,
    },
  ],
  schema: getNewMedicationSchema(createModuleTranslator("patients")),
  defaultValues: { medicationDisplay: "", dosageInstruction: "" },
}

export const buildEncounterFormConfig = (
  doctorOptions: ClinicalFormOption[]
): ClinicalFormConfig<NewEncounterFormData> => ({
  titleKey: "modals.encounter.title",
  descriptionKey: "modals.encounter.description",
  confirmKey: "modals.encounter.confirm",
  fields: [
    {
      name: "reasonDisplay",
      labelKey: "modals.encounter.reason",
      placeholderKey: "modals.encounter.reasonPlaceholder",
      descriptionKey: "modals.encounter.reasonHint",
      kind: "text",
      maxLength: 255,
    },
    {
      name: "practitionerId",
      labelKey: "modals.encounter.practitioner",
      placeholderKey: "modals.encounter.selectPractitioner",
      descriptionKey: "modals.encounter.practitionerHint",
      kind: "select",
      options: doctorOptions,
    },
  ],
  schema: getNewEncounterSchema(createModuleTranslator("patients")),
  defaultValues: {
    reasonDisplay: "",
    practitionerId: doctorOptions.length === 1 ? doctorOptions[0].value : "",
  },
})
