import { CircleAlert } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Alert, AlertDescription } from "../../../../shared/components/ui/alert"

export const FormSubmitError = () => {
  const { t } = useTranslation("patients")

  return (
    <Alert variant="destructive" data-testid="form-submit-error">
      <CircleAlert aria-hidden="true" />
      <AlertDescription>{t("form.submitError")}</AlertDescription>
    </Alert>
  )
}
