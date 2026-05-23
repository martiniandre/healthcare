import { create } from "zustand"
import { authApi } from "../services/auth_api"

interface AuthenticatedUserState {
  isAuthenticated: boolean
  token: string | null
  userId: string | null
  role: string | null
  email: string | null
  login: (sessionToken: string, sessionUserId: string, sessionRole: string, sessionEmail: string) => void
  logout: () => Promise<void>
}

export const useAuthStore = create<AuthenticatedUserState>((set) => ({
  isAuthenticated: false,
  token: null,
  userId: null,
  role: null,
  email: null,
  login: (sessionToken, sessionUserId, sessionRole, sessionEmail) =>
    set({
      isAuthenticated: true,
      token: sessionToken,
      userId: sessionUserId,
      role: sessionRole,
      email: sessionEmail,
    }),
  logout: async () => {
    try {
      await authApi.logout()
    } catch (logoutError) {
      console.error("Logout request failed:", logoutError)
    }
    set({
      isAuthenticated: false,
      token: null,
      userId: null,
      role: null,
      email: null,
    })
  },
}))
