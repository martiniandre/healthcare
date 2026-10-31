import { useTranslation } from "react-i18next"
import { TriangleAlert, Users } from "lucide-react"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { Button } from "@/shared/components/ui/button"
import { EmptyState } from "@/shared/components/ui/EmptyState"
import { Skeleton } from "@/shared/components/ui/skeleton"
import { Card, CardContent } from "@/shared/components/ui/card"
import { Checkbox } from "@/shared/components/ui/checkbox"
import { Label } from "@/shared/components/ui/label"
import { staffColorForIndex } from "../schedule_calendar_helpers"

interface StaffOverlaySidebarProps {
  staffMembers: Array<{ id: string; fullName: string; role: string }>
  selectedStaffIds: string[]
  onToggleStaff: (staffId: string) => void
  isLoading?: boolean
  isError?: boolean
  onRetry?: () => void
}

export const StaffOverlaySidebar = ({
  staffMembers,
  selectedStaffIds,
  onToggleStaff,
  isLoading = false,
  isError = false,
  onRetry,
}: StaffOverlaySidebarProps) => {
  const { t } = useTranslation("schedule")

  return (
    <Card className="py-4">
      <CardContent className="flex flex-col gap-2 px-4">
        <span className="text-xs font-semibold text-muted-foreground">{t("calendar.staffFilters")}</span>
        {isLoading ? (
          <div className="flex flex-col gap-2" role="status" aria-label={t("loading.staff")}>
            {Array.from({ length: 4 }).map((_, skeletonIndex) => (
              <Skeleton key={skeletonIndex} className="h-8 w-full rounded-lg" />
            ))}
          </div>
        ) : isError ? (
          <Alert variant="destructive">
            <TriangleAlert aria-hidden="true" />
            <AlertTitle>{t("errors.loadStaffTitle")}</AlertTitle>
            <AlertDescription>
              {onRetry ? (
                <Button type="button" variant="outline" size="sm" onClick={onRetry}>
                  {t("retry")}
                </Button>
              ) : null}
            </AlertDescription>
          </Alert>
        ) : staffMembers.length === 0 ? (
          <EmptyState icon={Users} title={t("emptyState.noStaff")} className="p-4" />
        ) : (
          staffMembers.map((staffMember, staffIndex) => {
            const isSelected = selectedStaffIds.includes(staffMember.id)
            const checkboxId = `staff-overlay-${staffMember.id}`
            return (
              <Label
                key={staffMember.id}
                htmlFor={checkboxId}
                className="rounded-lg px-2 py-1.5 hover:bg-muted cursor-pointer select-none"
              >
                <Checkbox
                  id={checkboxId}
                  checked={isSelected}
                  onCheckedChange={() => onToggleStaff(staffMember.id)}
                />
                <span
                  className="size-3 rounded-full shrink-0"
                  style={{ backgroundColor: staffColorForIndex(staffIndex) }}
                />
                <span className="text-sm text-foreground truncate">
                  {staffMember.fullName} — {staffMember.role}
                </span>
              </Label>
            )
          })
        )}
      </CardContent>
    </Card>
  )
}
