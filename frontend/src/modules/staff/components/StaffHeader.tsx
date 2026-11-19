import { useTranslation } from "react-i18next"
import { UserPlus } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import { PageHeader } from "@/shared/components/ui/PageContainer"

interface StaffHeaderProps {
  onAddStaff: () => void
}

export const StaffHeader = ({ onAddStaff }: StaffHeaderProps) => {
  const { t } = useTranslation("staff")

  return (
    <PageHeader
      title={t("title")}
      actions={
        <Button onClick={onAddStaff}>
          <UserPlus />
          {t("addStaff")}
        </Button>
      }
    />
  )
}
