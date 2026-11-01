import { useId, useState } from "react"
import { useTranslation } from "react-i18next"
import { AlertCircle, ChevronDown, ChevronUp, Download } from "lucide-react"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
import { Button } from "@/shared/components/ui/button"
import { Spinner } from "@/shared/components/ui/spinner"
import { useImagingStudyInstancesQuery } from "../queries"
import { shortenDicomUid } from "../utils/pacs_helpers"

interface ImagingStudyFilesProperties {
  studyId: string
}

export const ImagingStudyFiles = ({ studyId }: ImagingStudyFilesProperties) => {
  const { t } = useTranslation("imaging")
  const [isOpen, setIsOpen] = useState(false)
  const panelId = useId()

  return (
    <div className="flex flex-col gap-2">
      <Button
        type="button"
        variant="outline"
        size="sm"
        aria-expanded={isOpen}
        aria-controls={panelId}
        onClick={() => setIsOpen((previousValue) => !previousValue)}
        className="w-full justify-between max-sm:min-h-11"
      >
        {isOpen ? t("details.files.hide") : t("details.files.show")}
        {isOpen ? <ChevronUp aria-hidden="true" /> : <ChevronDown aria-hidden="true" />}
      </Button>
      {isOpen && (
        <div id={panelId}>
          <StudyFilesList studyId={studyId} />
        </div>
      )}
    </div>
  )
}

const StudyFilesList = ({ studyId }: ImagingStudyFilesProperties) => {
  const { t } = useTranslation("imaging")
  const { data: instances, isPending, isError, refetch } = useImagingStudyInstancesQuery(studyId)

  if (isPending) {
    return (
      <div className="flex items-center gap-2 py-2 text-xs text-muted-foreground">
        <Spinner />
        <span>{t("details.files.loading")}</span>
      </div>
    )
  }

  if (isError) {
    return (
      <Alert variant="destructive">
        <AlertCircle aria-hidden="true" />
        <AlertDescription className="flex flex-wrap items-center gap-3">
          <span className="text-xs">{t("details.files.loadError")}</span>
          <Button type="button" variant="outline" size="sm" onClick={() => refetch()} className="max-sm:min-h-11">
            {t("retry")}
          </Button>
        </AlertDescription>
      </Alert>
    )
  }

  if (instances.length === 0) {
    return <p className="text-xs text-muted-foreground">{t("details.files.empty")}</p>
  }

  return (
    <div className="flex flex-col gap-2">
      <p className="text-xs font-bold text-foreground">{t("details.files.count", { count: instances.length })}</p>
      <ul className="flex flex-col gap-2">
        {instances.map((instance) => {
          const displayName = instance.fileName || shortenDicomUid(instance.sopInstanceUid)
          return (
            <li
              key={instance.id}
              className="flex items-center justify-between gap-2 rounded-lg border border-border bg-card p-2"
            >
              <span className="min-w-0 truncate font-mono text-[11px] text-foreground" title={instance.sopInstanceUid}>
                {displayName}
              </span>
              <Button asChild variant="outline" size="sm" className="shrink-0 max-sm:min-h-11">
                <a
                  href={instance.downloadUrl}
                  download={instance.fileName || undefined}
                  aria-label={t("details.files.downloadLabel", { fileName: displayName })}
                >
                  <Download aria-hidden="true" />
                  {t("details.files.download")}
                </a>
              </Button>
            </li>
          )
        })}
      </ul>
    </div>
  )
}
