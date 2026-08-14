import { api } from "./api"

const interceptorHandlers = vi.hoisted(() => ({
  fulfilled: (response: unknown) => response,
  rejected: (error: unknown) => Promise.reject(error),
}))

vi.mock("axios", () => {
  const instance = {
    create: vi.fn(),
    interceptors: {
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
})
