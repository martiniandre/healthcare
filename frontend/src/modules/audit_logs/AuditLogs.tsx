import { useState } from "react"
import { useTranslation } from "react-i18next"
import { RefreshCw } from "lucide-react"
import { Card, CardContent } from "../../shared/components/ui/card"
import { Button } from "../../shared/components/ui/button"
import { Spinner } from "../../shared/components/ui/spinner"
import { PageContainer, PageHeader } from "../../shared/components/ui/PageContainer"
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
    <PageContainer className="select-none relative">
      <PageHeader
        title={t("title")}
        actions={
          <Button variant="outline" onClick={() => refetch()} disabled={isFetching}>
            {isFetching ? <Spinner /> : <RefreshCw aria-hidden="true" />}
            {t("refresh")}
          </Button>
        }
      />

      <Card className="py-4">
        <CardContent className="flex flex-col gap-4 px-4">
        <AuditLogsFilters
          filterAction={filterAction} onFilterActionChange={setFilterAction}
          filterStatus={filterStatus} onFilterStatusChange={setFilterStatus}
          userEmail={userEmail} onUserEmailChange={setUserEmail}
          startDate={startDate} onStartDateChange={setStartDate}
          endDate={endDate} onEndDateChange={setEndDate}
          onResetFilters={() => { setFilterAction("All"); setFilterStatus("All"); setUserEmail(""); setStartDate(""); setEndDate("") }}
          isError={isError}
          onRetry={() => refetch()}
        />
        <AuditLogsTable
          isLoading={isLoading} isError={isError} auditLogsList={auditLogsList}
          expandedLogId={expandedLogId}
          onToggleRowExpansion={(logId) => setExpandedLogId(expandedLogId === logId ? null : logId)}
        />
        </CardContent>
      </Card>
    </PageContainer>
  )
}
