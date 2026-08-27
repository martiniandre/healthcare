import { FileText, Trash2, Calendar, Database, Search } from "lucide-react"
import { useState } from "react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { formatDateTime } from "../../../shared/utils/dates"
import { Input } from "../../../shared/components/ui/Input"
import { Button } from "../../../shared/components/ui/Button"
import { ExamAnalysisStatus, type ExamAnalysis } from "../types"

interface AnalysisHistoryProperties {
  history: ExamAnalysis[]
  isLoading: boolean
  activeID: string | null
  onSelect: (analysis: ExamAnalysis) => void
  onDelete: (id: string) => void
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

  const filteredHistory = history.filter((item) => {
    const term = searchTerm.toLowerCase()
    return (
      item.file_name?.toLowerCase().includes(term) ||
      (item.exam_type && item.exam_type.toLowerCase().includes(term))
    )
  })

  return (
    <div className="w-full md:w-[320px] shrink-0 flex flex-col gap-4 bg-surface border border-border rounded-xl p-4 h-fit md:max-h-[calc(100vh-120px)] md:sticky md:top-6">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-bold text-foreground flex items-center gap-2">
          <Database className="w-4 h-4 text-primary" />
          {t("history.title")}
        </h3>
        <span className="text-[10px] text-muted font-bold bg-muted-soft border border-border/80 px-2 py-0.5 rounded-full">
          {t("history.exams", { count: filteredHistory.length })}
        </span>
      </div>

      <div className="flex items-center gap-2 bg-muted-soft border border-border rounded-lg px-3 py-1.5 shrink-0">
        <Search className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
        <Input
          type="text"
          placeholder={t("history.filterPlaceholder")}
          value={searchTerm}
          onChange={(event) => setSearchTerm(event.target.value)}
          className="w-full bg-transparent border-none shadow-none px-0 h-6 text-xs text-foreground placeholder:text-muted-foreground focus-visible:ring-0 focus-visible:border-transparent"
        />
      </div>

      <div className="flex-1 flex flex-col gap-2 overflow-y-auto max-h-[300px] md:max-h-none pr-1">
        {isLoading ? (
          <div className="py-10 text-center">
            <span className="text-xs text-muted">{t("history.loading")}</span>
          </div>
        ) : filteredHistory.length === 0 ? (
          <div className="py-10 text-center flex flex-col items-center justify-center gap-2">
            <FileText className="w-8 h-8 text-muted-soft" />
            <span className="text-xs text-muted block max-w-[180px] leading-normal mx-auto">
              {searchTerm ? t("history.noResults") : t("history.empty")}
            </span>
          </div>
        ) : (
          filteredHistory.map((item) => {
            const isCurrentlySelected = activeID === item.id

            const handleItemDelete = (event: React.MouseEvent) => {
              event.stopPropagation()
              onDelete(item.id)
            }

            return (
              <div
                key={item.id}
                onClick={() => onSelect(item)}
                className={`flex items-start justify-between gap-3 p-3 rounded-lg border transition-all duration-200 cursor-pointer select-none group ${isCurrentlySelected
                  ? "bg-primary/5 border-primary/30 shadow-sm"
                  : "bg-surface border-border/80 hover:bg-muted-soft hover:border-border-strong"
                  }`}
              >
                <div className="min-w-0 flex-1 flex flex-col gap-1 text-left">
                  <span className={`text-xs font-semibold block truncate group-hover:text-primary transition-colors ${isCurrentlySelected ? "text-primary" : "text-foreground"
                    }`}>
                    {item.file_name || t("history.unknownFile", "Unknown File")}
                  </span>

                  <span className="text-[10px] text-muted-foreground font-medium block">
                    {item.exam_type || (
                      item.status === ExamAnalysisStatus.PENDING || item.status === ExamAnalysisStatus.PROCESSING
                        ? t("history.processing")
                        : t("history.insufficient")
                    )}
                  </span>

                  <div className="flex items-center gap-3 mt-1 text-[9px] text-muted">
                    <span className="flex items-center gap-1">
                      <Calendar className="w-3 h-3 shrink-0" />
                      {item.created_at ? formatDateTime(item.created_at) : "—"}
                    </span>
                    <span className={`font-bold uppercase tracking-wider ${item.status === ExamAnalysisStatus.COMPLETED
                      ? "text-primary/70"
                      : item.status === ExamAnalysisStatus.PENDING || item.status === ExamAnalysisStatus.PROCESSING
                        ? "text-secondary/70 animate-pulse"
                        : "text-danger/70"
                      }`}>
                      {item.status === ExamAnalysisStatus.COMPLETED && t("history.statusCompleted")}
                      {(item.status === ExamAnalysisStatus.PENDING || item.status === ExamAnalysisStatus.PROCESSING) && t("history.statusQueue")}
                      {item.status === ExamAnalysisStatus.INSUFFICIENT_DATA && t("history.statusQuality")}
                      {item.status === ExamAnalysisStatus.FAILED && t("history.statusFailed")}
                    </span>
                  </div>
                </div>

                <Can I={Action.Delete} a={Feature.ExamAnalysis}>
                  <Button
                    type="button"
                    variantType="ghost"
                    size="sm"
                    onClick={handleItemDelete}
                    className="p-1 h-auto w-auto shadow-none text-muted-foreground hover:text-danger hover:bg-danger-soft transition-all opacity-0 group-hover:opacity-100 cursor-pointer shrink-0 rounded"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </Button>
                </Can>
              </div>
            )
          })
        )}
      </div>
    </div>
  )
}
