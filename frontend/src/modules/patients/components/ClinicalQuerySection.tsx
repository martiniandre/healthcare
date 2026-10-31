import type { ReactNode } from "react"
import type { LucideIcon } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Card, CardAction, CardContent, CardHeader, CardTitle } from "../../../shared/components/ui/card"
import { EmptyState } from "../../../shared/components/ui/EmptyState"
import { Skeleton } from "../../../shared/components/ui/skeleton"
import { QueryErrorAlert } from "./QueryErrorAlert"

interface ClinicalQuerySectionProps {
  title: string
  icon: ReactNode
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
  icon,
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
      <CardHeader className="border-b">
        <CardTitle className="flex items-center gap-2 text-base font-extrabold">
          {icon}
          {title}
        </CardTitle>
        {!isLoading && addButton ? <CardAction>{addButton}</CardAction> : null}
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
