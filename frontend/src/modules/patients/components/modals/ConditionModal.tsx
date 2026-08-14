import type { NewConditionFormData } from "../../patient_schemas"
import { ClinicalFormModal } from "../ClinicalFormModal/ClinicalFormModal"
import { conditionFormConfig } from "../ClinicalFormModal/clinicalFormConfigs"

interface ConditionModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewConditionFormData) => void
  isPending: boolean
}

export const ConditionModal = ({ isOpen, onClose, onSubmit, isPending }: ConditionModalProps) => {
  return (
    <ClinicalFormModal
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      isPending={isPending}
      config={conditionFormConfig}
    />
  )
}
