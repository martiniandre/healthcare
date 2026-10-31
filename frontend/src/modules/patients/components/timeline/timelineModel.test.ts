import { describe, it, expect } from "vitest"
import type { TimelineItem } from "../../types"
import {
  groupTimelineByDay,
  parseTimelineFilters,
  resolveTimelineRange,
  serializeTimelineFilters,
  DEFAULT_TIMELINE_FILTERS,
  TimelinePeriod,
} from "./timelineModel"

const buildItem = (identifier: string, occurredAt: string): TimelineItem => ({
  id: `encounter:${identifier}`,
  type: "encounter",
  occurredAt,
  clinicalDate: null,
  title: identifier,
  summary: "",
  status: "finished",
  resourceId: identifier,
  encounterId: null,
})

describe("groupTimelineByDay", () => {
  it("should group items sharing the same local day preserving newest-first order", () => {
    const groups = groupTimelineByDay([
      buildItem("c", new Date(2026, 4, 10, 18, 0).toISOString()),
      buildItem("b", new Date(2026, 4, 10, 9, 0).toISOString()),
      buildItem("a", new Date(2026, 4, 8, 12, 0).toISOString()),
    ])

    expect(groups.map((group) => group.dayKey)).toEqual(["2026-05-10", "2026-05-08"])
    expect(groups[0].items.map((groupItem) => groupItem.resourceId)).toEqual(["c", "b"])
    expect(groups[1].items).toHaveLength(1)
  })

  it("should return no groups for an empty list", () => {
    expect(groupTimelineByDay([])).toEqual([])
  })
})

describe("parseTimelineFilters", () => {
  it("should fall back to defaults when the url has no timeline params", () => {
    expect(parseTimelineFilters(new URLSearchParams("tab=timeline"))).toEqual(DEFAULT_TIMELINE_FILTERS)
  })

  it("should read valid types and period and ignore unknown values", () => {
    const filters = parseTimelineFilters(new URLSearchParams("types=report,bogus,allergy&period=last12Months"))

    expect(filters.types).toEqual(["allergy", "report"])
    expect(filters.period).toBe(TimelinePeriod.Last12Months)
  })

  it("should ignore malformed custom dates", () => {
    const filters = parseTimelineFilters(new URLSearchParams("period=custom&from=2026-01-05&to=nope"))

    expect(filters.from).toBe("2026-01-05")
    expect(filters.to).toBe("")
  })
})

describe("serializeTimelineFilters", () => {
  it("should omit default values", () => {
    expect(serializeTimelineFilters(DEFAULT_TIMELINE_FILTERS)).toEqual({})
  })

  it("should only include custom dates for the custom period", () => {
    expect(
      serializeTimelineFilters({ types: ["report"], period: TimelinePeriod.Custom, from: "2026-01-01", to: "2026-02-01" }),
    ).toEqual({ types: "report", period: "custom", from: "2026-01-01", to: "2026-02-01" })
    expect(serializeTimelineFilters({ types: [], period: TimelinePeriod.All, from: "2026-01-01", to: "" })).toEqual({
      period: "all",
    })
  })
})

describe("resolveTimelineRange", () => {
  const referenceDate = new Date("2026-06-15T12:00:00.000Z")

  it("should leave the window to the server for the default period", () => {
    expect(resolveTimelineRange(DEFAULT_TIMELINE_FILTERS, referenceDate)).toEqual({})
  })

  it("should compute a thirty day window", () => {
    const range = resolveTimelineRange({ ...DEFAULT_TIMELINE_FILTERS, period: TimelinePeriod.Last30Days }, referenceDate)

    expect(range.to).toBe("2026-06-15T12:00:00.000Z")
    expect(new Date(range.to!).getTime() - new Date(range.from!).getTime()).toBeGreaterThanOrEqual(29 * 24 * 3600 * 1000)
  })

  it("should compute a twelve month window", () => {
    const range = resolveTimelineRange({ ...DEFAULT_TIMELINE_FILTERS, period: TimelinePeriod.Last12Months }, referenceDate)

    expect(new Date(range.from!).getUTCFullYear()).toBe(2025)
  })

  it("should open the window to the whole history", () => {
    const range = resolveTimelineRange({ ...DEFAULT_TIMELINE_FILTERS, period: TimelinePeriod.All }, referenceDate)

    expect(range.from).toBe("1900-01-01T00:00:00.000Z")
  })

  it("should convert custom local days into an inclusive range", () => {
    const range = resolveTimelineRange(
      { types: [], period: TimelinePeriod.Custom, from: "2026-01-01", to: "2026-01-31" },
      referenceDate,
    )

    expect(range.from).toBe(new Date(2026, 0, 1, 0, 0, 0, 0).toISOString())
    expect(range.to).toBe(new Date(2026, 0, 31, 23, 59, 59, 999).toISOString())
  })
})
