import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { NotificationBell } from "./NotificationBell"
import { TooltipProvider } from "../../../shared/components/ui/tooltip"

const refetchList = vi.fn()
const queryState = vi.hoisted(() => ({
  list: { data: undefined as unknown, isLoading: false, isError: false },
}))

vi.mock("../queries", () => ({
  useUnreadCountQuery: () => ({ data: { count: 1200 }, isLoading: false }),
  useNotificationsQuery: () => ({ ...queryState.list, refetch: refetchList }),
  useMarkReadMutation: () => ({ mutate: vi.fn() }),
}))

const renderBell = () =>
  render(
    <TooltipProvider>
      <NotificationBell />
    </TooltipProvider>,
  )

describe("NotificationBell", () => {
  beforeEach(() => {
    refetchList.mockClear()
    queryState.list = { data: undefined, isLoading: false, isError: false }
  })

  it("shows the unread count formatted with Intl and 9+ cap", () => {
    renderBell()
    expect(screen.getByText(`${new Intl.NumberFormat("en-US").format(9)}+`)).toBeDefined()
  })

  it("shows skeleton while loading", () => {
    queryState.list = { data: undefined, isLoading: true, isError: false }
    renderBell()
    fireEvent.click(screen.getAllByRole("button")[0])
    expect(screen.getByTestId("notifications-skeleton")).toBeDefined()
  })

  it("shows error with retry", () => {
    queryState.list = { data: undefined, isLoading: false, isError: true }
    renderBell()
    fireEvent.click(screen.getAllByRole("button")[0])
    fireEvent.click(screen.getByRole("button", { name: "Try again" }))
    expect(refetchList).toHaveBeenCalledTimes(1)
  })

  it("shows empty state", () => {
    queryState.list = { data: { notifications: [], total: 0 }, isLoading: false, isError: false }
    renderBell()
    fireEvent.click(screen.getAllByRole("button")[0])
    expect(screen.getByText("No notifications yet")).toBeDefined()
  })
})
