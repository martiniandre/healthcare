import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { UnavailabilityCard } from "./UnavailabilityCard"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

const mockUnavailability = {
  id: "unavail-1",
  staff_id: "staff-1",
  starts_at: "2026-09-01T09:00:00Z",
  ends_at: "2026-09-01T12:00:00Z",
  reason: "Vacation",
  created_at: "2026-08-15T10:00:00Z",
}

describe("UnavailabilityCard", () => {
  it("should render the unavailability badge", () => {
    render(<UnavailabilityCard unavailability={mockUnavailability} onDelete={vi.fn()} />)
    expect(screen.getByText("unavailability.badge")).toBeInTheDocument()
  })

  it("should render the time range as separate time fragments", () => {
    render(<UnavailabilityCard unavailability={mockUnavailability} onDelete={vi.fn()} />)
    const timeSpan = screen.getByText(/—/)
    expect(timeSpan).toBeInTheDocument()
  })

  it("should render the reason", () => {
    render(<UnavailabilityCard unavailability={mockUnavailability} onDelete={vi.fn()} />)
    expect(screen.getByText("Vacation")).toBeInTheDocument()
  })

  it("should call onDelete with the unavailability id when delete button is clicked", () => {
    const onDelete = vi.fn()
    render(<UnavailabilityCard unavailability={mockUnavailability} onDelete={onDelete} />)
    fireEvent.click(screen.getByText("unavailability.delete"))
    expect(onDelete).toHaveBeenCalledWith("unavail-1")
  })
})
