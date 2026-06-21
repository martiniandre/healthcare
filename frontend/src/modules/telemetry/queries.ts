import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { telemetryApi } from "./api"
import type { UpdateBedConditionPayload } from "./types"

export const telemetryQueryKeys = {
  all: ["telemetry"] as const,
  rooms: () => [...telemetryQueryKeys.all, "rooms"] as const,
  beds: (roomIdValue: string) => [...telemetryQueryKeys.all, "beds", roomIdValue] as const,
}

export const useTelemetryRoomsQuery = () => {
  return useQuery({
    queryKey: telemetryQueryKeys.rooms(),
    queryFn: () => telemetryApi.getRooms(),
  })
}

export const useUnlockRoomMutation = () => {
  return useMutation({
    mutationFn: (payload: { roomIdValue: string; passcodeValue: string }) =>
      telemetryApi.unlockRoom(payload.roomIdValue, payload.passcodeValue),
  })
}

export const useTelemetryBedsQuery = (roomIdValue: string | null, isEnabled: boolean = true) => {
  return useQuery({
    queryKey: telemetryQueryKeys.beds(roomIdValue || ""),
    queryFn: () => telemetryApi.getBeds(roomIdValue || ""),
    enabled: !!roomIdValue && isEnabled,
    refetchInterval: 1000,
  })
}

export const useUpdateBedConditionMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: UpdateBedConditionPayload) =>
      telemetryApi.updateBedCondition(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: telemetryQueryKeys.all,
      })
    },
  })
}
