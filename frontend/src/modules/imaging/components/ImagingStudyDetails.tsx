import { useTranslation } from "react-i18next"
import { Eye } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/components/ui/card"
import { StatusBadge } from "@/shared/components/ui/status-badge"
import { useClinicalStatusLabels } from "../../../shared/hooks/useClinicalStatusLabels"
import { resolveStudyClinicalStatus } from "../utils/pacs_helpers"
import type { ImagingStudy } from "../types"

interface ImagingStudyDetailsProperties {
  study: ImagingStudy
}

const DETAIL_LABEL_CLASS = "block text-[10px] font-bold uppercase tracking-wider text-muted-foreground"
const DETAIL_BOX_CLASS = "rounded-xl border border-border bg-muted p-3.5"

export const ImagingStudyDetails = ({ study }: ImagingStudyDetailsProperties) => {
  const { t } = useTranslation("imaging")
  const { imagingStudyStatus } = useClinicalStatusLabels()

  return (
    <Card className="lg:col-span-1 text-left">
      <CardHeader className="border-b">
        <CardTitle className="flex items-center gap-2 text-base font-extrabold">
          <Eye className="size-4 text-primary" aria-hidden="true" />
          {t("details.title")}
        </CardTitle>
      </CardHeader>

      <CardContent className="flex flex-col gap-4">
        <div className={DETAIL_BOX_CLASS}>
          <span className={DETAIL_LABEL_CLASS}>{t("details.id")}</span>
          <span className="mt-1 block text-xs font-bold text-foreground">{study.id}</span>
        </div>

        <div className={DETAIL_BOX_CLASS}>
          <span className={DETAIL_LABEL_CLASS}>{t("details.modality")}</span>
          <span className="mt-1 block text-xs font-bold uppercase text-foreground">{study.modality}</span>
        </div>

        <div className={DETAIL_BOX_CLASS}>
          <span className={DETAIL_LABEL_CLASS}>{t("details.busStatus")}</span>
          <StatusBadge status={resolveStudyClinicalStatus(study.status)} className="mt-1">
            {imagingStudyStatus(study.status)}
          </StatusBadge>
        </div>

        <div className={DETAIL_BOX_CLASS}>
          <span className={DETAIL_LABEL_CLASS}>{t("details.seriesSlices")}</span>
          <span className="mt-1 block text-xs font-bold text-foreground">{t("details.activeSlice")}</span>
        </div>
      </CardContent>
    </Card>
  )
}
