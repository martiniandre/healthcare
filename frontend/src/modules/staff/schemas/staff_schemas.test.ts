import { describe, it, expect } from "vitest"
import { StaffRole } from "../../../shared/types"
import { staffFormSchema } from "./staff_schemas"

const validFormData = {
  fullName: "Dra. Marina Alves",
  role: StaffRole.Doctor,
  license: "CRM-SP 12345",
  email: "marina@clinica.com",
  departmentId: "0f9c2b1a-1b2c-4d5e-8f90-1234567890ab",
}

const expectFieldError = (data: Record<string, unknown>, field: string, expectedMessage: string) => {
  const parsedFormData = staffFormSchema.safeParse(data)

  expect(parsedFormData.success).toBe(false)
  if (!parsedFormData.success) {
    const fieldIssue = parsedFormData.error.issues.find((issue) => issue.path[0] === field)
    expect(fieldIssue?.message).toBe(expectedMessage)
  }
}

describe("staffFormSchema", () => {
  it("should accept a fully filled registration", () => {
    expect(staffFormSchema.safeParse(validFormData).success).toBe(true)
  })

  it("should require the full name", () => {
    expectFieldError({ ...validFormData, fullName: "  " }, "fullName", "validation.nameMinLength")
  })

  it("should require a minimum length on the full name", () => {
    expectFieldError({ ...validFormData, fullName: "Dr" }, "fullName", "validation.nameMinLength")
  })

  it("should require the professional license", () => {
    expectFieldError({ ...validFormData, license: "" }, "license", "validation.licenseRequired")
  })

  it("should reject an invalid license format", () => {
    expectFieldError({ ...validFormData, license: "registro" }, "license", "validation.licenseFormat")
  })

  it("should require the corporate email", () => {
    expectFieldError({ ...validFormData, email: "" }, "email", "validation.emailRequired")
  })

  it("should reject an invalid email", () => {
    expectFieldError({ ...validFormData, email: "marina" }, "email", "validation.emailInvalid")
  })

  it("should require the department", () => {
    expectFieldError({ ...validFormData, departmentId: "" }, "departmentId", "validation.departmentRequired")
  })

  it("should reject a role outside the allowed enum", () => {
    expectFieldError({ ...validFormData, role: "PHARMACIST" }, "role", "validation.roleRequired")
  })

  it("should accept a COREN license for nurses", () => {
    expect(
      staffFormSchema.safeParse({ ...validFormData, role: StaffRole.Nurse, license: "COREN-12345" }).success
    ).toBe(true)
  })
})
