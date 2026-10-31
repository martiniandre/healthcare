import { describe, it, expect, vi } from "vitest"
import { renderHook, act } from "@testing-library/react"
import { useDiscardGuard } from "./useDiscardGuard"

describe("useDiscardGuard", () => {
  it("should close immediately when the form is pristine", () => {
    const onClose = vi.fn()
    const { result } = renderHook(() => useDiscardGuard({ isDirty: false, onClose }))

    act(() => result.current.requestClose())

    expect(onClose).toHaveBeenCalledTimes(1)
    expect(result.current.isConfirmationOpen).toBe(false)
  })

  it("should ask for confirmation instead of closing when the form is dirty", () => {
    const onClose = vi.fn()
    const { result } = renderHook(() => useDiscardGuard({ isDirty: true, onClose }))

    act(() => result.current.requestClose())

    expect(onClose).not.toHaveBeenCalled()
    expect(result.current.isConfirmationOpen).toBe(true)
  })

  it("should keep the dialog open when the user chooses to keep editing", () => {
    const onClose = vi.fn()
    const { result } = renderHook(() => useDiscardGuard({ isDirty: true, onClose }))

    act(() => result.current.requestClose())
    act(() => result.current.keepEditing())

    expect(onClose).not.toHaveBeenCalled()
    expect(result.current.isConfirmationOpen).toBe(false)
  })

  it("should close after the user confirms the discard", () => {
    const onClose = vi.fn()
    const { result } = renderHook(() => useDiscardGuard({ isDirty: true, onClose }))

    act(() => result.current.requestClose())
    act(() => result.current.discardChanges())

    expect(onClose).toHaveBeenCalledTimes(1)
    expect(result.current.isConfirmationOpen).toBe(false)
  })
})
