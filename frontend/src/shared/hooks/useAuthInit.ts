import { useEffect, useState } from "react"
import { authApi } from "../../modules/auth/api"
import { useAuthStore } from "../store/auth_store"
import { refreshSession, setCsrfToken } from "../services/api"

export const useAuthInit = () => {
  const [isLoading, setIsLoading] = useState(true)
  const login = useAuthStore((state) => state.login)

  useEffect(() => {
    let cancelled = false

    authApi
      .me()
      .catch(async (meError: unknown) => {
        const status = (meError as { response?: { status?: number } })?.response?.status
        if (status !== 401) {
          throw meError
        }
        await refreshSession()
        return authApi.me()
      })
      .then((sessionData) => {
        if (!cancelled) {
          if (sessionData.csrfToken) {
            setCsrfToken(sessionData.csrfToken)
          }
          login(
            sessionData.userId,
            sessionData.role,
            sessionData.email ?? "",
            sessionData.fullName,
            sessionData.isActive,
          )
        }
      })
      .catch(() => {})
      .finally(() => {
        if (!cancelled) {
          setIsLoading(false)
        }
      })

    return () => {
      cancelled = true
    }
  }, [login])

  return { isLoading }
}
