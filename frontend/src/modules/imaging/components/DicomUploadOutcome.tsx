import { useTranslation } from "react-i18next"
import { AlertTriangle, CheckCircle2 } from "lucide-react"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { resolveInstanceRejectionReasonKey } from "../utils/dicom_upload_messages"
import type { DicomUploadResult } from "../types"

interface DicomUploadOutcomeProperties {
  result: DicomUploadResult
}

export const DicomUploadOutcome = ({ result }: DicomUploadOutcomeProperties) => {
  const { t } = useTranslation("imaging")
  const rejectedInstances = result.instancesRejected ?? []

  return (
    <>
      <Alert variant="success" role="status">
        <CheckCircle2 aria-hidden="true" />
        <AlertDescription className="flex flex-col gap-0.5 text-success">
          <span className="font-semibold">
            {result.studyCreated ? t("upload.result.created") : t("upload.result.appended")}
          </span>
          {typeof result.instancesAccepted === "number" && (
            <span className="text-xs">{t("upload.result.acceptedCount", { count: result.instancesAccepted })}</span>
          )}
          <span className="text-xs">{t("upload.successNextStep")}</span>
        </AlertDescription>
      </Alert>

      {rejectedInstances.length > 0 && (
        <Alert variant="warning">
          <AlertTriangle aria-hidden="true" />
          <AlertDescription className="flex flex-col gap-2 text-warning">
            <span className="font-semibold">
              {t("upload.result.rejectedTitle", { count: rejectedInstances.length })}
            </span>
            <ul className="flex flex-col gap-1.5 text-xs">
              {rejectedInstances.map((rejectedInstance, index) => (
                <li key={`${rejectedInstance.fileName}-${index}`} className="flex flex-col">
                  <span className="break-all font-semibold">{rejectedInstance.fileName}</span>
                  <span>{t(resolveInstanceRejectionReasonKey(rejectedInstance.code))}</span>
                </li>
              ))}
            </ul>
          </AlertDescription>
        </Alert>
      )}
    </>
  )
}
