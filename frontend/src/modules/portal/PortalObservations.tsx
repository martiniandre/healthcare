import { usePortalObservationsQuery } from "./queries"
import { Card } from "../../shared/components/ui/Card"
import { Loader2 } from "lucide-react"

export const PortalObservations = () => {
  const { data: observations, isLoading } = usePortalObservationsQuery()

  if (isLoading) {
    return (
      <Card className="flex items-center justify-center min-h-[300px]">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
      </Card>
    )
  }

  if (!observations || observations.length === 0) {
    return (
      <Card className="py-16 text-center">
        <p className="text-sm text-gray-500">Nenhum sinal vital registrado.</p>
      </Card>
    )
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      {observations.map((observation) => (
        <div
          key={observation.fhir_resource_id}
          className="bg-white border border-border rounded-xl p-5"
        >
          <p className="text-xs text-gray-500 font-medium mb-1">{observation.code_display}</p>
          <p className="text-2xl font-bold text-gray-900">
            {observation.value_quantity}{" "}
            <span className="text-sm font-normal text-gray-500">{observation.value_unit}</span>
          </p>
          <p className="text-xs text-gray-400 mt-2">
            {new Date(observation.observed_at).toLocaleDateString("pt-BR", {
              day: "numeric",
              month: "short",
              year: "numeric",
              hour: "2-digit",
              minute: "2-digit",
            })}
          </p>
        </div>
      ))}
    </div>
  )
}
