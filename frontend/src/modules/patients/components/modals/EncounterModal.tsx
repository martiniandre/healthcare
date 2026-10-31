import { useMemo } from "react"
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
  hasSubmitError?: boolean
}

export const EncounterModal = ({ isOpen, onClose, onSubmit, isPending, hasSubmitError }: EncounterModalProps) => {
  const { data: doctors = [] } = useStaffListQuery("", StaffRole.Doctor)

  const encounterFormConfig = useMemo(() => {
    const doctorOptions = doctors
      .filter((doctor) => doctor.fhirResourceId)
      .map((doctor) => ({ value: doctor.fhirResourceId, label: doctor.fullName }))
    return buildEncounterFormConfig(doctorOptions)
  }, [doctors])

  return (
    <ClinicalFormModal
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      isPending={isPending}
      hasSubmitError={hasSubmitError}
      config={encounterFormConfig}
    />
  )
}
