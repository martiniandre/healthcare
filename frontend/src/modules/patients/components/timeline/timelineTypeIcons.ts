import { Activity, FileText, HeartPulse, Pill, ScanLine, ShieldAlert, Stethoscope, type LucideIcon } from "lucide-react"
import { TimelineEntryType } from "../../types"

export const timelineTypeIcons: Record<TimelineEntryType, LucideIcon> = {
  [TimelineEntryType.Encounter]: Stethoscope,
  [TimelineEntryType.Observation]: HeartPulse,
  [TimelineEntryType.Condition]: Activity,
  [TimelineEntryType.Allergy]: ShieldAlert,
  [TimelineEntryType.Medication]: Pill,
  [TimelineEntryType.Report]: FileText,
  [TimelineEntryType.Imaging]: ScanLine,
}
