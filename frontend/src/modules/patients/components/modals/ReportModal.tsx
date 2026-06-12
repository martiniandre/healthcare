import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { Send } from "lucide-react"
import { Input } from "../../../../shared/components/ui/Input"
import { Button } from "../../../../shared/components/ui/Button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../../shared/components/ui/Dialog"
import { getNewReportSchema, type NewReportFormData } from "../../patient_schemas"

interface ReportModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewReportFormData) => void
  isPending: boolean
}

export const ReportModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: ReportModalProps) => {
  const { t } = useTranslation("patients")

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<NewReportFormData>({
    resolver: zodResolver(getNewReportSchema(t)),
  })

  if (!isOpen) {
    return null
  }

  return (
      <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle className="text-left">
              {t("modals.report.title")}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 text-left mt-4">
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("modals.report.exam")}
              </label>
              <Input
                type="text"
                placeholder={t("modals.report.examPlaceholder")}
                errorText={errors.reportDisplay?.message}
                {...register("reportDisplay")}
              />
            </div>

            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("modals.report.conclusion")}
              </label>
              <textarea
                placeholder={t("modals.report.conclusionPlaceholder")}
                className="w-full bg-gray-50 border border-border rounded-lg px-3.5 py-2.5 h-28 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-primary/50"
                {...register("conclusion")}
              />
              {errors.conclusion?.message && (
                <span className="text-xs text-red-500 font-medium px-1 mt-1">
                  {errors.conclusion.message}
                </span>
              )}
            </div>

            <div className="flex gap-3 justify-end mt-4">
              <Button variantType="outline" type="button" onClick={onClose}>
                {t("modal.cancel")}
              </Button>
              <Button type="submit" disabled={isPending} className="gap-2">
                <Send className="w-3.5 h-3.5" />
                {t("modals.report.confirm")}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
  )
}
