import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { History, Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/button"
import { ClinicalQueryTable } from "./ClinicalQueryTable"
import { EncounterModal } from "./modals/EncounterModal"
import { useEncounterHistoryColumns } from "./useEncounterHistoryColumns"
import { useCreateEncounterMutation, useEncountersQuery, useUpdateEncounterMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import { ConfirmActionDialog } from "../../../shared/components/ui/confirm-action-dialog"
import { formatDateTime } from "../../../shared/utils/dates"
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
  const [pendingStatusChange, setPendingStatusChange] = useState<{ encounter: Encounter; nextStatus: "finished" | "cancelled" } | null>(null)
  const { data: encounters = [], isLoading, isError, refetch } = useEncountersQuery(patientId)
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

  const handleConfirmStatusChange = async () => {
    if (!pendingStatusChange) return
    const { encounter, nextStatus } = pendingStatusChange
    const isFinishing = nextStatus === "finished"
    try {
      await updateEncounterMutation.mutateAsync({
        encounter_fhir_id: encounter.fhir_id,
        patient_fhir_id: patientId,
        status: nextStatus,
      })
      toast.success(t(isFinishing ? "toast.encounterFinishSuccess" : "toast.encounterCancelSuccess"))
      setPendingStatusChange(null)
    } catch {
      toast.error(t(isFinishing ? "toast.encounterFinishError" : "toast.encounterCancelError"))
    }
  }

  const isFinishingEncounter = pendingStatusChange?.nextStatus === "finished"

  const columns = useEncounterHistoryColumns({
    onFinishEncounter: (encounter) => setPendingStatusChange({ encounter, nextStatus: "finished" }),
    onCancelEncounter: (encounter) => setPendingStatusChange({ encounter, nextStatus: "cancelled" }),
    onOpenEncounter: (encounter) => navigate(`/patients/${patientId}/encounters/${encounter.fhir_id}`),
    isActionPending: updateEncounterMutation.isPending,
  })

  return (
    <>
      <ClinicalQueryTable
        title={t("details.encountersCard.title")}
        columns={columns}
        data={encounters}
        isLoading={isLoading}
        isError={isError}
        onRetry={() => void refetch()}
        emptyText={t("details.encountersCard.empty")}
        emptyIcon={History}
        addButton={
          <Can I={Action.Create} a={Feature.Encounter}>
            <Button size="sm" onClick={() => setIsModalOpen(true)}>
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
        hasSubmitError={createEncounterMutation.isError}
      />

      <ConfirmActionDialog
        isOpen={pendingStatusChange !== null}
        title={t(isFinishingEncounter ? "details.encountersCard.confirm.finishTitle" : "details.encountersCard.confirm.cancelTitle")}
        description={t(isFinishingEncounter ? "details.encountersCard.confirm.finishDescription" : "details.encountersCard.confirm.cancelDescription")}
        details={
          pendingStatusChange && (
            <>
              <span className="font-medium text-foreground">{pendingStatusChange.encounter.reason_display}</span>
              <span className="text-xs text-muted-foreground">{formatDateTime(pendingStatusChange.encounter.created_at)}</span>
            </>
          )
        }
        hint={t("details.encountersCard.confirm.irreversible")}
        confirmLabel={t(isFinishingEncounter ? "details.encountersCard.confirm.finishConfirm" : "details.encountersCard.confirm.cancelConfirm")}
        cancelLabel={t("details.encountersCard.confirm.back")}
        isDestructive={!isFinishingEncounter}
        isPending={updateEncounterMutation.isPending}
        onConfirm={() => void handleConfirmStatusChange()}
        onClose={() => setPendingStatusChange(null)}
      />
    </>
  )
}
