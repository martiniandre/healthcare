import { useState, useRef } from "react"
import { useTranslation } from "react-i18next"
import { UploadCloud, FileText, X, CheckCircle2, AlertCircle, ShieldAlert } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import { Progress } from "@/shared/components/ui/progress"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/components/ui/card"
import { Input } from "@/shared/components/ui/input"
import { Checkbox } from "@/shared/components/ui/checkbox"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { Badge } from "@/shared/components/ui/badge"
import { Spinner } from "@/shared/components/ui/spinner"
import { Field, FieldContent, FieldDescription, FieldLabel, FieldTitle } from "@/shared/components/ui/field"
import { cn } from "@/shared/utils/cn"
import {
  EXAM_ACCEPT_ATTRIBUTE,
  EXAM_MAX_FILE_MEGABYTES,
  formatFileSize,
  validateExamFile,
} from "../exam_file_validation"
import { ExamUploadErrorCode } from "../types"

export type UploadStatus = "success" | "failed" | "cancelled"

export interface UploadOutcome {
  status: UploadStatus
  errorCode?: string
}

type UploadResult = "idle" | UploadStatus

const SERVER_ERROR_MESSAGE_KEYS = {
  [ExamUploadErrorCode.FILE_CONTENT_MISMATCH]: {
    title: "uploader.serverError.contentMismatchTitle",
    help: "uploader.serverError.contentMismatchHelp",
  },
  [ExamUploadErrorCode.FILE_CORRUPTED]: {
    title: "uploader.serverError.corruptedTitle",
    help: "uploader.serverError.corruptedHelp",
  },
  [ExamUploadErrorCode.FILE_EMPTY]: {
    title: "uploader.serverError.emptyTitle",
    help: "uploader.serverError.emptyHelp",
  },
  [ExamUploadErrorCode.DEIDENTIFICATION_REQUIRED]: {
    title: "uploader.serverError.deidentificationTitle",
    help: "uploader.serverError.deidentificationHelp",
  },
} as const

const resolveServerErrorKeys = (errorCode: string | null) => {
  if (errorCode && errorCode in SERVER_ERROR_MESSAGE_KEYS) {
    return SERVER_ERROR_MESSAGE_KEYS[errorCode as keyof typeof SERVER_ERROR_MESSAGE_KEYS]
  }
  return { title: "uploader.failedTitle", help: "uploader.failedHelp" } as const
}

interface FileUploaderProperties {
  onUpload: (
    file: File,
    consent: boolean,
    anonymize: boolean,
    acknowledgeExternalProcessing: boolean
  ) => Promise<UploadOutcome>
  onCancel: () => void
  isPending: boolean
  uploadProgress: number | null
}

const FILE_HINT_ID = "exam-file-hint"
const FILE_ERROR_ID = "exam-file-error"
const REASONS_ID = "exam-submit-reasons"

export const FileUploader = ({ onUpload, onCancel, isPending, uploadProgress }: FileUploaderProperties) => {
  const { t } = useTranslation("examAnalyzer")
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [consentChecked, setConsentChecked] = useState(false)
  const [anonymizeChecked, setAnonymizeChecked] = useState(true)
  const [externalProcessingAcknowledged, setExternalProcessingAcknowledged] = useState(false)
  const [uploadErrorCode, setUploadErrorCode] = useState<string | null>(null)
  const [fileError, setFileError] = useState<string | null>(null)
  const [hasAttemptedSubmit, setHasAttemptedSubmit] = useState(false)
  const [uploadResult, setUploadResult] = useState<UploadResult>("idle")
  const [isDragActive, setIsDragActive] = useState(false)

  const fileInputReference = useRef<HTMLInputElement>(null)

  const isUploading = isPending || uploadProgress !== null

  const resetFileInput = () => {
    if (fileInputReference.current) {
      fileInputReference.current.value = ""
    }
  }

  const validateAndSetFile = (file: File) => {
    setUploadResult("idle")
    setUploadErrorCode(null)
    const validationError = validateExamFile(file)
    if (validationError) {
      const errorKeys = {
        empty: "uploader.errorEmpty",
        type: "uploader.errorType",
        tooLarge: "uploader.errorLimit",
      } as const
      setFileError(
        t(errorKeys[validationError], {
          fileName: file.name,
          size: formatFileSize(file.size),
          maxSize: EXAM_MAX_FILE_MEGABYTES,
        })
      )
      resetFileInput()
      return
    }
    setFileError(null)
    setSelectedFile(file)
  }

  const handleDragOver = (event: React.DragEvent<HTMLElement>) => {
    event.preventDefault()
    if (!isUploading) {
      setIsDragActive(true)
    }
  }

  const handleDragLeave = (event: React.DragEvent<HTMLElement>) => {
    event.preventDefault()
    setIsDragActive(false)
  }

  const handleDrop = (event: React.DragEvent<HTMLElement>) => {
    event.preventDefault()
    setIsDragActive(false)
    if (isUploading) {
      return
    }
    const droppedFile = event.dataTransfer.files?.[0]
    if (droppedFile) {
      validateAndSetFile(droppedFile)
    }
  }

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const chosenFile = event.target.files?.[0]
    if (chosenFile) {
      validateAndSetFile(chosenFile)
    }
  }

  const handleClearFile = () => {
    setSelectedFile(null)
    setFileError(null)
    setUploadResult("idle")
    setUploadErrorCode(null)
    resetFileInput()
  }

  const handleAnonymizeChange = (checked: boolean) => {
    setAnonymizeChecked(checked)
    if (checked) {
      setExternalProcessingAcknowledged(false)
    }
  }

  const requiresExternalProcessingAcknowledgement = !anonymizeChecked

  const missingRequirements = [
    ...(selectedFile ? [] : [t("uploader.missingFile")]),
    ...(consentChecked ? [] : [t("uploader.missingConsent")]),
    ...(requiresExternalProcessingAcknowledgement && !externalProcessingAcknowledged
      ? [t("uploader.missingAcknowledgement")]
      : []),
  ]

  const handleFormSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    if (isUploading) {
      return
    }
    setHasAttemptedSubmit(true)
    if (!selectedFile || missingRequirements.length > 0) {
      return
    }
    setUploadResult("idle")
    setUploadErrorCode(null)
    const outcome = await onUpload(
      selectedFile,
      consentChecked,
      anonymizeChecked,
      requiresExternalProcessingAcknowledgement && externalProcessingAcknowledged
    )
    setUploadResult(outcome.status)
    setUploadErrorCode(outcome.errorCode ?? null)
    if (outcome.status === "success") {
      setSelectedFile(null)
      setConsentChecked(false)
      setAnonymizeChecked(true)
      setExternalProcessingAcknowledged(false)
      setHasAttemptedSubmit(false)
      resetFileInput()
    }
  }

  const submitLabel = (() => {
    if (isUploading) return t("uploader.processing")
    if (uploadResult === "failed" || uploadResult === "cancelled") return t("uploader.retry")
    return t("uploader.submit")
  })()

  const liveStatus = (() => {
    if (isUploading) return t("uploader.uploading")
    if (uploadResult === "success") return t("uploader.successTitle")
    if (uploadResult === "failed") return t("uploader.failedTitle")
    if (uploadResult === "cancelled") return t("uploader.cancelledTitle")
    return ""
  })()

  const failureKeys = resolveServerErrorKeys(uploadErrorCode)

  const fileDescribedBy = fileError ? `${FILE_HINT_ID} ${FILE_ERROR_ID}` : FILE_HINT_ID

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base font-bold">{t("uploader.title")}</CardTitle>
        <CardDescription className="text-xs leading-relaxed">{t("uploader.subtitle")}</CardDescription>
      </CardHeader>

      <CardContent>
        <form onSubmit={handleFormSubmit} noValidate className="flex flex-col gap-6">
          <p role="status" aria-live="polite" className="sr-only">
            {liveStatus}
          </p>

          {uploadResult === "success" && (
            <Alert variant="success" role="status">
              <CheckCircle2 aria-hidden="true" />
              <AlertDescription className="flex flex-col gap-0.5 text-success">
                <span className="font-semibold">{t("uploader.successTitle")}</span>
                <span className="text-xs">{t("uploader.successNextStep")}</span>
              </AlertDescription>
            </Alert>
          )}

          {uploadResult === "failed" && (
            <Alert variant="destructive">
              <AlertCircle aria-hidden="true" />
              <AlertDescription className="flex flex-col gap-0.5">
                <span className="font-semibold">{t(failureKeys.title)}</span>
                <span className="text-xs">{t(failureKeys.help)}</span>
              </AlertDescription>
            </Alert>
          )}

          {uploadResult === "cancelled" && (
            <Alert variant="warning" role="status">
              <AlertCircle aria-hidden="true" />
              <AlertDescription className="flex flex-col gap-0.5 text-warning">
                <span className="font-semibold">{t("uploader.cancelledTitle")}</span>
                <span className="text-xs">{t("uploader.cancelledHelp")}</span>
              </AlertDescription>
            </Alert>
          )}

          <section aria-labelledby="exam-step-file" className="flex flex-col gap-3">
            <h3 id="exam-step-file" className="text-sm font-bold text-foreground">
              {t("uploader.stepFile")}
            </h3>

            <div
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              className={cn(
                "flex min-h-40 flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed p-6 text-center transition-all duration-300 sm:p-8",
                isDragActive ? "scale-[1.01] border-primary bg-primary/5" : "border-border bg-muted/30",
                fileError && "border-destructive/60"
              )}
            >
              <div className="flex size-12 items-center justify-center rounded-full bg-primary/10">
                <UploadCloud aria-hidden="true" className="size-6 text-primary" />
              </div>

              <p className="text-sm font-semibold text-foreground">
                {isDragActive ? t("uploader.dropActive") : t("uploader.dropTitle")}
              </p>

              <Input
                id="exam-file-input"
                data-testid="exam-file-input"
                type="file"
                ref={fileInputReference}
                onChange={handleFileSelect}
                className="hidden"
                accept={EXAM_ACCEPT_ATTRIBUTE}
                tabIndex={-1}
                aria-hidden="true"
              />

              <Button
                type="button"
                variant="outline"
                disabled={isUploading}
                onClick={() => fileInputReference.current?.click()}
                aria-describedby={fileDescribedBy}
                className="max-sm:min-h-11 max-sm:w-full"
              >
                <UploadCloud aria-hidden="true" />
                {selectedFile ? t("uploader.replaceFile") : t("uploader.chooseFile")}
              </Button>

              <p id={FILE_HINT_ID} className="text-[11px] leading-relaxed text-muted-foreground">
                {t("uploader.fileGuidelines", { maxSize: EXAM_MAX_FILE_MEGABYTES })}
              </p>
            </div>

            {fileError && (
              <Alert id={FILE_ERROR_ID} variant="destructive">
                <AlertCircle aria-hidden="true" />
                <AlertDescription className="text-xs font-semibold text-destructive">{fileError}</AlertDescription>
              </Alert>
            )}

            {selectedFile && (
              <div className="flex animate-fade-in items-center justify-between gap-3 rounded-lg border border-border bg-muted p-3">
                <div className="flex min-w-0 items-center gap-3">
                  <FileText aria-hidden="true" className="size-5 shrink-0 text-primary" />
                  <div className="min-w-0">
                    <span className="block text-[11px] font-medium text-muted-foreground">
                      {t("uploader.selectedFile")}
                    </span>
                    <span className="block truncate text-xs font-semibold text-foreground">{selectedFile.name}</span>
                    <span className="mt-0.5 block text-[11px] tabular-nums text-muted-foreground">
                      {formatFileSize(selectedFile.size)}
                    </span>
                  </div>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  disabled={isUploading}
                  onClick={handleClearFile}
                  aria-label={t("uploader.clearFile")}
                  className="text-muted-foreground hover:text-destructive max-sm:size-11"
                >
                  <X aria-hidden="true" />
                </Button>
              </div>
            )}
          </section>

          <section aria-labelledby="exam-step-options" className="flex flex-col gap-3">
            <h3 id="exam-step-options" className="text-sm font-bold text-foreground">
              {t("uploader.stepOptions")}
            </h3>

            <FieldLabel htmlFor="exam-consent" className="cursor-pointer select-none">
              <Field orientation="horizontal" className="items-start p-3!">
                <Checkbox
                  id="exam-consent"
                  checked={consentChecked}
                  disabled={isUploading}
                  aria-required="true"
                  aria-invalid={hasAttemptedSubmit && !consentChecked}
                  aria-describedby="exam-consent-description"
                  onCheckedChange={(checked) => setConsentChecked(checked === true)}
                />
                <FieldContent className="text-left">
                  <FieldTitle className="flex flex-wrap items-center gap-2 text-xs font-semibold">
                    {t("uploader.consentTitle")}
                    <Badge variant="destructive" className="text-[10px] uppercase">
                      {t("uploader.required")}
                    </Badge>
                  </FieldTitle>
                  <FieldDescription id="exam-consent-description" className="text-[11px] leading-relaxed">
                    {t("uploader.consentDesc")} {t("uploader.consentIfUnchecked")}
                  </FieldDescription>
                </FieldContent>
              </Field>
            </FieldLabel>

            <FieldLabel htmlFor="exam-anonymize" className="cursor-pointer select-none">
              <Field orientation="horizontal" className="items-start p-3!">
                <Checkbox
                  id="exam-anonymize"
                  checked={anonymizeChecked}
                  disabled={isUploading}
                  aria-describedby="exam-anonymize-description"
                  onCheckedChange={(checked) => handleAnonymizeChange(checked === true)}
                />
                <FieldContent className="text-left">
                  <FieldTitle className="flex flex-wrap items-center gap-2 text-xs font-semibold">
                    {t("uploader.anonymizeTitle")}
                    <Badge variant="secondary" className="text-[10px] uppercase">
                      {t("uploader.optional")}
                    </Badge>
                  </FieldTitle>
                  <FieldDescription id="exam-anonymize-description" className="text-[11px] leading-relaxed">
                    {t("uploader.anonymizeDesc")}
                  </FieldDescription>
                  <p className="text-[11px] font-semibold leading-relaxed text-foreground">
                    {t("uploader.anonymizeLimitation")}
                  </p>
                </FieldContent>
              </Field>
            </FieldLabel>

            {requiresExternalProcessingAcknowledgement && (
              <div className="flex animate-fade-in flex-col gap-3">
                <Alert variant="warning">
                  <ShieldAlert aria-hidden="true" />
                  <AlertTitle className="line-clamp-none text-xs font-bold">
                    {t("uploader.externalProcessingTitle")}
                  </AlertTitle>
                  <AlertDescription className="text-[11px] leading-relaxed text-warning">
                    {t("uploader.externalProcessingWarning")}
                  </AlertDescription>
                </Alert>

                <FieldLabel htmlFor="exam-acknowledge-external" className="cursor-pointer select-none">
                  <Field orientation="horizontal" className="items-start p-3!">
                    <Checkbox
                      id="exam-acknowledge-external"
                      checked={externalProcessingAcknowledged}
                      disabled={isUploading}
                      aria-required="true"
                      aria-invalid={hasAttemptedSubmit && !externalProcessingAcknowledged}
                      onCheckedChange={(checked) => setExternalProcessingAcknowledged(checked === true)}
                    />
                    <FieldContent className="text-left">
                      <FieldTitle className="flex flex-wrap items-center gap-2 text-xs font-semibold">
                        {t("uploader.acknowledgementTitle")}
                        <Badge variant="destructive" className="text-[10px] uppercase">
                          {t("uploader.required")}
                        </Badge>
                      </FieldTitle>
                    </FieldContent>
                  </Field>
                </FieldLabel>
              </div>
            )}
          </section>

          <section
            aria-labelledby="exam-step-submit"
            className="flex flex-col gap-3 max-sm:sticky max-sm:bottom-0 max-sm:-mx-6 max-sm:border-t max-sm:border-border max-sm:bg-card max-sm:px-6 max-sm:py-3"
          >
            <h3 id="exam-step-submit" className="text-sm font-bold text-foreground">
              {t("uploader.stepSubmit")}
            </h3>

            {isUploading && (
              <div className="flex animate-fade-in flex-col gap-1.5">
                <div className="flex items-center justify-between text-[11px] font-semibold text-muted-foreground">
                  <span>{t("uploader.uploading")}</span>
                  <span className="tabular-nums">{uploadProgress ?? 0}%</span>
                </div>
                <Progress
                  value={uploadProgress ?? 0}
                  aria-label={t("uploader.uploading")}
                  aria-valuetext={`${uploadProgress ?? 0}%`}
                  className="h-1.5"
                />
              </div>
            )}

            {!isUploading && missingRequirements.length > 0 && (
              <div
                id={REASONS_ID}
                aria-live="polite"
                className={cn(
                  "rounded-lg border p-3 text-xs",
                  hasAttemptedSubmit
                    ? "border-destructive/40 bg-destructive/5 text-destructive"
                    : "border-border bg-muted/40 text-muted-foreground"
                )}
              >
                <p className="font-semibold">
                  {hasAttemptedSubmit ? t("uploader.cannotSubmit") : t("uploader.stillNeeded")}
                </p>
                <ul className="mt-1 list-disc pl-5">
                  {missingRequirements.map((requirement) => (
                    <li key={requirement}>{requirement}</li>
                  ))}
                </ul>
              </div>
            )}

            <div className="flex flex-col gap-2 sm:flex-row">
              <Button
                type="submit"
                disabled={isUploading}
                aria-describedby={missingRequirements.length > 0 ? REASONS_ID : undefined}
                className="w-full font-bold max-sm:min-h-11 sm:flex-1"
              >
                {isUploading && <Spinner />}
                {submitLabel}
              </Button>
              {isUploading && (
                <Button type="button" variant="outline" onClick={onCancel} className="max-sm:min-h-11">
                  {t("uploader.cancel")}
                </Button>
              )}
            </div>
          </section>
        </form>
      </CardContent>
    </Card>
  )
}
