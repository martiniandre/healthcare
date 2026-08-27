import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { CheckCircle } from "lucide-react"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
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
          <span className="text-sm font-extrabold text-gray-900 bg-rose-50 border border-rose-100 px-2 py-1 rounded-md">
            {info.getValue()}
          </span>
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
          <span className="text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
            <CheckCircle className="w-3 h-3" />
            {info.getValue()}
          </span>
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
