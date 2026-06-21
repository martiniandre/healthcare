import { http } from "../utils/http"
import type { AuthResponseDto } from "../types"

export const authApi = {
  login: async (emailValue: string, passwordValue: string): Promise<AuthResponseDto> => {
    try {
      return await http.post<AuthResponseDto>("/auth/login", {
        email: emailValue,
        password: passwordValue,
      })
    } catch (apiError) {
      const axiosError = apiError as { response?: { data?: { error?: string } } }
      const errorMessage = axiosError.response?.data?.error || "Credenciais inválidas."
      throw new Error(errorMessage, { cause: apiError })
    }
  },

  me: async (): Promise<AuthResponseDto> => {
    return await http.get<AuthResponseDto>("/auth/me")
  },

  logout: async (): Promise<void> => {
    try {
      await http.post("/auth/logout")
    } catch (logoutError) {
      console.warn("Logout endpoint failed:", logoutError)
    }
  },
}
