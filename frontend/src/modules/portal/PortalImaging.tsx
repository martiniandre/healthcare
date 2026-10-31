import { Image } from "lucide-react"
import { useTranslation } from "react-i18next"
import { usePortalImagingQuery } from "./queries"
import { PortalSectionState } from "./components/PortalSectionState"
import { Card, CardContent } from "../../shared/components/ui/card"
import { Badge } from "../../shared/components/ui/badge"
import { StatusBadge } from "../../shared/components/ui/status-badge"
import { formatDate } from "../../shared/utils/dates"
import { useLocale } from "../../shared/hooks/useLocale"
import { useClinicalStatusLabels } from "../../shared/hooks/useClinicalStatusLabels"

export const PortalImaging = () => {
  const { t } = useTranslation("portal")
  const { data: imagingStudies, isLoading, isError, refetch } = usePortalImagingQuery()
  const locale = useLocale()
  const { dicomModality, imagingStudyStatus, notInformed } = useClinicalStatusLabels()

  return (
    <PortalSectionState
      isLoading={isLoading}
      isError={isError}
      onRetry={() => refetch()}
      isEmpty={!imagingStudies || imagingStudies.length === 0}
      emptyTitle={t("empty.imaging")}
      emptyIcon={Image}
      skeletonLayout="grid"
    >
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {imagingStudies?.map((study) => (
          <Card key={study.fhir_resource_id} className="py-4 sm:py-5">
            <CardContent className="px-4 sm:px-5">
              <div className="mb-3 flex flex-wrap items-start justify-between gap-2">
                <Badge variant="outline" className="rounded uppercase">
                  {study.modality ? dicomModality(study.modality) : notInformed()}
                </Badge>
                <span className="text-xs tabular-nums text-muted-foreground">
                  {study.created_at ? formatDate(study.created_at, locale) : ""}
                </span>
              </div>
              <p className="mb-2 break-words text-sm font-bold text-foreground">
                {study.title || t("fallback.imagingStudy")}
              </p>
              <StatusBadge status="info">{imagingStudyStatus(study.status)}</StatusBadge>
            </CardContent>
          </Card>
        ))}
      </div>
    </PortalSectionState>
  )
}
