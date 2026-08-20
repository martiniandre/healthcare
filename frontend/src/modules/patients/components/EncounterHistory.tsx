import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { History, Plus, AlertTriangle, ExternalLink, CheckCircle, XCircle } from "lucide-react"
import { useTranslation } from "react-i18next"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { EncounterModal } from "./modals/EncounterModal"
import { useCreateEncounterMutation, useEncountersQuery, useUpdateEncounterMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import { formatDateTime } from "../../../shared/utils/dates"
import type { Encounter } from "../types"
import type { NewEncounterFormData } from "../patient_schemas"

interface EncounterHistoryProps {
  patientId: string
}

const columnHelper = createColumnHelper<Encounter>()

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

  const columns = [
    columnHelper.accessor("reason_display", {
      header: t("details.encountersCard.reason"),
      cell: (info) => <span className="text-sm font-bold text-gray-800 block">{info.getValue()}</span>,
    }),
    columnHelper.accessor("status", {
      header: t("details.encountersCard.status"),
      cell: (info) => (
        <span className="text-[10px] bg-gray-100 text-gray-600 px-2.5 py-1 rounded font-bold uppercase tracking-wider">
          {info.getValue()}
        </span>
      ),
    }),
    columnHelper.accessor("created_at", {
      header: t("details.encountersCard.date"),
      cell: (info) => (
        <span className="text-xs text-gray-400 font-semibold">
          {formatDateTime(info.getValue())}
        </span>
      ),
    }),
    columnHelper.display({
      id: "actions",
      header: t("details.encountersCard.action"),
      cell: (info) => {
        const encounter = info.row.original
        return (
          <div className="flex justify-end gap-1.5 pr-6">
            {encounter.status === "in-progress" && (
              <>
                <Button
                  variantType="outline"
                  onClick={() => handleFinishEncounter(encounter)}
                  disabled={updateEncounterMutation.isPending}
                  className="px-2.5 py-1 text-[10px] font-bold gap-1"
                >
                  <CheckCircle className="w-3 h-3" />
                  {t("details.encountersCard.finish")}
                </Button>
                <Button
                  variantType="danger"
                  onClick={() => handleCancelEncounter(encounter)}
                  disabled={updateEncounterMutation.isPending}
                  className="px-2.5 py-1 text-[10px] font-bold gap-1"
                >
                  <XCircle className="w-3 h-3" />
                  {t("details.encountersCard.cancel")}
                </Button>
              </>
            )}
            <Button
              variantType="outline"
              onClick={() => navigate(`/patients/${patientId}/encounters/${encounter.fhir_id}`)}
              className="px-2.5 py-1 text-[10px] font-bold gap-1"
            >
              <ExternalLink className="w-3 h-3" />
              {t("details.openEncounter")}
            </Button>
          </div>
        )
      },
    }),
  ] as ColumnDef<Encounter>[]

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
