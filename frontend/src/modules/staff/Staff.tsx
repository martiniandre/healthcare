import { useState } from "react"
import { useTranslation } from "react-i18next"
import { useDebounce } from "../../shared/hooks/useDebounce"
import { Card, CardContent } from "@/shared/components/ui/card"
import { PageContainer } from "../../shared/components/ui/PageContainer"
import { toast } from "../../shared/store/toast_store"
import { useAbility } from "../../shared/auth/useAbility"
import { Action, Feature } from "../../shared/auth/types"
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

export const Staff = () => {
  const { t } = useTranslation("staff")
  const [filterRole, setFilterRole] = useState<string>("All")
  const [searchQuery, setSearchQuery] = useState<string>("")
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false)
  const [memberToDelete, setMemberToDelete] = useState<StaffMember | null>(null)
  const [pendingEmployeeId, setPendingEmployeeId] = useState<string | null>(null)

  const debouncedSearchQuery = useDebounce(searchQuery, 500)

  const ability = useAbility()
  const { data: staffList = [], isLoading, isError, refetch } = useStaffListQuery(debouncedSearchQuery, filterRole)
  const setEmployeeStatusMutation = useSetEmployeeStatusMutation()
  const deleteEmployeeMutation = useDeleteEmployeeMutation()

  const canManageStatus = ability.can(Action.Update, Feature.Staff)
  const canDelete = ability.can(Action.Delete, Feature.Staff)

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

      <Card className="py-4">
        <CardContent className="flex flex-col gap-4 px-4">
        <StaffFilters
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          filterRole={filterRole}
          onFilterChange={setFilterRole}
        />

        <StaffTable
          isLoading={isLoading}
          isError={isError}
          onRetry={() => refetch()}
          filteredStaff={staffList}
          onToggleStatus={handleToggleStatus}
          onDelete={setMemberToDelete}
          canManageStatus={canManageStatus}
          canDelete={canDelete}
          isPendingEmployeeId={pendingEmployeeId}
        />
        </CardContent>
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
