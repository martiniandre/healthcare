import { http } from "../../shared/utils/http"
import type { TelemetryRoom, BedPatient, UnlockRoomResponseDto, UpdateBedConditionPayload, UpdateBedConditionResponseDto } from "./types"

export const telemetryApi = {
  getRooms: async (): Promise<TelemetryRoom[]> => {
    return http.get<TelemetryRoom[]>("/telemetry/rooms")
  },

  unlockRoom: async (roomIdValue: string, passcodeValue: string): Promise<UnlockRoomResponseDto> => {
    return http.post<UnlockRoomResponseDto>(`/telemetry/rooms/${roomIdValue}/unlock`, {
      passcode: passcodeValue,
    })
  },

  getBeds: async (roomIdValue: string): Promise<BedPatient[]> => {
    return http.get<BedPatient[]>(`/telemetry/rooms/${roomIdValue}/beds`)
  },

  updateBedCondition: async (payload: UpdateBedConditionPayload): Promise<UpdateBedConditionResponseDto> => {
    return http.post<UpdateBedConditionResponseDto>(`/telemetry/beds/${payload.bedId}/condition`, {
      bpm: payload.bpm,
      spo2: payload.spo2,
      temperature: payload.temperature,
      status: payload.status,
      condition: payload.condition,
    })
  },
}
