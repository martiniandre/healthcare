import { useState } from "react"
import { Activity, Plus, ShieldAlert, CheckCircle } from "lucide-react"
import { useTranslation } from "react-i18next"
import { createColumnHelper } from "@tanstack/react-table"
import { Can, Action, Feature } from "../../../shared/auth/AbilityContext"
import { Button } from "../../../shared/components/ui/Button"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { ConditionModal } from "./modals/ConditionModal"
import { usePatientConditionsQuery, useCreateConditionMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"
import type { Condition } from "../types"

interface ClinicalConditionsProps {
  patientId: string
}

const columnHelper = createColumnHelper<Condition>()

export default function ClinicalConditions({ patientId }: ClinicalConditionsProps) {
  const { t } = useTranslation("patients")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const { data: conditions = [] } = usePatientConditionsQuery(patientId)
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

  const columns = [
    columnHelper.accessor("icd10_code", {
      header: t("details.conditionsCard.code"),
      cell: (info) => (
        <span className="text-sm font-extrabold text-gray-900 bg-rose-50 border border-rose-100 px-2 py-1 rounded-md">
          {info.getValue()}
        </span>
      ),
    }),
    columnHelper.accessor("code_display", {
      header: t("details.conditionsCard.display"),
      cell: (info) => <span className="text-sm font-bold text-gray-800 block">{info.getValue()}</span>,
    }),
    columnHelper.accessor("clinical_status", {
      header: t("details.conditionsCard.status"),
      cell: (info) => (
        <span className="text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
          <CheckCircle className="w-3 h-3" />
          {info.getValue()}
        </span>
      ),
    }),
    columnHelper.accessor("created_at", {
      header: t("details.conditionsCard.date"),
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
        title={t("details.conditionsCard.title")}
        icon={<Activity className="w-4 h-4 text-rose-500" />}
        columns={columns}
        data={conditions}
        isEmpty={conditions.length === 0}
        emptyIcon={<ShieldAlert className="w-8 h-8 text-gray-300" />}
        emptyText={t("details.conditionsCard.empty")}
        addButton={
          <Can I={Action.Create} a={Feature.Condition}>
            <Button onClick={() => setIsModalOpen(true)} className="px-3 py-2 text-xs">
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
