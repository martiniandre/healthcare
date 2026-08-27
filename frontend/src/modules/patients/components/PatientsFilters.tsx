import { useTranslation } from "react-i18next"
import { Search, X, Filter } from "lucide-react"
import { Input, Button } from "../../../shared/components/ui"

interface PatientsFiltersProps {
  searchTerm: string
  onSearchChange: (value: string) => void
  resultsCount: number
}

export const PatientsFilters = ({ searchTerm, onSearchChange, resultsCount }: PatientsFiltersProps) => {
  const { t } = useTranslation("patients")

  return (
    <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
      <div className="flex-1 flex items-center gap-2 bg-surface border border-border rounded-lg px-3.5">
        <Search className="w-4 h-4 text-muted-foreground/70 shrink-0" />
        <div className="relative w-full">
          <Input
            type="text"
            placeholder={t("searchPlaceholder")}
            value={searchTerm}
            onChange={(event) => onSearchChange(event.target.value)}
            className="px-0 h-9 bg-transparent border-transparent shadow-none focus-visible:border-transparent focus-visible:ring-0"
          />
        </div>
        {searchTerm && (
          <Button
            variantType="ghost"
            size="sm"
            aria-label={t("clearFilters")}
            onClick={() => onSearchChange("")}
            className="p-0 h-6 w-6 rounded-md text-muted-foreground hover:text-foreground"
          >
            <X className="w-3.5 h-3.5" />
          </Button>
        )}
      </div>

      <div className="flex items-center justify-center sm:justify-start gap-1.5 bg-surface border border-border rounded-lg px-3 py-2 shrink-0">
        <Filter className="w-3.5 h-3.5 text-muted-foreground/70" />
        <span className="text-[11px] text-muted-foreground font-medium">
          {t("filterResults", { count: resultsCount })}
        </span>
      </div>
    </div>
  )
}