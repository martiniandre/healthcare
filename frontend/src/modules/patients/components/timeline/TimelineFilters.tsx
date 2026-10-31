import { Check } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Input } from "../../../../shared/components/ui/input"
import { Label } from "../../../../shared/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../../../../shared/components/ui/select"
import { cn } from "../../../../shared/utils/cn"
import type { TimelineEntryType } from "../../types"
import { TIMELINE_ENTRY_TYPES, TimelinePeriod, type TimelineFilterState } from "./timelineModel"
import { timelineTypeIcons } from "./timelineTypeIcons"

interface TimelineFiltersProps {
  filters: TimelineFilterState
  onChange: (filters: TimelineFilterState) => void
}

const chipClassName =
  "inline-flex min-h-11 shrink-0 items-center gap-1.5 rounded-full border px-3 text-xs font-bold transition-colors outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 md:min-h-8"

const selectedChipClassName = "border-primary/40 bg-primary/10 text-primary"
const idleChipClassName = "border-border bg-card text-muted-foreground hover:bg-muted"

const periodOptions = Object.values(TimelinePeriod)

export const TimelineFilters = ({ filters, onChange }: TimelineFiltersProps) => {
  const { t } = useTranslation("patients")

  const toggleType = (entryType: TimelineEntryType) => {
    const nextTypes = filters.types.includes(entryType)
      ? filters.types.filter((selectedType) => selectedType !== entryType)
      : TIMELINE_ENTRY_TYPES.filter(
          (candidateType) => candidateType === entryType || filters.types.includes(candidateType),
        )
    onChange({ ...filters, types: nextTypes })
  }

  const allTypesSelected = filters.types.length === 0

  return (
    <div className="flex flex-col gap-3">
      <div
        role="group"
        aria-label={t("timeline.filters.typesLabel")}
        className="-mx-4 flex gap-2 overflow-x-auto px-4 pb-1 md:mx-0 md:flex-wrap md:overflow-visible md:px-0"
      >
        <button
          type="button"
          aria-pressed={allTypesSelected}
          className={cn(chipClassName, allTypesSelected ? selectedChipClassName : idleChipClassName)}
          onClick={() => onChange({ ...filters, types: [] })}
        >
          {allTypesSelected && <Check className="size-3.5" aria-hidden="true" />}
          {t("timeline.filters.allTypes")}
        </button>
        {TIMELINE_ENTRY_TYPES.map((entryType) => {
          const TypeIcon = timelineTypeIcons[entryType]
          const isSelected = filters.types.includes(entryType)
          return (
            <button
              key={entryType}
              type="button"
              aria-pressed={isSelected}
              className={cn(chipClassName, isSelected ? selectedChipClassName : idleChipClassName)}
              onClick={() => toggleType(entryType)}
            >
              {isSelected ? (
                <Check className="size-3.5" aria-hidden="true" />
              ) : (
                <TypeIcon className="size-3.5" aria-hidden="true" />
              )}
              {t(`timeline.types.${entryType}`)}
            </button>
          )
        })}
      </div>

      <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-end">
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="timeline-period">{t("timeline.filters.periodLabel")}</Label>
          <Select
            value={filters.period}
            onValueChange={(period) => onChange({ ...filters, period: period as TimelinePeriod })}
          >
            <SelectTrigger id="timeline-period" className="min-h-11 w-full sm:w-56 md:min-h-9">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {periodOptions.map((period) => (
                <SelectItem key={period} value={period} className="min-h-11 md:min-h-8">
                  {t(`timeline.filters.periods.${period}`)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {filters.period === TimelinePeriod.Custom && (
          <>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="timeline-from">{t("timeline.filters.from")}</Label>
              <Input
                id="timeline-from"
                type="date"
                value={filters.from}
                max={filters.to || undefined}
                className="min-h-11 tabular-nums sm:w-44 md:min-h-9"
                onChange={(event) => onChange({ ...filters, from: event.target.value })}
              />
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="timeline-to">{t("timeline.filters.to")}</Label>
              <Input
                id="timeline-to"
                type="date"
                value={filters.to}
                min={filters.from || undefined}
                className="min-h-11 tabular-nums sm:w-44 md:min-h-9"
                onChange={(event) => onChange({ ...filters, to: event.target.value })}
              />
            </div>
          </>
        )}
      </div>
    </div>
  )
}
