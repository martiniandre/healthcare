import { useTranslation } from "react-i18next"
import { AlertCircle, ArrowLeft, FileQuestion } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { EmptyState } from "@/shared/components/ui/EmptyState"
import { Spinner } from "@/shared/components/ui/spinner"
import { DicomViewport } from "./components/DicomViewport"
import { ImagingStudyDetails } from "./components/ImagingStudyDetails"
import { DicomUploadControl } from "./components/DicomUploadControl"
import { useDicomViewer } from "./hooks/useDicomViewer"
import { useImagingStudyQuery } from "./queries"

interface ImagingWorkspaceProps {
  studyId: string
  onBack: () => void
}

export const ImagingWorkspace = ({ studyId, onBack }: ImagingWorkspaceProps) => {
  const { t } = useTranslation("imaging")
  const {
    data: study,
    isLoading: isStudyLoading,
    isError: isStudyError,
    refetch: refetchStudy,
  } = useImagingStudyQuery(studyId)
  const dicomViewer = useDicomViewer(study)

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
      </div>

      <DicomUploadControl patientFhirId={study.patient_fhir_id} />

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
