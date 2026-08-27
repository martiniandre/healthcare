import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { History, Plus, AlertTriangle } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { EncounterModal } from "./modals/EncounterModal"
import { useEncounterHistoryColumns } from "./useEncounterHistoryColumns"
import { useCreateEncounterMutation, useEncountersQuery, useUpdateEncounterMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import type { Encounter } from "../types"
import type { NewEncounterFormData } from "../patient_schemas"

interface EncounterHistoryProps {
  patientId: string
}

export default function EncounterHistory({
  patientId,
}: EncounterHistoryProps) {
  const navigate = useNavigate()
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const { data: encounters = [] } = useEncountersQuery(patientId)
  const createEncounterMutation = useCreateEncounterMutation()
  const updateEncounterMutation = useUpdateEncounterMutation()

  const handleCreateEncounter = async (formData: NewEncounterFormData) => {
    try {
      const newEncounter = await createEncounterMutation.mutateAsync({
        patient_fhir_id: patientId,
        reason_code: formData.reasonCode ?? "",
        reason_display: formData.reasonDisplay,
        practitioner_id: formData.practitionerId,
      })
      setIsModalOpen(false)
      navigate(`/patients/${patientId}/encounters/${newEncounter.fhir_id}`)
      toast.success(t("toast.encounterSuccess"))
    } catch {
      toast.error(t("toast.encounterError"))
    }
  }

  const handleFinishEncounter = async (encounter: Encounter) => {
    try {
      await updateEncounterMutation.mutateAsync({
        encounter_fhir_id: encounter.fhir_id,
        patient_fhir_id: patientId,
        status: "finished",
      })
      toast.success(t("toast.encounterFinishSuccess"))
    } catch {
      toast.error(t("toast.encounterFinishError"))
    }
  }

  const handleCancelEncounter = async (encounter: Encounter) => {
    try {
      await updateEncounterMutation.mutateAsync({
        encounter_fhir_id: encounter.fhir_id,
        patient_fhir_id: patientId,
        status: "cancelled",
      })
      toast.success(t("toast.encounterCancelSuccess"))
    } catch {
      toast.error(t("toast.encounterCancelError"))
    }
  }

  const columns = useEncounterHistoryColumns({
    onFinishEncounter: handleFinishEncounter,
    onCancelEncounter: handleCancelEncounter,
    onOpenEncounter: (encounter) => navigate(`/patients/${patientId}/encounters/${encounter.fhir_id}`),
    isActionPending: updateEncounterMutation.isPending,
  })

  return (
    <>
      <ClinicalTable
        title={t("details.encountersCard.title")}
        icon={<History className="w-4 h-4 text-primary animate-pulse-glow" />}
        columns={columns}
        data={encounters}
        isEmpty={encounters.length === 0}
        emptyIcon={<AlertTriangle className="w-8 h-8 text-gray-300" />}
        emptyText={t("details.encountersCard.empty")}
        addButton={
          <Can I={Action.Create} a={Feature.Encounter}>
            <Button onClick={() => setIsModalOpen(true)} className="px-3 py-2 text-xs">
              <Plus className="w-3.5 h-3.5" />
              {t("details.encountersCard.add")}
            </Button>
          </Can>
        }
      />

      <EncounterModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateEncounter}
        isPending={createEncounterMutation.isPending}
      />
    </>
  )
}
