import { Ability, AbilityBuilder } from "@casl/ability"
import { Action, Feature } from "./types"
import type { AppAbility } from "./types"

export function defineAppAbility(userRole: string | null): AppAbility {
  const { can, cannot, build } = new AbilityBuilder<AppAbility>(Ability)

  if (userRole === "ADMIN") {
    can(Action.Manage, Feature.All)
    return build()
  }

  if (userRole === "DOCTOR") {
    can([Action.Create, Action.Read, Action.Update, Action.Delete], [
      Feature.Condition, Feature.Allergy, Feature.Observation, Feature.Encounter,
    ])
    can([Action.Create, Action.Read, Action.Update], [Feature.DiagnosticReport, Feature.TelemetryBed])
    can([Action.Create, Action.Read], [Feature.MedicationRequest])
    can([Action.Create, Action.Read, Action.Delete], [Feature.ExamAnalysis])
    can(Action.Read, [Feature.Patient, Feature.Staff, Feature.AuditLog, Feature.ImagingStudy])
    cannot(Action.Create, [Feature.Patient, Feature.Staff])
    return build()
  }

  if (userRole === "NURSE") {
    can([Action.Create, Action.Read], [
      Feature.Condition, Feature.Allergy, Feature.Observation, Feature.DiagnosticReport, Feature.TelemetryBed,
    ])
    can([Action.Create, Action.Read, Action.Update], [Feature.Encounter])
    can(Action.Read, [
      Feature.Patient, Feature.Staff, Feature.AuditLog, Feature.ImagingStudy, Feature.MedicationRequest, Feature.ExamAnalysis,
    ])
    cannot(Action.Create, [Feature.Patient, Feature.Staff, Feature.MedicationRequest])
    cannot(Action.Delete, [Feature.ExamAnalysis])
    return build()
  }

  if (userRole === "RECEPTION") {
    can([Action.Create, Action.Read], [Feature.Patient, Feature.Encounter])
    can(Action.Read, [Feature.Staff])
    return build()
  }

  if (userRole === "PATIENT") {
    can(Action.Read, [Feature.Portal])
    return build()
  }

  return build()
}
