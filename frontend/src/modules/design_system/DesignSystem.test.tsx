import { describe, it, expect } from "vitest"
import { fireEvent, render, screen, within } from "@testing-library/react"
import { DesignSystem } from "./DesignSystem"
import { COLOR_TOKENS, SECTION_IDS } from "./catalog"

describe("DesignSystem page", () => {
  it("should render the page title and the anchor navigation", () => {
    render(<DesignSystem />)

    expect(screen.getByRole("heading", { level: 1, name: "Design System" })).toBeInTheDocument()
    const navigation = screen.getByRole("navigation", { name: "Design system sections" })
    expect(within(navigation).getAllByRole("link")).toHaveLength(SECTION_IDS.length)
    expect(within(navigation).getByRole("link", { name: "Colors" })).toHaveAttribute("href", "#colors")
  })

  it("should render every section and a swatch per color token", () => {
    render(<DesignSystem />)

    for (const sectionId of SECTION_IDS) {
      expect(document.getElementById(sectionId)).toBeInTheDocument()
    }
    for (const token of COLOR_TOKENS) {
      expect(screen.getByTestId(`swatch-${token.name}`)).toBeInTheDocument()
    }
  })

  it("should show all five status badges", () => {
    render(<DesignSystem />)

    for (const label of ["Normal", "Attention", "Critical", "Info", "Neutral"]) {
      expect(document.querySelector(`[data-status="${label.toLowerCase()}"]`)).toBeInTheDocument()
    }
  })

  it("should toggle table density", () => {
    render(<DesignSystem />)

    const compactButton = screen.getByRole("button", { name: "Compact" })
    fireEvent.click(compactButton)

    expect(compactButton).toHaveAttribute("aria-pressed", "true")
    expect(document.querySelector('[data-slot="table"]')).toHaveAttribute("data-density", "compact")
  })
})
