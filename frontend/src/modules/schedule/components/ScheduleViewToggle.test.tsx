import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { ScheduleViewToggle } from "./ScheduleViewToggle"

describe("ScheduleViewToggle", () => {
  it("should render a view selector with the current value", () => {
    render(<ScheduleViewToggle value="week" onChange={vi.fn()} />)
    const viewSelector = screen.getByRole("combobox", { name: "Calendar view" })
    expect(viewSelector).toBeInTheDocument()
    expect(viewSelector).toHaveTextContent("Week")
  })

  it("should render the week, day, month and year options", async () => {
    render(<ScheduleViewToggle value="week" onChange={vi.fn()} />)
    fireEvent.click(screen.getByRole("combobox", { name: "Calendar view" }))
    expect(await screen.findByRole("option", { name: "Week" })).toBeInTheDocument()
    expect(screen.getByRole("option", { name: "Day" })).toBeInTheDocument()
    expect(screen.getByRole("option", { name: "Month" })).toBeInTheDocument()
    expect(screen.getByRole("option", { name: "Year" })).toBeInTheDocument()
  })

  it("should call onChange with the selected view when a new option is chosen", async () => {
    const onChange = vi.fn()
    render(<ScheduleViewToggle value="week" onChange={onChange} />)
    fireEvent.click(screen.getByRole("combobox", { name: "Calendar view" }))
    fireEvent.click(await screen.findByRole("option", { name: "Month" }))
    expect(onChange).toHaveBeenCalledWith("month")
  })
})
