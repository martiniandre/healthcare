import { describe, it, expect, vi } from "vitest"
import { renderHook } from "@testing-library/react"
import { StaffRole, StaffStatus } from "../../../shared/types"
import type { StaffMember } from "../types"
import { STAFF_ROLE_LABEL_KEYS, STAFF_ROLE_PILL_CLASSNAMES } from "../staff_roles"
import { useStaffColumns } from "./useStaffColumns"

const mockTranslateFunction = (key: string) => key

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: mockTranslateFunction,
  }),
}))

const buildMember = (overrides: Partial<StaffMember> = {}): StaffMember => ({
  id: "employee-1",
  fullName: "Dra. Marina Alves",
  role: StaffRole.Doctor,
  license: "CRM 123",
  email: "marina@clinica.com",
  status: StaffStatus.OnDuty,
  isActive: true,
  departmentId: "department-1",
  departmentName: "Cardiologia",
  fhirResourceId: "",
  ...overrides,
})

const stableToggleStatus = vi.fn()
const stableDelete = vi.fn()

const renderColumns = (overrides: Partial<Parameters<typeof useStaffColumns>[0]> = {}) =>
  renderHook(() =>
    useStaffColumns({
      onToggleStatus: stableToggleStatus,
      onDelete: stableDelete,
      canManageStatus: true,
      canDelete: true,
      ...overrides,
    })
  )

describe("useStaffColumns", () => {
  it("should return six column definitions", () => {
    const { result } = renderColumns()

    expect(result.current).toHaveLength(6)
  })

  it("should translate every column header", () => {
    const { result } = renderColumns()

    expect(result.current.map((column) => column.header)).toEqual([
      "table.professional",
      "table.role",
      "table.license",
      "table.department",
      "table.status",
      "table.actions",
    ])
  })

  it("should bind each column to its accessor", () => {
    const { result } = renderColumns()

    expect(result.current.map((column) => column.id)).toEqual([
      "fullName",
      "role",
      "license",
      "departmentName",
      "isActive",
      "actions",
    ])
  })

  it("should map every role enum to a dedicated lowercase translation key", () => {
    expect(STAFF_ROLE_LABEL_KEYS).toEqual({
      [StaffRole.Doctor]: "table.roles.doctor",
      [StaffRole.Nurse]: "table.roles.nurse",
      [StaffRole.Receptionist]: "table.roles.receptionist",
      [StaffRole.Admin]: "table.roles.admin",
    })
    Object.values(STAFF_ROLE_LABEL_KEYS).forEach((labelKey) => {
      expect(labelKey).toBe("table.roles." + labelKey.split(".")[2])
    })
  })

  it("should define a pill style for every role enum", () => {
    Object.values(StaffRole).forEach((role) => {
      expect(STAFF_ROLE_PILL_CLASSNAMES[role]).toBeTruthy()
    })
  })

  it("should bind the department and status columns to their row accessor keys", () => {
    const { result } = renderColumns()
    const departmentColumn = result.current.find((column) => column.id === "departmentName")
    const activeColumn = result.current.find((column) => column.id === "isActive")

    expect(departmentColumn?.accessorKey).toBe("departmentName")
    expect(activeColumn?.accessorKey).toBe("isActive")
    expect(buildMember().departmentName).toBe("Cardiologia")
    expect(buildMember().isActive).toBe(true)
  })

  it("should keep the same array reference across rerenders with stable dependencies", () => {
    const { result, rerender } = renderColumns()
    const firstReference = result.current

    rerender()

    expect(result.current).toBe(firstReference)
  })

  it("should rebuild the columns when the row pending identifier changes", () => {
    const { result, rerender } = renderHook(
      (pendingEmployeeId: string | null) =>
        useStaffColumns({
          onToggleStatus: stableToggleStatus,
          onDelete: stableDelete,
          canManageStatus: true,
          canDelete: true,
          isPendingEmployeeId: pendingEmployeeId,
        }),
      { initialProps: null as string | null }
    )
    const firstReference = result.current

    rerender("employee-1")

    expect(result.current).not.toBe(firstReference)
  })
})
