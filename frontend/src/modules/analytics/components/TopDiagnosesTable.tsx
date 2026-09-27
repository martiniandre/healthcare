import { useMemo } from "react"
import { useTranslation } from "react-i18next"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../../shared/components/ui/Table"
import { DashboardPanel, PanelMetric } from "./DashboardPanel"
import { toneFillVariables } from "../dashboard_thresholds"
import type { DiagnosisCount } from "../dashboard_types"

interface RankedDiagnosis extends DiagnosisCount {
  rank: number
  proportion: number
}

interface TopDiagnosesTableProps {
  topDiagnoses: DiagnosisCount[]
}

export const TopDiagnosesTable = ({ topDiagnoses }: TopDiagnosesTableProps) => {
  const { t } = useTranslation("analytics")

  const rankedDiagnoses = useMemo<RankedDiagnosis[]>(() => {
    const peakCount =
      topDiagnoses.length > 0
        ? Math.max(...topDiagnoses.map((diagnosis) => diagnosis.count))
        : 0

    return [...topDiagnoses]
      .map((diagnosis) => ({
        ...diagnosis,
        rank: 0,
        proportion: peakCount > 0 ? (diagnosis.count / peakCount) * 100 : 0,
      }))
      .sort((firstDiagnosis, secondDiagnosis) => secondDiagnosis.count - firstDiagnosis.count)
      .map((diagnosis, indexValue) => ({ ...diagnosis, rank: indexValue + 1 }))
  }, [topDiagnoses])

  return (
    <DashboardPanel
      title={t("dashboard.topDiagnoses.title")}
      subtitle={t("dashboard.topDiagnoses.subtitle")}
      trailing={
        rankedDiagnoses.length > 0 ? (
          <PanelMetric
            label={t("dashboard.topDiagnoses.codeCountLabel")}
            value={t("dashboard.topDiagnoses.codeCountUnit", { count: rankedDiagnoses.length })}
          />
        ) : null
      }
    >
      {rankedDiagnoses.length > 0 ? (
        <div className="-mx-5 -my-5 overflow-x-auto">
          <Table className="min-w-[560px]">
            <TableHeader>
              <TableRow>
                <TableHead className="w-14 pl-5">{t("dashboard.topDiagnoses.rank")}</TableHead>
                <TableHead className="w-32">{t("dashboard.topDiagnoses.code")}</TableHead>
                <TableHead>{t("dashboard.topDiagnoses.description")}</TableHead>
                <TableHead className="w-24 text-right">{t("dashboard.topDiagnoses.cases")}</TableHead>
                <TableHead className="w-48 pr-5 text-right">
                  {t("dashboard.topDiagnoses.proportion")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {rankedDiagnoses.map((diagnosis) => (
                <TableRow key={diagnosis.icd10_code}>
                  <TableCell className="pl-5 text-xs font-black tabular-nums text-muted-foreground">
                    {String(diagnosis.rank).padStart(2, "0")}
                  </TableCell>
                  <TableCell className="font-mono text-xs font-bold text-primary">
                    {diagnosis.icd10_code}
                  </TableCell>
                  <TableCell className="text-sm text-gray-700">{diagnosis.description}</TableCell>
                  <TableCell className="text-right text-sm font-black tabular-nums text-gray-900">
                    {diagnosis.count}
                  </TableCell>
                  <TableCell className="pr-5">
                    <div className="flex items-center justify-end gap-2.5">
                      <div className="h-2 w-28 overflow-hidden rounded-full bg-gray-200">
                        <div
                          className="h-full rounded-full"
                          style={{
                            width: `${diagnosis.proportion}%`,
                            backgroundColor: toneFillVariables.primary,
                          }}
                        />
                      </div>
                      <span className="w-9 text-right text-[11px] font-bold tabular-nums text-muted-foreground">
                        {diagnosis.proportion.toFixed(0)}%
                      </span>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      ) : (
        <EmptyState
          title={t("dashboard.topDiagnoses.noDataTitle")}
          description={t("dashboard.topDiagnoses.noDataDescription")}
        />
      )}
    </DashboardPanel>
  )
}
