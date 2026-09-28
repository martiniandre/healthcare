import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { staffApi } from "./api"
import type { CreateEmployeePayload, SetEmployeeStatusPayload } from "./types"

export const staffQueryKeys = {
  all: ["staff"] as const,
  lists: () => [...staffQueryKeys.all, "list"] as const,
  departments: () => [...staffQueryKeys.all, "departments"] as const,
}

export const useStaffListQuery = (search?: string, role?: string) => {
  return useQuery({
    queryKey: [...staffQueryKeys.lists(), { search, role }],
    queryFn: () => staffApi.listEmployees(search, role),
  })
}

export const useDepartmentsQuery = () => {
  return useQuery({
    queryKey: staffQueryKeys.departments(),
    queryFn: () => staffApi.listDepartments(),
    staleTime: 1000 * 60 * 30,
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

export const useSetEmployeeStatusMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: SetEmployeeStatusPayload) => staffApi.setEmployeeStatus(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: staffQueryKeys.lists(),
      })
    },
  })
}

export const useDeleteEmployeeMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (employeeId: string) => staffApi.deleteEmployee(employeeId),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: staffQueryKeys.lists(),
      })
    },
  })
}
