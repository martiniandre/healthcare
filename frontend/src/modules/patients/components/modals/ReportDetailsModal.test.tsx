import { describe, it, expect, vi } from "vitest"
import { render, screen } from "@testing-library/react"
import { ReportDetailsModal } from "./ReportDetailsModal"
import type { DiagnosticReport } from "../../types"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

const baseReport: DiagnosticReport = {
  fhir_id: "report-1",
  encounter_fhir_id: "encounter-1",
  patient_fhir_id: "patient-1",
  report_code: "58410-2",
  report_display: "Hemograma Completo",
  status: "final",
  conclusion: "Resultados dentro dos limites esperados.",
  version: "2",
  created_at: "2026-01-01T10:00:00Z",
}

describe("ReportDetailsModal", () => {
  it("should render nothing when closed", () => {
    const { container } = render(
      <ReportDetailsModal isOpen={false} onClose={vi.fn()} report={baseReport} />
    )
    expect(container).toBeEmptyDOMElement()
  })

  it("should render nothing when report is null", () => {
    const { container } = render(
      <ReportDetailsModal isOpen onClose={vi.fn()} report={null} />
    )
    expect(container).toBeEmptyDOMElement()
  })

  it("should render full report details", () => {
    render(<ReportDetailsModal isOpen onClose={vi.fn()} report={baseReport} />)
    expect(screen.getByText("modals.reportDetails.title — Hemograma Completo")).toBeDefined()
    expect(screen.getByText("Hemograma Completo")).toBeDefined()
    expect(screen.getByText("Resultados dentro dos limites esperados.")).toBeDefined()
    expect(screen.getByText("final")).toBeDefined()
    expect(screen.getByText("v2")).toBeDefined()
  })

  it("should render empty conclusion state", () => {
    const reportWithoutConclusion: DiagnosticReport = {
      ...baseReport,
      conclusion: "",
      version: undefined,
    }
    render(
      <ReportDetailsModal isOpen onClose={vi.fn()} report={reportWithoutConclusion} />
    )
    expect(screen.getByText("modals.reportDetails.noConclusion")).toBeDefined()
  })
})
