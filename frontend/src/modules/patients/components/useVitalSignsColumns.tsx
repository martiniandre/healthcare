import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { Activity } from "lucide-react"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { findVitalSignDisplay } from "./vitalSignDisplay"
import { VitalSignValueDisplay } from "./VitalSignValueDisplay"
import { formatDateTime } from "../../../shared/utils/dates"
import type { Observation } from "../types"

const columnHelper = createColumnHelper<Observation>()

export const useVitalSignsColumns = (): ColumnDef<Observation>[] => {
  const { t } = useTranslation("patients")

  return useMemo(
    () =>
      [
        columnHelper.accessor("code_display", {
          id: "code_display",
          header: t("details.vitalsCard.display"),
          cell: (info) => {
            const observation = info.row.original
            const displayMetadata = findVitalSignDisplay(observation.loinc_code)
            const IconComponent = displayMetadata?.IconComponent ?? Activity
            return (
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg border ${displayMetadata?.iconClassName ?? "bg-blue-50 border-blue-100 text-blue-600"}`}>
                  <IconComponent className="w-4 h-4" />
                </div>
                <span className="text-sm font-bold text-gray-800 block">
                  {displayMetadata?.labelKey ? t(displayMetadata.labelKey) : info.getValue()}
                </span>
              </div>
            )
          },
        }),
        columnHelper.accessor("loinc_code", {
          id: "loinc_code",
          header: t("details.vitalsCard.code"),
          cell: (info) => <span className="text-xs font-mono text-gray-500">{info.getValue()}</span>,
        }),
        columnHelper.accessor("value_quantity", {
          id: "value_quantity",
          header: t("details.vitalsCard.value"),
          cell: (info) => {
            const observation = info.row.original
            return (
              <VitalSignValueDisplay
                notPerformed={observation.not_performed}
                valueQuantity={info.getValue()}
                valueUnit={observation.value_unit}
              />
            )
          },
        }),
        columnHelper.accessor("created_at", {
          id: "created_at",
          header: t("details.vitalsCard.date"),
          cell: (info) => (
            <span className="text-xs text-gray-500 font-semibold">{formatDateTime(info.getValue())}</span>
          ),
        }),
      ] as ColumnDef<Observation>[],
    [t]
  )
}
