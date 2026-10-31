import { useEffect, useState } from "react"
import { useTranslation } from "react-i18next"
import { useForm, useWatch, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { CircleAlert, UserPlus } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/shared/components/ui/dialog"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { DiscardChangesDialog } from "@/shared/components/ui/discard-changes-dialog"
import {
  Field,
  FieldDescription,
  FieldError,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/shared/components/ui/field"
import { Spinner } from "@/shared/components/ui/spinner"
import { Input } from "@/shared/components/ui/input"
import { Button } from "@/shared/components/ui/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/components/ui/select"
import { useDiscardGuard } from "@/shared/hooks/useDiscardGuard"
import { StaffRole } from "../../../shared/types"
import { staffFormSchema, type StaffFormData } from "../schemas/staff_schemas"
import { STAFF_ROLE_LABEL_KEYS } from "../staff_roles"
import { toast } from "../../../shared/store/toast_store"
import { useCreateEmployeeMutation, useDepartmentsQuery } from "../queries"
import { useAuthStore } from "../../../shared/store/auth_store"

interface StaffModalProps {
  isOpen: boolean
  onClose: () => void
}

const labelClasses = "text-xs font-semibold"
const touchTargetClasses = "h-11 md:h-9"
const selectTouchTargetClasses = "w-full data-[size=default]:h-11 md:data-[size=default]:h-9"
const hintClasses = "text-xs"
const legendClasses = "mb-0 text-sm font-bold"

const ROLE_HINT_KEYS: Record<StaffRole, string> = {
  [StaffRole.Doctor]: "modal.roleHints.doctor",
  [StaffRole.Nurse]: "modal.roleHints.nurse",
  [StaffRole.Receptionist]: "modal.roleHints.receptionist",
  [StaffRole.Admin]: "modal.roleHints.admin",
}

const LICENSE_HINT_KEYS: Record<StaffRole, string> = {
  [StaffRole.Doctor]: "modal.licenseHints.doctor",
  [StaffRole.Nurse]: "modal.licenseHints.nurse",
  [StaffRole.Receptionist]: "modal.licenseHints.other",
  [StaffRole.Admin]: "modal.licenseHints.other",
}

const FIELD_IDS: Record<keyof StaffFormData, string> = {
  fullName: "staff-full-name",
  email: "staff-email",
  role: "staff-role",
  departmentId: "staff-department",
  license: "staff-license",
}

const FIELD_LABEL_KEYS: Record<keyof StaffFormData, string> = {
  fullName: "modal.name",
  email: "modal.email",
  role: "modal.category",
  departmentId: "modal.department",
  license: "modal.license",
}

const SUMMARY_ORDER: Array<keyof StaffFormData> = ["fullName", "email", "role", "departmentId", "license"]

const RequiredMark = () => {
  const { t } = useTranslation("staff")

  return (
    <>
      <span aria-hidden="true" className="text-destructive">*</span>
      <span className="sr-only">{t("modal.requiredLabel")}</span>
    </>
  )
}

const describedBy = (fieldId: string, hasError: boolean) =>
  hasError ? `${fieldId}-hint ${fieldId}-error` : `${fieldId}-hint`

export const StaffModal = ({ isOpen, onClose }: StaffModalProps) => {
  const { t } = useTranslation("staff")
  const createEmployeeMutation = useCreateEmployeeMutation()
  const { data: departments = [], isLoading: isLoadingDepartments } = useDepartmentsQuery()
  const authenticatedUserId = useAuthStore((state) => state.userId)
  const [isSubmitAttempted, setIsSubmitAttempted] = useState(false)

  const {
    register,
    handleSubmit,
    control,
    reset,
    setFocus,
    setValue,
    getValues,
    formState: { errors, isDirty },
  } = useForm<StaffFormData>({
    resolver: zodResolver(staffFormSchema),
    mode: "onTouched",
    defaultValues: {
      fullName: "",
      role: StaffRole.Doctor,
      license: "",
      email: "",
      departmentId: "",
    },
  })

  const selectedRole = useWatch({ control, name: "role" })
  const isPending = createEmployeeMutation.isPending

  const handleCloseModal = () => {
    setIsSubmitAttempted(false)
    onClose()
  }

  const { isConfirmationOpen, requestClose, keepEditing, discardChanges } = useDiscardGuard({
    isDirty,
    onClose: handleCloseModal,
  })

  useEffect(() => {
    if (!isOpen) {
      reset()
    }
  }, [isOpen, reset])

  useEffect(() => {
    if (isOpen && departments.length === 1 && !getValues("departmentId")) {
      setValue("departmentId", departments[0].id)
    }
  }, [isOpen, departments, getValues, setValue])

  const errorEntries = SUMMARY_ORDER.filter((fieldName) => errors[fieldName])
  const hasSummary = isSubmitAttempted && errorEntries.length > 0

  const translateFieldError = (message?: string) => (message ? t(message) : undefined)

  const handleAttemptClose = () => {
    if (isPending) {
      return
    }
    requestClose()
  }

  const handleRegisterStaff = async (formData: StaffFormData) => {
    try {
      if (!authenticatedUserId) {
        toast.error(t("toast.sessionError"))
        return
      }
      await createEmployeeMutation.mutateAsync({
        userId: authenticatedUserId,
        fullName: formData.fullName,
        email: formData.email,
        role: formData.role,
        crmNumber: formData.license,
        departmentId: formData.departmentId,
      })

      handleCloseModal()
      toast.success(t("toast.createSuccess"))
    } catch {
      toast.error(t("toast.createError"))
    }
  }

  const departmentPlaceholderKey = isLoadingDepartments
    ? "modal.departmentLoading"
    : "modal.departmentPlaceholder"

  return (
    <>
      <Dialog open={isOpen} onOpenChange={(open) => !open && handleAttemptClose()}>
        <DialogContent className="flex max-h-[90dvh] flex-col gap-0 p-0 sm:max-w-[640px]">
          <DialogHeader className="border-b border-border px-6 pb-3 pt-6">
            <DialogTitle className="flex items-center gap-2">
              <UserPlus className="size-5 text-primary" />
              {t("modal.title")}
            </DialogTitle>
            <DialogDescription>{t("requiredHint")}</DialogDescription>
          </DialogHeader>

          <form
            onSubmit={handleSubmit(handleRegisterStaff, () => setIsSubmitAttempted(true))}
            className="flex min-h-0 flex-1 flex-col"
            noValidate
            aria-busy={isPending}
          >
            <div className="flex flex-1 flex-col gap-6 overflow-y-auto px-6 py-4">
              {hasSummary && (
                <Alert variant="destructive">
                  <CircleAlert />
                  <AlertTitle>{t("modal.errorSummary.title", { count: errorEntries.length })}</AlertTitle>
                  <AlertDescription>
                    <ul className="flex flex-col gap-1">
                      {errorEntries.map((fieldName) => (
                        <li key={fieldName}>
                          <button
                            type="button"
                            onClick={() => setFocus(fieldName)}
                            className="min-h-11 text-left text-sm underline underline-offset-2 md:min-h-0"
                          >
                            {t(FIELD_LABEL_KEYS[fieldName])}: {translateFieldError(errors[fieldName]?.message)}
                          </button>
                        </li>
                      ))}
                    </ul>
                  </AlertDescription>
                </Alert>
              )}

              <FieldSet className="gap-4">
                <FieldLegend variant="label" className={legendClasses}>
                  {t("modal.sections.identity")}
                </FieldLegend>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Field data-invalid={!!errors.fullName} className="md:col-span-2">
                    <FieldLabel htmlFor={FIELD_IDS.fullName} className={labelClasses}>
                      {t("modal.name")} <RequiredMark />
                    </FieldLabel>
                    <Input
                      id={FIELD_IDS.fullName}
                      type="text"
                      autoFocus
                      autoComplete="off"
                      autoCapitalize="words"
                      maxLength={255}
                      placeholder={t("modal.namePlaceholder")}
                      aria-required="true"
                      aria-invalid={!!errors.fullName}
                      aria-describedby={describedBy(FIELD_IDS.fullName, !!errors.fullName)}
                      className={touchTargetClasses}
                      {...register("fullName")}
                    />
                    <FieldDescription id={`${FIELD_IDS.fullName}-hint`} className={hintClasses}>
                      {t("modal.nameHint")}
                    </FieldDescription>
                    <FieldError id={`${FIELD_IDS.fullName}-error`}>
                      {translateFieldError(errors.fullName?.message)}
                    </FieldError>
                  </Field>

                  <Field data-invalid={!!errors.email} className="md:col-span-2">
                    <FieldLabel htmlFor={FIELD_IDS.email} className={labelClasses}>
                      {t("modal.email")} <RequiredMark />
                    </FieldLabel>
                    <Input
                      id={FIELD_IDS.email}
                      type="email"
                      inputMode="email"
                      autoComplete="off"
                      autoCapitalize="none"
                      spellCheck={false}
                      maxLength={255}
                      placeholder={t("modal.emailPlaceholder")}
                      aria-required="true"
                      aria-invalid={!!errors.email}
                      aria-describedby={describedBy(FIELD_IDS.email, !!errors.email)}
                      className={touchTargetClasses}
                      {...register("email")}
                    />
                    <FieldDescription id={`${FIELD_IDS.email}-hint`} className={hintClasses}>
                      {t("modal.emailHint")}
                    </FieldDescription>
                    <FieldError id={`${FIELD_IDS.email}-error`}>
                      {translateFieldError(errors.email?.message)}
                    </FieldError>
                  </Field>
                </div>
              </FieldSet>

              <FieldSet className="gap-4">
                <FieldLegend variant="label" className={legendClasses}>
                  {t("modal.sections.assignment")}
                </FieldLegend>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Field data-invalid={!!errors.role}>
                    <FieldLabel htmlFor={FIELD_IDS.role} className={labelClasses}>
                      {t("modal.category")} <RequiredMark />
                    </FieldLabel>
                    <Controller
                      control={control}
                      name="role"
                      render={({ field }) => (
                        <Select onValueChange={field.onChange} value={field.value}>
                          <SelectTrigger
                            id={FIELD_IDS.role}
                            ref={field.ref}
                            onBlur={field.onBlur}
                            className={selectTouchTargetClasses}
                            aria-required="true"
                            aria-invalid={!!errors.role}
                            aria-describedby={describedBy(FIELD_IDS.role, !!errors.role)}
                          >
                            <SelectValue placeholder={t("modal.category")} />
                          </SelectTrigger>
                          <SelectContent>
                            {Object.values(StaffRole).map((roleOption) => (
                              <SelectItem key={roleOption} value={roleOption}>
                                {t(STAFF_ROLE_LABEL_KEYS[roleOption])}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      )}
                    />
                    <FieldDescription id={`${FIELD_IDS.role}-hint`} className={hintClasses}>
                      {t(ROLE_HINT_KEYS[selectedRole])}
                    </FieldDescription>
                    <FieldError id={`${FIELD_IDS.role}-error`}>
                      {translateFieldError(errors.role?.message)}
                    </FieldError>
                  </Field>

                  <Field data-invalid={!!errors.departmentId}>
                    <FieldLabel htmlFor={FIELD_IDS.departmentId} className={labelClasses}>
                      {t("modal.department")} <RequiredMark />
                    </FieldLabel>
                    <Controller
                      control={control}
                      name="departmentId"
                      render={({ field }) => (
                        <Select
                          onValueChange={field.onChange}
                          value={field.value}
                          disabled={isLoadingDepartments}
                        >
                          <SelectTrigger
                            id={FIELD_IDS.departmentId}
                            ref={field.ref}
                            onBlur={field.onBlur}
                            className={selectTouchTargetClasses}
                            aria-required="true"
                            aria-invalid={!!errors.departmentId}
                            aria-describedby={describedBy(FIELD_IDS.departmentId, !!errors.departmentId)}
                          >
                            <SelectValue placeholder={t(departmentPlaceholderKey)} />
                          </SelectTrigger>
                          <SelectContent>
                            {departments.map((department) => (
                              <SelectItem key={department.id} value={department.id}>
                                {department.name}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      )}
                    />
                    <FieldDescription id={`${FIELD_IDS.departmentId}-hint`} className={hintClasses}>
                      {t("modal.departmentHint")}
                    </FieldDescription>
                    <FieldError id={`${FIELD_IDS.departmentId}-error`}>
                      {translateFieldError(errors.departmentId?.message)}
                    </FieldError>
                  </Field>
                </div>
              </FieldSet>

              <FieldSet className="gap-4">
                <FieldLegend variant="label" className={legendClasses}>
                  {t("modal.sections.registration")}
                </FieldLegend>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Field data-invalid={!!errors.license}>
                    <FieldLabel htmlFor={FIELD_IDS.license} className={labelClasses}>
                      {t("modal.license")} <RequiredMark />
                    </FieldLabel>
                    <Input
                      id={FIELD_IDS.license}
                      type="text"
                      autoComplete="off"
                      autoCapitalize="characters"
                      spellCheck={false}
                      maxLength={50}
                      placeholder={t("modal.licensePlaceholder")}
                      aria-required="true"
                      aria-invalid={!!errors.license}
                      aria-describedby={describedBy(FIELD_IDS.license, !!errors.license)}
                      className={touchTargetClasses}
                      {...register("license")}
                    />
                    <FieldDescription id={`${FIELD_IDS.license}-hint`} className={hintClasses}>
                      {t(LICENSE_HINT_KEYS[selectedRole])}
                    </FieldDescription>
                    <FieldError id={`${FIELD_IDS.license}-error`}>
                      {translateFieldError(errors.license?.message)}
                    </FieldError>
                  </Field>
                </div>
              </FieldSet>
            </div>

            <DialogFooter className="border-t border-border bg-background px-6 py-4">
              <Button
                type="button"
                variant="outline"
                onClick={handleAttemptClose}
                disabled={isPending}
                className="h-11 text-xs font-bold md:h-9"
              >
                {t("modal.cancel")}
              </Button>
              <Button type="submit" disabled={isPending} className="h-11 text-xs font-bold md:h-9">
                {isPending && <Spinner />}
                {isPending ? t("modal.saving") : t("modal.confirm")}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      <DiscardChangesDialog
        open={isConfirmationOpen}
        onKeepEditing={keepEditing}
        onDiscard={discardChanges}
      />
    </>
  )
}
