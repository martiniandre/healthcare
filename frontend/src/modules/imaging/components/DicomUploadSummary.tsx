import { useTranslation } from "react-i18next"
import { Copy } from "lucide-react"
import { Button } from "@/shared/components/ui/button"
import { toast } from "../../../shared/store/toast_store"
import { formatDate } from "../../../shared/utils/dates"
import { normalizeDicomStudyDate, shortenDicomUid } from "../utils/pacs_helpers"
import type { DicomUploadResult } from "../types"

interface DicomUploadSummaryProperties {
  result: DicomUploadResult
}

interface SummaryEntry {
  key: string
  label: string
  value: string
  title?: string
  copyable?: boolean
}

export const DicomUploadSummary = ({ result }: DicomUploadSummaryProperties) => {
  const { t } = useTranslation("imaging")

  const studyUid = result.studyInstanceUid || result.study_instance_uid || ""

  const entries: SummaryEntry[] = [
    { key: "modality", label: t("upload.summary.modality"), value: result.modality || "" },
    {
      key: "studyDate",
      label: t("upload.summary.studyDate"),
      value: result.studyDate ? formatDate(normalizeDicomStudyDate(result.studyDate)) : "",
    },
    { key: "accessionNumber", label: t("upload.summary.accessionNumber"), value: result.accessionNumber || "" },
    {
      key: "studyUid",
      label: t("upload.summary.studyUid"),
      value: shortenDicomUid(studyUid),
      title: studyUid,
      copyable: true,
    },
    {
      key: "seriesUid",
      label: t("upload.summary.seriesUid"),
      value: shortenDicomUid(result.seriesInstanceUid || ""),
      title: result.seriesInstanceUid,
    },
    {
      key: "instanceUid",
      label: t("upload.summary.instanceUid"),
      value: shortenDicomUid(result.sopInstanceUid || ""),
      title: result.sopInstanceUid,
    },
  ].filter((entry) => entry.value !== "")

  if (entries.length === 0) {
    return null
  }

  const handleCopyStudyUid = async () => {
    try {
      await navigator.clipboard.writeText(studyUid)
      toast.success(t("upload.summary.copied"))
    } catch {
      toast.error(t("upload.summary.copyError"))
    }
  }

  return (
    <section
      aria-label={t("upload.summary.title")}
      data-testid="dicom-upload-summary"
      className="rounded-lg border border-border bg-muted/40 p-3"
    >
      <h4 className="text-[11px] font-bold uppercase tracking-wider text-muted-foreground">
        {t("upload.summary.title")}
      </h4>
      <dl className="mt-2 grid grid-cols-1 gap-x-4 gap-y-2 sm:grid-cols-2">
        {entries.map((entry) => (
          <div key={entry.key} className="min-w-0">
            <dt className="text-[11px] font-medium text-muted-foreground">{entry.label}</dt>
            <dd className="flex items-center gap-1 text-xs font-semibold text-foreground">
              <span className="min-w-0 truncate font-mono tabular-nums" title={entry.title}>
                {entry.value}
              </span>
              {entry.copyable && (
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={handleCopyStudyUid}
                  aria-label={t("upload.summary.copyUid")}
                  title={t("upload.summary.copyUid")}
                  className="size-8 shrink-0 max-sm:size-11"
                >
                  <Copy aria-hidden="true" />
                </Button>
              )}
            </dd>
          </div>
        ))}
      </dl>
    </section>
  )
}
