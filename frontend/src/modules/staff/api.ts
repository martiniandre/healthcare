import { http } from "../../shared/utils/http"
import type { StaffMember, CreateEmployeePayload, CreateEmployeeResponseDto } from "./types"
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
      department: "Geral",
      fhirResourceId: (emp.fhir_resource_id as string) || "",
    }))
  },

  createEmployee: async (payload: CreateEmployeePayload): Promise<CreateEmployeeResponseDto> => {
    return http.post<CreateEmployeeResponseDto>("/staff/employees", {
      created_by: payload.userId,
      full_name: payload.fullName,
      email: payload.email,
      role: payload.role,
      crm_number: payload.crmNumber,
    })
  },
}
