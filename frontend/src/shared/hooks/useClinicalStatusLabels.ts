import { useCallback } from "react"
import { useTranslation } from "react-i18next"

import type {
  AllergyClinicalStatus,
  CardiacCondition,
  ConditionClinicalStatus,
  DicomModality,
  DiagnosticReportStatus,
  EncounterStatus,
  ImagingStudyStatus,
  MedicationRequestStatus,
} from "../types"

type ClinicalStatusSet =
  | "encounter"
  | "condition"
  | "medicationRequest"
  | "diagnosticReport"
  | "imagingStudy"
  | "allergy"
  | "dicomModality"
  | "cardiacCondition"

export interface ClinicalStatusLabels {
  encounterStatus: (status: EncounterStatus) => string
  conditionStatus: (status: ConditionClinicalStatus) => string
  medicationRequestStatus: (status: MedicationRequestStatus) => string
  diagnosticReportStatus: (status: DiagnosticReportStatus) => string
  imagingStudyStatus: (status: ImagingStudyStatus) => string
  allergyClinicalStatus: (status: AllergyClinicalStatus) => string
  dicomModality: (modality: DicomModality) => string
  cardiacConditionLabel: (condition: CardiacCondition) => string
  notInformed: () => string
}

export const useClinicalStatusLabels = (): ClinicalStatusLabels => {
  const { t } = useTranslation("clinicalStatus")

  const resolveStatus = useCallback(
    (statusSet: ClinicalStatusSet, status: string): string =>
      t(`${statusSet}.${status}`, { defaultValue: status }),
    [t],
  )

  return {
    encounterStatus: (status) => resolveStatus("encounter", status),
    conditionStatus: (status) => resolveStatus("condition", status),
    medicationRequestStatus: (status) => resolveStatus("medicationRequest", status),
    diagnosticReportStatus: (status) => resolveStatus("diagnosticReport", status),
    imagingStudyStatus: (status) => resolveStatus("imagingStudy", status),
    allergyClinicalStatus: (status) => resolveStatus("allergy", status),
    dicomModality: (modality) => resolveStatus("dicomModality", modality),
    cardiacConditionLabel: (condition) => resolveStatus("cardiacCondition", condition),
    notInformed: () => t("notInformed"),
  }
}
