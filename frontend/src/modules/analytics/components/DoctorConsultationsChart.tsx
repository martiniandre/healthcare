import { Card } from "../../../shared/components/ui/Card"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import type { DoctorConsultation } from "../dashboard_types"

interface DoctorConsultationsChartProps {
  consultationsPerDoctor: DoctorConsultation[]
}

export const DoctorConsultationsChart = ({ consultationsPerDoctor }: DoctorConsultationsChartProps) => {
  const maxCount = consultationsPerDoctor.length > 0
    ? Math.max(...consultationsPerDoctor.map((item) => item.count))
    : 1

  return (
    <Card className="p-5 flex flex-col gap-5 text-left border border-border">
      <div>
        <h3 className="font-extrabold text-gray-900 text-md">Consultas por Médico</h3>
        <span className="text-xs text-muted block mt-1">Distribuição de consultas nos últimos 30 dias</span>
      </div>

      <div className="overflow-x-auto w-full">
        {consultationsPerDoctor.length > 0 ? (
          <div className="flex flex-col gap-3 min-w-[300px]">
            {consultationsPerDoctor.map((doctorConsultation) => {
              const barWidth = (doctorConsultation.count / maxCount) * 100
              return (
                <div key={doctorConsultation.doctor_name} className="flex items-center gap-3">
                  <div className="w-32 shrink-0 text-left">
                    <span className="text-xs font-bold text-gray-700 block truncate">{doctorConsultation.doctor_name}</span>
                    <span className="text-[10px] text-gray-500 font-semibold">{doctorConsultation.specialty}</span>
                  </div>
                  <div className="flex-1 flex items-center gap-2">
                    <div className="flex-1 bg-gray-100 rounded-full h-5 overflow-hidden">
                      <div
                        className="bg-primary rounded-full h-full transition-all duration-500"
                        style={{ width: `${barWidth}%` }}
                      />
                    </div>
                    <span className="text-xs font-black text-gray-900 w-8 text-right">{doctorConsultation.count}</span>
                  </div>
                </div>
              )
            })}
          </div>
        ) : (
          <EmptyState
            title="Sem dados"
            description="Nenhum dado de consultas por médico disponível para o período."
            className="h-40"
          />
        )}
      </div>
    </Card>
  )
}
