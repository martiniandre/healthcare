import { describe, it, expect, vi } from "vitest"
import { render, screen } from "@testing-library/react"
import { AnalysisCard } from "./AnalysisCard"
import { ExamAnalysisStatus, type ExamAnalysis } from "../types"

vi.mock("./AnalysisResult", () => ({
  AnalysisResult: ({ activeAnalysis }: { activeAnalysis: ExamAnalysis }) => (
    <div data-testid="analysis-result">{activeAnalysis.file_name}</div>
  ),
}))

const buildAnalysis = (status: ExamAnalysisStatus, message = ""): ExamAnalysis => ({
  id: "analysis-1",
  file_name: "chest-xray.png",
  exam_type: "X-ray",
  status,
  analysis_response: { status: ExamAnalysisStatus.INSUFFICIENT_DATA, message },
  consent_given: true,
  anonymized: false,
  created_at: "2026-01-01T10:00:00Z",
  updated_at: "2026-01-01T10:00:00Z",
})

describe("AnalysisCard", () => {
  it("should show the empty state when no analysis is selected", () => {
    render(<AnalysisCard activeAnalysis={null} />)
    expect(screen.queryByRole("status")).not.toBeInTheDocument()
    expect(screen.queryByTestId("analysis-result")).not.toBeInTheDocument()
  })

  it.each([ExamAnalysisStatus.PENDING, ExamAnalysisStatus.PROCESSING])(
    "should announce the processing state for %s analyses",
    (status) => {
      render(<AnalysisCard activeAnalysis={buildAnalysis(status)} />)
      expect(screen.getByRole("status")).toBeInTheDocument()
    }
  )

  it("should alert when the analysis failed", () => {
    render(<AnalysisCard activeAnalysis={buildAnalysis(ExamAnalysisStatus.FAILED)} />)
    expect(screen.getByRole("alert")).toBeInTheDocument()
  })

  it("should show the server message when the data is insufficient", () => {
    render(<AnalysisCard activeAnalysis={buildAnalysis(ExamAnalysisStatus.INSUFFICIENT_DATA, "Imagem desfocada")} />)
    expect(screen.getByText("Imagem desfocada")).toBeInTheDocument()
  })

  it("should render the result for completed analyses", () => {
    render(<AnalysisCard activeAnalysis={buildAnalysis(ExamAnalysisStatus.COMPLETED)} />)
    expect(screen.getByTestId("analysis-result")).toHaveTextContent("chest-xray.png")
  })
})
