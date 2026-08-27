import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { CheckCircle } from "lucide-react"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { formatDateTime } from "../../../shared/utils/dates"
import type { AllergyIntolerance } from "../types"

const columnHelper = createColumnHelper<AllergyIntolerance>()

export const useClinicalAllergiesColumns = (): ColumnDef<AllergyIntolerance>[] => {
  const { t } = useTranslation("patients")

  return useMemo(
    () => [
      columnHelper.accessor("allergen_code", {
        id: "allergen_code",
        header: t("details.allergiesCard.code"),
        cell: (info) => (
          <span className="text-xs font-mono font-bold text-gray-700 bg-amber-50 border border-amber-100 px-2 py-1 rounded-md">
            {info.getValue()}
          </span>
        ),
      }),
      columnHelper.accessor("allergen_display", {
        id: "allergen_display",
        header: t("details.allergiesCard.allergen"),
        cell: (info) => <span className="text-sm font-bold text-gray-800 block">{info.getValue()}</span>,
      }),
      columnHelper.accessor("reaction", {
        id: "reaction",
        header: t("details.allergiesCard.reaction"),
        cell: (info) => <span className="text-sm font-semibold text-red-600 block">{info.getValue()}</span>,
      }),
      columnHelper.accessor("clinical_status", {
        id: "clinical_status",
        header: t("details.allergiesCard.status"),
        cell: (info) => (
          <span className="text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
            <CheckCircle className="w-3 h-3" />
            {info.getValue()}
          </span>
        ),
      }),
      columnHelper.accessor("created_at", {
        id: "created_at",
        header: t("details.allergiesCard.date"),
        cell: (info) => (
          <span className="text-xs text-gray-500 font-semibold block">
            {formatDateTime(info.getValue())}
          </span>
        ),
      }),
    ] as ColumnDef<AllergyIntolerance>[],
    [t]
  )
}
