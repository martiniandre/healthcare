import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { AxiosError } from "axios"
import { LoginForm } from "./LoginForm"

const mutateAsync = vi.hoisted(() => vi.fn())
const loginToStore = vi.hoisted(() => vi.fn())

vi.mock("../queries", () => ({
  useLoginMutation: () => ({ isPending: false, mutateAsync }),
}))

vi.mock("../../../shared/store/auth_store", () => ({
  useAuthStore: (selector: (state: { login: typeof loginToStore }) => unknown) =>
    selector({ login: loginToStore }),
}))

const buildHttpError = (status: number) => {
  const httpError = new AxiosError("Request failed")
  httpError.response = { status, data: {}, statusText: "", headers: {}, config: {} as never }
  return httpError
}

const fillAndSubmit = (password = "wrong-password") => {
  fireEvent.change(screen.getByLabelText(/Professional E-mail/), { target: { value: "nurse@hospital.com" } })
  fireEvent.change(screen.getByLabelText(/Access Password/), { target: { value: password } })
  fireEvent.click(screen.getByRole("button", { name: "Enter Console" }))
}

describe("LoginForm", () => {
  beforeEach(() => {
    mutateAsync.mockReset()
    loginToStore.mockReset()
  })

  it("should expose browser friendly input semantics", () => {
    render(<LoginForm />)

    const emailInput = screen.getByLabelText(/Professional E-mail/)
    expect(emailInput).toHaveAttribute("autocomplete", "username")
    expect(emailInput).toHaveAttribute("autocapitalize", "none")
    expect(emailInput).toHaveAttribute("spellcheck", "false")
    expect(screen.getByLabelText(/Access Password/)).toHaveAttribute("autocomplete", "current-password")
  })

  it("should toggle password visibility with aria-pressed", () => {
    render(<LoginForm />)

    const toggleButton = screen.getByRole("button", { name: "Show password" })
    expect(toggleButton).toHaveAttribute("aria-pressed", "false")
    fireEvent.click(toggleButton)

    expect(toggleButton).toHaveAttribute("aria-pressed", "true")
    expect(screen.getByLabelText(/Access Password/)).toHaveAttribute("type", "text")
  })

  it("should show validation errors only after the field is touched", async () => {
    render(<LoginForm />)

    fireEvent.blur(screen.getByLabelText(/Professional E-mail/))

    expect(await screen.findByText("Email is required")).toBeInTheDocument()
    expect(screen.queryByText("Password is required")).not.toBeInTheDocument()
  })

  it("should show a generic credential error, keep the email and focus the password", async () => {
    mutateAsync.mockRejectedValue(buildHttpError(401))
    render(<LoginForm />)

    fillAndSubmit()

    const alert = await screen.findByRole("alert")
    expect(alert).toHaveTextContent("Invalid credentials. Check your e-mail and password and try again.")
    expect(screen.getByLabelText(/Professional E-mail/)).toHaveValue("nurse@hospital.com")
    expect(screen.getByLabelText(/Access Password/)).toHaveValue("")
    await waitFor(() => expect(screen.getByLabelText(/Access Password/)).toHaveFocus())
  })

  it("should show the connection error when the server cannot be reached", async () => {
    mutateAsync.mockRejectedValue(new AxiosError("Network Error"))
    render(<LoginForm />)

    fillAndSubmit()

    expect(await screen.findByText("Error establishing connection with the server.")).toBeInTheDocument()
  })

  it("should warn when Caps Lock is on while typing the password", () => {
    render(<LoginForm />)
    const passwordInput = screen.getByLabelText(/Access Password/)

    fireEvent.keyDown(passwordInput, { key: "A", modifierCapsLock: true })

    expect(screen.getByText("Caps Lock is on.")).toBeInTheDocument()
  })
})
