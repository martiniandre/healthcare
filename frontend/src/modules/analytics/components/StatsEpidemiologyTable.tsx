import { useMemo, useState } from "react"
import type { ReactNode } from "react"
import { useTranslation } from "react-i18next"
import {
  ArrowDown,
  ArrowDownRight,
  ArrowUp,
  ArrowUpRight,
  FileSpreadsheet,
  Minus,
} from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../../shared/components/ui/Table"
import { AnalyticsPanel, PanelMetric } from "./AnalyticsPanel"
import { resolveTrendReading } from "../stats_domain"
import { toneClassNames, toneFillVariables, type MetricTone } from "../metric_tones"
import { cn } from "../../../shared/utils/cn"
import type { PathologyData } from "../types"

interface StatsEpidemiologyTableProps {
  pathologies: PathologyData[]
}

type SortField = "code" | "activeCases" | "trendSeverity"
type SortDirection = "asc" | "desc"

const TREND_SEVERITY_RANK: Record<MetricTone, number> = {
  critical: 3,
  elevated: 2,
  neutral: 1,
  stable: 0,
  primary: 0,
  secondary: 0,
}

const TREND_ICONS = {
  rising: ArrowUpRight,
  falling: ArrowDownRight,
  flat: Minus,
  unknown: Minus,
} as const

interface EnrichedPathology extends PathologyData {
  proportion: number
  trendSeverity: MetricTone
  trendDirection: "rising" | "falling" | "flat" | "unknown"
}

type AriaSortState = "none" | "ascending" | "descending"

interface SortableHeadProps {
  field: SortField
  sortField: SortField
  sortDirection: SortDirection
  align: "left" | "right"
  onToggleSort: (field: SortField) => void
  children: ReactNode
  className?: string
}

const SortableHead = ({
  field,
  sortField,
  sortDirection,
  align,
  onToggleSort,
  children,
  className,
}: SortableHeadProps) => {
  const isActive = sortField === field
  const SortIcon = sortDirection === "asc" ? ArrowUp : ArrowDown
  const ariaSortState: AriaSortState = isActive
    ? sortDirection === "asc"
      ? "ascending"
      : "descending"
    : "none"

  return (
    <TableHead
      aria-sort={ariaSortState}
      className={cn(align === "right" && "text-right", className)}
    >
      <button
        type="button"
        onClick={() => onToggleSort(field)}
        className={cn(
          "inline-flex items-center gap-1 rounded transition-colors hover:text-gray-900",
          align === "right" && "flex-row-reverse",
          isActive ? "text-gray-900" : "text-muted-foreground"
        )}
      >
        {children}
        {isActive ? <SortIcon className="size-3 shrink-0" aria-hidden="true" /> : null}
      </button>
    </TableHead>
  )
}

export const StatsEpidemiologyTable = ({ pathologies }: StatsEpidemiologyTableProps) => {
  const { t } = useTranslation("analytics")
  const [sortField, setSortField] = useState<SortField>("activeCases")
  const [sortDirection, setSortDirection] = useState<SortDirection>("desc")

  const totalActiveCases = useMemo(
    () => pathologies.reduce((runningTotal, pathology) => runningTotal + pathology.activeCases, 0),
    [pathologies]
  )

  const stableTrendTokens = useMemo(() => {
    return new Set(["stable", t("epidemiology.table.stable").toLowerCase()])
  }, [t])

  const enrichedPathologies = useMemo<EnrichedPathology[]>(() => {
    const peakActiveCases =
      pathologies.length > 0
        ? Math.max(...pathologies.map((pathology) => pathology.activeCases))
        : 0

    return pathologies.map((pathology) => {
      const trendReading = resolveTrendReading(pathology.trend)
      return {
        ...pathology,
        proportion: peakActiveCases > 0 ? (pathology.activeCases / peakActiveCases) * 100 : 0,
        trendSeverity: trendReading.severity,
        trendDirection: trendReading.direction,
      }
    })
  }, [pathologies])

  const sortedPathologies = useMemo<EnrichedPathology[]>(() => {
    const sortMultiplier = sortDirection === "asc" ? 1 : -1

    return [...enrichedPathologies].sort((firstPathology, secondPathology) => {
      if (sortField === "code") {
        return firstPathology.code.localeCompare(secondPathology.code) * sortMultiplier
      }
      if (sortField === "trendSeverity") {
        const severityDelta =
          TREND_SEVERITY_RANK[firstPathology.trendSeverity] -
          TREND_SEVERITY_RANK[secondPathology.trendSeverity]
        if (severityDelta !== 0) return severityDelta * sortMultiplier
        return (firstPathology.activeCases - secondPathology.activeCases) * sortMultiplier
      }
      return (firstPathology.activeCases - secondPathology.activeCases) * sortMultiplier
    })
  }, [enrichedPathologies, sortField, sortDirection])

  const resolveTrendLabel = (trend: string): string => {
    if (stableTrendTokens.has(trend.trim().toLowerCase())) {
      return t("epidemiology.table.stable")
    }
    return trend
  }

  const toggleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection((currentDirection) => (currentDirection === "asc" ? "desc" : "asc"))
      return
    }
    setSortField(field)
    setSortDirection(field === "code" ? "asc" : "desc")
  }

  return (
    <AnalyticsPanel
      title={t("epidemiology.title")}
      subtitle={t("epidemiology.subtitle")}
      trailing={
        <div className="flex flex-wrap items-center justify-end gap-2">
          {pathologies.length > 0 ? (
            <PanelMetric
              label={t("epidemiology.activeCasesTotalLabel")}
              value={t("epidemiology.activeCasesUnit", { count: totalActiveCases })}
            />
          ) : null}
          <Button variantType="outline" className="gap-1.5 text-xs">
            <FileSpreadsheet className="size-4" aria-hidden="true" />
            {t("epidemiology.exportButton")}
          </Button>
        </div>
      }
    >
      {pathologies.length > 0 ? (
        <div className="-mx-5 -my-5 overflow-x-auto">
          <Table className="min-w-[720px]">
            <TableHeader>
              <TableRow>
                <SortableHead
                  field="code"
                  sortField={sortField}
                  sortDirection={sortDirection}
                  align="left"
                  onToggleSort={toggleSort}
                  className="w-28 pl-5"
                >
                  {t("epidemiology.table.code")}
                </SortableHead>
                <TableHead>{t("epidemiology.table.description")}</TableHead>
                <TableHead className="w-44">{t("epidemiology.table.category")}</TableHead>
                <SortableHead
                  field="activeCases"
                  sortField={sortField}
                  sortDirection={sortDirection}
                  align="right"
                  onToggleSort={toggleSort}
                  className="w-28"
                >
                  {t("epidemiology.table.activeCases")}
                </SortableHead>
                <SortableHead
                  field="trendSeverity"
                  sortField={sortField}
                  sortDirection={sortDirection}
                  align="left"
                  onToggleSort={toggleSort}
                  className="w-40"
                >
                  {t("epidemiology.table.trend")}
                </SortableHead>
                <TableHead className="w-44 pr-5 text-right">
                  {t("epidemiology.table.proportion")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {sortedPathologies.map((pathology) => {
                const trendTone = toneClassNames[pathology.trendSeverity]
                const TrendIcon = TREND_ICONS[pathology.trendDirection]
                const severityLabel = t(`epidemiology.trend.severity.${pathology.trendSeverity}`)

                return (
                  <TableRow key={pathology.code}>
                    <TableCell className="pl-5 font-mono text-xs font-bold text-primary">
                      {pathology.code}
                    </TableCell>
                    <TableCell className="text-sm text-gray-700">{pathology.description}</TableCell>
                    <TableCell>
                      <span className="inline-flex rounded-md border border-border bg-gray-50 px-2 py-0.5 text-[11px] font-semibold text-muted-foreground">
                        {pathology.category}
                      </span>
                    </TableCell>
                    <TableCell className="text-right text-sm font-black tabular-nums text-gray-900">
                      {pathology.activeCases}
                    </TableCell>
                    <TableCell>
                      <span className={cn("flex items-center gap-1 text-xs font-bold", trendTone.text)}>
                        <TrendIcon className="size-3.5 shrink-0" aria-hidden="true" />
                        {resolveTrendLabel(pathology.trend)}
                        <span className="sr-only"> — {severityLabel}</span>
                      </span>
                    </TableCell>
                    <TableCell className="pr-5">
                      <div className="flex items-center justify-end gap-2.5">
                        <div className="h-2 w-28 overflow-hidden rounded-full bg-gray-200">
                          <div
                            className="h-full rounded-full"
                            style={{
                              width: `${pathology.proportion}%`,
                              backgroundColor: toneFillVariables.primary,
                            }}
                          />
                        </div>
                        <span className="w-9 text-right text-[11px] font-bold tabular-nums text-muted-foreground">
                          {pathology.proportion.toFixed(0)}%
                        </span>
                      </div>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        </div>
      ) : (
        <EmptyState
          title={t("empty.epidemiology")}
          description={t("empty.epidemiologyDesc")}
        />
      )}
    </AnalyticsPanel>
  )
}
