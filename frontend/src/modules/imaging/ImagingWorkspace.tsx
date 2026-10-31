import { useState, useRef } from "react"
import { useTranslation } from "react-i18next"
import { AlertCircle, ArrowLeft, FileQuestion, UploadCloud } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import { Input } from "@/shared/components/ui/input"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { EmptyState } from "@/shared/components/ui/EmptyState"
import { Spinner } from "@/shared/components/ui/spinner"
import { DicomViewport } from "./components/DicomViewport"
import { ImagingStudyDetails } from "./components/ImagingStudyDetails"
import { ImagingUploadProgress } from "./components/ImagingUploadProgress"
import { useDicomViewer } from "./hooks/useDicomViewer"
import { useImagingStudyQuery, useUploadImagingStudyMutation } from "./queries"
import { waitForUploadFrame } from "./utils/pacs_helpers"
import { toast } from "../../shared/store/toast_store"

interface ImagingWorkspaceProps {
  studyId: string
  onBack: () => void
}

export const ImagingWorkspace = ({ studyId, onBack }: ImagingWorkspaceProps) => {
  const { t } = useTranslation("imaging")
  const [uploadState, setUploadState] = useState<{
    percentage: number | null
    status: string | null
  }>({ percentage: null, status: null })
  const fileInputReference = useRef<HTMLInputElement>(null)

  const {
    data: study,
    isLoading: isStudyLoading,
    isError: isStudyError,
    refetch: refetchStudy,
  } = useImagingStudyQuery(studyId)
  const uploadImagingStudyMutation = useUploadImagingStudyMutation()
  const dicomViewer = useDicomViewer(study)

  const handleButtonClick = () => {
    fileInputReference.current?.click()
  }

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = event.target.files?.[0]
    if (!selectedFile || !study) {
      return
    }

    setUploadState({ percentage: 0, status: t("uploadStatus.initial") })
    await waitForUploadFrame(300)
    setUploadState((prev) => ({ ...prev, percentage: 30 }))

    try {
      await uploadImagingStudyMutation.mutateAsync({
        patientFhirId: study.patient_fhir_id,
        title: selectedFile.name.replace(/\.[^/.]+$/, "") || "Nova Ressonância Magnética",
        modality: "MR",
        dicomBlob: selectedFile,
      })

      setUploadState({ percentage: 100, status: t("uploadStatus.grpcCompleted") })
      await waitForUploadFrame(500)
      toast.success(t("toast.uploadSuccess"))
    } catch {
      toast.error(t("toast.uploadError"))
    } finally {
      setUploadState({ percentage: null, status: null })
      if (fileInputReference.current) {
        fileInputReference.current.value = ""
      }
    }
  }

  if (!studyId) {
    return null
  }

  if (isStudyError) {
    return (
      <div className="py-16">
        <Alert variant="destructive">
          <AlertCircle aria-hidden="true" />
          <AlertDescription className="flex flex-wrap items-center gap-3">
            <span>{t("loadError")}</span>
            <Button variant="outline" size="sm" onClick={() => refetchStudy()}>
              {t("retry")}
            </Button>
          </AlertDescription>
        </Alert>
      </div>
    )
  }

  if (!isStudyLoading && !study) {
    return (
      <div className="py-16">
        <EmptyState
          icon={FileQuestion}
          title={t("studyNotFound")}
          description={t("studyNotFoundDesc")}
          action={
            <Button variant="outline" onClick={onBack}>
              <ArrowLeft className="w-4 h-4" />
              {t("backToRecord")}
            </Button>
          }
        />
      </div>
    )
  }

  if (isStudyLoading || !study) {
    return (
      <div className="flex items-center justify-center gap-2 py-16 text-sm text-muted-foreground">
        <Spinner />
        <span>{t("loading")}</span>
      </div>
    )
  }

  return (
    <div className="flex-1 p-8 flex flex-col gap-6 max-w-7xl mx-auto w-full select-none">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div className="flex items-center gap-4">
          <Button variant="outline" onClick={onBack}>
            <ArrowLeft className="w-4 h-4" />
            {t("backToRecord")}
          </Button>
          <div className="text-left">
            <h2 className="text-xl font-display font-bold text-foreground leading-none">
              {t("titleConsole")}
            </h2>
            <span className="text-xs text-muted-foreground mt-1.5 block">
              {t("studyLabel")}: {study.title} • {t("uidLabel")}: {study.study_instance_uid}
            </span>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Input
            type="file"
            accept=".dcm"
            ref={fileInputReference}
            onChange={handleFileChange}
            className="hidden"
          />
          <Button
            variant="outline"
            onClick={handleButtonClick}
            disabled={uploadImagingStudyMutation.isPending || uploadState.percentage !== null}
            className="border-primary/20 text-primary hover:bg-primary/5"
          >
            <UploadCloud className="w-4 h-4" />
            {t("uploadDcm")}
          </Button>
        </div>
      </div>

      {uploadState.percentage !== null && uploadState.status && (
        <ImagingUploadProgress percentage={uploadState.percentage} status={uploadState.status} />
      )}

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6 items-start">
        <ImagingStudyDetails study={study} />
        <DicomViewport
          activeTool={dicomViewer.activeTool}
          canvasReference={dicomViewer.canvasReference}
          onMouseDown={dicomViewer.handleMouseDown}
          onMouseMove={dicomViewer.handleMouseMove}
          onMouseUp={dicomViewer.handleMouseUp}
          onToolChange={dicomViewer.setActiveTool}
          onPresetChange={dicomViewer.applyPreset}
        />
      </div>
    </div>
  )
}
