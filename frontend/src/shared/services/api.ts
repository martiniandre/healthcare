import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios"

let csrfTokenFromLogin: string | null = null

export function setCsrfToken(token: string): void {
  csrfTokenFromLogin = token
}

export function clearCsrfToken(): void {
  csrfTokenFromLogin = null
}

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? "/api/v1",
  withCredentials: true,
})

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const method = config.method?.toUpperCase()
  if (method === "POST" || method === "PUT" || method === "PATCH" || method === "DELETE") {
    if (csrfTokenFromLogin) {
      config.headers.set("X-CSRF-Token", csrfTokenFromLogin)
    }
  }
  return config
})

api.interceptors.response.use(
  (successfulResponse) => successfulResponse,
  (requestError: AxiosError) => {
    if (requestError.response) {
      const errorPayload = requestError.response.data as { error?: string }
      if (errorPayload?.error) {
        requestError.message = errorPayload.error
      }
    }
    return Promise.reject(requestError)
  },
)
