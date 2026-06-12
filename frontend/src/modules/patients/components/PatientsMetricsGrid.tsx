import { useTranslation } from "react-i18next"
import { Users, Activity, Database } from "lucide-react"

interface PatientsMetricsGridProps {
  totalPatients: number
}

export const PatientsMetricsGrid = ({ totalPatients }: PatientsMetricsGridProps) => {
  const { t } = useTranslation()

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div className="bg-white border border-border rounded-xl p-4 flex items-center gap-4">
        <div className="w-10 h-10 rounded-lg bg-primary/8 flex items-center justify-center shrink-0">
          <Users className="w-5 h-5 text-primary" />
        </div>
        <div>
          <span className="text-[10px] text-muted font-semibold uppercase tracking-wider block">
            {t("patients.cards.total")}
          </span>
          <span className="text-2xl font-black text-gray-900 leading-none mt-0.5 block">
            {totalPatients}
          </span>
        </div>
      </div>

      <div className="bg-white border border-border rounded-xl p-4 flex items-center gap-4">
        <div className="w-10 h-10 rounded-lg bg-secondary/8 flex items-center justify-center shrink-0">
          <Activity className="w-5 h-5 text-secondary" />
        </div>
        <div>
          <span className="text-[10px] text-muted font-semibold uppercase tracking-wider block">
            {t("patients.cards.standard")}
          </span>
          <span className="text-sm font-bold text-gray-800 mt-0.5 block">FHIR R4 Compliant</span>
        </div>
      </div>

      <div className="bg-white border border-border rounded-xl p-4 flex items-center gap-4">
        <div className="w-10 h-10 rounded-lg bg-blue-50 flex items-center justify-center shrink-0">
          <Database className="w-5 h-5 text-blue-500" />
        </div>
        <div>
          <span className="text-[10px] text-muted font-semibold uppercase tracking-wider block">
            {t("patients.cards.integration")}
          </span>
          <span className="text-sm font-bold text-gray-800 mt-0.5 block">Cloud Healthcare API</span>
        </div>
      </div>
    </div>
  )
}
