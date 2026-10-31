import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios"

interface RetriableRequestConfig extends InternalAxiosRequestConfig {
  hasBeenReplayed?: boolean
}

export interface SessionRefreshResponse {
  userId: string
  role: string
  email: string
  csrfToken: string
}

let csrfTokenFromLogin: string | null = null

export function setCsrfToken(token: string): void {
  csrfTokenFromLogin = token
}

export function clearCsrfToken(): void {
  csrfTokenFromLogin = null
}

type UnauthorizedHandler = () => void

const SESSION_ENDPOINT_PATHS = ["/auth/login", "/auth/me", "/auth/logout", "/auth/register", "/auth/refresh"]

let unauthorizedHandler: UnauthorizedHandler | null = null

export function setUnauthorizedHandler(handler: UnauthorizedHandler | null): void {
  unauthorizedHandler = handler
}

function isSessionEndpoint(requestUrl: string | undefined): boolean {
  if (!requestUrl) {
    return false
  }
  return SESSION_ENDPOINT_PATHS.some((sessionPath) => requestUrl.endsWith(sessionPath))
}

function readCsrfTokenFromCookie(): string | null {
  const cookiePrefix = "csrf_token="
  for (const cookie of document.cookie.split(";")) {
    const trimmedCookie = cookie.trim()
    if (trimmedCookie.startsWith(cookiePrefix)) {
      return trimmedCookie.slice(cookiePrefix.length)
    }
  }
  return null
}

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? "/api/v1",
  withCredentials: true,
})

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const method = config.method?.toUpperCase()
  if (method === "POST" || method === "PUT" || method === "PATCH" || method === "DELETE") {
    const csrfToken = csrfTokenFromLogin ?? readCsrfTokenFromCookie()
    if (csrfToken) {
      config.headers.set("X-CSRF-Token", csrfToken)
    }
  }
  return config
})

let refreshInFlight: Promise<SessionRefreshResponse> | null = null
let guardedRefreshInFlight: Promise<SessionRefreshResponse> | null = null

export function refreshSession(): Promise<SessionRefreshResponse> {
  if (!refreshInFlight) {
    refreshInFlight = api
      .post<SessionRefreshResponse>("/auth/refresh")
      .then((refreshResponse) => {
        if (refreshResponse.data.csrfToken) {
          setCsrfToken(refreshResponse.data.csrfToken)
        }
        return refreshResponse.data
      })
      .finally(() => {
        refreshInFlight = null
      })
  }
  return refreshInFlight
}

function refreshOrEndSession(): Promise<SessionRefreshResponse> {
  if (!guardedRefreshInFlight) {
    guardedRefreshInFlight = refreshSession()
      .catch((refreshError: unknown) => {
        clearCsrfToken()
        unauthorizedHandler?.()
        throw refreshError
      })
      .finally(() => {
        guardedRefreshInFlight = null
      })
  }
  return guardedRefreshInFlight
}

function applyServerMessage(requestError: AxiosError): void {
  if (requestError.response) {
    const errorPayload = requestError.response.data as { error?: string }
    if (errorPayload?.error) {
      requestError.message = errorPayload.error
    }
  }
}

api.interceptors.response.use(
  (successfulResponse) => successfulResponse,
  async (requestError: AxiosError) => {
    const originalRequest = requestError.config as RetriableRequestConfig | undefined
    const isExpiredSession =
      requestError.response?.status === 401 && !isSessionEndpoint(originalRequest?.url)

    if (isExpiredSession && originalRequest) {
      if (originalRequest.hasBeenReplayed) {
        clearCsrfToken()
        unauthorizedHandler?.()
      } else {
        originalRequest.hasBeenReplayed = true
        try {
          await refreshOrEndSession()
        } catch {
          applyServerMessage(requestError)
          return Promise.reject(requestError)
        }
        return api.request(originalRequest)
      }
    } else if (isExpiredSession) {
      clearCsrfToken()
      unauthorizedHandler?.()
    }
    applyServerMessage(requestError)
    return Promise.reject(requestError)
  },
)
