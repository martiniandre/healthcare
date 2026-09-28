import { StaffRole, StaffStatus } from "../../shared/types"

export interface StaffMember {
  id: string
  fullName: string
  role: StaffRole
  license: string
  email: string
  status: StaffStatus
  isActive: boolean
  departmentId: string
  departmentName: string
  fhirResourceId: string
}

export interface Department {
  id: string
  name: string
}

export interface CreateEmployeePayload {
  userId: string
  fullName: string
  email: string
  role: StaffRole
  crmNumber: string
  departmentId: string
}

export interface CreateEmployeeResponseDto {
  employeeId: string
}

export interface SetEmployeeStatusPayload {
  employeeId: string
  isActive: boolean
}
