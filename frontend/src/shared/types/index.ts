export const CardiacCondition = {
  Normal: "Normal",
  Bradycardia: "Bradicardia",
  Tachycardia: "Taquicardia",
  CardiacArrest: "Parada Cardíaca"
} as const

export type CardiacCondition = typeof CardiacCondition[keyof typeof CardiacCondition]

export const BedStatus = {
  Normal: "normal",
  Warning: "warning",
  Danger: "danger"
} as const

export type BedStatus = typeof BedStatus[keyof typeof BedStatus]

export const StaffRole = {
  Doctor: "Médico",
  Nurse: "Enfermeiro",
  Receptionist: "Recepção",
  Admin: "Admin"
} as const

export type StaffRole = typeof StaffRole[keyof typeof StaffRole]

export const StaffStatus = {
  OnDuty: "Plantonista",
  OffDuty: "Fora de Escala"
} as const

export type StaffStatus = typeof StaffStatus[keyof typeof StaffStatus]

export interface AuthResponseDto {
  userId: string
  role: string
  email: string
}

export const EncounterStatus = {
  Planned: "planned",
  Arrived: "arrived",
  Triaged: "triaged",
  InProgress: "in-progress",
  OnLeave: "onleave",
  Finished: "finished",
  Cancelled: "cancelled",
  EnteredInError: "entered-in-error",
  Unknown: "unknown"
} as const

export type EncounterStatus = typeof EncounterStatus[keyof typeof EncounterStatus]

export const ConditionClinicalStatus = {
  Active: "active",
  Recurrence: "recurrence",
  Relapse: "relapse",
  Inactive: "inactive",
  Remission: "remission",
  Resolved: "resolved"
} as const

export type ConditionClinicalStatus = typeof ConditionClinicalStatus[keyof typeof ConditionClinicalStatus]

export const DiagnosticReportStatus = {
  Registered: "registered",
  Partial: "partial",
  Preliminary: "preliminary",
  Final: "final",
  Amended: "amended",
  Corrected: "corrected",
  Appended: "appended",
  Cancelled: "cancelled",
  EnteredInError: "entered-in-error",
  Unknown: "unknown"
} as const

export type DiagnosticReportStatus = typeof DiagnosticReportStatus[keyof typeof DiagnosticReportStatus]

export const AllergyClinicalStatus = {
  Active: "active",
  Inactive: "inactive",
  Resolved: "resolved"
} as const

export type AllergyClinicalStatus = typeof AllergyClinicalStatus[keyof typeof AllergyClinicalStatus]

export const MedicationRequestStatus = {
  Active: "active",
  OnHold: "on-hold",
  Cancelled: "cancelled",
  Completed: "completed",
  EnteredInError: "entered-in-error",
  Stopped: "stopped",
  Draft: "draft",
  Unknown: "unknown"
} as const

export type MedicationRequestStatus = typeof MedicationRequestStatus[keyof typeof MedicationRequestStatus]

export const LoincCode = {
  HeartRate: "8867-4",
  BodyTemperature: "8310-5",
  BloodPressure: "85354-9"
} as const

export type LoincCode = typeof LoincCode[keyof typeof LoincCode]

export const DicomModality = {
  CT: "CT",
  MR: "MR",
  US: "US",
  DX: "DX",
  CR: "CR",
  XA: "XA",
  PT: "PT",
  MG: "MG",
  SR: "SR"
} as const

export type DicomModality = typeof DicomModality[keyof typeof DicomModality]

export const ImagingStudyStatus = {
  Registered: "registered",
  Available: "available",
  Cancelled: "cancelled",
  EnteredInError: "entered-in-error",
  Unknown: "unknown"
} as const

export type ImagingStudyStatus = typeof ImagingStudyStatus[keyof typeof ImagingStudyStatus]

