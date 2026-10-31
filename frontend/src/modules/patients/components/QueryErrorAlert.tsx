import { AlertCircle, RotateCw } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Alert, AlertDescription, AlertTitle } from "../../../shared/components/ui/alert"
import { Button } from "../../../shared/components/ui/button"

interface QueryErrorAlertProps {
  onRetry: () => void
  title?: string
  description?: string
  isRetrying?: boolean
  className?: string
}

export const QueryErrorAlert = ({ onRetry, title, description, isRetrying, className }: QueryErrorAlertProps) => {
  const { t } = useTranslation("patients")

  return (
    <Alert variant="destructive" className={className}>
      <AlertCircle />
      <AlertTitle>{title ?? t("states.errorTitle")}</AlertTitle>
      <AlertDescription>
        <p>{description ?? t("states.errorDescription")}</p>
        <Button type="button" variant="outline" size="sm" onClick={onRetry} disabled={isRetrying} className="mt-2">
          <RotateCw />
          {t("states.retry")}
        </Button>
      </AlertDescription>
    </Alert>
  )
}
