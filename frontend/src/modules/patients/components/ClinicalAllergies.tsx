import { useState } from "react"
import { ShieldAlert, Plus, CheckCircle } from "lucide-react"
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
import { AllergyModal } from "./modals/AllergyModal"
import { usePatientAllergiesQuery, useCreateAllergyMutation } from "../queries"
import { toast } from "../../../shared/store/toast_store"

interface ClinicalAllergiesProps {
  patientId: string
}

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

  return (
    <>
      <Card className="flex flex-col gap-5 min-h-[450px]">
        <div className="flex items-center justify-between border-b border-border pb-4">
          <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2">
            <ShieldAlert className="w-4 h-4 text-amber-500 animate-pulse" />
            {t("details.allergiesCard.title")}
          </h3>
          <Button onClick={() => setIsModalOpen(true)} className="px-3 py-2 text-xs">
            <Plus className="w-3.5 h-3.5" />
            {t("details.allergiesCard.add")}
          </Button>
        </div>

        {allergies.length === 0 ? (
          <div className="flex-1 flex flex-col items-center justify-center gap-2 py-16">
            <ShieldAlert className="w-8 h-8 text-gray-300" />
            <span className="text-xs text-muted">
              {t("details.allergiesCard.empty")}
            </span>
          </div>
        ) : (
          <div className="overflow-x-auto w-full">
            <Table className="w-full text-left border-collapse">
              <TableHeader>
                <TableRow className="border-b border-border bg-gray-50/80">
                  <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                    {t("details.allergiesCard.code")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                    {t("details.allergiesCard.allergen")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                    {t("details.allergiesCard.reaction")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                    {t("details.allergiesCard.status")}
                  </TableHead>
                  <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                    {t("details.allergiesCard.date")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {allergies.map((allergyItem) => (
                  <TableRow key={allergyItem.fhir_id} className="border-b border-border/60 hover:bg-gray-50 transition-colors duration-300">
                    <TableCell className="py-4 px-4 align-top">
                      <span className="text-xs font-mono font-bold text-gray-700 bg-amber-50 border border-amber-100 px-2 py-1 rounded-md">
                        {allergyItem.allergen_code}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4 align-top">
                      <span className="text-sm font-bold text-gray-800 block">
                        {allergyItem.allergen_display}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4 align-top">
                      <span className="text-sm font-semibold text-red-600 block">
                        {allergyItem.reaction}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4 align-top">
                      <span className="text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
                        <CheckCircle className="w-3 h-3" />
                        {allergyItem.clinical_status}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4 align-top">
                      <span className="text-xs text-gray-500 font-semibold block">
                        {new Date(allergyItem.created_at).toLocaleString()}
                      </span>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </Card>

      <AllergyModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateAllergy}
        isPending={createAllergyMutation.isPending}
      />
    </>
  )
}
