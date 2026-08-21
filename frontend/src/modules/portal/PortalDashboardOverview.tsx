import type { PortalDashboard } from "./types"
import { Card } from "../../shared/components/ui/Card"
import { History, Activity, Pill, FileText } from "lucide-react"
import { formatDate } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { useTranslation } from "react-i18next"
import { findVitalSignDisplay } from "../patients/components/vitalSignDisplay"
import { VitalSignValueDisplay } from "../patients/components/VitalSignValueDisplay"

interface PortalDashboardOverviewProps {
  dashboard: PortalDashboard
}

export const PortalDashboardOverview = ({ dashboard }: PortalDashboardOverviewProps) => {
  const { t } = useTranslation("patients")
  const locale = useLocale()
  const summaryCards = [
    {
      icon: <History className="w-5 h-5 text-blue-600" />,
      label: "Total de Consultas",
      value: dashboard.upcoming_encounters.length + dashboard.recent_reports.length,
      bgColor: "bg-blue-50",
    },
    {
      icon: <Activity className="w-5 h-5 text-amber-600" />,
      label: "Condições Ativas",
      value: dashboard.active_conditions.length,
      bgColor: "bg-amber-50",
    },
    {
      icon: <Pill className="w-5 h-5 text-green-600" />,
      label: "Medicamentos Ativos",
      value: dashboard.active_medications.length,
      bgColor: "bg-green-50",
    },
    {
      icon: <FileText className="w-5 h-5 text-purple-600" />,
      label: "Exames Recentes",
      value: dashboard.recent_reports.length,
      bgColor: "bg-purple-50",
    },
  ]

  return (
    <div className="flex flex-col gap-6">
      <div className="bg-white border border-border rounded-xl p-6">
        <h2 className="text-lg font-bold text-gray-900 mb-1">
          Olá, {dashboard.patient_info.full_name.split(" ")[0]}!
        </h2>
        <p className="text-sm text-gray-500">
          Bem-vindo ao seu portal de saúde. Aqui você encontra todas as suas informações clínicas.
        </p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {summaryCards.map((card) => (
          <Card key={card.label} className={`${card.bgColor} border-0 p-5`}>
            <div className="flex items-center gap-3">
              {card.icon}
              <div>
                <p className="text-2xl font-bold text-gray-900">{card.value}</p>
                <p className="text-xs text-gray-600 font-medium">{card.label}</p>
              </div>
            </div>
          </Card>
        ))}
      </div>

      {dashboard.upcoming_encounters.length > 0 && (
        <div className="bg-white border border-border rounded-xl p-6">
          <h3 className="text-sm font-bold text-gray-900 mb-4">Próximas Consultas</h3>
          <div className="space-y-3">
            {dashboard.upcoming_encounters.slice(0, 5).map((encounter) => (
              <div
                key={encounter.fhir_resource_id}
                className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
              >
                <div>
                  <p className="text-sm font-semibold text-gray-900">
                    {encounter.reason_display || "Consulta"}
                  </p>
                  <p className="text-xs text-gray-500">
                    {formatDate(encounter.started_at, locale)}
                  </p>
                </div>
                <span className="text-xs font-bold px-2 py-1 rounded-full bg-blue-100 text-blue-700 capitalize">
                  {encounter.status}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}

      {dashboard.recent_observations.length > 0 && (
        <div className="bg-white border border-border rounded-xl p-6">
          <h3 className="text-sm font-bold text-gray-900 mb-4">Últimos Sinais Vitais</h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {dashboard.recent_observations.slice(0, 8).map((observation) => {
              const displayMetadata = findVitalSignDisplay(observation.loinc_code)
              return (
                <div
                  key={observation.fhir_resource_id}
                  className="p-3 bg-gray-50 rounded-lg"
                >
                  <p className="text-xs text-gray-500 font-medium">
                    {displayMetadata?.labelKey ? t(displayMetadata.labelKey) : observation.code_display}
                  </p>
                  <VitalSignValueDisplay
                    notPerformed={observation.not_performed}
                    valueQuantity={observation.value_quantity}
                    valueUnit={observation.value_unit}
                    valueClassName="text-lg font-bold text-gray-900"
                  />
                </div>
              )
            })}
          </div>
        </div>
      )}
    </div>
  )
}
