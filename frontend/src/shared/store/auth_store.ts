import { create } from "zustand"
import { authApi } from "../../modules/auth/api"
import { clearCsrfToken } from "../services/api"

interface AuthenticatedUserState {
  isAuthenticated: boolean
  userId: string | null
  role: string | null
  email: string | null
  fullName: string | null
  isActive: boolean | null
  login: (sessionUserId: string, sessionRole: string, sessionEmail: string, sessionFullName?: string, sessionIsActive?: boolean) => void
  clearSession: () => void
  logout: () => Promise<void>
}

const anonymousSessionState = {
  isAuthenticated: false,
  userId: null,
  role: null,
  email: null,
  fullName: null,
  isActive: null,
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
  clearSession: () => {
    clearCsrfToken()
    set(anonymousSessionState)
  },
  logout: async () => {
    try {
      await authApi.logout()
    } catch (logoutError) {
      console.error("Logout request failed:", logoutError)
    }
    clearCsrfToken()
    set(anonymousSessionState)
  },
}))

