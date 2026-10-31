import { ClinicalFormModal } from "../ClinicalFormModal/ClinicalFormModal"
import { vitalSignsPanelFormConfig } from "../ClinicalFormModal/clinicalFormConfigs"
import type { NewVitalSignsPanelFormData } from "../../patient_schemas"

interface ObservationModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (panelFormData: NewVitalSignsPanelFormData) => void
  isPending: boolean
  hasSubmitError?: boolean
}

export const ObservationModal = ({ isOpen, onClose, onSubmit, isPending, hasSubmitError }: ObservationModalProps) => {
  return (
    <ClinicalFormModal
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      isPending={isPending}
      hasSubmitError={hasSubmitError}
      config={vitalSignsPanelFormConfig}
    />
  )
}
