import { http } from "../../shared/utils/http"
import type {
  StaffMember,
  Department,
  CreateEmployeePayload,
  CreateEmployeeResponseDto,
  SetEmployeeStatusPayload,
} from "./types"
import { StaffRole, StaffStatus } from "../../shared/types"

const mapRole = (role: string): StaffRole => {
  switch (role) {
    case 'DOCTOR': return StaffRole.Doctor
    case 'NURSE': return StaffRole.Nurse
    case 'RECEPTION': return StaffRole.Receptionist
    case 'ADMIN': return StaffRole.Admin
    default: return role as StaffRole
  }
}

export const staffApi = {
  listEmployees: async (search?: string, role?: string): Promise<StaffMember[]> => {
    const params = new URLSearchParams()
    if (search) params.append("search", search)
    if (role && role !== "All") params.append("role", role)
    const queryString = params.toString()

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const response = await http.get<any[]>(`/staff/employees${queryString ? `?${queryString}` : ""}`)

    return response.map((emp: Record<string, unknown>) => ({
      id: emp.id as string,
      fullName: emp.full_name as string,
      role: mapRole(emp.role as string),
      license: (emp.crm_number as string) || "-",
      email: emp.email as string,
      status: emp.is_active ? StaffStatus.OnDuty : StaffStatus.OffDuty,
      isActive: Boolean(emp.is_active),
      departmentId: (emp.department_id as string) || "",
      departmentName: (emp.department_name as string) || "-",
      fhirResourceId: (emp.fhir_resource_id as string) || "",
    }))
  },

  listDepartments: async (): Promise<Department[]> => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const response = await http.get<any[]>("/staff/departments")

    return response.map((department: Record<string, unknown>) => ({
      id: department.id as string,
      name: department.name as string,
    }))
  },

  createEmployee: async (payload: CreateEmployeePayload): Promise<CreateEmployeeResponseDto> => {
    return http.post<CreateEmployeeResponseDto>("/staff/employees", {
      created_by: payload.userId,
      full_name: payload.fullName,
      email: payload.email,
      role: payload.role,
      crm_number: payload.crmNumber,
      department_id: payload.departmentId,
    })
  },

  setEmployeeStatus: async (payload: SetEmployeeStatusPayload): Promise<{ isActive: boolean }> => {
    const response = await http.patch<{ is_active: boolean }>(
      `/staff/employees/${payload.employeeId}/status`,
      { is_active: payload.isActive },
    )
    return { isActive: response.is_active }
  },

  deleteEmployee: async (employeeId: string): Promise<{ isActive: boolean }> => {
    const response = await http.delete<{ is_active: boolean }>(`/staff/employees/${employeeId}`)
    return { isActive: response.is_active }
  },
}
