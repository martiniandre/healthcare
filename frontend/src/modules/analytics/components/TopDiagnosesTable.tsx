import { Card } from "../../../shared/components/ui/Card"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "../../../shared/components/ui/Table"
import type { DiagnosisCount } from "../dashboard_types"

interface TopDiagnosesTableProps {
  topDiagnoses: DiagnosisCount[]
}

export const TopDiagnosesTable = ({ topDiagnoses }: TopDiagnosesTableProps) => {
  const maxCount = topDiagnoses.length > 0
    ? Math.max(...topDiagnoses.map((diagnosis) => diagnosis.count))
    : 1

  return (
    <Card className="p-5 flex flex-col gap-4 text-left border border-border">
      <div>
        <h3 className="font-extrabold text-gray-900 text-md">Principais Diagnósticos</h3>
        <span className="text-xs text-muted block mt-1">CID-10 mais frequentes nos últimos 30 dias</span>
      </div>

      <div className="overflow-x-auto w-full">
        {topDiagnoses.length > 0 ? (
          <Table className="min-w-[450px]">
            <TableHeader>
              <TableRow>
                <TableHead>Código</TableHead>
                <TableHead>Descrição</TableHead>
                <TableHead>Casos</TableHead>
                <TableHead className="text-right">Proporção</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody className="text-gray-700 font-medium">
              {topDiagnoses.map((diagnosis) => {
                const proportion = (diagnosis.count / maxCount) * 100
                return (
                  <TableRow key={diagnosis.icd10_code}>
                    <TableCell className="font-mono font-bold text-primary">{diagnosis.icd10_code}</TableCell>
                    <TableCell>{diagnosis.description}</TableCell>
                    <TableCell className="font-bold text-gray-900">{diagnosis.count}</TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-2">
                        <div className="w-20 bg-gray-100 rounded-full h-2 overflow-hidden">
                          <div
                            className="bg-primary rounded-full h-full"
                            style={{ width: `${proportion}%` }}
                          />
                        </div>
                        <span className="text-[10px] text-gray-500 font-semibold w-8">{proportion.toFixed(0)}%</span>
                      </div>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        ) : (
          <EmptyState
            title="Sem dados"
            description="Nenhum diagnóstico registrado no período."
          />
        )}
      </div>
    </Card>
  )
}
