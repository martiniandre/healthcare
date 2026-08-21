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

export const vitalSignsPanelFormConfig: ClinicalFormConfig<NewVitalSignsPanelFormData> = {
  titleKey: "modals.observation.title",
  confirmKey: "modals.observation.confirm",
  fields: vitalSignMetricDefinitions.map((metricDefinition) => ({
    name: metricDefinition.formFieldName,
    labelKey: metricDefinition.labelKey,
    kind: "number" as const,
  })),
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
