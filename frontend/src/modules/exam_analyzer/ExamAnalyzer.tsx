import { useState, useEffect } from "react"
import { useQueryClient } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"
import { Sparkles } from "lucide-react"
import { FileUploader } from "./components/FileUploader"
import { AnalysisCard } from "./components/AnalysisCard"
import { AnalysisHistory } from "./components/AnalysisHistory"
import {
  useExamAnalysesQuery,
  useExamAnalysisQuery,
  useUploadExamMutation,
  useDeleteAnalysisMutation,
  examAnalyzerKeys,
} from "./queries"
import { toast } from "../../shared/store/toast_store"
import { PageContainer } from "../../shared/components/ui/PageContainer"
import type { ExamAnalysis } from "./types"

export interface ExamAnalyzerProps {
  patientFhirId?: string
}

export const ExamAnalyzer = ({ patientFhirId }: ExamAnalyzerProps = {}) => {
  const { t } = useTranslation("examAnalyzer")
  const [selectedAnalysisID, setSelectedAnalysisID] = useState<string | null>(null)
  const [uploadPercentageValue, setUploadPercentageValue] = useState<number | null>(null)

  const queryClient = useQueryClient()

  const { data: rawAnalysesHistory = [], isLoading: isHistoryLoading } = useExamAnalysesQuery(patientFhirId)
  const analysesHistory = rawAnalysesHistory || []
  const uploadExamMutation = useUploadExamMutation()
  const deleteAnalysisMutation = useDeleteAnalysisMutation()

  const activeRecordInHistory = analysesHistory.find(
    (item: ExamAnalysis) => item.id === selectedAnalysisID
  )

  const shouldPollForUpdates =
    !!selectedAnalysisID &&
    (activeRecordInHistory?.status === "pending" ||
      activeRecordInHistory?.status === "processing")

  const { data: polledAnalysisDetails } = useExamAnalysisQuery(
    selectedAnalysisID || "",
    {
      enabled: !!selectedAnalysisID,
      refetchInterval: shouldPollForUpdates ? 2500 : undefined,
    }
  )

  useEffect(() => {
    if (polledAnalysisDetails) {
      if (
        polledAnalysisDetails.status === "completed" ||
        polledAnalysisDetails.status === "failed" ||
        polledAnalysisDetails.status === "insufficient_data"
      ) {
        toast.info(t("toast.analysisCompleted", { fileName: polledAnalysisDetails.file_name }))
        queryClient.invalidateQueries({ queryKey: examAnalyzerKeys.all })
      }
    }
  }, [polledAnalysisDetails, queryClient, t])

  const handleFileUpload = async (file: File, consent: boolean, anonymize: boolean) => {
    setUploadPercentageValue(0)
    try {
      const createdRecord = await uploadExamMutation.mutateAsync({
        file,
        consent,
        anonymize,
        patientFhirId,
        onUploadProgress: (percentage) => {
          setUploadPercentageValue(percentage)
        },
      })
      
      toast.success(t("toast.uploadSuccess"))
      setSelectedAnalysisID(createdRecord.id)
    } catch {
      toast.error(t("toast.uploadError"))
    } finally {
      setUploadPercentageValue(null)
    }
  }

  const handleSelectAnalysis = (analysis: ExamAnalysis) => {
    setSelectedAnalysisID(analysis.id)
  }

  const handleDeleteAnalysis = async (id: string) => {
    try {
      await deleteAnalysisMutation.mutateAsync(id)
      toast.success(t("toast.deleteSuccess"))
      if (selectedAnalysisID === id) {
        setSelectedAnalysisID(null)
      }
    } catch {
      toast.error(t("toast.deleteError"))
    }
  }

  const activeAnalysisToRender = (() => {
    if (polledAnalysisDetails) {
      return polledAnalysisDetails
    }
    return activeRecordInHistory || null
  })()

  return (
    <PageContainer className="gap-6 select-none">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h2 className="text-xl font-display font-bold text-gray-900 tracking-tight leading-none">
              {t("title")}
            </h2>
            <div className="flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-primary/8 border border-primary/10">
              <Sparkles className="w-3 h-3 text-primary animate-pulse" />
              <span className="text-[10px] font-bold text-primary">{t("badge")}</span>
            </div>
          </div>
          <span className="text-xs text-muted mt-1.5 block">
            {t("subtitle")}
          </span>
        </div>
      </div>

      <div className="flex flex-col md:flex-row gap-6 items-stretch">
        <AnalysisHistory
          history={analysesHistory}
          isLoading={isHistoryLoading}
          activeID={selectedAnalysisID}
          onSelect={handleSelectAnalysis}
          onDelete={handleDeleteAnalysis}
        />

        <div className="flex-1 flex flex-col gap-6">
          <FileUploader
            onUpload={handleFileUpload}
            isPending={uploadExamMutation.isPending}
            uploadProgress={uploadPercentageValue}
          />

          <AnalysisCard activeAnalysis={activeAnalysisToRender} />
        </div>
      </div>
    </PageContainer>
  )
}
export default ExamAnalyzer
