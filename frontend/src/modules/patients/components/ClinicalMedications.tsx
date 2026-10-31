import { useState } from "react"
import { Pill, Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/button"
import { ClinicalQueryTable } from "./ClinicalQueryTable"
import { MedicationModal } from "./modals/MedicationModal"
import { useClinicalMedicationsColumns } from "./useClinicalMedicationsColumns"
import { useMedicationsQuery, useCreateMedicationMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"

interface ClinicalMedicationsProps {
  patientId: string
  encounterId: string
}

export default function ClinicalMedications({ patientId, encounterId }: ClinicalMedicationsProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const columns = useClinicalMedicationsColumns()
  const { data: medications = [], isLoading, isError, refetch } = useMedicationsQuery(encounterId)
  const createMedicationMutation = useCreateMedicationMutation()

  const handleCreateMedication = async (formData: { medicationDisplay: string; dosageInstruction: string }) => {
    try {
      await createMedicationMutation.mutateAsync({
        encounter_fhir_id: encounterId,
        patient_fhir_id: patientId,
        medication_name: formData.medicationDisplay,
        dosage_instructions: formData.dosageInstruction,
      })
      setIsModalOpen(false)
      toast.success(t("toast.medicationSuccess"))
    } catch {
      toast.error(t("toast.medicationError"))
    }
  }

  return (
    <>
      <ClinicalQueryTable
        title={t("details.medicationsCard.title")}
        icon={<Pill className="w-4 h-4 text-purple-500" />}
        columns={columns}
        data={medications}
        isLoading={isLoading}
        isError={isError}
        onRetry={() => void refetch()}
        emptyText={t("details.medicationsCard.empty")}
        emptyIcon={Pill}
        addButton={
          <Can I={Action.Create} a={Feature.MedicationRequest}>
            <Button size="sm" onClick={() => setIsModalOpen(true)}>
              <Plus className="w-3.5 h-3.5" />
              {t("details.medicationsCard.add")}
            </Button>
          </Can>
        }
      />

      <MedicationModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateMedication}
        isPending={createMedicationMutation.isPending}
      />
    </>
  )
}
