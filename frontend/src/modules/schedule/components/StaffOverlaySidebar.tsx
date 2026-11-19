import type { ReactNode } from "react"
import { useTranslation } from "react-i18next"
import { TriangleAlert, Users } from "lucide-react"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { Button } from "@/shared/components/ui/button"
import { EmptyState } from "@/shared/components/ui/EmptyState"
import { Skeleton } from "@/shared/components/ui/skeleton"
import { Card, CardContent } from "@/shared/components/ui/card"
import { Checkbox } from "@/shared/components/ui/checkbox"
import { Label } from "@/shared/components/ui/label"
import { cn } from "@/shared/utils/cn"
import { staffColorForIndex } from "../schedule_calendar_helpers"

interface StaffOverlaySidebarProps {
  staffMembers: Array<{ id: string; fullName: string; role: string }>
  selectedStaffIds: string[]
  onToggleStaff: (staffId: string) => void
  isLoading?: boolean
  isError?: boolean
  onRetry?: () => void
  viewControl?: ReactNode
}

export const StaffOverlaySidebar = ({
  staffMembers,
  selectedStaffIds,
  onToggleStaff,
  isLoading = false,
  isError = false,
  onRetry,
  viewControl,
}: StaffOverlaySidebarProps) => {
  const { t } = useTranslation("schedule")

  return (
    <Card className="gap-5 py-5 lg:sticky lg:top-20">
      {viewControl && (
        <CardContent className="flex flex-col gap-2 px-5">
          <span className="text-xs font-semibold text-muted-foreground">{t("calendar.viewLabel")}</span>
          {viewControl}
        </CardContent>
      )}
      <CardContent className="flex flex-col gap-1.5 px-3">
        <span className="px-2 text-xs font-semibold text-muted-foreground">{t("calendar.staffFilters")}</span>
        {isLoading ? (
          <div className="flex flex-col gap-2 px-2" role="status" aria-label={t("loading.staff")}>
            {Array.from({ length: 4 }).map((_, skeletonIndex) => (
              <Skeleton key={skeletonIndex} className="h-12 w-full" />
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
            const staffColor = staffColorForIndex(staffIndex)
            return (
              <Label
                key={staffMember.id}
                htmlFor={checkboxId}
                className={cn(
                  "cursor-pointer select-none items-start gap-3 rounded-md px-2 py-2 transition-colors hover:bg-accent",
                  isSelected && "bg-muted"
                )}
              >
                <Checkbox
                  id={checkboxId}
                  checked={isSelected}
                  onCheckedChange={() => onToggleStaff(staffMember.id)}
                  className="mt-0.5"
                />
                <span className="flex min-w-0 flex-1 flex-col gap-0.5">
                  <span className="flex items-center gap-2 text-sm font-semibold leading-snug text-foreground">
                    <span aria-hidden="true" className="size-2.5 shrink-0 rounded-full" style={{ backgroundColor: staffColor }} />
                    <span className="break-words">{staffMember.fullName}</span>
                  </span>
                  <span className="pl-[18px] text-xs font-normal text-muted-foreground first-letter:uppercase">{staffMember.role}</span>
                </span>
              </Label>
            )
          })
        )}
      </CardContent>
    </Card>
  )
}
