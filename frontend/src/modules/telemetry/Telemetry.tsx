import { useState, useEffect, useRef } from "react"
import type { FormEvent } from "react"
import { Card } from "../../shared/components/ui/Card"
import { Button } from "../../shared/components/ui/Button"
import { 
  Activity, 
  Heart, 
  Thermometer, 
  AlertTriangle, 
  Volume2, 
  VolumeX, 
  Bell,
  Lock,
  Unlock,
  KeyRound,
  Hospital
} from "lucide-react"

import { CardiacCondition, BedStatus } from "../../shared/types"
import { 
  useTelemetryRoomsQuery, 
  useUnlockRoomMutation, 
  useTelemetryBedsQuery, 
  useUpdateBedConditionMutation 
} from "./queries"
import { toast } from "../../shared/store/toast_store"

export const Telemetry = () => {
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
      toast.success("Simulação clínica propagada com sucesso!")
    } catch {
      toast.error("Falha ao propagar simulação clínica.")
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
      setPasscodeError("Senha de Acesso incorreta. Verifique a escala do plantão.")
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
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div className="text-left">
          <h2 className="text-xl font-black text-gray-900 leading-none flex items-center gap-2">
            <Activity className="w-5 h-5 text-primary animate-pulse-glow" />
            Central de Monitoramento e Telemetria
          </h2>
          <span className="text-xs text-muted mt-1.5 block">
            Monitoramento de sinais vitais organizado por alas e salas protegidas por senha de plantão
          </span>
        </div>

        <div className="flex items-center gap-3">
          <Button
            variantType="outline"
            onClick={() => setIsMuted(!isMuted)}
            className="px-3 gap-2 text-xs"
          >
            {isMuted ? (
              <>
                <VolumeX className="w-4 h-4 text-red-500" />
                Alarmes Silenciados
              </>
            ) : (
              <>
                <Volume2 className="w-4 h-4 text-emerald-500 animate-pulse" />
                Alarme Sonoro Ativo
              </>
            )}
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {rooms.map((roomItem) => {
          const isSelected = roomItem.id === selectedRoomId
          const isUnlocked = unlockedRoomIds.includes(roomItem.id)

          return (
            <Card
              key={roomItem.id}
              onClick={() => handleSelectRoom(roomItem.id)}
              className={`p-4 cursor-pointer text-left transition-all duration-200 border relative ${
                isSelected 
                  ? "border-primary bg-primary/5 shadow-sm" 
                  : "border-border hover:border-gray-300 bg-white"
              }`}
            >
              <div className="flex items-center justify-between">
                <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">
                  Ala Monitorada
                </span>
                <span className={`inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full border ${
                  isUnlocked 
                    ? "bg-emerald-50 text-emerald-600 border-emerald-100" 
                    : "bg-amber-50 text-amber-600 border-amber-100"
                }`}>
                  {isUnlocked ? (
                    <>
                      <Unlock className="w-3 h-3" />
                      Desbloqueada
                    </>
                  ) : (
                    <>
                      <Lock className="w-3 h-3" />
                      Protegida
                    </>
                  )}
                </span>
              </div>

              <h4 className="text-sm font-extrabold text-gray-900 mt-2 block">
                {roomItem.name}
              </h4>
              <span className="text-[11px] text-gray-400 block mt-1 leading-normal">
                {roomItem.description}
              </span>

              {isUnlocked && (
                <button
                  onClick={(event) => {
                    event.stopPropagation()
                    handleLockRoom(roomItem.id)
                  }}
                  className="absolute bottom-4 right-4 text-xs text-red-500 hover:text-red-700 transition-colors font-bold"
                >
                  Bloquear
                </button>
              )}
            </Card>
          )
        })}
      </div>

      {!isCurrentRoomUnlocked ? (
        <Card className="flex-1 p-8 border border-border bg-gray-50/50 flex flex-col items-center justify-center text-center gap-5 min-h-[400px]">
          <div className="bg-amber-50 p-4 rounded-full border border-amber-100 text-amber-500 animate-pulse-glow">
            <Lock className="w-10 h-10" />
          </div>

          <div className="max-w-md flex flex-col gap-1">
            <h3 className="text-md font-extrabold text-gray-900">
              Sala com Acesso Restrito
            </h3>
            <p className="text-xs text-gray-500 leading-normal">
              Esta ala contém dados clínicos protegidos (HIPAA/LGPD). Digite a senha da escala de plantão da <strong className="text-gray-800">{activeRoom?.name || "Sala Selecionada"}</strong> para liberar a telemetria.
            </p>
          </div>

          <form onSubmit={handleUnlockRoom} className="w-full max-w-[320px] flex flex-col gap-3">
            <div className="flex flex-col gap-1.5 text-left">
              <label className="text-[10px] font-bold text-gray-500 uppercase tracking-wider">Senha de Escala</label>
              <div className="relative">
                <KeyRound className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
                <input
                  type="password"
                  placeholder="Digite o código da sala..."
                  value={passcodeInput}
                  onChange={(e) => setPasscodeInput(e.target.value)}
                  className="w-full bg-white border border-border rounded-lg pl-9 pr-4 py-2.5 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200 font-mono"
                  required
                />
              </div>
              {passcodeError && (
                <span className="text-[10px] text-red-500 font-bold mt-1 block">
                  {passcodeError}
                </span>
              )}
            </div>

            <Button
              type="submit"
              variantType="primary"
              disabled={unlockRoomMutation.isPending}
              className="w-full py-2.5 text-xs font-bold gap-2"
            >
              <Unlock className="w-4 h-4" />
              Liberar Monitoramento
            </Button>
          </form>

          <div className="text-[10px] text-gray-400 font-semibold border-t border-border/80 pt-4 w-full max-w-sm mt-2 flex items-center justify-center gap-1.5">
            <Hospital className="w-3.5 h-3.5" />
            Conselho Técnico de Enfermagem e Medicina Intensiva
          </div>
        </Card>
      ) : beds.length === 0 ? (
        <Card className="flex-1 p-8 border border-border bg-white flex flex-col items-center justify-center text-center gap-3 min-h-[400px]">
          <Activity className="w-8 h-8 text-gray-400 animate-pulse" />
          <span className="text-xs text-gray-500 font-bold">Nenhum leito ativo nesta sala de monitoramento.</span>
        </Card>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 animate-fade-in">
          <div className="lg:col-span-2 flex flex-col gap-6">
            {activeBed ? (
              <Card className="p-4 flex flex-col gap-4">
                <div className="flex items-center justify-between border-b border-border pb-3 flex-wrap gap-2">
                  <div className="text-left">
                    <span className="text-xs font-bold text-primary uppercase tracking-wider">{activeRoom?.name}</span>
                    <h3 className="text-md font-bold text-gray-900">{activeBed.bedNumber} • {activeBed.patientName}</h3>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded font-bold">
                      {activeBed.age} anos · {activeBed.gender}
                    </span>
                    {activeBed.status !== BedStatus.Normal && (
                      <span className="text-xs bg-red-50 text-red-600 px-2.5 py-1 rounded-full font-bold flex items-center gap-1 border border-red-200">
                        <AlertTriangle className="w-3.5 h-3.5" />
                        Crítico
                      </span>
                    )}
                  </div>
                </div>

                <div className="relative border border-border rounded-xl overflow-hidden bg-slate-950 p-1">
                  <canvas
                    ref={canvasRef}
                    width={700}
                    height={260}
                    className="block w-full max-w-full rounded-lg"
                  />
                  <div className="absolute top-4 left-4 bg-black/60 backdrop-blur-md px-3 py-1.5 rounded text-[10px] text-emerald-400 font-mono border border-emerald-500/20">
                    ECG - DERIVAÇÃO II (DII)
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  <div className="bg-gray-50 border border-border p-4 rounded-xl flex items-center justify-between">
                    <div className="text-left">
                      <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Freq. Cardíaca</span>
                      <span className="text-2xl font-black text-gray-900 mt-1 block">
                        {activeBed.bpm > 0 ? `${activeBed.bpm} ` : "--- "}
                        <span className="text-xs font-normal text-gray-400">BPM</span>
                      </span>
                    </div>
                    <Heart className={`w-8 h-8 text-red-500 shrink-0 ${activeBed.bpm > 110 ? "animate-pulse" : ""}`} />
                  </div>

                  <div className="bg-gray-50 border border-border p-4 rounded-xl flex items-center justify-between">
                    <div className="text-left">
                      <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Saturação O₂</span>
                      <span className="text-2xl font-black text-gray-900 mt-1 block">
                        {activeBed.spo2 > 0 ? `${activeBed.spo2}%` : "---%"}
                      </span>
                    </div>
                    <Activity className="w-8 h-8 text-emerald-500 shrink-0" />
                  </div>

                  <div className="bg-gray-50 border border-border p-4 rounded-xl flex items-center justify-between">
                    <div className="text-left">
                      <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Temperatura</span>
                      <span className="text-2xl font-black text-gray-900 mt-1 block">
                        {activeBed.temperature.toFixed(1)}°C
                      </span>
                    </div>
                    <Thermometer className="w-8 h-8 text-sky-500 shrink-0" />
                  </div>
                </div>

                <div className="border-t border-border pt-4 text-left flex flex-col gap-2">
                  <span className="text-xs font-bold text-gray-600 block">Injetar Simulação Clínica:</span>
                  <div className="flex gap-2 flex-wrap">
                    <Button 
                      variantType={activeBed.condition === CardiacCondition.Normal ? "primary" : "outline"} 
                      onClick={() => updateSelectedBedCondition(CardiacCondition.Normal)}
                      className="px-3.5 py-2 text-xs"
                    >
                      Simular Ritmo Normal
                    </Button>
                    <Button 
                      variantType={activeBed.condition === CardiacCondition.Bradycardia ? "primary" : "outline"} 
                      onClick={() => updateSelectedBedCondition(CardiacCondition.Bradycardia)}
                      className="px-3.5 py-2 text-xs"
                    >
                      Simular Bradicardia
                    </Button>
                    <Button 
                      variantType={activeBed.condition === CardiacCondition.Tachycardia ? "primary" : "outline"} 
                      onClick={() => updateSelectedBedCondition(CardiacCondition.Tachycardia)}
                      className="px-3.5 py-2 text-xs"
                    >
                      Simular Taquicardia
                    </Button>
                    <Button 
                      variantType={activeBed.condition === CardiacCondition.CardiacArrest ? "danger" : "outline"} 
                      onClick={() => updateSelectedBedCondition(CardiacCondition.CardiacArrest)}
                      className="px-3.5 py-2 text-xs text-red-500 hover:text-white"
                    >
                      Simular Parada Cardíaca
                    </Button>
                  </div>
                </div>
              </Card>
            ) : (
              <Card className="flex-1 p-8 border border-border bg-gray-50/50 flex flex-col items-center justify-center text-center gap-4 min-h-[400px]">
                <div className="bg-primary/5 p-4 rounded-full border border-primary/10 text-primary">
                  <Activity className="w-8 h-8 animate-pulse" />
                </div>
                <div className="max-w-md flex flex-col gap-1">
                  <h3 className="text-md font-extrabold text-gray-900">
                    Nenhum Leito Selecionado
                  </h3>
                  <p className="text-xs text-gray-500 leading-normal">
                    Selecione um leito na lista lateral para monitorar os sinais vitais em tempo real e injetar simulações clínicas.
                  </p>
                </div>
              </Card>
            )}
          </div>

          <div className="flex flex-col gap-4 lg:col-span-1 text-left">
            <Card className="p-4 flex flex-col gap-4">
              <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2 border-b border-border pb-3">
                <Bell className="w-4 h-4 text-primary animate-pulse-glow" />
                Leitos Disponíveis ({beds.length})
              </h3>

              <div className="flex flex-col gap-3">
                {beds.map((bedItem) => {
                  const isSelected = bedItem.id === selectedBedId
                  return (
                    <div
                      key={bedItem.id}
                      onClick={() => setSelectedBedId(bedItem.id)}
                      className={`cursor-pointer border p-3.5 rounded-xl transition-all duration-200 ${
                        isSelected 
                          ? "bg-primary/5 border-primary" 
                          : bedItem.status === BedStatus.Danger
                            ? "bg-red-50 border-red-200 hover:border-red-300"
                            : "bg-white border-border hover:border-gray-300"
                      }`}
                    >
                      <div className="flex items-center justify-between">
                        <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider">
                          {bedItem.bedNumber}
                        </span>
                        <span className={`text-[10px] px-2 py-0.5 rounded font-black uppercase ${
                          bedItem.status === BedStatus.Danger 
                            ? "bg-red-100 text-red-600" 
                            : bedItem.status === BedStatus.Warning
                              ? "bg-amber-100 text-amber-600"
                              : "bg-emerald-100 text-emerald-600"
                        }`}>
                          {bedItem.condition}
                        </span>
                      </div>

                      <h4 className="text-sm font-bold text-gray-800 mt-1">
                        {bedItem.patientName}
                      </h4>

                      <div className="flex items-center justify-between mt-3 text-xs">
                        <span className="text-gray-500">FC: <strong className="text-gray-800">{bedItem.bpm > 0 ? `${bedItem.bpm} BPM` : "---"}</strong></span>
                        <span className="text-gray-500">SpO₂: <strong className="text-gray-800">{bedItem.spo2 > 0 ? `${bedItem.spo2}%` : "---"}</strong></span>
                      </div>
                    </div>
                  )
                })}
              </div>
            </Card>
          </div>
        </div>
      )}
    </div>
  )
}
