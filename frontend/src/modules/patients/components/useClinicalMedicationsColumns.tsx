import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { StatusBadge } from "../../../shared/components/ui/status-badge"
import { resolveClinicalStatus } from "./clinicalStatus"
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
          <StatusBadge status={resolveClinicalStatus("medication", info.getValue())} className="uppercase">
            {info.getValue()}
          </StatusBadge>
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
