import { useTranslation } from "react-i18next"
import { usePortalMedicationsQuery } from "./queries"
import { Card } from "../../shared/components/ui/Card"
import { Loader2 } from "lucide-react"
import { formatDate } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { useClinicalStatusLabels } from "../../shared/hooks/useClinicalStatusLabels"
import type { MedicationRequestStatus } from "../../shared/types"

const statusBadgeClass = (status: MedicationRequestStatus) => {
  switch (status) {
    case "active":
      return "bg-green-100 text-green-700"
    case "on-hold":
      return "bg-amber-100 text-amber-700"
    case "completed":
    case "stopped":
      return "bg-gray-100 text-gray-700"
    case "cancelled":
      return "bg-red-100 text-red-700"
    default:
      return "bg-gray-100 text-gray-700"
  }
}

export const PortalMedications = () => {
  const { t } = useTranslation("portal")
  const { data: medications, isLoading } = usePortalMedicationsQuery()
  const locale = useLocale()
  const { medicationRequestStatus, notInformed } = useClinicalStatusLabels()

  if (isLoading) {
    return (
      <Card className="flex items-center justify-center min-h-[300px]">
        <Loader2 className="w-8 h-8 text-primary animate-spin" />
      </Card>
    )
  }

  if (!medications || medications.length === 0) {
    return (
      <Card className="py-16 text-center">
        <p className="text-sm text-gray-500">{t("empty.medications")}</p>
      </Card>
    )
  }

  return (
    <div className="flex flex-col gap-3">
      {medications.map((medication) => (
        <div
          key={medication.fhir_resource_id}
          className="bg-card border border-border rounded-xl p-5"
        >
          <div className="flex items-start justify-between mb-2">
            <div>
              <p className="text-sm font-bold text-gray-900">{medication.medication_name || t("fallback.medication")}</p>
              <p className="text-xs text-gray-500 mt-0.5">
                {t("table.prescribedAt")}{" "}
                {medication.issued_at
                  ? formatDate(medication.issued_at, locale)
                  : notInformed()}
              </p>
            </div>
            <span
              className={`text-xs font-bold px-2.5 py-1 rounded-full ${statusBadgeClass(medication.status)}`}
            >
              {medicationRequestStatus(medication.status)}
            </span>
          </div>
          {medication.dosage_instructions && (
            <p className="text-xs text-gray-600 bg-gray-50 rounded-lg p-3 mt-2">
              {medication.dosage_instructions}
            </p>
          )}
        </div>
      ))}
    </div>
  )
}
