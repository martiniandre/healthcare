import { create } from "zustand"

interface AuthenticatedUserState {
  isAuthenticated: boolean
  token: string | null
  userId: string | null
  role: string | null
  email: string | null
  login: (sessionToken: string, sessionUserId: string, sessionRole: string, sessionEmail: string) => void
  logout: () => void
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
  logout: () =>
    set({
      isAuthenticated: false,
      token: null,
      userId: null,
      role: null,
      email: null,
    }),
}))
