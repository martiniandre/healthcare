import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Card } from "../../../../shared/components/ui/Card"
import { Input } from "../../../../shared/components/ui/Input"
import { Button } from "../../../../shared/components/ui/Button"
import { newAllergySchema, type NewAllergyFormData } from "../../patient_schemas"

interface AllergyModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewAllergyFormData) => void
  isPending: boolean
}

export const AllergyModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: AllergyModalProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewAllergyFormData>({
    resolver: zodResolver(newAllergySchema),
  })

  if (!isOpen) {
    return null
  }

  const handleFormSubmit = (formData: NewAllergyFormData) => {
    onSubmit(formData)
    reset()
  }

  return (
    <div className="fixed inset-0 z-50 bg-black/30 backdrop-blur-[2px] flex items-center justify-center p-4">
      <Card glowingType="cyan" className="w-full max-w-[460px] p-8 relative">
        <h3 className="text-xl font-bold text-gray-900 mb-6 text-left">Registrar Alergia</h3>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="flex flex-col gap-4 text-left">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Código do Alérgeno</label>
            <Input
              type="text"
              placeholder="Ex: 300916003, 716185002"
              errorText={errors.allergenCode?.message}
              {...register("allergenCode")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Descrição do Alérgeno</label>
            <Input
              type="text"
              placeholder="Ex: Alergia a Penicilina, Alergia a Amendoim"
              errorText={errors.allergenDisplay?.message}
              {...register("allergenDisplay")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">Reação Adversa Relatada</label>
            <Input
              type="text"
              placeholder="Ex: Urticária severa, Choque anafilático"
              errorText={errors.reaction?.message}
              {...register("reaction")}
            />
          </div>

          <div className="flex gap-3 justify-end mt-4">
            <Button variantType="outline" type="button" onClick={onClose}>
              Cancelar
            </Button>
            <Button type="submit" disabled={isPending}>
              Registrar Alergia
            </Button>
          </div>
        </form>
      </Card>
    </div>
  )
}
