import { useTranslation } from "react-i18next"
import { Shield, CheckCircle2, XCircle, ChevronDown, ChevronUp } from "lucide-react"
import * as React from "react"
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from "../../../shared/components/ui/Table"
import { Badge } from "../../../shared/components/ui/Badge"
import { Skeleton } from "../../../shared/components/ui/Skeleton"
import type { AuditLog } from "../types"

interface AuditLogsTableProps {
  isLoading: boolean
  auditLogsList: AuditLog[]
  expandedLogId: string | null
  onToggleRowExpansion: (logId: string) => void
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

export const AuditLogsTable = ({
  isLoading,
  auditLogsList,
  expandedLogId,
  onToggleRowExpansion,
}: AuditLogsTableProps) => {
  const { t } = useTranslation("auditLogs")

  return (
    <div className="overflow-x-auto border border-border rounded-xl w-full bg-white">
      <Table className="min-w-[800px]">
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            <TableHead className="w-[180px]">{t("tableTime")}</TableHead>
            <TableHead>{t("tableUser")}</TableHead>
            <TableHead>{t("tableRole")}</TableHead>
            <TableHead>{t("tableAction")}</TableHead>
            <TableHead>{t("tableStatus")}</TableHead>
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
                {t("noLogsFound")}
              </TableCell>
            </TableRow>
          ) : (
            auditLogsList.map((log: AuditLog) => {
              const isExpanded = expandedLogId === log.id
              const isSuccess = log.access_granted
              return (
                <React.Fragment key={log.id}>
                  <TableRow
                    onClick={() => onToggleRowExpansion(log.id)}
                    className="cursor-pointer hover:bg-gray-50/80 transition-colors"
                  >
                    <TableCell className="text-xs text-gray-600 font-medium">
                      {formatTimestamp(log.created_at)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <div className="bg-gray-100 p-1.5 rounded-lg border border-gray-200">
                          <Shield className="w-3.5 h-3.5 text-gray-500" />
                        </div>
                        <span className="font-bold text-gray-900 text-xs">{log.caller_user_id}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant={getRoleBadgeVariant(log.caller_role)}>
                        {log.caller_role}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-mono text-xs text-gray-700">
                      {log.method}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1.5">
                        {isSuccess ? (
                          <CheckCircle2 className="w-4 h-4 text-green-500 shrink-0" />
                        ) : (
                          <XCircle className="w-4 h-4 text-red-500 shrink-0" />
                        )}
                        <span className={`text-xs font-bold ${isSuccess ? "text-green-600" : "text-red-600"}`}>
                          {isSuccess ? t("statusSuccess") : t("statusFailure")}
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
                              {t("logId")}
                            </span>
                            <span className="text-gray-800">{log.id}</span>
                          </div>
                          <div className="flex flex-col gap-1 mt-2">
                            <span className="font-bold text-gray-500 uppercase text-[9px] tracking-wider">
                              {t("details")}
                            </span>
                            <pre className="whitespace-pre-wrap font-mono text-gray-800 bg-gray-50 p-2.5 rounded border border-border">
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
  )
}
