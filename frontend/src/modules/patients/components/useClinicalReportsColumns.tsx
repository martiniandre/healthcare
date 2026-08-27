import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { FileCheck, CheckCircle, History } from "lucide-react"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { Button } from "../../../shared/components/ui/Button"
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
          <button
            type="button"
            onClick={() => onOpenDetails(info.row.original)}
            title={t("details.reportsCard.viewDetails")}
            className="text-left w-full max-w-[240px] cursor-pointer"
          >
            <p className="text-xs text-gray-600 leading-relaxed bg-gray-50 border border-border p-3 rounded-lg line-clamp-2 hover:bg-gray-100 transition-colors">
              {info.getValue() || t("modals.reportDetails.noConclusion")}
            </p>
          </button>
        ),
      }),
      columnHelper.accessor("status", {
        id: "status",
        header: t("details.reportsCard.status"),
        cell: (info) => (
          <span className="text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
            <CheckCircle className="w-3 h-3" />
            {info.getValue()}
          </span>
        ),
      }),
      columnHelper.accessor("created_at", {
        id: "created_at",
        header: t("details.reportsCard.date"),
        cell: (info) => (
          <span className="text-xs text-gray-500 font-semibold block mt-1">
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
              <span className="text-xs text-gray-400">—</span>
            )
          }
          return (
            <span className="text-[9px] bg-blue-50 border border-blue-100 text-blue-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
              <FileCheck className="w-3 h-3" />
              v{versionValue}
            </span>
          )
        },
      }),
      columnHelper.display({
        id: "history",
        header: t("details.reportsCard.history"),
        cell: (info) => (
          <Button
            variantType="outline"
            onClick={() => onOpenVersions(info.row.original)}
            className="px-2 py-1 text-xs"
            title={t("details.reportsCard.history")}
          >
            <History className="w-3.5 h-3.5 text-blue-500" />
          </Button>
        ),
      }),
    ] as ColumnDef<DiagnosticReport>[],
    [t, onOpenDetails, onOpenVersions]
  )
}
