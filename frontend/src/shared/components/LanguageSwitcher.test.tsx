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

const openMenu = () => {
  fireEvent.keyDown(screen.getByRole("button", { name: "English" }), { key: "Enter" })
}

describe("LanguageSwitcher", () => {
  beforeEach(() => {
    mockChangeLanguage.mockReset()
  })

  it("should render the active language trigger for the default layout", () => {
    render(<LanguageSwitcher />)

    const trigger = screen.getByRole("button", { name: "English" })
    expect(trigger).toBeInTheDocument()
  })

  it("should open the menu and list every language option", () => {
    render(<LanguageSwitcher />)

    openMenu()

    expect(screen.getByRole("menuitemradio", { name: /Português/ })).toBeInTheDocument()
    expect(screen.getByRole("menuitemradio", { name: /English/ })).toBeChecked()
    expect(screen.getByRole("menuitemradio", { name: /Español/ })).toBeInTheDocument()
  })

  it("should change the language when a new option is selected", () => {
    render(<LanguageSwitcher />)

    openMenu()
    fireEvent.click(screen.getByRole("menuitemradio", { name: /Português/ }))

    expect(mockChangeLanguage).toHaveBeenCalledWith("pt-BR")
  })

  it("should render the full width trigger for the sidebar layout", () => {
    render(<LanguageSwitcher sidebarLayout />)

    const trigger = screen.getByRole("button", { name: "English" })
    expect(trigger.className).toContain("w-full")
    expect(trigger.className).toContain("text-[13px]")
  })
})
