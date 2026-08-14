import { staffApi } from "./api"
import { StaffRole, StaffStatus } from "../../shared/types"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("staffApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should list employees without query params when no filters are provided", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await staffApi.listEmployees()

    expect(http.get).toHaveBeenCalledWith("/staff/employees")
  })

  it("should map backend employee rows to the staff member shape", async () => {
    vi.mocked(http.get).mockResolvedValue([
      {
        id: "employee-1",
        full_name: "Dra. Marina",
        role: "DOCTOR",
        crm_number: "CRM 123",
        email: "marina@clinica.com",
        is_active: true,
        fhir_resource_id: "practitioner-1",
      },
    ])

    const employees = await staffApi.listEmployees()

    expect(employees).toEqual([
      {
        id: "employee-1",
        fullName: "Dra. Marina",
        role: StaffRole.Doctor,
        license: "CRM 123",
        email: "marina@clinica.com",
        status: StaffStatus.OnDuty,
        department: "Geral",
        fhirResourceId: "practitioner-1",
      },
    ])
  })

  it("should fall back to a dash license and empty fhir id when missing", async () => {
    vi.mocked(http.get).mockResolvedValue([
      {
        id: "employee-2",
        full_name: "Enf. Carla",
        role: "NURSE",
        email: "carla@clinica.com",
        is_active: false,
      },
    ])

    const employees = await staffApi.listEmployees()

    expect(employees[0]).toMatchObject({
      role: StaffRole.Nurse,
      license: "-",
      status: StaffStatus.OffDuty,
      fhirResourceId: "",
    })
  })

  it("should pass search and role filters as query params", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await staffApi.listEmployees("Marina", "DOCTOR")

    expect(http.get).toHaveBeenCalledWith("/staff/employees?search=Marina&role=DOCTOR")
  })

  it("should skip the role filter when it is All", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await staffApi.listEmployees("", "All")

    expect(http.get).toHaveBeenCalledWith("/staff/employees")
  })

  it("should create an employee with the payload mapping", async () => {
    vi.mocked(http.post).mockResolvedValue({ employee_id: "employee-1", fhir_resource_id: "practitioner-1" })

    const result = await staffApi.createEmployee({
      userId: "user-1",
      fullName: "Dr. Paulo",
      email: "paulo@clinica.com",
      role: StaffRole.Doctor,
      crmNumber: "CRM 456",
    })

    expect(http.post).toHaveBeenCalledWith("/staff/employees", {
      created_by: "user-1",
      full_name: "Dr. Paulo",
      email: "paulo@clinica.com",
      role: "DOCTOR",
      crm_number: "CRM 456",
    })
    expect(result).toEqual({ employee_id: "employee-1", fhir_resource_id: "practitioner-1" })
  })
})
