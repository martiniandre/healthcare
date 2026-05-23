import { http } from "../utils/http"

export const authApi = {
  login: async (emailValue: string, passwordValue: string): Promise<{ token: string; userId: string; role: string; email: string }> => {
    try {
      return await http.post<{ token: string; userId: string; role: string; email: string }>("/api/auth/login", {
        email: emailValue,
        password: passwordValue,
      })
    } catch (apiError) {
      const axiosError = apiError as { response?: { data?: { error?: string } } }
      const errorMessage = axiosError.response?.data?.error || "Credenciais inválidas."
      throw new Error(errorMessage, { cause: apiError })
    }
  },

  logout: async (): Promise<void> => {
    try {
      await http.post("/api/auth/logout")
    } catch (logoutError) {
      console.warn("Logout endpoint failed:", logoutError)
    }
  },
}
