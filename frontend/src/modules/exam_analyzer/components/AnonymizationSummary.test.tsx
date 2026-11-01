import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import { AnonymizationSummary } from "./AnonymizationSummary"
import { AnalysisCard } from "./AnalysisCard"
import { ExamAnalysisStatus, type ExamAnalysis } from "../types"

const buildCompletedAnalysis = (overrides: Partial<ExamAnalysis> = {}): ExamAnalysis => ({
  id: "analysis-1",
  file_name: "labs.pdf",
  status: ExamAnalysisStatus.COMPLETED,
  analysis_response: {
    examType: "Chest X-ray",
    qualityAssessment: { score: 0.9, warnings: [] },
    detectedFindings: [],
    possibleInterpretations: [],
    recommendation: { urgency: "normal", nextSteps: [] },
    limitations: [],
    disclaimer: "Assistive analysis only.",
  },
  consent_given: true,
  anonymized: true,
  created_at: "2026-01-01T00:00:00Z",
  updated_at: "2026-01-01T00:00:00Z",
  ...overrides,
})

describe("AnonymizationSummary", () => {
  it("should show the anonymized badge with the steps and the limitations note", () => {
    render(
      <AnonymizationSummary
        report={{
          applied: true,
          steps: ["metadata_removed", "filename_neutralized"],
          limitations: ["burned_in_text_not_removed", "pdf_metadata_not_removed"],
        }}
      />
    )

    expect(screen.getByText("Anonymized")).toBeInTheDocument()
    expect(screen.getByText("File metadata removed")).toBeInTheDocument()
    expect(screen.getByText("Original file name replaced")).toBeInTheDocument()
    expect(screen.getByText(/Text printed inside the image \(names, ID numbers\) was not removed/)).toBeInTheDocument()
    expect(screen.getByText("PDF metadata was not removed.")).toBeInTheDocument()
  })

  it("should fall back to a generic text for unknown steps and limitations", () => {
    render(<AnonymizationSummary report={{ applied: true, steps: ["new_step"], limitations: ["new_limit"] }} />)

    expect(screen.getByText("Additional anonymization step applied")).toBeInTheDocument()
    expect(screen.getByText("Additional limitation reported by the server")).toBeInTheDocument()
  })

  it("should show the sent without anonymization badge when there is no report", () => {
    render(<AnonymizationSummary report={null} />)

    expect(screen.getByText("Sent without anonymization")).toBeInTheDocument()
    expect(screen.queryByText("Anonymized")).not.toBeInTheDocument()
  })

  it("should treat a report that was not applied as sent without anonymization", () => {
    render(<AnonymizationSummary report={{ applied: false, steps: [], limitations: [] }} />)

    expect(screen.getByText("Sent without anonymization")).toBeInTheDocument()
  })
})

describe("AnalysisCard anonymization summary", () => {
  it("should render the summary when the analysis carries a report", () => {
    render(
      <AnalysisCard
        activeAnalysis={buildCompletedAnalysis({
          anonymizationReport: { applied: true, steps: ["metadata_removed"], limitations: ["burned_in_text_not_removed"] },
        })}
      />
    )

    expect(screen.getByTestId("anonymization-summary")).toBeInTheDocument()
    expect(screen.getByText("Anonymized")).toBeInTheDocument()
  })

  it("should render the not anonymized badge when the report is null", () => {
    render(<AnalysisCard activeAnalysis={buildCompletedAnalysis({ anonymized: false, anonymizationReport: null })} />)

    expect(screen.getByText("Sent without anonymization")).toBeInTheDocument()
  })

  it("should not render the summary for analyses without the report field", () => {
    render(<AnalysisCard activeAnalysis={buildCompletedAnalysis()} />)

    expect(screen.queryByTestId("anonymization-summary")).not.toBeInTheDocument()
  })
})
