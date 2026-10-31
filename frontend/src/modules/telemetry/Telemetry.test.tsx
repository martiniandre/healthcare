import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { Telemetry } from "./Telemetry"

const roomsState = vi.hoisted(() => ({ current: {} as Record<string, unknown> }))

vi.mock("./queries", () => ({
  useTelemetryRoomsQuery: () => roomsState.current,
  useTelemetryBedsQuery: () => ({ data: [], isLoading: false, isError: false, refetch: vi.fn(), dataUpdatedAt: 0 }),
  useUnlockRoomMutation: () => ({ isPending: false, mutateAsync: vi.fn() }),
  useUpdateBedConditionMutation: () => ({ isPending: false, mutateAsync: vi.fn() }),
}))

describe("Telemetry rooms screen states", () => {
  beforeEach(() => {
    roomsState.current = {}
  })

  it("should show skeletons while rooms load", () => {
    roomsState.current = { data: undefined, isLoading: true, isError: false, refetch: vi.fn() }

    render(<Telemetry />)

    expect(screen.getByRole("status", { name: "Loading rooms..." })).toBeInTheDocument()
  })

  it("should show an error with retry when rooms fail", () => {
    const refetch = vi.fn()
    roomsState.current = { data: undefined, isLoading: false, isError: true, refetch }

    render(<Telemetry />)
    fireEvent.click(screen.getByRole("button", { name: "Try again" }))

    expect(screen.getByRole("alert")).toBeInTheDocument()
    expect(refetch).toHaveBeenCalledTimes(1)
  })

  it("should show an empty state when there are no rooms", () => {
    roomsState.current = { data: [], isLoading: false, isError: false, refetch: vi.fn() }

    render(<Telemetry />)

    expect(screen.getByText("No monitoring rooms available.")).toBeInTheDocument()
  })
})
