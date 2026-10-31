import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { Badge } from "../../../shared/components/ui/badge"
import { StatusBadge } from "../../../shared/components/ui/status-badge"
import { resolveClinicalStatus } from "./clinicalStatus"
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
          <Badge variant="warning" className="font-mono rounded-md">
            {info.getValue()}
          </Badge>
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
        cell: (info) => <span className="text-sm font-semibold text-destructive block">{info.getValue()}</span>,
      }),
      columnHelper.accessor("clinical_status", {
        id: "clinical_status",
        header: t("details.allergiesCard.status"),
        cell: (info) => (
          <StatusBadge status={resolveClinicalStatus("allergy", info.getValue())} className="uppercase">
            {info.getValue()}
          </StatusBadge>
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
