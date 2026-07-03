import { http } from "../../shared/utils/http"
import type { LoginRequest, LoginResponse, RegisterRequest } from "./types"

export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    return http.post<LoginResponse>("/auth/login", data)
  },

  me: async (): Promise<LoginResponse> => {
    return http.get<LoginResponse>("/auth/me")
  },

  logout: async (): Promise<void> => {
    try {
      await http.post("/auth/logout")
    } catch (logoutError) {
      console.warn("Logout endpoint failed:", logoutError)
    }
  },

  register: async (data: RegisterRequest): Promise<void> => {
    await http.post("/auth/register", data)
  },
}
