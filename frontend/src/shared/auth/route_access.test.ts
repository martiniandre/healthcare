import { getHomePathForRole, isRoleAllowed, routeAllowedRoles } from "./route_access"

describe("route access", () => {
  it("should deny every route when there is no role", () => {
    expect(isRoleAllowed(null, routeAllowedRoles.patients)).toBe(false)
  })

  it("should allow only admins into audit logs", () => {
    expect(isRoleAllowed("ADMIN", routeAllowedRoles.auditLogs)).toBe(true)
    expect(isRoleAllowed("DOCTOR", routeAllowedRoles.auditLogs)).toBe(false)
  })

  it("should keep reception out of clinical-only routes", () => {
    expect(isRoleAllowed("RECEPTION", routeAllowedRoles.analytics)).toBe(false)
    expect(isRoleAllowed("RECEPTION", routeAllowedRoles.examAnalyzer)).toBe(false)
    expect(isRoleAllowed("RECEPTION", routeAllowedRoles.patients)).toBe(true)
  })

  it("should let every staff role open the design system but not patients", () => {
    expect(isRoleAllowed("RECEPTION", routeAllowedRoles.designSystem)).toBe(true)
    expect(isRoleAllowed("DOCTOR", routeAllowedRoles.designSystem)).toBe(true)
    expect(isRoleAllowed("PATIENT", routeAllowedRoles.designSystem)).toBe(false)
  })

  it("should keep patients out of every staff route and into the portal", () => {
    expect(isRoleAllowed("PATIENT", routeAllowedRoles.patients)).toBe(false)
    expect(isRoleAllowed("PATIENT", routeAllowedRoles.telemetry)).toBe(false)
    expect(isRoleAllowed("PATIENT", routeAllowedRoles.portal)).toBe(true)
  })

  it("should keep staff out of the patient portal", () => {
    expect(isRoleAllowed("NURSE", routeAllowedRoles.portal)).toBe(false)
  })

  it("should resolve the home path by role", () => {
    expect(getHomePathForRole("PATIENT")).toBe("/portal")
    expect(getHomePathForRole("DOCTOR")).toBe("/")
  })
})
