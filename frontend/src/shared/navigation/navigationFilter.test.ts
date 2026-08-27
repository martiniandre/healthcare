import { describe, it, expect } from "vitest"
import { navigationGroups } from "./navigationConfig"
import {
  isNavigationItemVisible,
  isNavigationItemActive,
  getVisibleNavigationGroups,
} from "./navigationFilter"

const findByKey = (sourceKey: string) =>
  navigationGroups.find((navigationGroup) => navigationGroup.key === sourceKey)!

describe("isNavigationItemVisible", () => {
  it("should show staff items to every staff role", () => {
    const patientsItem = findByKey("groups.care").items[0]
    for (const staffRole of ["DOCTOR", "NURSE", "RECEPTION", "ADMIN"]) {
      expect(isNavigationItemVisible(patientsItem, staffRole)).toBe(true)
    }
  })

  it("should hide staff items from patients", () => {
    const patientsItem = findByKey("groups.care").items[0]
    expect(isNavigationItemVisible(patientsItem, "PATIENT")).toBe(false)
  })

  it("should show the portal only to patients", () => {
    const portalItem = findByKey("groups.patientAccess").items[0]
    expect(isNavigationItemVisible(portalItem, "PATIENT")).toBe(true)
    expect(isNavigationItemVisible(portalItem, "DOCTOR")).toBe(false)
    expect(isNavigationItemVisible(portalItem, null)).toBe(false)
  })

  it("should show admin-only items only to admins", () => {
    const auditLogsItem = findByKey("groups.administration").items[1]
    expect(isNavigationItemVisible(auditLogsItem, "ADMIN")).toBe(true)
    expect(isNavigationItemVisible(auditLogsItem, "DOCTOR")).toBe(false)
    expect(isNavigationItemVisible(auditLogsItem, "PATIENT")).toBe(false)
  })

  it("should keep disabled items visible", () => {
    const settingsItem = findByKey("groups.administration").items[2]
    expect(isNavigationItemVisible(settingsItem, "DOCTOR")).toBe(true)
    expect(isNavigationItemVisible(settingsItem, "ADMIN")).toBe(true)
  })
})

describe("getVisibleNavigationGroups", () => {
  it("should return all populated groups for an admin", () => {
    const visibleGroups = getVisibleNavigationGroups("ADMIN")
    const visibleItemKeys = visibleGroups.flatMap((navigationGroup) =>
      navigationGroup.items.map((navigationItem) => navigationItem.key)
    )
    expect(visibleItemKeys).toContain("auditLogs")
    expect(visibleItemKeys).not.toContain("portal")
    expect(visibleGroups.length).toBe(3)
  })

  it("should hide the admin-only item for non-admin staff", () => {
    const visibleGroups = getVisibleNavigationGroups("DOCTOR")
    const visibleItemKeys = visibleGroups.flatMap((navigationGroup) =>
      navigationGroup.items.map((navigationItem) => navigationItem.key)
    )
    expect(visibleItemKeys).toContain("patients")
    expect(visibleItemKeys).toContain("settings")
    expect(visibleItemKeys).not.toContain("auditLogs")
    expect(visibleItemKeys).not.toContain("portal")
  })

  it("should leave only the patient access group for patients", () => {
    const visibleGroups = getVisibleNavigationGroups("PATIENT")
    expect(visibleGroups.length).toBe(1)
    expect(visibleGroups[0].key).toBe("groups.patientAccess")
    expect(visibleGroups[0].items[0].key).toBe("portal")
  })

  it("should drop groups that become empty after filtering", () => {
    const visibleGroups = getVisibleNavigationGroups("PATIENT")
    expect(visibleGroups.map((navigationGroup) => navigationGroup.key)).toEqual([
      "groups.patientAccess",
    ])
  })
})

describe("isNavigationItemActive", () => {
  it("should flag the home patients item for the root path", () => {
    const patientsItem = findByKey("groups.care").items[0]
    expect(isNavigationItemActive(patientsItem, "/")).toBe(true)
  })

  it("should flag the home patients item on patient detail routes", () => {
    const patientsItem = findByKey("groups.care").items[0]
    expect(isNavigationItemActive(patientsItem, "/patients/fhir-pat-1")).toBe(true)
    expect(isNavigationItemActive(patientsItem, "/patients/fhir-pat-1/encounters/enc-1")).toBe(true)
  })

  it("should match exact paths", () => {
    const dashboardItem = findByKey("groups.intelligence").items[0]
    expect(isNavigationItemActive(dashboardItem, "/dashboard")).toBe(true)
    expect(isNavigationItemActive(dashboardItem, "/analytics")).toBe(false)
  })

  it("should match nested paths under an item prefix", () => {
    const portalItem = findByKey("groups.patientAccess").items[0]
    expect(isNavigationItemActive(portalItem, "/portal/reports")).toBe(true)
  })

  it("should not flag items for unrelated routes", () => {
    const analyticsItem = findByKey("groups.intelligence").items[2]
    expect(isNavigationItemActive(analyticsItem, "/patients")).toBe(false)
    expect(isNavigationItemActive(analyticsItem, "/portal")).toBe(false)
    expect(isNavigationItemActive(analyticsItem, "/staff")).toBe(false)
  })

  it("should flag an item when another route shares its path prefix", () => {
    const analyticsItem = findByKey("groups.intelligence").items[2]
    expect(isNavigationItemActive(analyticsItem, "/analytics-service")).toBe(true)
  })
})