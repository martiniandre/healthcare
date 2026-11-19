import type { NotificationItem, NotificationPriority } from "./types"

const PRIORITY_RANK: Record<NotificationPriority, number> = {
  critical: 0,
  high: 1,
  medium: 2,
  low: 3,
}

const priorityRankOf = (priority: string): number =>
  priority in PRIORITY_RANK ? PRIORITY_RANK[priority as NotificationPriority] : Object.keys(PRIORITY_RANK).length

export const sortNotificationsByUrgency = (notifications: NotificationItem[]): NotificationItem[] =>
  [...notifications].sort((first, second) => {
    if (first.is_read !== second.is_read) return first.is_read ? 1 : -1
    const priorityDifference = priorityRankOf(first.priority) - priorityRankOf(second.priority)
    if (priorityDifference !== 0) return priorityDifference
    return new Date(second.created_at).getTime() - new Date(first.created_at).getTime()
  })

const readPatientId = (notification: NotificationItem): string | null => {
  const patientId = notification.params?.patient_id
  return typeof patientId === "string" && patientId ? patientId : null
}

export const resolveNotificationTarget = (notification: NotificationItem, isStaffUser: boolean): string | null => {
  const patientId = readPatientId(notification)

  if (!isStaffUser) {
    if (notification.type === "report_ready") return "/portal?tab=reports"
    if (notification.type === "encounter_created" || notification.type === "encounter_updated") return "/portal?tab=encounters"
    return null
  }

  switch (notification.type) {
    case "telemetry_alert":
      return "/telemetry"
    case "exam_complete":
      return "/exam-analyzer"
    case "encounter_created":
    case "encounter_updated":
      if (patientId && notification.resource_id) return `/patients/${patientId}/encounters/${notification.resource_id}`
      return patientId ? `/patients/${patientId}` : null
    case "report_ready":
      return patientId ? `/patients/${patientId}` : null
    case "patient_created":
    case "patient_updated":
      return "/"
    case "audit_alert":
      return "/audit-logs"
    default:
      return null
  }
}
