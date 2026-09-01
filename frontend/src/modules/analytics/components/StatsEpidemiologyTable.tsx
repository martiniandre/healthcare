import { useTranslation } from "react-i18next"
import { Card } from "../../../shared/components/ui/Card"
import { Button } from "../../../shared/components/ui/Button"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "../../../shared/components/ui/Table"
import { FileSpreadsheet } from "lucide-react"

interface Pathology {
  code: string
  description: string
  category: string
  activeCases: number
  trend: string
}

interface StatsEpidemiologyTableProps {
  pathologies: Pathology[]
}

export const StatsEpidemiologyTable = ({ pathologies }: StatsEpidemiologyTableProps) => {
  const { t: translate } = useTranslation("analytics")

  const getTrendStyle = (pathologyCode: string): string => {
    if (pathologyCode === "E11.9") {
      return "text-red-500 font-bold"
    }
    if (pathologyCode === "J45.9") {
      return "text-emerald-600 font-bold"
    }
    return "text-gray-400 font-bold"
  }

  return (
    <Card className="p-5 flex flex-col gap-4 text-left border border-border">
      <div className="flex items-center justify-between border-b border-border pb-3 flex-wrap gap-2">
        <div>
          <h3 className="font-extrabold text-gray-900 text-md">{translate("epidemiology.title")}</h3>
          <span className="text-xs text-muted block mt-1">{translate("epidemiology.subtitle")}</span>
        </div>
        <Button variantType="outline" className="px-3 gap-1.5 text-xs">
          <FileSpreadsheet className="w-4 h-4" />
          {translate("epidemiology.exportButton")}
        </Button>
      </div>

      <div className="overflow-x-auto w-full">
          {pathologies?.length ? (
            <Table className="min-w-[500px] md:min-w-0">
              <TableHeader>
                <TableRow>
                  <TableHead>{translate("epidemiology.table.code")}</TableHead>
                  <TableHead>{translate("epidemiology.table.description")}</TableHead>
                  <TableHead>{translate("epidemiology.table.category")}</TableHead>
                  <TableHead>{translate("epidemiology.table.activeCases")}</TableHead>
                  <TableHead className="text-right">{translate("epidemiology.table.trend")}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody className="text-gray-700 font-medium">
                {pathologies.map((pathologyItem) => {
                  const isStableTrend = pathologyItem.trend.toLowerCase() === "stable"

                  return (
                    <TableRow key={pathologyItem.code}>
                      <TableCell className="font-mono font-bold text-primary">{pathologyItem.code}</TableCell>
                      <TableCell>{pathologyItem.description}</TableCell>
                      <TableCell>{pathologyItem.category}</TableCell>
                      <TableCell className="font-bold text-gray-900">{pathologyItem.activeCases}</TableCell>
                      <TableCell className="text-right">
                        <span className={getTrendStyle(pathologyItem.code)}>
                          {isStableTrend ? translate("epidemiology.table.stable") : pathologyItem.trend}
                        </span>
                      </TableCell>
                    </TableRow>
                  )
                })}
              </TableBody>
            </Table>
          ) : (
            <EmptyState 
              title={translate("empty.epidemiology")} 
              description={translate("empty.epidemiologyDesc")} 
            />
          )}
      </div>
    </Card>
  )
}
