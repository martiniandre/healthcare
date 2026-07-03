import { useTranslation } from "react-i18next"
import { Search, Filter, XCircle } from "lucide-react"
import { Alert, AlertTitle, AlertDescription } from "../../../shared/components/ui/Alert"

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
          <Search className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            placeholder={t("filterEmailPlaceholder")}
            value={userEmail}
            onChange={(event) => onUserEmailChange(event.target.value)}
            className="w-full bg-white border border-border rounded-lg pl-9 pr-4 py-2 text-xs text-gray-800 placeholder-gray-400 focus:outline-none focus:border-primary/50 transition-all duration-200 h-9"
          />
        </div>

        <div>
          <select
            value={filterAction}
            onChange={(event) => onFilterActionChange(event.target.value)}
            className="w-full bg-white border border-border rounded-lg px-3 py-2 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200 h-9 cursor-pointer"
          >
            <option value="All">{t("filterAllActions")}</option>
            <option value="PAGE_VIEW">{t("actionPageView")}</option>
            <option value="LOGIN">{t("actionLogin")}</option>
            <option value="LOGOUT">{t("actionLogout")}</option>
            <option value="API_REQUEST">{t("actionApiRequest")}</option>
          </select>
        </div>

        <div>
          <select
            value={filterStatus}
            onChange={(event) => onFilterStatusChange(event.target.value)}
            className="w-full bg-white border border-border rounded-lg px-3 py-2 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200 h-9 cursor-pointer"
          >
            <option value="All">{t("filterAllStatuses")}</option>
            <option value="SUCCESS">{t("statusSuccess")}</option>
            <option value="FAILURE">{t("statusFailure")}</option>
          </select>
        </div>

        <div>
          <input
            type="date"
            value={startDate}
            onChange={(event) => onStartDateChange(event.target.value)}
            className="w-full bg-white border border-border rounded-lg px-3 py-2 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200 h-9 cursor-pointer"
          />
        </div>

        <div>
          <input
            type="date"
            value={endDate}
            onChange={(event) => onEndDateChange(event.target.value)}
            className="w-full bg-white border border-border rounded-lg px-3 py-2 text-xs text-gray-800 focus:outline-none focus:border-primary/50 transition-all duration-200 h-9 cursor-pointer"
          />
        </div>
      </div>

      {hasActiveFilters && (
        <div className="flex justify-end">
          <button
            onClick={onResetFilters}
            className="text-xs text-primary hover:underline font-bold flex items-center gap-1"
          >
            <Filter className="w-3 h-3" />
            {t("clearFilters")}
          </button>
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
