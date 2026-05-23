import { Activity, Plus, ShieldAlert, CheckCircle } from "lucide-react"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"
import type { Condition } from "../types"

interface ClinicalConditionsProps {
  conditions: Condition[]
  onAdd: () => void
}

export const ClinicalConditions = ({ conditions, onAdd }: ClinicalConditionsProps) => {
  return (
    <Card className="flex flex-col gap-5 min-h-[450px]">
      <div className="flex items-center justify-between border-b border-border pb-4">
        <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2">
          <Activity className="w-4 h-4 text-rose-500" />
          Diagnósticos Ativos & Comorbidades (Conditions - FHIR)
        </h3>
        <Button onClick={onAdd} className="px-3 py-2 text-xs">
          <Plus className="w-3.5 h-3.5" />
          Adicionar Diagnóstico
        </Button>
      </div>

      {conditions.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center gap-2 py-16">
          <ShieldAlert className="w-8 h-8 text-gray-300" />
          <span className="text-xs text-muted">Nenhum diagnóstico ativo registrado para este paciente</span>
        </div>
      ) : (
        <div className="overflow-x-auto w-full">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-border bg-gray-50/80">
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Código CID-10</th>
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Descrição Clínico-Científica</th>
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Status Clínico</th>
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Data de Registro</th>
              </tr>
            </thead>
            <tbody>
              {conditions.map((conditionItem) => (
                <tr key={conditionItem.fhir_id} className="border-b border-border/60 hover:bg-gray-50 transition-colors duration-300">
                  <td className="py-4 px-4 align-top">
                    <span className="text-sm font-extrabold text-gray-900 bg-rose-50 border border-rose-100 px-2 py-1 rounded-md">
                      {conditionItem.icd10_code}
                    </span>
                  </td>
                  <td className="py-4 px-4 align-top">
                    <span className="text-sm font-bold text-gray-800 block">
                      {conditionItem.code_display}
                    </span>
                  </td>
                  <td className="py-4 px-4 align-top">
                    <span className="text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
                      <CheckCircle className="w-3 h-3" />
                      {conditionItem.clinical_status}
                    </span>
                  </td>
                  <td className="py-4 px-4 align-top">
                    <span className="text-xs text-gray-500 font-semibold block">
                      {new Date(conditionItem.created_at).toLocaleString()}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Card>
  )
}
