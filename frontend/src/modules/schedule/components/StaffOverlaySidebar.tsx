import { useTranslation } from "react-i18next"
import { staffColorForIndex } from "../schedule_calendar_helpers"

interface StaffOverlaySidebarProps {
  staffMembers: Array<{ id: string; fullName: string; role: string }>
  selectedStaffIds: string[]
  onToggleStaff: (staffId: string) => void
}

export const StaffOverlaySidebar = ({ staffMembers, selectedStaffIds, onToggleStaff }: StaffOverlaySidebarProps) => {
  const { t } = useTranslation("schedule")

  return (
    <div className="bg-white border border-border rounded-xl p-4 flex flex-col gap-2">
      <span className="text-xs font-semibold text-gray-600">{t("calendar.staffFilters")}</span>
      {staffMembers.length === 0 ? (
        <p className="text-sm text-gray-400">{t("emptyState.noStaff")}</p>
      ) : (
        staffMembers.map((staffMember, staffIndex) => {
          const isSelected = selectedStaffIds.includes(staffMember.id)
          return (
            <label
              key={staffMember.id}
              className="flex items-center gap-2 rounded-lg px-2 py-1.5 hover:bg-gray-50 cursor-pointer select-none"
            >
              <input
                type="checkbox"
                checked={isSelected}
                onChange={() => onToggleStaff(staffMember.id)}
                className="w-4 h-4 accent-primary"
              />
              <span
                className="w-3 h-3 rounded-full shrink-0"
                style={{ backgroundColor: staffColorForIndex(staffIndex) }}
              />
              <span className="text-sm text-gray-800 truncate">
                {staffMember.fullName} — {staffMember.role}
              </span>
            </label>
          )
        })
      )}
    </div>
  )
}
