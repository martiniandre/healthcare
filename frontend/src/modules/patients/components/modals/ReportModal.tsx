import type { NewReportFormData } from "../../patient_schemas"
import { ClinicalFormModal } from "../ClinicalFormModal/ClinicalFormModal"
import { reportFormConfig } from "../ClinicalFormModal/clinicalFormConfigs"

interface ReportModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewReportFormData) => void
  isPending: boolean
  hasSubmitError?: boolean
}

export const ReportModal = ({ isOpen, onClose, onSubmit, isPending, hasSubmitError }: ReportModalProps) => {
  return (
    <ClinicalFormModal
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      isPending={isPending}
      hasSubmitError={hasSubmitError}
      config={reportFormConfig}
    />
  )
}
