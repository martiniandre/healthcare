import { Heart, Thermometer, Activity, Plus } from "lucide-react"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"

interface ObservationRepresentation {
  fhir_id: string
  encounter_fhir_id: string
  patient_fhir_id: string
  loinc_code: string
  code_display: string
  value_quantity: number
  value_unit: string
  created_at: string
}

interface VitalSignsProps {
  observations: ObservationRepresentation[]
  onAdd: () => void
}

export const VitalSigns = ({ observations, onAdd }: VitalSignsProps) => {
  return (
    <Card className="flex flex-col gap-5 min-h-[450px]">
      <div className="flex items-center justify-between border-b border-border pb-4">
        <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2">
          <Heart className="w-4 h-4 text-red-500 animate-pulse-glow" />
          Sinais Vitais (Observations)
        </h3>
        <Button onClick={onAdd} className="px-3 py-2 text-xs">
          <Plus className="w-3.5 h-3.5" />
          Adicionar Sinal
        </Button>
      </div>

      {observations.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center gap-2 py-16">
          <Heart className="w-8 h-8 text-gray-300" />
          <span className="text-xs text-muted">Nenhuma medição registrada para esta consulta</span>
        </div>
      ) : (
        <div className="overflow-x-auto w-full">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-border bg-gray-50/80">
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Métrica de Aferição</th>
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Código LOINC</th>
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Valor Clínico</th>
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Data da Aferição</th>
              </tr>
            </thead>
            <tbody>
              {observations.map((observation) => {
                const isHeartRate = observation.loinc_code === "8867-4"
                const isTemp = observation.loinc_code === "8310-5"
                return (
                  <tr key={observation.fhir_id} className="border-b border-border/60 hover:bg-gray-50 transition-colors duration-300">
                    <td className="py-4 px-4">
                      <div className="flex items-center gap-3">
                        <div className={`p-2 rounded-lg border ${
                          isHeartRate 
                            ? "bg-red-50 border-red-100 text-red-600"
                            : isTemp 
                              ? "bg-amber-50 border-amber-100 text-amber-600"
                              : "bg-blue-50 border-blue-100 text-blue-600"
                        }`}>
                          {isHeartRate ? (
                            <Heart className="w-4 h-4" />
                          ) : isTemp ? (
                            <Thermometer className="w-4 h-4" />
                          ) : (
                            <Activity className="w-4 h-4" />
                          )}
                        </div>
                        <span className="text-sm font-bold text-gray-800 block">
                          {observation.code_display}
                        </span>
                      </div>
                    </td>
                    <td className="py-4 px-4">
                      <span className="text-xs font-mono text-gray-500">
                        {observation.loinc_code}
                      </span>
                    </td>
                    <td className="py-4 px-4">
                      <span className="text-sm font-extrabold text-gray-800">
                        {observation.value_quantity}
                        <span className="text-xs text-muted font-normal ml-1">
                          {observation.value_unit}
                        </span>
                      </span>
                    </td>
                    <td className="py-4 px-4">
                      <span className="text-xs text-gray-500 font-semibold">
                        {new Date(observation.created_at).toLocaleString()}
                      </span>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}
    </Card>
  )
}
