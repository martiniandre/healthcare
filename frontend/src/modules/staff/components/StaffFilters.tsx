import { useTranslation } from "react-i18next"
import { Search } from "lucide-react"
import { Input } from "../../../shared/components/ui/Input"
import { Button } from "../../../shared/components/ui/Button"
import { StaffRole } from "../../../shared/types"

interface StaffFiltersProps {
  searchQuery: string
  onSearchChange: (value: string) => void
  filterRole: string
  onFilterChange: (role: string) => void
}

export const StaffFilters = ({
  searchQuery,
  onSearchChange,
  filterRole,
  onFilterChange,
}: StaffFiltersProps) => {
  const { t } = useTranslation("staff")

  return (
    <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div className="relative flex-1 max-w-md">
        <Search className="w-4 h-4 text-muted-foreground/70 absolute left-3 top-1/2 -translate-y-1/2" />
        <div className="relative">
          <Input
            type="text"
            placeholder={t("searchPlaceholder")}
            value={searchQuery === " " ? "" : searchQuery}
            onChange={(event) => onSearchChange(event.target.value)}
            className="pl-9"
          />
        </div>
      </div>

      <div className="flex gap-2 flex-wrap">
        {["All", StaffRole.Doctor, StaffRole.Nurse, StaffRole.Receptionist, StaffRole.Admin].map((roleOption) => {
          const getRoleLabel = (role: string) => {
            switch (role) {
              case "All": return t("filterAll")
              case StaffRole.Doctor: return t("table.roles.doctor", "Médico")
              case StaffRole.Nurse: return t("table.roles.nurse", "Enfermeiro")
              case StaffRole.Receptionist: return t("table.roles.receptionist", "Recepção")
              case StaffRole.Admin: return t("table.roles.admin", "Admin")
              default: return role
            }
          }

          return (
            <Button
              key={roleOption}
              variantType="outline"
              size="sm"
              onClick={() => onFilterChange(roleOption)}
              className={
                filterRole === roleOption
                  ? "border-primary bg-primary-soft text-primary hover:bg-primary-soft"
                  : "text-muted-foreground hover:text-foreground"
              }
            >
              {getRoleLabel(roleOption)}
            </Button>
          )
        })}
      </div>
    </div>
  )
}
