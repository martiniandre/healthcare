import { useTranslation } from "react-i18next"
import { Users } from "lucide-react"
import { EmptyState } from "../../../shared/components/ui/EmptyState"

interface PatientsEmptyStateProps {
  hasSearchTerm: boolean
  searchTerm: string
}

export const PatientsEmptyState = ({ hasSearchTerm, searchTerm }: PatientsEmptyStateProps) => {
  const { t } = useTranslation()

  return (
    <EmptyState
      icon={Users}
      title={hasSearchTerm ? t("patients.noResults") : t("patients.noPatients")}
      description={hasSearchTerm ? t("patients.noResultsDesc", { searchTerm }) : t("patients.noPatientsDesc")}
      className="flex-1 py-20"
    />
  )
}
