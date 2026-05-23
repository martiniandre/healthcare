import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { staffApi } from "./api"
import type { CreateEmployeePayload } from "./types"

export const staffQueryKeys = {
  all: ["staff"] as const,
  lists: () => [...staffQueryKeys.all, "list"] as const,
}

export const useStaffListQuery = () => {
  return useQuery({
    queryKey: staffQueryKeys.lists(),
    queryFn: () => staffApi.listEmployees(),
  })
}

export const useCreateEmployeeMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateEmployeePayload) => staffApi.createEmployee(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: staffQueryKeys.lists(),
      })
    },
  })
}
