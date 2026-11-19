import { useState, useRef } from "react"
import { useEcgTrace } from "./hooks/useEcgTrace"
import { simulatedVitalsFor } from "./utils/ecg_waveform"
import type { FormEvent } from "react"
import { useTranslation } from "react-i18next"
import { Activity, AlertCircle, Info } from "lucide-react"

import { CardiacCondition } from "../../shared/types"
import {
  useTelemetryRoomsQuery,
  useUnlockRoomMutation,
  useTelemetryBedsQuery,
  useUpdateBedConditionMutation
} from "./queries"
import { toast } from "../../shared/store/toast_store"
import { PageContainer } from "../../shared/components/ui/PageContainer"
import { EmptyState } from "../../shared/components/ui/EmptyState"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { Button } from "@/shared/components/ui/button"
import { Skeleton } from "@/shared/components/ui/skeleton"
import { Spinner } from "@/shared/components/ui/spinner"

import { TelemetryHeader } from "./components/TelemetryHeader"
import { TelemetryRoomList } from "./components/TelemetryRoomList"
import { TelemetryRestrictedState } from "./components/TelemetryRestrictedState"
import { TelemetryBedMonitor } from "./components/TelemetryBedMonitor"
import { TelemetryBedList } from "./components/TelemetryBedList"
import { useAlarmMute } from "./hooks/useAlarmMute"
import { useBedHistory } from "./hooks/useBedHistory"

const LAST_TELEMETRY_ROOM_STORAGE_KEY = "healthcare.telemetry.lastRoom"

const readLastTelemetryRoom = (): string | null => {
  try {
    return window.localStorage.getItem(LAST_TELEMETRY_ROOM_STORAGE_KEY)
  } catch {
    return null
  }
}

const rememberLastTelemetryRoom = (roomId: string) => {
  try {
    window.localStorage.setItem(LAST_TELEMETRY_ROOM_STORAGE_KEY, roomId)
  } catch {
    return
  }
}

export const Telemetry = () => {
  const { t } = useTranslation("telemetry")
  const [selectedRoomId, setSelectedRoomId] = useState<string | null>(readLastTelemetryRoom)
  const [unlockedRoomIds, setUnlockedRoomIds] = useState<string[]>([])
  const [selectedBedId, setSelectedBedId] = useState<string | null>(null)
  
  const [passcode, setPasscode] = useState<{
    input: string
    error: string
  }>({ input: "", error: "" })
  const { isMuted, setIsMuted } = useAlarmMute()

  const {
    data: rooms = [],
    isLoading: isRoomsLoading,
    isError: isRoomsError,
    refetch: refetchRooms,
  } = useTelemetryRoomsQuery()
  
  const isCurrentRoomUnlocked = selectedRoomId ? unlockedRoomIds.includes(selectedRoomId) : false
  const {
    data: beds = [],
    isLoading: isBedsLoading,
    isError: isBedsError,
    refetch: refetchBeds,
    dataUpdatedAt: bedsUpdatedAt,
  } = useTelemetryBedsQuery(selectedRoomId, isCurrentRoomUnlocked)
  const bedHistory = useBedHistory(beds, bedsUpdatedAt)
  
  const unlockRoomMutation = useUnlockRoomMutation()
  const updateBedConditionMutation = useUpdateBedConditionMutation()

  const activeRoom = rooms.find((roomItem) => roomItem.id === selectedRoomId) || null

  const activeBed = beds.find((bedItem) => bedItem.id === selectedBedId) || null

  const canvasRef = useRef<HTMLCanvasElement | null>(null)
  useEcgTrace(canvasRef, isCurrentRoomUnlocked && activeBed !== null, activeBed?.condition)

  const updateSelectedBedCondition = async (newCondition: CardiacCondition) => {
    if (!activeBed) {
      return
    }

    const simulatedVitals = simulatedVitalsFor(newCondition)

    try {
      await updateBedConditionMutation.mutateAsync({
        bedId: activeBed.id,
        bpm: simulatedVitals.bpm,
        spo2: simulatedVitals.spo2,
        temperature: activeBed.temperature,
        status: simulatedVitals.status,
        condition: newCondition,
      })
      toast.success(t("toast.simulateSuccess"))
    } catch {
      toast.error(t("toast.simulateError"))
    }
  }

  const handleUnlockRoom = async (event: FormEvent) => {
    event.preventDefault()
    if (!activeRoom || !selectedRoomId) {
      return
    }

    try {
      await unlockRoomMutation.mutateAsync({
        roomIdValue: selectedRoomId,
        passcodeValue: passcode.input,
      })

      setUnlockedRoomIds((previousUnlocked) => [...previousUnlocked, selectedRoomId])
      setPasscode({ input: "", error: "" })
    } catch {
      setPasscode((prev) => ({ ...prev, error: t("toast.unlockError") }))
    }
  }

  const handleLockRoom = (roomId: string) => {
    setUnlockedRoomIds((previousUnlocked) => previousUnlocked.filter((id) => id !== roomId))
  }

  const handleSelectRoom = (roomId: string) => {
    rememberLastTelemetryRoom(roomId)
    setSelectedRoomId(roomId)
    setPasscode({ input: "", error: "" })
    setSelectedBedId(null)
  }

  return (
    <PageContainer className="select-none">
      <TelemetryHeader
        isMuted={isMuted}
        setIsMuted={setIsMuted}
      />

      {isRoomsLoading ? (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3" role="status" aria-label={t("loadingRooms")}>
          <Skeleton className="h-28 rounded-lg" />
          <Skeleton className="h-28 rounded-lg" />
          <Skeleton className="h-28 rounded-lg" />
        </div>
      ) : isRoomsError ? (
        <Alert variant="destructive">
          <AlertCircle aria-hidden="true" />
          <AlertDescription className="flex flex-wrap items-center gap-3">
            <span>{t("roomsError")}</span>
            <Button variant="outline" size="sm" onClick={() => refetchRooms()}>
              {t("retry")}
            </Button>
          </AlertDescription>
        </Alert>
      ) : rooms.length === 0 ? (
        <EmptyState icon={Activity} title={t("noRooms")} />
      ) : (
        <div className="flex flex-col gap-2">
          <TelemetryRoomList
            rooms={rooms}
            selectedRoomId={selectedRoomId}
            unlockedRoomIds={unlockedRoomIds}
            handleSelectRoom={handleSelectRoom}
            handleLockRoom={handleLockRoom}
          />
          <p className="flex items-start gap-1.5 text-xs text-muted-foreground">
            <Info className="mt-0.5 size-3.5 shrink-0" aria-hidden="true" />
            {t("protectedHint")}
          </p>
        </div>
      )}

      {!selectedRoomId || (!isRoomsLoading && !activeRoom) ? (
        <EmptyState
          icon={Activity}
          title={t("noRoomSelected")}
          description={t("noRoomSelectedDesc")}
          className="py-10"
        />
      ) : !isCurrentRoomUnlocked ? (
        <TelemetryRestrictedState
          activeRoomName={activeRoom?.name}
          passcodeInput={passcode.input}
          setPasscodeInput={(val) => setPasscode((prev) => ({ ...prev, input: val }))}
          passcodeError={passcode.error}
          isPending={unlockRoomMutation.isPending}
          handleUnlockRoom={handleUnlockRoom}
        />
      ) : isBedsLoading ? (
        <div className="flex min-h-[400px] flex-1 items-center justify-center gap-2 text-sm text-muted-foreground">
          <Spinner />
          <span>{t("loadingBeds")}</span>
        </div>
      ) : isBedsError ? (
        <Alert variant="destructive">
          <AlertCircle aria-hidden="true" />
          <AlertDescription className="flex flex-wrap items-center gap-3">
            <span>{t("bedsError")}</span>
            <Button variant="outline" size="sm" onClick={() => refetchBeds()}>
              {t("retry")}
            </Button>
          </AlertDescription>
        </Alert>
      ) : beds.length === 0 ? (
        <EmptyState
          icon={Activity}
          title={t("noBeds")}
          className="py-10"
        />
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 animate-fade-in">
          <div className="lg:col-span-2 flex flex-col gap-6">
            <TelemetryBedMonitor 
              activeBed={activeBed}
              activeRoomName={activeRoom?.name}
              canvasRef={canvasRef}
              updateSelectedBedCondition={updateSelectedBedCondition}
            />
          </div>

          <TelemetryBedList
            beds={beds}
            selectedBedId={selectedBedId}
            setSelectedBedId={setSelectedBedId}
            history={bedHistory}
          />
        </div>
      )}
    </PageContainer>
  )
}

