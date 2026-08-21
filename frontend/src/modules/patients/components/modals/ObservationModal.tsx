import { ClinicalFormModal } from "../ClinicalFormModal/ClinicalFormModal"
import { vitalSignsPanelFormConfig } from "../ClinicalFormModal/clinicalFormConfigs"
import type { NewVitalSignsPanelFormData } from "../../patient_schemas"

interface ObservationModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (panelFormData: NewVitalSignsPanelFormData) => void
  isPending: boolean
}

export const ObservationModal = ({ isOpen, onClose, onSubmit, isPending }: ObservationModalProps) => {
  return (
    <ClinicalFormModal
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      isPending={isPending}
      config={vitalSignsPanelFormConfig}
    />
  )
}
