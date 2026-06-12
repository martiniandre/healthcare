import { FileText, Trash2, Calendar, Database, Search } from "lucide-react"
import { useState } from "react"
import { useTranslation } from "react-i18next"
import type { ExamAnalysis } from "../types"

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
  const { t, i18n } = useTranslation("examAnalyzer")
  const [searchTerm, setSearchTerm] = useState<string>("")

  const filteredHistory = history.filter((item) => {
    const term = searchTerm.toLowerCase()
    return (
      item.file_name.toLowerCase().includes(term) ||
      (item.exam_type && item.exam_type.toLowerCase().includes(term))
    )
  })

  return (
    <div className="w-full md:w-[320px] shrink-0 flex flex-col gap-4 bg-white border border-border rounded-xl p-4 h-fit md:max-h-[calc(100vh-120px)] md:sticky md:top-6">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-bold text-gray-900 flex items-center gap-2">
          <Database className="w-4 h-4 text-primary" />
          {t("history.title")}
        </h3>
        <span className="text-[10px] text-muted font-bold bg-gray-50 border border-border/80 px-2 py-0.5 rounded-full">
          {t("history.exams", { count: filteredHistory.length })}
        </span>
      </div>

      <div className="flex items-center gap-2 bg-gray-50 border border-border rounded-lg px-3 py-1.5 shrink-0">
        <Search className="w-3.5 h-3.5 text-gray-400 shrink-0" />
        <input
          type="text"
          placeholder={t("history.filterPlaceholder")}
          value={searchTerm}
          onChange={(event) => setSearchTerm(event.target.value)}
          className="w-full bg-transparent text-xs text-gray-800 placeholder-gray-400 focus:outline-none"
        />
      </div>

      <div className="flex-1 flex flex-col gap-2 overflow-y-auto max-h-[300px] md:max-h-none pr-1">
        {isLoading ? (
          <div className="py-10 text-center">
            <span className="text-xs text-muted">{t("history.loading")}</span>
          </div>
        ) : filteredHistory.length === 0 ? (
          <div className="py-10 text-center flex flex-col items-center justify-center gap-2">
            <FileText className="w-8 h-8 text-gray-200" />
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
                className={`flex items-start justify-between gap-3 p-3 rounded-lg border transition-all duration-200 cursor-pointer select-none group ${
                  isCurrentlySelected
                    ? "bg-primary/5 border-primary/30 shadow-sm"
                    : "bg-white border-border/80 hover:bg-gray-50/50 hover:border-gray-300"
                }`}
              >
                <div className="min-w-0 flex-1 flex flex-col gap-1 text-left">
                  <span className={`text-xs font-semibold block truncate group-hover:text-primary transition-colors ${
                    isCurrentlySelected ? "text-primary" : "text-gray-800"
                  }`}>
                    {item.file_name}
                  </span>
                  
                  <span className="text-[10px] text-gray-500 font-medium block">
                    {item.exam_type || (
                      item.status === "pending" || item.status === "processing" 
                        ? t("history.processing") 
                        : t("history.insufficient")
                    )}
                  </span>

                  <div className="flex items-center gap-3 mt-1 text-[9px] text-muted">
                    <span className="flex items-center gap-1">
                      <Calendar className="w-3 h-3 shrink-0" />
                      {new Date(item.created_at).toLocaleDateString(i18n.language, {
                        day: "2-digit",
                        month: "2-digit",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </span>
                    <span className={`font-bold uppercase tracking-wider ${
                      item.status === "completed" 
                        ? "text-primary/70"
                        : item.status === "pending" || item.status === "processing"
                          ? "text-secondary/70 animate-pulse"
                          : "text-red-500/70"
                    }`}>
                      {item.status === "completed" && t("history.statusCompleted")}
                      {(item.status === "pending" || item.status === "processing") && t("history.statusQueue")}
                      {item.status === "insufficient_data" && t("history.statusQuality")}
                      {item.status === "failed" && t("history.statusFailed")}
                    </span>
                  </div>
                </div>

                <button
                  type="button"
                  onClick={handleItemDelete}
                  className="p-1 rounded text-gray-400 hover:text-red-500 hover:bg-red-50 transition-all opacity-0 group-hover:opacity-100 cursor-pointer shrink-0"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                </button>
              </div>
            )
          })
        )}
      </div>
    </div>
  )
}
