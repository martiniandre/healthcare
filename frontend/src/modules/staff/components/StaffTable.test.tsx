import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { StaffTable } from "./StaffTable"
import type { StaffMember } from "../types"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({ t: (key: string) => key, i18n: { language: "en-US" } }),
}))

const buildProperties = (overrides: Partial<React.ComponentProps<typeof StaffTable>> = {}) => ({
  isLoading: false,
  isError: false,
  onRetry: vi.fn(),
  filteredStaff: [] as StaffMember[],
  onToggleStatus: vi.fn(),
  onDelete: vi.fn(),
  canManageStatus: true,
  canDelete: true,
  isPendingEmployeeId: null,
  ...overrides,
})

describe("StaffTable states", () => {
  it("shows the empty state when there are no professionals", () => {
    render(<StaffTable {...buildProperties()} />)
    expect(screen.getByText("emptyState")).toBeInTheDocument()
  })

  it("shows an error alert with a retry action", () => {
    const onRetry = vi.fn()
    render(<StaffTable {...buildProperties({ isError: true, onRetry })} />)
    expect(screen.getByRole("alert")).toBeInTheDocument()
    fireEvent.click(screen.getByRole("button", { name: "error.retry" }))
    expect(onRetry).toHaveBeenCalledTimes(1)
  })

  it("keeps the table skeleton visible while loading", () => {
    render(<StaffTable {...buildProperties({ isLoading: true })} />)
    expect(screen.getByRole("table")).toBeInTheDocument()
    expect(screen.queryByText("emptyState")).not.toBeInTheDocument()
  })
})
