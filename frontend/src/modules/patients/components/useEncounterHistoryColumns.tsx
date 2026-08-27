import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { CheckCircle, XCircle, ExternalLink } from "lucide-react"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { Button } from "../../../shared/components/ui/Button"
import { formatDateTime } from "../../../shared/utils/dates"
import type { Encounter } from "../types"

interface EncounterHistoryColumnsDependencies {
  onFinishEncounter: (encounter: Encounter) => void
  onCancelEncounter: (encounter: Encounter) => void
  onOpenEncounter: (encounter: Encounter) => void
  isActionPending: boolean
}

const columnHelper = createColumnHelper<Encounter>()

export const useEncounterHistoryColumns = ({
  onFinishEncounter,
  onCancelEncounter,
  onOpenEncounter,
  isActionPending,
}: EncounterHistoryColumnsDependencies): ColumnDef<Encounter>[] => {
  const { t } = useTranslation("patients")

  return useMemo(
    () => [
      columnHelper.accessor("reason_display", {
        id: "reason_display",
        header: t("details.encountersCard.reason"),
        cell: (info) => <span className="text-sm font-bold text-gray-800 block">{info.getValue()}</span>,
      }),
      columnHelper.accessor("status", {
        id: "status",
        header: t("details.encountersCard.status"),
        cell: (info) => (
          <span className="text-[10px] bg-gray-100 text-gray-600 px-2.5 py-1 rounded font-bold uppercase tracking-wider">
            {info.getValue()}
          </span>
        ),
      }),
      columnHelper.accessor("created_at", {
        id: "created_at",
        header: t("details.encountersCard.date"),
        cell: (info) => (
          <span className="text-xs text-gray-400 font-semibold">
            {formatDateTime(info.getValue())}
          </span>
        ),
      }),
      columnHelper.display({
        id: "actions",
        header: t("details.encountersCard.action"),
        cell: (info) => {
          const encounter = info.row.original
          return (
            <div className="flex justify-end gap-1.5 pr-6">
              {encounter.status === "in-progress" && (
                <>
                  <Button
                    variantType="outline"
                    onClick={() => onFinishEncounter(encounter)}
                    disabled={isActionPending}
                    className="px-2.5 py-1 text-[10px] font-bold gap-1"
                  >
                    <CheckCircle className="w-3 h-3" />
                    {t("details.encountersCard.finish")}
                  </Button>
                  <Button
                    variantType="danger"
                    onClick={() => onCancelEncounter(encounter)}
                    disabled={isActionPending}
                    className="px-2.5 py-1 text-[10px] font-bold gap-1"
                  >
                    <XCircle className="w-3 h-3" />
                    {t("details.encountersCard.cancel")}
                  </Button>
                </>
              )}
              <Button
                variantType="outline"
                onClick={() => onOpenEncounter(encounter)}
                className="px-2.5 py-1 text-[10px] font-bold gap-1"
              >
                <ExternalLink className="w-3 h-3" />
                {t("details.openEncounter")}
              </Button>
            </div>
          )
        },
      }),
    ] as ColumnDef<Encounter>[],
    [t, isActionPending, onFinishEncounter, onCancelEncounter, onOpenEncounter]
  )
}
