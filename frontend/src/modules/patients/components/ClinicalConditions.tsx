import { useState } from "react"
import { Activity, Plus, ShieldAlert, CheckCircle } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "../../../shared/components/ui/Table"
import { ConditionModal } from "./modals/ConditionModal"
import { usePatientConditionsQuery, useCreateConditionMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"

interface ClinicalConditionsProps {
  patientId: string
}

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

  return (
    <>
      <Card className="flex flex-col gap-5 min-h-[450px]">
        <div className="flex items-center justify-between border-b border-border pb-4">
          <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2">
            <Activity className="w-4 h-4 text-rose-500" />
            {t("details.conditionsCard.title")}
          </h3>
          <Button onClick={() => setIsModalOpen(true)} className="px-3 py-2 text-xs">
            <Plus className="w-3.5 h-3.5" />
            {t("details.conditionsCard.add")}
          </Button>
        </div>

        {conditions.length === 0 ? (
          <div className="flex-1 flex flex-col items-center justify-center gap-2 py-16">
            <ShieldAlert className="w-8 h-8 text-gray-300" />
            <span className="text-xs text-muted">
              {t("details.conditionsCard.empty")}
            </span>
          </div>
        ) : (
          <div className="overflow-x-auto w-full">
            <Table className="w-full text-left border-collapse">
              <TableHeader>
                <TableRow className="border-b border-border bg-gray-50/80">
                  <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                    {t("details.conditionsCard.code")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                    {t("details.conditionsCard.display")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                    {t("details.conditionsCard.status")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                    {t("details.conditionsCard.date")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {conditions.map((conditionItem) => (
                  <TableRow key={conditionItem.fhir_id} className="border-b border-border/60 hover:bg-gray-50 transition-colors duration-300">
                    <TableCell className="py-4 px-4 align-top">
                      <span className="text-sm font-extrabold text-gray-900 bg-rose-50 border border-rose-100 px-2 py-1 rounded-md">
                        {conditionItem.icd10_code}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4 align-top">
                      <span className="text-sm font-bold text-gray-800 block">
                        {conditionItem.code_display}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4 align-top">
                      <span className="text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
                        <CheckCircle className="w-3 h-3" />
                        {conditionItem.clinical_status}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4 align-top">
                      <span className="text-xs text-gray-500 font-semibold block">
                        {new Date(conditionItem.created_at).toLocaleString()}
                      </span>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </Card>

      <ConditionModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateCondition}
        isPending={createConditionMutation.isPending}
      />
    </>
  )
}
