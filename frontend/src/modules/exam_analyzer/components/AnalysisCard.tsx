import { useTranslation } from "react-i18next"
import { AlertTriangle, Clock, ShieldAlert, Sparkles, HelpCircle, Activity } from "lucide-react"
import { Card } from "../../../shared/components/ui/Card"
import { ExamAnalysisStatus, type ExamAnalysis, type MedicalAnalysisResponse } from "../types"

interface AnalysisCardProperties {
  activeAnalysis: ExamAnalysis | null
}

export const AnalysisCard = ({ activeAnalysis }: AnalysisCardProperties) => {
  const { t } = useTranslation("examAnalyzer")

  if (!activeAnalysis) {
    return (
      <div className="flex-1 flex flex-col items-center justify-center border border-dashed border-gray-200 rounded-xl p-16 text-center bg-white min-h-[300px]">
        <div className="w-12 h-12 rounded-xl bg-gray-50 flex items-center justify-center mb-4 border border-border">
          <HelpCircle className="w-6 h-6 text-gray-300" />
        </div>
        <h4 className="text-sm font-bold text-gray-800">
          {t("card.noExam")}
        </h4>
        <span className="text-xs text-muted max-w-xs block mt-1.5 leading-normal">
          {t("card.noExamDesc")}
        </span>
      </div>
    )
  }

  if (activeAnalysis.status === ExamAnalysisStatus.PENDING || activeAnalysis.status === ExamAnalysisStatus.PROCESSING) {
    return (
      <Card glowingType="amethyst" className="flex-1 flex flex-col items-center justify-center p-16 text-center bg-white min-h-[300px] animate-pulse">
        <div className="w-12 h-12 rounded-xl bg-secondary/8 flex items-center justify-center mb-4 animate-spin">
          <Clock className="w-6 h-6 text-secondary" />
        </div>
        <h4 className="text-sm font-bold text-gray-800">
          {t("card.processing")}
        </h4>
        <span className="text-xs text-muted max-w-xs block mt-1.5 leading-normal">
          {t("card.processingDesc")}
        </span>
      </Card>
    )
  }

  if (activeAnalysis.status === ExamAnalysisStatus.FAILED) {
    return (
      <Card glowingType="none" className="flex-1 flex flex-col items-center justify-center p-16 text-center bg-white min-h-[300px]">
        <div className="w-12 h-12 rounded-xl bg-red-50 flex items-center justify-center mb-4 border border-red-100">
          <ShieldAlert className="w-6 h-6 text-danger" />
        </div>
        <h4 className="text-sm font-bold text-gray-800">
          {t("card.failed")}
        </h4>
        <span className="text-xs text-muted max-w-xs block mt-1.5 leading-normal">
          {t("card.failedDesc")}
        </span>
      </Card>
    )
  }

  if (activeAnalysis.status === ExamAnalysisStatus.INSUFFICIENT_DATA) {
    const insufficientMessage = (activeAnalysis.analysis_response as { message: string })?.message || 
      t("card.insufficientDefault")

    return (
      <Card glowingType="none" className="flex-1 border-l-4 border-l-danger p-6 bg-red-50/30 rounded-xl">
        <div className="flex items-start gap-4">
          <div className="w-10 h-10 rounded-lg bg-red-50 flex items-center justify-center shrink-0 border border-red-100">
            <AlertTriangle className="w-5 h-5 text-danger" />
          </div>
          <div>
            <h4 className="text-sm font-black text-gray-900 leading-none">
              {t("card.insufficient")}
            </h4>
            <p className="text-xs text-gray-600 mt-2 leading-relaxed">
              {insufficientMessage}
            </p>
            <div className="mt-4 p-3.5 bg-white border border-red-100 rounded-lg">
              <span className="text-[11px] font-bold text-gray-800 block">
                {t("card.possibleCauses")}
              </span>
              <ul className="list-disc list-inside text-[10px] text-gray-500 mt-1.5 flex flex-col gap-1">
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

  return (
    <div className="flex-1 flex flex-col gap-6 animate-fade-in">
      <Card glowingType="cyan" className="p-6 bg-white border border-border rounded-xl">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-border/80 pb-4 mb-4">
          <div>
            <div className="flex items-center gap-2">
              <Sparkles className="w-4.5 h-4.5 text-primary" />
              <h3 className="text-base font-black text-gray-900">{analysisPayload.examType}</h3>
            </div>
            <span className="text-[10px] text-muted font-mono mt-1.5 block">
              {t("card.analysisId")} {activeAnalysis.id}
            </span>
          </div>

          <div className="flex items-center gap-4.5 bg-gray-50 px-4 py-2 border border-border/80 rounded-xl self-start sm:self-auto shrink-0">
            <div>
              <span className="text-[9px] text-muted font-bold uppercase tracking-wider block">
                {t("card.qualityAssessment")}
              </span>
              <span className="text-sm font-black text-gray-900 mt-0.5 block">
                {(analysisPayload.qualityAssessment.score * 100).toFixed(0)}%
              </span>
            </div>
            <div className="w-1.5 h-8 bg-border rounded-full overflow-hidden shrink-0">
              <div 
                className="h-full bg-secondary transition-all"
                style={{ height: `${analysisPayload.qualityAssessment.score * 100}%` }}
              />
            </div>
          </div>
        </div>

        {analysisPayload.qualityAssessment.warnings.length > 0 && (
          <div className="p-3 bg-amber-50 border border-amber-100 rounded-lg flex items-start gap-2.5 mb-5">
            <AlertTriangle className="w-4 h-4 text-amber-500 mt-0.5 shrink-0" />
            <div className="flex flex-col gap-0.5">
              <span className="text-[10px] font-bold text-amber-800">
                {t("card.warnings")}
              </span>
              {analysisPayload.qualityAssessment.warnings.map((warningItem, index) => (
                <span key={index} className="text-[10px] text-amber-700 leading-normal">
                  • {warningItem}
                </span>
              ))}
            </div>
          </div>
        )}

        <div className="flex flex-col gap-4">
          <div>
            <h4 className="text-xs font-bold text-gray-800 mb-2.5">
              {t("card.findings")}
            </h4>
            <div className="grid grid-cols-1 gap-3">
              {analysisPayload.detectedFindings.map((findingItem, index) => (
                <div key={index} className="border border-border/80 rounded-lg p-3.5 bg-gray-50/50 flex flex-col gap-2.5">
                  <div className="flex items-start justify-between gap-4">
                    <span className="text-xs font-semibold text-gray-800 leading-relaxed">
                      {findingItem.finding}
                    </span>
                    <span className={`text-[9px] font-bold px-2 py-0.5 rounded-full border uppercase shrink-0 ${
                      findingItem.severity === "high" 
                        ? "bg-orange-50 text-orange-600 border-orange-100"
                        : findingItem.severity === "medium"
                          ? "bg-amber-50 text-amber-600 border-amber-100"
                          : "bg-blue-50 text-blue-600 border-blue-100"
                    }`}>
                      {findingItem.severity === "high" && t("card.severityHigh")}
                      {findingItem.severity === "medium" && t("card.severityMedium")}
                      {findingItem.severity === "low" && t("card.severityLow")}
                    </span>
                  </div>

                  <div className="flex items-center gap-3">
                    <span className="text-[10px] text-muted shrink-0">
                      {t("card.confidence")}
                    </span>
                    <div className="flex-1 h-2 bg-gray-100 rounded-full overflow-hidden">
                      <div 
                        className="h-full bg-primary transition-all"
                        style={{ width: `${findingItem.confidence * 100}%` }}
                      />
                    </div>
                    <span className="text-[10px] font-bold text-gray-700 shrink-0">
                      {(findingItem.confidence * 100).toFixed(0)}%
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="h-px bg-border/60" />

          <div>
            <h4 className="text-xs font-bold text-gray-800 mb-2">
              {t("card.interpretations")}
            </h4>
            <div className="flex flex-col gap-2">
              {analysisPayload.possibleInterpretations.map((interpretationItem, index) => (
                <div key={index} className="flex items-start gap-2">
                  <div className="w-1.5 h-1.5 rounded-full bg-gray-400 mt-1.5 shrink-0" />
                  <span className="text-xs text-gray-600 leading-relaxed">{interpretationItem}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="h-px bg-border/60" />

          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div>
              <div className="flex items-center gap-2 mb-2.5">
                <Activity className="w-4 h-4 text-secondary" />
                <h4 className="text-xs font-bold text-gray-800">
                  {t("card.recommendations")}
                </h4>
              </div>
              <div className="flex flex-col gap-3">
                <div className="flex items-center gap-2">
                  <span className="text-[10px] text-muted">
                    {t("card.urgency")}
                  </span>
                  <span className={`text-[9px] font-black px-2 py-0.5 rounded border uppercase ${
                    analysisPayload.recommendation.urgency === "urgent"
                      ? "bg-red-50 text-red-600 border-red-200"
                      : analysisPayload.recommendation.urgency === "medical_followup"
                        ? "bg-amber-50 text-amber-600 border-amber-200"
                        : "bg-blue-50 text-blue-600 border-blue-200"
                  }`}>
                    {analysisPayload.recommendation.urgency === "urgent" && t("card.urgencyUrgent")}
                    {analysisPayload.recommendation.urgency === "medical_followup" && t("card.urgencyFollowup")}
                    {analysisPayload.recommendation.urgency === "normal" && t("card.urgencyNormal")}
                  </span>
                </div>
                <div className="flex flex-col gap-2">
                  {analysisPayload.recommendation.nextSteps.map((stepItem, index) => (
                    <div key={index} className="flex items-start gap-2.5 p-2 bg-gray-50 rounded border border-border/50">
                      <div className="w-4.5 h-4.5 rounded bg-white border border-border flex items-center justify-center text-[10px] font-bold text-gray-500 shrink-0 select-none">
                        {index + 1}
                      </div>
                      <span className="text-xs text-gray-600 leading-normal">{stepItem}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div>
              <div className="flex items-center gap-2 mb-2.5">
                <AlertTriangle className="w-4 h-4 text-muted" />
                <h4 className="text-xs font-bold text-gray-800">
                  {t("card.limitations")}
                </h4>
              </div>
              <div className="flex flex-col gap-2">
                {analysisPayload.limitations.map((limitationItem, index) => (
                  <div key={index} className="flex items-start gap-2">
                    <div className="w-1.5 h-1.5 rounded-full bg-gray-400 mt-1.5 shrink-0" />
                    <span className="text-[11px] text-gray-500 leading-relaxed">{limitationItem}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="p-4 bg-gray-50 border border-border rounded-xl mt-2">
            <span className="text-[10px] font-bold text-gray-800 uppercase tracking-wider block mb-1">
              {t("card.disclaimerTitle")}
            </span>
            <p className="text-[10px] text-gray-500 leading-relaxed font-medium">
              {analysisPayload.disclaimer}
            </p>
          </div>
        </div>
      </Card>
    </div>
  )
}
