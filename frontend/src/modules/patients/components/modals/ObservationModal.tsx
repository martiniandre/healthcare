import { ClinicalFormModal } from "../ClinicalFormModal/ClinicalFormModal"
import {
  observationFormConfig,
  type SubmittedObservationFormData,
} from "../ClinicalFormModal/clinicalFormConfigs"

interface ObservationModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: SubmittedObservationFormData) => void
  isPending: boolean
}

export const ObservationModal = ({ isOpen, onClose, onSubmit, isPending }: ObservationModalProps) => {
  return (
    <ClinicalFormModal
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      isPending={isPending}
      config={observationFormConfig}
    />
  )
}
