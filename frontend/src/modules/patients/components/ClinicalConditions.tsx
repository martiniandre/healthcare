import { useState } from "react"
import { Activity, Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/button"
import { ClinicalQueryTable } from "./ClinicalQueryTable"
import { ConditionModal } from "./modals/ConditionModal"
import { useClinicalConditionsColumns } from "./useClinicalConditionsColumns"
import { usePatientConditionsQuery, useCreateConditionMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"

interface ClinicalConditionsProps {
  patientId: string
}

export default function ClinicalConditions({ patientId }: ClinicalConditionsProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const columns = useClinicalConditionsColumns()
  const { data: conditions = [], isLoading, isError, refetch } = usePatientConditionsQuery(patientId)
  const createConditionMutation = useCreateConditionMutation()

  const handleCreateCondition = async (formData: { icd10Code: string; codeDisplay: string }) => {
    try {
      await createConditionMutation.mutateAsync({
        patient_fhir_id: patientId,
        icd10_code: formData.icd10Code,
        code_display: formData.codeDisplay,
      })
      setIsModalOpen(false)
      toast.success(t("toast.conditionSuccess"))
    } catch {
      toast.error(t("toast.conditionError"))
    }
  }

  return (
    <>
      <ClinicalQueryTable
        title={t("details.conditionsCard.title")}
        icon={<Activity className="w-4 h-4 text-rose-500" />}
        columns={columns}
        data={conditions}
        isLoading={isLoading}
        isError={isError}
        onRetry={() => void refetch()}
        emptyText={t("details.conditionsCard.empty")}
        emptyIcon={Activity}
        addButton={
          <Can I={Action.Create} a={Feature.Condition}>
            <Button size="sm" onClick={() => setIsModalOpen(true)}>
              <Plus className="w-3.5 h-3.5" />
              {t("details.conditionsCard.add")}
            </Button>
          </Can>
        }
      />

      <ConditionModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateCondition}
        isPending={createConditionMutation.isPending}
      />
    </>
  )
}
