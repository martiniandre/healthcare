import { usePortalObservationsQuery } from "./queries"
import { Card } from "../../shared/components/ui/Card"
import { Loader2 } from "lucide-react"
import { useTranslation } from "react-i18next"
import { formatDateTime } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { findVitalSignDisplay } from "../patients/components/vitalSignDisplay"
import { VitalSignValueDisplay } from "../patients/components/VitalSignValueDisplay"

export const PortalObservations = () => {
  const { t } = useTranslation("patients")
  const { data: observations, isLoading } = usePortalObservationsQuery()
  const locale = useLocale()

  if (isLoading) {
    return (
      <Card className="flex items-center justify-center min-h-[300px]">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
      </Card>
    )
  }

  if (!observations || observations.length === 0) {
    return (
      <Card className="py-16 text-center">
        <p className="text-sm text-gray-500">Nenhum sinal vital registrado.</p>
      </Card>
    )
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      {observations.map((observation) => {
        const displayMetadata = findVitalSignDisplay(observation.loinc_code)
        return (
          <div
            key={observation.fhir_resource_id}
            className="bg-white border border-border rounded-xl p-5"
          >
            <p className="text-xs text-gray-500 font-medium mb-1">
              {displayMetadata?.labelKey ? t(displayMetadata.labelKey) : observation.code_display}
            </p>
            <VitalSignValueDisplay
              notPerformed={observation.not_performed}
              valueQuantity={observation.value_quantity}
              valueUnit={observation.value_unit}
              valueClassName="text-2xl font-bold text-gray-900"
            />
            <p className="text-xs text-gray-400 mt-2">
              {formatDateTime(observation.observed_at, locale)}
            </p>
          </div>
        )
      })}
    </div>
  )
}
