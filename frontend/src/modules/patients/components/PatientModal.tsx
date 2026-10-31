import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../shared/components/ui/dialog"
import { MaskedInput } from "../../../shared/components/ui/MaskedInput"
import { Button } from "../../../shared/components/ui/button"
import { Input } from "../../../shared/components/ui/input"
import { Field, FieldError, FieldLabel } from "../../../shared/components/ui/field"
import { Spinner } from "../../../shared/components/ui/spinner"
import { getNewPatientSchema, type NewPatientFormData } from "../patient_schemas"
import { createModuleTranslator } from "../../../shared/i18n/i18n"
import { localizedDateToIso } from "../../../shared/utils/dates"
import { useCreatePatientMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import { useEffect } from "react"

interface PatientModalProps {
  isOpen: boolean
  onOpenChange: (open: boolean) => void
}

export const PatientModal = ({ isOpen, onOpenChange }: PatientModalProps) => {
  const { t, i18n } = useTranslation("patients")
  const createPatientMutation = useCreatePatientMutation()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<NewPatientFormData>({
    resolver: zodResolver(getNewPatientSchema(createModuleTranslator("patients"))),
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
        birth_date: localizedDateToIso(formData.birthDate, i18n.language),
        document_id: formData.documentId,
        phone_number: formData.phoneNumber,
      })
      onOpenChange(false)
      toast.success(t("toast.createSuccess"))
    } catch {
      toast.error(t("toast.createError"))
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle className="text-left">{t("modal.title")}</DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 text-left mt-4" noValidate>
          <Field className="gap-1" data-invalid={Boolean(errors.fullName)}>
            <FieldLabel className="text-xs font-semibold text-muted-foreground">{t("modal.name")}</FieldLabel>
            <Input
              type="text"
              placeholder={t("modal.namePlaceholder")}
              aria-invalid={Boolean(errors.fullName)}
              {...register("fullName")}
            />
            {errors.fullName && <FieldError>{errors.fullName.message}</FieldError>}
          </Field>

          <Field className="gap-1" data-invalid={Boolean(errors.birthDate)}>
            <FieldLabel className="text-xs font-semibold text-muted-foreground">{t("modal.birthDate")}</FieldLabel>
            <MaskedInput
              mask="99/99/9999"
              placeholder={t("modal.birthDatePlaceholder")}
              errorText={errors.birthDate?.message}
              {...register("birthDate")}
            />
          </Field>

          <Field className="gap-1" data-invalid={Boolean(errors.documentId)}>
            <FieldLabel className="text-xs font-semibold text-muted-foreground">{t("modal.document")}</FieldLabel>
            <MaskedInput
              mask="999.999.999-99"
              placeholder={t("modal.documentPlaceholder")}
              errorText={errors.documentId?.message}
              {...register("documentId")}
            />
          </Field>

          <Field className="gap-1" data-invalid={Boolean(errors.phoneNumber)}>
            <FieldLabel className="text-xs font-semibold text-muted-foreground">{t("modal.phone")}</FieldLabel>
            <MaskedInput
              mask="(99) 99999-9999"
              placeholder={t("modal.phonePlaceholder")}
              errorText={errors.phoneNumber?.message}
              {...register("phoneNumber")}
            />
          </Field>

          <div className="flex gap-3 justify-end mt-3">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              {t("modal.cancel")}
            </Button>
            <Button type="submit" disabled={createPatientMutation.isPending}>
              {createPatientMutation.isPending && <Spinner />}
              {t("modal.confirm")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
