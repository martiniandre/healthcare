import axios, { type AxiosError } from "axios"

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? "/api/v1",
  withCredentials: true,
  xsrfCookieName: "csrf_token",
  xsrfHeaderName: "X-CSRF-Token",
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

