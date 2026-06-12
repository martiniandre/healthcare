import { FileText, Plus, FileCheck, CheckCircle } from "lucide-react"
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

interface DiagnosticReportRepresentation {
  fhir_id: string
  encounter_fhir_id: string
  patient_fhir_id: string
  report_display: string
  status: string
  conclusion: string
  created_at: string
}

interface ClinicalReportsProps {
  reports: DiagnosticReportRepresentation[]
  onAdd: () => void
}

export const ClinicalReports = ({ reports, onAdd }: ClinicalReportsProps) => {
  const { t } = useTranslation("patients")

  return (
    <Card className="flex flex-col gap-5 min-h-[450px]">
      <div className="flex items-center justify-between border-b border-border pb-4">
        <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2">
          <FileText className="w-4 h-4 text-emerald-500" />
          {t("details.reportsCard.title")}
        </h3>
        <Button onClick={onAdd} className="px-3 py-2 text-xs">
          <Plus className="w-3.5 h-3.5" />
          {t("details.reportsCard.add")}
        </Button>
      </div>

      {reports.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center gap-2 py-16">
          <FileText className="w-8 h-8 text-gray-300" />
          <span className="text-xs text-muted">
            {t("details.reportsCard.empty")}
          </span>
        </div>
      ) : (
        <div className="overflow-x-auto w-full">
          <Table className="w-full text-left border-collapse">
            <TableHeader>
              <TableRow className="border-b border-border bg-gray-50/80">
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.reportsCard.display")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.reportsCard.conclusion")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.reportsCard.status")}
                </TableHead>
                <TableHead className="py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider">
                  {t("details.reportsCard.date")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {reports.map((report) => (
                <TableRow key={report.fhir_id} className="border-b border-border/60 hover:bg-gray-50 transition-colors duration-300">
                  <TableCell className="py-4 px-4 align-top">
                    <div className="flex items-center gap-3">
                      <div className="bg-emerald-50 border border-emerald-100 p-2 rounded-lg text-emerald-600">
                        <FileCheck className="w-4 h-4" />
                      </div>
                      <span className="text-sm font-bold text-gray-800 block">
                        {report.report_display}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="py-4 px-4 max-w-xs align-top">
                    <p className="text-xs text-gray-600 leading-relaxed bg-gray-50 border border-border p-3 rounded-lg max-h-24 overflow-y-auto">
                      {report.conclusion}
                    </p>
                  </TableCell>
                  <TableCell className="py-4 px-4 align-top">
                    <span className="text-[9px] bg-emerald-50 border border-emerald-100 text-emerald-600 px-2 py-0.5 rounded font-bold uppercase inline-flex items-center gap-1">
                      <CheckCircle className="w-3 h-3" />
                      {report.status}
                    </span>
                  </TableCell>
                  <TableCell className="py-4 px-4 align-top">
                    <span className="text-xs text-gray-500 font-semibold block mt-1">
                      {new Date(report.created_at).toLocaleString()}
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
