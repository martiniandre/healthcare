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

  const isSubmitDisabled = !uploaderState.file || !uploaderState.consentChecked || isPending

  const submitHint = (() => {
    if (isPending || uploadProgress !== null) return null
    if (!uploaderState.file) return t("uploader.hintNeedsFile")
    if (!uploaderState.consentChecked) return t("uploader.hintNeedsConsent")
    return null
  })()

  return (
    <Card glowingType="cyan" className="bg-card">
      <h3 className="text-base font-bold text-gray-900">{t("uploader.title")}</h3>
      <span className="mt-1.5 mb-4 block text-xs leading-relaxed text-muted">
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
          <div
            role="alert"
            className="rounded-lg border border-danger/25 bg-danger/8 p-3 text-center text-xs font-semibold text-danger"
          >
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

        <div className="flex flex-col gap-2.5">
          <label className="group flex cursor-pointer select-none items-start gap-3 rounded-lg border border-border bg-gray-50/60 p-3 transition-colors hover:border-primary/30">
            <input
              type="checkbox"
              checked={uploaderState.consentChecked}
              onChange={(event) => setUploaderState((prev) => ({ ...prev, consentChecked: event.target.checked }))}
              className="sr-only"
            />
            <span className="mt-0.5 shrink-0 text-primary">
              {uploaderState.consentChecked ? (
                <CheckSquare className="h-4 w-4 transition-transform group-hover:scale-105" aria-hidden="true" />
              ) : (
                <Square className="h-4 w-4 text-gray-400 transition-transform group-hover:scale-105" aria-hidden="true" />
              )}
            </span>
            <span className="flex-1 text-left">
              <span className="block text-xs font-semibold text-gray-900">
                {t("uploader.consentTitle")}
              </span>
              <span className="mt-0.5 block text-[11px] leading-relaxed text-muted-foreground">
                {t("uploader.consentDesc")}
              </span>
            </span>
          </label>

          <label className="group flex cursor-pointer select-none items-start gap-3 rounded-lg border border-border bg-gray-50/60 p-3 transition-colors hover:border-secondary/30">
            <input
              type="checkbox"
              checked={uploaderState.anonymizeChecked}
              onChange={(event) => setUploaderState((prev) => ({ ...prev, anonymizeChecked: event.target.checked }))}
              className="sr-only"
            />
            <span className="mt-0.5 shrink-0 text-secondary">
              {uploaderState.anonymizeChecked ? (
                <CheckSquare className="h-4 w-4 transition-transform group-hover:scale-105" aria-hidden="true" />
              ) : (
                <Square className="h-4 w-4 text-gray-400 transition-transform group-hover:scale-105" aria-hidden="true" />
              )}
            </span>
            <span className="flex-1 text-left">
              <span className="block text-xs font-semibold text-gray-900">
                {t("uploader.anonymizeTitle")}
              </span>
              <span className="mt-0.5 block text-[11px] leading-relaxed text-muted-foreground">
                {t("uploader.anonymizeDesc")}
              </span>
            </span>
          </label>
        </div>

        {uploadProgress !== null && (
          <div className="mt-2 flex flex-col gap-1.5 animate-fade-in">
            <div className="flex items-center justify-between text-[11px] font-semibold text-muted-foreground">
              <span>{t("uploader.uploading")}</span>
              <span className="tabular-nums">{uploadProgress}%</span>
            </div>
            <div
              role="progressbar"
              aria-label={t("uploader.uploading")}
              aria-valuenow={uploadProgress}
              aria-valuemin={0}
              aria-valuemax={100}
              className="h-1.5 w-full overflow-hidden rounded-full bg-black/10"
            >
              <div
                className="h-full rounded-full bg-primary transition-all duration-300"
                style={{ width: `${uploadProgress}%` }}
              />
            </div>
          </div>
        )}

        {submitHint && (
          <p className="text-center text-[11px] leading-normal text-muted-foreground">
            {submitHint}
          </p>
        )}

        <Button
          type="submit"
          disabled={isSubmitDisabled}
          className="mt-1 w-full py-2.5 font-bold"
        >
          {isPending ? t("uploader.processing") : t("uploader.submit")}
        </Button>
      </form>
    </Card>
  )
}
