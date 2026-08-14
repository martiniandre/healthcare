import type { NewAllergyFormData } from "../../patient_schemas"
import { ClinicalFormModal } from "../ClinicalFormModal/ClinicalFormModal"
import { allergyFormConfig } from "../ClinicalFormModal/clinicalFormConfigs"

interface AllergyModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewAllergyFormData) => void
  isPending: boolean
}

export const AllergyModal = ({ isOpen, onClose, onSubmit, isPending }: AllergyModalProps) => {
  return (
    <ClinicalFormModal
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      isPending={isPending}
      config={allergyFormConfig}
    />
  )
}
