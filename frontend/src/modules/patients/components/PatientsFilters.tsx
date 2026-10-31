import { useTranslation } from "react-i18next"
import { Search, X, Filter } from "lucide-react"
import { Button } from "../../../shared/components/ui/button"
import { Input } from "../../../shared/components/ui/input"

interface PatientsFiltersProps {
  searchTerm: string
  onSearchChange: (value: string) => void
  resultsCount: number
}

export const PatientsFilters = ({ searchTerm, onSearchChange, resultsCount }: PatientsFiltersProps) => {
  const { t } = useTranslation("patients")

  return (
    <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
      <div className="relative flex-1">
        <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          type="text"
          placeholder={t("searchPlaceholder")}
          value={searchTerm}
          onChange={(event) => onSearchChange(event.target.value)}
          className="bg-card pl-9 pr-9"
        />
        {searchTerm && (
          <Button
            type="button"
            variant="ghost"
            size="icon-xs"
            onClick={() => onSearchChange("")}
            aria-label={t("clearSearch")}
            className="absolute right-2 top-1/2 -translate-y-1/2"
          >
            <X />
          </Button>
        )}
      </div>

      <div className="flex items-center justify-center sm:justify-start gap-1.5 bg-card border border-border rounded-lg px-3 py-2.5 shrink-0">
        <Filter className="w-3.5 h-3.5 text-muted-foreground" />
        <span className="text-[11px] text-muted-foreground font-medium">
          {t("filterResults", { count: resultsCount })}
        </span>
      </div>
    </div>
  )
}
