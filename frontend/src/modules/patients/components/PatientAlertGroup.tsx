import { RotateCw } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Button } from "../../../shared/components/ui/button"
import { Skeleton } from "../../../shared/components/ui/skeleton"
import { StatusBadge, type ClinicalStatus } from "../../../shared/components/ui/status-badge"

const MAX_VISIBLE_ITEMS = 3

interface PatientAlertGroupProps {
  label: string
  items: string[]
  itemStatus: ClinicalStatus
  emptyLabel: string
  isLoading: boolean
  isError: boolean
  onRetry: () => void
  showAllItems?: boolean
  testId: string
}

export const PatientAlertGroup = ({
  label,
  items,
  itemStatus,
  emptyLabel,
  isLoading,
  isError,
  onRetry,
  showAllItems = false,
  testId,
}: PatientAlertGroupProps) => {
  const { t } = useTranslation("patients")
  const visibleItems = showAllItems ? items : items.slice(0, MAX_VISIBLE_ITEMS)
  const hiddenCount = items.length - visibleItems.length

  return (
    <div data-testid={testId} className="flex flex-wrap items-center gap-1.5">
      <span className="text-[10px] font-black uppercase tracking-wider text-muted-foreground">{label}</span>
      {isLoading ? (
        <Skeleton data-testid={`${testId}-loading`} className="h-5 w-24" />
      ) : isError ? (
        <Button type="button" variant="ghost" size="xs" onClick={onRetry} className="text-destructive">
          <RotateCw />
          {t("states.retry")}
        </Button>
      ) : items.length === 0 ? (
        <StatusBadge status="neutral">{emptyLabel}</StatusBadge>
      ) : (
        <>
          {visibleItems.map((itemName) => (
            <StatusBadge key={itemName} status={itemStatus}>
              {itemName}
            </StatusBadge>
          ))}
          {hiddenCount > 0 && (
            <StatusBadge status={itemStatus}>{t("header.moreItems", { count: hiddenCount })}</StatusBadge>
          )}
        </>
      )}
    </div>
  )
}
