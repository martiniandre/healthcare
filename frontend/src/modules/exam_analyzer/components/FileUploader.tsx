import { useState, useRef } from "react"
import { useTranslation } from "react-i18next"
import { UploadCloud, CheckSquare, Square, FileText, X } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"
import { Card } from "../../../shared/components/ui/Card"

interface FileUploaderProperties {
  onUpload: (file: File, consent: boolean, anonymize: boolean) => void
  isPending: boolean
  uploadProgress: number | null
}

export const FileUploader = ({ onUpload, isPending, uploadProgress }: FileUploaderProperties) => {
  const { t } = useTranslation("examAnalyzer")
  const [uploaderState, setUploaderState] = useState<{
    file: File | null
    consentChecked: boolean
    anonymizeChecked: boolean
    error: string | null
  }>({
    file: null,
    consentChecked: false,
    anonymizeChecked: false,
    error: null,
  })
  const [isDragActive, setIsDragActive] = useState<boolean>(false)
  
  const fileInputReference = useRef<HTMLInputElement>(null)

  const validateAndSetFile = (file: File) => {
    setUploaderState((prev) => ({ ...prev, error: null }))
    const fifteenMegaBytes = 15 * 1024 * 1024
    if (file.size > fifteenMegaBytes) {
      setUploaderState((prev) => ({ ...prev, error: t("uploader.errorLimit") }))
      return
    }
    const allowedMimeTypes = ["image/jpeg", "image/png", "image/gif", "image/webp", "application/pdf"]
    if (!allowedMimeTypes.includes(file.type)) {
      setUploaderState((prev) => ({ ...prev, error: t("uploader.errorType") }))
      return
    }
    setUploaderState((prev) => ({ ...prev, file }))
  }

  const handleDragOver = (event: React.DragEvent<HTMLElement>) => {
    event.preventDefault()
    setIsDragActive(true)
  }

  const handleDragLeave = (event: React.DragEvent<HTMLElement>) => {
    event.preventDefault()
    setIsDragActive(false)
  }

  const handleDrop = (event: React.DragEvent<HTMLElement>) => {
    event.preventDefault()
    setIsDragActive(false)
    const file = event.dataTransfer.files?.[0]
    if (file) {
      validateAndSetFile(file)
    }
  }

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      validateAndSetFile(file)
    }
  }

  const handleClearFile = () => {
    setUploaderState((prev) => ({ ...prev, file: null, error: null }))
    if (fileInputReference.current) {
      fileInputReference.current.value = ""
    }
  }

  const handleFormSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    if (uploaderState.file && uploaderState.consentChecked) {
      onUpload(uploaderState.file, uploaderState.consentChecked, uploaderState.anonymizeChecked)
    }
  }

  return (
    <Card glowingType="cyan" className="p-6 bg-white border border-border rounded-xl">
      <h3 className="text-base font-bold text-gray-900 mb-2">
        {t("uploader.title")}
      </h3>
      <span className="text-xs text-muted block mb-5 leading-normal">
        {t("uploader.subtitle")}
      </span>

      <form onSubmit={handleFormSubmit} className="flex flex-col gap-5">
        <label
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
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
              {t("uploader.selectFile")}
            </span>
            <span className="text-[11px] text-muted block mt-1">
              {t("uploader.fileGuidelines")}
            </span>
          </div>
        </label>

        {uploaderState.error && (
          <div className="text-xs font-semibold text-red-500 bg-red-50 border border-red-200 rounded-lg p-3 text-center">
            {uploaderState.error}
          </div>
        )}

        {uploaderState.file && (
          <div className="flex items-center justify-between p-3.5 bg-gray-50 border border-border/80 rounded-lg animate-fade-in">
            <div className="flex items-center gap-3 min-w-0">
              <FileText className="w-5 h-5 text-primary shrink-0" />
              <div className="min-w-0">
                <span className="text-xs font-semibold text-gray-800 block truncate">
                  {uploaderState.file.name}
                </span>
                <span className="text-[10px] text-muted block mt-0.5">
                  {(uploaderState.file.size / (1024 * 1024)).toFixed(2)} MB
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
          <label className="flex items-start gap-3 cursor-pointer select-none group">
            <input
              type="checkbox"
              checked={uploaderState.consentChecked}
              onChange={(event) => setUploaderState((prev) => ({ ...prev, consentChecked: event.target.checked }))}
              className="sr-only"
            />
            <div className="mt-0.5 text-primary">
              {uploaderState.consentChecked ? (
                <CheckSquare className="w-4.5 h-4.5 group-hover:scale-105 transition-transform" />
              ) : (
                <Square className="w-4.5 h-4.5 text-gray-400 group-hover:scale-105 transition-transform" />
              )}
            </div>
            <div className="flex-1 text-left">
              <span className="text-xs font-semibold text-gray-700 block">
                {t("uploader.consentTitle")}
              </span>
              <span className="text-[10px] text-muted block mt-0.5 leading-normal">
                {t("uploader.consentDesc")}
              </span>
            </div>
          </label>

          <label className="flex items-start gap-3 cursor-pointer select-none group">
            <input
              type="checkbox"
              checked={uploaderState.anonymizeChecked}
              onChange={(event) => setUploaderState((prev) => ({ ...prev, anonymizeChecked: event.target.checked }))}
              className="sr-only"
            />
            <div className="mt-0.5 text-secondary">
              {uploaderState.anonymizeChecked ? (
                <CheckSquare className="w-4.5 h-4.5 group-hover:scale-105 transition-transform" />
              ) : (
                <Square className="w-4.5 h-4.5 text-gray-400 group-hover:scale-105 transition-transform" />
              )}
            </div>
            <div className="flex-1 text-left">
              <span className="text-xs font-semibold text-gray-700 block">
                {t("uploader.anonymizeTitle")}
              </span>
              <span className="text-[10px] text-muted block mt-0.5 leading-normal">
                {t("uploader.anonymizeDesc")}
              </span>
            </div>
          </label>
        </div>

        {uploadProgress !== null && (
          <div className="flex flex-col gap-1.5 mt-2 animate-fade-in">
            <div className="flex items-center justify-between text-[10px] font-semibold text-muted">
              <span>{t("uploader.uploading")}</span>
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
          disabled={!uploaderState.file || !uploaderState.consentChecked || isPending}
          className="w-full py-2.5 mt-2 font-bold"
        >
          {isPending ? t("uploader.processing") : t("uploader.submit")}
        </Button>
      </form>
    </Card>
  )
}
