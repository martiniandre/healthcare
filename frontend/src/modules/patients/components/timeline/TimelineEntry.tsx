import { ExternalLink } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Button } from "../../../../shared/components/ui/button"
import { Card } from "../../../../shared/components/ui/card"
import { StatusBadge } from "../../../../shared/components/ui/status-badge"
import { useLocale } from "../../../../shared/hooks/useLocale"
import { formatDateTime, formatTime } from "../../../../shared/utils/dates"
import { TimelineEntryType, type TimelineItem } from "../../types"
import { resolveClinicalStatus } from "../clinicalStatus"
import { isSameMoment, statusKindByType, statusNamespaceByType } from "./timelineModel"
import { timelineTypeIcons } from "./timelineTypeIcons"

interface TimelineEntryProps {
  item: TimelineItem
  onOpenEncounter: (encounterId: string) => void
}

export const TimelineEntry = ({ item, onOpenEncounter }: TimelineEntryProps) => {
  const { t } = useTranslation("patients")
  const { t: translateStatus } = useTranslation("clinicalStatus")
  const locale = useLocale()
  const TypeIcon = timelineTypeIcons[item.type]
  const statusKind = statusKindByType[item.type]
  const relatedEncounterId = item.encounterId ?? (item.type === TimelineEntryType.Encounter ? item.resourceId : null)
  const showClinicalDate = item.clinicalDate !== null && !isSameMoment(item.clinicalDate, item.occurredAt)

  return (
    <Card className="gap-2 p-4" data-testid="timeline-entry" data-entry-type={item.type}>
      <div className="flex flex-wrap items-center gap-x-2 gap-y-1">
        <span className="inline-flex items-center gap-1.5 text-[11px] font-bold uppercase tracking-wider text-primary">
          <TypeIcon className="size-4" aria-hidden="true" />
          {t(`timeline.types.${item.type}`)}
        </span>
        <time dateTime={item.occurredAt} className="text-xs font-semibold tabular-nums text-muted-foreground">
          {formatTime(item.occurredAt, locale)}
        </time>
        {item.status && (
          <StatusBadge
            status={statusKind ? resolveClinicalStatus(statusKind, item.status) : "neutral"}
            className="ml-auto uppercase tracking-wider"
          >
            {translateStatus(`${statusNamespaceByType[item.type]}.${item.status}`, { defaultValue: item.status })}
          </StatusBadge>
        )}
      </div>
      <h4 className="text-sm font-bold text-foreground">{item.title}</h4>
      {item.summary && <p className="text-sm text-muted-foreground">{item.summary}</p>}
      {showClinicalDate && item.clinicalDate && (
        <p className="text-xs tabular-nums text-muted-foreground">
          {t("timeline.clinicalDate", { date: formatDateTime(item.clinicalDate, locale) })}
        </p>
      )}
      {relatedEncounterId && (
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="mt-1 min-h-11 self-start md:min-h-8"
          onClick={() => onOpenEncounter(relatedEncounterId)}
        >
          <ExternalLink aria-hidden="true" />
          {t("timeline.openEncounter")}
        </Button>
      )}
    </Card>
  )
}
