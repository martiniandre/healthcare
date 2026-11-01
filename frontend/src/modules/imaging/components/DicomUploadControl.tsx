import { useEffect, useRef, useState } from "react"
import { useTranslation } from "react-i18next"
import { AlertCircle, AlertTriangle, CheckCircle2, UploadCloud } from "lucide-react"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { Button } from "@/shared/components/ui/button"
import { Input } from "@/shared/components/ui/input"
import { cn } from "@/shared/utils/cn"
import { toast } from "../../../shared/store/toast_store"
import { extractApiErrorDetails } from "../../../shared/utils/api_error"
import { useUploadImagingStudyMutation } from "../queries"
import {
  DICOM_ACCEPT_ATTRIBUTE,
  DICOM_MAX_FILE_GIGABYTES,
  formatDicomFileSize,
  validateDicomFile,
} from "../utils/dicom_file_validation"
import { ImagingUploadProgress } from "./ImagingUploadProgress"
import { DicomUploadSummary } from "./DicomUploadSummary"
import { DicomUploadErrorCode, DicomUploadWarningCode, type DicomUploadResult } from "../types"

interface DicomUploadControlProperties {
  patientFhirId: string
}

type UploadFeedback =
  | { kind: "idle" }
  | { kind: "validation"; message: string }
  | { kind: "success"; result: DicomUploadResult }
  | { kind: "failed"; errorCode?: string; reason?: string }
  | { kind: "cancelled" }

const WARNING_MESSAGE_KEYS = {
  [DicomUploadWarningCode.PATIENT_IDENTITY_MISSING]: "upload.warnings.patient_identity_missing",
  [DicomUploadWarningCode.ANONYMIZED_STUDY]: "upload.warnings.anonymized_study",
  [DicomUploadWarningCode.DUPLICATE_INSTANCE]: "upload.warnings.duplicate_instance",
  [DicomUploadWarningCode.MODALITY_OVERRIDDEN]: "upload.warnings.modality_overridden",
} as const

const REJECTION_MESSAGE_KEYS = {
  [DicomUploadErrorCode.PATIENT_MISMATCH]: {
    title: "upload.errors.patientMismatchTitle",
    help: "upload.errors.patientMismatchHelp",
  },
  [DicomUploadErrorCode.INVALID_DICOM]: {
    title: "upload.errors.invalidDicomTitle",
    help: "upload.errors.invalidDicomHelp",
  },
  [DicomUploadErrorCode.DUPLICATE_INSTANCE]: {
    title: "upload.errors.duplicateTitle",
    help: "upload.errors.duplicateHelp",
  },
} as const

const resolveWarningMessageKey = (warningCode: string): string => {
  return warningCode in WARNING_MESSAGE_KEYS
    ? WARNING_MESSAGE_KEYS[warningCode as keyof typeof WARNING_MESSAGE_KEYS]
    : "upload.warnings.unknown"
}

const resolveRejectionMessageKeys = (errorCode: string | undefined) => {
  return errorCode && errorCode in REJECTION_MESSAGE_KEYS
    ? REJECTION_MESSAGE_KEYS[errorCode as keyof typeof REJECTION_MESSAGE_KEYS]
    : null
}

const HINT_ID = "dicom-upload-hint"
const FEEDBACK_ID = "dicom-upload-feedback"

export const DicomUploadControl = ({ patientFhirId }: DicomUploadControlProperties) => {
  const { t } = useTranslation("imaging")
  const [uploadPercentage, setUploadPercentage] = useState<number | null>(null)
  const [lastAttemptedFile, setLastAttemptedFile] = useState<File | null>(null)
  const [feedback, setFeedback] = useState<UploadFeedback>({ kind: "idle" })
  const [isDragActive, setIsDragActive] = useState(false)
  const fileInputReference = useRef<HTMLInputElement>(null)
  const abortControllerReference = useRef<AbortController | null>(null)
  const uploadImagingStudyMutation = useUploadImagingStudyMutation()

  const isUploading = uploadPercentage !== null || uploadImagingStudyMutation.isPending

  useEffect(() => {
    if (!isUploading) {
      return
    }
    const warnBeforeLeaving = (event: BeforeUnloadEvent) => {
      event.preventDefault()
      event.returnValue = ""
    }
    window.addEventListener("beforeunload", warnBeforeLeaving)
    return () => window.removeEventListener("beforeunload", warnBeforeLeaving)
  }, [isUploading])

  const resetFileInput = () => {
    if (fileInputReference.current) {
      fileInputReference.current.value = ""
    }
  }

  const uploadFile = async (file: File) => {
    const abortController = new AbortController()
    abortControllerReference.current = abortController
    setLastAttemptedFile(file)
    setFeedback({ kind: "idle" })
    setUploadPercentage(0)
    try {
      const uploadResult = await uploadImagingStudyMutation.mutateAsync({
        patientFhirId,
        dicomBlob: file,
        signal: abortController.signal,
        onUploadProgress: setUploadPercentage,
      })
      setLastAttemptedFile(null)
      setFeedback({ kind: "success", result: uploadResult })
      toast.success(t("toast.uploadSuccess"))
    } catch (uploadError) {
      if (abortController.signal.aborted) {
        setFeedback({ kind: "cancelled" })
      } else {
        const { code, message, reason } = extractApiErrorDetails(uploadError)
        setFeedback({ kind: "failed", errorCode: code, reason: reason ?? message })
        toast.error(t("toast.uploadError"))
      }
    } finally {
      abortControllerReference.current = null
      setUploadPercentage(null)
      resetFileInput()
    }
  }

  const handleSelectedFile = (file: File) => {
    const validationError = validateDicomFile(file)
    if (validationError) {
      const errorKeys = {
        empty: "upload.errorEmpty",
        type: "upload.errorType",
        tooLarge: "upload.errorLimit",
      } as const
      setFeedback({
        kind: "validation",
        message: t(errorKeys[validationError], {
          fileName: file.name,
          size: formatDicomFileSize(file.size),
          maxSize: DICOM_MAX_FILE_GIGABYTES,
        }),
      })
      resetFileInput()
      return
    }
    void uploadFile(file)
  }

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const chosenFile = event.target.files?.[0]
    if (chosenFile) {
      handleSelectedFile(chosenFile)
    }
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
      handleSelectedFile(droppedFile)
    }
  }

  const handleCancel = () => {
    abortControllerReference.current?.abort()
  }

  const handleRetry = () => {
    if (lastAttemptedFile) {
      void uploadFile(lastAttemptedFile)
    }
  }

  const progressStatus =
    uploadPercentage !== null && uploadPercentage >= 100 ? t("uploadStatus.grpcCompleted") : t("uploadStatus.initial")

  const rejectionKeys = feedback.kind === "failed" ? resolveRejectionMessageKeys(feedback.errorCode) : null
  const invalidDicomReason =
    feedback.kind === "failed" && feedback.errorCode === DicomUploadErrorCode.INVALID_DICOM ? feedback.reason : undefined

  const describedBy = feedback.kind === "validation" ? `${HINT_ID} ${FEEDBACK_ID}` : HINT_ID

  return (
    <section aria-label={t("upload.sectionLabel")} className="flex flex-col gap-3">
      {uploadPercentage !== null ? (
        <ImagingUploadProgress
          percentage={uploadPercentage}
          status={progressStatus}
          fileName={lastAttemptedFile?.name}
          fileSize={lastAttemptedFile ? formatDicomFileSize(lastAttemptedFile.size) : undefined}
          onCancel={handleCancel}
        />
      ) : (
        <div
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          className={cn(
            "flex flex-col items-center gap-3 rounded-xl border-2 border-dashed p-4 text-center transition-all duration-300 sm:flex-row sm:justify-between sm:text-left",
            isDragActive ? "border-primary bg-primary/5" : "border-border bg-muted/30"
          )}
        >
          <div className="flex min-w-0 flex-col gap-0.5">
            <p className="text-sm font-semibold text-foreground">
              {isDragActive ? t("upload.dropActive") : t("upload.dropTitle")}
            </p>
            <p id={HINT_ID} className="text-xs text-muted-foreground">
              {t("upload.guidelines", { maxSize: DICOM_MAX_FILE_GIGABYTES })}
            </p>
          </div>

          <Input
            id="dicom-file-input"
            data-testid="dicom-file-input"
            type="file"
            accept={DICOM_ACCEPT_ATTRIBUTE}
            ref={fileInputReference}
            onChange={handleFileChange}
            className="hidden"
            tabIndex={-1}
            aria-hidden="true"
          />
          <Button
            type="button"
            variant="outline"
            onClick={() => fileInputReference.current?.click()}
            disabled={isUploading}
            aria-describedby={describedBy}
            className="border-primary/20 text-primary hover:bg-primary/5 max-sm:min-h-11 max-sm:w-full"
          >
            <UploadCloud aria-hidden="true" className="size-4" />
            {t("uploadDcm")}
          </Button>
        </div>
      )}

      {feedback.kind === "validation" && (
        <Alert id={FEEDBACK_ID} variant="destructive">
          <AlertCircle aria-hidden="true" />
          <AlertDescription className="text-xs font-semibold text-destructive">{feedback.message}</AlertDescription>
        </Alert>
      )}

      {feedback.kind === "success" && (
        <>
          <Alert variant="success" role="status">
            <CheckCircle2 aria-hidden="true" />
            <AlertDescription className="flex flex-col gap-0.5 text-success">
              <span className="font-semibold">{t("upload.successTitle")}</span>
              <span className="text-xs">{t("upload.successNextStep")}</span>
            </AlertDescription>
          </Alert>

          <DicomUploadSummary result={feedback.result} />

          {(feedback.result.warnings ?? []).map((warningItem, index) => (
            <Alert key={`${warningItem.code}-${index}`} variant="warning">
              <AlertTriangle aria-hidden="true" />
              <AlertDescription className="text-xs font-semibold text-warning">
                {t(resolveWarningMessageKey(warningItem.code))}
              </AlertDescription>
            </Alert>
          ))}
        </>
      )}

      {feedback.kind === "cancelled" && (
        <Alert variant="warning" role="status">
          <AlertCircle aria-hidden="true" />
          <AlertDescription className="flex flex-wrap items-center gap-3 text-warning">
            <span className="flex flex-col gap-0.5">
              <span className="font-semibold">{t("upload.cancelledTitle")}</span>
              <span className="text-xs">{t("upload.cancelledHelp")}</span>
            </span>
            {lastAttemptedFile && (
              <Button type="button" variant="outline" size="sm" onClick={handleRetry} className="max-sm:min-h-11">
                {t("upload.retry")}
              </Button>
            )}
          </AlertDescription>
        </Alert>
      )}

      {feedback.kind === "failed" && (
        <Alert variant="destructive">
          <AlertCircle aria-hidden="true" />
          <AlertDescription className="flex flex-wrap items-center gap-3">
            <span className="flex flex-col gap-0.5">
              <span className="font-semibold">
                {rejectionKeys ? t(rejectionKeys.title) : t("upload.failedTitle")}
              </span>
              <span className="text-xs">
                {rejectionKeys
                  ? t(rejectionKeys.help)
                  : t("upload.failedHelp", { fileName: lastAttemptedFile?.name ?? "" })}
              </span>
              {invalidDicomReason && (
                <span className="text-xs">{t("upload.errors.invalidDicomReason", { reason: invalidDicomReason })}</span>
              )}
            </span>
            {lastAttemptedFile && !rejectionKeys && (
              <Button type="button" variant="outline" size="sm" onClick={handleRetry} className="max-sm:min-h-11">
                {t("upload.retry")}
              </Button>
            )}
          </AlertDescription>
        </Alert>
      )}
    </section>
  )
}
