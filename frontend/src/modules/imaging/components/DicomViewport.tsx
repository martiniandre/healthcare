import * as React from "react"
import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { DicomToolControls } from "./DicomToolControls"
import type { ImagingPreset, ImagingTool } from "../hooks/useDicomViewer"

interface DicomViewportProperties {
  activeTool: ImagingTool
  canvasReference: React.RefObject<HTMLCanvasElement | null>
  onMouseDown: (mouseEvent: React.MouseEvent<HTMLCanvasElement>) => void
  onMouseMove: (mouseEvent: React.MouseEvent<HTMLCanvasElement>) => void
  onMouseUp: () => void
  onToolChange: (imagingTool: ImagingTool) => void
  onPresetChange: (imagingPreset: ImagingPreset) => void
}

export const DicomViewport = ({
  activeTool,
  canvasReference,
  onMouseDown,
  onMouseMove,
  onMouseUp,
  onToolChange,
  onPresetChange,
}: DicomViewportProperties) => {
  const { t } = useTranslation("imaging")

  const getContextualCursorClass = (): string => {
    if (activeTool === "zoom") {
      return "cursor-zoom-in"
    }
    if (activeTool === "contrast") {
      return "cursor-ew-resize"
    }
    return "cursor-crosshair"
  }

  return (
    <Card className="lg:col-span-3 flex flex-col items-center gap-6 p-4">
      <div className="relative border border-border rounded-2xl overflow-hidden bg-slate-950 w-full flex justify-center">
        <canvas
          ref={canvasReference}
          width={640}
          height={440}
          onMouseDown={onMouseDown}
          onMouseMove={onMouseMove}
          onMouseUp={onMouseUp}
          onMouseLeave={onMouseUp}
          role="img"
          aria-label={t("details.viewportAriaLabel")}
          className={`block w-full max-w-full ${getContextualCursorClass()}`}
        />
      </div>

      <DicomToolControls
        activeTool={activeTool}
        onToolChange={onToolChange}
        onPresetChange={onPresetChange}
      />
    </Card>
  )
}
