import { useEffect, useId } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { Dialog } from "../../../shared/components/ui/dialog"
import { DiscardChangesDialog } from "../../../shared/components/ui/discard-changes-dialog"
import { FieldLegend, FieldSet } from "../../../shared/components/ui/field"
import { Input } from "../../../shared/components/ui/input"
import { MaskedInput } from "../../../shared/components/ui/MaskedInput"
import { useDiscardGuard } from "../../../shared/hooks/useDiscardGuard"
import { getNewPatientSchema, type NewPatientFormData } from "../patient_schemas"
import { createModuleTranslator } from "../../../shared/i18n/i18n"
import { localizedDateToIso } from "../../../shared/utils/dates"
import { useCreatePatientMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import { FormDialogLayout } from "./form/FormDialogLayout"
import { FormField } from "./form/FormField"
import { FormSubmitError } from "./form/FormSubmitError"
import { useSubmissionLock } from "./form/useSubmissionLock"

interface PatientModalProps {
  isOpen: boolean
  onOpenChange: (open: boolean) => void
}

const EMPTY_PATIENT_VALUES: NewPatientFormData = {
  fullName: "",
  birthDate: "",
  documentId: "",
  phoneNumber: "",
}

export const PatientModal = ({ isOpen, onOpenChange }: PatientModalProps) => {
  const { t, i18n } = useTranslation("patients")
  const formId = useId()
  const runWithSubmissionLock = useSubmissionLock()
  const {
    mutateAsync: createPatient,
    reset: resetCreatePatient,
    isPending,
    isError,
  } = useCreatePatientMutation()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isDirty },
  } = useForm<NewPatientFormData>({
    resolver: zodResolver(getNewPatientSchema(createModuleTranslator("patients"))),
    defaultValues: EMPTY_PATIENT_VALUES,
    mode: "onTouched",
  })

  const closeDialog = () => onOpenChange(false)

  const { isConfirmationOpen, requestClose, keepEditing, discardChanges } = useDiscardGuard({
    isDirty,
    onClose: closeDialog,
  })

  useEffect(() => {
    if (isOpen) {
      reset(EMPTY_PATIENT_VALUES)
      resetCreatePatient()
    }
  }, [isOpen, reset, resetCreatePatient])

  const onSubmit = async (formData: NewPatientFormData) => {
    if (isPending) {
      return
    }
    await runWithSubmissionLock(async () => {
      try {
        await createPatient({
          full_name: formData.fullName,
          birth_date: localizedDateToIso(formData.birthDate, i18n.language),
          document_id: formData.documentId,
          phone_number: formData.phoneNumber,
        })
        closeDialog()
        toast.success(t("toast.createSuccess"))
      } catch {
        toast.error(t("toast.createError"))
      }
    })
  }

  return (
    <>
      <Dialog open={isOpen} onOpenChange={(open) => !open && requestClose()}>
        <FormDialogLayout
          title={t("modal.title")}
          description={t("modal.description")}
          widthClassName="sm:max-w-[560px]"
          onSubmit={handleSubmit(onSubmit)}
          cancelLabel={t("modal.cancel")}
          onCancel={requestClose}
          submitLabel={t("modal.confirm")}
          isPending={isPending}
        >
          {isError ? <FormSubmitError /> : null}
          <p className="text-xs text-muted-foreground">{t("form.requiredLegend")}</p>

          <FieldSet className="gap-3">
            <FieldLegend variant="label" className="mb-0 text-sm font-semibold">
              {t("modal.sections.identification")}
            </FieldLegend>
            <div className="grid gap-4 md:grid-cols-2">
              <FormField
                id={`${formId}-fullName`}
                label={t("modal.name")}
                description={t("modal.nameHint")}
                error={errors.fullName?.message}
                className="md:col-span-2"
              >
                {(controlProperties) => (
                  <Input
                    type="text"
                    placeholder={t("modal.namePlaceholder")}
                    autoComplete="off"
                    autoCapitalize="words"
                    className="min-h-11 md:min-h-9"
                    {...controlProperties}
                    {...register("fullName")}
                  />
                )}
              </FormField>

              <FormField
                id={`${formId}-birthDate`}
                label={t("modal.birthDate")}
                description={t("modal.birthDateHint")}
                error={errors.birthDate?.message}
              >
                {(controlProperties) => (
                  <MaskedInput
                    mask="99/99/9999"
                    placeholder={t("modal.birthDatePlaceholder")}
                    inputMode="numeric"
                    autoComplete="off"
                    className="min-h-11 md:min-h-9"
                    {...controlProperties}
                    {...register("birthDate")}
                  />
                )}
              </FormField>

              <FormField
                id={`${formId}-documentId`}
                label={t("modal.document")}
                description={t("modal.documentHint")}
                error={errors.documentId?.message}
              >
                {(controlProperties) => (
                  <MaskedInput
                    mask="999.999.999-99"
                    placeholder={t("modal.documentPlaceholder")}
                    inputMode="numeric"
                    autoComplete="off"
                    spellCheck={false}
                    className="min-h-11 md:min-h-9"
                    {...controlProperties}
                    {...register("documentId")}
                  />
                )}
              </FormField>
            </div>
          </FieldSet>

          <FieldSet className="gap-3">
            <FieldLegend variant="label" className="mb-0 text-sm font-semibold">
              {t("modal.sections.contact")}
            </FieldLegend>
            <div className="grid gap-4 md:grid-cols-2">
              <FormField
                id={`${formId}-phoneNumber`}
                label={t("modal.phone")}
                description={t("modal.phoneHint")}
                error={errors.phoneNumber?.message}
              >
                {(controlProperties) => (
                  <MaskedInput
                    mask="(99) 99999-9999"
                    type="tel"
                    placeholder={t("modal.phonePlaceholder")}
                    inputMode="tel"
                    autoComplete="off"
                    className="min-h-11 md:min-h-9"
                    {...controlProperties}
                    {...register("phoneNumber")}
                  />
                )}
              </FormField>
            </div>
          </FieldSet>
        </FormDialogLayout>
      </Dialog>
      <DiscardChangesDialog open={isConfirmationOpen} onKeepEditing={keepEditing} onDiscard={discardChanges} />
    </>
  )
}
