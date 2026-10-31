import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { PortalSectionState } from "./PortalSectionState"

const renderState = (overrides: Partial<React.ComponentProps<typeof PortalSectionState>> = {}) => {
  const onRetry = vi.fn()
  render(
    <PortalSectionState
      isLoading={false}
      isError={false}
      onRetry={onRetry}
      isEmpty={false}
      emptyTitle="Nothing here"
      {...overrides}
    >
      <p>content</p>
    </PortalSectionState>,
  )
  return onRetry
}

describe("PortalSectionState", () => {
  it("renders skeleton while loading", () => {
    renderState({ isLoading: true })
    expect(screen.getByTestId("portal-section-skeleton")).toBeDefined()
    expect(screen.queryByText("content")).toBeNull()
  })

  it("renders error alert and retries", () => {
    const onRetry = renderState({ isError: true })
    expect(screen.getByRole("alert")).toBeDefined()
    fireEvent.click(screen.getByRole("button"))
    expect(onRetry).toHaveBeenCalledTimes(1)
  })

  it("renders empty state", () => {
    renderState({ isEmpty: true })
    expect(screen.getByText("Nothing here")).toBeDefined()
  })

  it("renders children when data is available", () => {
    renderState()
    expect(screen.getByText("content")).toBeDefined()
  })
})
