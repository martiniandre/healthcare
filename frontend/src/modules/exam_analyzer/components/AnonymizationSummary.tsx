import { useTranslation } from "react-i18next"
import { ShieldAlert } from "lucide-react"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { StatusBadge } from "@/shared/components/ui/status-badge"
import type { AnonymizationReport } from "../types"

interface AnonymizationSummaryProperties {
  report: AnonymizationReport | null
}

const STEP_LABEL_KEYS: Record<string, string> = {
  metadata_removed: "card.anonymization.steps.metadata_removed",
  filename_neutralized: "card.anonymization.steps.filename_neutralized",
}

const LIMITATION_LABEL_KEYS: Record<string, string> = {
  burned_in_text_not_removed: "card.anonymization.limitations.burned_in_text_not_removed",
  pdf_metadata_not_removed: "card.anonymization.limitations.pdf_metadata_not_removed",
}

const UNKNOWN_STEP_KEY = "card.anonymization.unknownStep"
const UNKNOWN_LIMITATION_KEY = "card.anonymization.unknownLimitation"

export const AnonymizationSummary = ({ report }: AnonymizationSummaryProperties) => {
  const { t } = useTranslation("examAnalyzer")

  if (!report || !report.applied) {
    return (
      <div className="flex flex-col gap-2" data-testid="anonymization-summary">
        <StatusBadge status="attention" className="w-fit text-[11px] font-bold uppercase tracking-wider">
          {t("card.anonymization.notAppliedBadge")}
        </StatusBadge>
        <p className="text-[11px] leading-relaxed text-muted-foreground">
          {t("card.anonymization.notAppliedDescription")}
        </p>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-3" data-testid="anonymization-summary">
      <StatusBadge status="normal" className="w-fit text-[11px] font-bold uppercase tracking-wider">
        {t("card.anonymization.appliedBadge")}
      </StatusBadge>

      {report.steps.length > 0 && (
        <div>
          <span className="block text-[11px] font-bold text-foreground">{t("card.anonymization.stepsTitle")}</span>
          <ul className="mt-1 flex list-disc flex-col gap-1 pl-4 text-[11px] leading-normal text-foreground/80">
            {report.steps.map((stepCode, index) => (
              <li key={`${stepCode}-${index}`}>{t(STEP_LABEL_KEYS[stepCode] ?? UNKNOWN_STEP_KEY)}</li>
            ))}
          </ul>
        </div>
      )}

      {report.limitations.length > 0 && (
        <Alert variant="warning">
          <ShieldAlert aria-hidden="true" />
          <AlertTitle className="line-clamp-none text-[11px] font-bold">
            {t("card.anonymization.limitationsTitle")}
          </AlertTitle>
          <AlertDescription>
            <ul className="flex list-disc flex-col gap-1 pl-4 text-[11px] leading-normal text-foreground/80">
              {report.limitations.map((limitationCode, index) => (
                <li key={`${limitationCode}-${index}`}>
                  {t(LIMITATION_LABEL_KEYS[limitationCode] ?? UNKNOWN_LIMITATION_KEY)}
                </li>
              ))}
            </ul>
          </AlertDescription>
        </Alert>
      )}
    </div>
  )
}
