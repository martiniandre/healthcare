import { CheckCircle, Eye } from "lucide-react"
import { Card } from "../../../shared/components/ui/Card"
import type { ImagingStudy } from "../types"

interface ImagingStudyDetailsProperties {
  study: ImagingStudy
}

export const ImagingStudyDetails = ({ study }: ImagingStudyDetailsProperties) => {
  return (
    <Card className="flex flex-col gap-5 lg:col-span-1 text-left">
      <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2 border-b border-border pb-4">
        <Eye className="w-4 h-4 text-primary animate-pulse-glow" />
        Detalhes do Estudo
      </h3>

      <div className="flex flex-col gap-4">
        <div className="bg-gray-50 border border-border p-3.5 rounded-xl">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">ID do Estudo</span>
          <span className="text-xs font-bold text-gray-800 mt-1 block">{study.id}</span>
        </div>

        <div className="bg-gray-50 border border-border p-3.5 rounded-xl">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Modalidade Clínica</span>
          <span className="text-xs font-bold text-gray-800 mt-1 block uppercase">{study.modality}</span>
        </div>

        <div className="bg-gray-50 border border-border p-3.5 rounded-xl">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Status do Barramento</span>
          <span className="text-xs font-bold text-emerald-600 mt-1 flex items-center gap-1.5">
            <CheckCircle className="w-4 h-4" />
            {study.status}
          </span>
        </div>

        <div className="bg-gray-50 border border-border p-3.5 rounded-xl">
          <span className="text-[10px] text-gray-500 font-bold uppercase tracking-wider block">Série / Fatias</span>
          <span className="text-xs font-bold text-gray-800 mt-1 block">Slice #1 (Visualização Ativa)</span>
        </div>
      </div>
    </Card>
  )
}
