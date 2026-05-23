import * as React from "react"
import { useCallback, useEffect, useRef, useState } from "react"
import type { ImagingStudy } from "../types"

export type ImagingTool = "zoom" | "contrast" | "ruler"
export type ImagingPreset = "bone" | "lung" | "soft"

interface DicomViewerConfiguration {
  contrast: number
  brightness: number
  zoom: number
}

export const useDicomViewer = (imagingStudy: ImagingStudy | null | undefined) => {
  const [activeTool, setActiveTool] = useState<ImagingTool>("zoom")
  
  const [interfaceConfiguration, setInterfaceConfiguration] = useState<DicomViewerConfiguration>({
    contrast: 1.0,
    brightness: 1.0,
    zoom: 1.0,
  })

  const canvasReference = useRef<HTMLCanvasElement | null>(null)
  const isUserDraggingReference = useRef(false)
  const dragStartPointCoordinate = useRef({ coordinateX: 0, coordinateY: 0 })
  const rulerStartPointCoordinate = useRef<{ coordinateX: number; coordinateY: number } | null>(null)
  const rulerEndPointCoordinate = useRef<{ coordinateX: number; coordinateY: number } | null>(null)

  const contrastValueReference = useRef(1.0)
  const brightnessValueReference = useRef(1.0)
  const zoomValueReference = useRef(1.0)

  const drawDicomImageOnCanvas = useCallback(() => {
    const activeCanvasElement = canvasReference.current
    if (!activeCanvasElement) {
      return
    }

    const canvasRenderingContext = activeCanvasElement.getContext("2d")
    if (!canvasRenderingContext) {
      return
    }

    const canvasWidthPixels = activeCanvasElement.width
    const canvasHeightPixels = activeCanvasElement.height

    canvasRenderingContext.fillStyle = "#090d16"
    canvasRenderingContext.fillRect(0, 0, canvasWidthPixels, canvasHeightPixels)

    canvasRenderingContext.save()
    canvasRenderingContext.translate(canvasWidthPixels / 2, canvasHeightPixels / 2)
    canvasRenderingContext.scale(zoomValueReference.current, zoomValueReference.current)
    canvasRenderingContext.translate(-canvasWidthPixels / 2, -canvasHeightPixels / 2)
    
    const contrastPercentageValue = contrastValueReference.current * 100
    const brightnessPercentageValue = brightnessValueReference.current * 100
    canvasRenderingContext.filter = `contrast(${contrastPercentageValue}%) brightness(${brightnessPercentageValue}%)`

    canvasRenderingContext.fillStyle = "#1e293b"
    canvasRenderingContext.beginPath()
    canvasRenderingContext.arc(canvasWidthPixels / 2, canvasHeightPixels / 2, 140, 0, Math.PI * 2)
    canvasRenderingContext.fill()

    canvasRenderingContext.fillStyle = "#334155"
    canvasRenderingContext.fillRect(canvasWidthPixels / 2 - 30, canvasHeightPixels / 2 - 120, 60, 240)
    canvasRenderingContext.fillRect(canvasWidthPixels / 2 - 120, canvasHeightPixels / 2 - 30, 240, 60)

    canvasRenderingContext.strokeStyle = "#475569"
    canvasRenderingContext.lineWidth = 4
    canvasRenderingContext.strokeRect(canvasWidthPixels / 2 - 80, canvasHeightPixels / 2 - 80, 160, 160)

    canvasRenderingContext.fillStyle = "#64748b"
    canvasRenderingContext.beginPath()
    canvasRenderingContext.arc(canvasWidthPixels / 2, canvasHeightPixels / 2, 45, 0, Math.PI * 2)
    canvasRenderingContext.fill()

    canvasRenderingContext.restore()

    canvasRenderingContext.fillStyle = "rgba(255, 255, 255, 0.4)"
    canvasRenderingContext.font = "10px Outfit, sans-serif"
    canvasRenderingContext.fillText("HOSPITAL GERAL - MOCK PACS", 15, 25)
    canvasRenderingContext.fillText(`SERIES: 1 • MODALITY: ${imagingStudy?.modality || "CT"}`, 15, 40)
    canvasRenderingContext.fillText(`UID: ${imagingStudy?.study_instance_uid || "1.2.840.10008"}`, 15, 55)
    
    canvasRenderingContext.fillText(`ZOOM: ${Math.round(zoomValueReference.current * 100)}%`, canvasWidthPixels - 110, 25)
    canvasRenderingContext.fillText(`CONTRASTE: ${Math.round(contrastValueReference.current * 100)}%`, canvasWidthPixels - 110, 40)
    canvasRenderingContext.fillText(`BRILHO: ${Math.round(brightnessValueReference.current * 100)}%`, canvasWidthPixels - 110, 55)

    if (rulerStartPointCoordinate.current) {
      canvasRenderingContext.strokeStyle = "#0ea5e9"
      canvasRenderingContext.lineWidth = 2
      canvasRenderingContext.beginPath()
      canvasRenderingContext.moveTo(rulerStartPointCoordinate.current.coordinateX, rulerStartPointCoordinate.current.coordinateY)

      const activeEndPointCoordinate = rulerEndPointCoordinate.current || rulerStartPointCoordinate.current
      canvasRenderingContext.lineTo(activeEndPointCoordinate.coordinateX, activeEndPointCoordinate.coordinateY)
      canvasRenderingContext.stroke()

      canvasRenderingContext.fillStyle = "#0ea5e9"
      canvasRenderingContext.beginPath()
      canvasRenderingContext.arc(rulerStartPointCoordinate.current.coordinateX, rulerStartPointCoordinate.current.coordinateY, 4, 0, Math.PI * 2)
      canvasRenderingContext.arc(activeEndPointCoordinate.coordinateX, activeEndPointCoordinate.coordinateY, 4, 0, Math.PI * 2)
      canvasRenderingContext.fill()

      const deltaCoordinateX = activeEndPointCoordinate.coordinateX - rulerStartPointCoordinate.current.coordinateX
      const deltaCoordinateY = activeEndPointCoordinate.coordinateY - rulerStartPointCoordinate.current.coordinateY
      const calculatedDistancePixels = Math.sqrt(deltaCoordinateX * deltaCoordinateX + deltaCoordinateY * deltaCoordinateY)
      const calculatedDistanceMillimeters = (calculatedDistancePixels * 0.28).toFixed(1)

      canvasRenderingContext.fillStyle = "#ffffff"
      canvasRenderingContext.font = "12px Outfit, sans-serif"
      canvasRenderingContext.fillText(`${calculatedDistanceMillimeters} mm`, activeEndPointCoordinate.coordinateX + 10, activeEndPointCoordinate.coordinateY + 10)
    }
  }, [imagingStudy])

  useEffect(() => {
    drawDicomImageOnCanvas()
  }, [drawDicomImageOnCanvas])

  const handleMouseDown = (mouseEvent: React.MouseEvent<HTMLCanvasElement>) => {
    const activeCanvasElement = canvasReference.current
    if (!activeCanvasElement) {
      return
    }

    const canvasElementBounds = activeCanvasElement.getBoundingClientRect()
    const relativeCoordinateX = mouseEvent.clientX - canvasElementBounds.left
    const relativeCoordinateY = mouseEvent.clientY - canvasElementBounds.top

    isUserDraggingReference.current = true
    dragStartPointCoordinate.current = { coordinateX: mouseEvent.clientX, coordinateY: mouseEvent.clientY }

    if (activeTool === "ruler") {
      rulerStartPointCoordinate.current = { coordinateX: relativeCoordinateX, coordinateY: relativeCoordinateY }
      rulerEndPointCoordinate.current = { coordinateX: relativeCoordinateX, coordinateY: relativeCoordinateY }
      drawDicomImageOnCanvas()
    }
  }

  const handleMouseMove = (mouseEvent: React.MouseEvent<HTMLCanvasElement>) => {
    if (!isUserDraggingReference.current) {
      return
    }

    const activeCanvasElement = canvasReference.current
    if (!activeCanvasElement) {
      return
    }

    const canvasElementBounds = activeCanvasElement.getBoundingClientRect()
    const currentCoordinateX = mouseEvent.clientX - canvasElementBounds.left
    const currentCoordinateY = mouseEvent.clientY - canvasElementBounds.top
    
    const deltaCoordinateX = mouseEvent.clientX - dragStartPointCoordinate.current.coordinateX
    const deltaCoordinateY = mouseEvent.clientY - dragStartPointCoordinate.current.coordinateY

    if (activeTool === "zoom") {
      const zoomScaleIncrement = deltaCoordinateY * 0.005
      zoomValueReference.current = Math.max(0.2, Math.min(5.0, zoomValueReference.current - zoomScaleIncrement))
      dragStartPointCoordinate.current = { coordinateX: mouseEvent.clientX, coordinateY: mouseEvent.clientY }
      drawDicomImageOnCanvas()
      return
    }

    if (activeTool === "contrast") {
      const contrastValueIncrement = deltaCoordinateX * 0.005
      const brightnessValueIncrement = deltaCoordinateY * 0.005
      contrastValueReference.current = Math.max(0.1, Math.min(3.0, contrastValueReference.current + contrastValueIncrement))
      brightnessValueReference.current = Math.max(0.1, Math.min(3.0, brightnessValueReference.current - brightnessValueIncrement))
      dragStartPointCoordinate.current = { coordinateX: mouseEvent.clientX, coordinateY: mouseEvent.clientY }
      drawDicomImageOnCanvas()
      return
    }

    rulerEndPointCoordinate.current = { coordinateX: currentCoordinateX, coordinateY: currentCoordinateY }
    drawDicomImageOnCanvas()
  }

  const handleMouseUp = () => {
    isUserDraggingReference.current = false
    setInterfaceConfiguration({
      contrast: contrastValueReference.current,
      brightness: brightnessValueReference.current,
      zoom: zoomValueReference.current,
    })
  }

  const applyPreset = (imagingPreset: ImagingPreset) => {
    rulerStartPointCoordinate.current = null
    rulerEndPointCoordinate.current = null

    if (imagingPreset === "bone") {
      contrastValueReference.current = 2.0
      brightnessValueReference.current = 0.8
      zoomValueReference.current = 1.0
    } else if (imagingPreset === "lung") {
      contrastValueReference.current = 0.6
      brightnessValueReference.current = 1.4
      zoomValueReference.current = 1.0
    } else {
      contrastValueReference.current = 1.0
      brightnessValueReference.current = 1.0
      zoomValueReference.current = 1.0
    }

    setInterfaceConfiguration({
      contrast: contrastValueReference.current,
      brightness: brightnessValueReference.current,
      zoom: zoomValueReference.current,
    })

    drawDicomImageOnCanvas()
  }

  return {
    activeTool,
    canvasReference,
    interfaceConfiguration,
    setActiveTool,
    applyPreset,
    handleMouseDown,
    handleMouseMove,
    handleMouseUp,
  }
}
