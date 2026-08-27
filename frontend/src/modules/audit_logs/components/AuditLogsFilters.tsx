import { useTranslation } from "react-i18next"
import { Search, Filter, XCircle } from "lucide-react"
import {
  Alert,
  AlertTitle,
  AlertDescription,
} from "../../../shared/components/ui/Alert"
import { Input } from "../../../shared/components/ui/Input"
import { Button } from "../../../shared/components/ui/Button"
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "../../../shared/components/ui/Select"

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
}: AuditLogsFiltersProps) => {
  const { t } = useTranslation("auditLogs")

  const hasActiveFilters = filterAction !== "All" || filterStatus !== "All" || userEmail || startDate || endDate

  return (
    <>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-3">
        <div className="relative">
          <Search className="w-4 h-4 text-muted-foreground/70 absolute left-3 top-1/2 -translate-y-1/2" />
          <div className="relative">
            <Input
              type="text"
              placeholder={t("filterEmailPlaceholder")}
              value={userEmail}
              onChange={(event) => onUserEmailChange(event.target.value)}
              className="pl-9"
            />
          </div>
        </div>

        <div>
          <Select value={filterAction} onValueChange={onFilterActionChange}>
            <SelectTrigger>
              <SelectValue placeholder={t("filterAllActions")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="All">{t("filterAllActions")}</SelectItem>
              <SelectItem value="LOGIN">{t("actionLogin")}</SelectItem>
              <SelectItem value="LOGOUT">{t("actionLogout")}</SelectItem>
              <SelectItem value="API_REQUEST">{t("actionApiRequest")}</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div>
          <Select value={filterStatus} onValueChange={onFilterStatusChange}>
            <SelectTrigger>
              <SelectValue placeholder={t("filterAllStatuses")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="All">{t("filterAllStatuses")}</SelectItem>
              <SelectItem value="SUCCESS">{t("statusSuccess")}</SelectItem>
              <SelectItem value="FAILURE">{t("statusFailure")}</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div>
          <Input
            type="date"
            value={startDate}
            onChange={(event) => onStartDateChange(event.target.value)}
          />
        </div>

        <div>
          <Input
            type="date"
            value={endDate}
            onChange={(event) => onEndDateChange(event.target.value)}
          />
        </div>
      </div>

      {hasActiveFilters && (
        <div className="flex justify-end">
          <Button
            variantType="ghost"
            size="sm"
            onClick={onResetFilters}
            className="text-primary hover:text-primary px-0"
          >
            <Filter className="w-3 h-3" />
            {t("clearFilters")}
          </Button>
        </div>
      )}
      {isError && (
        <Alert variant="destructive">
          <XCircle className="w-4 h-4" />
          <AlertTitle>{t("errorTitle")}</AlertTitle>
          <AlertDescription>
            {t("errorDescription")}
          </AlertDescription>
        </Alert>
      )}
    </>
  )
}
