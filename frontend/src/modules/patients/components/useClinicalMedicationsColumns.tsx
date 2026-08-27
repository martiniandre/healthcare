import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { CheckCircle } from "lucide-react"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { formatDateTime } from "../../../shared/utils/dates"
import type { MedicationRequest } from "../types"

const columnHelper = createColumnHelper<MedicationRequest>()

export const useClinicalMedicationsColumns = (): ColumnDef<MedicationRequest>[] => {
  const { t } = useTranslation("patients")

  return useMemo(
    () => [
      columnHelper.accessor("medication_name", {
        id: "medication_name",
        header: t("details.medicationsCard.display"),
        cell: (info) => <span className="text-sm font-extrabold text-gray-900 block">{info.getValue()}</span>,
      }),
      columnHelper.accessor("dosage_instructions", {
        id: "dosage_instructions",
        header: t("details.medicationsCard.dosage"),
        cell: (info) => (
          <span className="text-sm font-bold text-gray-800 block whitespace-pre-line">{info.getValue()}</span>
        ),
      }),
      columnHelper.accessor("status", {
        id: "status",
        header: t("details.medicationsCard.status"),
        cell: (info) => (
          <span className="text-[9px] bg-purple-50 border border-purple-100 text-purple-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
            <CheckCircle className="w-3 h-3" />
            {info.getValue()}
          </span>
        ),
      }),
      columnHelper.accessor("created_at", {
        id: "created_at",
        header: t("details.medicationsCard.date"),
        cell: (info) => (
          <span className="text-xs text-gray-500 font-semibold block">
            {formatDateTime(info.getValue())}
          </span>
        ),
      }),
    ] as ColumnDef<MedicationRequest>[],
    [t]
  )
}
