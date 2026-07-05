import { useState } from "react"
import { History, Plus, AlertTriangle, CheckCircle } from "lucide-react"
import { useTranslation } from "react-i18next"
import { createColumnHelper } from "@tanstack/react-table"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { EncounterModal } from "./modals/EncounterModal"
import { useCreateEncounterMutation, useEncountersQuery } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import type { Encounter } from "../types"

interface EncounterHistoryProps {
  patientId: string
  selectedEncounterId: string | null
  onSelect: (id: string) => void
}

const columnHelper = createColumnHelper<Encounter>()

export default function EncounterHistory({
  patientId,
  selectedEncounterId,
  onSelect,
}: EncounterHistoryProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const { data: encounters = [] } = useEncountersQuery(patientId)
  const createEncounterMutation = useCreateEncounterMutation()

  const handleCreateEncounter = async (formData: { reasonDisplay: string; practitionerId?: string }) => {
    try {
      const newEncounter = await createEncounterMutation.mutateAsync({
        patient_fhir_id: patientId,
        reason_display: formData.reasonDisplay,
        practitioner_id: formData.practitionerId || undefined,
      })
      setIsModalOpen(false)
      onSelect(newEncounter.fhir_id)
      toast.success(t("toast.encounterSuccess"))
    } catch {
      toast.error(t("toast.encounterError"))
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
          {new Date(info.getValue()).toLocaleString()}
        </span>
      ),
    }),
    columnHelper.display({
      id: "actions",
      header: t("details.encountersCard.action"),
      cell: (info) => {
        const encounter = info.row.original
        const isActive = selectedEncounterId === encounter.fhir_id
        return (
          <div className="text-right pr-6">
            <Button
              variantType={isActive ? "primary" : "outline"}
              onClick={() => onSelect(encounter.fhir_id)}
              className="px-2.5 py-1 text-[10px] font-bold gap-1"
            >
              {isActive && <CheckCircle className="w-3 h-3 text-white" />}
              {t("details.focus")}
            </Button>
          </div>
        )
      },
    }),
  ]

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
