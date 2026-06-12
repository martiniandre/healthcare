import { useTranslation } from "react-i18next"
import { Search, X, Filter } from "lucide-react"

interface PatientsFiltersProps {
  searchTerm: string
  onSearchChange: (value: string) => void
  resultsCount: number
}

export const PatientsFilters = ({ searchTerm, onSearchChange, resultsCount }: PatientsFiltersProps) => {
  const { t } = useTranslation()

  return (
    <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
      <div className="flex-1 flex items-center gap-2.5 bg-white border border-border rounded-lg px-4 py-2.5">
        <Search className="w-4 h-4 text-gray-400 shrink-0" />
        <input
          type="text"
          placeholder={t("patients.searchPlaceholder")}
          value={searchTerm}
          onChange={(event) => onSearchChange(event.target.value)}
          className="w-full bg-transparent text-sm text-gray-800 placeholder-gray-400 focus:outline-none"
        />
        {searchTerm && (
          <button
            onClick={() => onSearchChange("")}
            className="p-0.5 rounded hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors"
          >
            <X className="w-3.5 h-3.5" />
          </button>
        )}
      </div>

      <div className="flex items-center justify-center sm:justify-start gap-1.5 bg-white border border-border rounded-lg px-3 py-2.5 shrink-0">
        <Filter className="w-3.5 h-3.5 text-gray-400" />
        <span className="text-[11px] text-muted font-medium">
          {t("patients.filterResults", { count: resultsCount })}
        </span>
      </div>
    </div>
  )
}
