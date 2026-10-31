import { useTranslation } from "react-i18next"
import { Users, UserPlus } from "lucide-react"
import { Button } from "@/shared/components/ui/button"

interface StaffHeaderProps {
  onAddStaff: () => void
}

export const StaffHeader = ({ onAddStaff }: StaffHeaderProps) => {
  const { t } = useTranslation("staff")

  return (
    <div className="flex items-center justify-between flex-wrap gap-4">
      <div className="text-left">
        <h2 className="text-xl font-display font-bold text-foreground leading-none flex items-center gap-2">
          <Users className="w-5 h-5 text-primary animate-pulse-glow" />
          {t("title")}
        </h2>
        <span className="text-xs text-muted-foreground mt-1.5 block">
          {t("subtitle")}
        </span>
      </div>

      <Button onClick={onAddStaff} className="text-xs font-bold">
        <UserPlus />
        {t("addStaff")}
      </Button>
    </div>
  )
}
