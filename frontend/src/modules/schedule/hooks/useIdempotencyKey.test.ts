import { describe, it, expect } from "vitest"
import { renderHook, act } from "@testing-library/react"
import { useIdempotencyKey } from "./useIdempotencyKey"

describe("useIdempotencyKey", () => {
  it("should return the same key on repeated calls", () => {
    const { result } = renderHook(() => useIdempotencyKey())
    const firstKey = result.current.getOrCreateKey()
    const secondKey = result.current.getOrCreateKey()
    expect(secondKey).toBe(firstKey)
  })

  it("should generate a new key after reset", () => {
    const { result } = renderHook(() => useIdempotencyKey())
    const firstKey = result.current.getOrCreateKey()
    act(() => {
      result.current.resetKey()
    })
    const secondKey = result.current.getOrCreateKey()
    expect(secondKey).not.toBe(firstKey)
  })

  it("should generate unique keys across different hook instances", () => {
    const { result: firstHook } = renderHook(() => useIdempotencyKey())
    const { result: secondHook } = renderHook(() => useIdempotencyKey())
    expect(firstHook.current.getOrCreateKey()).not.toBe(secondHook.current.getOrCreateKey())
  })
})
