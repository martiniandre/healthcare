import { useTranslation } from "react-i18next"
import { AlertTriangle, Loader2, ShieldAlert } from "lucide-react"
import { Card, CardContent } from "@/shared/components/ui/card"
import type { ExamAnalysis } from "../types"

const CENTERED_STATE_CLASS = "min-h-[300px] flex-1 items-center justify-center gap-0 p-16 text-center"

export const AnalysisProcessingState = () => {
  const { t } = useTranslation("examAnalyzer")

  return (
    <Card className={CENTERED_STATE_CLASS} role="status" aria-live="polite">
      <Loader2 aria-hidden="true" className="mb-3 size-6 animate-spin text-muted-foreground motion-reduce:animate-none" />
      <h4 className="text-sm font-semibold text-foreground">{t("card.processing")}</h4>
      <span className="mt-1.5 block max-w-xs text-xs leading-normal text-muted-foreground">{t("card.processingDesc")}</span>
    </Card>
  )
}

export const AnalysisFailedState = () => {
  const { t } = useTranslation("examAnalyzer")

  return (
    <Card className={CENTERED_STATE_CLASS} role="alert">
      <ShieldAlert aria-hidden="true" className="mb-3 size-6 text-destructive" />
      <h4 className="text-sm font-semibold text-foreground">{t("card.failed")}</h4>
      <span className="mt-1.5 block max-w-xs text-xs leading-normal text-muted-foreground">{t("card.failedDesc")}</span>
    </Card>
  )
}

interface AnalysisInsufficientStateProperties {
  activeAnalysis: ExamAnalysis
}

export const AnalysisInsufficientState = ({ activeAnalysis }: AnalysisInsufficientStateProperties) => {
  const { t } = useTranslation("examAnalyzer")
  const insufficientMessage =
    (activeAnalysis.analysis_response as { message?: string } | null)?.message ?? t("card.insufficientDefault")

  return (
    <Card className="flex-1 border-l-4 border-l-destructive bg-destructive/5">
      <CardContent className="flex items-start gap-4">
        <AlertTriangle aria-hidden="true" className="mt-0.5 size-5 shrink-0 text-destructive" />
        <div>
          <h4 className="text-sm font-semibold leading-none text-foreground">{t("card.insufficient")}</h4>
          <p className="mt-2 text-xs leading-relaxed text-foreground/80">{insufficientMessage}</p>
          <div className="mt-4 rounded-lg border border-destructive/15 bg-card p-3.5">
            <span className="block text-2xs font-semibold text-foreground">{t("card.possibleCauses")}</span>
            <ul className="mt-1.5 flex list-disc flex-col gap-1 pl-4 text-2xs leading-normal text-muted-foreground">
              <li>{t("card.cause1")}</li>
              <li>{t("card.cause2")}</li>
              <li>{t("card.cause3")}</li>
              <li>{t("card.cause4")}</li>
            </ul>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
