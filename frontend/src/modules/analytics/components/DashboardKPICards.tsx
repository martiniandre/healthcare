import { Card } from "../../../shared/components/ui/Card"
import { Users, Clock, Activity, ArrowUpRight, Stethoscope, FileText } from "lucide-react"

interface DashboardKPICardsProps {
  consultationsToday: number
  consultationsTrend: string
  occupancyRate: number
  avgWaitTimeMinutes: number
  activePatients: number
  examsToday: number
  newDiagnosesToday: number
}

export const DashboardKPICards = ({
  consultationsToday,
  consultationsTrend,
  occupancyRate,
  avgWaitTimeMinutes,
  activePatients,
  examsToday,
  newDiagnosesToday,
}: DashboardKPICardsProps) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4">
      <Card className="p-4 flex items-center justify-between border border-border">
        <div className="text-left">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Consultas hoje</span>
          <span className="text-2xl font-black text-gray-900 mt-1 block">{consultationsToday}</span>
          <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
            <ArrowUpRight className="w-3.5 h-3.5" />
            {consultationsTrend}
          </span>
        </div>
        <div className="bg-primary/8 p-3 rounded-xl">
          <Stethoscope className="w-6 h-6 text-primary" />
        </div>
      </Card>

      <Card className="p-4 flex flex-col items-start justify-between border border-border">
        <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider">Taxa de ocupação</span>
        <div className="w-full mt-2">
          <span className="text-2xl font-black text-gray-900">{occupancyRate.toFixed(1)}%</span>
          <div className="w-full bg-gray-200 rounded-full h-2 mt-2">
            <div
              className="bg-primary rounded-full h-2 transition-all duration-500"
              style={{ width: `${Math.min(occupancyRate, 100)}%` }}
            />
          </div>
        </div>
      </Card>

      <Card className="p-4 flex items-center justify-between border border-border">
        <div className="text-left">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Tempo médio de espera</span>
          <span className="text-2xl font-black text-gray-900 mt-1 block">{avgWaitTimeMinutes.toFixed(0)} min</span>
          <span className="text-[10px] text-amber-600 font-bold flex items-center gap-1 mt-1.5">
            <Clock className="w-3.5 h-3.5" />
            Por departamento
          </span>
        </div>
        <div className="bg-amber-50 p-3 rounded-xl">
          <Clock className="w-6 h-6 text-amber-500" />
        </div>
      </Card>

      <Card className="p-4 flex items-center justify-between border border-border">
        <div className="text-left">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Pacientes ativos</span>
          <span className="text-2xl font-black text-gray-900 mt-1 block">{activePatients}</span>
          <span className="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1.5">
            <Users className="w-3.5 h-3.5" />
            Em atendimento
          </span>
        </div>
        <div className="bg-emerald-50 p-3 rounded-xl">
          <Users className="w-6 h-6 text-emerald-600" />
        </div>
      </Card>

      <Card className="p-4 flex items-center justify-between border border-border">
        <div className="text-left">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Exames hoje</span>
          <span className="text-2xl font-black text-gray-900 mt-1 block">{examsToday}</span>
          <span className="text-[10px] text-purple-600 font-bold flex items-center gap-1 mt-1.5">
            <Activity className="w-3.5 h-3.5" />
            Realizados
          </span>
        </div>
        <div className="bg-purple-50 p-3 rounded-xl">
          <Activity className="w-6 h-6 text-purple-600" />
        </div>
      </Card>

      <Card className="p-4 flex items-center justify-between border border-border">
        <div className="text-left">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Novos diagnósticos</span>
          <span className="text-2xl font-black text-gray-900 mt-1 block">{newDiagnosesToday}</span>
          <span className="text-[10px] text-sky-600 font-bold flex items-center gap-1 mt-1.5">
            <FileText className="w-3.5 h-3.5" />
            Últimos 30 dias
          </span>
        </div>
        <div className="bg-sky-50 p-3 rounded-xl">
          <FileText className="w-6 h-6 text-sky-600" />
        </div>
      </Card>
    </div>
  )
}
