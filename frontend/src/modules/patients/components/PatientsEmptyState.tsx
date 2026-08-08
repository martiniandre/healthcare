import { useTranslation } from "react-i18next"
import { Users } from "lucide-react"
import { EmptyState } from "../../../shared/components/ui/EmptyState"

interface PatientsEmptyStateProps {
  hasSearchTerm: boolean
  searchTerm: string
}

export const PatientsEmptyState = ({ hasSearchTerm, searchTerm }: PatientsEmptyStateProps) => {
  const { t } = useTranslation("patients")

  return (
    <EmptyState
      icon={Users}
      title={hasSearchTerm ? t("noResults") : t("noPatients")}
      description={hasSearchTerm ? t("noResultsDesc", { searchTerm }) : t("noPatientsDesc")}
      className="flex-1 py-20"
    />
  )
}
