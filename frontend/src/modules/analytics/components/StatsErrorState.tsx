import { useTranslation } from "react-i18next"
import { TriangleAlert } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"
import { Card } from "../../../shared/components/ui/Card"
import { PageContainer } from "../../../shared/components/ui/PageContainer"

interface StatsErrorStateProps {
  onRetry: () => void
}

export const StatsErrorState = ({ onRetry }: StatsErrorStateProps) => {
  const { t } = useTranslation("analytics")

  return (
    <PageContainer className="flex items-center justify-center">
      <Card className="flex w-full max-w-md flex-col items-center gap-4 border-danger/20 p-8 text-center">
        <span className="grid size-14 place-items-center rounded-full bg-danger/10">
          <TriangleAlert className="size-7 text-danger" aria-hidden="true" />
        </span>
        <h3 className="font-display text-lg font-bold text-gray-900">{t("errorTitle")}</h3>
        <p className="text-xs leading-relaxed text-muted">{t("errorDescription")}</p>
        <Button variantType="danger" onClick={onRetry} className="mt-2 w-full">
          {t("retryButton")}
        </Button>
      </Card>
    </PageContainer>
  )
}
