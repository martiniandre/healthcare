import { ArrowLeft } from "lucide-react"
import { useTranslation } from "react-i18next"
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
  const { t } = useTranslation("patients")

  return (
    <div className="flex flex-col sm:flex-row sm:items-center gap-4 text-left">
      <Button variantType="outline" onClick={onBack} className="px-3 self-start sm:self-auto gap-2">
        <ArrowLeft className="w-4 h-4" />
        {t("details.back")}
      </Button>
      <div>
        <h2 className="text-xl font-black text-gray-900 leading-none">
          {patient.full_name}
        </h2>
        <span className="text-xs text-muted mt-2.5 flex flex-wrap items-center gap-x-2.5 gap-y-1.5">
          <span className="font-semibold text-gray-700">
            {t("details.cpf")}: {patient.document_id}
          </span>
          <span className="hidden sm:inline text-gray-300">•</span>
          <span className="text-gray-500">
            {t("details.birth")}: {patient.birth_date}
          </span>
          <span className="hidden sm:inline text-gray-300">•</span>
          <span className="font-mono text-[10px] text-gray-400 bg-gray-50 border border-border/80 px-2 py-0.5 rounded-md">
            {t("details.fhirId")}: {patient.fhir_resource_id}
          </span>
        </span>
      </div>
    </div>
  )
}
