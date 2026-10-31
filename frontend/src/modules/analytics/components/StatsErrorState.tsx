import { useTranslation } from "react-i18next"
import { TriangleAlert } from "lucide-react"
import { Alert, AlertDescription, AlertTitle } from "../../../shared/components/ui/alert"
import { Button } from "../../../shared/components/ui/button"
import { PageContainer } from "../../../shared/components/ui/PageContainer"

interface StatsErrorStateProps {
  onRetry: () => void
  titleKey?: string
  descriptionKey?: string
  retryKey?: string
}

export const StatsErrorState = ({
  onRetry,
  titleKey = "errorTitle",
  descriptionKey = "errorDescription",
  retryKey = "retryButton",
}: StatsErrorStateProps) => {
  const { t } = useTranslation("analytics")

  return (
    <PageContainer className="flex items-center justify-center">
      <Alert variant="destructive" className="max-w-md">
        <TriangleAlert aria-hidden="true" />
        <AlertTitle>{t(titleKey)}</AlertTitle>
        <AlertDescription>
          <p>{t(descriptionKey)}</p>
          <Button type="button" variant="outline" size="sm" onClick={onRetry} className="mt-2">
            {t(retryKey)}
          </Button>
        </AlertDescription>
      </Alert>
    </PageContainer>
  )
}
