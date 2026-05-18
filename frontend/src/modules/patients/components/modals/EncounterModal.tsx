import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Card } from "../../../../shared/components/ui/Card"
import { Input } from "../../../../shared/components/ui/Input"
import { Button } from "../../../../shared/components/ui/Button"
import { newEncounterSchema, type NewEncounterFormData } from "../../patient_schemas"

interface EncounterModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewEncounterFormData) => void
  isPending: boolean
}

export const EncounterModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: EncounterModalProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<NewEncounterFormData>({
    resolver: zodResolver(newEncounterSchema),
  })

  if (!isOpen) {
    return null
  }

  return (
    <div className="fixed inset-0 z-50 bg-black/30 backdrop-blur-[2px] flex items-center justify-center p-4">
      <Card glowingType="cyan" className="w-full max-w-[460px] p-8 relative">
        <h3 className="text-xl font-bold text-gray-900 mb-6 text-left">Iniciar Nova Consulta</h3>
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 text-left">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Motivo / Queixa Principal</label>
            <Input
              type="text"
              placeholder="Ex: Consulta de Rotina, Dor Abdominal"
              errorText={errors.reasonDisplay?.message}
              {...register("reasonDisplay")}
            />
          </div>
          <div className="flex gap-3 justify-end mt-4">
            <Button variantType="outline" type="button" onClick={onClose}>
              Cancelar
            </Button>
            <Button type="submit" disabled={isPending}>
              Registrar Consulta
            </Button>
          </div>
        </form>
      </Card>
    </div>
  )
}
