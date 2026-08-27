import { useState } from "react"
import { Heart, Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { ObservationModal } from "./modals/ObservationModal"
import { useVitalSignsColumns } from "./useVitalSignsColumns"
import { useObservationsQuery, useCreateVitalSignsPanelMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import type { NewVitalSignsPanelFormData } from "../patient_schemas"

interface VitalSignsProps {
  patientId: string
  encounterId: string
}

export default function VitalSigns({ patientId, encounterId }: VitalSignsProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const columns = useVitalSignsColumns()
  const { data: observations = [] } = useObservationsQuery(encounterId)
  const createVitalSignsPanelMutation = useCreateVitalSignsPanelMutation()

  const handleCreateVitalSignsPanel = async (panelFormData: NewVitalSignsPanelFormData) => {
    try {
      await createVitalSignsPanelMutation.mutateAsync({
        encounter_fhir_id: encounterId,
        patient_fhir_id: patientId,
        panel_form_data: panelFormData,
      })
      setIsModalOpen(false)
      toast.success(t("toast.observationSuccess"))
    } catch {
      toast.error(t("toast.observationError"))
    }
  }

  return (
    <>
      <ClinicalTable
        title={t("details.vitalsCard.title")}
        icon={<Heart className="w-4 h-4 text-red-500 animate-pulse-glow" />}
        columns={columns}
        data={observations}
        isEmpty={observations.length === 0}
        emptyIcon={<Heart className="w-8 h-8 text-gray-300" />}
        emptyText={t("details.vitalsCard.empty")}
        addButton={
          <Can I={Action.Create} a={Feature.Observation}>
            <Button onClick={() => setIsModalOpen(true)} className="px-3 py-2 text-xs">
              <Plus className="w-3.5 h-3.5" />
              {t("details.vitalsCard.add")}
            </Button>
          </Can>
        }
      />

      <ObservationModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateVitalSignsPanel}
        isPending={createVitalSignsPanelMutation.isPending}
      />
    </>
  )
}
