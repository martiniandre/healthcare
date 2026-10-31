import { History } from "lucide-react"
import { useTranslation } from "react-i18next"
import { usePortalEncountersQuery } from "./queries"
import { PortalSectionState } from "./components/PortalSectionState"
import { Card, CardContent } from "../../shared/components/ui/card"
import { StatusBadge } from "../../shared/components/ui/status-badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../shared/components/ui/table"
import { formatDate } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { useClinicalStatusLabels } from "../../shared/hooks/useClinicalStatusLabels"

export const PortalEncounters = () => {
  const { t } = useTranslation("portal")
  const { data: encounters, isLoading, isError, refetch } = usePortalEncountersQuery()
  const locale = useLocale()
  const { encounterStatus } = useClinicalStatusLabels()

  return (
    <PortalSectionState
      isLoading={isLoading}
      isError={isError}
      onRetry={() => refetch()}
      isEmpty={!encounters || encounters.length === 0}
      emptyTitle={t("emptyEncounters")}
      emptyIcon={History}
    >
      <div className="flex flex-col gap-3 md:hidden">
        {encounters?.map((encounter) => (
          <Card key={encounter.fhir_resource_id} className="py-4">
            <CardContent className="flex flex-col gap-2 px-4">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <p className="text-sm font-semibold tabular-nums text-foreground">
                  {formatDate(encounter.started_at, locale)}
                </p>
                <StatusBadge status="info">{encounterStatus(encounter.status)}</StatusBadge>
              </div>
              <p className="break-words text-sm text-muted-foreground">
                {encounter.reason_display || t("fallback.encounter")}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card className="hidden gap-0 overflow-hidden py-0 md:block">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="p-4">{t("table.date")}</TableHead>
              <TableHead className="p-4">{t("table.reason")}</TableHead>
              <TableHead className="p-4">{t("table.status")}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {encounters?.map((encounter) => (
              <TableRow key={encounter.fhir_resource_id}>
                <TableCell className="p-4 font-medium tabular-nums whitespace-nowrap">
                  {formatDate(encounter.started_at, locale)}
                </TableCell>
                <TableCell className="p-4 whitespace-normal break-words">
                  {encounter.reason_display || t("fallback.encounter")}
                </TableCell>
                <TableCell className="p-4">
                  <StatusBadge status="info">{encounterStatus(encounter.status)}</StatusBadge>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Card>
    </PortalSectionState>
  )
}
