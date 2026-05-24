import { Pill, Plus, ShieldAlert, CheckCircle } from "lucide-react"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"
import type { MedicationRequest } from "../types"

interface ClinicalMedicationsProps {
  medications: MedicationRequest[]
  onAdd: () => void
}

export const ClinicalMedications = ({ medications, onAdd }: ClinicalMedicationsProps) => {
  return (
    <Card className="flex flex-col gap-5 min-h-[450px]">
      <div className="flex items-center justify-between border-b border-border pb-4">
        <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2">
          <Pill className="w-4 h-4 text-purple-500" />
          Prescrições (MedicationRequest - FHIR)
        </h3>
        <Button onClick={onAdd} className="px-3 py-2 text-xs">
          <Plus className="w-3.5 h-3.5" />
          Nova Prescrição
        </Button>
      </div>

      {medications.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center gap-2 py-16">
          <ShieldAlert className="w-8 h-8 text-gray-300" />
          <span className="text-xs text-muted">Nenhuma medicação prescrita para este atendimento</span>
        </div>
      ) : (
        <div className="overflow-x-auto w-full">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-border bg-gray-50/80">
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Medicação</th>
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Instrução de Dosagem</th>
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Status</th>
                <th className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">Data de Prescrição</th>
              </tr>
            </thead>
            <tbody>
              {medications.map((medicationItem) => (
                <tr key={medicationItem.fhir_id} className="border-b border-border/60 hover:bg-gray-50 transition-colors duration-300">
                  <td className="py-4 px-4 align-top">
                    <span className="text-sm font-extrabold text-gray-900 block">
                      {medicationItem.medication_display}
                    </span>
                  </td>
                  <td className="py-4 px-4 align-top">
                    <span className="text-sm font-bold text-gray-800 block whitespace-pre-line">
                      {medicationItem.dosage_instruction}
                    </span>
                  </td>
                  <td className="py-4 px-4 align-top">
                    <span className="text-[9px] bg-purple-50 border border-purple-100 text-purple-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
                      <CheckCircle className="w-3 h-3" />
                      {medicationItem.status}
                    </span>
                  </td>
                  <td className="py-4 px-4 align-top">
                    <span className="text-xs text-gray-500 font-semibold block">
                      {new Date(medicationItem.created_at).toLocaleString()}
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
