import { usePortalConditionsQuery } from "./queries"
import { Card } from "../../shared/components/ui/Card"
import { Loader2 } from "lucide-react"

const statusBadgeClass = (status: string) => {
  switch (status) {
    case "active":
      return "bg-red-100 text-red-700"
    case "recurrence":
      return "bg-orange-100 text-orange-700"
    case "relapse":
      return "bg-amber-100 text-amber-700"
    case "inactive":
    case "remission":
      return "bg-yellow-100 text-yellow-700"
    case "resolved":
      return "bg-green-100 text-green-700"
    default:
      return "bg-gray-100 text-gray-700"
  }
}

export const PortalConditions = () => {
  const { data: conditions, isLoading } = usePortalConditionsQuery()

  if (isLoading) {
    return (
      <Card className="flex items-center justify-center min-h-[300px]">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
      </Card>
    )
  }

  if (!conditions || conditions.length === 0) {
    return (
      <Card className="py-16 text-center">
        <p className="text-sm text-gray-500">Nenhuma condição registrada.</p>
      </Card>
    )
  }

  return (
    <div className="flex flex-col gap-3">
      {conditions.map((condition) => (
        <div
          key={condition.fhir_resource_id}
          className="bg-white border border-border rounded-xl p-5 flex items-start justify-between"
        >
          <div>
            <p className="text-sm font-bold text-gray-900">{condition.code_display}</p>
            <p className="text-xs text-gray-500 mt-0.5">
              {condition.icd10_code && <span className="font-mono">{condition.icd10_code} | </span>}
              Início: {condition.onset_at ? new Date(condition.onset_at).toLocaleDateString("pt-BR") : "N/I"}
            </p>
          </div>
          <span
            className={`text-xs font-bold px-2.5 py-1 rounded-full capitalize ${statusBadgeClass(condition.clinical_status)}`}
          >
            {condition.clinical_status}
          </span>
        </div>
      ))}
    </div>
  )
}
