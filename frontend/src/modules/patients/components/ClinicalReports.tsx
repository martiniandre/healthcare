import { useState } from "react"
import { FileText, Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { ReportModal } from "./modals/ReportModal"
import { ReportVersionsModal } from "./modals/ReportVersionsModal"
import { ReportDetailsModal } from "./modals/ReportDetailsModal"
import { useClinicalReportsColumns } from "./useClinicalReportsColumns"
import { useDiagnosticReportsQuery, useCreateDiagnosticReportMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import type { DiagnosticReport } from "../types"
import type { NewReportFormData } from "../patient_schemas"

interface ClinicalReportsProps {
  patientId: string
  encounterId: string
}

export default function ClinicalReports({ patientId, encounterId }: ClinicalReportsProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [versionsReport, setVersionsReport] = useState<DiagnosticReport | null>(null)
  const [detailsReport, setDetailsReport] = useState<DiagnosticReport | null>(null)
  const columns = useClinicalReportsColumns({
    onOpenDetails: (report) => setDetailsReport(report),
    onOpenVersions: (report) => setVersionsReport(report),
  })
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
