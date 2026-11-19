import { useEffect, useRef, type RefObject } from "react"
import { CardiacCondition } from "@/shared/types"
import { computeEcgSample, ECG_COLORS, ECG_GRID_SIZE, ECG_PLOT_POINTS } from "../utils/ecg_waveform"

const drawGrid = (canvasContext: CanvasRenderingContext2D, canvasWidth: number, canvasHeight: number) => {
  canvasContext.fillStyle = ECG_COLORS.background
  canvasContext.fillRect(0, 0, canvasWidth, canvasHeight)
  canvasContext.strokeStyle = ECG_COLORS.grid
  canvasContext.lineWidth = 1

  for (let xCoordinate = 0; xCoordinate < canvasWidth; xCoordinate += ECG_GRID_SIZE) {
    canvasContext.beginPath()
    canvasContext.moveTo(xCoordinate, 0)
    canvasContext.lineTo(xCoordinate, canvasHeight)
    canvasContext.stroke()
  }

  for (let yCoordinate = 0; yCoordinate < canvasHeight; yCoordinate += ECG_GRID_SIZE) {
    canvasContext.beginPath()
    canvasContext.moveTo(0, yCoordinate)
    canvasContext.lineTo(canvasWidth, yCoordinate)
    canvasContext.stroke()
  }
}

const drawTrace = (
  canvasContext: CanvasRenderingContext2D,
  tracePoints: number[],
  canvasWidth: number,
  condition: CardiacCondition | undefined
) => {
  canvasContext.beginPath()
  canvasContext.strokeStyle = condition === CardiacCondition.CardiacArrest ? ECG_COLORS.arrestTrace : ECG_COLORS.trace
  canvasContext.lineWidth = 2.5
  canvasContext.lineJoin = "round"

  tracePoints.forEach((yPixelCoordinate, pointIndex) => {
    const xPixelCoordinate = (pointIndex / (tracePoints.length - 1)) * canvasWidth
    if (pointIndex === 0) {
      canvasContext.moveTo(xPixelCoordinate, yPixelCoordinate)
    } else {
      canvasContext.lineTo(xPixelCoordinate, yPixelCoordinate)
    }
  })
  canvasContext.stroke()
}

export const useEcgTrace = (
  canvasRef: RefObject<HTMLCanvasElement | null>,
  isActive: boolean,
  condition: CardiacCondition | undefined
) => {
  const animationFrameIdRef = useRef<number | null>(null)
  const traceOffsetRef = useRef(0)
  const tracePointsRef = useRef<number[]>([])

  useEffect(() => {
    if (!isActive) {
      return
    }

    const renderFrame = () => {
      const activeCanvas = canvasRef.current
      const canvasContext = activeCanvas?.getContext("2d")
      if (!activeCanvas || !canvasContext) {
        return
      }

      const { width: canvasWidth, height: canvasHeight } = activeCanvas
      drawGrid(canvasContext, canvasWidth, canvasHeight)

      if (tracePointsRef.current.length < ECG_PLOT_POINTS) {
        tracePointsRef.current = new Array(ECG_PLOT_POINTS).fill(canvasHeight / 2)
      }
      tracePointsRef.current.push(computeEcgSample(traceOffsetRef.current, canvasHeight / 2, condition))
      tracePointsRef.current.shift()

      drawTrace(canvasContext, tracePointsRef.current, canvasWidth, condition)

      traceOffsetRef.current += 1
      animationFrameIdRef.current = requestAnimationFrame(renderFrame)
    }

    animationFrameIdRef.current = requestAnimationFrame(renderFrame)
    return () => {
      if (animationFrameIdRef.current) {
        cancelAnimationFrame(animationFrameIdRef.current)
      }
    }
  }, [canvasRef, isActive, condition])
}
