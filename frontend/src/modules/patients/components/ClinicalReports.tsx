import { useState } from "react"
import { FileText, Plus, FileCheck, CheckCircle, History } from "lucide-react"
import { useTranslation } from "react-i18next"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { ReportModal } from "./modals/ReportModal"
import { ReportVersionsModal } from "./modals/ReportVersionsModal"
import { ReportDetailsModal } from "./modals/ReportDetailsModal"
import { useDiagnosticReportsQuery, useCreateDiagnosticReportMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import { formatDateTime } from "../../../shared/utils/dates"
import type { DiagnosticReport } from "../types"
import type { NewReportFormData } from "../patient_schemas"

interface ClinicalReportsProps {
  patientId: string
  encounterId: string
}

const columnHelper = createColumnHelper<DiagnosticReport>()

export default function ClinicalReports({ patientId, encounterId }: ClinicalReportsProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [versionsReport, setVersionsReport] = useState<DiagnosticReport | null>(null)
  const [detailsReport, setDetailsReport] = useState<DiagnosticReport | null>(null)
  const { data: reports = [] } = useDiagnosticReportsQuery(encounterId)
  const createReportMutation = useCreateDiagnosticReportMutation()

  const handleCreateReport = async (formData: NewReportFormData) => {
    try {
      await createReportMutation.mutateAsync({
        encounter_fhir_id: encounterId,
        patient_fhir_id: patientId,
        report_code: formData.reportCode,
        report_display: formData.reportDisplay,
        conclusion: formData.conclusion,
      })
      setIsModalOpen(false)
      toast.success(t("toast.reportSuccess"))
    } catch {
      toast.error(t("toast.reportError"))
    }
  }

  const columns = [
    columnHelper.accessor("report_display", {
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
      header: t("details.reportsCard.conclusion"),
      cell: (info) => (
        <button
          type="button"
          onClick={() => setDetailsReport(info.row.original)}
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
      header: t("details.reportsCard.status"),
      cell: (info) => (
        <span className="text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
          <CheckCircle className="w-3 h-3" />
          {info.getValue()}
        </span>
      ),
    }),
    columnHelper.accessor("created_at", {
      header: t("details.reportsCard.date"),
      cell: (info) => (
        <span className="text-xs text-gray-500 font-semibold block mt-1">
          {formatDateTime(info.getValue())}
        </span>
      ),
    }),
    columnHelper.accessor("version", {
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
          onClick={() => setVersionsReport(info.row.original)}
          className="px-2 py-1 text-xs"
          title={t("details.reportsCard.history")}
        >
          <History className="w-3.5 h-3.5 text-blue-500" />
        </Button>
      ),
    }),
  ] as ColumnDef<DiagnosticReport>[]

  return (
    <>
      <ClinicalTable
        title={t("details.reportsCard.title")}
        icon={<FileText className="w-4 h-4 text-blue-500" />}
        columns={columns}
        data={reports}
        isEmpty={reports.length === 0}
        emptyIcon={<FileText className="w-8 h-8 text-gray-300" />}
        emptyText={t("details.reportsCard.empty")}
        addButton={
          <Can I={Action.Create} a={Feature.DiagnosticReport}>
            <Button onClick={() => setIsModalOpen(true)} className="px-3 py-2 text-xs">
              <Plus className="w-3.5 h-3.5" />
              {t("details.reportsCard.add")}
            </Button>
          </Can>
        }
      />

      <ReportModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateReport}
        isPending={createReportMutation.isPending}
      />

      <ReportVersionsModal
        isOpen={versionsReport !== null}
        onClose={() => setVersionsReport(null)}
        reportFhirId={versionsReport?.fhir_id ?? ""}
        reportDisplay={versionsReport?.report_display ?? ""}
      />

      <ReportDetailsModal
        isOpen={detailsReport !== null}
        onClose={() => setDetailsReport(null)}
        report={detailsReport}
      />
    </>
  )
}
