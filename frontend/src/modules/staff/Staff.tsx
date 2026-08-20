import { useState } from "react"
import { useDebounce } from "../../shared/hooks/useDebounce"
import { Card } from "../../shared/components/ui/Card"
import { PageContainer } from "../../shared/components/ui/PageContainer"
import { useStaffListQuery } from "./queries"
import { StaffHeader } from "./components/StaffHeader"
import { StaffFilters } from "./components/StaffFilters"
import { StaffTable } from "./components/StaffTable"
import { StaffModal } from "./components/StaffModal"

export const Staff = () => {
  const [filterRole, setFilterRole] = useState<string>("All")
  const [searchQuery, setSearchQuery] = useState<string>("")
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false)

  const debouncedSearchQuery = useDebounce(searchQuery, 500)

  const { data: staffList = [], isLoading } = useStaffListQuery(debouncedSearchQuery, filterRole)

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

        <StaffTable isLoading={isLoading} filteredStaff={staffList} />
      </Card>

      <StaffModal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />
    </PageContainer>
  )
}
