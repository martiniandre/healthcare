export const AppRole = {
  Admin: "ADMIN",
  Doctor: "DOCTOR",
  Nurse: "NURSE",
  Reception: "RECEPTION",
  Patient: "PATIENT",
} as const

export type AppRole = typeof AppRole[keyof typeof AppRole]

const STAFF_ROLES: readonly AppRole[] = [AppRole.Admin, AppRole.Doctor, AppRole.Nurse, AppRole.Reception]
const CLINICAL_ROLES: readonly AppRole[] = [AppRole.Admin, AppRole.Doctor, AppRole.Nurse]

export const routeAllowedRoles = {
  patients: STAFF_ROLES,
  dashboard: CLINICAL_ROLES,
  schedule: STAFF_ROLES,
  telemetry: STAFF_ROLES,
  examAnalyzer: CLINICAL_ROLES,
  analytics: CLINICAL_ROLES,
  staff: STAFF_ROLES,
  auditLogs: [AppRole.Admin],
  portal: [AppRole.Patient],
  designSystem: STAFF_ROLES,
} as const satisfies Record<string, readonly AppRole[]>

export type RouteAccessKey = keyof typeof routeAllowedRoles

export function isRoleAllowed(userRole: string | null, allowedRoles: readonly string[]): boolean {
  return userRole !== null && allowedRoles.includes(userRole)
}

export function getHomePathForRole(userRole: string | null): string {
  return userRole === AppRole.Patient ? "/portal" : "/"
}
