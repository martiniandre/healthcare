import { describe, it, expect, vi } from "vitest"
import { render, screen } from "@testing-library/react"
import { NotificationItem } from "./NotificationItem"
import type { NotificationItem as NotificationItemType } from "../types"

const mockNotification: NotificationItemType = {
  id: "notif-1",
  type: "telemetry_alert",
  priority: "critical",
  title: "Alerta Clínico - Leito 01",
  body: "Paciente apresenta condição crítica.",
  resource_type: "bed",
  resource_id: "bed-1",
  is_read: false,
  created_at: new Date().toISOString(),
}

describe("NotificationItem", () => {
  it("renders title and body", () => {
    render(<NotificationItem notification={mockNotification} onMarkRead={vi.fn()} />)
    expect(screen.getByText("Alerta Clínico - Leito 01")).toBeDefined()
    expect(screen.getByText("Paciente apresenta condição crítica.")).toBeDefined()
  })

  it("calls onMarkRead when clicked", () => {
    const onMarkRead = vi.fn()
    render(<NotificationItem notification={mockNotification} onMarkRead={onMarkRead} />)
    screen.getByRole("button").click()
    expect(onMarkRead).toHaveBeenCalledWith("notif-1")
  })
})
