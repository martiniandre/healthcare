import { usePortalImagingQuery } from "./queries"
import { Card } from "../../shared/components/ui/Card"
import { Loader2 } from "lucide-react"

export const PortalImaging = () => {
  const { data: imagingStudies, isLoading } = usePortalImagingQuery()

  if (isLoading) {
    return (
      <Card className="flex items-center justify-center min-h-[300px]">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
      </Card>
    )
  }

  if (!imagingStudies || imagingStudies.length === 0) {
    return (
      <Card className="py-16 text-center">
        <p className="text-sm text-gray-500">Nenhum exame de imagem encontrado.</p>
      </Card>
    )
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      {imagingStudies.map((study) => (
        <div
          key={study.fhir_resource_id}
          className="bg-white border border-border rounded-xl p-5"
        >
          <div className="flex items-start justify-between mb-3">
            <span className="text-xs font-bold px-2 py-0.5 rounded bg-gray-100 text-gray-600 uppercase">
              {study.modality || "N/I"}
            </span>
            <span className="text-xs text-gray-400">
              {study.created_at
                ? new Date(study.created_at).toLocaleDateString("pt-BR")
                : ""}
            </span>
          </div>
          <p className="text-sm font-bold text-gray-900 mb-2">
            {study.title || "Estudo de Imagem"}
          </p>
          <span className="text-xs font-bold px-2.5 py-1 rounded-full capitalize bg-blue-100 text-blue-700">
            {study.status}
          </span>
        </div>
      ))}
    </div>
  )
}
