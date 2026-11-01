import { useState, useEffect, useRef } from "react"
import { useQueryClient } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"
import { AlertCircle, Sparkles } from "lucide-react"
import { FileUploader, type UploadOutcome } from "./components/FileUploader"
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
import { Badge } from "@/shared/components/ui/badge"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { Button } from "@/shared/components/ui/button"
import { PageContainer } from "../../shared/components/ui/PageContainer"
import { extractApiErrorDetails } from "@/shared/utils/api_error"
import type { ExamAnalysis } from "./types"

export interface ExamAnalyzerProps {
  patientFhirId?: string
}

export const ExamAnalyzer = ({ patientFhirId }: ExamAnalyzerProps = {}) => {
  const { t } = useTranslation("examAnalyzer")
  const [selectedAnalysisID, setSelectedAnalysisID] = useState<string | null>(null)
  const [uploadPercentageValue, setUploadPercentageValue] = useState<number | null>(null)

  const uploadAbortControllerReference = useRef<AbortController | null>(null)

  const queryClient = useQueryClient()

  const {
    data: rawAnalysesHistory = [],
    isLoading: isHistoryLoading,
    isError: isHistoryError,
    refetch: refetchHistory,
  } = useExamAnalysesQuery(patientFhirId)
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

  const {
    data: polledAnalysisDetails,
    isError: isDetailsError,
    refetch: refetchDetails,
  } = useExamAnalysisQuery(
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

  useEffect(() => {
    if (uploadPercentageValue === null) {
      return
    }
    const warnBeforeLeaving = (event: BeforeUnloadEvent) => {
      event.preventDefault()
      event.returnValue = ""
    }
    window.addEventListener("beforeunload", warnBeforeLeaving)
    return () => window.removeEventListener("beforeunload", warnBeforeLeaving)
  }, [uploadPercentageValue])

  const handleCancelUpload = () => {
    uploadAbortControllerReference.current?.abort()
  }

  const handleFileUpload = async (
    file: File,
    consent: boolean,
    anonymize: boolean,
    acknowledgeExternalProcessing: boolean
  ): Promise<UploadOutcome> => {
    const abortController = new AbortController()
    uploadAbortControllerReference.current = abortController
    setUploadPercentageValue(0)
    try {
      const createdRecord = await uploadExamMutation.mutateAsync({
        file,
        consent,
        anonymize,
        acknowledgeExternalProcessing,
        patientFhirId,
        signal: abortController.signal,
        onUploadProgress: (percentage) => {
          setUploadPercentageValue(percentage)
        },
      })
      
      toast.success(t("toast.uploadSuccess"))
      setSelectedAnalysisID(createdRecord.id)
      return { status: "success" }
    } catch (uploadError) {
      if (abortController.signal.aborted) {
        return { status: "cancelled" }
      }
      toast.error(t("toast.uploadError"))
      return { status: "failed", errorCode: extractApiErrorDetails(uploadError).code }
    } finally {
      uploadAbortControllerReference.current = null
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
    <PageContainer className="gap-5">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between sm:gap-4">
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <h2 className="font-display text-xl font-bold leading-none tracking-tight text-foreground">
              {t("title")}
            </h2>
            <Badge variant="info" className="gap-1.5 text-[11px] font-bold uppercase tracking-wider">
              <Sparkles aria-hidden="true" />
              {t("badge")}
            </Badge>
          </div>
          <p className="mt-1.5 block max-w-2xl text-xs leading-relaxed text-muted-foreground">
            {t("subtitle")}
          </p>
        </div>
      </div>

      <div className="flex flex-col items-stretch gap-5 md:flex-row">
        <AnalysisHistory
          history={analysesHistory}
          isLoading={isHistoryLoading}
          isError={isHistoryError}
          onRetry={() => refetchHistory()}
          activeID={selectedAnalysisID}
          onSelect={handleSelectAnalysis}
          onDelete={handleDeleteAnalysis}
        />

        <div className="flex flex-1 flex-col gap-5">
          <FileUploader
            onUpload={handleFileUpload}
            onCancel={handleCancelUpload}
            isPending={uploadExamMutation.isPending}
            uploadProgress={uploadPercentageValue}
          />

          {isDetailsError && selectedAnalysisID ? (
            <Alert variant="destructive">
              <AlertCircle aria-hidden="true" />
              <AlertDescription className="flex flex-wrap items-center gap-3">
                <span>{t("card.loadError")}</span>
                <Button variant="outline" size="sm" onClick={() => refetchDetails()}>
                  {t("retry")}
                </Button>
              </AlertDescription>
            </Alert>
          ) : null}

          <AnalysisCard activeAnalysis={activeAnalysisToRender} />
        </div>
      </div>
    </PageContainer>
  )
}
export default ExamAnalyzer
