import { useState, useRef } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { ArrowLeft, UploadCloud, Image as ImageIcon } from "lucide-react"
import { Button } from "../../shared/components/ui/Button"
import { DicomViewport } from "./components/DicomViewport"
import { ImagingStudyDetails } from "./components/ImagingStudyDetails"
import { ImagingUploadProgress } from "./components/ImagingUploadProgress"
import { useDicomViewer } from "./hooks/useDicomViewer"
import { useImagingStudyQuery, useUploadImagingStudyMutation } from "./queries"
import { waitForUploadFrame } from "./utils/pacs_helpers"
import { toast } from "../../shared/store/toast_store"

export const ImagingWorkspace = () => {
  const { t } = useTranslation()
  const { studyId = "" } = useParams<{ studyId: string }>()
  const navigate = useNavigate()
  const [uploadState, setUploadState] = useState<{
    percentage: number | null
    status: string | null
  }>({ percentage: null, status: null })
  const fileInputReference = useRef<HTMLInputElement>(null)

  const { data: study, isLoading: isStudyLoading } = useImagingStudyQuery(studyId)
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

    setUploadState({ percentage: 0, status: t("imaging.uploadStatus.initial") })
    await waitForUploadFrame(300)
    setUploadState((prev) => ({ ...prev, percentage: 30 }))

    try {
      await uploadImagingStudyMutation.mutateAsync({
        patientFhirId: study.patient_fhir_id,
        title: selectedFile.name.replace(/\.[^/.]+$/, "") || "Nova Ressonância Magnética",
        modality: "MR",
        dicomBlob: selectedFile,
      })

      setUploadState({ percentage: 100, status: t("imaging.uploadStatus.grpcCompleted") })
      await waitForUploadFrame(500)
      toast.success(t("imaging.toast.uploadSuccess"))
      window.alert(t("imaging.alert.uploadSuccess"))
    } catch {
      toast.error(t("imaging.toast.uploadError"))
    } finally {
      setUploadState({ percentage: null, status: null })
      if (fileInputReference.current) {
        fileInputReference.current.value = ""
      }
    }
  }

  if (!studyId) {
    return (
      <div className="flex-1 p-8 flex flex-col items-center justify-center text-center">
        <div className="bg-white p-8 rounded-2xl border border-border shadow-sm max-w-md w-full">
          <ImageIcon className="w-12 h-12 text-primary/40 mx-auto mb-4" />
          <h2 className="text-xl font-bold text-gray-900 mb-2">PACS Viewer</h2>
          <p className="text-sm text-muted mb-6">
            Please select a patient from the patients list, and then select an imaging study to view it in the PACS viewer.
          </p>
          <Button variantType="primary" onClick={() => navigate("/")}>
            Go to Patients
          </Button>
        </div>
      </div>
    )
  }

  if (isStudyLoading || !study) {
    return (
      <div className="text-center py-16">
        <span className="text-sm text-muted">{t("imaging.loading")}</span>
      </div>
    )
  }

  return (
    <div className="flex-1 p-8 flex flex-col gap-6 max-w-7xl mx-auto w-full select-none">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div className="flex items-center gap-4">
          <Button variantType="outline" onClick={() => navigate(`/patients/${study.patient_fhir_id}`)} className="px-3">
            <ArrowLeft className="w-4 h-4" />
            {t("imaging.backToRecord")}
          </Button>
          <div className="text-left">
            <h2 className="text-xl font-black text-gray-900 leading-none">
              {t("imaging.titleConsole")}
            </h2>
            <span className="text-xs text-muted mt-1.5 block">
              Estudo: {study.title} • UID: {study.study_instance_uid}
            </span>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <input
            type="file"
            accept=".dcm"
            ref={fileInputReference}
            onChange={handleFileChange}
            className="hidden"
          />
          <Button
            variantType="outline"
            onClick={handleButtonClick}
            disabled={uploadImagingStudyMutation.isPending || uploadState.percentage !== null}
            className="px-3.5 gap-2 border-primary/20 text-primary hover:bg-primary/5"
          >
            <UploadCloud className="w-4 h-4" />
            {t("imaging.uploadDcm")}
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
