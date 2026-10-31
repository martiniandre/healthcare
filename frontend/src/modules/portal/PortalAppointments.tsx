import { CalendarClock } from "lucide-react"
import { useTranslation } from "react-i18next"
import { useMyAppointmentsQuery } from "../schedule/queries"
import { PortalSectionState } from "./components/PortalSectionState"
import { Card, CardContent } from "../../shared/components/ui/card"
import { StatusBadge, type ClinicalStatus } from "../../shared/components/ui/status-badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../shared/components/ui/table"
import { formatDate, formatTime } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"

const appointmentStatusMap: Record<string, ClinicalStatus> = {
  scheduled: "info",
  confirmed: "normal",
  cancelled: "critical",
  finished: "neutral",
}

export const PortalAppointments = () => {
  const { t } = useTranslation(["schedule", "portal"])
  const locale = useLocale()
  const { data: appointments, isLoading, isError, refetch } = useMyAppointmentsQuery()

  return (
    <PortalSectionState
      isLoading={isLoading}
      isError={isError}
      onRetry={() => refetch()}
      isEmpty={!appointments || appointments.length === 0}
      emptyTitle={t("schedule:portal.empty")}
      emptyIcon={CalendarClock}
    >
      <div className="flex flex-col gap-3 md:hidden">
        {appointments?.map((appointment) => (
          <Card key={appointment.id} className="py-4">
            <CardContent className="flex flex-col gap-2 px-4">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <p className="text-sm font-semibold tabular-nums text-foreground">
                  {formatDate(appointment.starts_at, locale)} - {formatTime(appointment.starts_at, locale)}
                </p>
                <StatusBadge status={appointmentStatusMap[appointment.status] ?? "neutral"}>
                  {t(`schedule:status.${appointment.status}`)}
                </StatusBadge>
              </div>
              <p className="break-words text-sm text-muted-foreground">{appointment.reason || "-"}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card className="hidden gap-0 overflow-hidden py-0 md:block">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="p-4">{t("schedule:portal.date")}</TableHead>
              <TableHead className="p-4">{t("schedule:portal.time")}</TableHead>
              <TableHead className="p-4">{t("schedule:portal.reason")}</TableHead>
              <TableHead className="p-4">{t("schedule:portal.status")}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {appointments?.map((appointment) => (
              <TableRow key={appointment.id}>
                <TableCell className="p-4 font-medium tabular-nums whitespace-nowrap">
                  {formatDate(appointment.starts_at, locale)}
                </TableCell>
                <TableCell className="p-4 tabular-nums whitespace-nowrap">
                  {formatTime(appointment.starts_at, locale)}
                </TableCell>
                <TableCell className="p-4 whitespace-normal break-words">{appointment.reason || "-"}</TableCell>
                <TableCell className="p-4">
                  <StatusBadge status={appointmentStatusMap[appointment.status] ?? "neutral"}>
                    {t(`schedule:status.${appointment.status}`)}
                  </StatusBadge>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Card>
    </PortalSectionState>
  )
}
