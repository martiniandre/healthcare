import { useTranslation } from "react-i18next"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { UserPlus } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
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
import { toast } from "../../../shared/store/toast_store"
import { useCreateEmployeeMutation } from "../queries"
import { useEffect } from "react"

interface StaffModalProps {
  isOpen: boolean
  onClose: () => void
}

export const StaffModal = ({ isOpen, onClose }: StaffModalProps) => {
  const { t } = useTranslation()
  const createEmployeeMutation = useCreateEmployeeMutation()

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
      department: "",
    },
  })

  useEffect(() => {
    if (!isOpen) {
      reset()
    }
  }, [isOpen, reset])

  const handleRegisterStaff = async (formData: StaffFormData) => {
    try {
      const temporaryRandomUserId = crypto.randomUUID()
      await createEmployeeMutation.mutateAsync({
        userId: temporaryRandomUserId,
        fullName: formData.fullName,
        email: formData.email,
        role: formData.role,
        crmNumber: formData.license || "N/A",
      })

      onClose()
      toast.success(t("staff.toast.createSuccess"))
    } catch {
      toast.error(t("staff.toast.createError"))
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[480px]">
        <DialogHeader className="border-b border-border pb-3">
          <DialogTitle className="flex items-center gap-2">
            <UserPlus className="w-5 h-5 text-primary" />
            {t("staff.modal.title")}
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit(handleRegisterStaff)} className="flex flex-col gap-4" noValidate>
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">{t("staff.modal.name")}</label>
            <Input
              type="text"
              placeholder={t("staff.modal.namePlaceholder")}
              errorText={errors.fullName?.message}
              {...register("fullName")}
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">{t("staff.modal.category")}</label>
              <Controller
                control={control}
                name="role"
                render={({ field }) => (
                  <Select onValueChange={field.onChange} defaultValue={field.value}>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder={t("staff.modal.category")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value={StaffRole.Doctor}>{StaffRole.Doctor}</SelectItem>
                      <SelectItem value={StaffRole.Nurse}>{StaffRole.Nurse}</SelectItem>
                      <SelectItem value={StaffRole.Receptionist}>{StaffRole.Receptionist}</SelectItem>
                      <SelectItem value={StaffRole.Admin}>{StaffRole.Admin}</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.role?.message && (
                <span className="text-xs text-red-500 font-medium px-1 mt-1">
                  {errors.role.message}
                </span>
              )}
            </div>

            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">{t("staff.modal.license")}</label>
              <Input
                type="text"
                placeholder={t("staff.modal.licensePlaceholder")}
                errorText={errors.license?.message}
                {...register("license")}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">{t("staff.modal.email")}</label>
              <Input
                type="email"
                placeholder={t("staff.modal.emailPlaceholder")}
                errorText={errors.email?.message}
                {...register("email")}
              />
            </div>

            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">{t("staff.modal.department")}</label>
              <Input
                type="text"
                placeholder={t("staff.modal.departmentPlaceholder")}
                errorText={errors.department?.message}
                {...register("department")}
              />
            </div>
          </div>

          <div className="flex gap-3 justify-end border-t border-border pt-4 mt-2">
            <Button
              type="button"
              variantType="outline"
              onClick={onClose}
              className="px-4 py-2 text-xs font-bold"
            >
              {t("staff.modal.cancel")}
            </Button>
            <Button
              type="submit"
              disabled={createEmployeeMutation.isPending}
              variantType="primary"
              className="px-4 py-2 text-xs font-bold"
            >
              {t("staff.modal.confirm")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
