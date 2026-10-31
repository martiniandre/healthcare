import { useTranslation } from "react-i18next"
import { Search } from "lucide-react"
import { Input } from "@/shared/components/ui/input"
import { Button } from "@/shared/components/ui/button"
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
        <Search className="size-4 text-muted-foreground absolute left-3 top-1/2 -translate-y-1/2" />
        <Input
          type="text"
          placeholder={t("searchPlaceholder")}
          value={searchQuery === " " ? "" : searchQuery}
          onChange={(event) => onSearchChange(event.target.value)}
          className="pl-9 text-xs"
        />
      </div>

      <div className="flex gap-2 flex-wrap">
        {roleFilterOptions.map((roleOption) => (
          <Button
            key={roleOption}
            type="button"
            size="sm"
            variant={filterRole === roleOption ? "default" : "outline"}
            aria-pressed={filterRole === roleOption}
            onClick={() => onFilterChange(roleOption)}
            className="text-xs font-bold"
          >
            {getRoleLabel(roleOption)}
          </Button>
        ))}
      </div>
    </div>
  )
}
