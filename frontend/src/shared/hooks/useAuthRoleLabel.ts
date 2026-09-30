import { useCallback } from "react"
import { useTranslation } from "react-i18next"

const authRoleLabelKeys: Record<string, string> = {
  ADMIN: "roles.RoleAdmin",
  DOCTOR: "roles.RoleDoctor",
  NURSE: "roles.RoleNurse",
  RECEPTION: "roles.RoleReception",
  PATIENT: "roles.RolePatient",
  RoleAdmin: "roles.RoleAdmin",
  RoleDoctor: "roles.RoleDoctor",
  RoleNurse: "roles.RoleNurse",
  RoleReception: "roles.RoleReception",
  RolePatient: "roles.RolePatient",
}

export const useAuthRoleLabel = (): ((userRole: string | null) => string) => {
  const { t } = useTranslation("header")

  return useCallback(
    (userRole: string | null): string => {
      if (!userRole) {
        return t("roles.RoleDefault")
      }
      const labelKey = authRoleLabelKeys[userRole]
      return labelKey ? t(labelKey) : t("roles.RoleDefault")
    },
    [t],
  )
}
