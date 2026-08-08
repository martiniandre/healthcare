import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { ReportModal } from "./ReportModal"
import { createModuleTranslator } from "../../../../shared/i18n/i18n"

vi.mock("react-i18next", async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useTranslation: () => ({
      t: (key: string) => key,
    }),
  }
})

describe("ReportModal", () => {
  it("should show validation error when exam type is not selected", async () => {
    const onSubmit = vi.fn()
    render(<ReportModal isOpen onClose={vi.fn()} onSubmit={onSubmit} isPending={false} />)
    fireEvent.change(screen.getByPlaceholderText("modals.report.examPlaceholder"), {
      target: { value: "Hemograma Completo" },
    })
    fireEvent.change(screen.getByPlaceholderText("modals.report.conclusionPlaceholder"), {
      target: { value: "Resultados dentro dos limites esperados para o exame." },
    })
    fireEvent.click(screen.getByRole("button", { name: "modals.report.confirm" }))
    expect(await screen.findByText(createModuleTranslator("patients")("validation.reportCodeReq"))).toBeDefined()
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it("should submit report with selected exam type", async () => {
    const onSubmit = vi.fn()
    render(<ReportModal isOpen onClose={vi.fn()} onSubmit={onSubmit} isPending={false} />)
    fireEvent.change(screen.getByPlaceholderText("modals.report.examPlaceholder"), {
      target: { value: "Hemograma Completo" },
    })
    fireEvent.change(screen.getByPlaceholderText("modals.report.conclusionPlaceholder"), {
      target: { value: "Resultados dentro dos limites esperados para o exame." },
    })
    fireEvent.click(screen.getByRole("combobox"))
    const visibleOption = (await screen.findAllByRole("option", { name: "modals.report.examTypes.completeBloodCount" }))
      .find((option) => option.tagName === "DIV")
    if (!visibleOption) {
      throw new Error("Exam type option not found")
    }
    fireEvent.click(visibleOption)
    fireEvent.click(screen.getByRole("button", { name: "modals.report.confirm" }))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())
    expect(onSubmit.mock.calls[0][0]).toEqual({
      reportCode: "58410-2",
      reportDisplay: "Hemograma Completo",
      conclusion: "Resultados dentro dos limites esperados para o exame.",
    })
  })
})
