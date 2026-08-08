import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { EncounterModal } from "./EncounterModal"
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

vi.mock("../../../staff/queries", () => ({
  useStaffListQuery: () => ({
    data: [
      { id: "staff-1", fullName: "Dr. Ana Souza", fhirResourceId: "fhir-doctor-1" },
    ],
  }),
}))

describe("EncounterModal", () => {
  it("should show validation error when practitioner is not selected", async () => {
    const onSubmit = vi.fn()
    render(<EncounterModal isOpen onClose={vi.fn()} onSubmit={onSubmit} isPending={false} />)
    fireEvent.change(screen.getByPlaceholderText("modals.encounter.reasonPlaceholder"), {
      target: { value: "Routine checkup" },
    })
    fireEvent.click(screen.getByRole("button", { name: "modals.encounter.confirm" }))
    expect(await screen.findByText(createModuleTranslator("patients")("validation.practitionerReq"))).toBeDefined()
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it("should submit encounter with selected practitioner", async () => {
    const onSubmit = vi.fn()
    render(<EncounterModal isOpen onClose={vi.fn()} onSubmit={onSubmit} isPending={false} />)
    fireEvent.change(screen.getByPlaceholderText("modals.encounter.reasonPlaceholder"), {
      target: { value: "Routine checkup" },
    })
    fireEvent.click(screen.getByRole("combobox"))
    const visibleOption = (await screen.findAllByRole("option", { name: "Dr. Ana Souza" }))
      .find((option) => option.tagName === "DIV")
    if (!visibleOption) {
      throw new Error("Practitioner option not found")
    }
    fireEvent.click(visibleOption)
    fireEvent.click(screen.getByRole("button", { name: "modals.encounter.confirm" }))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())
    expect(onSubmit.mock.calls[0][0]).toEqual({
      reasonDisplay: "Routine checkup",
      practitionerId: "fhir-doctor-1",
    })
  })
})
