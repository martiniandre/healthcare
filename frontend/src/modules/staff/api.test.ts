import { staffApi } from "./api"
import { StaffRole, StaffStatus } from "../../shared/types"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
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
        department_id: "department-1",
        department_name: "Cardiologia",
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
        isActive: true,
        departmentId: "department-1",
        departmentName: "Cardiologia",
        fhirResourceId: "practitioner-1",
      },
    ])
  })

  it("should map every backend role to the staff role enum", async () => {
    vi.mocked(http.get).mockResolvedValue([
      { id: "1", full_name: "A", role: "DOCTOR", email: "a@clinica.com", is_active: true },
      { id: "2", full_name: "B", role: "NURSE", email: "b@clinica.com", is_active: true },
      { id: "3", full_name: "C", role: "RECEPTION", email: "c@clinica.com", is_active: true },
      { id: "4", full_name: "D", role: "ADMIN", email: "d@clinica.com", is_active: true },
    ])

    const employees = await staffApi.listEmployees()

    expect(employees.map((employee) => employee.role)).toEqual([
      StaffRole.Doctor,
      StaffRole.Nurse,
      StaffRole.Receptionist,
      StaffRole.Admin,
    ])
  })

  it("should keep inactive employees listed so they can be reactivated", async () => {
    vi.mocked(http.get).mockResolvedValue([
      {
        id: "employee-2",
        full_name: "Enf. Carla",
        role: "NURSE",
        email: "carla@clinica.com",
        is_active: false,
        department_id: "department-2",
        department_name: "Pediatria",
      },
    ])

    const employees = await staffApi.listEmployees()

    expect(employees[0]).toMatchObject({
      role: StaffRole.Nurse,
      status: StaffStatus.OffDuty,
      isActive: false,
      departmentName: "Pediatria",
    })
  })

  it("should fall back to a dash license, dash department and empty fhir id when missing", async () => {
    vi.mocked(http.get).mockResolvedValue([
      {
        id: "employee-3",
        full_name: "Recep. Joana",
        role: "RECEPTION",
        email: "joana@clinica.com",
        is_active: true,
      },
    ])

    const employees = await staffApi.listEmployees()

    expect(employees[0]).toMatchObject({
      license: "-",
      departmentId: "",
      departmentName: "-",
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

  it("should list the department catalog from the staff endpoint", async () => {
    vi.mocked(http.get).mockResolvedValue([
      { id: "department-1", name: "Clínica Geral" },
      { id: "department-2", name: "Dermatologia" },
    ])

    const departments = await staffApi.listDepartments()

    expect(http.get).toHaveBeenCalledWith("/staff/departments")
    expect(departments).toEqual([
      { id: "department-1", name: "Clínica Geral" },
      { id: "department-2", name: "Dermatologia" },
    ])
  })

  it("should create an employee with the department in the payload mapping", async () => {
    vi.mocked(http.post).mockResolvedValue({ employee_id: "employee-1", fhir_resource_id: "practitioner-1" })

    const result = await staffApi.createEmployee({
      userId: "user-1",
      fullName: "Dr. Paulo",
      email: "paulo@clinica.com",
      role: StaffRole.Doctor,
      crmNumber: "CRM 456",
      departmentId: "department-1",
    })

    expect(http.post).toHaveBeenCalledWith("/staff/employees", {
      created_by: "user-1",
      full_name: "Dr. Paulo",
      email: "paulo@clinica.com",
      role: "DOCTOR",
      crm_number: "CRM 456",
      department_id: "department-1",
    })
    expect(result).toEqual({ employee_id: "employee-1", fhir_resource_id: "practitioner-1" })
  })

  it("should patch the employee status endpoint when deactivating", async () => {
    vi.mocked(http.patch).mockResolvedValue({ is_active: false })

    const result = await staffApi.setEmployeeStatus({ employeeId: "employee-1", isActive: false })

    expect(http.patch).toHaveBeenCalledWith("/staff/employees/employee-1/status", { is_active: false })
    expect(result).toEqual({ isActive: false })
  })

  it("should patch the employee status endpoint when reactivating", async () => {
    vi.mocked(http.patch).mockResolvedValue({ is_active: true })

    const result = await staffApi.setEmployeeStatus({ employeeId: "employee-1", isActive: true })

    expect(http.patch).toHaveBeenCalledWith("/staff/employees/employee-1/status", { is_active: true })
    expect(result).toEqual({ isActive: true })
  })

  it("should soft delete an employee through the delete endpoint", async () => {
    vi.mocked(http.delete).mockResolvedValue({ is_active: false })

    const result = await staffApi.deleteEmployee("employee-1")

    expect(http.delete).toHaveBeenCalledWith("/staff/employees/employee-1")
    expect(result).toEqual({ isActive: false })
  })
})
