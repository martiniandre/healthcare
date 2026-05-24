import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Card } from "../../../../shared/components/ui/Card"
import { Input } from "../../../../shared/components/ui/Input"
import { Button } from "../../../../shared/components/ui/Button"
import { newMedicationSchema, type NewMedicationFormData } from "../../patient_schemas"

interface MedicationModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewMedicationFormData) => void
  isPending: boolean
}

export const MedicationModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: MedicationModalProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewMedicationFormData>({
    resolver: zodResolver(newMedicationSchema),
  })

  if (!isOpen) {
    return null
  }

  const handleFormSubmit = (formData: NewMedicationFormData) => {
    onSubmit(formData)
    reset()
  }

  return (
    <div className="fixed inset-0 z-50 bg-black/30 backdrop-blur-[2px] flex items-center justify-center p-4">
      <Card glowingType="amethyst" className="w-full max-w-[460px] p-8 relative">
        <h3 className="text-xl font-bold text-gray-900 mb-6 text-left">Nova Prescrição</h3>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="flex flex-col gap-4 text-left">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Medicação</label>
            <Input
              type="text"
              placeholder="Ex: Dipirona 500mg, Amoxicilina 875mg"
              errorText={errors.medicationDisplay?.message}
              {...register("medicationDisplay")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Instrução de Dosagem</label>
            <textarea
              className="w-full h-24 px-3 py-2 text-sm border border-border rounded-lg outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all resize-none"
              placeholder="Ex: Tomar 1 comprimido de 8 em 8 horas por 5 dias."
              {...register("dosageInstruction")}
            />
            {errors.dosageInstruction?.message && (
              <span className="text-xs text-rose-500 font-medium">
                {errors.dosageInstruction.message}
              </span>
            )}
          </div>

          <div className="flex gap-3 justify-end mt-4">
            <Button variantType="outline" type="button" onClick={onClose}>
              Cancelar
            </Button>
            <Button type="submit" disabled={isPending}>
              Prescrever
            </Button>
          </div>
        </form>
      </Card>
    </div>
  )
}
