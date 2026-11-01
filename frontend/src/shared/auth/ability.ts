import { Ability, AbilityBuilder } from "@casl/ability"
import { Capability } from "./capabilities"
import { Action, Feature } from "./types"
import type { AppAbility } from "./types"

type AbilityRule = readonly [Action, Feature]

export const capabilityAbilityRules: Readonly<Record<Capability, AbilityRule>> = {
  [Capability.PatientRead]: [Action.Read, Feature.Patient],
  [Capability.PatientCreate]: [Action.Create, Feature.Patient],
  [Capability.StaffRead]: [Action.Read, Feature.Staff],
  [Capability.StaffCreate]: [Action.Create, Feature.Staff],
  [Capability.StaffUpdateStatus]: [Action.Update, Feature.Staff],
  [Capability.StaffDelete]: [Action.Delete, Feature.Staff],
  [Capability.EncounterRead]: [Action.Read, Feature.Encounter],
  [Capability.EncounterCreate]: [Action.Create, Feature.Encounter],
  [Capability.EncounterUpdate]: [Action.Update, Feature.Encounter],
  [Capability.EncounterDelete]: [Action.Delete, Feature.Encounter],
  [Capability.ObservationRead]: [Action.Read, Feature.Observation],
  [Capability.ObservationCreate]: [Action.Create, Feature.Observation],
  [Capability.ObservationUpdate]: [Action.Update, Feature.Observation],
  [Capability.ObservationDelete]: [Action.Delete, Feature.Observation],
  [Capability.ConditionRead]: [Action.Read, Feature.Condition],
  [Capability.ConditionCreate]: [Action.Create, Feature.Condition],
  [Capability.ConditionUpdate]: [Action.Update, Feature.Condition],
  [Capability.ConditionDelete]: [Action.Delete, Feature.Condition],
  [Capability.AllergyRead]: [Action.Read, Feature.Allergy],
  [Capability.AllergyCreate]: [Action.Create, Feature.Allergy],
  [Capability.AllergyUpdate]: [Action.Update, Feature.Allergy],
  [Capability.AllergyDelete]: [Action.Delete, Feature.Allergy],
  [Capability.MedicationRead]: [Action.Read, Feature.MedicationRequest],
  [Capability.MedicationCreate]: [Action.Create, Feature.MedicationRequest],
  [Capability.ReportRead]: [Action.Read, Feature.DiagnosticReport],
  [Capability.ReportCreate]: [Action.Create, Feature.DiagnosticReport],
  [Capability.ReportUpdate]: [Action.Update, Feature.DiagnosticReport],
  [Capability.ImagingRead]: [Action.Read, Feature.ImagingStudy],
  [Capability.ImagingCreate]: [Action.Create, Feature.ImagingStudy],
  [Capability.ExamAnalysisRead]: [Action.Read, Feature.ExamAnalysis],
  [Capability.ExamAnalysisCreate]: [Action.Create, Feature.ExamAnalysis],
  [Capability.ExamAnalysisDelete]: [Action.Delete, Feature.ExamAnalysis],
  [Capability.TelemetryRead]: [Action.Read, Feature.TelemetryBed],
  [Capability.TelemetryUpdateBed]: [Action.Update, Feature.TelemetryBed],
  [Capability.AppointmentRead]: [Action.Read, Feature.Appointment],
  [Capability.AppointmentCreate]: [Action.Create, Feature.Appointment],
  [Capability.AppointmentUpdate]: [Action.Update, Feature.Appointment],
  [Capability.AppointmentCancel]: [Action.Delete, Feature.Appointment],
  [Capability.AppointmentReadOwn]: [Action.Read, Feature.OwnAppointment],
  [Capability.UnavailabilityRead]: [Action.Read, Feature.Unavailability],
  [Capability.UnavailabilityCreate]: [Action.Create, Feature.Unavailability],
  [Capability.UnavailabilityDelete]: [Action.Delete, Feature.Unavailability],
  [Capability.AnalyticsRead]: [Action.Read, Feature.Analytics],
  [Capability.AuditLogRead]: [Action.Read, Feature.AuditLog],
  [Capability.AuditLogCreate]: [Action.Create, Feature.AuditLog],
  [Capability.NotificationRead]: [Action.Read, Feature.Notification],
  [Capability.PortalRead]: [Action.Read, Feature.Portal],
  [Capability.DesignSystemRead]: [Action.Read, Feature.DesignSystem],
  [Capability.TimelineRead]: [Action.Read, Feature.Timeline],
}

const abilityRuleByCapability = new Map<string, AbilityRule>(Object.entries(capabilityAbilityRules))

export function defineAppAbility(capabilities: readonly string[]): AppAbility {
  const { can, build } = new AbilityBuilder<AppAbility>(Ability)

  for (const capability of capabilities) {
    const abilityRule = abilityRuleByCapability.get(capability)
    if (abilityRule) {
      can(abilityRule[0], abilityRule[1])
    }
  }

  return build()
}
