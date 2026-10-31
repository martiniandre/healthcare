import { describe, it, expect, beforeEach, vi } from "vitest"
import { act, renderHook } from "@testing-library/react"
import { useAlarmMute } from "./useAlarmMute"

describe("useAlarmMute", () => {
  beforeEach(() => {
    vi.restoreAllMocks()
    window.localStorage.clear()
  })

  it("should default to muted and persist changes", () => {
    const { result } = renderHook(() => useAlarmMute())
    expect(result.current.isMuted).toBe(true)

    act(() => result.current.setIsMuted(false))

    expect(result.current.isMuted).toBe(false)
    expect(window.localStorage.getItem("telemetry.alarmsMuted")).toBe("false")
  })

  it("should restore the persisted preference", () => {
    window.localStorage.setItem("telemetry.alarmsMuted", "false")

    const { result } = renderHook(() => useAlarmMute())

    expect(result.current.isMuted).toBe(false)
  })

  it("should keep working when storage throws", () => {
    vi.spyOn(Storage.prototype, "getItem").mockImplementation(() => {
      throw new Error("blocked")
    })
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {
      throw new Error("blocked")
    })

    const { result } = renderHook(() => useAlarmMute())
    act(() => result.current.setIsMuted(false))

    expect(result.current.isMuted).toBe(false)
  })
})
