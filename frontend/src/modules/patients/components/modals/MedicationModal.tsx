import type { NewMedicationFormData } from "../../patient_schemas"
import { ClinicalFormModal } from "../ClinicalFormModal/ClinicalFormModal"
import { medicationFormConfig } from "../ClinicalFormModal/clinicalFormConfigs"

interface MedicationModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewMedicationFormData) => void
  isPending: boolean
}

export const MedicationModal = ({ isOpen, onClose, onSubmit, isPending }: MedicationModalProps) => {
  return (
    <ClinicalFormModal
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      isPending={isPending}
      config={medicationFormConfig}
    />
  )
}
