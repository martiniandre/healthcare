import { useTranslation } from "react-i18next"
import { Ruler, Sun, ZoomIn, type LucideIcon } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import { Separator } from "@/shared/components/ui/separator"
import type { ImagingPreset, ImagingTool } from "../hooks/useDicomViewer"

interface DicomToolControlsProperties {
  activeTool: ImagingTool
  onToolChange: (tool: ImagingTool) => void
  onPresetChange: (preset: ImagingPreset) => void
}

const TOOL_OPTIONS: { tool: ImagingTool; labelKey: string; Icon: LucideIcon }[] = [
  { tool: "zoom", labelKey: "tools.zoom", Icon: ZoomIn },
  { tool: "contrast", labelKey: "tools.contrast", Icon: Sun },
  { tool: "ruler", labelKey: "tools.measure", Icon: Ruler },
]

const PRESET_OPTIONS: { preset: ImagingPreset; labelKey: string }[] = [
  { preset: "soft", labelKey: "presets.softTissue" },
  { preset: "bone", labelKey: "presets.bone" },
  { preset: "lung", labelKey: "presets.lung" },
]

export const DicomToolControls = ({ activeTool, onToolChange, onPresetChange }: DicomToolControlsProperties) => {
  const { t } = useTranslation("imaging")

  return (
    <div className="flex w-full flex-col gap-4">
      <Separator />
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="flex gap-2.5">
          {TOOL_OPTIONS.map(({ tool, labelKey, Icon }) => (
            <Button
              key={tool}
              variant={activeTool === tool ? "default" : "outline"}
              aria-pressed={activeTool === tool}
              onClick={() => onToolChange(tool)}
              className="text-xs font-bold"
            >
              <Icon />
              {t(labelKey)}
            </Button>
          ))}
        </div>

        <div className="flex gap-2">
          {PRESET_OPTIONS.map(({ preset, labelKey }) => (
            <Button
              key={preset}
              variant="outline"
              size="sm"
              onClick={() => onPresetChange(preset)}
              className="text-xs"
            >
              {t(labelKey)}
            </Button>
          ))}
        </div>
      </div>
    </div>
  )
}
