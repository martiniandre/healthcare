import type { NewEncounterFormData } from "../../patient_schemas"
import { useStaffListQuery } from "../../../staff/queries"
import { StaffRole } from "../../../../shared/types"
import { ClinicalFormModal } from "../ClinicalFormModal/ClinicalFormModal"
import { buildEncounterFormConfig } from "../ClinicalFormModal/clinicalFormConfigs"

interface EncounterModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewEncounterFormData) => void
  isPending: boolean
}

export const EncounterModal = ({ isOpen, onClose, onSubmit, isPending }: EncounterModalProps) => {
  const { data: doctors = [] } = useStaffListQuery("", StaffRole.Doctor)

  const doctorOptions = doctors
    .filter((doctor) => doctor.fhirResourceId)
    .map((doctor) => ({ value: doctor.fhirResourceId, label: doctor.fullName }))

  return (
    <ClinicalFormModal
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      isPending={isPending}
      config={buildEncounterFormConfig(doctorOptions)}
    />
  )
}
