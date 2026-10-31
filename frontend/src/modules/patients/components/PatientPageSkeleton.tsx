import { useTranslation } from "react-i18next"
import { PageContainer } from "../../../shared/components/ui/PageContainer"
import { Skeleton } from "../../../shared/components/ui/skeleton"

export const PatientPageSkeleton = () => {
  const { t } = useTranslation("patients")

  return (
    <PageContainer>
      <div role="status" aria-busy="true" aria-label={t("details.loadingDetails")} className="flex flex-col gap-6">
        <Skeleton className="h-16 w-full" />
        <div className="flex flex-col md:flex-row gap-6">
          <Skeleton className="h-64 w-full md:w-64 shrink-0" />
          <Skeleton className="h-[450px] w-full" />
        </div>
      </div>
    </PageContainer>
  )
}
