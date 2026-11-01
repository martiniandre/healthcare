import { Capability } from "./capabilities"
import { getHomePathForRole, isCapabilityGranted, routeRequiredCapability } from "./route_access"
import { backendMatrixCapabilitiesByRoleFixture, type RoleFixtureName } from "../../test/backend_matrix_role_fixtures"

const grantedFor = (role: RoleFixtureName) => new Set<string>(backendMatrixCapabilitiesByRoleFixture[role])

describe("route access", () => {
  it("should deny every route when no capability is granted", () => {
    expect(isCapabilityGranted(new Set(), routeRequiredCapability.patients)).toBe(false)
  })

  it("should require the declared capability for each route", () => {
    expect(routeRequiredCapability).toEqual({
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
    })
  })

  it("should allow only admins into audit logs", () => {
    expect(isCapabilityGranted(grantedFor("ADMIN"), routeRequiredCapability.auditLogs)).toBe(true)
    expect(isCapabilityGranted(grantedFor("DOCTOR"), routeRequiredCapability.auditLogs)).toBe(false)
  })

  it("should keep reception out of clinical-only routes", () => {
    const reception = grantedFor("RECEPTION")
    expect(isCapabilityGranted(reception, routeRequiredCapability.analytics)).toBe(false)
    expect(isCapabilityGranted(reception, routeRequiredCapability.examAnalyzer)).toBe(false)
    expect(isCapabilityGranted(reception, routeRequiredCapability.patients)).toBe(true)
  })

  it("should let every staff role open the design system but not patients", () => {
    expect(isCapabilityGranted(grantedFor("RECEPTION"), routeRequiredCapability.designSystem)).toBe(true)
    expect(isCapabilityGranted(grantedFor("DOCTOR"), routeRequiredCapability.designSystem)).toBe(true)
    expect(isCapabilityGranted(grantedFor("PATIENT"), routeRequiredCapability.designSystem)).toBe(false)
  })

  it("should keep patients out of every staff route and into the portal", () => {
    const patient = grantedFor("PATIENT")
    expect(isCapabilityGranted(patient, routeRequiredCapability.patients)).toBe(false)
    expect(isCapabilityGranted(patient, routeRequiredCapability.telemetry)).toBe(false)
    expect(isCapabilityGranted(patient, routeRequiredCapability.portal)).toBe(true)
  })

  it("should keep staff out of the patient portal", () => {
    expect(isCapabilityGranted(grantedFor("NURSE"), routeRequiredCapability.portal)).toBe(false)
  })

  it("should resolve the home path from the portal capability", () => {
    expect(getHomePathForRole(grantedFor("PATIENT"))).toBe("/portal")
    expect(getHomePathForRole(grantedFor("DOCTOR"))).toBe("/")
    expect(getHomePathForRole(new Set())).toBe("/")
  })
})
