import { Ruler, Sun, ZoomIn } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"
import type { ImagingPreset, ImagingTool } from "../hooks/useDicomViewer"

interface DicomToolControlsProperties {
  activeTool: ImagingTool
  onToolChange: (tool: ImagingTool) => void
  onPresetChange: (preset: ImagingPreset) => void
}

export const DicomToolControls = ({ activeTool, onToolChange, onPresetChange }: DicomToolControlsProperties) => {
  return (
    <div className="flex flex-wrap items-center justify-between w-full border-t border-border pt-4 gap-4">
      <div className="flex gap-2.5">
        <Button
          variantType={activeTool === "zoom" ? "primary" : "outline"}
          onClick={() => onToolChange("zoom")}
          className="px-3.5 gap-2 text-xs font-bold"
        >
          <ZoomIn className="w-4 h-4" />
          Zoom (Arrastar)
        </Button>
        <Button
          variantType={activeTool === "contrast" ? "primary" : "outline"}
          onClick={() => onToolChange("contrast")}
          className="px-3.5 gap-2 text-xs font-bold"
        >
          <Sun className="w-4 h-4" />
          Luminosidade
        </Button>
        <Button
          variantType={activeTool === "ruler" ? "primary" : "outline"}
          onClick={() => onToolChange("ruler")}
          className="px-3.5 gap-2 text-xs font-bold"
        >
          <Ruler className="w-4 h-4" />
          Régua (Medir)
        </Button>
      </div>

      <div className="flex gap-2">
        <Button variantType="outline" onClick={() => onPresetChange("soft")} className="px-3 py-1.5 text-xs">
          Tecido Mole
        </Button>
        <Button variantType="outline" onClick={() => onPresetChange("bone")} className="px-3 py-1.5 text-xs">
          Osso
        </Button>
        <Button variantType="outline" onClick={() => onPresetChange("lung")} className="px-3 py-1.5 text-xs">
          Pulmão
        </Button>
      </div>
    </div>
  )
}
