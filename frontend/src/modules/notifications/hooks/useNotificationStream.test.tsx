import { renderHook } from "@testing-library/react"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { useNotificationStream } from "./useNotificationStream"
import { useAuthStore } from "../../../shared/store/auth_store"

class MockEventSource {
  static instances: MockEventSource[] = []
  listeners: Record<string, Array<() => void>> = {}
  closed = false
  onerror: (() => void) | null = null

  constructor(public url: string) {
    MockEventSource.instances.push(this)
  }

  addEventListener(type: string, listener: () => void) {
    this.listeners[type] = this.listeners[type] ?? []
    this.listeners[type].push(listener)
  }

  emit(type: string) {
    this.listeners[type]?.forEach((listener) => listener())
  }

  close() {
    this.closed = true
  }
}

describe("useNotificationStream", () => {
  let queryClient: QueryClient

  beforeEach(() => {
    MockEventSource.instances = []
    useAuthStore.setState({ isAuthenticated: false })
    vi.stubGlobal("EventSource", MockEventSource)
    queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    })
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  function renderStream() {
    return renderHook(() => useNotificationStream(), {
      wrapper: ({ children }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
      ),
    })
  }

  it("should open the stream when the user is authenticated", () => {
    useAuthStore.setState({ isAuthenticated: true })

    renderStream()

    expect(MockEventSource.instances).toHaveLength(1)
    expect(MockEventSource.instances[0].url).toBe("/api/v1/notifications/stream")
  })

  it("should not open the stream for unauthenticated users", () => {
    useAuthStore.setState({ isAuthenticated: false })

    renderStream()

    expect(MockEventSource.instances).toHaveLength(0)
  })

  it("should invalidate notification queries when a notification event arrives", () => {
    useAuthStore.setState({ isAuthenticated: true })
    const invalidateSpy = vi.spyOn(queryClient, "invalidateQueries")

    renderStream()
    MockEventSource.instances[0].emit("notification")

    expect(invalidateSpy).toHaveBeenCalledWith({ queryKey: expect.anything() })
  })

  it("should close the stream on unmount", () => {
    useAuthStore.setState({ isAuthenticated: true })

    const { unmount } = renderStream()
    const eventSource = MockEventSource.instances[0]

    unmount()

    expect(eventSource.closed).toBe(true)
  })

  it("should close the stream when the server reports an error", () => {
    useAuthStore.setState({ isAuthenticated: true })

    renderStream()
    const eventSource = MockEventSource.instances[0]

    eventSource.onerror?.()

    expect(eventSource.closed).toBe(true)
  })
})
