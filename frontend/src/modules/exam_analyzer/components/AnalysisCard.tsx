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
import { Card, CardContent, CardHeader } from "@/shared/components/ui/card"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { Badge } from "@/shared/components/ui/badge"
import { StatusBadge, type ClinicalStatus } from "@/shared/components/ui/status-badge"
import { Separator } from "@/shared/components/ui/separator"
import { EmptyState } from "@/shared/components/ui/EmptyState"
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
import { cn } from "@/shared/utils/cn"
import { AnonymizationSummary } from "./AnonymizationSummary"

interface AnalysisCardProperties {
  activeAnalysis: ExamAnalysis | null
}

const QUALITY_TONE_STYLES: Record<QualityTone, { text: string; surface: string; meter: string }> = {
  good: { text: "text-success", surface: "bg-success/10", meter: "bg-success" },
  fair: { text: "text-warning", surface: "bg-warning/10", meter: "bg-warning" },
  poor: { text: "text-destructive", surface: "bg-destructive/10", meter: "bg-destructive" },
}

const TRIAGE_TIER_STATUS: Record<FindingTriageTier, ClinicalStatus> = {
  priority: "critical",
  verify: "attention",
  review: "info",
  informative: "neutral",
}

const CONFIDENCE_BAND_STYLES: Record<ConfidenceBand, { text: string; meter: string }> = {
  strong: { text: "text-success", meter: "bg-success" },
  moderate: { text: "text-warning", meter: "bg-warning" },
  weak: { text: "text-destructive", meter: "bg-destructive" },
}

const URGENCY_TONE_STYLES: Record<
  UrgencyTone,
  { variant: "destructive" | "warning" | "success"; Icon: LucideIcon }
> = {
  urgent: { variant: "destructive", Icon: AlertTriangle },
  followup: { variant: "warning", Icon: Clock },
  routine: { variant: "success", Icon: CheckCircle2 },
}

const SECTION_TITLE_CLASS = "text-[11px] font-bold uppercase tracking-wider text-muted-foreground"
const BODY_TEXT_CLASS = "text-xs leading-relaxed text-foreground/80"

const SectionHeading = ({ title, caption }: { title: string; caption?: string }) => (
  <div className="mb-2.5">
    <h4 className={SECTION_TITLE_CLASS}>{title}</h4>
    {caption ? <span className="mt-1 block text-[11px] leading-normal text-muted-foreground">{caption}</span> : null}
  </div>
)

const EmptySection = ({ message }: { message: string }) => (
  <p className="rounded-lg border border-dashed border-border bg-muted/60 px-3 py-2.5 text-[11px] leading-normal text-muted-foreground">
    {message}
  </p>
)

export const AnalysisCard = ({ activeAnalysis }: AnalysisCardProperties) => {
  const { t } = useTranslation("examAnalyzer")

  if (!activeAnalysis) {
    return (
      <EmptyState
        icon={HelpCircle}
        title={t("card.noExam")}
        description={t("card.noExamDesc")}
        className="min-h-[300px] flex-1 p-16"
      />
    )
  }

  const isAnalysisInFlight =
    activeAnalysis.status === ExamAnalysisStatus.PENDING ||
    activeAnalysis.status === ExamAnalysisStatus.PROCESSING

  if (isAnalysisInFlight) {
    return (
      <Card className="min-h-[300px] flex-1 items-center justify-center gap-0 p-16 text-center">
        <div className="mb-4 flex size-12 items-center justify-center rounded-xl bg-secondary">
          <Clock className="size-6 animate-spin text-secondary-foreground" />
        </div>
        <h4 className="text-sm font-bold text-foreground">{t("card.processing")}</h4>
        <span className="mt-1.5 block max-w-xs text-xs leading-normal text-muted-foreground">
          {t("card.processingDesc")}
        </span>
      </Card>
    )
  }

  if (activeAnalysis.status === ExamAnalysisStatus.FAILED) {
    return (
      <Card className="min-h-[300px] flex-1 items-center justify-center gap-0 p-16 text-center">
        <div className="mb-4 flex size-12 items-center justify-center rounded-xl border border-destructive/20 bg-destructive/10">
          <ShieldAlert className="size-6 text-destructive" />
        </div>
        <h4 className="text-sm font-bold text-foreground">{t("card.failed")}</h4>
        <span className="mt-1.5 block max-w-xs text-xs leading-normal text-muted-foreground">
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
      <Card className="flex-1 border-l-4 border-l-destructive bg-destructive/5">
        <CardContent className="flex items-start gap-4">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-lg border border-destructive/20 bg-destructive/10">
            <AlertTriangle className="size-5 text-destructive" />
          </div>
          <div>
            <h4 className="text-sm font-black leading-none text-foreground">
              {t("card.insufficient")}
            </h4>
            <p className="mt-2 text-xs leading-relaxed text-foreground/80">{insufficientMessage}</p>
            <div className="mt-4 rounded-lg border border-destructive/15 bg-card p-3.5">
              <span className="block text-[11px] font-bold text-foreground">
                {t("card.possibleCauses")}
              </span>
              <ul className="mt-1.5 flex list-disc flex-col gap-1 pl-4 text-[11px] leading-normal text-muted-foreground">
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
    <div className="flex flex-1 animate-fade-in flex-col gap-4">
      <Card>
        <CardHeader className="gap-4 border-b sm:flex sm:items-start sm:justify-between">
          <div className="min-w-0">
            <div className="flex items-center gap-2">
              <Sparkles className="size-4 shrink-0 text-primary" />
              <h3 className="truncate text-base font-black text-foreground">
                {analysisPayload.examType}
              </h3>
            </div>
            <span
              className="mt-1.5 block font-mono text-[11px] text-muted-foreground"
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
              className="flex h-9 w-1.5 shrink-0 flex-col justify-end overflow-hidden rounded-full bg-muted"
              role="img"
              aria-label={`${t("card.qualityAssessment")}: ${qualityVerdict.percentage}%`}
            >
              <div
                className={cn("w-full rounded-full", qualityStyles.meter)}
                style={{ height: `${qualityVerdict.percentage}%` }}
              />
            </div>
          </div>
        </CardHeader>

        <CardContent className="flex flex-col gap-4">
          <Alert variant={urgencyStyles.variant}>
            <UrgencyIcon aria-hidden="true" />
            <AlertTitle className="line-clamp-none flex flex-wrap items-baseline gap-x-2 gap-y-1">
              <span className="text-xs font-black uppercase tracking-wider">
                {t(urgencyReading.translationKey)}
              </span>
              <span className="text-[11px] font-normal text-muted-foreground">{t("card.urgencyTitle")}</span>
            </AlertTitle>
            <AlertDescription className="text-[11px] leading-normal">
              <span>{t("card.urgencyCaption")}</span>
              {nextSteps.length > 0 ? (
                <ol className="mt-1.5 flex flex-col gap-1.5">
                  {nextSteps.map((stepItem, index) => (
                    <li key={index} className="flex items-start gap-2 text-xs leading-relaxed text-foreground">
                      <span className="mt-0.5 flex size-4 shrink-0 select-none items-center justify-center rounded border border-border bg-card text-[11px] font-bold text-muted-foreground">
                        {index + 1}
                      </span>
                      {stepItem}
                    </li>
                  ))}
                </ol>
              ) : (
                <p className="mt-1.5 text-[11px] text-muted-foreground">{t("card.noNextSteps")}</p>
              )}
            </AlertDescription>
          </Alert>

          {activeAnalysis.anonymizationReport !== undefined && (
            <AnonymizationSummary report={activeAnalysis.anonymizationReport} />
          )}

          {qualityWarnings.length > 0 && (
            <Alert>
              <Info aria-hidden="true" />
              <AlertTitle className="text-[11px] font-bold">{t("card.qualityObservations")}</AlertTitle>
              <AlertDescription>
                <ul className="flex list-disc flex-col gap-1 pl-4 text-[11px] leading-normal text-foreground/80">
                  {qualityWarnings.map((warningItem, index) => (
                    <li key={index}>{warningItem}</li>
                  ))}
                </ul>
              </AlertDescription>
            </Alert>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardContent>
          <SectionHeading title={t("card.findings")} caption={t("card.findingsCaption")} />

          {analysisPayload.detectedFindings.length === 0 ? (
            <Alert variant="success">
              <CheckCircle2 aria-hidden="true" />
              <AlertDescription className="text-xs text-foreground/80">{t("card.noFindings")}</AlertDescription>
            </Alert>
          ) : (
            <ul className="flex flex-col gap-2.5">
              {analysisPayload.detectedFindings.map((findingItem, index) => {
                const triageReading = resolveFindingTriage(
                  findingItem.severity,
                  findingItem.confidence
                )
                const triageStatus = TRIAGE_TIER_STATUS[triageReading.tier]
                const confidenceReading = resolveConfidenceBand(findingItem.confidence)
                const confidenceStyles = CONFIDENCE_BAND_STYLES[confidenceReading.band]
                const significanceTranslationKey = `card.significance${
                  findingItem.severity.charAt(0).toUpperCase()}${findingItem.severity.slice(1)
                }`

                return (
                  <li
                    key={index}
                    className="rounded-lg border border-border bg-muted/60 p-3.5"
                  >
                    <div className="flex flex-wrap items-start justify-between gap-2">
                      <span className="min-w-0 flex-1 text-xs font-semibold leading-relaxed text-foreground">
                        {findingItem.finding}
                      </span>
                      <div className="flex shrink-0 items-center gap-1.5">
                        <Badge
                          variant="outline"
                          className="rounded-md bg-card text-[11px] font-bold uppercase tracking-wider text-muted-foreground"
                        >
                          {t("card.significance")}: {t(significanceTranslationKey)}
                        </Badge>
                        <StatusBadge
                          status={triageStatus}
                          className="rounded-md text-[11px] font-bold uppercase tracking-wider"
                        >
                          {t(triageReading.translationKey)}
                        </StatusBadge>
                      </div>
                    </div>

                    <div className="mt-2.5 flex items-center gap-2.5">
                      <span className="shrink-0 text-[11px] font-semibold text-muted-foreground">
                        {t("card.confidence")}
                      </span>
                      <div
                        className="h-1.5 flex-1 overflow-hidden rounded-full bg-muted"
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
        </CardContent>
      </Card>

      <Card>
        <CardContent>
          <SectionHeading title={t("card.interpretations")} />

          {analysisPayload.possibleInterpretations.length === 0 ? (
            <EmptySection message={t("card.noInterpretations")} />
          ) : (
            <ul className="flex flex-col gap-2">
              {analysisPayload.possibleInterpretations.map((interpretationItem, index) => (
                <li key={index} className="flex items-start gap-2">
                  <ArrowUpRight className="mt-0.5 size-3.5 shrink-0 text-muted-foreground" aria-hidden="true" />
                  <span className={BODY_TEXT_CLASS}>{interpretationItem}</span>
                </li>
              ))}
            </ul>
          )}

          <Separator className="my-4" />

          <SectionHeading title={t("card.limitations")} />

          {analysisPayload.limitations.length === 0 ? (
            <EmptySection message={t("card.noLimitations")} />
          ) : (
            <ul className="flex flex-col gap-2">
              {analysisPayload.limitations.map((limitationItem, index) => (
                <li key={index} className="flex items-start gap-2">
                  <Info className="mt-0.5 size-3.5 shrink-0 text-muted-foreground" aria-hidden="true" />
                  <span className="text-[11px] leading-relaxed text-muted-foreground">{limitationItem}</span>
                </li>
              ))}
            </ul>
          )}

          <Alert className="mt-4">
            <ShieldAlert aria-hidden="true" />
            <AlertTitle className="text-[11px] font-bold uppercase tracking-wider">
              {t("card.disclaimerTitle")}
            </AlertTitle>
            <AlertDescription className="text-[11px] font-medium leading-relaxed">
              {analysisPayload.disclaimer}
            </AlertDescription>
          </Alert>
        </CardContent>
      </Card>
    </div>
  )
}
