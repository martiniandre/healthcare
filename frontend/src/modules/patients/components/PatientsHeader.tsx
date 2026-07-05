import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { UserPlus } from "lucide-react"

interface PatientsHeaderProps {
  onNewPatient: () => void
}

export const PatientsHeader = ({ onNewPatient }: PatientsHeaderProps) => {
  const { t } = useTranslation()

  return (
    <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 className="text-xl font-black text-gray-900 tracking-tight leading-none">
          {t("patients.title")}
        </h2>
        <span className="text-xs text-muted mt-1 block">
          {t("patients.subtitle")}
        </span>
      </div>
      <Can I={Action.Create} a={Feature.Patient}>
        <Button onClick={onNewPatient} className="py-2 px-4 self-start sm:self-auto gap-2">
          <UserPlus className="w-4 h-4" />
          {t("patients.newPatient")}
        </Button>
      </Can>
    </div>
  )
}
