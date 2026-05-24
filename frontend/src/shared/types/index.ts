export const CardiacCondition = {
  Normal: "Normal",
  Bradycardia: "Bradicardia",
  Tachycardia: "Taquicardia",
  CardiacArrest: "Parada Cardíaca"
} as const

export type CardiacCondition = typeof CardiacCondition[keyof typeof CardiacCondition]

export const BedStatus = {
  Normal: "normal",
  Warning: "warning",
  Danger: "danger"
} as const

export type BedStatus = typeof BedStatus[keyof typeof BedStatus]

export const StaffRole = {
  Doctor: "Médico",
  Nurse: "Enfermeiro",
  Receptionist: "Recepção",
  Admin: "Admin"
} as const

export type StaffRole = typeof StaffRole[keyof typeof StaffRole]

export const StaffStatus = {
  OnDuty: "Plantonista",
  OffDuty: "Fora de Escala"
} as const

export type StaffStatus = typeof StaffStatus[keyof typeof StaffStatus]

export interface AuthResponseDto {
  userId: string
  role: string
  email: string
}
