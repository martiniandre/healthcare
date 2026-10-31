import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { Power, PowerOff, Trash2, Users } from "lucide-react"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { Badge } from "@/shared/components/ui/badge"
import { StatusBadge } from "@/shared/components/ui/status-badge"
import { Button } from "@/shared/components/ui/button"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/shared/components/ui/tooltip"
import { STAFF_ROLE_LABEL_KEYS, STAFF_ROLE_BADGE_VARIANTS } from "../staff_roles"
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
                  <span className={`font-extrabold truncate ${member.isActive ? "text-foreground" : "text-muted-foreground line-through"}`}>
                    {member.fullName}
                  </span>
                  <span className="text-[10px] text-muted-foreground truncate mt-0.5">{member.email}</span>
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
              <Badge variant={STAFF_ROLE_BADGE_VARIANTS[role] ?? "secondary"} className="lowercase tracking-wide">
                {t(STAFF_ROLE_LABEL_KEYS[role] ?? "table.roles.doctor")}
              </Badge>
            )
          },
        }),
        columnHelper.accessor("license", {
          id: "license",
          header: t("table.license"),
          cell: (info) => <span className="font-mono text-xs text-muted-foreground">{info.getValue()}</span>,
        }),
        columnHelper.accessor("departmentName", {
          id: "departmentName",
          header: t("table.department"),
          cell: (info) => <span className="text-xs font-medium text-muted-foreground">{info.getValue()}</span>,
        }),
        columnHelper.accessor("isActive", {
          id: "isActive",
          header: t("table.status"),
          cell: (info) => {
            const isActive = info.getValue()
            return (
              <StatusBadge status={isActive ? "normal" : "neutral"}>
                {isActive ? t("table.active") : t("table.inactive")}
              </StatusBadge>
            )
          },
        }),
        columnHelper.display({
          id: "actions",
          header: t("table.actions"),
          cell: (info) => {
            const member = info.row.original
            const isRowPending = isPendingEmployeeId === member.id
            const toggleLabel = member.isActive ? t("table.disable") : t("table.enable")
            const showToggle = canManageStatus || canDelete
            if (!showToggle) {
              return null
            }

            return (
              <div className="flex items-center gap-1.5">
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      type="button"
                      variant="outline"
                      size="icon-sm"
                      onClick={() => onToggleStatus(member)}
                      disabled={!canManageStatus || isRowPending}
                      aria-label={toggleLabel}
                    >
                      {member.isActive ? <PowerOff /> : <Power />}
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>{toggleLabel}</TooltipContent>
                </Tooltip>
                {canDelete && (
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        type="button"
                        variant="outline"
                        size="icon-sm"
                        onClick={() => onDelete(member)}
                        disabled={isRowPending}
                        aria-label={t("table.delete")}
                        className="hover:text-destructive"
                      >
                        <Trash2 />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>{t("table.delete")}</TooltipContent>
                  </Tooltip>
                )}
              </div>
            )
          },
        }),
      ] as ColumnDef<StaffMember>[],
    [t, onToggleStatus, onDelete, canManageStatus, canDelete, isPendingEmployeeId]
  )
}
