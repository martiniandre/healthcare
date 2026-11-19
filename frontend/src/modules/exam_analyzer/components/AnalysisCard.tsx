import { useTranslation } from "react-i18next"
import { HelpCircle } from "lucide-react"
import { EmptyState } from "@/shared/components/ui/EmptyState"
import { ExamAnalysisStatus, type ExamAnalysis } from "../types"
import { AnalysisFailedState, AnalysisInsufficientState, AnalysisProcessingState } from "./AnalysisStatusStates"
import { AnalysisResult } from "./AnalysisResult"

interface AnalysisCardProperties {
  activeAnalysis: ExamAnalysis | null
}

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

  switch (activeAnalysis.status) {
    case ExamAnalysisStatus.PENDING:
    case ExamAnalysisStatus.PROCESSING:
      return <AnalysisProcessingState />
    case ExamAnalysisStatus.FAILED:
      return <AnalysisFailedState />
    case ExamAnalysisStatus.INSUFFICIENT_DATA:
      return <AnalysisInsufficientState activeAnalysis={activeAnalysis} />
    default:
      return <AnalysisResult activeAnalysis={activeAnalysis} />
  }
}
