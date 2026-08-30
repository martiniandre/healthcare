import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { ScheduleViewToggle } from "./ScheduleViewToggle"

describe("ScheduleViewToggle", () => {
  it("should render a view selector with the current value", () => {
    render(<ScheduleViewToggle value="week" onChange={vi.fn()} />)
    const viewSelector = screen.getByRole("combobox", { name: "Calendar view" })
    expect(viewSelector).toBeInTheDocument()
    expect(viewSelector).toHaveValue("week")
  })

  it("should render the week, day, month and year options", () => {
    render(<ScheduleViewToggle value="week" onChange={vi.fn()} />)
    expect(screen.getByRole("option", { name: "week" })).toBeInTheDocument()
    expect(screen.getByRole("option", { name: "day" })).toBeInTheDocument()
    expect(screen.getByRole("option", { name: "month" })).toBeInTheDocument()
    expect(screen.getByRole("option", { name: "year" })).toBeInTheDocument()
  })

  it("should call onChange with the selected view when a new option is chosen", () => {
    const onChange = vi.fn()
    render(<ScheduleViewToggle value="week" onChange={onChange} />)
    const viewSelector = screen.getByRole("combobox", { name: "Calendar view" })
    fireEvent.change(viewSelector, { target: { value: "month" } })
    expect(onChange).toHaveBeenCalledWith("month")
  })
})
