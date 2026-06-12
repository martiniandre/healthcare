import { useState, useEffect, useRef } from "react"
import type { FormEvent } from "react"
import { useTranslation } from "react-i18next"
import { Card } from "../../shared/components/ui/Card"
import { Activity } from "lucide-react"

import { CardiacCondition, BedStatus } from "../../shared/types"
import { 
  useTelemetryRoomsQuery, 
  useUnlockRoomMutation, 
  useTelemetryBedsQuery, 
  useUpdateBedConditionMutation 
} from "./queries"
import { toast } from "../../shared/store/toast_store"

import { TelemetryHeader } from "./components/TelemetryHeader"
import { TelemetryRoomList } from "./components/TelemetryRoomList"
import { TelemetryRestrictedState } from "./components/TelemetryRestrictedState"
import { TelemetryBedMonitor } from "./components/TelemetryBedMonitor"
import { TelemetryBedList } from "./components/TelemetryBedList"

export const Telemetry = () => {
  const { t } = useTranslation()
  const [selectedRoomId, setSelectedRoomId] = useState<string>("room-1")
  const [unlockedRoomIds, setUnlockedRoomIds] = useState<string[]>(["room-1"])
  const [selectedBedId, setSelectedBedId] = useState<string | null>(null)
  
  const [passcodeInput, setPasscodeInput] = useState<string>("")
  const [passcodeError, setPasscodeError] = useState<string>("")
  const [isMuted, setIsMuted] = useState<boolean>(true)

  const { data: rooms = [] } = useTelemetryRoomsQuery()
  const { data: beds = [] } = useTelemetryBedsQuery(selectedRoomId)
  
  const unlockRoomMutation = useUnlockRoomMutation()
  const updateBedConditionMutation = useUpdateBedConditionMutation()

  const activeRoom = rooms.find((roomItem) => roomItem.id === selectedRoomId) || rooms[0]
  const isCurrentRoomUnlocked = unlockedRoomIds.includes(selectedRoomId)

  const activeBed = beds.find((bedItem) => bedItem.id === selectedBedId) || null

  const canvasRef = useRef<HTMLCanvasElement | null>(null)
  const animationFrameIdRef = useRef<number | null>(null)
  const traceOffsetRef = useRef<number>(0)
  const ecgPointsRef = useRef<number[]>([])

  const updateSelectedBedCondition = async (newCondition: CardiacCondition) => {
    if (!activeBed) {
      return
    }

    let dynamicBpm = 75
    let dynamicSpo2 = 98
    let dynamicStatus: BedStatus = BedStatus.Normal

    if (newCondition === CardiacCondition.Bradycardia) {
      dynamicBpm = 45
      dynamicSpo2 = 93
      dynamicStatus = BedStatus.Warning
    } else if (newCondition === CardiacCondition.Tachycardia) {
      dynamicBpm = 135
      dynamicSpo2 = 95
      dynamicStatus = BedStatus.Warning
    } else if (newCondition === CardiacCondition.CardiacArrest) {
      dynamicBpm = 0
      dynamicSpo2 = 0
      dynamicStatus = BedStatus.Danger
    }

    try {
      await updateBedConditionMutation.mutateAsync({
        bedId: activeBed.id,
        bpm: dynamicBpm,
        spo2: dynamicSpo2,
        temperature: activeBed.temperature,
        status: dynamicStatus,
        condition: newCondition,
      })
      toast.success(t("telemetry.toast.simulateSuccess"))
    } catch {
      toast.error(t("telemetry.toast.simulateError"))
    }
  }

  useEffect(() => {
    const renderEcgWaveform = () => {
      const activeCanvas = canvasRef.current
      if (!activeCanvas || !isCurrentRoomUnlocked || !activeBed) {
        return
      }

      const canvasContext = activeCanvas.getContext("2d")
      if (!canvasContext) {
        return
      }

      const canvasWidth = activeCanvas.width
      const canvasHeight = activeCanvas.height

      canvasContext.fillStyle = "#090d16"
      canvasContext.fillRect(0, 0, canvasWidth, canvasHeight)

      canvasContext.strokeStyle = "rgba(16, 185, 129, 0.15)"
      canvasContext.lineWidth = 1
      const gridSize = 15

      for (let xCoordinate = 0; xCoordinate < canvasWidth; xCoordinate += gridSize) {
        canvasContext.beginPath()
        canvasContext.moveTo(xCoordinate, 0)
        canvasContext.lineTo(xCoordinate, canvasHeight)
        canvasContext.stroke()
      }

      for (let yCoordinate = 0; yCoordinate < canvasHeight; yCoordinate += gridSize) {
        canvasContext.beginPath()
        canvasContext.moveTo(0, yCoordinate)
        canvasContext.lineTo(canvasWidth, yCoordinate)
        canvasContext.stroke()
      }

      const totalPlotPoints = 300
      if (ecgPointsRef.current.length < totalPlotPoints) {
        ecgPointsRef.current = new Array(totalPlotPoints).fill(canvasHeight / 2)
      }

      let cyclePeriod = 60
      if (activeBed.condition === CardiacCondition.Bradycardia) {
        cyclePeriod = 110
      } else if (activeBed.condition === CardiacCondition.Tachycardia) {
        cyclePeriod = 35
      }

      let ecgSignalValue: number

      if (activeBed.condition !== CardiacCondition.CardiacArrest) {
        const currentPhase = traceOffsetRef.current % cyclePeriod

        if (currentPhase < 6) {
          ecgSignalValue = canvasHeight / 2
        } else if (currentPhase >= 6 && currentPhase < 12) {
          const sineProgress = (currentPhase - 6) * Math.PI / 6
          ecgSignalValue = canvasHeight / 2 - Math.sin(sineProgress) * 12
        } else if (currentPhase >= 12 && currentPhase < 16) {
          ecgSignalValue = canvasHeight / 2
        } else if (currentPhase >= 16 && currentPhase < 18) {
          ecgSignalValue = canvasHeight / 2 + 8
        } else if (currentPhase >= 18 && currentPhase < 21) {
          const rWaveProgress = (currentPhase - 18) / 3
          ecgSignalValue = canvasHeight / 2 + 8 - rWaveProgress * 65
        } else if (currentPhase >= 21 && currentPhase < 24) {
          const sWaveProgress = (currentPhase - 21) / 3
          ecgSignalValue = canvasHeight / 2 - 57 + sWaveProgress * 78
        } else if (currentPhase >= 24 && currentPhase < 27) {
          const baselineReturn = (currentPhase - 24) / 3
          ecgSignalValue = canvasHeight / 2 + 21 - baselineReturn * 21
        } else if (currentPhase >= 27 && currentPhase < 33) {
          ecgSignalValue = canvasHeight / 2
        } else if (currentPhase >= 33 && currentPhase < 43) {
          const tWaveProgress = (currentPhase - 33) * Math.PI / 10
          ecgSignalValue = canvasHeight / 2 - Math.sin(tWaveProgress) * 18
        } else {
          ecgSignalValue = canvasHeight / 2
        }
      } else {
        const flatlineFluctuation = Math.sin(traceOffsetRef.current * 0.15) * 0.8
        ecgSignalValue = canvasHeight / 2 + flatlineFluctuation
      }

      ecgPointsRef.current.push(ecgSignalValue)
      ecgPointsRef.current.shift()

      canvasContext.beginPath()
      canvasContext.strokeStyle = activeBed.condition === CardiacCondition.CardiacArrest ? "#ef4444" : "#10b981"
      canvasContext.lineWidth = 2.5
      canvasContext.lineJoin = "round"

      for (let pointIndex = 0; pointIndex < ecgPointsRef.current.length; pointIndex++) {
        const xPixelCoordinate = (pointIndex / (ecgPointsRef.current.length - 1)) * canvasWidth
        const yPixelCoordinate = ecgPointsRef.current[pointIndex]

        if (pointIndex === 0) {
          canvasContext.moveTo(xPixelCoordinate, yPixelCoordinate)
        } else {
          canvasContext.lineTo(xPixelCoordinate, yPixelCoordinate)
        }
      }
      canvasContext.stroke()

      traceOffsetRef.current += 1
      animationFrameIdRef.current = requestAnimationFrame(renderEcgWaveform)
    }

    if (isCurrentRoomUnlocked && activeBed) {
      animationFrameIdRef.current = requestAnimationFrame(renderEcgWaveform)
    }
    return () => {
      if (animationFrameIdRef.current) {
        cancelAnimationFrame(animationFrameIdRef.current)
      }
    }
  }, [isCurrentRoomUnlocked, activeBed])

  const handleUnlockRoom = async (event: FormEvent) => {
    event.preventDefault()
    if (!activeRoom) {
      return
    }

    try {
      await unlockRoomMutation.mutateAsync({
        roomIdValue: selectedRoomId,
        passcodeValue: passcodeInput,
      })

      setUnlockedRoomIds((previousUnlocked) => [...previousUnlocked, selectedRoomId])
      setPasscodeInput("")
      setPasscodeError("")
    } catch {
      setPasscodeError(t("telemetry.toast.unlockError"))
    }
  }

  const handleLockRoom = (roomId: string) => {
    setUnlockedRoomIds((previousUnlocked) => previousUnlocked.filter((id) => id !== roomId))
  }

  const handleSelectRoom = (roomId: string) => {
    setSelectedRoomId(roomId)
    setPasscodeError("")
    setPasscodeInput("")
    setSelectedBedId(null)
  }

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full select-none">
      <TelemetryHeader 
        isMuted={isMuted} 
        setIsMuted={setIsMuted} 
      />

      <TelemetryRoomList 
        rooms={rooms}
        selectedRoomId={selectedRoomId}
        unlockedRoomIds={unlockedRoomIds}
        handleSelectRoom={handleSelectRoom}
        handleLockRoom={handleLockRoom}
      />

      {!isCurrentRoomUnlocked ? (
        <TelemetryRestrictedState 
          activeRoomName={activeRoom?.name}
          passcodeInput={passcodeInput}
          setPasscodeInput={setPasscodeInput}
          passcodeError={passcodeError}
          isPending={unlockRoomMutation.isPending}
          handleUnlockRoom={handleUnlockRoom}
        />
      ) : beds.length === 0 ? (
        <Card className="flex-1 p-8 border border-border bg-white flex flex-col items-center justify-center text-center gap-3 min-h-[400px]">
          <Activity className="w-8 h-8 text-gray-400 animate-pulse" />
          <span className="text-xs text-gray-500 font-bold">{t("telemetry.noBeds")}</span>
        </Card>
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
          />
        </div>
      )}
    </div>
  )
}

