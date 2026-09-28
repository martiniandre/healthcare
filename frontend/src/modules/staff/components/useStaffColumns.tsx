import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { Power, PowerOff, Trash2, Users } from "lucide-react"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { STAFF_ROLE_LABEL_KEYS, STAFF_ROLE_PILL_CLASSNAMES } from "../staff_roles"
import type { StaffMember } from "../types"

const columnHelper = createColumnHelper<StaffMember>()

interface UseStaffColumnsOptions {
  onToggleStatus: (member: StaffMember) => void
  onDelete: (member: StaffMember) => void
  canManageStatus: boolean
  canDelete: boolean
  isPendingEmployeeId?: string | null
}

export const useStaffColumns = ({
  onToggleStatus,
  onDelete,
  canManageStatus,
  canDelete,
  isPendingEmployeeId = null,
}: UseStaffColumnsOptions): ColumnDef<StaffMember>[] => {
  const { t } = useTranslation("staff")

  return useMemo(
    () =>
      [
        columnHelper.accessor("fullName", {
          id: "fullName",
          header: t("table.professional"),
          cell: (info) => {
            const member = info.row.original
            return (
              <div className="flex items-center gap-3">
                <div className="bg-primary/8 p-2 rounded-lg border border-primary/10">
                  <Users className="w-4 h-4 text-primary" />
                </div>
                <div className="flex flex-col min-w-0">
                  <span className={`font-extrabold truncate ${member.isActive ? "text-gray-900" : "text-gray-400 line-through"}`}>
                    {member.fullName}
                  </span>
                  <span className="text-[10px] text-gray-500 truncate mt-0.5">{member.email}</span>
                </div>
              </div>
            )
          },
        }),
        columnHelper.accessor("role", {
          id: "role",
          header: t("table.role"),
          cell: (info) => {
            const role = info.getValue()
            return (
              <span
                className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-[11px] font-semibold lowercase tracking-wide ${STAFF_ROLE_PILL_CLASSNAMES[role] ?? "bg-gray-50 text-gray-600 border-gray-100"}`}
              >
                {t(STAFF_ROLE_LABEL_KEYS[role] ?? "table.roles.doctor")}
              </span>
            )
          },
        }),
        columnHelper.accessor("license", {
          id: "license",
          header: t("table.license"),
          cell: (info) => <span className="font-mono text-xs text-gray-600">{info.getValue()}</span>,
        }),
        columnHelper.accessor("departmentName", {
          id: "departmentName",
          header: t("table.department"),
          cell: (info) => <span className="text-xs font-medium text-gray-600">{info.getValue()}</span>,
        }),
        columnHelper.accessor("isActive", {
          id: "isActive",
          header: t("table.status"),
          cell: (info) => {
            const isActive = info.getValue()
            return (
              <span
                className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-[11px] font-bold ${
                  isActive ? "bg-primary/8 text-primary" : "bg-gray-100 text-gray-500"
                }`}
              >
                <span className={`w-1.5 h-1.5 rounded-full ${isActive ? "bg-primary" : "bg-gray-400"}`} />
                {isActive ? t("table.active") : t("table.inactive")}
              </span>
            )
          },
        }),
        columnHelper.display({
          id: "actions",
          header: t("table.actions"),
          cell: (info) => {
            const member = info.row.original
            const isRowPending = isPendingEmployeeId === member.id
            const showToggle = canManageStatus || canDelete
            if (!showToggle) {
              return null
            }

            return (
              <div className="flex items-center gap-1.5">
                <button
                  type="button"
                  onClick={() => onToggleStatus(member)}
                  disabled={!canManageStatus || isRowPending}
                  title={member.isActive ? t("table.disable") : t("table.enable")}
                  aria-label={member.isActive ? t("table.disable") : t("table.enable")}
                  className="p-1.5 rounded-lg border border-border bg-card text-gray-500 hover:text-primary hover:border-primary/40 transition-colors disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:text-gray-500 disabled:hover:border-border"
                >
                  {member.isActive ? <PowerOff className="w-3.5 h-3.5" /> : <Power className="w-3.5 h-3.5" />}
                </button>
                {canDelete && (
                  <button
                    type="button"
                    onClick={() => onDelete(member)}
                    disabled={isRowPending}
                    title={t("table.delete")}
                    aria-label={t("table.delete")}
                    className="p-1.5 rounded-lg border border-border bg-card text-gray-400 hover:text-red-600 hover:border-red-200 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                )}
              </div>
            )
          },
        }),
      ] as ColumnDef<StaffMember>[],
    [t, onToggleStatus, onDelete, canManageStatus, canDelete, isPendingEmployeeId]
  )
}
