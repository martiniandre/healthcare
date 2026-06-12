import { useTranslation } from "react-i18next"
import { Users, UserPlus } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"

interface StaffHeaderProps {
  onAddStaff: () => void
}

export const StaffHeader = ({ onAddStaff }: StaffHeaderProps) => {
  const { t } = useTranslation()

  return (
    <div className="flex items-center justify-between flex-wrap gap-4">
      <div className="text-left">
        <h2 className="text-xl font-black text-gray-900 leading-none flex items-center gap-2">
          <Users className="w-5 h-5 text-primary animate-pulse-glow" />
          {t("staff.title")}
        </h2>
        <span className="text-xs text-muted mt-1.5 block">
          {t("staff.subtitle")}
        </span>
      </div>

      <Button
        variantType="primary"
        onClick={onAddStaff}
        className="px-4 gap-2 text-xs font-bold"
      >
        <UserPlus className="w-4 h-4" />
        {t("staff.addStaff")}
      </Button>
    </div>
  )
}
