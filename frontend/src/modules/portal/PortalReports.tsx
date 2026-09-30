import { useTranslation } from "react-i18next"
import { usePortalReportsQuery } from "./queries"
import { Card } from "../../shared/components/ui/Card"
import { Loader2 } from "lucide-react"
import { formatDate } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { useClinicalStatusLabels } from "../../shared/hooks/useClinicalStatusLabels"
import type { DiagnosticReportStatus } from "../../shared/types"

const statusBadgeClass = (status: DiagnosticReportStatus) => {
  switch (status) {
    case "final":
      return "bg-green-100 text-green-700"
    case "preliminary":
    case "partial":
      return "bg-amber-100 text-amber-700"
    case "amended":
    case "corrected":
    case "appended":
      return "bg-blue-100 text-blue-700"
    case "cancelled":
      return "bg-red-100 text-red-700"
    default:
      return "bg-gray-100 text-gray-700"
  }
}

export const PortalReports = () => {
  const { t } = useTranslation("portal")
  const { data: reports, isLoading } = usePortalReportsQuery()
  const locale = useLocale()
  const { diagnosticReportStatus, notInformed } = useClinicalStatusLabels()

  if (isLoading) {
    return (
      <Card className="flex items-center justify-center min-h-[300px]">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
      </Card>
    )
  }

  if (!reports || reports.length === 0) {
    return (
      <Card className="py-16 text-center">
        <p className="text-sm text-gray-500">{t("empty.reports")}</p>
      </Card>
    )
  }

  return (
    <div className="flex flex-col gap-3">
      {reports.map((report) => (
        <div
          key={report.fhir_resource_id}
          className="bg-card border border-border rounded-xl p-5"
        >
          <div className="flex items-start justify-between mb-2">
            <div>
              <p className="text-sm font-bold text-gray-900">{report.report_display || t("fallback.exam")}</p>
              <p className="text-xs text-gray-500 mt-0.5">
                {report.issued_at
                  ? formatDate(report.issued_at, locale)
                  : notInformed()}
              </p>
            </div>
            <div className="flex items-center gap-2">
              {report.version && (
                <span className="text-[10px] bg-blue-50 border border-blue-100 text-blue-600 px-2 py-0.5 rounded font-bold uppercase">
                  v{report.version}
                </span>
              )}
              <span
                className={`text-xs font-bold px-2.5 py-1 rounded-full ${statusBadgeClass(report.status)}`}
              >
                {diagnosticReportStatus(report.status)}
              </span>
            </div>
          </div>
          {report.conclusion && (
            <p className="text-xs text-gray-600 bg-gray-50 rounded-lg p-3 mt-2">
              {report.conclusion}
            </p>
          )}
        </div>
      ))}
    </div>
  )
}
