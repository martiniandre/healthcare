import { CardiacCondition, BedStatus } from "../../shared/types"

export interface TelemetryRoom {
  id: string
  name: string
  description: string
}

export interface BedPatient {
  id: string
  roomId: string
  bedNumber: string
  patientName: string
  age: number
  gender: string
  bpm: number
  spo2: number
  temperature: number
  status: BedStatus
  condition: CardiacCondition
}

export interface UnlockRoomResponseDto {
  success: boolean
  roomName: string
}

export interface UpdateBedConditionPayload {
  bedId: string
  bpm: number
  spo2: number
  temperature: number
  status: BedStatus
  condition: CardiacCondition
}

export interface UpdateBedConditionResponseDto {
  success: boolean
}
