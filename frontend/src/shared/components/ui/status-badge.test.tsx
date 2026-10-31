import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import { StatusBadge } from "./status-badge"

describe("StatusBadge", () => {
  it("should render the label together with a decorative icon", () => {
    render(<StatusBadge status="critical">Critical</StatusBadge>)

    const statusBadge = screen.getByText("Critical")
    expect(statusBadge).toHaveAttribute("data-status", "critical")
    expect(statusBadge.querySelector("svg")).toHaveAttribute("aria-hidden", "true")
  })

  it("should use the token classes of each clinical status", () => {
    render(
      <>
        <StatusBadge status="normal">Normal</StatusBadge>
        <StatusBadge status="attention">Attention</StatusBadge>
      </>
    )

    expect(screen.getByText("Normal").className).toContain("text-success")
    expect(screen.getByText("Attention").className).toContain("text-warning")
  })
})
