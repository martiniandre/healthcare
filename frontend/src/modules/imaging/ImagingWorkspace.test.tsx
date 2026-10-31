import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { ImagingWorkspace } from "./ImagingWorkspace"

const studyQueryState = vi.hoisted(() => ({
  current: {} as Record<string, unknown>,
}))

vi.mock("./queries", () => ({
  useImagingStudyQuery: () => studyQueryState.current,
  useUploadImagingStudyMutation: () => ({ isPending: false, mutateAsync: vi.fn() }),
}))

vi.mock("./hooks/useDicomViewer", () => ({
  useDicomViewer: () => ({
    activeTool: "zoom",
    canvasReference: { current: null },
    handleMouseDown: vi.fn(),
    handleMouseMove: vi.fn(),
    handleMouseUp: vi.fn(),
    setActiveTool: vi.fn(),
    applyPreset: vi.fn(),
  }),
}))

describe("ImagingWorkspace screen states", () => {
  beforeEach(() => {
    studyQueryState.current = {}
  })

  it("should show a loading indicator while the study loads", () => {
    studyQueryState.current = { data: undefined, isLoading: true, isError: false, refetch: vi.fn() }

    render(<ImagingWorkspace studyId="study-1" onBack={vi.fn()} />)

    expect(screen.getAllByRole("status").length).toBeGreaterThan(0)
  })

  it("should show an error with a retry action that refetches", () => {
    const refetch = vi.fn()
    studyQueryState.current = { data: undefined, isLoading: false, isError: true, refetch }

    render(<ImagingWorkspace studyId="study-1" onBack={vi.fn()} />)
    fireEvent.click(screen.getByRole("button", { name: "Try again" }))

    expect(screen.getByRole("alert")).toBeInTheDocument()
    expect(refetch).toHaveBeenCalledTimes(1)
  })

  it("should show an empty state when the study does not exist", () => {
    studyQueryState.current = { data: undefined, isLoading: false, isError: false, refetch: vi.fn() }

    render(<ImagingWorkspace studyId="study-1" onBack={vi.fn()} />)

    expect(screen.getByText("Study not found")).toBeInTheDocument()
  })
})
