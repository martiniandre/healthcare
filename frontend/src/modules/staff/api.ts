import { http } from "../../shared/utils/http"
import type { StaffMember, CreateEmployeePayload, CreateEmployeeResponseDto } from "./types"

export const staffApi = {
  listEmployees: async (search?: string, role?: string): Promise<StaffMember[]> => {
    const params = new URLSearchParams()
    if (search) params.append("search", search)
    if (role && role !== "All") params.append("role", role)
    const queryString = params.toString()
    return http.get<StaffMember[]>(`/staff/employees${queryString ? `?${queryString}` : ""}`)
  },

  createEmployee: async (payload: CreateEmployeePayload): Promise<CreateEmployeeResponseDto> => {
    return http.post<CreateEmployeeResponseDto>("/staff/employees", {
      user_id: payload.userId,
      full_name: payload.fullName,
      email: payload.email,
      role: payload.role,
      crm_number: payload.crmNumber,
    })
  },
}
