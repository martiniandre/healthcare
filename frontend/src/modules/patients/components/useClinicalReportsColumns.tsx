import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { FileCheck, History } from "lucide-react"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { Badge } from "../../../shared/components/ui/badge"
import { StatusBadge } from "../../../shared/components/ui/status-badge"
import { resolveClinicalStatus } from "./clinicalStatus"
import { Button } from "../../../shared/components/ui/button"
import { Tooltip, TooltipContent, TooltipTrigger } from "../../../shared/components/ui/tooltip"
import { formatDateTime } from "../../../shared/utils/dates"
import type { DiagnosticReport } from "../types"

interface ClinicalReportsColumnsDependencies {
  onOpenDetails: (report: DiagnosticReport) => void
  onOpenVersions: (report: DiagnosticReport) => void
}

const columnHelper = createColumnHelper<DiagnosticReport>()

export const useClinicalReportsColumns = ({
  onOpenDetails,
  onOpenVersions,
}: ClinicalReportsColumnsDependencies): ColumnDef<DiagnosticReport>[] => {
  const { t } = useTranslation("patients")

  return useMemo(
    () => [
      columnHelper.accessor("report_display", {
        id: "report_display",
        header: t("details.reportsCard.display"),
        cell: (info) => (
          <div className="flex items-center gap-3">
            <div className="bg-emerald-50 border border-emerald-100 p-2 rounded-lg text-emerald-600">
              <FileCheck className="w-4 h-4" />
            </div>
            <span className="text-sm font-bold text-gray-800 block">{info.getValue()}</span>
          </div>
        ),
      }),
      columnHelper.accessor("conclusion", {
        id: "conclusion",
        header: t("details.reportsCard.conclusion"),
        cell: (info) => (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                type="button"
                variant="ghost"
                onClick={() => onOpenDetails(info.row.original)}
                aria-label={t("details.reportsCard.viewDetails")}
                className="h-auto w-full max-w-[240px] justify-start whitespace-normal bg-muted p-3 text-left text-xs font-normal leading-relaxed text-muted-foreground hover:bg-accent"
              >
                <span className="line-clamp-2">
                  {info.getValue() || t("modals.reportDetails.noConclusion")}
                </span>
              </Button>
            </TooltipTrigger>
            <TooltipContent>{t("details.reportsCard.viewDetails")}</TooltipContent>
          </Tooltip>
        ),
      }),
      columnHelper.accessor("status", {
        id: "status",
        header: t("details.reportsCard.status"),
        cell: (info) => (
          <StatusBadge status={resolveClinicalStatus("report", info.getValue())} className="uppercase">
            {info.getValue()}
          </StatusBadge>
        ),
      }),
      columnHelper.accessor("created_at", {
        id: "created_at",
        header: t("details.reportsCard.date"),
        cell: (info) => (
          <span className="text-xs text-muted-foreground font-semibold block mt-1">
            {formatDateTime(info.getValue())}
          </span>
        ),
      }),
      columnHelper.accessor("version", {
        id: "version",
        header: t("details.reportsCard.version"),
        cell: (info) => {
          const versionValue = info.getValue()
          if (!versionValue) {
            return (
              <span className="text-xs text-muted-foreground">—</span>
            )
          }
          return (
            <Badge variant="info" className="uppercase">
              <FileCheck />
              v{versionValue}
            </Badge>
          )
        },
      }),
      columnHelper.display({
        id: "history",
        header: t("details.reportsCard.history"),
        cell: (info) => (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="outline"
                size="icon-sm"
                onClick={() => onOpenVersions(info.row.original)}
                aria-label={t("details.reportsCard.history")}
              >
                <History className="text-primary" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>{t("details.reportsCard.history")}</TooltipContent>
          </Tooltip>
        ),
      }),
    ] as ColumnDef<DiagnosticReport>[],
    [t, onOpenDetails, onOpenVersions]
  )
}
