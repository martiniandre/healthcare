import { useState } from "react"
import { useTranslation } from "react-i18next"
import { History, RefreshCw } from "lucide-react"
import { Card } from "../../shared/components/ui/Card"
import { Button } from "../../shared/components/ui/Button"
import { useAuditLogsQuery } from "./queries"
import { AuditLogsFilters } from "./components/AuditLogsFilters"
import { AuditLogsTable } from "./components/AuditLogsTable"
import { AuditLogsLoadingState } from "./components/AuditLogsLoadingState"
import type { AuditLogsFilter } from "./types"

export const AuditLogs = () => {
  const { t } = useTranslation("auditLogs")
  const [filterAction, setFilterAction] = useState("All")
  const [filterStatus, setFilterStatus] = useState("All")
  const [userEmail, setUserEmail] = useState("")
  const [startDate, setStartDate] = useState("")
  const [endDate, setEndDate] = useState("")
  const [expandedLogId, setExpandedLogId] = useState<string | null>(null)

  const activeFilters: AuditLogsFilter = { action: filterAction, email: userEmail, status: filterStatus, startDate, endDate }
  const { data, isLoading, isError, refetch, isFetching } = useAuditLogsQuery(activeFilters)
  const auditLogsList = data?.audit_logs ?? []

  if (isLoading && !data) {
    return <AuditLogsLoadingState />
  }

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full select-none relative animate-fade-in">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <div className="bg-primary/8 p-2 rounded-xl border border-primary/10">
              <History className="w-5 h-5 text-primary" />
            </div>
            <h1 className="text-xl md:text-2xl font-black text-gray-900 tracking-tight">{t("title")}</h1>
          </div>
          <p className="text-xs text-gray-500 font-medium mt-1">{t("subtitle")}</p>
        </div>
        <Button variantType="outline" onClick={() => refetch()} disabled={isFetching} className="flex items-center gap-2 border-border h-9">
          <RefreshCw className={`w-3.5 h-3.5 ${isFetching ? "animate-spin" : ""}`} />
          {t("refresh")}
        </Button>
      </div>

      <Card className="p-4 flex flex-col gap-4">
        <AuditLogsFilters
          filterAction={filterAction} onFilterActionChange={setFilterAction}
          filterStatus={filterStatus} onFilterStatusChange={setFilterStatus}
          userEmail={userEmail} onUserEmailChange={setUserEmail}
          startDate={startDate} onStartDateChange={setStartDate}
          endDate={endDate} onEndDateChange={setEndDate}
          onResetFilters={() => { setFilterAction("All"); setFilterStatus("All"); setUserEmail(""); setStartDate(""); setEndDate("") }}
          isError={isError}
        />
        <AuditLogsTable
          isLoading={isLoading} auditLogsList={auditLogsList}
          expandedLogId={expandedLogId}
          onToggleRowExpansion={(logId) => setExpandedLogId(expandedLogId === logId ? null : logId)}
        />
      </Card>
    </div>
  )
}
