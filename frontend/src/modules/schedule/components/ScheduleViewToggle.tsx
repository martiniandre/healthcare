import { useTranslation } from "react-i18next"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/shared/components/ui/select"
import type { ScheduleViewMode } from "./ScheduleCalendar"

interface ScheduleViewToggleProps {
  value: ScheduleViewMode
  onChange: (mode: ScheduleViewMode) => void
}

const VIEW_MODES: ScheduleViewMode[] = ["week", "day", "month", "year"]

export const ScheduleViewToggle = ({ value, onChange }: ScheduleViewToggleProps) => {
  const { t } = useTranslation("schedule")

  return (
    <Select value={value} onValueChange={(mode) => onChange(mode as ScheduleViewMode)}>
      <SelectTrigger aria-label={t("calendar.viewLabel")} className="w-40 font-medium">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        {VIEW_MODES.map((mode) => (
          <SelectItem key={mode} value={mode}>
            {t(`calendar.views.${mode}`)}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
