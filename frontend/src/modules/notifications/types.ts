export type NotificationPriority = "critical" | "high" | "medium" | "low"

export type NotificationType =
  | "telemetry_alert"
  | "exam_complete"
  | "encounter_created"
  | "encounter_updated"
  | "patient_created"
  | "patient_updated"
  | "audit_alert"
  | "report_ready"
  | "system"

export interface NotificationItem {
  id: string
  type: NotificationType
  priority: NotificationPriority
  title: string
  body: string
  title_key?: string
  body_key?: string
  params?: Record<string, unknown>
  resource_type: string
  resource_id: string
  is_read: boolean
  created_at: string
}

export interface NotificationListResponse {
  notifications: NotificationItem[]
  total: number
}

export interface UnreadCountResponse {
  count: number
}
