import { useTranslation } from "react-i18next"
import { useMyAppointmentsQuery } from "../schedule/queries"
import { Card } from "../../shared/components/ui/Card"
import { Badge } from "../../shared/components/ui/Badge"
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from "../../shared/components/ui/Table"
import { Loader2 } from "lucide-react"
import { formatDate, formatTime } from "../../shared/utils/dates"

const statusBadgeVariant: Record<string, "secondary" | "success" | "muted" | "info" | "warning"> = {
  scheduled: "info",
  confirmed: "success",
  cancelled: "muted",
  finished: "muted",
}

export const PortalAppointments = () => {
  const { t } = useTranslation("schedule")
  const { data: appointments, isLoading } = useMyAppointmentsQuery()

  if (isLoading) {
    return (
      <Card className="flex items-center justify-center min-h-[300px]">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
      </Card>
    )
  }

  if (!appointments || appointments.length === 0) {
    return (
      <Card className="py-16 text-center">
        <p className="text-sm text-muted-foreground">{t("portal.empty")}</p>
      </Card>
    )
  }

  return (
    <div className="bg-surface border border-border rounded-xl overflow-hidden">
      <div className="overflow-x-auto">
        <Table className="w-full text-sm">
          <TableHeader>
            <TableRow className="bg-muted-soft border-b border-border">
              <TableHead className="text-left p-4 text-xs font-bold uppercase tracking-wider">{t("portal.date")}</TableHead>
              <TableHead className="text-left p-4 text-xs font-bold uppercase tracking-wider">{t("portal.time")}</TableHead>
              <TableHead className="text-left p-4 text-xs font-bold uppercase tracking-wider">{t("portal.reason")}</TableHead>
              <TableHead className="text-left p-4 text-xs font-bold uppercase tracking-wider">{t("portal.status")}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody className="divide-y divide-border">
            {appointments.map((appointment) => (
              <TableRow key={appointment.id} className="hover:bg-muted-soft">
                <TableCell className="p-4 text-foreground font-medium whitespace-nowrap">
                  {formatDate(appointment.starts_at)}
                </TableCell>
                <TableCell className="p-4 font-mono text-foreground whitespace-nowrap">
                  {formatTime(appointment.starts_at)}
                </TableCell>
                <TableCell className="p-4 text-foreground">{appointment.reason || "-"}</TableCell>
                <TableCell className="p-4">
                  <Badge variant={statusBadgeVariant[appointment.status] ?? "secondary"} className="capitalize">
                    {t(`status.${appointment.status}`)}
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
