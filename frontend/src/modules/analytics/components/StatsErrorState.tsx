import { useTranslation } from "react-i18next"
import { Activity } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"

export const StatsErrorState = () => {
  const { t: translate } = useTranslation()

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col items-center justify-center gap-4 max-w-7xl mx-auto w-full select-none">
      <div className="text-center p-8 bg-white border border-red-100 shadow-xl rounded-2xl max-w-md w-full flex flex-col items-center gap-4">
        <div className="bg-red-50 p-4 rounded-full">
          <Activity className="w-10 h-10 text-red-500 animate-bounce" />
        </div>
        <h3 className="text-lg font-black text-gray-900">{translate("analytics.errorTitle") || "Erro ao carregar dados"}</h3>
        <p className="text-xs text-gray-500 leading-relaxed">
          {translate("analytics.errorDescription") || "Não foi possível estabelecer conexão com o serviço de analytics FHIR."}
        </p>
        <Button 
          onClick={() => window.location.reload()} 
          className="w-full bg-red-600 hover:bg-red-700 text-white font-bold py-2 rounded-xl transition-all duration-200 mt-2"
        >
          {translate("analytics.retryButton") || "Tentar Novamente"}
        </Button>
      </div>
    </div>
  )
}
