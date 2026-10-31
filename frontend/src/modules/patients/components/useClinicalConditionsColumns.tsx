import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { Badge } from "../../../shared/components/ui/badge"
import { StatusBadge } from "../../../shared/components/ui/status-badge"
import { resolveClinicalStatus } from "./clinicalStatus"
import { formatDateTime } from "../../../shared/utils/dates"
import type { Condition } from "../types"

const columnHelper = createColumnHelper<Condition>()

export const useClinicalConditionsColumns = (): ColumnDef<Condition>[] => {
  const { t } = useTranslation("patients")

  return useMemo(
    () => [
      columnHelper.accessor("icd10_code", {
        id: "icd10_code",
        header: t("details.conditionsCard.code"),
        cell: (info) => (
          <Badge variant="destructive" className="rounded-md">
            {info.getValue()}
          </Badge>
        ),
      }),
      columnHelper.accessor("code_display", {
        id: "code_display",
        header: t("details.conditionsCard.display"),
        cell: (info) => <span className="text-sm font-bold text-gray-800 block">{info.getValue()}</span>,
      }),
      columnHelper.accessor("clinical_status", {
        id: "clinical_status",
        header: t("details.conditionsCard.status"),
        cell: (info) => (
          <StatusBadge status={resolveClinicalStatus("condition", info.getValue())} className="uppercase">
            {info.getValue()}
          </StatusBadge>
        ),
      }),
      columnHelper.accessor("created_at", {
        id: "created_at",
        header: t("details.conditionsCard.date"),
        cell: (info) => (
          <span className="text-xs text-gray-500 font-semibold block">
            {formatDateTime(info.getValue())}
          </span>
        ),
      }),
    ] as ColumnDef<Condition>[],
    [t]
  )
}
