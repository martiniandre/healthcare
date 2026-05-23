import { http } from "../../shared/utils/http"
import type { StaffMember, CreateEmployeePayload, CreateEmployeeResponseDto } from "./types"

export const staffApi = {
  listEmployees: async (): Promise<StaffMember[]> => {
    return http.get<StaffMember[]>("/staff/employees")
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
