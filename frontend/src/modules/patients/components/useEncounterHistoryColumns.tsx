import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { CheckCircle, XCircle, ExternalLink } from "lucide-react"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { Button } from "../../../shared/components/ui/button"
import { StatusBadge } from "../../../shared/components/ui/status-badge"
import { resolveClinicalStatus } from "./clinicalStatus"
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
          <StatusBadge status={resolveClinicalStatus("encounter", info.getValue())} className="uppercase tracking-wider">
            {info.getValue()}
          </StatusBadge>
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
                    variant="outline"
                    size="xs"
                    onClick={() => onFinishEncounter(encounter)}
                    disabled={isActionPending}
                  >
                    <CheckCircle />
                    {t("details.encountersCard.finish")}
                  </Button>
                  <Button
                    variant="destructive"
                    size="xs"
                    onClick={() => onCancelEncounter(encounter)}
                    disabled={isActionPending}
                  >
                    <XCircle />
                    {t("details.encountersCard.cancel")}
                  </Button>
                </>
              )}
              <Button
                variant="outline"
                size="xs"
                onClick={() => onOpenEncounter(encounter)}
              >
                <ExternalLink />
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
