import type { ReactNode } from "react"
import { AlertCircle, type LucideIcon } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Alert, AlertDescription, AlertTitle } from "../../../shared/components/ui/alert"
import { Button } from "../../../shared/components/ui/button"
import { Card, CardContent } from "../../../shared/components/ui/card"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import { Skeleton } from "../../../shared/components/ui/skeleton"

type SkeletonLayout = "list" | "grid"

interface PortalSectionStateProps {
  isLoading: boolean
  isError: boolean
  onRetry: () => void
  isEmpty: boolean
  emptyTitle: string
  emptyIcon?: LucideIcon
  skeletonLayout?: SkeletonLayout
  children: ReactNode
}

const skeletonItems = [0, 1, 2]

const SectionSkeleton = ({ layout }: { layout: SkeletonLayout }) => {
  const containerClassName =
    layout === "grid" ? "grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3" : "flex flex-col gap-3"

  return (
    <div className={containerClassName} data-testid="portal-section-skeleton" aria-busy="true">
      {skeletonItems.map((item) => (
        <Card key={item} className="py-5">
          <CardContent className="flex flex-col gap-3 px-4 sm:px-5">
            <Skeleton className="h-4 w-2/3" />
            <Skeleton className="h-3 w-1/3" />
            <Skeleton className="h-10 w-full" />
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

export const PortalSectionState = ({
  isLoading,
  isError,
  onRetry,
  isEmpty,
  emptyTitle,
  emptyIcon,
  skeletonLayout = "list",
  children,
}: PortalSectionStateProps) => {
  const { t } = useTranslation("portal")

  if (isLoading) {
    return <SectionSkeleton layout={skeletonLayout} />
  }

  if (isError) {
    return (
      <Alert variant="destructive">
        <AlertCircle aria-hidden="true" />
        <AlertTitle>{t("error.title")}</AlertTitle>
        <AlertDescription>
          <p>{t("error.description")}</p>
          <Button type="button" variant="outline" className="mt-2 min-h-11 min-w-11" onClick={onRetry}>
            {t("error.retry")}
          </Button>
        </AlertDescription>
      </Alert>
    )
  }

  if (isEmpty) {
    return <EmptyState icon={emptyIcon} title={emptyTitle} className="py-12" />
  }

  return <>{children}</>
}
