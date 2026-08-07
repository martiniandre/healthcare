export const Action = {
  Create: "create",
  Read: "read",
  Update: "update",
  Delete: "delete",
  Manage: "manage",
} as const

export type Action = typeof Action[keyof typeof Action]

export const Feature = {
  All: "all",
  Patient: "Patient",
  Condition: "Condition",
  Allergy: "Allergy",
  Observation: "Observation",
  DiagnosticReport: "DiagnosticReport",
  MedicationRequest: "MedicationRequest",
  Encounter: "Encounter",
  TelemetryBed: "TelemetryBed",
  Staff: "Staff",
  AuditLog: "AuditLog",
  ExamAnalysis: "ExamAnalysis",
  ImagingStudy: "ImagingStudy",
  Portal: "Portal",
  Appointment: "Appointment",
} as const

export type Feature = typeof Feature[keyof typeof Feature]

export type AppAbility = import("@casl/ability").Ability<[Action, Feature]>
