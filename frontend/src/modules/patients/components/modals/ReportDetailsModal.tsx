import { useTranslation } from "react-i18next"
import { FileText, CheckCircle } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../../shared/components/ui/dialog"
import { Badge } from "../../../../shared/components/ui/badge"
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
  const { t: translateStatus } = useTranslation("clinicalStatus")

  if (!isOpen || !report) {
    return null
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-xl">
        <DialogHeader>
          <DialogTitle className="text-left flex items-center gap-2">
            <FileText className="w-4 h-4 text-primary" />
            {t("modals.reportDetails.title")}: {report.report_display}
          </DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-3 mt-4 max-h-[60vh] overflow-y-auto">
          <div className="flex items-center justify-between gap-2 flex-wrap">
            <Badge variant="success">
              <CheckCircle />
              {translateStatus(`diagnosticReport.${report.status}`, { defaultValue: report.status })}
            </Badge>
            {report.version && (
              <Badge variant="info">
                <FileText />
                v{report.version}
              </Badge>
            )}
            <span className="text-xs text-muted-foreground font-semibold">
                {formatDateTime(report.created_at)}
            </span>
          </div>
          <div>
            <p className="text-xs font-medium text-muted-foreground mb-1">
              {t("modals.reportDetails.examLabel")}
            </p>
            <p className="text-sm font-semibold text-foreground bg-muted border border-border p-3 rounded-lg">
              {report.report_display}
            </p>
          </div>
          <div>
            <p className="text-xs font-medium text-muted-foreground mb-1">
              {t("modals.reportDetails.conclusionLabel")}
            </p>
            {report.conclusion ? (
              <p className="text-sm text-muted-foreground leading-relaxed bg-muted border border-border p-3 rounded-lg whitespace-pre-line">
                {report.conclusion}
              </p>
            ) : (
              <p className="text-sm text-muted-foreground italic">
                {t("modals.reportDetails.noConclusion")}
              </p>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
