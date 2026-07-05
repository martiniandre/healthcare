import { StaffRole, StaffStatus } from "../../shared/types"

export interface StaffMember {
  id: string
  fullName: string
  role: StaffRole
  license: string
  email: string
  status: StaffStatus
  department: string
  fhirResourceId: string
}

export interface CreateEmployeePayload {
  userId: string
  fullName: string
  email: string
  role: StaffRole
  crmNumber: string
}

export interface CreateEmployeeResponseDto {
  employeeId: string
}
