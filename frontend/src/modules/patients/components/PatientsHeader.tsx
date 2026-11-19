import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/button"
import { PageHeader } from "../../../shared/components/ui/PageContainer"
import { UserPlus } from "lucide-react"

interface PatientsHeaderProps {
  onNewPatient: () => void
}

export const PatientsHeader = ({ onNewPatient }: PatientsHeaderProps) => {
  const { t } = useTranslation("patients")

  return (
    <PageHeader
      title={t("title")}
      actions={
        <Can I={Action.Create} a={Feature.Patient}>
          <Button onClick={onNewPatient}>
            <UserPlus />
            {t("newPatient")}
          </Button>
        </Can>
      }
    />
  )
}
