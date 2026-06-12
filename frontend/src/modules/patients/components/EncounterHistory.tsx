import { History, Plus, AlertTriangle, CheckCircle } from "lucide-react"
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

interface EncounterRepresentation {
  fhir_id: string
  patient_fhir_id: string
  status: string
  reason_display: string
  practitioner_id?: string
  created_at: string
}

interface EncounterHistoryProps {
  encounters: EncounterRepresentation[]
  selectedEncounterId: string | null
  onSelect: (id: string) => void
  onNew: () => void
}

export const EncounterHistory = ({
  encounters,
  selectedEncounterId,
  onSelect,
  onNew
}: EncounterHistoryProps) => {
  const { t } = useTranslation("patients")

  return (
    <Card className="flex flex-col gap-5 min-h-[450px]">
      <div className="flex items-center justify-between border-b border-border pb-4">
        <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2">
          <History className="w-4 h-4 text-primary animate-pulse-glow" />
          {t("details.encountersCard.title")}
        </h3>
        <Button onClick={onNew} className="px-3 py-2 text-xs">
          <Plus className="w-3.5 h-3.5" />
          {t("details.encountersCard.add")}
        </Button>
      </div>

      {encounters.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center gap-2 py-16">
          <AlertTriangle className="w-8 h-8 text-gray-300" />
          <span className="text-xs text-gray-500 font-bold">
            {t("details.encountersCard.empty")}
          </span>
        </div>
      ) : (
        <div className="overflow-x-auto w-full">
          <Table className="w-full text-left border-collapse">
            <TableHeader>
              <TableRow className="border-b border-border bg-gray-50/80">
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.encountersCard.reason")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.encountersCard.status")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.encountersCard.date")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-right text-xs font-black text-gray-400 uppercase tracking-wider pr-6">
                  {t("details.encountersCard.action")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {encounters.map((encounter) => {
                const isActive = selectedEncounterId === encounter.fhir_id
                return (
                  <TableRow
                    key={encounter.fhir_id}
                    className={`border-b border-border/60 transition-colors duration-300 ${
                      isActive 
                        ? "bg-primary/5 border-primary/20" 
                        : "hover:bg-gray-50"
                    }`}
                  >
                    <TableCell className="py-4 px-4">
                      <span className="text-sm font-bold text-gray-800 block">
                        {encounter.reason_display}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4">
                      <span className="text-[10px] bg-gray-100 text-gray-600 px-2.5 py-1 rounded font-bold uppercase tracking-wider">
                        {encounter.status}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4">
                      <span className="text-xs text-gray-400 font-semibold">
                        {new Date(encounter.created_at).toLocaleString()}
                      </span>
                    </TableCell>
                    <TableCell className="py-4 px-4 text-right pr-6">
                      <Button
                        variantType={isActive ? "primary" : "outline"}
                        onClick={() => onSelect(encounter.fhir_id)}
                        className="px-2.5 py-1 text-[10px] font-bold gap-1"
                      >
                        {isActive && <CheckCircle className="w-3 h-3 text-white" />}
                        {t("details.focus")}
                      </Button>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        </div>
      )}
    </Card>
  )
}
