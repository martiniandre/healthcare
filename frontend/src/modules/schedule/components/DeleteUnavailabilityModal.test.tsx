import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { DeleteUnavailabilityModal } from "./DeleteUnavailabilityModal"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

describe("DeleteUnavailabilityModal", () => {
  it("should render nothing when closed", () => {
    const { container } = render(
      <DeleteUnavailabilityModal
        isOpen={false}
        onClose={vi.fn()}
        onConfirm={vi.fn()}
        isPending={false}
      />
    )
    expect(container.innerHTML).toBe("")
  })

  it("should render the modal title and description when open", () => {
    render(
      <DeleteUnavailabilityModal
        isOpen={true}
        onClose={vi.fn()}
        onConfirm={vi.fn()}
        isPending={false}
      />
    )
    expect(screen.getByText("unavailability.deleteModal.title")).toBeInTheDocument()
    expect(screen.getByText("unavailability.deleteModal.description")).toBeInTheDocument()
  })

  it("should call onClose when the back button is clicked", () => {
    const onClose = vi.fn()
    render(
      <DeleteUnavailabilityModal
        isOpen={true}
        onClose={onClose}
        onConfirm={vi.fn()}
        isPending={false}
      />
    )
    fireEvent.click(screen.getByText("unavailability.deleteModal.back"))
    expect(onClose).toHaveBeenCalled()
  })

  it("should call onConfirm when the confirm button is clicked", () => {
    const onConfirm = vi.fn()
    render(
      <DeleteUnavailabilityModal
        isOpen={true}
        onClose={vi.fn()}
        onConfirm={onConfirm}
        isPending={false}
      />
    )
    fireEvent.click(screen.getByText("unavailability.deleteModal.confirm"))
    expect(onConfirm).toHaveBeenCalled()
  })
})
