import { api, setCsrfToken, clearCsrfToken } from "./api"

const interceptorHandlers = vi.hoisted(() => ({
  requestFulfilled: (config: unknown) => config,
  fulfilled: (response: unknown) => response,
  rejected: (error: unknown) => Promise.reject(error),
}))

vi.mock("axios", () => {
  const instance = {
    create: vi.fn(),
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
})
