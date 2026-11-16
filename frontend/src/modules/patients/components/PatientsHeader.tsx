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
  const { t: translateSidebar } = useTranslation("sidebar")

  return (
    <PageHeader
      eyebrow={translateSidebar("topics.clinical")}
      title={t("title")}
      description={t("subtitle")}
      actions={
        <Can I={Action.Create} a={Feature.Patient}>
          <Button onClick={onNewPatient} size="lg">
            <UserPlus />
            {t("newPatient")}
          </Button>
        </Can>
      }
    />
  )
}
