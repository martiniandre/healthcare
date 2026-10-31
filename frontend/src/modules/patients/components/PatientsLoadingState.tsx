import { useTranslation } from "react-i18next"
import { Spinner } from "../../../shared/components/ui/spinner"

export const PatientsLoadingState = () => {
  const { t } = useTranslation("patients")

  return (
    <div className="flex-1 flex items-center justify-center py-20">
      <div className="flex items-center gap-3 text-muted-foreground">
        <Spinner className="size-5" />
        <span className="text-sm font-medium">{t("loading")}</span>
      </div>
    </div>
  )
}
