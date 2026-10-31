import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { TooltipProvider } from "@/shared/components/ui/tooltip"
import { AnalysisHistory } from "./AnalysisHistory"
import { ExamAnalysisStatus, type ExamAnalysis } from "../types"

vi.mock("../../../shared/auth/AbilityContext", () => ({
  Can: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  Action: { Delete: "delete" },
  Feature: { ExamAnalysis: "ExamAnalysis" },
}))

const buildAnalysis = (status: ExamAnalysisStatus): ExamAnalysis => ({
  id: "analysis-1",
  file_name: "chest-xray.png",
  exam_type: "X-ray",
  status,
  analysis_response: { status: ExamAnalysisStatus.INSUFFICIENT_DATA, message: "" },
  consent_given: true,
  anonymized: false,
  created_at: "2026-01-01T10:00:00Z",
  updated_at: "2026-01-01T10:00:00Z",
})

const renderHistory = (properties: Partial<React.ComponentProps<typeof AnalysisHistory>>) =>
  render(
    <TooltipProvider>
      <AnalysisHistory
        history={[]}
        isLoading={false}
        activeID={null}
        onSelect={vi.fn()}
        onDelete={vi.fn()}
        {...properties}
      />
    </TooltipProvider>
  )

describe("AnalysisHistory screen states", () => {
  it("should show loading skeletons", () => {
    renderHistory({ isLoading: true })

    expect(screen.getByRole("status")).toBeInTheDocument()
  })

  it("should show an error with a working retry button", () => {
    const onRetry = vi.fn()
    renderHistory({ isError: true, onRetry })

    fireEvent.click(screen.getByRole("button", { name: "Try again" }))

    expect(screen.getByRole("alert")).toBeInTheDocument()
    expect(onRetry).toHaveBeenCalledTimes(1)
  })

  it("should show the empty message when there are no analyses", () => {
    renderHistory({})

    expect(screen.getByText("No exams analyzed yet")).toBeInTheDocument()
  })

  it("should render the status as a labeled badge", () => {
    renderHistory({ history: [buildAnalysis(ExamAnalysisStatus.FAILED)] })

    expect(document.querySelector('[data-status="critical"]')).not.toBeNull()
  })
})
