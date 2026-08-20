import { useTranslation } from "react-i18next"
import { Activity } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"
import { PageContainer } from "../../../shared/components/ui/PageContainer"

interface StatsErrorStateProps {
  onRetry: () => void
}

export const StatsErrorState = ({ onRetry }: StatsErrorStateProps) => {
  const { t: translate } = useTranslation("analytics")

  return (
    <PageContainer className="flex items-center justify-center gap-4 select-none">
      <div className="text-center p-8 bg-white border border-red-100 shadow-xl rounded-2xl max-w-md w-full flex flex-col items-center gap-4">
        <div className="bg-red-50 p-4 rounded-full">
          <Activity className="w-10 h-10 text-red-500 animate-bounce" />
        </div>
        <h3 className="text-lg font-display font-bold text-gray-900">{translate("errorTitle")}</h3>
        <p className="text-xs text-gray-500 leading-relaxed">
          {translate("errorDescription")}
        </p>
        <Button
          onClick={onRetry}
          className="w-full bg-red-600 hover:bg-red-700 text-white font-bold py-2 rounded-xl transition-all duration-200 mt-2"
        >
          {translate("retryButton")}
        </Button>
      </div>
    </PageContainer>
  )
}
