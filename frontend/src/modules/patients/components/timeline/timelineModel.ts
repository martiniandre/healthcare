import { TimelineEntryType, type TimelineItem } from "../../types"
import type { ClinicalStatusKind } from "../clinicalStatus"

export const TimelinePeriod = {
  Last30Days: "last30Days",
  Last6Months: "last6Months",
  Last12Months: "last12Months",
  All: "all",
  Custom: "custom",
} as const

export type TimelinePeriod = (typeof TimelinePeriod)[keyof typeof TimelinePeriod]

export interface TimelineFilterState {
  types: TimelineEntryType[]
  period: TimelinePeriod
  from: string
  to: string
}

export interface TimelineDayGroup {
  dayKey: string
  items: TimelineItem[]
}

export const TIMELINE_ENTRY_TYPES: readonly TimelineEntryType[] = Object.values(TimelineEntryType)

export const DEFAULT_TIMELINE_PERIOD: TimelinePeriod = TimelinePeriod.Last6Months

export const DEFAULT_TIMELINE_FILTERS: TimelineFilterState = {
  types: [],
  period: DEFAULT_TIMELINE_PERIOD,
  from: "",
  to: "",
}

export const TIMELINE_PAGE_SIZE = 30

const periodValues: readonly string[] = Object.values(TimelinePeriod)
const dayKeyPattern = /^\d{4}-\d{2}-\d{2}$/

export const parseTimelineFilters = (searchParameters: URLSearchParams): TimelineFilterState => {
  const rawPeriod = searchParameters.get("period") ?? ""
  const rawTypes = (searchParameters.get("types") ?? "").split(",")
  const rawFrom = searchParameters.get("from") ?? ""
  const rawTo = searchParameters.get("to") ?? ""
  return {
    types: TIMELINE_ENTRY_TYPES.filter((entryType) => rawTypes.includes(entryType)),
    period: periodValues.includes(rawPeriod) ? (rawPeriod as TimelinePeriod) : DEFAULT_TIMELINE_PERIOD,
    from: dayKeyPattern.test(rawFrom) ? rawFrom : "",
    to: dayKeyPattern.test(rawTo) ? rawTo : "",
  }
}

export const serializeTimelineFilters = (filters: TimelineFilterState): Record<string, string> => {
  const serialized: Record<string, string> = {}
  if (filters.types.length > 0) serialized.types = filters.types.join(",")
  if (filters.period !== DEFAULT_TIMELINE_PERIOD) serialized.period = filters.period
  if (filters.period === TimelinePeriod.Custom) {
    if (filters.from) serialized.from = filters.from
    if (filters.to) serialized.to = filters.to
  }
  return serialized
}

const ALL_HISTORY_START = "1900-01-01T00:00:00.000Z"

export const resolveTimelineRange = (
  filters: TimelineFilterState,
  referenceDate: Date,
): { from?: string; to?: string } => {
  if (filters.period === TimelinePeriod.Last6Months) {
    return {}
  }
  if (filters.period === TimelinePeriod.All) {
    return { from: ALL_HISTORY_START, to: referenceDate.toISOString() }
  }
  if (filters.period === TimelinePeriod.Last30Days) {
    const windowStart = new Date(referenceDate)
    windowStart.setDate(windowStart.getDate() - 30)
    return { from: windowStart.toISOString(), to: referenceDate.toISOString() }
  }
  if (filters.period === TimelinePeriod.Custom) {
    const range: { from?: string; to?: string } = {}
    if (filters.from) range.from = new Date(`${filters.from}T00:00:00`).toISOString()
    if (filters.to) range.to = new Date(`${filters.to}T23:59:59.999`).toISOString()
    return range
  }
  const windowStart = new Date(referenceDate)
  windowStart.setMonth(windowStart.getMonth() - 12)
  return { from: windowStart.toISOString(), to: referenceDate.toISOString() }
}

const padNumber = (value: number): string => String(value).padStart(2, "0")

export const toLocalDayKey = (isoDateTime: string): string => {
  const parsedDate = new Date(isoDateTime)
  return `${parsedDate.getFullYear()}-${padNumber(parsedDate.getMonth() + 1)}-${padNumber(parsedDate.getDate())}`
}

export const groupTimelineByDay = (items: TimelineItem[]): TimelineDayGroup[] => {
  const groupsByDay = new Map<string, TimelineItem[]>()
  items.forEach((timelineItem) => {
    const dayKey = toLocalDayKey(timelineItem.occurredAt)
    const dayItems = groupsByDay.get(dayKey)
    if (dayItems) {
      dayItems.push(timelineItem)
    } else {
      groupsByDay.set(dayKey, [timelineItem])
    }
  })
  return Array.from(groupsByDay, ([dayKey, dayItems]) => ({ dayKey, items: dayItems }))
}

export const statusKindByType: Partial<Record<TimelineEntryType, ClinicalStatusKind>> = {
  [TimelineEntryType.Encounter]: "encounter",
  [TimelineEntryType.Condition]: "condition",
  [TimelineEntryType.Allergy]: "allergy",
  [TimelineEntryType.Medication]: "medication",
  [TimelineEntryType.Report]: "report",
}

export const statusNamespaceByType: Record<TimelineEntryType, string> = {
  [TimelineEntryType.Encounter]: "encounter",
  [TimelineEntryType.Observation]: "observation",
  [TimelineEntryType.Condition]: "condition",
  [TimelineEntryType.Allergy]: "allergy",
  [TimelineEntryType.Medication]: "medicationRequest",
  [TimelineEntryType.Report]: "diagnosticReport",
  [TimelineEntryType.Imaging]: "imagingStudy",
}

export const isSameMoment = (firstDateTime: string, secondDateTime: string): boolean =>
  new Date(firstDateTime).getTime() === new Date(secondDateTime).getTime()
