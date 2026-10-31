import { api, setCsrfToken, clearCsrfToken, setUnauthorizedHandler, refreshSession } from "./api"

const apiMocks = vi.hoisted(() => ({
  post: vi.fn(),
  request: vi.fn(),
}))

const interceptorHandlers = vi.hoisted(() => ({
  requestFulfilled: (config: unknown) => config,
  fulfilled: (response: unknown) => response,
  rejected: (error: unknown) => Promise.reject(error),
}))

vi.mock("axios", () => {
  const instance = {
    create: vi.fn(),
    post: apiMocks.post,
    request: apiMocks.request,
    interceptors: {
      request: {
        use: (fulfilled: (config: unknown) => unknown) => {
          interceptorHandlers.requestFulfilled = fulfilled
        },
      },
      response: {
        use: (fulfilled: (response: unknown) => unknown, rejected: (error: unknown) => unknown) => {
          interceptorHandlers.fulfilled = fulfilled
          interceptorHandlers.rejected = rejected
        },
      },
    },
  }
  instance.create.mockReturnValue(instance)
  return {
    default: { create: instance.create },
  }
})

describe("api client", () => {
  it("should keep the api instance registered after module import", () => {
    expect(api).toBeDefined()
  })

  it("should replace the error message with the backend error payload", async () => {
    const serverError = {
      response: { data: { error: "patient already exists" } },
    }

    await expect(interceptorHandlers.rejected(serverError)).rejects.toMatchObject({
      message: "patient already exists",
    })
  })

  it("should preserve errors without a response body", async () => {
    const networkError = { message: "network down" }

    await expect(interceptorHandlers.rejected(networkError)).rejects.toEqual(networkError)
  })

  it("should pass successful responses through unchanged", () => {
    const successfulResponse = { data: { id: "1" } }

    expect(interceptorHandlers.fulfilled(successfulResponse)).toBe(successfulResponse)
  })

  it("should attach the CSRF token on mutating requests", () => {
    setCsrfToken("csrf-token-123")
    const headers = { set: vi.fn() }
    const config = { method: "post", headers }

    interceptorHandlers.requestFulfilled(config)

    expect(headers.set).toHaveBeenCalledWith("X-CSRF-Token", "csrf-token-123")
    clearCsrfToken()
  })

  it("should not attach the CSRF token on read requests", () => {
    setCsrfToken("csrf-token-123")
    const headers = { set: vi.fn() }
    const config = { method: "get", headers }

    interceptorHandlers.requestFulfilled(config)

    expect(headers.set).not.toHaveBeenCalled()
    clearCsrfToken()
  })

  it("should not attach the CSRF token when none is stored", () => {
    clearCsrfToken()
    const headers = { set: vi.fn() }
    const config = { method: "delete", headers }

    interceptorHandlers.requestFulfilled(config)

    expect(headers.set).not.toHaveBeenCalled()
  })

  it("should fall back to the csrf cookie when no token is stored", () => {
    clearCsrfToken()
    document.cookie = "csrf_token=cookie-token-456; Path=/"
    const headers = { set: vi.fn() }
    const config = { method: "post", headers }

    interceptorHandlers.requestFulfilled(config)

    expect(headers.set).toHaveBeenCalledWith("X-CSRF-Token", "cookie-token-456")
    document.cookie = "csrf_token=; Path=/; Max-Age=0"
  })

  it("should prefer the stored token over the csrf cookie", () => {
    setCsrfToken("stored-token-789")
    document.cookie = "csrf_token=cookie-token-456; Path=/"
    const headers = { set: vi.fn() }
    const config = { method: "put", headers }

    interceptorHandlers.requestFulfilled(config)

    expect(headers.set).toHaveBeenCalledWith("X-CSRF-Token", "stored-token-789")
    document.cookie = "csrf_token=; Path=/; Max-Age=0"
    clearCsrfToken()
  })


  describe("silent refresh", () => {
    const refreshedSession = { userId: "user-1", role: "DOCTOR", email: "a@b.com", csrfToken: "csrf-new" }

    function buildExpiredError(url: string, extra: Record<string, unknown> = {}) {
      return { config: { url, ...extra }, response: { status: 401, data: {} } }
    }

    beforeEach(() => {
      apiMocks.post.mockReset()
      apiMocks.request.mockReset()
    })

    afterEach(() => {
      setUnauthorizedHandler(null)
      clearCsrfToken()
    })

    it("should refresh once and replay the original request", async () => {
      apiMocks.post.mockResolvedValue({ data: refreshedSession })
      apiMocks.request.mockResolvedValue({ data: { id: "patient-1" } })

      const replayed = await interceptorHandlers.rejected(buildExpiredError("/patients"))

      expect(apiMocks.post).toHaveBeenCalledTimes(1)
      expect(apiMocks.post).toHaveBeenCalledWith("/auth/refresh")
      expect(apiMocks.request).toHaveBeenCalledTimes(1)
      expect(apiMocks.request.mock.calls[0][0]).toMatchObject({ url: "/patients", hasBeenReplayed: true })
      expect(replayed).toEqual({ data: { id: "patient-1" } })
    })

    it("should store the csrf token returned by the refresh", async () => {
      apiMocks.post.mockResolvedValue({ data: refreshedSession })
      apiMocks.request.mockResolvedValue({ data: {} })

      await interceptorHandlers.rejected(buildExpiredError("/patients"))
      const headers = { set: vi.fn() }
      interceptorHandlers.requestFulfilled({ method: "post", headers })

      expect(headers.set).toHaveBeenCalledWith("X-CSRF-Token", "csrf-new")
    })

    it("should share one refresh between concurrent 401 responses", async () => {
      let resolveRefresh: (value: unknown) => void = () => {}
      apiMocks.post.mockImplementation(
        () => new Promise((resolve) => {
          resolveRefresh = resolve
        }),
      )
      apiMocks.request.mockResolvedValue({ data: {} })

      const pendingReplays = [
        interceptorHandlers.rejected(buildExpiredError("/patients")),
        interceptorHandlers.rejected(buildExpiredError("/appointments")),
        interceptorHandlers.rejected(buildExpiredError("/staff")),
      ]
      resolveRefresh({ data: refreshedSession })
      await Promise.all(pendingReplays)

      expect(apiMocks.post).toHaveBeenCalledTimes(1)
      expect(apiMocks.request).toHaveBeenCalledTimes(3)
    })

    it("should notify the unauthorized handler once when the refresh fails", async () => {
      const unauthorizedHandler = vi.fn()
      setUnauthorizedHandler(unauthorizedHandler)
      setCsrfToken("stale-token")
      apiMocks.post.mockRejectedValue(buildExpiredError("/auth/refresh"))

      const results = await Promise.allSettled([
        interceptorHandlers.rejected(buildExpiredError("/patients")),
        interceptorHandlers.rejected(buildExpiredError("/appointments")),
      ])

      expect(results.every((result) => result.status === "rejected")).toBe(true)
      expect(unauthorizedHandler).toHaveBeenCalledTimes(1)
      expect(apiMocks.request).not.toHaveBeenCalled()
      const headers = { set: vi.fn() }
      interceptorHandlers.requestFulfilled({ method: "post", headers })
      expect(headers.set).not.toHaveBeenCalledWith("X-CSRF-Token", "stale-token")
    })

    it("should not refresh for auth endpoints", async () => {
      const unauthorizedHandler = vi.fn()
      setUnauthorizedHandler(unauthorizedHandler)

      for (const url of ["/auth/login", "/auth/me", "/auth/logout", "/auth/register", "/auth/refresh"]) {
        await expect(interceptorHandlers.rejected(buildExpiredError(url))).rejects.toBeDefined()
      }

      expect(apiMocks.post).not.toHaveBeenCalled()
      expect(unauthorizedHandler).not.toHaveBeenCalled()
    })

    it("should not loop when the replayed request returns 401 again", async () => {
      const unauthorizedHandler = vi.fn()
      setUnauthorizedHandler(unauthorizedHandler)

      await expect(
        interceptorHandlers.rejected(buildExpiredError("/patients", { hasBeenReplayed: true })),
      ).rejects.toBeDefined()

      expect(apiMocks.post).not.toHaveBeenCalled()
      expect(apiMocks.request).not.toHaveBeenCalled()
      expect(unauthorizedHandler).toHaveBeenCalledTimes(1)
    })

    it("should not refresh for non-401 failures", async () => {
      await expect(
        interceptorHandlers.rejected({ config: { url: "/patients" }, response: { status: 403, data: {} } }),
      ).rejects.toBeDefined()

      expect(apiMocks.post).not.toHaveBeenCalled()
    })

    it("should expose a shared refreshSession that returns the session", async () => {
      apiMocks.post.mockResolvedValue({ data: refreshedSession })

      const [first, second] = await Promise.all([refreshSession(), refreshSession()])

      expect(apiMocks.post).toHaveBeenCalledTimes(1)
      expect(first).toEqual(refreshedSession)
      expect(second).toEqual(refreshedSession)
    })
  })
})
