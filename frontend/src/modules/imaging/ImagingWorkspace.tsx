import { useState, useRef } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { ArrowLeft, UploadCloud } from "lucide-react"
import { Button } from "../../shared/components/ui/Button"
import { DicomViewport } from "./components/DicomViewport"
import { ImagingStudyDetails } from "./components/ImagingStudyDetails"
import { ImagingUploadProgress } from "./components/ImagingUploadProgress"
import { useDicomViewer } from "./hooks/useDicomViewer"
import { useImagingStudyQuery, useUploadImagingStudyMutation } from "./queries"
import { waitForUploadFrame } from "./utils/pacs_helpers"
import { toast } from "../../shared/store/toast_store"

export const ImagingWorkspace = () => {
  const { studyId = "" } = useParams<{ studyId: string }>()
  const navigate = useNavigate()
  const [uploadPercentage, setUploadPercentage] = useState<number | null>(null)
  const [uploadStatus, setUploadStatus] = useState<string | null>(null)
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

    setUploadPercentage(0)
    setUploadStatus("Iniciando upload e validação de assinatura DICOM...")
    await waitForUploadFrame(300)
    setUploadPercentage(30)

    try {
      await uploadImagingStudyMutation.mutateAsync({
        patientFhirId: study.patient_fhir_id,
        title: selectedFile.name.replace(/\.[^/.]+$/, "") || "Nova Ressonância Magnética",
        modality: "MR",
        dicomBlob: selectedFile,
      })

      setUploadPercentage(100)
      setUploadStatus("Transmissão gRPC concluída. Processando DICOM metadata...")
      await waitForUploadFrame(500)
      toast.success("DICOM carregado e processado com sucesso no barramento do PACS!")
    } catch {
      toast.error("Falha no upload do arquivo DICOM.")
    } finally {
      setUploadPercentage(null)
      setUploadStatus(null)
      if (fileInputReference.current) {
        fileInputReference.current.value = ""
      }
    }
  }

  if (isStudyLoading || !study) {
    return (
      <div className="text-center py-16">
        <span className="text-sm text-muted">Carregando visualizador PACS...</span>
      </div>
    )
  }

  return (
    <div className="flex-1 p-8 flex flex-col gap-6 max-w-7xl mx-auto w-full select-none">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div className="flex items-center gap-4">
          <Button variantType="outline" onClick={() => navigate(`/patients/${study.patient_fhir_id}`)} className="px-3">
            <ArrowLeft className="w-4 h-4" />
            Voltar Prontuário
          </Button>
          <div className="text-left">
            <h2 className="text-xl font-black text-gray-900 leading-none">
              Console Cirúrgico PACS
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
            disabled={uploadImagingStudyMutation.isPending || uploadPercentage !== null}
            className="px-3.5 gap-2 border-primary/20 text-primary hover:bg-primary/5"
          >
            <UploadCloud className="w-4 h-4" />
            Upload Novo .DCM
          </Button>
        </div>
      </div>

      {uploadPercentage !== null && uploadStatus && (
        <ImagingUploadProgress percentage={uploadPercentage} status={uploadStatus} />
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
