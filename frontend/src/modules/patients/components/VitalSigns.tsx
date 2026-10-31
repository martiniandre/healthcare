import { useState } from "react"
import { Heart, Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/button"
import { ClinicalQueryTable } from "./ClinicalQueryTable"
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
  const { data: observations = [], isLoading, isError, refetch } = useObservationsQuery(encounterId)
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
      <ClinicalQueryTable
        title={t("details.vitalsCard.title")}
        icon={<Heart className="w-4 h-4 text-red-500 animate-pulse-glow" />}
        columns={columns}
        data={observations}
        isLoading={isLoading}
        isError={isError}
        onRetry={() => void refetch()}
        emptyText={t("details.vitalsCard.empty")}
        emptyIcon={Heart}
        addButton={
          <Can I={Action.Create} a={Feature.Observation}>
            <Button size="sm" onClick={() => setIsModalOpen(true)}>
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
