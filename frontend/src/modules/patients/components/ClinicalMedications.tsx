import { Pill, Plus, ShieldAlert, CheckCircle } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"
import type { MedicationRequest } from "../types"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "../../../shared/components/ui/Table"

interface ClinicalMedicationsProps {
  medications: MedicationRequest[]
  onAdd: () => void
}

export const ClinicalMedications = ({ medications, onAdd }: ClinicalMedicationsProps) => {
  const { t } = useTranslation("patients")

  return (
    <Card className="flex flex-col gap-5 min-h-[450px]">
      <div className="flex items-center justify-between border-b border-border pb-4">
        <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2">
          <Pill className="w-4 h-4 text-purple-500" />
          {t("details.medicationsCard.title")}
        </h3>
        <Button onClick={onAdd} className="px-3 py-2 text-xs">
          <Plus className="w-3.5 h-3.5" />
          {t("details.medicationsCard.add")}
        </Button>
      </div>

      {medications.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center gap-2 py-16">
          <ShieldAlert className="w-8 h-8 text-gray-300" />
          <span className="text-xs text-muted">
            {t("details.medicationsCard.empty")}
          </span>
        </div>
      ) : (
        <div className="overflow-x-auto w-full">
          <Table className="w-full text-left border-collapse">
            <TableHeader>
              <TableRow className="border-b border-border bg-gray-50/80">
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.medicationsCard.display")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.medicationsCard.dosage")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.medicationsCard.status")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.medicationsCard.date")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {medications.map((medicationItem) => (
                <TableRow key={medicationItem.fhir_id} className="border-b border-border/60 hover:bg-gray-50 transition-colors duration-300">
                  <TableCell className="py-4 px-4 align-top">
                    <span className="text-sm font-extrabold text-gray-900 block">
                      {medicationItem.medication_display}
                    </span>
                  </TableCell>
                  <TableCell className="py-4 px-4 align-top">
                    <span className="text-sm font-bold text-gray-800 block whitespace-pre-line">
                      {medicationItem.dosage_instruction}
                    </span>
                  </TableCell>
                  <TableCell className="py-4 px-4 align-top">
                    <span className="text-[9px] bg-purple-50 border border-purple-100 text-purple-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
                      <CheckCircle className="w-3 h-3" />
                      {medicationItem.status}
                    </span>
                  </TableCell>
                  <TableCell className="py-4 px-4 align-top">
                    <span className="text-xs text-gray-500 font-semibold block">
                      {new Date(medicationItem.created_at).toLocaleString()}
                    </span>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </Card>
  )
}
