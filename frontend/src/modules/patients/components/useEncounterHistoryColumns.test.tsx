import { describe, it, expect, vi } from "vitest"
import { renderHook } from "@testing-library/react"
import { useEncounterHistoryColumns } from "./useEncounterHistoryColumns"

const mockTranslateFunction = (key: string) => key

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: mockTranslateFunction,
  }),
}))

describe("useEncounterHistoryColumns", () => {
  const buildDependencies = () => ({
    onFinishEncounter: vi.fn(),
    onCancelEncounter: vi.fn(),
    onOpenEncounter: vi.fn(),
    isActionPending: false,
  })

  it("should return four column definitions", () => {
    const { result } = renderHook(() => useEncounterHistoryColumns(buildDependencies()))
    expect(result.current).toHaveLength(4)
  })

  it("should translate every column header", () => {
    const { result } = renderHook(() => useEncounterHistoryColumns(buildDependencies()))
    const headers = result.current.map((column) => column.header as string)
    expect(headers).toEqual([
      "details.encountersCard.reason",
      "details.encountersCard.status",
      "details.encountersCard.date",
      "details.encountersCard.action",
    ])
  })

  it("should bind each column to its encounter accessor or display id", () => {
    const { result } = renderHook(() => useEncounterHistoryColumns(buildDependencies()))
    const accessorKeys = result.current.map((column) => column.id)
    expect(accessorKeys).toEqual(["reason_display", "status", "created_at", "actions"])
  })

  it("should keep the same array reference across rerenders with stable dependencies", () => {
    const dependencies = buildDependencies()
    const { result, rerender } = renderHook(() => useEncounterHistoryColumns(dependencies))
    const firstReference = result.current
    rerender()
    expect(result.current).toBe(firstReference)
  })
})
