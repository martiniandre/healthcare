import { useMemo } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { AlertTriangle, History, RotateCw } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Alert, AlertDescription, AlertTitle } from "../../../../shared/components/ui/alert"
import { Button } from "../../../../shared/components/ui/button"
import { Card } from "../../../../shared/components/ui/card"
import { EmptyState } from "../../../../shared/components/ui/EmptyState"
import { Skeleton } from "../../../../shared/components/ui/skeleton"
import { Spinner } from "../../../../shared/components/ui/spinner"
import { useLocale } from "../../../../shared/hooks/useLocale"
import { useTimelineQuery } from "../../queries"
import { QueryErrorAlert } from "../QueryErrorAlert"
import { TimelineFilters } from "./TimelineFilters"
import { TimelineList } from "./TimelineList"
import { parseTimelineFilters, serializeTimelineFilters, type TimelineFilterState } from "./timelineModel"

interface ClinicalTimelineProps {
  patientId: string
}

const TimelineSkeleton = () => {
  const { t } = useTranslation("patients")

  return (
    <div role="status" aria-busy="true" aria-label={t("timeline.loadingList")} className="flex flex-col gap-2">
      {[0, 1, 2, 3].map((placeholderIndex) => (
        <Skeleton key={placeholderIndex} className="h-24 w-full rounded-xl" />
      ))}
    </div>
  )
}

export default function ClinicalTimeline({ patientId }: ClinicalTimelineProps) {
  const { t } = useTranslation("patients")
  const locale = useLocale()
  const navigate = useNavigate()
  const [searchParameters, setSearchParameters] = useSearchParams()
  const filters = useMemo(() => parseTimelineFilters(searchParameters), [searchParameters])

  const { data, isLoading, isError, refetch, isRefetching, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useTimelineQuery(patientId, filters)

  const updateFilters = (nextFilters: TimelineFilterState) => {
    setSearchParameters(
      { tab: searchParameters.get("tab") ?? "timeline", ...serializeTimelineFilters(nextFilters) },
      { replace: true },
    )
  }

  const items = useMemo(() => data?.pages.flatMap((page) => page.items) ?? [], [data])
  const failedTypes = useMemo(
    () => Array.from(new Set(data?.pages.flatMap((page) => page.partial?.failedTypes ?? []) ?? [])),
    [data],
  )
  const failedTypesLabel = new Intl.ListFormat(locale, { style: "long", type: "conjunction" }).format(
    failedTypes.map((failedType) => t(`timeline.types.${failedType}`)),
  )

  const renderBody = () => {
    if (isLoading) {
      return <TimelineSkeleton />
    }
    if (isError && !data) {
      return <QueryErrorAlert onRetry={() => void refetch()} isRetrying={isRefetching} />
    }
    if (items.length === 0 && failedTypes.length === 0) {
      return (
        <EmptyState icon={History} title={t("timeline.empty.title")} description={t("timeline.empty.description")} />
      )
    }
    return (
      <>
        {items.length > 0 && (
          <TimelineList
            items={items}
            onOpenEncounter={(encounterId) => navigate(`/patients/${patientId}/encounters/${encounterId}`)}
          />
        )}
        {isError && <QueryErrorAlert onRetry={() => void fetchNextPage()} isRetrying={isFetchingNextPage} />}
        {hasNextPage && (
          <Button
            type="button"
            variant="outline"
            className="min-h-11 self-center md:min-h-9"
            disabled={isFetchingNextPage}
            onClick={() => void fetchNextPage()}
          >
            {isFetchingNextPage && <Spinner />}
            {isFetchingNextPage ? t("timeline.loadingMore") : t("timeline.loadMore")}
          </Button>
        )}
      </>
    )
  }

  return (
    <Card className="gap-4 p-4 md:p-6">
      <div className="flex items-center gap-2">
        <History className="size-4 text-primary" aria-hidden="true" />
        <h2 className="text-sm font-extrabold uppercase tracking-wider text-foreground">{t("timeline.title")}</h2>
      </div>

      <TimelineFilters filters={filters} onChange={updateFilters} />

      {failedTypes.length > 0 && (
        <Alert variant="warning">
          <AlertTriangle />
          <AlertTitle>{t("timeline.partial.title")}</AlertTitle>
          <AlertDescription>
            <p>{t("timeline.partial.description", { types: failedTypesLabel })}</p>
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="mt-2 min-h-11 md:min-h-8"
              disabled={isRefetching}
              onClick={() => void refetch()}
            >
              <RotateCw />
              {t("timeline.partial.retry")}
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {renderBody()}
    </Card>
  )
}
