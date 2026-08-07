export interface AuditLog {
  id: string
  correlation_id: string
  caller_user_id: string
  caller_role: string
  method: string
  access_granted: boolean
  resource_type?: string
  resource_id?: string
  action?: string
  payload_diff?: Record<string, Record<string, unknown>>
  created_at: string
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
