import { useTranslation } from "react-i18next"
import type { ScheduleViewMode } from "./ScheduleCalendar"

interface ScheduleViewToggleProps {
  value: ScheduleViewMode
  onChange: (mode: ScheduleViewMode) => void
}

const VIEW_MODES: ScheduleViewMode[] = ["week", "day", "month", "year"]

export const ScheduleViewToggle = ({ value, onChange }: ScheduleViewToggleProps) => {
  const { t } = useTranslation("schedule")

  return (
    <select
      aria-label={t("calendar.viewLabel")}
      value={value}
      onChange={(event) => onChange(event.target.value as ScheduleViewMode)}
      className="h-10 w-40 cursor-pointer rounded-md border border-input bg-background px-3 py-2 text-sm font-medium text-gray-700 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
    >
      {VIEW_MODES.map((mode) => (
        <option key={mode} value={mode}>
          {t(`calendar.views.${mode}`)}
        </option>
      ))}
    </select>
  )
}
