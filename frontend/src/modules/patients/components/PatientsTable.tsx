import { useTranslation } from "react-i18next"
import { useNavigate } from "react-router-dom"
import { ArrowRight, ArrowDown, ArrowUp } from "lucide-react"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "../../../shared/components/ui/table"
import { Button } from "../../../shared/components/ui/button"
import { Card } from "../../../shared/components/ui/card"
import { Skeleton } from "../../../shared/components/ui/skeleton"
import { formatDate } from "../../../shared/utils/dates"
import { cn } from "../../../shared/utils/cn"
import type { Patient } from "../queries"

export type SortDirection = "asc" | "desc"
export type SortField = "full_name" | "birth_date" | "document_id"

interface PatientsTableProps {
  isLoading?: boolean
  patients: Patient[]
  sortField: SortField
  sortDirection: SortDirection
  onSort: (field: SortField) => void
  currentPage: number
  totalPages: number
  totalPatients: number
  onPageChange: (page: number) => void
}

const AVATAR_TONES = [
  "bg-sage-100 text-sage-800 dark:bg-primary/15 dark:text-primary",
  "bg-clay-100 text-clay-800 dark:bg-secondary/15 dark:text-secondary",
  "bg-sand-200 text-gray-700 dark:bg-muted dark:text-foreground",
]

const buildInitials = (fullName: string): string => {
  const nameParts = fullName.trim().split(/\s+/).filter(Boolean)
  if (nameParts.length === 0) return "?"
  const firstInitial = nameParts[0].charAt(0)
  const lastInitial = nameParts.length > 1 ? nameParts[nameParts.length - 1].charAt(0) : ""
  return `${firstInitial}${lastInitial}`.toUpperCase()
}

const pickAvatarTone = (seed: string): string => {
  const characterSum = Array.from(seed).reduce((total, character) => total + character.charCodeAt(0), 0)
  return AVATAR_TONES[characterSum % AVATAR_TONES.length]
}

export const PatientsTable = ({
  isLoading,
  patients,
  sortField,
  sortDirection,
  onSort,
  currentPage,
  totalPages,
  totalPatients,
  onPageChange
}: PatientsTableProps) => {
  const { t } = useTranslation("patients")
  const navigate = useNavigate()

  const renderSortableHead = (field: SortField, label: string, className?: string) => {
    const isSorted = sortField === field
    const SortIcon = sortDirection === "asc" ? ArrowUp : ArrowDown
    return (
      <TableHead
        className={cn("cursor-pointer select-none transition-colors hover:text-foreground", isSorted && "text-foreground", className)}
        onClick={() => onSort(field)}
        aria-sort={isSorted ? (sortDirection === "asc" ? "ascending" : "descending") : undefined}
      >
        <span className="inline-flex items-center gap-1">
          {label}
          {isSorted && <SortIcon className="size-3.5" aria-hidden="true" />}
        </span>
      </TableHead>
    )
  }

  return (
    <Card className="gap-0 overflow-hidden py-0">
      <Table>
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            {renderSortableHead("full_name", t("table.patient"), "pl-6")}
            {renderSortableHead("document_id", t("table.document"), "hidden sm:table-cell")}
            {renderSortableHead("birth_date", t("table.birthDate"), "hidden md:table-cell")}
            <TableHead className="hidden lg:table-cell">{t("table.phone")}</TableHead>
            <TableHead className="pr-6 text-right">
              <span className="sr-only sm:not-sr-only">{t("table.action")}</span>
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading ? (
            Array.from({ length: 5 }).map((_, index) => (
              <TableRow key={`skeleton-${index}`}>
                <TableCell className="pl-6">
                  <div className="flex items-center gap-3">
                    <Skeleton className="size-9 rounded-full" />
                    <Skeleton className="h-4 w-36" />
                  </div>
                </TableCell>
                <TableCell className="hidden sm:table-cell"><Skeleton className="h-4 w-28" /></TableCell>
                <TableCell className="hidden md:table-cell"><Skeleton className="h-4 w-20" /></TableCell>
                <TableCell className="hidden lg:table-cell"><Skeleton className="h-4 w-24" /></TableCell>
                <TableCell className="pr-6 text-right">
                  <Skeleton className="ml-auto h-8 w-28" />
                </TableCell>
              </TableRow>
            ))
          ) : (
            patients.map((patient) => (
              <TableRow key={patient.patient_id} className="group">
                <TableCell className="pl-6">
                  <div className="flex min-w-0 items-center gap-3">
                    <span
                      aria-hidden="true"
                      className={cn(
                        "flex size-9 shrink-0 items-center justify-center rounded-full text-xs font-semibold",
                        pickAvatarTone(patient.full_name)
                      )}
                    >
                      {buildInitials(patient.full_name)}
                    </span>
                    <span className="block whitespace-normal break-words text-sm font-semibold text-foreground">
                      {patient.full_name}
                    </span>
                  </div>
                </TableCell>
                <TableCell className="hidden sm:table-cell">
                  <span className="font-mono text-xs tabular-nums text-muted-foreground">{patient.document_id}</span>
                </TableCell>
                <TableCell className="hidden md:table-cell">
                  <span className="text-sm tabular-nums text-muted-foreground">{formatDate(patient.birth_date)}</span>
                </TableCell>
                <TableCell className="hidden lg:table-cell">
                  <span className="text-sm tabular-nums text-muted-foreground">{patient.phone_number}</span>
                </TableCell>
                <TableCell className="pr-6 text-right">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => navigate(`/patients/${patient.fhir_resource_id}`)}
                    aria-label={t("table.medicalRecord")}
                    className="group-hover:border-primary/40 group-hover:bg-primary group-hover:text-primary-foreground"
                  >
                    <span className="hidden sm:inline">{t("table.medicalRecord")}</span>
                    <ArrowRight className="transition-transform group-hover:translate-x-0.5" />
                  </Button>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
      {totalPages > 1 && (
        <div className="flex flex-col gap-3 border-t border-border/60 px-6 py-4 sm:flex-row sm:items-center sm:justify-between">
          <span className="text-xs font-medium tabular-nums text-muted-foreground">
            {t("pagination.pageInfo", { currentPage, totalPages, total: totalPatients })}
          </span>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={currentPage === 1}
              onClick={() => onPageChange(Math.max(currentPage - 1, 1))}
            >
              {t("pagination.prev")}
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={currentPage === totalPages}
              onClick={() => onPageChange(Math.min(currentPage + 1, totalPages))}
            >
              {t("pagination.next")}
            </Button>
          </div>
        </div>
      )}
    </Card>
  )
}
