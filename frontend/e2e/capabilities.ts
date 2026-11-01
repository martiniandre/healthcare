// Espelha backend/internal/app/policy/capability.go (capabilityRoles).
// Mantenha em sincronia: é o que o endpoint /auth/permissions devolve por papel.

export type MockRole = "ADMIN" | "DOCTOR" | "NURSE" | "RECEPTION" | "PATIENT"

const everyRole: MockRole[] = ["ADMIN", "DOCTOR", "NURSE", "RECEPTION", "PATIENT"]
const staffRoles: MockRole[] = ["ADMIN", "DOCTOR", "NURSE", "RECEPTION"]
const clinicalRoles: MockRole[] = ["ADMIN", "DOCTOR", "NURSE"]
const prescriberRoles: MockRole[] = ["ADMIN", "DOCTOR"]
const frontDeskRoles: MockRole[] = ["ADMIN", "RECEPTION"]
const adminOnlyRoles: MockRole[] = ["ADMIN"]
const patientOnly: MockRole[] = ["PATIENT"]

const capabilityRoles: Record<string, MockRole[]> = {
  "auth:permissions": everyRole,

  "patient:read": staffRoles,
  "patient:create": frontDeskRoles,

  "staff:read": staffRoles,
  "staff:create": adminOnlyRoles,
  "staff:update_status": clinicalRoles,
  "staff:delete": adminOnlyRoles,

  "encounter:read": clinicalRoles,
  "encounter:create": staffRoles,
  "encounter:update": clinicalRoles,
  "encounter:delete": clinicalRoles,

  "observation:read": clinicalRoles,
  "observation:create": clinicalRoles,
  "observation:update": clinicalRoles,
  "observation:delete": clinicalRoles,

  "condition:read": clinicalRoles,
  "condition:create": clinicalRoles,
  "condition:update": clinicalRoles,
  "condition:delete": clinicalRoles,

  "allergy:read": clinicalRoles,
  "allergy:create": clinicalRoles,
  "allergy:update": clinicalRoles,
  "allergy:delete": clinicalRoles,

  "medication:read": clinicalRoles,
  "medication:create": prescriberRoles,
  "medication:check": prescriberRoles,

  "report:read": clinicalRoles,
  "report:create": clinicalRoles,
  "report:update": clinicalRoles,

  "imaging:read": clinicalRoles,
  "imaging:create": clinicalRoles,

  "exam_analysis:read": clinicalRoles,
  "exam_analysis:create": clinicalRoles,
  "exam_analysis:delete": prescriberRoles,

  "telemetry:read": staffRoles,
  "telemetry:read_beds": clinicalRoles,
  "telemetry:unlock": clinicalRoles,
  "telemetry:update_bed": clinicalRoles,

  "appointment:read": staffRoles,
  "appointment:create": frontDeskRoles,
  "appointment:update": frontDeskRoles,
  "appointment:cancel": frontDeskRoles,
  "appointment:read_own": patientOnly,

  "unavailability:read": prescriberRoles,
  "unavailability:create": prescriberRoles,
  "unavailability:delete": prescriberRoles,

  "analytics:read": clinicalRoles,
  "timeline:read": clinicalRoles,

  "audit_log:read": adminOnlyRoles,
  "audit_log:create": adminOnlyRoles,

  "notification:read": everyRole,
  "notification:update": everyRole,
  "notification_event:read": staffRoles,

  "portal:read": patientOnly,
  "design_system:read": staffRoles,
}

export const capabilitiesForRole = (role: MockRole): string[] =>
  Object.entries(capabilityRoles)
    .filter(([, allowedRoles]) => allowedRoles.includes(role))
    .map(([capability]) => capability)
    .sort()
