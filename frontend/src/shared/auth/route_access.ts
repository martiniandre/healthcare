import { Capability } from "./capabilities"

export const routeRequiredCapability = {
  patients: Capability.PatientRead,
  dashboard: Capability.AnalyticsRead,
  schedule: Capability.AppointmentRead,
  telemetry: Capability.TelemetryRead,
  examAnalyzer: Capability.ExamAnalysisRead,
  analytics: Capability.AnalyticsRead,
  staff: Capability.StaffRead,
  auditLogs: Capability.AuditLogRead,
  portal: Capability.PortalRead,
  designSystem: Capability.DesignSystemRead,
} as const satisfies Record<string, Capability>

export type RouteAccessKey = keyof typeof routeRequiredCapability

export function isCapabilityGranted(
  grantedCapabilities: ReadonlySet<string>,
  requiredCapability: Capability,
): boolean {
  return grantedCapabilities.has(requiredCapability)
}

export function getHomePathForRole(grantedCapabilities: ReadonlySet<string>): string {
  return grantedCapabilities.has(Capability.PortalRead) ? "/portal" : "/"
}
