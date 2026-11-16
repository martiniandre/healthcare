import { useTranslation } from "react-i18next"
import { Search, X } from "lucide-react"
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
    <div className="relative">
      <Search className="pointer-events-none absolute left-3.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
      <Input
        type="text"
        placeholder={t("searchPlaceholder")}
        value={searchTerm}
        onChange={(event) => onSearchChange(event.target.value)}
        className="h-11 bg-card pl-10 pr-12 sm:pr-36"
      />
      <div className="absolute right-2 top-1/2 flex -translate-y-1/2 items-center gap-1">
        {searchTerm && (
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={() => onSearchChange("")}
            aria-label={t("clearSearch")}
          >
            <X />
          </Button>
        )}
        <span
          aria-live="polite"
          className="sr-only rounded px-2 py-1 text-xs font-medium tabular-nums text-muted-foreground sm:not-sr-only"
        >
          {t("filterResults", { count: resultsCount })}
        </span>
      </div>
    </div>
  )
}
