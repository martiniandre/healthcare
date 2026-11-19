import type { ReactNode } from "react"
import type { LucideIcon } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Card, CardContent, CardHeader, CardTitle } from "../../../shared/components/ui/card"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import { Skeleton } from "../../../shared/components/ui/skeleton"
import { QueryErrorAlert } from "@/shared/components/ui/QueryErrorAlert"

interface ClinicalQuerySectionProps {
  title: string
  isLoading: boolean
  isError: boolean
  onRetry: () => void
  isEmpty: boolean
  emptyTitle: string
  emptyIcon?: LucideIcon
  addButton?: ReactNode
  children: ReactNode
}

const skeletonRows = ["first", "second", "third", "fourth"]

export const ClinicalQuerySection = ({
  title,
  isLoading,
  isError,
  onRetry,
  isEmpty,
  emptyTitle,
  emptyIcon,
  addButton,
  children,
}: ClinicalQuerySectionProps) => {
  const { t } = useTranslation("patients")

  if (!isLoading && !isError && !isEmpty) {
    return <>{children}</>
  }

  return (
    <Card className="min-h-[450px] gap-5">
      <CardHeader className="flex flex-col gap-3 border-b sm:flex-row sm:items-center sm:justify-between">
        <CardTitle className="flex min-w-0 items-center gap-2 text-base font-semibold [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:text-muted-foreground">
          <span className="min-w-0">{title}</span>
        </CardTitle>
        {!isLoading && addButton ? <div className="flex flex-wrap items-center gap-2 sm:shrink-0">{addButton}</div> : null}
      </CardHeader>
      <CardContent className="flex flex-1 flex-col gap-3">
        {isLoading ? (
          <div role="status" aria-busy="true" aria-label={t("states.loading")} className="flex flex-col gap-3">
            {skeletonRows.map((rowKey) => (
              <Skeleton key={rowKey} className="h-10 w-full" />
            ))}
          </div>
        ) : isError ? (
          <QueryErrorAlert onRetry={onRetry} />
        ) : (
          <EmptyState icon={emptyIcon} title={emptyTitle} className="flex-1" />
        )}
      </CardContent>
    </Card>
  )
}
