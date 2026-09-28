import { useTranslation } from "react-i18next"
import { Search } from "lucide-react"
import { StaffRole } from "../../../shared/types"
import { STAFF_ROLE_LABEL_KEYS } from "../staff_roles"

interface StaffFiltersProps {
  searchQuery: string
  onSearchChange: (value: string) => void
  filterRole: string
  onFilterChange: (role: string) => void
}

const roleFilterOptions = ["All", StaffRole.Doctor, StaffRole.Nurse, StaffRole.Receptionist, StaffRole.Admin]

export const StaffFilters = ({
  searchQuery,
  onSearchChange,
  filterRole,
  onFilterChange,
}: StaffFiltersProps) => {
  const { t } = useTranslation("staff")

  const getRoleLabel = (role: string) => {
    if (role === "All") {
      return t("filterAll")
    }
    const labelKey = STAFF_ROLE_LABEL_KEYS[role as StaffRole]
    return labelKey ? t(labelKey) : role
  }

  return (
    <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div className="relative flex-1 max-w-md">
        <Search className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input
          type="text"
          placeholder={t("searchPlaceholder")}
          value={searchQuery === " " ? "" : searchQuery}
          onChange={(event) => onSearchChange(event.target.value)}
          className="w-full bg-card border border-border rounded-lg pl-9 pr-4 py-2 text-xs text-gray-800 placeholder-gray-400 focus:outline-none focus:border-primary/50 transition-all duration-200"
        />
      </div>

      <div className="flex gap-2 flex-wrap">
        {roleFilterOptions.map((roleOption) => (
          <button
            key={roleOption}
            onClick={() => onFilterChange(roleOption)}
            className={`px-3 py-1.5 rounded-lg text-xs font-bold transition-all duration-200 border ${
              filterRole === roleOption
                ? "bg-primary/5 text-primary border-primary"
                : "bg-card text-gray-500 border-border hover:bg-gray-50 hover:text-gray-900"
            }`}
          >
            {getRoleLabel(roleOption)}
          </button>
        ))}
      </div>
    </div>
  )
}
