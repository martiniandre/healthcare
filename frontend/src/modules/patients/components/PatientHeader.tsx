import { ArrowLeft } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"

interface PatientRepresentation {
  patient_id: string
  fhir_resource_id: string
  full_name: string
  birth_date: string
  document_id: string
  phone_number: string
}

interface PatientHeaderProps {
  patient: PatientRepresentation
  onBack: () => void
}

export const PatientHeader = ({ patient, onBack }: PatientHeaderProps) => {
  return (
    <div className="flex items-center gap-4">
      <Button variantType="outline" onClick={onBack} className="px-3">
        <ArrowLeft className="w-4 h-4" />
        Voltar
      </Button>
      <div className="text-left">
        <h2 className="text-xl font-black text-gray-900 leading-none">
          {patient.full_name}
        </h2>
        <span className="text-xs text-muted mt-1.5 block">
          CPF: {patient.document_id} • Nascimento: {patient.birth_date} • ID FHIR: {patient.fhir_resource_id}
        </span>
      </div>
    </div>
  )
}
