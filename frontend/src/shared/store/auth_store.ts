import { create } from "zustand"
import { authApi } from "../services/auth_api"

interface AuthenticatedUserState {
  isAuthenticated: boolean
  userId: string | null
  role: string | null
  email: string | null
  fullName: string | null
  isActive: boolean | null
  login: (sessionUserId: string, sessionRole: string, sessionEmail: string, sessionFullName?: string, sessionIsActive?: boolean) => void
  logout: () => Promise<void>
}

export const useAuthStore = create<AuthenticatedUserState>((set) => ({
  isAuthenticated: false,
  userId: null,
  role: null,
  email: null,
  fullName: null,
  isActive: null,
  login: (sessionUserId, sessionRole, sessionEmail, sessionFullName, sessionIsActive) =>
    set({
      isAuthenticated: true,
      userId: sessionUserId,
      role: sessionRole,
      email: sessionEmail,
      fullName: sessionFullName ?? null,
      isActive: sessionIsActive ?? null,
    }),
  logout: async () => {
    try {
      await authApi.logout()
    } catch (logoutError) {
      console.error("Logout request failed:", logoutError)
    }
    set({
      isAuthenticated: false,
      userId: null,
      role: null,
      email: null,
      fullName: null,
      isActive: null,
    })
  },
}))

