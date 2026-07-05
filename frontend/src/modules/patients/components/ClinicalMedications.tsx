import { useState } from "react"
import { Pill, Plus, ShieldAlert, CheckCircle } from "lucide-react"
import { useTranslation } from "react-i18next"
import { createColumnHelper } from "@tanstack/react-table"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { MedicationModal } from "./modals/MedicationModal"
import { useMedicationsQuery, useCreateMedicationMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import type { MedicationRequest } from "../types"

interface ClinicalMedicationsProps {
  patientId: string
  encounterId: string
}

const columnHelper = createColumnHelper<MedicationRequest>()

export default function ClinicalMedications({ patientId, encounterId }: ClinicalMedicationsProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const { data: medications = [] } = useMedicationsQuery(encounterId)
  const createMedicationMutation = useCreateMedicationMutation()

  const handleCreateMedication = async (formData: { medicationDisplay: string; dosageInstruction: string }) => {
    try {
      await createMedicationMutation.mutateAsync({
        encounter_fhir_id: encounterId,
        patient_fhir_id: patientId,
        medication_display: formData.medicationDisplay,
        dosage_instruction: formData.dosageInstruction,
      })
      setIsModalOpen(false)
      toast.success(t("toast.medicationSuccess"))
    } catch {
      toast.error(t("toast.medicationError"))
    }
  }

  const columns = [
    columnHelper.accessor("medication_display", {
      header: t("details.medicationsCard.display"),
      cell: (info) => <span className="text-sm font-extrabold text-gray-900 block">{info.getValue()}</span>,
    }),
    columnHelper.accessor("dosage_instruction", {
      header: t("details.medicationsCard.dosage"),
      cell: (info) => (
        <span className="text-sm font-bold text-gray-800 block whitespace-pre-line">{info.getValue()}</span>
      ),
    }),
    columnHelper.accessor("status", {
      header: t("details.medicationsCard.status"),
      cell: (info) => (
        <span className="text-[9px] bg-purple-50 border border-purple-100 text-purple-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
          <CheckCircle className="w-3 h-3" />
          {info.getValue()}
        </span>
      ),
    }),
    columnHelper.accessor("created_at", {
      header: t("details.medicationsCard.date"),
      cell: (info) => (
        <span className="text-xs text-gray-500 font-semibold block">
          {new Date(info.getValue()).toLocaleString()}
        </span>
      ),
    }),
  ]

  return (
    <>
      <ClinicalTable
        title={t("details.medicationsCard.title")}
        icon={<Pill className="w-4 h-4 text-purple-500" />}
        columns={columns}
        data={medications}
        isEmpty={medications.length === 0}
        emptyIcon={<ShieldAlert className="w-8 h-8 text-gray-300" />}
        emptyText={t("details.medicationsCard.empty")}
        addButton={
          <Can I={Action.Create} a={Feature.MedicationRequest}>
            <Button onClick={() => setIsModalOpen(true)} className="px-3 py-2 text-xs">
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
