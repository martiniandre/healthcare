import { StaffRole } from "../../shared/types"

export const STAFF_ROLE_LABEL_KEYS: Record<StaffRole, string> = {
  [StaffRole.Doctor]: "table.roles.doctor",
  [StaffRole.Nurse]: "table.roles.nurse",
  [StaffRole.Receptionist]: "table.roles.receptionist",
  [StaffRole.Admin]: "table.roles.admin",
}

export const STAFF_ROLE_BADGE_VARIANTS: Record<StaffRole, "info" | "success" | "warning" | "secondary"> = {
  [StaffRole.Doctor]: "info",
  [StaffRole.Nurse]: "success",
  [StaffRole.Receptionist]: "warning",
  [StaffRole.Admin]: "secondary",
}
