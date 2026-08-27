import { useState } from "react"
import { ShieldAlert, Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { AllergyModal } from "./modals/AllergyModal"
import { useClinicalAllergiesColumns } from "./useClinicalAllergiesColumns"
import { usePatientAllergiesQuery, useCreateAllergyMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"

interface ClinicalAllergiesProps {
  patientId: string
}

export default function ClinicalAllergies({ patientId }: ClinicalAllergiesProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const columns = useClinicalAllergiesColumns()
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
