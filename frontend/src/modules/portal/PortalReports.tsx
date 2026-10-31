import { FileText } from "lucide-react"
import { useTranslation } from "react-i18next"
import { usePortalReportsQuery } from "./queries"
import { PortalSectionState } from "./components/PortalSectionState"
import { Card, CardContent } from "../../shared/components/ui/card"
import { Badge } from "../../shared/components/ui/badge"
import { StatusBadge, type ClinicalStatus } from "../../shared/components/ui/status-badge"
import { formatDate } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { useClinicalStatusLabels } from "../../shared/hooks/useClinicalStatusLabels"
import type { DiagnosticReportStatus } from "../../shared/types"

const reportStatusToClinicalStatus = (status: DiagnosticReportStatus): ClinicalStatus => {
  switch (status) {
    case "final":
      return "normal"
    case "preliminary":
    case "partial":
      return "attention"
    case "amended":
    case "corrected":
    case "appended":
      return "info"
    case "cancelled":
      return "critical"
    default:
      return "neutral"
  }
}

export const PortalReports = () => {
  const { t } = useTranslation("portal")
  const { data: reports, isLoading, isError, refetch } = usePortalReportsQuery()
  const locale = useLocale()
  const { diagnosticReportStatus, notInformed } = useClinicalStatusLabels()

  return (
    <PortalSectionState
      isLoading={isLoading}
      isError={isError}
      onRetry={() => refetch()}
      isEmpty={!reports || reports.length === 0}
      emptyTitle={t("empty.reports")}
      emptyIcon={FileText}
    >
      <div className="flex flex-col gap-3">
        {reports?.map((report) => (
          <Card key={report.fhir_resource_id} className="py-4 sm:py-5">
            <CardContent className="px-4 sm:px-5">
              <div className="mb-2 flex flex-wrap items-start justify-between gap-2">
                <div className="min-w-0 flex-1">
                  <p className="break-words text-sm font-bold text-foreground">
                    {report.report_display || t("fallback.exam")}
                  </p>
                  <p className="mt-0.5 text-xs tabular-nums text-muted-foreground">
                    {report.issued_at ? formatDate(report.issued_at, locale) : notInformed()}
                  </p>
                </div>
                <div className="flex flex-wrap items-center gap-2">
                  {report.version && (
                    <Badge variant="info" className="rounded text-[10px] font-bold uppercase tabular-nums">
                      v{report.version}
                    </Badge>
                  )}
                  <StatusBadge status={reportStatusToClinicalStatus(report.status)}>
                    {diagnosticReportStatus(report.status)}
                  </StatusBadge>
                </div>
              </div>
              {report.conclusion && (
                <p className="mt-2 break-words rounded-lg bg-muted p-3 text-xs text-muted-foreground">
                  {report.conclusion}
                </p>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </PortalSectionState>
  )
}
