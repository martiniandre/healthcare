import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { authApi } from "./api"
import type { LoginRequest, RegisterRequest } from "./types"
import { setCsrfToken, clearCsrfToken } from "../../shared/services/api"

export const authKeys = {
  all: ["auth"] as const,
  me: () => [...authKeys.all, "me"] as const,
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
    },
  })
}

export const useRegisterMutation = () => {
  return useMutation({
    mutationFn: (data: RegisterRequest) => authApi.register(data),
  })
}

export const useLogoutMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => authApi.logout(),
    onSuccess: () => {
      clearCsrfToken()
      queryClient.clear()
    },
  })
}

export const useCurrentUserQuery = () => {
  return useQuery({
    queryKey: authKeys.me(),
    queryFn: async () => {
      const response = await authApi.me()
      if (response.csrfToken) {
        setCsrfToken(response.csrfToken)
      }
      return response
    },
  })
}
