import { AlertCircle, FileText, Trash2, Calendar, Database, Search } from "lucide-react"
import { Card, CardAction, CardContent, CardHeader, CardTitle } from "@/shared/components/ui/card"
import { Badge } from "@/shared/components/ui/badge"
import { Button } from "@/shared/components/ui/button"
import { Input } from "@/shared/components/ui/input"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { Skeleton } from "@/shared/components/ui/skeleton"
import { StatusBadge, type ClinicalStatus } from "@/shared/components/ui/status-badge"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/shared/components/ui/tooltip"
import { useState } from "react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { formatDateTime } from "../../../shared/utils/dates"
import { ExamAnalysisStatus, type ExamAnalysis } from "../types"
import { resolveHistoryStatus, type HistoryStatusTone } from "../exam_analysis_domain"
import { cn } from "../../../shared/utils/cn"

interface AnalysisHistoryProperties {
  history: ExamAnalysis[]
  isLoading: boolean
  isError?: boolean
  onRetry?: () => void
  activeID: string | null
  onSelect: (analysis: ExamAnalysis) => void
  onDelete: (id: string) => void
}

const HISTORY_TONE_TO_STATUS: Record<HistoryStatusTone, ClinicalStatus> = {
  completed: "normal",
  inProgress: "info",
  attention: "attention",
  failed: "critical",
}

export const AnalysisHistory = ({
  history,
  isLoading,
  isError = false,
  onRetry,
  activeID,
  onSelect,
  onDelete,
}: AnalysisHistoryProperties) => {
  const { t } = useTranslation("examAnalyzer")
  const [searchTerm, setSearchTerm] = useState<string>("")

  const normalizedSearchTerm = searchTerm.trim().toLowerCase()

  const matchesSearchTerm = (item: ExamAnalysis): boolean => {
    if (normalizedSearchTerm.length === 0) return true

    const statusLabel = t(resolveHistoryStatus(item.status).translationKey).toLowerCase()

    return (
      item.file_name?.toLowerCase().includes(normalizedSearchTerm) === true ||
      item.exam_type?.toLowerCase().includes(normalizedSearchTerm) === true ||
      statusLabel.includes(normalizedSearchTerm)
    )
  }

  const filteredHistory = history.filter(matchesSearchTerm)

  return (
    <Card className="h-fit w-full shrink-0 gap-4 py-4 md:sticky md:top-6 md:max-h-[calc(100vh-120px)] md:w-[320px]">
      <CardHeader className="px-4">
        <CardTitle className="flex items-center gap-2 text-sm font-bold">
          <Database className="size-4 text-primary" aria-hidden="true" />
          {t("history.title")}
        </CardTitle>
        <CardAction>
          <Badge variant="outline" className="text-[11px] font-bold text-muted-foreground">
            {t("history.exams", { count: filteredHistory.length })}
          </Badge>
        </CardAction>
      </CardHeader>

      <CardContent className="flex min-h-0 flex-1 flex-col gap-4 px-4">
        <div className="relative shrink-0">
          <Search
            className="absolute left-3 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground"
            aria-hidden="true"
          />
          <Input
            type="text"
            placeholder={t("history.filterPlaceholder")}
            aria-label={t("history.filterPlaceholder")}
            value={searchTerm}
            onChange={(event) => setSearchTerm(event.target.value)}
            className="pl-9 text-xs md:text-xs"
          />
        </div>

        <div className="flex max-h-[300px] flex-1 flex-col gap-2 overflow-y-auto pr-1 md:max-h-none">
          {isLoading ? (
            <div className="flex flex-col gap-2" role="status" aria-label={t("history.loading")}>
              <Skeleton className="h-16 rounded-lg" />
              <Skeleton className="h-16 rounded-lg" />
              <Skeleton className="h-16 rounded-lg" />
            </div>
          ) : isError ? (
            <Alert variant="destructive">
              <AlertCircle aria-hidden="true" />
              <AlertDescription className="flex flex-col items-start gap-2">
                <span>{t("history.loadError")}</span>
                <Button type="button" variant="outline" size="xs" onClick={onRetry}>
                  {t("retry")}
                </Button>
              </AlertDescription>
            </Alert>
          ) : filteredHistory.length === 0 ? (
            <div className="flex flex-col items-center justify-center gap-2 py-10 text-center">
              <FileText className="size-8 text-muted-foreground/40" aria-hidden="true" />
              <span className="mx-auto block max-w-[180px] text-xs leading-normal text-muted-foreground">
                {normalizedSearchTerm.length > 0 ? t("history.noResults") : t("history.empty")}
              </span>
            </div>
          ) : (
            filteredHistory.map((item) => {
              const isCurrentlySelected = activeID === item.id
              const fileName = item.file_name || t("history.unknownFile")
              const statusReading = resolveHistoryStatus(item.status)
              const isInFlight =
                item.status === ExamAnalysisStatus.PENDING ||
                item.status === ExamAnalysisStatus.PROCESSING

              const secondaryLabel =
                item.exam_type ||
                (isInFlight ? t("history.processing") : t("history.insufficient"))

              return (
                <div
                  key={item.id}
                  className={cn(
                    "group flex items-start gap-1 rounded-lg border p-1 transition-all duration-200",
                    isCurrentlySelected
                      ? "border-primary/30 bg-primary/5 shadow-sm"
                      : "border-border bg-card hover:bg-muted/60"
                  )}
                >
                  <Button
                    type="button"
                    variant="ghost"
                    onClick={() => onSelect(item)}
                    aria-current={isCurrentlySelected ? "true" : undefined}
                    aria-label={t("history.selectAnalysis", { fileName })}
                    className="h-auto min-w-0 flex-1 flex-col items-start justify-start gap-0 whitespace-normal px-2 py-1.5 text-left font-normal hover:bg-transparent"
                  >
                    <span
                      className={cn(
                        "block w-full truncate text-xs font-semibold transition-colors",
                        isCurrentlySelected ? "text-primary" : "text-foreground"
                      )}
                    >
                      {fileName}
                    </span>

                    <span className="mt-1 block text-[11px] font-medium text-muted-foreground">
                      {secondaryLabel}
                    </span>

                    <span className="mt-1.5 flex items-center gap-2.5 text-[11px] text-muted-foreground">
                      <span className="flex items-center gap-1">
                        <Calendar className="size-3 shrink-0" aria-hidden="true" />
                        {item.created_at ? formatDateTime(item.created_at) : "—"}
                      </span>
                      <StatusBadge
                        status={HISTORY_TONE_TO_STATUS[statusReading.tone]}
                        className="text-[10px] font-bold uppercase tracking-wider"
                      >
                        {t(statusReading.translationKey)}
                      </StatusBadge>
                    </span>
                  </Button>

                  <Can I={Action.Delete} a={Feature.ExamAnalysis}>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon-xs"
                          onClick={() => onDelete(item.id)}
                          aria-label={t("history.deleteAnalysis", { fileName })}
                          className="mt-1.5 shrink-0 text-muted-foreground hover:bg-destructive/10 hover:text-destructive focus-visible:opacity-100 md:opacity-0 md:group-hover:opacity-100"
                        >
                          <Trash2 aria-hidden="true" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>{t("history.deleteAnalysis", { fileName })}</TooltipContent>
                    </Tooltip>
                  </Can>
                </div>
              )
            })
          )}
        </div>
      </CardContent>
    </Card>
  )
}
