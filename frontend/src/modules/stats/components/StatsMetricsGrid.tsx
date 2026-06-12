import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { Users, CheckSquare, Clock, ArrowUpRight, Activity } from "lucide-react"

interface StatsMetricsGridProps {
  totalRegisteredPatients: number
  fhirComplianceRate: number
  averageServiceDurationMinutes: number
  activeConsultationsTotal: number
}

export const StatsMetricsGrid = ({
  totalRegisteredPatients,
  fhirComplianceRate,
  averageServiceDurationMinutes,
  activeConsultationsTotal
}: StatsMetricsGridProps) => {
  const { t: translate } = useTranslation()

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
      <Card className="p-4 flex items-center justify-between border border-border">
        <div className="text-left">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{translate("stats.metrics.activePatients")}</span>
          <span className="text-2xl font-black text-gray-900 mt-1 block">{totalRegisteredPatients}</span>
          <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
            <ArrowUpRight className="w-3.5 h-3.5" />
            {translate("stats.metrics.admissionsGrowth")}
          </span>
        </div>
        <div className="bg-primary/8 p-3 rounded-xl">
          <Users className="w-6 h-6 text-primary" />
        </div>
      </Card>

      <Card className="p-4 flex items-center justify-between border border-border">
        <div className="text-left">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{translate("stats.metrics.fhirCompliance")}</span>
          <span className="text-2xl font-black text-gray-900 mt-1 block">{fhirComplianceRate}%</span>
          <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
            <CheckSquare className="w-3.5 h-3.5" />
            {translate("stats.metrics.fhirCompliantDesc")}
          </span>
        </div>
        <div className="bg-emerald-50 p-3 rounded-xl">
          <CheckSquare className="w-6 h-6 text-emerald-600" />
        </div>
      </Card>

      <Card className="p-4 flex items-center justify-between border border-border">
        <div className="text-left">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{translate("stats.metrics.avgConsultationTime")}</span>
          <span className="text-2xl font-black text-gray-900 mt-1 block">{averageServiceDurationMinutes} min</span>
          <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
            <Clock className="w-3.5 h-3.5" />
            {translate("stats.metrics.consultationTrend")}
          </span>
        </div>
        <div className="bg-purple-50 p-3 rounded-xl">
          <Clock className="w-6 h-6 text-purple-600" />
        </div>
      </Card>

      <Card className="p-4 flex items-center justify-between border border-border">
        <div className="text-left">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">{translate("stats.metrics.weeklyConsultations")}</span>
          <span className="text-2xl font-black text-gray-900 mt-1 block">{activeConsultationsTotal}</span>
          <span className="text-[10px] text-amber-600 font-bold flex items-center gap-1 mt-1.5">
            <Activity className="w-3.5 h-3.5" />
            {translate("stats.metrics.activeIcuBeds")}
          </span>
        </div>
        <div className="bg-amber-50 p-3 rounded-xl">
          <Activity className="w-6 h-6 text-amber-500" />
        </div>
      </Card>
    </div>
  )
}
