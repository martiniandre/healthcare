import { useState } from "react"
import { ShieldAlert, Plus, CheckCircle } from "lucide-react"
import { useTranslation } from "react-i18next"
import { createColumnHelper } from "@tanstack/react-table"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { AllergyModal } from "./modals/AllergyModal"
import { usePatientAllergiesQuery, useCreateAllergyMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import type { AllergyIntolerance } from "../types"

interface ClinicalAllergiesProps {
  patientId: string
}

const columnHelper = createColumnHelper<AllergyIntolerance>()

export default function ClinicalAllergies({ patientId }: ClinicalAllergiesProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const { data: allergies = [] } = usePatientAllergiesQuery(patientId)
  const createAllergyMutation = useCreateAllergyMutation()

  const handleCreateAllergy = async (formData: { allergenCode: string; allergenDisplay: string; reaction: string }) => {
    try {
      await createAllergyMutation.mutateAsync({
        patient_fhir_id: patientId,
        allergen_code: formData.allergenCode,
        allergen_display: formData.allergenDisplay,
        reaction: formData.reaction,
      })
      setIsModalOpen(false)
      toast.success(t("toast.allergySuccess"))
    } catch {
      toast.error(t("toast.allergyError"))
    }
  }

  const columns = [
    columnHelper.accessor("allergen_code", {
      header: t("details.allergiesCard.code"),
      cell: (info) => (
        <span className="text-xs font-mono font-bold text-gray-700 bg-amber-50 border border-amber-100 px-2 py-1 rounded-md">
          {info.getValue()}
        </span>
      ),
    }),
    columnHelper.accessor("allergen_display", {
      header: t("details.allergiesCard.allergen"),
      cell: (info) => <span className="text-sm font-bold text-gray-800 block">{info.getValue()}</span>,
    }),
    columnHelper.accessor("reaction", {
      header: t("details.allergiesCard.reaction"),
      cell: (info) => <span className="text-sm font-semibold text-red-600 block">{info.getValue()}</span>,
    }),
    columnHelper.accessor("clinical_status", {
      header: t("details.allergiesCard.status"),
      cell: (info) => (
        <span className="text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
          <CheckCircle className="w-3 h-3" />
          {info.getValue()}
        </span>
      ),
    }),
    columnHelper.accessor("created_at", {
      header: t("details.allergiesCard.date"),
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
        title={t("details.allergiesCard.title")}
        icon={<ShieldAlert className="w-4 h-4 text-amber-500 animate-pulse" />}
        columns={columns}
        data={allergies}
        isEmpty={allergies.length === 0}
        emptyIcon={<ShieldAlert className="w-8 h-8 text-gray-300" />}
        emptyText={t("details.allergiesCard.empty")}
        addButton={
          <Can I={Action.Create} a={Feature.Allergy}>
            <Button onClick={() => setIsModalOpen(true)} className="px-3 py-2 text-xs">
              <Plus className="w-3.5 h-3.5" />
              {t("details.allergiesCard.add")}
            </Button>
          </Can>
        }
      />

      <AllergyModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateAllergy}
        isPending={createAllergyMutation.isPending}
      />
    </>
  )
}
