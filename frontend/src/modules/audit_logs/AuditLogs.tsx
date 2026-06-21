import * as React from "react"
import { useState } from "react"
import { useTranslation } from "react-i18next"
import { History, Search, Filter, Shield, CheckCircle2, XCircle, ChevronDown, ChevronUp, RefreshCw } from "lucide-react"
import { Card } from "../../shared/components/ui/Card"
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from "../../shared/components/ui/Table"
import { Badge } from "../../shared/components/ui/Badge"
import { Button } from "../../shared/components/ui/Button"
import { Skeleton } from "../../shared/components/ui/Skeleton"
import { useAuditLogsQuery } from "./queries"
import type { AuditLogsFilter, AuditLog } from "./types"

export const AuditLogs = () => {
  const { t } = useTranslation()
  
  const [filterAction, setFilterAction] = useState<string>("All")
  const [filterStatus, setFilterStatus] = useState<string>("All")
  const [userEmail, setUserEmail] = useState<string>("")
  const [startDate, setStartDate] = useState<string>("")
  const [endDate, setEndDate] = useState<string>("")
  
  const [expandedLogId, setExpandedLogId] = useState<string | null>(null)

  const activeFilters: AuditLogsFilter = {
    action: filterAction,
    email: userEmail,
    status: filterStatus,
    startDate: startDate,
    endDate: endDate,
  }

  const { data: auditLogsList = [], isLoading, refetch, isFetching } = useAuditLogsQuery(activeFilters)

  const toggleRowExpansion = (logId: string) => {
    setExpandedLogId(expandedLogId === logId ? null : logId)
  }

  const handleResetFilters = () => {
    setFilterAction("All")
    setFilterStatus("All")
    setUserEmail("")
    setStartDate("")
    setEndDate("")
  }

  const formatTimestamp = (timestampString: string) => {
    try {
      const dateObject = new Date(timestampString)
      return dateObject.toLocaleString()
    } catch {
      return timestampString
    }
  }

  const getRoleBadgeVariant = (userRole: string) => {
    const uppercaseRole = userRole.toUpperCase()
    if (uppercaseRole === "ADMIN") {
      return "destructive"
    }
    if (uppercaseRole === "DOCTOR") {
      return "default"
    }
    if (uppercaseRole === "NURSE") {
      return "secondary"
    }
    if (uppercaseRole === "RECEPTION") {
      return "warning"
    }
    return "outline"
  }

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full select-none relative animate-fade-in">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <div className="bg-primary/8 p-2 rounded-xl border border-primary/10">
              <History className="w-5 h-5 text-primary" />
            </div>
            <h1 className="text-xl md:text-2xl font-black text-gray-900 tracking-tight">
              {t("auditLogs.title", "Audit Logs")}
            </h1>
          </div>
          <p className="text-xs text-gray-500 font-medium mt-1">
            {t("auditLogs.subtitle", "Track system actions, page views, and API calls across the organization.")}
          </p>
        </div>

        <div className="flex gap-2">
          <Button
            variantType="outline"
            onClick={() => refetch()}
            disabled={isLoading || isFetching}
            className="flex items-center gap-2 border-border h-9"
          >
            <RefreshCw className={`w-3.5 h-3.5 ${isFetching ? "animate-spin" : ""}`} />
            {t("auditLogs.refresh", "Refresh")}
          </Button>
        </div>
      </div>

      <Card className="p-4 flex flex-col gap-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-3">
          <div className="relative">
            <Search className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              type="text"
              placeholder={t("auditLogs.filterEmailPlaceholder", "Search User Email...")}
              value={userEmail}
              onChange={(event) => setUserEmail(event.target.value)}
              className="w-full bg-white border border-border rounded-lg pl-9 pr-4 py-2 text-xs text-gray-800 placeholder-gray-400 focus:outline-none focus:border-primary/50 transition-all duration-200 h-9"
            />
          </div>

          <div>
            <select
              value={filterAction}
              onChange={(event) => setFilterAction(event.target.value)}
              className="w-full bg-white border border-border rounded-lg px-3 py-2 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200 h-9 cursor-pointer"
            >
              <option value="All">{t("auditLogs.filterAllActions", "All Actions")}</option>
              <option value="PAGE_VIEW">{t("auditLogs.actionPageView", "Page View")}</option>
              <option value="LOGIN">{t("auditLogs.actionLogin", "Login")}</option>
              <option value="LOGOUT">{t("auditLogs.actionLogout", "Logout")}</option>
              <option value="API_REQUEST">{t("auditLogs.actionApiRequest", "API Request")}</option>
            </select>
          </div>

          <div>
            <select
              value={filterStatus}
              onChange={(event) => setFilterStatus(event.target.value)}
              className="w-full bg-white border border-border rounded-lg px-3 py-2 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200 h-9 cursor-pointer"
            >
              <option value="All">{t("auditLogs.filterAllStatuses", "All Statuses")}</option>
              <option value="SUCCESS">{t("auditLogs.statusSuccess", "Success")}</option>
              <option value="FAILURE">{t("auditLogs.statusFailure", "Failure")}</option>
            </select>
          </div>

          <div>
            <input
              type="date"
              value={startDate}
              onChange={(event) => setStartDate(event.target.value)}
              className="w-full bg-white border border-border rounded-lg px-3 py-2 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200 h-9 cursor-pointer"
              placeholder={t("auditLogs.startDate", "Start Date")}
            />
          </div>

          <div>
            <input
              type="date"
              value={endDate}
              onChange={(event) => setEndDate(event.target.value)}
              className="w-full bg-white border border-border rounded-lg px-3 py-2 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200 h-9 cursor-pointer"
              placeholder={t("auditLogs.endDate", "End Date")}
            />
          </div>
        </div>

        {(filterAction !== "All" || filterStatus !== "All" || userEmail || startDate || endDate) && (
          <div className="flex justify-end">
            <button
              onClick={handleResetFilters}
              className="text-xs text-primary hover:underline font-bold flex items-center gap-1"
            >
              <Filter className="w-3 h-3" />
              {t("auditLogs.clearFilters", "Clear Filters")}
            </button>
          </div>
        )}

        <div className="overflow-x-auto border border-border rounded-xl w-full bg-white">
          <Table className="min-w-[800px]">
            <TableHeader>
              <TableRow className="hover:bg-transparent">
                <TableHead className="w-[180px]">{t("auditLogs.tableTime", "Date & Time")}</TableHead>
                <TableHead>{t("auditLogs.tableUser", "User")}</TableHead>
                <TableHead>{t("auditLogs.tableRole", "Role")}</TableHead>
                <TableHead>{t("auditLogs.tableAction", "Action")}</TableHead>
                <TableHead>{t("auditLogs.tableStatus", "Status")}</TableHead>
                <TableHead className="w-[80px]"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                Array.from({ length: 5 }).map((_, skeletonIndex) => (
                  <TableRow key={`skeleton-${skeletonIndex}`}>
                    <TableCell><Skeleton className="h-4 w-28" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-36" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-28" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-6" /></TableCell>
                  </TableRow>
                ))
              ) : auditLogsList.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-8 text-gray-500 font-semibold text-xs">
                    {t("auditLogs.noLogsFound", "No audit logs found matching the criteria.")}
                  </TableCell>
                </TableRow>
              ) : (
                auditLogsList.map((log: AuditLog) => {
                  const isExpanded = expandedLogId === log.id
                  const isSuccess = log.status.toUpperCase() === "SUCCESS"
                  return (
                    <React.Fragment key={log.id}>
                      <TableRow
                        onClick={() => toggleRowExpansion(log.id)}
                        className="cursor-pointer hover:bg-gray-50/80 transition-colors"
                      >
                        <TableCell className="text-xs text-gray-600 font-medium">
                          {formatTimestamp(log.timestamp)}
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <div className="bg-gray-100 p-1.5 rounded-lg border border-gray-200">
                              <Shield className="w-3.5 h-3.5 text-gray-500" />
                            </div>
                            <span className="font-bold text-gray-900 text-xs">{log.email || log.userId}</span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge variant={getRoleBadgeVariant(log.role)}>
                            {log.role}
                          </Badge>
                        </TableCell>
                        <TableCell className="font-mono text-xs text-gray-700">
                          {log.action}
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-1.5">
                            {isSuccess ? (
                              <CheckCircle2 className="w-4 h-4 text-green-500 shrink-0" />
                            ) : (
                              <XCircle className="w-4 h-4 text-red-500 shrink-0" />
                            )}
                            <span className={`text-xs font-bold ${isSuccess ? "text-green-600" : "text-red-600"}`}>
                              {isSuccess ? t("auditLogs.statusSuccess", "Success") : t("auditLogs.statusFailure", "Failure")}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          {isExpanded ? (
                            <ChevronUp className="w-4 h-4 text-gray-400 ml-auto" />
                          ) : (
                            <ChevronDown className="w-4 h-4 text-gray-400 ml-auto" />
                          )}
                        </TableCell>
                      </TableRow>
                      {isExpanded && (
                        <TableRow className="bg-gray-50/50 hover:bg-gray-50/50">
                          <TableCell colSpan={6} className="p-4 border-t border-border">
                            <div className="flex flex-col gap-2 bg-white p-3 rounded-lg border border-border text-xs font-mono text-gray-700 overflow-x-auto max-w-full">
                              <div className="flex flex-col gap-1">
                                <span className="font-bold text-gray-500 uppercase text-[9px] tracking-wider">
                                  {t("auditLogs.logId", "Log ID")}
                                </span>
                                <span className="text-gray-800">{log.id}</span>
                              </div>
                              <div className="flex flex-col gap-1 mt-2">
                                <span className="font-bold text-gray-500 uppercase text-[9px] tracking-wider">
                                  {t("auditLogs.details", "Details / Metadata")}
                                </span>
                                <pre className="whitespace-pre-wrap font-mono text-gray-800 bg-gray-50 p-2.5 rounded border border-border">
                                  {log.details || t("auditLogs.noDetails", "No additional details available")}
                                </pre>
                              </div>
                            </div>
                          </TableCell>
                        </TableRow>
                      )}
                    </React.Fragment>
                  )
                })
              )}
            </TableBody>
          </Table>
        </div>
      </Card>
    </div>
  )
}
