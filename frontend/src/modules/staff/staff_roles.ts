import { StaffRole } from "../../shared/types"

export const STAFF_ROLE_LABEL_KEYS: Record<StaffRole, string> = {
  [StaffRole.Doctor]: "table.roles.doctor",
  [StaffRole.Nurse]: "table.roles.nurse",
  [StaffRole.Receptionist]: "table.roles.receptionist",
  [StaffRole.Admin]: "table.roles.admin",
}

export const STAFF_ROLE_PILL_CLASSNAMES: Record<StaffRole, string> = {
  [StaffRole.Doctor]: "bg-sky-50 text-sky-700 border-sky-100",
  [StaffRole.Nurse]: "bg-emerald-50 text-emerald-700 border-emerald-100",
  [StaffRole.Receptionist]: "bg-amber-50 text-amber-700 border-amber-100",
  [StaffRole.Admin]: "bg-violet-50 text-violet-700 border-violet-100",
}
