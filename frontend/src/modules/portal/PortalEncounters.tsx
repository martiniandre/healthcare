import { useTranslation } from "react-i18next"
import { usePortalEncountersQuery } from "./queries"
import { Card } from "../../shared/components/ui/Card"
import { Badge } from "../../shared/components/ui/Badge"
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from "../../shared/components/ui/Table"
import { Loader2 } from "lucide-react"
import { formatDate } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"

export const PortalEncounters = () => {
  const { t } = useTranslation("portal")
  const { data: encounters, isLoading } = usePortalEncountersQuery()
  const locale = useLocale()

  if (isLoading) {
    return (
      <Card className="flex items-center justify-center min-h-[300px]">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
      </Card>
    )
  }

  if (!encounters || encounters.length === 0) {
    return (
      <Card className="py-16 text-center">
        <p className="text-sm text-muted-foreground">{t("emptyEncounters")}</p>
      </Card>
    )
  }

  return (
    <div className="bg-surface border border-border rounded-xl overflow-hidden">
      <div className="overflow-x-auto">
        <Table className="w-full text-sm">
          <TableHeader>
            <TableRow className="bg-muted-soft border-b border-border">
              <TableHead className="text-left p-4 text-xs font-bold uppercase tracking-wider">Data</TableHead>
              <TableHead className="text-left p-4 text-xs font-bold uppercase tracking-wider">Motivo</TableHead>
              <TableHead className="text-left p-4 text-xs font-bold uppercase tracking-wider">Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody className="divide-y divide-border">
            {encounters.map((encounter) => (
              <TableRow key={encounter.fhir_resource_id} className="hover:bg-muted-soft">
                <TableCell className="p-4 text-foreground font-medium whitespace-nowrap">
                  {formatDate(encounter.started_at, locale)}
                </TableCell>
                <TableCell className="p-4 text-foreground">{encounter.reason_display || "-"}</TableCell>
                <TableCell className="p-4">
                  <Badge variant="info" className="capitalize">
                    {encounter.status}
                  </Badge>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
