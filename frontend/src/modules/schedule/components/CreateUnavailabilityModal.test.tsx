import { describe, it, expect, vi } from "vitest"
import { render, screen } from "@testing-library/react"
import { CreateUnavailabilityModal } from "./CreateUnavailabilityModal"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock("react-hook-form", () => ({
  useForm: () => ({
    register: vi.fn(() => ({ name: "mock", ref: vi.fn() })),
    handleSubmit: vi.fn((fn) => fn),
    control: {},
    reset: vi.fn(),
    formState: { errors: {} },
  }),
}))

describe("CreateUnavailabilityModal", () => {
  it("should render nothing when closed", () => {
    const { container } = render(
      <CreateUnavailabilityModal
        isOpen={false}
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        staffId="staff-1"
        defaultDate="2026-09-01"
      />
    )
    expect(container.innerHTML).toBe("")
  })

  it("should render the modal title when open", () => {
    render(
      <CreateUnavailabilityModal
        isOpen={true}
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        staffId="staff-1"
        defaultDate="2026-09-01"
      />
    )
    expect(screen.getByText("unavailability.modal.title")).toBeInTheDocument()
  })

  it("should render date, time, and reason inputs", () => {
    render(
      <CreateUnavailabilityModal
        isOpen={true}
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        staffId="staff-1"
        defaultDate="2026-09-01"
      />
    )
    expect(screen.getByText("unavailability.modal.date")).toBeInTheDocument()
    expect(screen.getByText("unavailability.modal.startTime")).toBeInTheDocument()
    expect(screen.getByText("unavailability.modal.endTime")).toBeInTheDocument()
    expect(screen.getByText("unavailability.modal.reason")).toBeInTheDocument()
  })

  it("should render cancel and confirm buttons", () => {
    render(
      <CreateUnavailabilityModal
        isOpen={true}
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        staffId="staff-1"
        defaultDate="2026-09-01"
      />
    )
    expect(screen.getByText("unavailability.modal.cancel")).toBeInTheDocument()
    expect(screen.getByText("unavailability.modal.confirm")).toBeInTheDocument()
  })
})
