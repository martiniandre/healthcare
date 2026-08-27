import { useTranslation } from "react-i18next"
import { staffColorForIndex } from "../schedule_calendar_helpers"
import { Label } from "../../../shared/components/ui/Label"
import { Checkbox } from "../../../shared/components/ui/Checkbox"

interface StaffOverlaySidebarProps {
  staffMembers: Array<{ id: string; fullName: string; role: string }>
  selectedStaffIds: string[]
  onToggleStaff: (staffId: string) => void
}

export const StaffOverlaySidebar = ({ staffMembers, selectedStaffIds, onToggleStaff }: StaffOverlaySidebarProps) => {
  const { t } = useTranslation("schedule")

  return (
    <div className="bg-surface border border-border rounded-xl p-4 flex flex-col gap-2">
      <span className="text-xs font-semibold text-muted-foreground">{t("calendar.staffFilters")}</span>
      {staffMembers.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t("emptyState.noStaff")}</p>
      ) : (
        staffMembers.map((staffMember, staffIndex) => {
          const isSelected = selectedStaffIds.includes(staffMember.id)
          return (
            <Label
              key={staffMember.id}
              onClick={() => onToggleStaff(staffMember.id)}
              className="flex items-center gap-2 rounded-lg px-2 py-1.5 hover:bg-muted-soft cursor-pointer select-none"
            >
              <Checkbox checked={isSelected} />
              <span
                className="w-3 h-3 rounded-full shrink-0"
                style={{ backgroundColor: staffColorForIndex(staffIndex) }}
              />
              <span className="text-sm text-foreground truncate">
                {staffMember.fullName} — {staffMember.role}
              </span>
            </Label>
          )
        })
      )}
    </div>
  )
}
