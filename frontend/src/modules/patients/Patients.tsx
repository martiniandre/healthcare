import { useState } from "react"
import { useDebounce } from "../../shared/hooks/useDebounce"
import { useSearchParams } from "react-router-dom"
import { PatientsHeader } from "./components/PatientsHeader"
import { PatientsMetricsGrid } from "./components/PatientsMetricsGrid"
import { PatientsFilters } from "./components/PatientsFilters"
import { PatientsLoadingState } from "./components/PatientsLoadingState"
import { PatientsEmptyState } from "./components/PatientsEmptyState"
import { PatientsTable, type SortField, type SortDirection } from "./components/PatientsTable"
import { PatientModal } from "./components/PatientModal"
import { usePatientsQuery } from "./queries"

export const Patients = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [searchTerm, setSearchTerm] = useState("")
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [sort, setSort] = useState<{
    field: SortField
    direction: SortDirection
  }>({ field: "full_name", direction: "asc" })
  
  const currentPage = parseInt(searchParams.get("page") || "1", 10)
  const itemsPerPage = 5

  const setCurrentPage = (page: number) => {
    setSearchParams((prev) => {
      if (page <= 1) {
        prev.delete("page")
      } else {
        prev.set("page", page.toString())
      }
      return prev
    }, { replace: true })
  }

  const debouncedSearchTerm = useDebounce(searchTerm, 500)

  const { data: patients = [], isLoading } = usePatientsQuery(
    debouncedSearchTerm,
    sort.field,
    sort.direction,
    currentPage,
    itemsPerPage
  )

  const handleSortToggle = (field: SortField) => {
    setSort((prev) => {
      if (prev.field === field) {
        return { field, direction: prev.direction === "asc" ? "desc" : "asc" }
      }
      return { field, direction: "asc" }
    })
    setCurrentPage(1)
  }

  const handleSearchChange = (termValue: string) => {
    setSearchTerm(termValue)
    setCurrentPage(1)
  }

  const totalPages = patients.length === itemsPerPage ? currentPage + 1 : currentPage

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-5 animate-fade-in">
      <PatientsHeader onNewPatient={() => setIsModalOpen(true)} />
      
      <PatientsMetricsGrid totalPatients={patients.length} />

      <PatientsFilters
        searchTerm={searchTerm}
        onSearchChange={handleSearchChange}
        resultsCount={patients.length}
      />

      {isLoading ? (
        <PatientsLoadingState />
      ) : patients.length === 0 ? (
        <PatientsEmptyState 
          hasSearchTerm={!!searchTerm} 
          searchTerm={searchTerm} 
        />
      ) : (
        <PatientsTable
          patients={patients}
          sortField={sort.field}
          sortDirection={sort.direction}
          onSort={handleSortToggle}
          currentPage={currentPage}
          totalPages={totalPages}
          totalPatients={-1}
          onPageChange={setCurrentPage}
        />
      )}

      <PatientModal 
        isOpen={isModalOpen} 
        onOpenChange={setIsModalOpen} 
      />
    </div>
  )
}
