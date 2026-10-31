import { useState, useRef } from "react"
import { useTranslation } from "react-i18next"
import { UploadCloud, FileText, X } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import { Progress } from "@/shared/components/ui/progress"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/components/ui/card"
import { Input } from "@/shared/components/ui/input"
import { Label } from "@/shared/components/ui/label"
import { Checkbox } from "@/shared/components/ui/checkbox"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { Spinner } from "@/shared/components/ui/spinner"
import { Field, FieldContent, FieldDescription, FieldLabel, FieldTitle } from "@/shared/components/ui/field"
import { cn } from "@/shared/utils/cn"

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
    <Card>
      <CardHeader>
        <CardTitle className="text-base font-bold">{t("uploader.title")}</CardTitle>
        <CardDescription className="text-xs leading-relaxed">{t("uploader.subtitle")}</CardDescription>
      </CardHeader>

      <CardContent>
        <form onSubmit={handleFormSubmit} className="flex flex-col gap-5">
          <Label
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
            className={cn(
              "flex cursor-pointer flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed p-8 transition-all duration-300",
              isDragActive
                ? "scale-[1.01] border-primary bg-primary/5"
                : "border-border hover:border-primary/50 hover:bg-muted/50"
            )}
          >
            <Input
              type="file"
              ref={fileInputReference}
              onChange={handleFileSelect}
              className="hidden"
              accept="image/*,.pdf"
            />

            <div className="flex size-12 items-center justify-center rounded-full bg-primary/10">
              <UploadCloud className="size-6 text-primary" />
            </div>

            <div className="text-center">
              <span className="block text-sm font-semibold text-foreground">
                {t("uploader.selectFile")}
              </span>
              <span className="mt-1 block text-[11px] font-normal text-muted-foreground">
                {t("uploader.fileGuidelines")}
              </span>
            </div>
          </Label>

          {uploaderState.error && (
            <Alert variant="destructive" className="justify-items-center text-center">
              <AlertDescription className="justify-items-center text-xs font-semibold text-destructive">
                {uploaderState.error}
              </AlertDescription>
            </Alert>
          )}

          {uploaderState.file && (
            <div className="flex animate-fade-in items-center justify-between rounded-lg border border-border bg-muted p-3.5">
              <div className="flex min-w-0 items-center gap-3">
                <FileText className="size-5 shrink-0 text-primary" />
                <div className="min-w-0">
                  <span className="block truncate text-xs font-semibold text-foreground">
                    {uploaderState.file.name}
                  </span>
                  <span className="mt-0.5 block text-[11px] text-muted-foreground">
                    {(uploaderState.file.size / (1024 * 1024)).toFixed(2)} MB
                  </span>
                </div>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="icon-sm"
                onClick={handleClearFile}
                aria-label={t("uploader.clearFile")}
                className="text-muted-foreground hover:text-destructive"
              >
                <X />
              </Button>
            </div>
          )}

          <div className="flex flex-col gap-2.5">
            <FieldLabel htmlFor="exam-consent" className="cursor-pointer select-none">
              <Field orientation="horizontal" className="items-start p-3!">
                <Checkbox
                  id="exam-consent"
                  checked={uploaderState.consentChecked}
                  onCheckedChange={(checked) =>
                    setUploaderState((prev) => ({ ...prev, consentChecked: checked === true }))
                  }
                />
                <FieldContent className="text-left">
                  <FieldTitle className="text-xs font-semibold">{t("uploader.consentTitle")}</FieldTitle>
                  <FieldDescription className="text-[11px] leading-relaxed">
                    {t("uploader.consentDesc")}
                  </FieldDescription>
                </FieldContent>
              </Field>
            </FieldLabel>

            <FieldLabel htmlFor="exam-anonymize" className="cursor-pointer select-none">
              <Field orientation="horizontal" className="items-start p-3!">
                <Checkbox
                  id="exam-anonymize"
                  checked={uploaderState.anonymizeChecked}
                  onCheckedChange={(checked) =>
                    setUploaderState((prev) => ({ ...prev, anonymizeChecked: checked === true }))
                  }
                />
                <FieldContent className="text-left">
                  <FieldTitle className="text-xs font-semibold">{t("uploader.anonymizeTitle")}</FieldTitle>
                  <FieldDescription className="text-[11px] leading-relaxed">
                    {t("uploader.anonymizeDesc")}
                  </FieldDescription>
                </FieldContent>
              </Field>
            </FieldLabel>
          </div>

          {uploadProgress !== null && (
            <div className="mt-2 flex animate-fade-in flex-col gap-1.5">
              <div className="flex items-center justify-between text-[11px] font-semibold text-muted-foreground">
                <span>{t("uploader.uploading")}</span>
                <span className="tabular-nums">{uploadProgress}%</span>
              </div>
              <Progress value={uploadProgress} aria-label={t("uploader.uploading")} className="h-1.5" />
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
            className="mt-1 w-full font-bold"
          >
            {isPending && <Spinner />}
            {isPending ? t("uploader.processing") : t("uploader.submit")}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
