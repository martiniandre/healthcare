import { useTranslation } from "react-i18next"
import { Shield, ChevronDown, ChevronUp, History } from "lucide-react"
import * as React from "react"
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from "../../../shared/components/ui/table"
import { Badge } from "../../../shared/components/ui/badge"
import { StatusBadge } from "../../../shared/components/ui/status-badge"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import { useLayoutStore } from "../../../shared/store/layout_store"
import { useLocale } from "../../../shared/hooks/useLocale"
import { TableDensityToggle } from "@/shared/components/clinical/TableDensityToggle"
import { Skeleton } from "../../../shared/components/ui/skeleton"
import { formatDateTime } from "../../../shared/utils/dates"
import type { AuditLog } from "../types"

interface AuditLogsTableProps {
  isLoading: boolean
  isError?: boolean
  auditLogsList: AuditLog[]
  expandedLogId: string | null
  onToggleRowExpansion: (logId: string) => void
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

export const AuditLogsTable = ({
  isLoading,
  isError = false,
  auditLogsList,
  expandedLogId,
  onToggleRowExpansion,
}: AuditLogsTableProps) => {
  const { t } = useTranslation("auditLogs")
  const locale = useLocale()
  const tableDensity = useLayoutStore((state) => state.tableDensity)

  if (isError) {
    return null
  }

  if (!isLoading && auditLogsList.length === 0) {
    return <EmptyState icon={History} title={t("noLogsFound")} className="py-12" />
  }

  return (
    <div className="flex flex-col gap-2">
    <div className="flex justify-end">
      <TableDensityToggle />
    </div>
    <div className="w-full overflow-hidden rounded-xl border border-border bg-card">
      <Table
        className="min-w-[800px]"
        density={tableDensity}
        containerClassName="max-h-[560px] overflow-y-auto"
      >
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            <TableHead className="sticky top-0 z-10 min-w-40 bg-card">{t("tableTime")}</TableHead>
            <TableHead className="sticky top-0 z-10 bg-card">{t("tableUser")}</TableHead>
            <TableHead className="sticky top-0 z-10 bg-card">{t("tableRole")}</TableHead>
            <TableHead className="sticky top-0 z-10 bg-card">{t("tableAction")}</TableHead>
            <TableHead className="sticky top-0 z-10 bg-card">{t("tableResource")}</TableHead>
            <TableHead className="sticky top-0 z-10 bg-card">{t("tableStatus")}</TableHead>
            <TableHead className="sticky top-0 z-10 w-12 bg-card"><span className="sr-only">{t("tableDetails")}</span></TableHead>
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
                <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                <TableCell><Skeleton className="h-4 w-6" /></TableCell>
              </TableRow>
            ))
          ) : (
            auditLogsList.map((log: AuditLog) => {
              const isExpanded = expandedLogId === log.id
              const isSuccess = log.access_granted
              return (
                <React.Fragment key={log.id}>
                  <TableRow
                    onClick={() => onToggleRowExpansion(log.id)}
                    className="cursor-pointer hover:bg-muted/50 transition-colors"
                  >
                    <TableCell className="text-xs text-muted-foreground font-medium tabular-nums">
                      {formatDateTime(log.created_at, locale)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <div className="bg-muted p-1.5 rounded-lg border border-border">
                          <Shield className="w-3.5 h-3.5 text-muted-foreground" />
                        </div>
                        <span className="font-bold text-foreground text-xs">{log.caller_user_id}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant={getRoleBadgeVariant(log.caller_role)}>
                        {log.caller_role}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-mono text-xs text-foreground">
                      {log.method}
                    </TableCell>
                    <TableCell>
                      {log.resource_type ? (
                        <div className="flex flex-col gap-0.5">
                          <span className="text-[9px] font-bold uppercase text-muted-foreground">
                            {log.resource_type}
                          </span>
                          {log.resource_id && (
                            <span className="font-mono text-[10px] text-primary">
                              {log.resource_id}
                            </span>
                          )}
                        </div>
                      ) : (
                        <span className="text-xs text-muted-foreground">—</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <StatusBadge status={isSuccess ? "normal" : "critical"}>
                        {isSuccess ? t("statusSuccess") : t("statusFailure")}
                      </StatusBadge>
                    </TableCell>
                    <TableCell className="text-right">
                      {isExpanded ? (
                        <ChevronUp className="w-4 h-4 text-muted-foreground ml-auto" />
                      ) : (
                        <ChevronDown className="w-4 h-4 text-muted-foreground ml-auto" />
                      )}
                    </TableCell>
                  </TableRow>
                  {isExpanded && (
                    <TableRow className="bg-muted/50 hover:bg-muted/50">
                      <TableCell colSpan={7} className="p-4 border-t border-border">
                        <div className="flex flex-col gap-2 bg-card p-3 rounded-lg border border-border text-xs font-mono text-foreground overflow-x-auto max-w-full">
                          <div className="flex flex-col gap-1">
                            <span className="font-bold text-muted-foreground uppercase text-[9px] tracking-wider">
                              {t("logId")}
                            </span>
                            <span className="text-foreground">{log.id}</span>
                          </div>
                          {log.action && (
                            <div className="flex flex-col gap-1 mt-2">
                              <span className="font-bold text-muted-foreground uppercase text-[9px] tracking-wider">
                                {t("auditAction")}
                              </span>
                              <span className="text-foreground">{log.action}</span>
                            </div>
                          )}
                          {log.payload_diff && Object.keys(log.payload_diff).length > 0 && (
                            <div className="flex flex-col gap-1 mt-2">
                              <span className="font-bold text-muted-foreground uppercase text-[9px] tracking-wider">
                                {t("payloadDiff")}
                              </span>
                              <pre className="whitespace-pre-wrap font-mono text-foreground bg-muted p-2.5 rounded border border-border max-h-56 overflow-y-auto">
                                {JSON.stringify(log.payload_diff, null, 2)}
                              </pre>
                            </div>
                          )}
                          <div className="flex flex-col gap-1 mt-2">
                            <span className="font-bold text-muted-foreground uppercase text-[9px] tracking-wider">
                              {t("details")}
                            </span>
                            <pre className="whitespace-pre-wrap font-mono text-foreground bg-muted p-2.5 rounded border border-border">
                              {log.correlation_id || t("noDetails")}
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
    </div>
  )
}
