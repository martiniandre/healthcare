export interface AuditLog {
  id: string
  timestamp: string
  userId: string
  email: string
  role: string
  action: string
  status: string
  details: string
}

export interface AuditLogsFilter {
  action: string
  email: string
  status: string
  startDate: string
  endDate: string
}

export interface AuditLogsResponse {
  audit_logs: AuditLog[]
  total: number
}
