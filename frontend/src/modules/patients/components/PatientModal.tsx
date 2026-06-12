import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../shared/components/ui/Dialog"
import { MaskedInput } from "../../../shared/components/ui/MaskedInput"
import { Button } from "../../../shared/components/ui/Button"
import { Input } from "../../../shared/components/ui/Input"
import { getNewPatientSchema, type NewPatientFormData } from "../patient_schemas"
import { useCreatePatientMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import { useEffect } from "react"

interface PatientModalProps {
  isOpen: boolean
  onOpenChange: (open: boolean) => void
}

export const PatientModal = ({ isOpen, onOpenChange }: PatientModalProps) => {
  const { t } = useTranslation()
  const createPatientMutation = useCreatePatientMutation()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<NewPatientFormData>({
    resolver: zodResolver(getNewPatientSchema(t)),
    defaultValues: {
      fullName: "",
      birthDate: "",
      documentId: "",
      phoneNumber: "",
    },
  })

  useEffect(() => {
    if (isOpen) {
      reset()
    }
  }, [isOpen, reset])

  const onSubmit = async (formData: NewPatientFormData) => {
    try {
      await createPatientMutation.mutateAsync({
        full_name: formData.fullName,
        birth_date: formData.birthDate,
        document_id: formData.documentId,
        phone_number: formData.phoneNumber,
      })
      onOpenChange(false)
      toast.success(t("patients.toast.createSuccess"))
    } catch {
      toast.error(t("patients.toast.createError"))
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle className="text-left">{t("patients.modal.title")}</DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 text-left mt-4" noValidate>
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">{t("patients.modal.name")}</label>
            <Input
              type="text"
              placeholder={t("patients.modal.namePlaceholder")}
              errorText={errors.fullName?.message}
              {...register("fullName")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">{t("patients.modal.birthDate")}</label>
            <MaskedInput
              mask="9999-99-99"
              placeholder={t("patients.modal.birthDatePlaceholder")}
              errorText={errors.birthDate?.message}
              {...register("birthDate")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">{t("patients.modal.document")}</label>
            <MaskedInput
              mask="999.999.999-99"
              placeholder={t("patients.modal.documentPlaceholder")}
              errorText={errors.documentId?.message}
              {...register("documentId")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">{t("patients.modal.phone")}</label>
            <MaskedInput
              mask="(99) 99999-9999"
              placeholder={t("patients.modal.phonePlaceholder")}
              errorText={errors.phoneNumber?.message}
              {...register("phoneNumber")}
            />
          </div>

          <div className="flex gap-3 justify-end mt-3">
            <Button variantType="outline" type="button" onClick={() => onOpenChange(false)}>
              {t("patients.modal.cancel")}
            </Button>
            <Button type="submit" disabled={createPatientMutation.isPending}>
              {t("patients.modal.confirm")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
