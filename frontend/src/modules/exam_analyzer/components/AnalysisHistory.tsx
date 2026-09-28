import { FileText, Trash2, Calendar, Database, Search } from "lucide-react"
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
  activeID: string | null
  onSelect: (analysis: ExamAnalysis) => void
  onDelete: (id: string) => void
}

const STATUS_DOT_STYLES: Record<HistoryStatusTone, string> = {
  completed: "bg-success",
  inProgress: "bg-secondary animate-pulse",
  attention: "bg-warning",
  failed: "bg-danger",
}

const STATUS_TEXT_STYLES: Record<HistoryStatusTone, string> = {
  completed: "text-success",
  inProgress: "text-secondary",
  attention: "text-warning",
  failed: "text-danger",
}

export const AnalysisHistory = ({
  history,
  isLoading,
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
    <div className="flex h-fit w-full shrink-0 flex-col gap-4 rounded-xl border border-border bg-card p-4 md:sticky md:top-6 md:max-h-[calc(100vh-120px)] md:w-[320px]">
      <div className="flex items-center justify-between">
        <h3 className="flex items-center gap-2 text-sm font-bold text-gray-900">
          <Database className="h-4 w-4 text-primary" aria-hidden="true" />
          {t("history.title")}
        </h3>
        <span className="rounded-full border border-border bg-gray-50 px-2 py-0.5 text-[11px] font-bold text-muted-foreground">
          {t("history.exams", { count: filteredHistory.length })}
        </span>
      </div>

      <div className="flex shrink-0 items-center gap-2 rounded-lg border border-border bg-gray-50 px-3 py-1.5">
        <Search className="h-3.5 w-3.5 shrink-0 text-gray-400" aria-hidden="true" />
        <input
          type="text"
          placeholder={t("history.filterPlaceholder")}
          aria-label={t("history.filterPlaceholder")}
          value={searchTerm}
          onChange={(event) => setSearchTerm(event.target.value)}
          className="w-full bg-transparent text-xs text-gray-800 placeholder-gray-400 focus:outline-none"
        />
      </div>

      <div className="flex max-h-[300px] flex-1 flex-col gap-2 overflow-y-auto pr-1 md:max-h-none">
        {isLoading ? (
          <div className="py-10 text-center">
            <span className="text-xs text-muted">{t("history.loading")}</span>
          </div>
        ) : filteredHistory.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-2 py-10 text-center">
            <FileText className="h-8 w-8 text-gray-200" aria-hidden="true" />
            <span className="mx-auto block max-w-[180px] text-xs leading-normal text-muted">
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
                    : "border-border bg-card hover:border-border hover:bg-gray-50/60"
                )}
              >
                <button
                  type="button"
                  onClick={() => onSelect(item)}
                  aria-current={isCurrentlySelected ? "true" : undefined}
                  aria-label={t("history.selectAnalysis", { fileName })}
                  className="min-w-0 flex-1 cursor-pointer rounded-md px-2 py-1.5 text-left"
                >
                  <span
                    className={cn(
                      "block truncate text-xs font-semibold transition-colors",
                      isCurrentlySelected ? "text-primary" : "text-gray-900"
                    )}
                  >
                    {fileName}
                  </span>

                  <span className="mt-1 block text-[11px] font-medium text-gray-600">
                    {secondaryLabel}
                  </span>

                  <span className="mt-1.5 flex items-center gap-2.5 text-[11px] text-muted-foreground">
                    <span className="flex items-center gap-1">
                      <Calendar className="h-3 w-3 shrink-0" aria-hidden="true" />
                      {item.created_at ? formatDateTime(item.created_at) : "—"}
                    </span>
                    <span
                      className={cn(
                        "flex items-center gap-1 font-bold uppercase tracking-wider",
                        STATUS_TEXT_STYLES[statusReading.tone]
                      )}
                    >
                      <span
                        className={cn(
                          "h-1.5 w-1.5 shrink-0 rounded-full",
                          STATUS_DOT_STYLES[statusReading.tone]
                        )}
                        aria-hidden="true"
                      />
                      {t(statusReading.translationKey)}
                    </span>
                  </span>
                </button>

                <Can I={Action.Delete} a={Feature.ExamAnalysis}>
                  <button
                    type="button"
                    onClick={() => onDelete(item.id)}
                    aria-label={t("history.deleteAnalysis", { fileName })}
                    className="mt-1.5 shrink-0 cursor-pointer rounded p-1.5 text-gray-400 transition-all hover:bg-danger/10 hover:text-danger focus-visible:opacity-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-danger/40 md:opacity-0 md:group-hover:opacity-100"
                  >
                    <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                  </button>
                </Can>
              </div>
            )
          })
        )}
      </div>
    </div>
  )
}
