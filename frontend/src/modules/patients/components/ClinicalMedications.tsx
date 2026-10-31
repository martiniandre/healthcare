import { useState } from "react"
import { Pill, Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/button"
import { ClinicalQueryTable } from "./ClinicalQueryTable"
import { MedicationModal, type MedicationSubmission, type MedicationSubmitOutcome } from "./modals/MedicationModal"
import { buildAlertMessage, extractAllergyBlockedAlerts, type AlertTranslateFunction } from "../medicationSafety"
import type { MedicationSafetyAlert } from "../types"
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

  const buildWarningsMessage = (warnings: MedicationSafetyAlert[]): string => {
    const warningMessages = warnings.map((warning) => buildAlertMessage(t as AlertTranslateFunction, warning))
    return `${t("medicationSafety.warningsSummary", { count: warnings.length })} ${warningMessages.join(" ")}`
  }

  const handleCreateMedication = async (formData: MedicationSubmission): Promise<MedicationSubmitOutcome> => {
    try {
      const createdMedication = await createMedicationMutation.mutateAsync({
        encounter_fhir_id: encounterId,
        patient_fhir_id: patientId,
        medication_name: formData.medicationDisplay,
        dosage_instructions: formData.dosageInstruction,
        overrideJustification: formData.overrideJustification,
      })
      setIsModalOpen(false)
      toast.success(t("toast.medicationSuccess"))
      const warnings = createdMedication?.warnings ?? []
      if (warnings.length > 0) {
        toast.info(buildWarningsMessage(warnings), 10000)
      }
      return { status: "success" }
    } catch (error) {
      const blockedAlerts = extractAllergyBlockedAlerts(error)
      if (blockedAlerts) {
        toast.error(t("medicationSafety.blockedByServer"))
        return { status: "blocked", alerts: blockedAlerts }
      }
      toast.error(t("toast.medicationError"))
      return { status: "failed" }
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
        patientId={patientId}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateMedication}
        isPending={createMedicationMutation.isPending}
      />
    </>
  )
}
