import type { ReactElement, ReactNode } from "react"
import type { ColumnDef } from "@tanstack/react-table"
import type { LucideIcon } from "lucide-react"
import { ClinicalTable } from "../../../shared/components/clinical/ClinicalTable"
import { ClinicalQuerySection } from "./ClinicalQuerySection"

interface ClinicalQueryTableProps<T> {
  title: string
  icon: ReactNode
  columns: ColumnDef<T>[]
  data: T[]
  isLoading: boolean
  isError: boolean
  onRetry: () => void
  emptyText: string
  emptyIcon: LucideIcon
  addButton?: ReactNode
}

export function ClinicalQueryTable<T>({
  title,
  icon,
  columns,
  data,
  isLoading,
  isError,
  onRetry,
  emptyText,
  emptyIcon: EmptyIcon,
  addButton,
}: ClinicalQueryTableProps<T>) {
  return (
    <ClinicalQuerySection
      title={title}
      icon={icon}
      isLoading={isLoading}
      isError={isError}
      onRetry={onRetry}
      isEmpty={data.length === 0}
      emptyTitle={emptyText}
      emptyIcon={EmptyIcon}
      addButton={addButton}
    >
      <ClinicalTable
        title={title}
        icon={icon as ReactElement}
        columns={columns}
        data={data}
        isEmpty={false}
        emptyIcon={<EmptyIcon />}
        emptyText={emptyText}
        addButton={addButton}
      />
    </ClinicalQuerySection>
  )
}
