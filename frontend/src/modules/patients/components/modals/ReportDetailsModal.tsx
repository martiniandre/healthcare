import { useTranslation } from "react-i18next"
import { FileText, CheckCircle } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../../shared/components/ui/Dialog"
import type { DiagnosticReport } from "../../types"
import { formatDateTime } from "../../../../shared/utils/dates"

interface ReportDetailsModalProps {
  isOpen: boolean
  onClose: () => void
  report: DiagnosticReport | null
}

export const ReportDetailsModal = ({
  isOpen,
  onClose,
  report,
}: ReportDetailsModalProps) => {
  const { t } = useTranslation("patients")

  if (!isOpen || !report) {
    return null
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[560px]">
        <DialogHeader>
          <DialogTitle className="text-left flex items-center gap-2">
            <FileText className="w-4 h-4 text-blue-500" />
            {t("modals.reportDetails.title")} — {report.report_display}
          </DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-3 mt-4 max-h-[60vh] overflow-y-auto">
          <div className="flex items-center justify-between gap-2 flex-wrap">
            <span className="inline-flex items-center gap-1.5 text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase">
              <CheckCircle className="w-3 h-3" />
              {report.status}
            </span>
            {report.version && (
              <span className="inline-flex items-center gap-1.5 text-[9px] bg-blue-50 border border-blue-100 text-blue-600 px-2 py-0.5 rounded font-bold uppercase">
                <FileText className="w-3 h-3" />
                v{report.version}
              </span>
            )}
            <span className="text-xs text-gray-500 font-semibold">
                {formatDateTime(report.created_at)}
            </span>
          </div>
          <div>
            <p className="text-[10px] uppercase font-bold text-gray-400 mb-1">
              {t("modals.reportDetails.examLabel")}
            </p>
            <p className="text-sm font-bold text-gray-800 bg-gray-50 border border-border p-3 rounded-lg">
              {report.report_display}
            </p>
          </div>
          <div>
            <p className="text-[10px] uppercase font-bold text-gray-400 mb-1">
              {t("modals.reportDetails.conclusionLabel")}
            </p>
            {report.conclusion ? (
              <p className="text-sm text-gray-700 leading-relaxed bg-gray-50 border border-border p-3 rounded-lg whitespace-pre-line">
                {report.conclusion}
              </p>
            ) : (
              <p className="text-sm text-gray-400 italic">
                {t("modals.reportDetails.noConclusion")}
              </p>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
