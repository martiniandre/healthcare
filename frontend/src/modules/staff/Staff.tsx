import { useState } from "react"
import { useTranslation } from "react-i18next"
import { useDebounce } from "../../shared/hooks/useDebounce"
import { Card } from "../../shared/components/ui/Card"
import { PageContainer } from "../../shared/components/ui/PageContainer"
import { toast } from "../../shared/store/toast_store"
import { useAuthStore } from "../../shared/store/auth_store"
import {
  useStaffListQuery,
  useSetEmployeeStatusMutation,
  useDeleteEmployeeMutation,
} from "./queries"
import type { StaffMember } from "./types"
import { StaffHeader } from "./components/StaffHeader"
import { StaffFilters } from "./components/StaffFilters"
import { StaffTable } from "./components/StaffTable"
import { StaffModal } from "./components/StaffModal"
import { DeleteStaffModal } from "./components/DeleteStaffModal"

const rolesAllowedToManageStatus = ["ADMIN", "DOCTOR", "NURSE"]

export const Staff = () => {
  const { t } = useTranslation("staff")
  const [filterRole, setFilterRole] = useState<string>("All")
  const [searchQuery, setSearchQuery] = useState<string>("")
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false)
  const [memberToDelete, setMemberToDelete] = useState<StaffMember | null>(null)
  const [pendingEmployeeId, setPendingEmployeeId] = useState<string | null>(null)

  const debouncedSearchQuery = useDebounce(searchQuery, 500)

  const sessionRole = useAuthStore((state) => state.role)
  const { data: staffList = [], isLoading } = useStaffListQuery(debouncedSearchQuery, filterRole)
  const setEmployeeStatusMutation = useSetEmployeeStatusMutation()
  const deleteEmployeeMutation = useDeleteEmployeeMutation()

  const canManageStatus = sessionRole !== null && rolesAllowedToManageStatus.includes(sessionRole)
  const canDelete = sessionRole === "ADMIN"

  const handleToggleStatus = async (member: StaffMember) => {
    const nextIsActive = !member.isActive
    setPendingEmployeeId(member.id)
    try {
      await setEmployeeStatusMutation.mutateAsync({ employeeId: member.id, isActive: nextIsActive })
      toast.success(t(nextIsActive ? "toast.activateSuccess" : "toast.deactivateSuccess"))
    } catch {
      toast.error(t(nextIsActive ? "toast.activateError" : "toast.deactivateError"))
    } finally {
      setPendingEmployeeId(null)
    }
  }

  const handleConfirmDelete = async () => {
    if (!memberToDelete) {
      return
    }
    setPendingEmployeeId(memberToDelete.id)
    try {
      await deleteEmployeeMutation.mutateAsync(memberToDelete.id)
      setMemberToDelete(null)
      toast.success(t("toast.deleteSuccess"))
    } catch {
      toast.error(t("toast.deleteError"))
    } finally {
      setPendingEmployeeId(null)
    }
  }

  return (
    <PageContainer className="select-none relative">
      <StaffHeader onAddStaff={() => setIsModalOpen(true)} />

      <Card className="p-4 flex flex-col gap-4">
        <StaffFilters
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          filterRole={filterRole}
          onFilterChange={setFilterRole}
        />

        <StaffTable
          isLoading={isLoading}
          filteredStaff={staffList}
          onToggleStatus={handleToggleStatus}
          onDelete={setMemberToDelete}
          canManageStatus={canManageStatus}
          canDelete={canDelete}
          isPendingEmployeeId={pendingEmployeeId}
        />
      </Card>

      <StaffModal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />

      <DeleteStaffModal
        member={memberToDelete}
        onClose={() => setMemberToDelete(null)}
        onConfirm={handleConfirmDelete}
        isPending={deleteEmployeeMutation.isPending}
      />
    </PageContainer>
  )
}
