import { useAuthStore } from "./auth_store"

vi.mock("../../modules/auth/api", () => ({
  authApi: { logout: vi.fn().mockResolvedValue(undefined) },
}))

describe("auth store", () => {
  it("should reset the whole session when clearSession is called", () => {
    useAuthStore.getState().login("user-1", "DOCTOR", "doctor@example.com", "Doctor", true)

    useAuthStore.getState().clearSession()

    expect(useAuthStore.getState()).toMatchObject({
      isAuthenticated: false,
      userId: null,
      role: null,
      email: null,
    })
  })
})
