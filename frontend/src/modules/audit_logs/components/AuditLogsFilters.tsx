import { useTranslation } from "react-i18next"
import { Search, Filter, XCircle } from "lucide-react"
import { Button } from "../../../shared/components/ui/button"
import { Input } from "../../../shared/components/ui/input"
import { Label } from "../../../shared/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../../shared/components/ui/select"
import { Alert, AlertTitle, AlertDescription } from "../../../shared/components/ui/alert"

interface AuditLogsFiltersProps {
  filterAction: string
  onFilterActionChange: (value: string) => void
  filterStatus: string
  onFilterStatusChange: (value: string) => void
  userEmail: string
  onUserEmailChange: (value: string) => void
  startDate: string
  onStartDateChange: (value: string) => void
  endDate: string
  onEndDateChange: (value: string) => void
  onResetFilters: () => void
  isError: boolean
  onRetry: () => void
}

export const AuditLogsFilters = ({
  filterAction,
  onFilterActionChange,
  filterStatus,
  onFilterStatusChange,
  userEmail,
  onUserEmailChange,
  startDate,
  onStartDateChange,
  endDate,
  onEndDateChange,
  onResetFilters,
  isError,
  onRetry,
}: AuditLogsFiltersProps) => {
  const { t } = useTranslation("auditLogs")

  const hasActiveFilters = filterAction !== "All" || filterStatus !== "All" || userEmail || startDate || endDate

  return (
    <>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-3">
        <div className="relative">
          <Search className="w-4 h-4 text-muted-foreground absolute left-3 top-1/2 -translate-y-1/2" aria-hidden="true" />
          <Input
            type="text"
            aria-label={t("filterEmailPlaceholder")}
            placeholder={t("filterEmailPlaceholder")}
            value={userEmail}
            onChange={(event) => onUserEmailChange(event.target.value)}
            className="pl-9 text-xs"
          />
        </div>

        <Select value={filterAction} onValueChange={onFilterActionChange}>
          <SelectTrigger className="w-full text-xs" aria-label={t("filterAllActions")}>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="All">{t("filterAllActions")}</SelectItem>
            <SelectItem value="LOGIN">{t("actionLogin")}</SelectItem>
            <SelectItem value="LOGOUT">{t("actionLogout")}</SelectItem>
            <SelectItem value="API_REQUEST">{t("actionApiRequest")}</SelectItem>
          </SelectContent>
        </Select>

        <Select value={filterStatus} onValueChange={onFilterStatusChange}>
          <SelectTrigger className="w-full text-xs" aria-label={t("filterAllStatuses")}>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="All">{t("filterAllStatuses")}</SelectItem>
            <SelectItem value="SUCCESS">{t("statusSuccess")}</SelectItem>
            <SelectItem value="FAILURE">{t("statusFailure")}</SelectItem>
          </SelectContent>
        </Select>

        <div className="relative">
          <Label htmlFor="audit-start-date" className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-xs font-medium text-muted-foreground">
            {t("startDate")}
          </Label>
          <Input
            id="audit-start-date"
            type="date"
            value={startDate}
            max={endDate || undefined}
            onChange={(event) => onStartDateChange(event.target.value)}
            className="pl-12 text-xs"
          />
        </div>

        <div className="relative">
          <Label htmlFor="audit-end-date" className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-xs font-medium text-muted-foreground">
            {t("endDate")}
          </Label>
          <Input
            id="audit-end-date"
            type="date"
            value={endDate}
            min={startDate || undefined}
            onChange={(event) => onEndDateChange(event.target.value)}
            className="pl-12 text-xs"
          />
        </div>
      </div>

      {hasActiveFilters && (
        <div className="flex justify-end">
          <Button type="button" variant="link" size="sm" onClick={onResetFilters} className="text-xs font-semibold">
            <Filter aria-hidden="true" />
            {t("clearFilters")}
          </Button>
        </div>
      )}
      {isError && (
        <Alert variant="destructive">
          <XCircle className="w-4 h-4" />
          <AlertTitle>{t("errorTitle")}</AlertTitle>
          <AlertDescription>
            <p>{t("errorDescription")}</p>
            <Button type="button" variant="outline" size="sm" onClick={onRetry} className="mt-2">
              {t("retry")}
            </Button>
          </AlertDescription>
        </Alert>
      )}
    </>
  )
}
