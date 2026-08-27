import { describe, it, expect, vi } from "vitest"
import { renderHook } from "@testing-library/react"
import { useClinicalConditionsColumns } from "./useClinicalConditionsColumns"

const mockTranslateFunction = (key: string) => key

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: mockTranslateFunction,
  }),
}))

describe("useClinicalConditionsColumns", () => {
  it("should return four column definitions", () => {
    const { result } = renderHook(() => useClinicalConditionsColumns())
    expect(result.current).toHaveLength(4)
  })

  it("should translate every column header", () => {
    const { result } = renderHook(() => useClinicalConditionsColumns())
    const headers = result.current.map((column) => column.header as string)
    expect(headers).toEqual([
      "details.conditionsCard.code",
      "details.conditionsCard.display",
      "details.conditionsCard.status",
      "details.conditionsCard.date",
    ])
  })

  it("should bind each column to its condition accessor", () => {
    const { result } = renderHook(() => useClinicalConditionsColumns())
    const accessorKeys = result.current.map((column) => column.id)
    expect(accessorKeys).toEqual(["icd10_code", "code_display", "clinical_status", "created_at"])
  })

  it("should keep the same array reference across rerenders with stable dependencies", () => {
    const { result, rerender } = renderHook(() => useClinicalConditionsColumns())
    const firstReference = result.current
    rerender()
    expect(result.current).toBe(firstReference)
  })
})
