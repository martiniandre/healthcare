import { useTranslation } from "react-i18next"
import {
  AlertTriangle,
  ArrowUpRight,
  CheckCircle2,
  Clock,
  HelpCircle,
  Info,
  ShieldAlert,
  Sparkles,
  type LucideIcon,
} from "lucide-react"
import { Card } from "../../../shared/components/ui/Card"
import { ExamAnalysisStatus, type ExamAnalysis, type MedicalAnalysisResponse } from "../types"
import {
  resolveConfidenceBand,
  resolveFindingTriage,
  resolveQualityVerdict,
  resolveUrgencyTone,
  truncateIdentifier,
  type ConfidenceBand,
  type FindingTriageTier,
  type QualityTone,
  type UrgencyTone,
} from "../exam_analysis_domain"
import { cn } from "../../../shared/utils/cn"

interface AnalysisCardProperties {
  activeAnalysis: ExamAnalysis | null
}

const QUALITY_TONE_STYLES: Record<QualityTone, { text: string; surface: string; meter: string }> = {
  good: { text: "text-success", surface: "bg-success/10", meter: "bg-success" },
  fair: { text: "text-warning", surface: "bg-warning/10", meter: "bg-warning" },
  poor: { text: "text-danger", surface: "bg-danger/10", meter: "bg-danger" },
}

const TRIAGE_TIER_STYLES: Record<FindingTriageTier, { text: string; surface: string }> = {
  priority: { text: "text-danger", surface: "bg-danger/10" },
  verify: { text: "text-warning", surface: "bg-warning/10" },
  review: { text: "text-primary", surface: "bg-primary/10" },
  informative: { text: "text-muted-foreground", surface: "bg-gray-100" },
}

const CONFIDENCE_BAND_STYLES: Record<ConfidenceBand, { text: string; meter: string }> = {
  strong: { text: "text-success", meter: "bg-success" },
  moderate: { text: "text-warning", meter: "bg-warning" },
  weak: { text: "text-danger", meter: "bg-danger" },
}

const URGENCY_TONE_STYLES: Record<
  UrgencyTone,
  { text: string; surface: string; border: string; Icon: LucideIcon }
> = {
  urgent: {
    text: "text-danger",
    surface: "bg-danger/8",
    border: "border-danger/30",
    Icon: AlertTriangle,
  },
  followup: {
    text: "text-warning",
    surface: "bg-warning/8",
    border: "border-warning/30",
    Icon: Clock,
  },
  routine: {
    text: "text-success",
    surface: "bg-success/8",
    border: "border-success/30",
    Icon: CheckCircle2,
  },
}

const SECTION_TITLE_CLASS = "text-[11px] font-bold uppercase tracking-wider text-muted-foreground"
const BODY_TEXT_CLASS = "text-xs leading-relaxed text-gray-700"

const SectionHeading = ({ title, caption }: { title: string; caption?: string }) => (
  <div className="mb-2.5">
    <h4 className={SECTION_TITLE_CLASS}>{title}</h4>
    {caption ? <span className="mt-1 block text-[11px] leading-normal text-muted">{caption}</span> : null}
  </div>
)

const EmptySection = ({ message }: { message: string }) => (
  <p className="rounded-lg border border-dashed border-border bg-gray-50/60 px-3 py-2.5 text-[11px] leading-normal text-muted">
    {message}
  </p>
)

export const AnalysisCard = ({ activeAnalysis }: AnalysisCardProperties) => {
  const { t } = useTranslation("examAnalyzer")

  if (!activeAnalysis) {
    return (
      <div className="flex min-h-[300px] flex-1 flex-col items-center justify-center rounded-xl border border-dashed border-border bg-card p-16 text-center">
        <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl border border-border bg-gray-50">
          <HelpCircle className="h-6 w-6 text-gray-300" />
        </div>
        <h4 className="text-sm font-bold text-gray-900">{t("card.noExam")}</h4>
        <span className="mt-1.5 block max-w-xs text-xs leading-normal text-muted">
          {t("card.noExamDesc")}
        </span>
      </div>
    )
  }

  const isAnalysisInFlight =
    activeAnalysis.status === ExamAnalysisStatus.PENDING ||
    activeAnalysis.status === ExamAnalysisStatus.PROCESSING

  if (isAnalysisInFlight) {
    return (
      <Card
        glowingType="amethyst"
        className="flex min-h-[300px] flex-1 flex-col items-center justify-center bg-card p-16 text-center"
      >
        <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-secondary/10">
          <Clock className="h-6 w-6 animate-spin text-secondary" />
        </div>
        <h4 className="text-sm font-bold text-gray-900">{t("card.processing")}</h4>
        <span className="mt-1.5 block max-w-xs text-xs leading-normal text-muted">
          {t("card.processingDesc")}
        </span>
      </Card>
    )
  }

  if (activeAnalysis.status === ExamAnalysisStatus.FAILED) {
    return (
      <Card className="flex min-h-[300px] flex-1 flex-col items-center justify-center bg-card p-16 text-center">
        <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl border border-danger/20 bg-danger/10">
          <ShieldAlert className="h-6 w-6 text-danger" />
        </div>
        <h4 className="text-sm font-bold text-gray-900">{t("card.failed")}</h4>
        <span className="mt-1.5 block max-w-xs text-xs leading-normal text-muted">
          {t("card.failedDesc")}
        </span>
      </Card>
    )
  }

  if (activeAnalysis.status === ExamAnalysisStatus.INSUFFICIENT_DATA) {
    const insufficientMessage =
      (activeAnalysis.analysis_response as { message?: string })?.message ??
      t("card.insufficientDefault")

    return (
      <Card className="flex-1 border-l-4 border-l-danger bg-danger/5">
        <div className="flex items-start gap-4">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border border-danger/20 bg-danger/10">
            <AlertTriangle className="h-5 w-5 text-danger" />
          </div>
          <div>
            <h4 className="text-sm font-black leading-none text-gray-900">
              {t("card.insufficient")}
            </h4>
            <p className="mt-2 text-xs leading-relaxed text-gray-700">{insufficientMessage}</p>
            <div className="mt-4 rounded-lg border border-danger/15 bg-card p-3.5">
              <span className="block text-[11px] font-bold text-gray-900">
                {t("card.possibleCauses")}
              </span>
              <ul className="mt-1.5 flex list-disc flex-col gap-1 pl-4 text-[11px] leading-normal text-gray-600">
                <li>{t("card.cause1")}</li>
                <li>{t("card.cause2")}</li>
                <li>{t("card.cause3")}</li>
                <li>{t("card.cause4")}</li>
              </ul>
            </div>
          </div>
        </div>
      </Card>
    )
  }

  const analysisPayload = activeAnalysis.analysis_response as MedicalAnalysisResponse

  if (!analysisPayload) {
    return null
  }

  const qualityVerdict = resolveQualityVerdict(analysisPayload.qualityAssessment.score)
  const qualityStyles = QUALITY_TONE_STYLES[qualityVerdict.tone]
  const urgencyReading = resolveUrgencyTone(analysisPayload.recommendation.urgency)
  const urgencyStyles = URGENCY_TONE_STYLES[urgencyReading.tone]
  const UrgencyIcon = urgencyStyles.Icon
  const qualityWarnings = analysisPayload.qualityAssessment.warnings
  const nextSteps = analysisPayload.recommendation.nextSteps

  return (
    <div className="flex flex-1 flex-col gap-4 animate-fade-in">
      <Card glowingType="cyan" className="bg-card">
        <div className="flex flex-col gap-4 border-b border-border/70 pb-4 sm:flex-row sm:items-start sm:justify-between">
          <div className="min-w-0">
            <div className="flex items-center gap-2">
              <Sparkles className="h-4 w-4 shrink-0 text-primary" />
              <h3 className="truncate text-base font-black text-gray-900">
                {analysisPayload.examType}
              </h3>
            </div>
            <span
              className="mt-1.5 block font-mono text-[11px] text-muted"
              title={activeAnalysis.id}
            >
              {t("card.analysisId")} {truncateIdentifier(activeAnalysis.id)}
            </span>
          </div>

          <div
            className={cn(
              "flex shrink-0 items-center gap-3 rounded-xl border px-3.5 py-2",
              qualityStyles.surface
            )}
          >
            <div className="min-w-0">
              <span className="block text-[11px] font-bold uppercase tracking-wider text-muted-foreground">
                {t("card.qualityAssessment")}
              </span>
              <span className={cn("mt-0.5 block text-[11px] font-bold", qualityStyles.text)}>
                {t(qualityVerdict.translationKey)}
              </span>
            </div>
            <span
              className={cn("text-lg font-black tabular-nums", qualityStyles.text)}
              aria-label={`${qualityVerdict.percentage}%`}
            >
              {qualityVerdict.percentage}%
            </span>
            <div
              className="flex h-9 w-1.5 shrink-0 flex-col justify-end overflow-hidden rounded-full bg-black/10"
              role="img"
              aria-label={`${t("card.qualityAssessment")}: ${qualityVerdict.percentage}%`}
            >
              <div
                className={cn("w-full rounded-full", qualityStyles.meter)}
                style={{ height: `${qualityVerdict.percentage}%` }}
              />
            </div>
          </div>
        </div>

        <div
          className={cn(
            "mt-4 flex items-start gap-3 rounded-xl border px-4 py-3",
            urgencyStyles.surface,
            urgencyStyles.border
          )}
        >
          <UrgencyIcon className={cn("mt-0.5 h-4 w-4 shrink-0", urgencyStyles.text)} aria-hidden="true" />
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-baseline gap-x-2 gap-y-1">
              <span className={cn("text-xs font-black uppercase tracking-wider", urgencyStyles.text)}>
                {t(urgencyReading.translationKey)}
              </span>
              <span className="text-[11px] text-muted-foreground">{t("card.urgencyTitle")}</span>
            </div>
            <span className="mt-1 block text-[11px] leading-normal text-muted-foreground">
              {t("card.urgencyCaption")}
            </span>
            {nextSteps.length > 0 ? (
              <ol className="mt-2.5 flex flex-col gap-1.5">
                {nextSteps.map((stepItem, index) => (
                  <li key={index} className="flex items-start gap-2 text-xs leading-relaxed text-gray-800">
                    <span className="mt-0.5 flex h-4 w-4 shrink-0 select-none items-center justify-center rounded border border-border bg-card text-[11px] font-bold text-muted-foreground">
                      {index + 1}
                    </span>
                    {stepItem}
                  </li>
                ))}
              </ol>
            ) : (
              <p className="mt-2.5 text-[11px] text-muted-foreground">{t("card.noNextSteps")}</p>
            )}
          </div>
        </div>

        {qualityWarnings.length > 0 && (
          <div className="mt-4 flex items-start gap-2.5 rounded-lg border border-border bg-gray-50 px-3.5 py-3">
            <Info className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground" aria-hidden="true" />
            <div className="min-w-0">
              <span className="block text-[11px] font-bold text-gray-900">
                {t("card.qualityObservations")}
              </span>
              <ul className="mt-1 flex list-disc flex-col gap-1 pl-4 text-[11px] leading-normal text-gray-700">
                {qualityWarnings.map((warningItem, index) => (
                  <li key={index}>{warningItem}</li>
                ))}
              </ul>
            </div>
          </div>
        )}
      </Card>

      <Card className="bg-card">
        <SectionHeading title={t("card.findings")} caption={t("card.findingsCaption")} />

        {analysisPayload.detectedFindings.length === 0 ? (
          <div className="flex items-center gap-2.5 rounded-lg border border-success/25 bg-success/8 px-3.5 py-3">
            <CheckCircle2 className="h-4 w-4 shrink-0 text-success" aria-hidden="true" />
            <p className="text-xs text-gray-700">{t("card.noFindings")}</p>
          </div>
        ) : (
          <ul className="flex flex-col gap-2.5">
            {analysisPayload.detectedFindings.map((findingItem, index) => {
              const triageReading = resolveFindingTriage(
                findingItem.severity,
                findingItem.confidence
              )
              const triageStyles = TRIAGE_TIER_STYLES[triageReading.tier]
              const confidenceReading = resolveConfidenceBand(findingItem.confidence)
              const confidenceStyles = CONFIDENCE_BAND_STYLES[confidenceReading.band]
              const significanceTranslationKey = `card.significance${
                findingItem.severity.charAt(0).toUpperCase()}${findingItem.severity.slice(1)
              }`

              return (
                <li
                  key={index}
                  className="rounded-lg border border-border bg-gray-50/60 p-3.5"
                >
                  <div className="flex flex-wrap items-start justify-between gap-2">
                    <span className="min-w-0 flex-1 text-xs font-semibold leading-relaxed text-gray-900">
                      {findingItem.finding}
                    </span>
                    <div className="flex shrink-0 items-center gap-1.5">
                      <span
                        className={cn(
                          "rounded-md border border-border bg-card px-2 py-0.5 text-[11px] font-bold uppercase tracking-wider text-muted-foreground"
                        )}
                      >
                        {t("card.significance")}: {t(significanceTranslationKey)}
                      </span>
                      <span
                        className={cn(
                          "rounded-md px-2 py-0.5 text-[11px] font-bold uppercase tracking-wider",
                          triageStyles.surface,
                          triageStyles.text
                        )}
                      >
                        {t(triageReading.translationKey)}
                      </span>
                    </div>
                  </div>

                  <div className="mt-2.5 flex items-center gap-2.5">
                    <span className="shrink-0 text-[11px] font-semibold text-muted-foreground">
                      {t("card.confidence")}
                    </span>
                    <div
                      className="h-1.5 flex-1 overflow-hidden rounded-full bg-black/10"
                      role="img"
                      aria-label={`${t("card.confidence")}: ${confidenceReading.percentage}%`}
                    >
                      <div
                        className={cn("h-full rounded-full", confidenceStyles.meter)}
                        style={{ width: `${confidenceReading.percentage}%` }}
                      />
                    </div>
                    <span
                      className={cn("shrink-0 text-[11px] font-bold tabular-nums", confidenceStyles.text)}
                    >
                      {confidenceReading.percentage}%
                    </span>
                    <span className="sr-only">{t(confidenceReading.translationKey)}</span>
                  </div>
                </li>
              )
            })}
          </ul>
        )}
      </Card>

      <Card className="bg-card">
        <SectionHeading title={t("card.interpretations")} />

        {analysisPayload.possibleInterpretations.length === 0 ? (
          <EmptySection message={t("card.noInterpretations")} />
        ) : (
          <ul className="flex flex-col gap-2">
            {analysisPayload.possibleInterpretations.map((interpretationItem, index) => (
              <li key={index} className="flex items-start gap-2">
                <ArrowUpRight className="mt-0.5 h-3.5 w-3.5 shrink-0 text-muted-foreground" aria-hidden="true" />
                <span className={BODY_TEXT_CLASS}>{interpretationItem}</span>
              </li>
            ))}
          </ul>
        )}

        <div className="my-4 h-px bg-border/60" />

        <SectionHeading title={t("card.limitations")} />

        {analysisPayload.limitations.length === 0 ? (
          <EmptySection message={t("card.noLimitations")} />
        ) : (
          <ul className="flex flex-col gap-2">
            {analysisPayload.limitations.map((limitationItem, index) => (
              <li key={index} className="flex items-start gap-2">
                <Info className="mt-0.5 h-3.5 w-3.5 shrink-0 text-muted-foreground" aria-hidden="true" />
                <span className="text-[11px] leading-relaxed text-gray-600">{limitationItem}</span>
              </li>
            ))}
          </ul>
        )}

        <div className="mt-4 flex items-start gap-2.5 rounded-xl border border-border bg-gray-50 p-3.5">
          <ShieldAlert className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground" aria-hidden="true" />
          <div>
            <span className="block text-[11px] font-bold uppercase tracking-wider text-gray-900">
              {t("card.disclaimerTitle")}
            </span>
            <p className="mt-1 text-[11px] font-medium leading-relaxed text-gray-600">
              {analysisPayload.disclaimer}
            </p>
          </div>
        </div>
      </Card>
    </div>
  )
}
