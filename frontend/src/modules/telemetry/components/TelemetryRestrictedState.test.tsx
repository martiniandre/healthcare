import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { TelemetryRestrictedState } from "./TelemetryRestrictedState"

const renderState = (overrides: Partial<React.ComponentProps<typeof TelemetryRestrictedState>> = {}) => {
  const handleUnlockRoom = vi.fn((event: { preventDefault: () => void }) => event.preventDefault())
  const setPasscodeInput = vi.fn()
  render(
    <TelemetryRestrictedState
      activeRoomName="ICU"
      passcodeInput=""
      setPasscodeInput={setPasscodeInput}
      passcodeError=""
      isPending={false}
      handleUnlockRoom={handleUnlockRoom}
      {...overrides}
    />,
  )
  return { handleUnlockRoom, setPasscodeInput }
}

describe("TelemetryRestrictedState", () => {
  it("should use a numeric keypad friendly passcode field without autofill", () => {
    renderState()

    const passcodeInput = screen.getByLabelText(/Duty passcode/)
    expect(passcodeInput).toHaveAttribute("inputmode", "numeric")
    expect(passcodeInput).toHaveAttribute("autocomplete", "off")
    expect(passcodeInput).toBeRequired()
  })

  it("should block an empty submit with an inline message", () => {
    const { handleUnlockRoom } = renderState()

    fireEvent.click(screen.getByRole("button", { name: "Release monitoring" }))

    expect(handleUnlockRoom).not.toHaveBeenCalled()
    expect(screen.getByRole("alert")).toHaveTextContent("Enter the shift passcode to unlock the room.")
    expect(screen.getByLabelText(/Duty passcode/)).toHaveAttribute("aria-invalid", "true")
  })

  it("should submit when a passcode is typed", () => {
    const { handleUnlockRoom } = renderState({ passcodeInput: "1234" })

    fireEvent.click(screen.getByRole("button", { name: "Release monitoring" }))

    expect(handleUnlockRoom).toHaveBeenCalledTimes(1)
  })

  it("should reveal the passcode on demand", () => {
    renderState({ passcodeInput: "1234" })

    const toggleButton = screen.getByRole("button", { name: "Show" })
    fireEvent.click(toggleButton)

    expect(toggleButton).toHaveAttribute("aria-pressed", "true")
    expect(screen.getByLabelText(/Duty passcode/)).toHaveAttribute("type", "text")
  })

  it("should link the server error to the field and refocus it", () => {
    renderState({ passcodeInput: "0000", passcodeError: "Incorrect passcode" })

    const passcodeInput = screen.getByLabelText(/Duty passcode/)
    expect(passcodeInput).toHaveAttribute("aria-describedby", expect.stringContaining("telemetry-passcode-error"))
    expect(passcodeInput).toHaveFocus()
  })
})
