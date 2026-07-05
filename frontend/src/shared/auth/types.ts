export enum Action {
  Create = "create",
  Read = "read",
  Update = "update",
  Delete = "delete",
  Manage = "manage",
}

export enum Feature {
  All = "all",
  Patient = "Patient",
  Condition = "Condition",
  Allergy = "Allergy",
  Observation = "Observation",
  DiagnosticReport = "DiagnosticReport",
  MedicationRequest = "MedicationRequest",
  Encounter = "Encounter",
  TelemetryBed = "TelemetryBed",
  Staff = "Staff",
  AuditLog = "AuditLog",
  ExamAnalysis = "ExamAnalysis",
  ImagingStudy = "ImagingStudy",
  Portal = "Portal",
}

export type AppAbility = import("@casl/ability").Ability<[Action, Feature]>
