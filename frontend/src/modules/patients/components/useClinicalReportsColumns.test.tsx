import { describe, it, expect, vi } from "vitest"
import { renderHook } from "@testing-library/react"
import { useClinicalReportsColumns } from "./useClinicalReportsColumns"

const mockTranslateFunction = (key: string) => key

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: mockTranslateFunction,
  }),
}))

describe("useClinicalReportsColumns", () => {
  it("should return six column definitions", () => {
    const { result } = renderHook(() =>
      useClinicalReportsColumns({
        onOpenDetails: vi.fn(),
        onOpenVersions: vi.fn(),
      })
    )
    expect(result.current).toHaveLength(6)
  })

  it("should translate every column header", () => {
    const { result } = renderHook(() =>
      useClinicalReportsColumns({
        onOpenDetails: vi.fn(),
        onOpenVersions: vi.fn(),
      })
    )
    const headers = result.current.map((column) => column.header as string)
    expect(headers).toEqual([
      "details.reportsCard.display",
      "details.reportsCard.conclusion",
      "details.reportsCard.status",
      "details.reportsCard.date",
      "details.reportsCard.version",
      "details.reportsCard.history",
    ])
  })

  it("should bind each column to its report accessor or display id", () => {
    const { result } = renderHook(() =>
      useClinicalReportsColumns({
        onOpenDetails: vi.fn(),
        onOpenVersions: vi.fn(),
      })
    )
    const accessorKeys = result.current.map((column) => column.id)
    expect(accessorKeys).toEqual(["report_display", "conclusion", "status", "created_at", "version", "history"])
  })

  it("should keep the same array reference across rerenders with stable callbacks", () => {
    const onOpenDetails = vi.fn()
    const onOpenVersions = vi.fn()
    const { result, rerender } = renderHook(() =>
      useClinicalReportsColumns({ onOpenDetails, onOpenVersions })
    )
    const firstReference = result.current
    rerender()
    expect(result.current).toBe(firstReference)
  })
})
