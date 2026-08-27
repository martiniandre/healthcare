import { describe, it, expect, vi } from "vitest"
import { renderHook } from "@testing-library/react"
import { useClinicalAllergiesColumns } from "./useClinicalAllergiesColumns"

const mockTranslateFunction = (key: string) => key

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: mockTranslateFunction,
  }),
}))

describe("useClinicalAllergiesColumns", () => {
  it("should return five column definitions", () => {
    const { result } = renderHook(() => useClinicalAllergiesColumns())
    expect(result.current).toHaveLength(5)
  })

  it("should translate every column header", () => {
    const { result } = renderHook(() => useClinicalAllergiesColumns())
    const headers = result.current.map((column) => column.header as string)
    expect(headers).toEqual([
      "details.allergiesCard.code",
      "details.allergiesCard.allergen",
      "details.allergiesCard.reaction",
      "details.allergiesCard.status",
      "details.allergiesCard.date",
    ])
  })

  it("should bind each column to its allergy accessor", () => {
    const { result } = renderHook(() => useClinicalAllergiesColumns())
    const accessorKeys = result.current.map((column) => column.id)
    expect(accessorKeys).toEqual(["allergen_code", "allergen_display", "reaction", "clinical_status", "created_at"])
  })

  it("should keep the same array reference across rerenders with stable dependencies", () => {
    const { result, rerender } = renderHook(() => useClinicalAllergiesColumns())
    const firstReference = result.current
    rerender()
    expect(result.current).toBe(firstReference)
  })
})
