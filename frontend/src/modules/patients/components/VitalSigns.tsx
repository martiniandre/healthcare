import { useState } from "react"
import { Activity, Heart, Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { createColumnHelper, type ColumnDef } from "@tanstack/react-table"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { ObservationModal } from "./modals/ObservationModal"
import { findVitalSignDisplay } from "./vitalSignDisplay"
import { VitalSignValueDisplay } from "./VitalSignValueDisplay"
import { useObservationsQuery, useCreateVitalSignsPanelMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import { formatDateTime } from "../../../shared/utils/dates"
import type { Observation } from "../types"
import type { NewVitalSignsPanelFormData } from "../patient_schemas"

interface VitalSignsProps {
  patientId: string
  encounterId: string
}

const columnHelper = createColumnHelper<Observation>()

export default function VitalSigns({ patientId, encounterId }: VitalSignsProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
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

  const columns = [
    columnHelper.accessor("code_display", {
      header: t("details.vitalsCard.display"),
      cell: (info) => {
        const observation = info.row.original
        const displayMetadata = findVitalSignDisplay(observation.loinc_code)
        const IconComponent = displayMetadata?.IconComponent ?? Activity
        return (
          <div className="flex items-center gap-3">
            <div className={`p-2 rounded-lg border ${displayMetadata?.iconClassName ?? "bg-blue-50 border-blue-100 text-blue-600"}`}>
              <IconComponent className="w-4 h-4" />
            </div>
            <span className="text-sm font-bold text-gray-800 block">
              {displayMetadata?.labelKey ? t(displayMetadata.labelKey) : info.getValue()}
            </span>
          </div>
        )
      },
    }),
    columnHelper.accessor("loinc_code", {
      header: t("details.vitalsCard.code"),
      cell: (info) => <span className="text-xs font-mono text-gray-500">{info.getValue()}</span>,
    }),
    columnHelper.accessor("value_quantity", {
      header: t("details.vitalsCard.value"),
      cell: (info) => {
        const observation = info.row.original
        return (
          <VitalSignValueDisplay
            notPerformed={observation.not_performed}
            valueQuantity={info.getValue()}
            valueUnit={observation.value_unit}
          />
        )
      },
    }),
    columnHelper.accessor("created_at", {
      header: t("details.vitalsCard.date"),
      cell: (info) => (
        <span className="text-xs text-gray-500 font-semibold">{formatDateTime(info.getValue())}</span>
      ),
    }),
  ] as ColumnDef<Observation>[]

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
