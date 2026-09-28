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
} from "../../../shared/components/ui/Dialog"
import { Input } from "../../../shared/components/ui/Input"
import { Button } from "../../../shared/components/ui/Button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../../shared/components/ui/Select"
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

const requiredLabelClasses = "text-xs font-semibold text-gray-600"

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
            <UserPlus className="w-5 h-5 text-primary" />
            {t("modal.title")}
          </DialogTitle>
          <DialogDescription>{t("requiredHint")}</DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(handleRegisterStaff)} className="flex flex-col gap-4" noValidate>
          <div className="flex flex-col gap-1">
            <label htmlFor="staff-full-name" className={requiredLabelClasses}>
              {t("modal.name")} <span className="text-red-500">*</span>
            </label>
            <Input
              id="staff-full-name"
              type="text"
              placeholder={t("modal.namePlaceholder")}
              errorText={translateFieldError(errors.fullName?.message)}
              {...register("fullName")}
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="flex flex-col gap-1">
              <label htmlFor="staff-role" className={requiredLabelClasses}>
                {t("modal.category")} <span className="text-red-500">*</span>
              </label>
              <Controller
                control={control}
                name="role"
                render={({ field }) => (
                  <Select onValueChange={field.onChange} value={field.value}>
                    <SelectTrigger id="staff-role" className="w-full">
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
              {errors.role?.message && (
                <span className="text-xs text-red-500 font-medium px-1 mt-1">
                  {translateFieldError(errors.role.message)}
                </span>
              )}
            </div>

            <div className="flex flex-col gap-1">
              <label htmlFor="staff-license" className={requiredLabelClasses}>
                {t("modal.license")} <span className="text-red-500">*</span>
              </label>
              <Input
                id="staff-license"
                type="text"
                placeholder={t("modal.licensePlaceholder")}
                errorText={translateFieldError(errors.license?.message)}
                {...register("license")}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="flex flex-col gap-1">
              <label htmlFor="staff-email" className={requiredLabelClasses}>
                {t("modal.email")} <span className="text-red-500">*</span>
              </label>
              <Input
                id="staff-email"
                type="email"
                placeholder={t("modal.emailPlaceholder")}
                errorText={translateFieldError(errors.email?.message)}
                {...register("email")}
              />
            </div>

            <div className="flex flex-col gap-1">
              <label htmlFor="staff-department" className={requiredLabelClasses}>
                {t("modal.department")} <span className="text-red-500">*</span>
              </label>
              <Controller
                control={control}
                name="departmentId"
                render={({ field }) => (
                  <Select onValueChange={field.onChange} value={field.value}>
                    <SelectTrigger id="staff-department" className="w-full">
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
              {errors.departmentId?.message && (
                <span className="text-xs text-red-500 font-medium px-1 mt-1">
                  {translateFieldError(errors.departmentId.message)}
                </span>
              )}
            </div>
          </div>

          <div className="flex gap-3 justify-end border-t border-border pt-4 mt-2">
            <Button
              type="button"
              variantType="outline"
              onClick={onClose}
              className="px-4 py-2 text-xs font-bold"
            >
              {t("modal.cancel")}
            </Button>
            <Button
              type="submit"
              disabled={createEmployeeMutation.isPending || isLoadingDepartments}
              variantType="primary"
              className="px-4 py-2 text-xs font-bold"
            >
              {t("modal.confirm")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
