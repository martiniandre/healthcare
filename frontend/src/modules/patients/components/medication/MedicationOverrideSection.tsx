import type { ComponentProps } from "react"
import { Lock, ShieldAlert } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Alert, AlertDescription, AlertTitle } from "../../../../shared/components/ui/alert"
import { Checkbox } from "../../../../shared/components/ui/checkbox"
import { Field, FieldError, FieldLabel } from "../../../../shared/components/ui/field"
import { Label } from "../../../../shared/components/ui/label"
import { Textarea } from "../../../../shared/components/ui/textarea"

interface MedicationOverrideSectionProps {
  canOverride: boolean
  justificationProps: ComponentProps<typeof Textarea>
  justificationError?: string
  riskAcknowledged: boolean
  onRiskAcknowledgedChange: (acknowledged: boolean) => void
  riskAcknowledgedError?: string
  blockedByServer: boolean
}

export function MedicationOverrideSection({
  canOverride,
  justificationProps,
  justificationError,
  riskAcknowledged,
  onRiskAcknowledgedChange,
  riskAcknowledgedError,
  blockedByServer,
}: MedicationOverrideSectionProps) {
  const { t } = useTranslation("patients")

  if (!canOverride) {
    return (
      <Alert variant="destructive" data-testid="override-not-permitted">
        <Lock aria-hidden="true" />
        <AlertTitle>{t("medicationSafety.override.blockedTitle")}</AlertTitle>
        <AlertDescription>
          <p>{t("medicationSafety.override.notPermitted")}</p>
        </AlertDescription>
      </Alert>
    )
  }

  return (
    <section className="flex flex-col gap-3 rounded-lg border border-destructive/30 p-3" data-testid="override-section">
      <div className="flex items-start gap-2 text-sm font-semibold text-destructive">
        <ShieldAlert className="mt-0.5 size-4 shrink-0" aria-hidden="true" />
        <div className="flex flex-col gap-1">
          <span>{t("medicationSafety.override.title")}</span>
          <span className="text-xs font-normal text-muted-foreground">
            {blockedByServer ? t("medicationSafety.override.serverBlocked") : t("medicationSafety.override.description")}
          </span>
        </div>
      </div>
      <Field className="gap-1" data-invalid={Boolean(justificationError)}>
        <FieldLabel htmlFor="override-justification" className="text-xs font-semibold text-muted-foreground">
          {t("medicationSafety.override.justification")}
        </FieldLabel>
        <Textarea
          id="override-justification"
          className="h-24 resize-none"
          aria-invalid={Boolean(justificationError)}
          placeholder={t("medicationSafety.override.justificationPlaceholder")}
          {...justificationProps}
        />
        {justificationError && <FieldError>{justificationError}</FieldError>}
      </Field>
      <Field className="gap-1" data-invalid={Boolean(riskAcknowledgedError)}>
        <div className="flex min-h-11 items-center gap-2 md:min-h-0">
          <Checkbox
            id="override-risk-acknowledgement"
            checked={riskAcknowledged}
            aria-invalid={Boolean(riskAcknowledgedError)}
            onCheckedChange={(checkedState) => onRiskAcknowledgedChange(checkedState === true)}
          />
          <Label htmlFor="override-risk-acknowledgement" className="text-sm">
            {t("medicationSafety.override.acknowledgeRisk")}
          </Label>
        </div>
        {riskAcknowledgedError && <FieldError>{riskAcknowledgedError}</FieldError>}
      </Field>
    </section>
  )
}
