import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Send } from "lucide-react"
import { Card } from "../../../../shared/components/ui/Card"
import { Input } from "../../../../shared/components/ui/Input"
import { Button } from "../../../../shared/components/ui/Button"
import { newReportSchema, type NewReportFormData } from "../../patient_schemas"

interface ReportModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewReportFormData) => void
  isPending: boolean
}

export const ReportModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: ReportModalProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<NewReportFormData>({
    resolver: zodResolver(newReportSchema),
  })

  if (!isOpen) {
    return null
  }

  return (
    <div className="fixed inset-0 z-50 bg-black/30 backdrop-blur-[2px] flex items-center justify-center p-4">
      <Card glowingType="cyan" className="w-full max-w-[500px] p-8 relative">
        <h3 className="text-xl font-bold text-gray-900 mb-6 text-left">Registrar Laudo Diagnóstico</h3>
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 text-left">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Nome do Exame / Laudo</label>
            <Input
              type="text"
              placeholder="Ex: Hemograma Completo, ECG"
              errorText={errors.reportDisplay?.message}
              {...register("reportDisplay")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Conclusão Médica</label>
            <textarea
              placeholder="Redija a conclusão clínica circunstanciada..."
              className="w-full bg-gray-50 border border-border rounded-lg px-3.5 py-2.5 h-28 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-primary/50"
              {...register("conclusion")}
            />
            {errors.conclusion?.message && (
              <span className="text-xs text-red-500 font-medium px-1 mt-1">
                {errors.conclusion.message}
              </span>
            )}
          </div>

          <div className="flex gap-3 justify-end mt-4">
            <Button variantType="outline" type="button" onClick={onClose}>
              Cancelar
            </Button>
            <Button type="submit" disabled={isPending} className="gap-2">
              <Send className="w-3.5 h-3.5" />
              Assinar Laudo
            </Button>
          </div>
        </form>
      </Card>
    </div>
  )
}
