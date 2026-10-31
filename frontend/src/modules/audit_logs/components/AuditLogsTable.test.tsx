import { describe, it, expect, vi } from "vitest"
import { render, screen } from "@testing-library/react"
import { AuditLogsTable } from "./AuditLogsTable"
import type { AuditLog } from "../types"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({ t: (key: string) => key, i18n: { language: "en-US" } }),
}))

const buildLog = (overrides: Partial<AuditLog>): AuditLog => ({
  id: "log-1",
  caller_user_id: "user-1",
  caller_role: "ADMIN",
  method: "GET /patients",
  access_granted: true,
  created_at: "2026-05-20T09:00:00Z",
  ...overrides,
} as AuditLog)

const baseProperties = {
  isLoading: false,
  expandedLogId: null,
  onToggleRowExpansion: vi.fn(),
}

describe("AuditLogsTable states", () => {
  it("shows the empty state when no logs match", () => {
    render(<AuditLogsTable {...baseProperties} auditLogsList={[]} />)
    expect(screen.getByText("noLogsFound")).toBeInTheDocument()
  })

  it("renders nothing while the query is in error", () => {
    render(<AuditLogsTable {...baseProperties} isError auditLogsList={[]} />)
    expect(screen.queryByText("noLogsFound")).not.toBeInTheDocument()
    expect(screen.queryByRole("table")).not.toBeInTheDocument()
  })

  it("marks the result with a status badge for success and failure", () => {
    render(
      <AuditLogsTable
        {...baseProperties}
        auditLogsList={[buildLog({}), buildLog({ id: "log-2", access_granted: false })]}
      />
    )
    expect(screen.getByText("statusSuccess").closest("[data-status]")).toHaveAttribute("data-status", "normal")
    expect(screen.getByText("statusFailure").closest("[data-status]")).toHaveAttribute("data-status", "critical")
  })
})
