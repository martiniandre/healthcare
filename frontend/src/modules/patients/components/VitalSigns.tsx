import { useState } from "react"
import { Heart, Thermometer, Activity, Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { createColumnHelper } from "@tanstack/react-table"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { ObservationModal } from "./modals/ObservationModal"
import { useObservationsQuery, useCreateObservationMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import type { Observation } from "../types"

interface VitalSignsProps {
  patientId: string
  encounterId: string
}

const columnHelper = createColumnHelper<Observation>()

export default function VitalSigns({ patientId, encounterId }: VitalSignsProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const { data: observations = [] } = useObservationsQuery(encounterId)
  const createObservationMutation = useCreateObservationMutation()

  const handleCreateObservation = async (formData: { loincCode: string; valueQuantity: number }) => {
    let display = "Frequência Cardíaca"
    let unit = "bpm"

    if (formData.loincCode === "8310-5") {
      display = "Temperatura Corporal"
      unit = "°C"
    } else if (formData.loincCode === "85354-9") {
      display = "Pressão Arterial Sistólica"
      unit = "mmHg"
    }

    try {
      await createObservationMutation.mutateAsync({
        encounter_fhir_id: encounterId,
        patient_fhir_id: patientId,
        loinc_code: formData.loincCode,
        code_display: display,
        value_quantity: formData.valueQuantity,
        value_unit: unit,
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
        const isHeartRate = observation.loinc_code === "8867-4"
        const isTemp = observation.loinc_code === "8310-5"
        return (
          <div className="flex items-center gap-3">
            <div className={`p-2 rounded-lg border ${
              isHeartRate
                ? "bg-red-50 border-red-100 text-red-600"
                : isTemp
                  ? "bg-amber-50 border-amber-100 text-amber-600"
                  : "bg-blue-50 border-blue-100 text-blue-600"
            }`}>
              {isHeartRate ? <Heart className="w-4 h-4" /> : isTemp ? <Thermometer className="w-4 h-4" /> : <Activity className="w-4 h-4" />}
            </div>
            <span className="text-sm font-bold text-gray-800 block">{info.getValue()}</span>
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
          <span className="text-sm font-extrabold text-gray-800">
            {info.getValue()}
            <span className="text-xs text-muted font-normal ml-1">{observation.value_unit}</span>
          </span>
        )
      },
    }),
    columnHelper.accessor("created_at", {
      header: t("details.vitalsCard.date"),
      cell: (info) => (
        <span className="text-xs text-gray-500 font-semibold">{new Date(info.getValue()).toLocaleString()}</span>
      ),
    }),
  ]

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
        onSubmit={handleCreateObservation}
        isPending={createObservationMutation.isPending}
      />
    </>
  )
}
