import { Activity } from "lucide-react"
import { useTranslation } from "react-i18next"
import { usePortalConditionsQuery } from "./queries"
import { PortalSectionState } from "./components/PortalSectionState"
import { Card, CardContent } from "../../shared/components/ui/card"
import { StatusBadge, type ClinicalStatus } from "../../shared/components/ui/status-badge"
import { formatDate } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { useClinicalStatusLabels } from "../../shared/hooks/useClinicalStatusLabels"
import type { ConditionClinicalStatus } from "../../shared/types"

const conditionStatusToClinicalStatus = (status: ConditionClinicalStatus): ClinicalStatus => {
  switch (status) {
    case "active":
      return "critical"
    case "recurrence":
    case "relapse":
    case "inactive":
    case "remission":
      return "attention"
    case "resolved":
      return "normal"
    default:
      return "neutral"
  }
}

export const PortalConditions = () => {
  const { t } = useTranslation("portal")
  const { data: conditions, isLoading, isError, refetch } = usePortalConditionsQuery()
  const locale = useLocale()
  const { conditionStatus, notInformed } = useClinicalStatusLabels()

  return (
    <PortalSectionState
      isLoading={isLoading}
      isError={isError}
      onRetry={() => refetch()}
      isEmpty={!conditions || conditions.length === 0}
      emptyTitle={t("emptyConditions")}
      emptyIcon={Activity}
    >
      <div className="flex flex-col gap-3">
        {conditions?.map((condition) => (
          <Card key={condition.fhir_resource_id} className="py-4 sm:py-5">
            <CardContent className="flex flex-wrap items-start justify-between gap-2 px-4 sm:px-5">
              <div className="min-w-0 flex-1">
                <p className="break-words text-sm font-bold text-foreground">{condition.code_display}</p>
                <p className="mt-0.5 text-xs text-muted-foreground">
                  {condition.icd10_code && <span className="font-mono">{condition.icd10_code} | </span>}
                  {t("table.onset")}{" "}
                  <span className="tabular-nums">
                    {condition.onset_at ? formatDate(condition.onset_at, locale) : notInformed()}
                  </span>
                </p>
              </div>
              <StatusBadge status={conditionStatusToClinicalStatus(condition.clinical_status)}>
                {conditionStatus(condition.clinical_status)}
              </StatusBadge>
            </CardContent>
          </Card>
        ))}
      </div>
    </PortalSectionState>
  )
}
