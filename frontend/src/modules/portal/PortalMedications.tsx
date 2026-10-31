import { Pill } from "lucide-react"
import { useTranslation } from "react-i18next"
import { usePortalMedicationsQuery } from "./queries"
import { PortalSectionState } from "./components/PortalSectionState"
import { Card, CardContent } from "../../shared/components/ui/card"
import { StatusBadge, type ClinicalStatus } from "../../shared/components/ui/status-badge"
import { formatDate } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { useClinicalStatusLabels } from "../../shared/hooks/useClinicalStatusLabels"
import type { MedicationRequestStatus } from "../../shared/types"

const medicationStatusToClinicalStatus = (status: MedicationRequestStatus): ClinicalStatus => {
  switch (status) {
    case "active":
      return "normal"
    case "on-hold":
      return "attention"
    case "cancelled":
      return "critical"
    default:
      return "neutral"
  }
}

export const PortalMedications = () => {
  const { t } = useTranslation("portal")
  const { data: medications, isLoading, isError, refetch } = usePortalMedicationsQuery()
  const locale = useLocale()
  const { medicationRequestStatus, notInformed } = useClinicalStatusLabels()

  return (
    <PortalSectionState
      isLoading={isLoading}
      isError={isError}
      onRetry={() => refetch()}
      isEmpty={!medications || medications.length === 0}
      emptyTitle={t("empty.medications")}
      emptyIcon={Pill}
    >
      <div className="flex flex-col gap-3">
        {medications?.map((medication) => (
          <Card key={medication.fhir_resource_id} className="py-4 sm:py-5">
            <CardContent className="px-4 sm:px-5">
              <div className="mb-2 flex flex-wrap items-start justify-between gap-2">
                <div className="min-w-0 flex-1">
                  <p className="break-words text-sm font-bold text-foreground">
                    {medication.medication_name || t("fallback.medication")}
                  </p>
                  <p className="mt-0.5 text-xs text-muted-foreground">
                    {t("table.prescribedAt")}{" "}
                    <span className="tabular-nums">
                      {medication.issued_at ? formatDate(medication.issued_at, locale) : notInformed()}
                    </span>
                  </p>
                </div>
                <StatusBadge status={medicationStatusToClinicalStatus(medication.status)}>
                  {medicationRequestStatus(medication.status)}
                </StatusBadge>
              </div>
              {medication.dosage_instructions && (
                <p className="mt-2 break-words rounded-lg bg-muted p-3 text-xs text-muted-foreground">
                  {medication.dosage_instructions}
                </p>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </PortalSectionState>
  )
}
