import { useEffect } from "react"
import { useTranslation } from "react-i18next"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { UserPlus } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/shared/components/ui/dialog"
import { Field, FieldError, FieldLabel } from "@/shared/components/ui/field"
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

const requiredLabelClasses = "text-xs font-semibold"

export const StaffModal = ({ isOpen, onClose }: StaffModalProps) => {
  const { t } = useTranslation("staff")
  const createEmployeeMutation = useCreateEmployeeMutation()
  const { data: departments = [], isLoading: isLoadingDepartments } = useDepartmentsQuery()
  const authenticatedUserId = useAuthStore((state) => state.userId)

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors },
  } = useForm<StaffFormData>({
    resolver: zodResolver(staffFormSchema),
    defaultValues: {
      fullName: "",
      role: StaffRole.Doctor,
      license: "",
      email: "",
      departmentId: "",
    },
  })

  useEffect(() => {
    if (!isOpen) {
      reset()
    }
  }, [isOpen, reset])

  const translateFieldError = (message?: string) => (message ? t(message) : undefined)

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

      onClose()
      toast.success(t("toast.createSuccess"))
    } catch {
      toast.error(t("toast.createError"))
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[520px]">
        <DialogHeader className="border-b border-border pb-3">
          <DialogTitle className="flex items-center gap-2">
            <UserPlus className="size-5 text-primary" />
            {t("modal.title")}
          </DialogTitle>
          <DialogDescription>{t("requiredHint")}</DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(handleRegisterStaff)} className="flex flex-col gap-4" noValidate>
          <Field data-invalid={!!errors.fullName}>
            <FieldLabel htmlFor="staff-full-name" className={requiredLabelClasses}>
              {t("modal.name")} <span className="text-destructive">*</span>
            </FieldLabel>
            <Input
              id="staff-full-name"
              type="text"
              placeholder={t("modal.namePlaceholder")}
              aria-invalid={!!errors.fullName}
              {...register("fullName")}
            />
            <FieldError>{translateFieldError(errors.fullName?.message)}</FieldError>
          </Field>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <Field data-invalid={!!errors.role}>
              <FieldLabel htmlFor="staff-role" className={requiredLabelClasses}>
                {t("modal.category")} <span className="text-destructive">*</span>
              </FieldLabel>
              <Controller
                control={control}
                name="role"
                render={({ field }) => (
                  <Select onValueChange={field.onChange} value={field.value}>
                    <SelectTrigger id="staff-role" className="w-full" aria-invalid={!!errors.role}>
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
              <FieldError>{translateFieldError(errors.role?.message)}</FieldError>
            </Field>

            <Field data-invalid={!!errors.license}>
              <FieldLabel htmlFor="staff-license" className={requiredLabelClasses}>
                {t("modal.license")} <span className="text-destructive">*</span>
              </FieldLabel>
              <Input
                id="staff-license"
                type="text"
                placeholder={t("modal.licensePlaceholder")}
                aria-invalid={!!errors.license}
                {...register("license")}
              />
              <FieldError>{translateFieldError(errors.license?.message)}</FieldError>
            </Field>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <Field data-invalid={!!errors.email}>
              <FieldLabel htmlFor="staff-email" className={requiredLabelClasses}>
                {t("modal.email")} <span className="text-destructive">*</span>
              </FieldLabel>
              <Input
                id="staff-email"
                type="email"
                placeholder={t("modal.emailPlaceholder")}
                aria-invalid={!!errors.email}
                {...register("email")}
              />
              <FieldError>{translateFieldError(errors.email?.message)}</FieldError>
            </Field>

            <Field data-invalid={!!errors.departmentId}>
              <FieldLabel htmlFor="staff-department" className={requiredLabelClasses}>
                {t("modal.department")} <span className="text-destructive">*</span>
              </FieldLabel>
              <Controller
                control={control}
                name="departmentId"
                render={({ field }) => (
                  <Select onValueChange={field.onChange} value={field.value}>
                    <SelectTrigger id="staff-department" className="w-full" aria-invalid={!!errors.departmentId}>
                      <SelectValue placeholder={t("modal.departmentPlaceholder")} />
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
              <FieldError>{translateFieldError(errors.departmentId?.message)}</FieldError>
            </Field>
          </div>

          <DialogFooter className="border-t border-border pt-4 mt-2">
            <Button type="button" variant="outline" onClick={onClose} className="text-xs font-bold">
              {t("modal.cancel")}
            </Button>
            <Button
              type="submit"
              disabled={createEmployeeMutation.isPending || isLoadingDepartments}
              className="text-xs font-bold"
            >
              {createEmployeeMutation.isPending && <Spinner />}
              {t("modal.confirm")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
