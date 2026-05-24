import { useState, useRef } from "react"
import { UploadCloud, CheckSquare, Square, FileText, X } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"
import { Card } from "../../../shared/components/ui/Card"

interface FileUploaderProperties {
  onUpload: (file: File, consent: boolean, anonymize: boolean) => void
  isPending: boolean
  uploadProgress: number | null
}

export const FileUploader = ({ onUpload, isPending, uploadProgress }: FileUploaderProperties) => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [isDragActive, setIsDragActive] = useState<boolean>(false)
  const [consentChecked, setConsentChecked] = useState<boolean>(false)
  const [anonymizeChecked, setAnonymizeChecked] = useState<boolean>(false)
  
  const fileInputReference = useRef<HTMLInputElement>(null)

  const handleDragOver = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault()
    setIsDragActive(true)
  }

  const handleDragLeave = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault()
    setIsDragActive(false)
  }

  const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault()
    setIsDragActive(false)
    const file = event.dataTransfer.files?.[0]
    if (file) {
      setSelectedFile(file)
    }
  }

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      setSelectedFile(file)
    }
  }

  const handleTriggerSelect = () => {
    fileInputReference.current?.click()
  }

  const handleClearFile = () => {
    setSelectedFile(null)
    if (fileInputReference.current) {
      fileInputReference.current.value = ""
    }
  }

  const handleFormSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    if (selectedFile && consentChecked) {
      onUpload(selectedFile, consentChecked, anonymizeChecked)
    }
  }

  return (
    <Card glowingType="cyan" className="p-6 bg-white border border-border rounded-xl">
      <h3 className="text-base font-bold text-gray-900 mb-2">Enviar Exame para Análise</h3>
      <span className="text-xs text-muted block mb-5 leading-normal">
        Arraste arquivos de exames radiológicos, fotos clínicas ou PDFs laboratoriais. A análise é assistiva e probabilística.
      </span>

      <form onSubmit={handleFormSubmit} className="flex flex-col gap-5">
        <div
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onClick={handleTriggerSelect}
          className={`border-2 border-dashed rounded-xl p-8 flex flex-col items-center justify-center gap-3 cursor-pointer transition-all duration-300 ${
            isDragActive
              ? "border-primary bg-primary/5 scale-[1.01]"
              : "border-gray-200 hover:border-primary/50 hover:bg-gray-50/50"
          }`}
        >
          <input
            type="file"
            ref={fileInputReference}
            onChange={handleFileSelect}
            className="hidden"
            accept="image/*,.pdf"
          />

          <div className="w-12 h-12 rounded-full bg-primary/8 flex items-center justify-center">
            <UploadCloud className="w-6 h-6 text-primary" />
          </div>

          <div className="text-center">
            <span className="text-sm font-semibold text-gray-800 block">
              Selecione ou solte o arquivo aqui
            </span>
            <span className="text-[11px] text-muted block mt-1">
              Imagens (PNG, JPG, DICOM) ou arquivos PDF de até 15MB
            </span>
          </div>
        </div>

        {selectedFile && (
          <div className="flex items-center justify-between p-3.5 bg-gray-50 border border-border/80 rounded-lg animate-fade-in">
            <div className="flex items-center gap-3 min-w-0">
              <FileText className="w-5 h-5 text-primary shrink-0" />
              <div className="min-w-0">
                <span className="text-xs font-semibold text-gray-800 block truncate">
                  {selectedFile.name}
                </span>
                <span className="text-[10px] text-muted block mt-0.5">
                  {(selectedFile.size / (1024 * 1024)).toFixed(2)} MB
                </span>
              </div>
            </div>
            <button
              type="button"
              onClick={handleClearFile}
              className="p-1 rounded-md text-gray-400 hover:text-red-500 hover:bg-red-50 transition-all cursor-pointer"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        )}

        <div className="flex flex-col gap-3">
          <div
            onClick={() => setConsentChecked((previous) => !previous)}
            className="flex items-start gap-3 cursor-pointer select-none group"
          >
            <div className="mt-0.5 text-primary">
              {consentChecked ? (
                <CheckSquare className="w-4.5 h-4.5 group-hover:scale-105 transition-transform" />
              ) : (
                <Square className="w-4.5 h-4.5 text-gray-400 group-hover:scale-105 transition-transform" />
              )}
            </div>
            <div className="flex-1">
              <span className="text-xs font-semibold text-gray-700 block">
                Consentimento do Paciente
              </span>
              <span className="text-[10px] text-muted block mt-0.5 leading-normal">
                Confirmo que possuo a autorização expressa do paciente para submeter seus dados e imagens para processamento clínico assistivo.
              </span>
            </div>
          </div>

          <div
            onClick={() => setAnonymizeChecked((previous) => !previous)}
            className="flex items-start gap-3 cursor-pointer select-none group"
          >
            <div className="mt-0.5 text-secondary">
              {anonymizeChecked ? (
                <CheckSquare className="w-4.5 h-4.5 group-hover:scale-105 transition-transform" />
              ) : (
                <Square className="w-4.5 h-4.5 text-gray-400 group-hover:scale-105 transition-transform" />
              )}
            </div>
            <div className="flex-1">
              <span className="text-xs font-semibold text-gray-700 block">
                Anonimização de Segurança (Recomendado)
              </span>
              <span className="text-[10px] text-muted block mt-0.5 leading-normal">
                Substituir o nome do arquivo enviado por um identificador UUID criptográfico seguro antes da gravação no armazenamento temporário.
              </span>
            </div>
          </div>
        </div>

        {uploadProgress !== null && (
          <div className="flex flex-col gap-1.5 mt-2 animate-fade-in">
            <div className="flex items-center justify-between text-[10px] font-semibold text-muted">
              <span>Transmitindo arquivo para o servidor...</span>
              <span>{uploadProgress}%</span>
            </div>
            <div className="w-full h-1.5 bg-gray-100 rounded-full overflow-hidden">
              <div
                className="h-full bg-primary transition-all duration-300"
                style={{ width: `${uploadProgress}%` }}
              />
            </div>
          </div>
        )}

        <Button
          type="submit"
          disabled={!selectedFile || !consentChecked || isPending}
          className="w-full py-2.5 mt-2 font-bold"
        >
          {isPending ? "Processando Análise..." : "Enviar e Analisar Exame"}
        </Button>
      </form>
    </Card>
  )
}
