import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { authApi } from "./api"
import type { LoginRequest } from "./types"
import { setCsrfToken } from "../../shared/services/api"
import { useAuthStore } from "../../shared/store/auth_store"

export const authKeys = {
  all: ["auth"] as const,
  me: () => [...authKeys.all, "me"] as const,
  permissions: () => [...authKeys.all, "permissions"] as const,
}

export const useLoginMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: LoginRequest) => authApi.login(data),
    onSuccess: (response) => {
      if (response.csrfToken) {
        setCsrfToken(response.csrfToken)
      }
      queryClient.invalidateQueries({ queryKey: authKeys.me() })
      queryClient.invalidateQueries({ queryKey: authKeys.permissions() })
    },
  })
}

export const usePermissionsQuery = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  return useQuery({
    queryKey: authKeys.permissions(),
    queryFn: () => authApi.getPermissions(),
    enabled: isAuthenticated,
    staleTime: Infinity,
  })
}
