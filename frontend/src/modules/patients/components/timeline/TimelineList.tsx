import { useTranslation } from "react-i18next"
import { useLocale } from "../../../../shared/hooks/useLocale"
import { formatLongDate } from "../../../../shared/utils/dates"
import type { TimelineItem } from "../../types"
import { TimelineEntry } from "./TimelineEntry"
import { groupTimelineByDay } from "./timelineModel"

interface TimelineListProps {
  items: TimelineItem[]
  onOpenEncounter: (encounterId: string) => void
}

export const TimelineList = ({ items, onOpenEncounter }: TimelineListProps) => {
  const { t } = useTranslation("patients")
  const locale = useLocale()
  const dayGroups = groupTimelineByDay(items)

  return (
    <div className="flex flex-col gap-4" aria-label={t("timeline.title")} role="group">
      {dayGroups.map((dayGroup) => (
        <section key={dayGroup.dayKey} aria-labelledby={`timeline-day-${dayGroup.dayKey}`} data-testid="timeline-day">
          <h3
            id={`timeline-day-${dayGroup.dayKey}`}
            className="sticky top-0 z-10 -mx-1 bg-background/95 px-1 py-2 text-xs font-black uppercase tracking-widest tabular-nums text-muted-foreground backdrop-blur"
          >
            {formatLongDate(dayGroup.dayKey, locale)}
          </h3>
          <ol className="flex flex-col gap-2">
            {dayGroup.items.map((timelineItem) => (
              <li key={timelineItem.id}>
                <TimelineEntry item={timelineItem} onOpenEncounter={onOpenEncounter} />
              </li>
            ))}
          </ol>
        </section>
      ))}
    </div>
  )
}
