import { useTranslation } from "react-i18next"
import { History, FileText } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../../shared/components/ui/dialog"
import { Badge } from "../../../../shared/components/ui/badge"
import { Spinner } from "../../../../shared/components/ui/spinner"
import { QueryErrorAlert } from "../QueryErrorAlert"
import { useDiagnosticReportVersionsQuery } from "../../queries"
import { formatDateTime } from "../../../../shared/utils/dates"

interface ReportVersionsModalProps {
  isOpen: boolean
  onClose: () => void
  reportFhirId: string
  reportDisplay: string
}

export const ReportVersionsModal = ({
  isOpen,
  onClose,
  reportFhirId,
  reportDisplay,
}: ReportVersionsModalProps) => {
  const { t } = useTranslation("patients")
  const { data: versions = [], isLoading, isError, refetch } = useDiagnosticReportVersionsQuery(
    isOpen ? reportFhirId : ""
  )

  if (!isOpen) {
    return null
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[560px]">
        <DialogHeader>
          <DialogTitle className="text-left flex items-center gap-2">
            <History className="w-4 h-4 text-blue-500" />
            {t("modals.reportVersions.title")} — {reportDisplay}
          </DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-3 mt-4 max-h-[60vh] overflow-y-auto">
          {isLoading ? (
            <div className="flex items-center justify-center py-10">
              <Spinner className="size-6 text-primary" />
            </div>
          ) : isError ? (
            <QueryErrorAlert onRetry={() => void refetch()} />
          ) : versions.length === 0 ? (
            <p className="text-sm text-muted-foreground text-center py-10">
              {t("modals.reportVersions.empty")}
            </p>
          ) : (
            versions.map((versionEntry) => (
              <div
                key={versionEntry.version}
                className="bg-muted border border-border rounded-xl p-4"
              >
                <div className="flex items-center justify-between mb-2">
                  <Badge variant="info" className="uppercase">
                    <FileText />
                    {t("modals.reportVersions.versionLabel")} {versionEntry.version}
                  </Badge>
                  <span className="text-xs text-muted-foreground font-semibold">
                      {formatDateTime(versionEntry.changed_at)}
                  </span>
                </div>
                {versionEntry.snapshot?.report_display && (
                  <p className="text-sm font-bold text-gray-800">
                    {versionEntry.snapshot.report_display}
                  </p>
                )}
                {versionEntry.snapshot?.conclusion && (
                  <p className="text-xs text-gray-600 leading-relaxed bg-card border border-border p-3 rounded-lg mt-2">
                    {versionEntry.snapshot.conclusion}
                  </p>
                )}
                {!versionEntry.snapshot?.conclusion && (
                  <p className="text-xs text-muted-foreground italic mt-1">
                    {t("modals.reportVersions.noConclusion")}
                  </p>
                )}
              </div>
            ))
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
