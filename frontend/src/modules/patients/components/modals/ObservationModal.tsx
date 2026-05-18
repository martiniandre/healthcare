import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Card } from "../../../../shared/components/ui/Card"
import { Input } from "../../../../shared/components/ui/Input"
import { Button } from "../../../../shared/components/ui/Button"
import { newObservationSchema, type NewObservationFormData } from "../../patient_schemas"

interface ObservationModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewObservationFormData) => void
  isPending: boolean
}

export const ObservationModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: ObservationModalProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<NewObservationFormData>({
    resolver: zodResolver(newObservationSchema),
    defaultValues: {
      loincCode: "8867-4",
    },
  })

  if (!isOpen) {
    return null
  }

  return (
    <div className="fixed inset-0 z-50 bg-black/30 backdrop-blur-[2px] flex items-center justify-center p-4">
      <Card glowingType="cyan" className="w-full max-w-[460px] p-8 relative">
        <h3 className="text-xl font-bold text-gray-900 mb-6 text-left">Adicionar Sinal Vital</h3>
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 text-left">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Sinal Vital (LOINC)</label>
            <select
              className="w-full bg-gray-50 border border-border rounded-lg px-3.5 py-2.5 text-sm text-gray-900 focus:outline-none focus:border-primary/50"
              {...register("loincCode")}
            >
              <option value="8867-4">Frequência Cardíaca (8867-4)</option>
              <option value="8310-5">Temperatura Corporal (8310-5)</option>
              <option value="85354-9">Pressão Arterial Sistólica (85354-9)</option>
            </select>
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Valor Quantitativo</label>
            <Input
              type="number"
              step="any"
              placeholder="Insira o valor numérico"
              errorText={errors.valueQuantity?.message}
              {...register("valueQuantity", { valueAsNumber: true })}
            />
          </div>

          <div className="flex gap-3 justify-end mt-4">
            <Button variantType="outline" type="button" onClick={onClose}>
              Cancelar
            </Button>
            <Button type="submit" disabled={isPending}>
              Gravar Métrica
            </Button>
          </div>
        </form>
      </Card>
    </div>
  )
}
