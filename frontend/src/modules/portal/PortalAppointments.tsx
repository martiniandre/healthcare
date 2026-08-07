import { useTranslation } from "react-i18next"
import { useMyAppointmentsQuery } from "../schedule/queries"
import { Card } from "../../shared/components/ui/Card"
import { Loader2 } from "lucide-react"

const statusBadgeClassNames: Record<string, string> = {
  scheduled: "bg-blue-100 text-blue-700",
  confirmed: "bg-emerald-100 text-emerald-700",
  cancelled: "bg-gray-200 text-gray-500",
  finished: "bg-violet-100 text-violet-700",
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
        <p className="text-sm text-gray-500">{t("portal.empty")}</p>
      </Card>
    )
  }

  return (
    <div className="bg-white border border-border rounded-xl overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-gray-50 border-b border-border">
              <th className="text-left p-4 text-xs font-bold text-gray-500 uppercase tracking-wider">{t("portal.date")}</th>
              <th className="text-left p-4 text-xs font-bold text-gray-500 uppercase tracking-wider">{t("portal.time")}</th>
              <th className="text-left p-4 text-xs font-bold text-gray-500 uppercase tracking-wider">{t("portal.reason")}</th>
              <th className="text-left p-4 text-xs font-bold text-gray-500 uppercase tracking-wider">{t("portal.status")}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {appointments.map((appointment) => (
              <tr key={appointment.id} className="hover:bg-gray-50">
                <td className="p-4 text-gray-900 font-medium whitespace-nowrap">
                  {new Date(appointment.starts_at).toLocaleDateString("pt-BR")}
                </td>
                <td className="p-4 font-mono text-gray-700 whitespace-nowrap">
                  {new Date(appointment.starts_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                </td>
                <td className="p-4 text-gray-700">{appointment.reason || "-"}</td>
                <td className="p-4">
                  <span className={`text-xs font-bold px-2 py-1 rounded-full capitalize ${statusBadgeClassNames[appointment.status] ?? "bg-gray-100 text-gray-500"}`}>
                    {t(`status.${appointment.status}`)}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
