import { describe, it, expect, beforeEach, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { LanguageSwitcher } from "./LanguageSwitcher"

const mockChangeLanguage = vi.fn()

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    i18n: {
      language: "en-US",
      changeLanguage: mockChangeLanguage,
    },
  }),
}))

describe("LanguageSwitcher", () => {
  beforeEach(() => {
    mockChangeLanguage.mockReset()
  })

  it("should render the active language trigger for the default layout", () => {
    render(<LanguageSwitcher />)

    const trigger = screen.getByRole("button", { name: "English" })
    expect(trigger).toBeInTheDocument()
  })

  it("should open the dropdown and list every language option", () => {
    render(<LanguageSwitcher />)

    fireEvent.click(screen.getByRole("button", { name: "English" }))

    expect(screen.getByRole("button", { name: /Português/ })).toBeInTheDocument()
    expect(screen.getAllByRole("button", { name: /English/ })).toHaveLength(2)
    expect(screen.getByRole("button", { name: /Español/ })).toBeInTheDocument()
  })

  it("should change the language and close the dropdown when a new option is selected", () => {
    render(<LanguageSwitcher />)

    fireEvent.click(screen.getByRole("button", { name: "English" }))
    fireEvent.click(screen.getByRole("button", { name: /Português/ }))

    expect(mockChangeLanguage).toHaveBeenCalledWith("pt-BR")
    expect(screen.queryByRole("button", { name: /Español/ })).not.toBeInTheDocument()
  })

  it("should render the full width trigger for the sidebar layout", () => {
    render(<LanguageSwitcher sidebarLayout />)

    const trigger = screen.getByRole("button", { name: "English" })
    expect(trigger.className).toContain("w-full")
    expect(trigger.className).toContain("text-[13px]")
  })
})