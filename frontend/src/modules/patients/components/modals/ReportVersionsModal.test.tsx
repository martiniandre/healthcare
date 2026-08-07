import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen } from "@testing-library/react"
import { ReportVersionsModal } from "./ReportVersionsModal"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock("../../queries", () => ({
  useDiagnosticReportVersionsQuery: vi.fn(),
}))

import { useDiagnosticReportVersionsQuery } from "../../queries"

const mockedUseDiagnosticReportVersionsQuery = vi.mocked(useDiagnosticReportVersionsQuery)

describe("ReportVersionsModal", () => {
  beforeEach(() => {
    mockedUseDiagnosticReportVersionsQuery.mockReturnValue({
      data: [],
      isLoading: false,
    } as ReturnType<typeof useDiagnosticReportVersionsQuery>)
  })

  it("should render nothing when closed", () => {
    mockedUseDiagnosticReportVersionsQuery.mockReturnValue({
      data: [],
      isLoading: false,
    } as ReturnType<typeof useDiagnosticReportVersionsQuery>)
    const { container } = render(
      <ReportVersionsModal isOpen={false} onClose={vi.fn()} reportFhirId="report-1" reportDisplay="Hemograma Completo" />
    )
    expect(container).toBeEmptyDOMElement()
  })

  it("should render empty state when no versions exist", () => {
    render(
      <ReportVersionsModal isOpen onClose={vi.fn()} reportFhirId="report-1" reportDisplay="Hemograma Completo" />
    )
    expect(
      screen.getByText("modals.reportVersions.title — Hemograma Completo")
    ).toBeDefined()
    expect(screen.getByText("modals.reportVersions.empty")).toBeDefined()
  })

  it("should render version entries with snapshot data", () => {
    mockedUseDiagnosticReportVersionsQuery.mockReturnValue({
      data: [
        {
          version: "1",
          changed_at: "2026-01-01T10:00:00Z",
          snapshot: {
            report_display: "Hemograma Completo",
            conclusion: "Resultados dentro dos limites esperados.",
          },
        },
        {
          version: "2",
          changed_at: "2026-01-02T10:00:00Z",
          snapshot: {
            report_display: "Hemograma Completo",
            conclusion: "Resultados alterados.",
          },
        },
      ],
      isLoading: false,
    } as ReturnType<typeof useDiagnosticReportVersionsQuery>)
    render(
      <ReportVersionsModal isOpen onClose={vi.fn()} reportFhirId="report-1" reportDisplay="Hemograma Completo" />
    )
    expect(screen.getByText("modals.reportVersions.versionLabel 1")).toBeDefined()
    expect(screen.getByText("modals.reportVersions.versionLabel 2")).toBeDefined()
    expect(screen.getByText("Resultados dentro dos limites esperados.")).toBeDefined()
    expect(screen.getByText("Resultados alterados.")).toBeDefined()
  })

  it("should render loading state", () => {
    mockedUseDiagnosticReportVersionsQuery.mockReturnValue({
      data: [],
      isLoading: true,
    } as ReturnType<typeof useDiagnosticReportVersionsQuery>)
    render(
      <ReportVersionsModal isOpen onClose={vi.fn()} reportFhirId="report-1" reportDisplay="Hemograma Completo" />
    )
    expect(document.querySelector(".animate-spin")).not.toBeNull()
  })
})
