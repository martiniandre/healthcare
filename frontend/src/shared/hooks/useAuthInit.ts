import { useEffect, useState } from "react"
import { authApi } from "../../modules/auth/api"
import { useAuthStore } from "../store/auth_store"

export const useAuthInit = () => {
  const [isLoading, setIsLoading] = useState(true)
  const login = useAuthStore((state) => state.login)

  useEffect(() => {
    let cancelled = false

    authApi
      .me()
      .then((sessionData) => {
        if (!cancelled) {
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
