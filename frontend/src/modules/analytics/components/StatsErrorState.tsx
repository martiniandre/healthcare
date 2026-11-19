import { useTranslation } from "react-i18next"
import { QueryErrorAlert } from "@/shared/components/ui/QueryErrorAlert"
import { PageContainer } from "@/shared/components/ui/PageContainer"

interface StatsErrorStateProps {
  onRetry: () => void
  titleKey?: string
  descriptionKey?: string
}

export const StatsErrorState = ({
  onRetry,
  titleKey = "errorTitle",
  descriptionKey = "errorDescription",
}: StatsErrorStateProps) => {
  const { t } = useTranslation("analytics")

  return (
    <PageContainer className="flex items-center justify-center">
      <QueryErrorAlert className="max-w-md" title={t(titleKey)} description={t(descriptionKey)} onRetry={onRetry} />
    </PageContainer>
  )
}
