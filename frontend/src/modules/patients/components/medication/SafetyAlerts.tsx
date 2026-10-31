import { AlertTriangle, ShieldAlert } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Alert, AlertDescription, AlertTitle } from "../../../../shared/components/ui/alert"
import { Button } from "../../../../shared/components/ui/button"
import { Checkbox } from "../../../../shared/components/ui/checkbox"
import { Field, FieldError } from "../../../../shared/components/ui/field"
import { Label } from "../../../../shared/components/ui/label"
import { Skeleton } from "../../../../shared/components/ui/skeleton"
import { Spinner } from "../../../../shared/components/ui/spinner"
import { StatusBadge } from "../../../../shared/components/ui/status-badge"
import { buildAlertMessage, resolveAlertBadgeStatus, sortAlertsByKind, type AlertTranslateFunction } from "../../medicationSafety"
import type { MedicationSafetyStatus } from "../../queries"
import type { MedicationSafetyAlert } from "../../types"

interface SafetyAlertsProps {
  status: MedicationSafetyStatus
  alerts: MedicationSafetyAlert[]
  onRetry: () => void
  unverifiedAcknowledged: boolean
  onUnverifiedAcknowledgedChange: (acknowledged: boolean) => void
  unverifiedAcknowledgementError?: string
}

export function SafetyAlerts({
  status,
  alerts,
  onRetry,
  unverifiedAcknowledged,
  onUnverifiedAcknowledgedChange,
  unverifiedAcknowledgementError,
}: SafetyAlertsProps) {
  const { t } = useTranslation("patients")

  if (status === "idle") return null

  if (status === "loading") {
    return (
      <div className="flex flex-col gap-2" role="status" aria-live="polite">
        <div className="flex items-center gap-2 text-xs font-semibold text-muted-foreground">
          <Spinner />
          {t("medicationSafety.checking")}
        </div>
        <Skeleton className="h-12 w-full" />
      </div>
    )
  }

  if (status === "error") {
    return (
      <Alert variant="warning" data-testid="safety-check-failed">
        <AlertTriangle aria-hidden="true" />
        <AlertTitle>{t("medicationSafety.checkFailed.title")}</AlertTitle>
        <AlertDescription>
          <p>{t("medicationSafety.checkFailed.description")}</p>
          <Button type="button" variant="outline" size="sm" className="min-h-11 md:min-h-8" onClick={onRetry}>
            {t("medicationSafety.checkFailed.retry")}
          </Button>
          <Field data-invalid={Boolean(unverifiedAcknowledgementError)}>
            <div className="flex min-h-11 items-center gap-2 md:min-h-0">
              <Checkbox
                id="unverified-acknowledgement"
                checked={unverifiedAcknowledged}
                aria-invalid={Boolean(unverifiedAcknowledgementError)}
                onCheckedChange={(checkedState) => onUnverifiedAcknowledgedChange(checkedState === true)}
              />
              <Label htmlFor="unverified-acknowledgement" className="text-sm text-foreground">
                {t("medicationSafety.checkFailed.acknowledge")}
              </Label>
            </div>
            {unverifiedAcknowledgementError && <FieldError>{unverifiedAcknowledgementError}</FieldError>}
          </Field>
        </AlertDescription>
      </Alert>
    )
  }

  if (alerts.length === 0) return null

  return (
    <ul className="flex flex-col gap-2" aria-label={t("medicationSafety.title")}>
      {sortAlertsByKind(alerts).map((alert) => {
        const isAllergy = alert.kind === "allergy"
        const AlertIcon = isAllergy ? ShieldAlert : AlertTriangle
        return (
          <li key={`${alert.kind}-${alert.messageKey}-${alert.relatedResourceId}`}>
            <Alert variant={isAllergy ? "destructive" : "warning"} data-testid={`safety-alert-${alert.kind}`}>
              <AlertIcon aria-hidden="true" />
              <AlertTitle className="flex flex-wrap items-center gap-2">
                {t(`medicationSafety.kind.${alert.kind}`)}
                <StatusBadge status={resolveAlertBadgeStatus(alert)}>
                  {t(`medicationSafety.severity.${alert.severity}`)}
                </StatusBadge>
              </AlertTitle>
              <AlertDescription>
                <p>{buildAlertMessage(t as AlertTranslateFunction, alert)}</p>
              </AlertDescription>
            </Alert>
          </li>
        )
      })}
    </ul>
  )
}
