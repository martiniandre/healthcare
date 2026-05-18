import * as React from "react"
import { useState, useEffect, useRef, useCallback } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { Card } from "../../shared/components/ui/Card"
import { Button } from "../../shared/components/ui/Button"
import { useImagingStudyQuery } from "./queries"
import { 
  ArrowLeft, 
  ZoomIn, 
  Sun, 
  Ruler, 
  UploadCloud, 
  CheckCircle,
  Eye
} from "lucide-react"

export const ImagingWorkspace = () => {
  const { studyId = "" } = useParams<{ studyId: string }>()
  const navigate = useNavigate()

  const [activeTool, setActiveTool] = useState<"zoom" | "contrast" | "ruler">("zoom")
  const [contrastSetting, setContrastSetting] = useState(1.0)
  const [brightnessSetting, setBrightnessSetting] = useState(1.0)
  const [zoomSetting, setZoomSetting] = useState(1.0)
  
  const [uploadPercentage, setUploadPercentage] = useState<number | null>(null)
  const [uploadStatus, setUploadStatus] = useState<string | null>(null)

  const { data: study, isLoading: isStudyLoading } = useImagingStudyQuery(studyId)

  const canvasReference = useRef<HTMLCanvasElement | null>(null)
  const isDraggingReference = useRef(false)
  const dragStartPoint = useRef({ x: 0, y: 0 })

  const rulerStartPoint = useRef<{ x: number; y: number } | null>(null)
  const rulerEndPoint = useRef<{ x: number; y: number } | null>(null)

  const drawDICOMImage = useCallback(() => {
    const activeCanvas = canvasReference.current
    if (!activeCanvas) {
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

    canvasContext.save()
    canvasContext.translate(canvasWidth / 2, canvasHeight / 2)
    canvasContext.scale(zoomSetting, zoomSetting)
    canvasContext.translate(-canvasWidth / 2, -canvasHeight / 2)

    canvasContext.filter = `contrast(${contrastSetting * 100}%) brightness(${brightnessSetting * 100}%)`

    canvasContext.fillStyle = "#1e293b"
    canvasContext.beginPath()
    canvasContext.arc(canvasWidth / 2, canvasHeight / 2, 140, 0, Math.PI * 2)
    canvasContext.fill()

    canvasContext.fillStyle = "#334155"
    canvasContext.fillRect(canvasWidth / 2 - 30, canvasHeight / 2 - 120, 60, 240)
    canvasContext.fillRect(canvasWidth / 2 - 120, canvasHeight / 2 - 30, 240, 60)

    canvasContext.strokeStyle = "#475569"
    canvasContext.lineWidth = 4
    canvasContext.strokeRect(canvasWidth / 2 - 80, canvasHeight / 2 - 80, 160, 160)

    canvasContext.fillStyle = "#64748b"
    canvasContext.beginPath()
    canvasContext.arc(canvasWidth / 2, canvasHeight / 2, 45, 0, Math.PI * 2)
    canvasContext.fill()

    canvasContext.restore()

    canvasContext.fillStyle = "rgba(255, 255, 255, 0.4)"
    canvasContext.font = "10px Outfit, sans-serif"
    canvasContext.fillText("HOSPITAL GERAL - MOCK PACS", 15, 25)
    canvasContext.fillText(`SERIES: 1 • MODALITY: ${study?.modality || "CT"}`, 15, 40)
    canvasContext.fillText(`UID: ${study?.study_instance_uid || "1.2.840.10008"}`, 15, 55)

    canvasContext.fillText(`ZOOM: ${Math.round(zoomSetting * 100)}%`, canvasWidth - 110, 25)
    canvasContext.fillText(`CONTRASTE: ${Math.round(contrastSetting * 100)}%`, canvasWidth - 110, 40)
    canvasContext.fillText(`BRILHO: ${Math.round(brightnessSetting * 100)}%`, canvasWidth - 110, 55)

    if (rulerStartPoint.current) {
      canvasContext.strokeStyle = "#0ea5e9"
      canvasContext.lineWidth = 2
      canvasContext.beginPath()
      canvasContext.moveTo(rulerStartPoint.current.x, rulerStartPoint.current.y)
      
      const activeEnd = rulerEndPoint.current || rulerStartPoint.current
      canvasContext.lineTo(activeEnd.x, activeEnd.y)
      canvasContext.stroke()

      canvasContext.fillStyle = "#0ea5e9"
      canvasContext.beginPath()
      canvasContext.arc(rulerStartPoint.current.x, rulerStartPoint.current.y, 4, 0, Math.PI * 2)
      canvasContext.arc(activeEnd.x, activeEnd.y, 4, 0, Math.PI * 2)
      canvasContext.fill()

      const deltaX = activeEnd.x - rulerStartPoint.current.x
      const deltaY = activeEnd.y - rulerStartPoint.current.y
      const distancePixels = Math.sqrt(deltaX * deltaX + deltaY * deltaY)
      const distanceMm = (distancePixels * 0.28).toFixed(1)

      canvasContext.fillStyle = "#ffffff"
      canvasContext.font = "12px Outfit, sans-serif"
      canvasContext.fillText(`${distanceMm} mm`, activeEnd.x + 10, activeEnd.y + 10)
    }
  }, [study, contrastSetting, brightnessSetting, zoomSetting])

  useEffect(() => {
    drawDICOMImage()
  }, [drawDICOMImage])

  const handleMouseDown = (event: React.MouseEvent<HTMLCanvasElement>) => {
    const activeCanvas = canvasReference.current
    if (!activeCanvas) {
      return
    }
    const canvasBounds = activeCanvas.getBoundingClientRect()
    const clickX = event.clientX - canvasBounds.left
    const clickY = event.clientY - canvasBounds.top

    isDraggingReference.current = true
    dragStartPoint.current = { x: event.clientX, y: event.clientY }

    if (activeTool === "ruler") {
      rulerStartPoint.current = { x: clickX, y: clickY }
      rulerEndPoint.current = { x: clickX, y: clickY }
      drawDICOMImage()
    }
  }

  const handleMouseMove = (event: React.MouseEvent<HTMLCanvasElement>) => {
    if (!isDraggingReference.current) {
      return
    }
    const activeCanvas = canvasReference.current
    if (!activeCanvas) {
      return
    }
    const canvasBounds = activeCanvas.getBoundingClientRect()
    const currentX = event.clientX - canvasBounds.left
    const currentY = event.clientY - canvasBounds.top

    const deltaX = event.clientX - dragStartPoint.current.x
    const deltaY = event.clientY - dragStartPoint.current.y

    if (activeTool === "zoom") {
      const increment = deltaY * 0.005
      setZoomSetting((previous) => Math.max(0.2, Math.min(5.0, previous - increment)))
      dragStartPoint.current = { x: event.clientX, y: event.clientY }
    } else if (activeTool === "contrast") {
      const contrastIncrement = deltaX * 0.005
      const brightnessIncrement = deltaY * 0.005
      setContrastSetting((previous) => Math.max(0.1, Math.min(3.0, previous + contrastIncrement)))
      setBrightnessSetting((previous) => Math.max(0.1, Math.min(3.0, previous - brightnessIncrement)))
      dragStartPoint.current = { x: event.clientX, y: event.clientY }
    } else if (activeTool === "ruler") {
      rulerEndPoint.current = { x: currentX, y: currentY }
      drawDICOMImage()
    }
  }

  const handleMouseUp = () => {
    isDraggingReference.current = false
  }

  const handlePreset = (presetType: "bone" | "lung" | "soft") => {
    rulerStartPoint.current = null
    rulerEndPoint.current = null
    if (presetType === "bone") {
      setContrastSetting(2.0)
      setBrightnessSetting(0.8)
      setZoomSetting(1.0)
    } else if (presetType === "lung") {
      setContrastSetting(0.6)
      setBrightnessSetting(1.4)
      setZoomSetting(1.0)
    } else {
      setContrastSetting(1.0)
      setBrightnessSetting(1.0)
      setZoomSetting(1.0)
    }
  }

  const simulateDICOMUpload = async () => {
    if (!study) {
      return
    }
    setUploadPercentage(0)
    setUploadStatus("Iniciando upload de chunks gRPC-Web...")

    for (let currentProgress = 10; currentProgress <= 100; currentProgress += 10) {
      await new Promise((resolve) => setTimeout(resolve, 300))
      setUploadPercentage(currentProgress)
      if (currentProgress < 100) {
        setUploadStatus(`Enviando chunk ${currentProgress / 10} de 10...`)
      } else {
        setUploadStatus("Transmissão gRPC concluída. Processando DICOM metadata...")
      }
    }

    await new Promise((resolve) => setTimeout(resolve, 500))
    setUploadPercentage(null)
    setUploadStatus(null)
    alert("DICOM carregado e processado com sucesso no barramento do PACS!")
  }

  if (isStudyLoading || !study) {
    return (
      <div className="text-center py-16">
        <span className="text-sm text-muted">Carregando visualizador PACS...</span>
      </div>
    )
  }

  return (
    <div className="flex-1 p-8 flex flex-col gap-6 max-w-7xl mx-auto w-full select-none">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div className="flex items-center gap-4">
          <Button variantType="outline" onClick={() => navigate(`/patients/${study.patient_fhir_id}`)} className="px-3">
            <ArrowLeft className="w-4 h-4" />
            Voltar Prontuário
          </Button>
          <div className="text-left">
            <h2 className="text-xl font-black text-gray-900 leading-none">
              Console Cirúrgico PACS
            </h2>
            <span className="text-xs text-muted mt-1.5 block">
              Estudo: {study.title} • UID: {study.study_instance_uid}
            </span>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button variantType="outline" onClick={simulateDICOMUpload} className="px-3.5 gap-2 border-primary/20 text-primary hover:bg-primary/5">
            <UploadCloud className="w-4 h-4" />
            Upload Novo .DCM
          </Button>
        </div>
      </div>

      {uploadPercentage !== null && (
        <Card className="p-4 bg-primary/5 border border-primary/20 flex flex-col gap-2.5 text-left">
          <div className="flex justify-between items-center text-xs">
            <span className="text-primary font-bold">{uploadStatus}</span>
            <span className="text-gray-500 font-bold">{uploadPercentage}%</span>
          </div>
          <div className="w-full bg-gray-100 rounded-full h-2">
            <div
              className="bg-primary h-2 rounded-full transition-all duration-300"
              style={{ width: `${uploadPercentage}%` }}
            />
          </div>
        </Card>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6 items-start">
        
        <Card className="flex flex-col gap-5 lg:col-span-1 text-left">
          <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2 border-b border-border pb-4">
            <Eye className="w-4 h-4 text-primary animate-pulse-glow" />
            Detalhes do Estudo
          </h3>

          <div className="flex flex-col gap-4">
            <div className="bg-gray-50 border border-border p-3.5 rounded-xl">
              <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">ID do Estudo</span>
              <span className="text-xs font-bold text-gray-800 mt-1 block">{study.imaging_study_id}</span>
            </div>

            <div className="bg-gray-50 border border-border p-3.5 rounded-xl">
              <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Modalidade Clínica</span>
              <span className="text-xs font-bold text-gray-800 mt-1 block uppercase">{study.modality}</span>
            </div>

            <div className="bg-gray-50 border border-border p-3.5 rounded-xl">
              <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Status do Barramento</span>
              <span className="text-xs font-bold text-emerald-600 mt-1 flex items-center gap-1.5">
                <CheckCircle className="w-4 h-4" />
                {study.status}
              </span>
            </div>

            <div className="bg-gray-50 border border-border p-3.5 rounded-xl">
              <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Série / Fatias</span>
              <span className="text-xs font-bold text-gray-800 mt-1 block">Slice #1 (Visualização Ativa)</span>
            </div>
          </div>
        </Card>

        <Card className="lg:col-span-3 flex flex-col items-center gap-6 p-4">
          <div className="relative border border-border rounded-2xl overflow-hidden bg-slate-950">
            <canvas
              ref={canvasReference}
              width={640}
              height={440}
              onMouseDown={handleMouseDown}
              onMouseMove={handleMouseMove}
              onMouseUp={handleMouseUp}
              onMouseLeave={handleMouseUp}
              className="cursor-crosshair block w-full max-w-full"
            />
          </div>

          <div className="flex flex-wrap items-center justify-between w-full border-t border-border pt-4 gap-4">
            <div className="flex gap-2.5">
              <Button
                variantType={activeTool === "zoom" ? "primary" : "outline"}
                onClick={() => setActiveTool("zoom")}
                className="px-3.5 gap-2 text-xs font-bold"
              >
                <ZoomIn className="w-4 h-4" />
                Zoom (Arrastar)
              </Button>
              <Button
                variantType={activeTool === "contrast" ? "primary" : "outline"}
                onClick={() => setActiveTool("contrast")}
                className="px-3.5 gap-2 text-xs font-bold"
              >
                <Sun className="w-4 h-4" />
                Luminosidade
              </Button>
              <Button
                variantType={activeTool === "ruler" ? "primary" : "outline"}
                onClick={() => setActiveTool("ruler")}
                className="px-3.5 gap-2 text-xs font-bold"
              >
                <Ruler className="w-4 h-4" />
                Régua (Medir)
              </Button>
            </div>

            <div className="flex gap-2">
              <Button variantType="outline" onClick={() => handlePreset("soft")} className="px-3 py-1.5 text-xs">
                Tecido Mole
              </Button>
              <Button variantType="outline" onClick={() => handlePreset("bone")} className="px-3 py-1.5 text-xs">
                Osso
              </Button>
              <Button variantType="outline" onClick={() => handlePreset("lung")} className="px-3 py-1.5 text-xs">
                Pulmão
              </Button>
            </div>
          </div>
        </Card>
      </div>
    </div>
  )
}
