import { useState } from "react"
import { FileText, Plus, FileCheck, CheckCircle } from "lucide-react"
import { useTranslation } from "react-i18next"
import { createColumnHelper } from "@tanstack/react-table"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { ReportModal } from "./modals/ReportModal"
import { useDiagnosticReportsQuery, useCreateDiagnosticReportMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import type { DiagnosticReport } from "../types"

interface ClinicalReportsProps {
  patientId: string
  encounterId: string
}

const columnHelper = createColumnHelper<DiagnosticReport>()

export default function ClinicalReports({ patientId, encounterId }: ClinicalReportsProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const { data: reports = [] } = useDiagnosticReportsQuery(encounterId)
  const createReportMutation = useCreateDiagnosticReportMutation()

  const handleCreateReport = async (formData: { reportDisplay: string; conclusion: string }) => {
    try {
      await createReportMutation.mutateAsync({
        encounter_fhir_id: encounterId,
        patient_fhir_id: patientId,
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
        <p className="text-xs text-gray-600 leading-relaxed bg-gray-50 border border-border p-3 rounded-lg max-h-24 overflow-y-auto">
          {info.getValue()}
        </p>
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
          {new Date(info.getValue()).toLocaleString()}
        </span>
      ),
    }),
  ]

  return (
    <>
      <ClinicalTable
        title={t("details.reportsCard.title")}
        icon={<FileText className="w-4 h-4 text-emerald-500" />}
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
    </>
  )
}
