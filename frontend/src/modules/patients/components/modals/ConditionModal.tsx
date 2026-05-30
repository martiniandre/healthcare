import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Card } from "../../../../shared/components/ui/Card"
import { Input } from "../../../../shared/components/ui/Input"
import { Button } from "../../../../shared/components/ui/Button"
import { newConditionSchema, type NewConditionFormData } from "../../patient_schemas"

interface ConditionModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewConditionFormData) => void
  isPending: boolean
}

export const ConditionModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: ConditionModalProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewConditionFormData>({
    resolver: zodResolver(newConditionSchema),
  })

  if (!isOpen) {
    return null
  }

  const handleFormSubmit = (formData: NewConditionFormData) => {
    onSubmit({
      ...formData,
      icd10Code: formData.icd10Code.toUpperCase().trim(),
    })
    reset()
  }

  return (
    <div className="fixed inset-0 z-50 bg-black/30 backdrop-blur-[2px] flex items-center justify-center p-4">
      <Card glowingType="cyan" className="w-full max-w-[460px] p-8 relative">
        <h3 className="text-xl font-bold text-gray-900 mb-6 text-left">Registrar Comorbidade</h3>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="flex flex-col gap-4 text-left">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Código CID-10</label>
            <Input
              type="text"
              placeholder="Ex: I10, E11.9, J45"
              errorText={errors.icd10Code?.message}
              {...register("icd10Code")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Descrição do Diagnóstico</label>
            <Input
              type="text"
              placeholder="Ex: Hipertensão essencial, Diabetes mellitus tipo 2"
              errorText={errors.codeDisplay?.message}
              {...register("codeDisplay")}
            />
          </div>

          <div className="flex gap-3 justify-end mt-4">
            <Button variantType="outline" type="button" onClick={onClose}>
              Cancelar
            </Button>
            <Button type="submit" disabled={isPending}>
              Registrar Diagnóstico
            </Button>
          </div>
        </form>
      </Card>
    </div>
  )
}
