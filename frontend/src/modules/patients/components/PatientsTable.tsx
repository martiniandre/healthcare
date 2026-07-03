import { useTranslation } from "react-i18next"
import { useNavigate } from "react-router-dom"
import { User, ArrowRight } from "lucide-react"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "../../../shared/components/ui/Table"
import { Button } from "../../../shared/components/ui/Button"
import { Skeleton } from "../../../shared/components/ui/Skeleton"
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
  const { t } = useTranslation()
  const navigate = useNavigate()

  const sortIndicator = (field: SortField) => {
    if (sortField !== field) return ""
    return sortDirection === "asc" ? " ↑" : " ↓"
  }

  return (
    <div className="bg-white border border-border rounded-xl overflow-hidden">
      <div className="overflow-x-auto w-full">
        <Table className="min-w-[650px] md:min-w-0">
          <TableHeader>
            <TableRow>
              <TableHead
                className="cursor-pointer hover:text-gray-700 transition-colors select-none"
                onClick={() => onSort("full_name")}
              >
                {t("patients.table.patient")}{sortIndicator("full_name")}
              </TableHead>
              <TableHead
                className="cursor-pointer hover:text-gray-700 transition-colors select-none"
                onClick={() => onSort("document_id")}
              >
                {t("patients.table.document")}{sortIndicator("document_id")}
              </TableHead>
              <TableHead
                className="cursor-pointer hover:text-gray-700 transition-colors select-none"
                onClick={() => onSort("birth_date")}
              >
                {t("patients.table.birthDate")}{sortIndicator("birth_date")}
              </TableHead>
              <TableHead>
                {t("patients.table.phone")}
              </TableHead>
              <TableHead className="text-right">
                {t("patients.table.action")}
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: 5 }).map((_, index) => (
                <TableRow key={`skeleton-${index}`}>
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <Skeleton className="w-8 h-8 rounded-lg" />
                      <div className="flex flex-col gap-1.5 min-w-0">
                        <Skeleton className="h-4 w-32" />
                        <Skeleton className="h-3 w-48 mt-0.5" />
                      </div>
                    </div>
                  </TableCell>
                  <TableCell><Skeleton className="h-4 w-28" /></TableCell>
                  <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                  <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                  <TableCell className="text-right">
                    <Skeleton className="h-6 w-24 ml-auto" />
                  </TableCell>
                </TableRow>
              ))
            ) : (
              patients.map((patient) => (
                <TableRow
                  key={patient.patient_id}
                  className="group"
                >
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-lg bg-primary/6 flex items-center justify-center shrink-0 group-hover:bg-primary/10 transition-colors">
                        <User className="w-4 h-4 text-primary/60 group-hover:text-primary transition-colors" />
                      </div>
                      <div className="min-w-0">
                        <span className="text-[13px] font-semibold text-gray-800 block truncate group-hover:text-primary transition-colors">
                          {patient.full_name}
                        </span>
                        <span className="text-[10px] text-gray-400 font-mono block mt-0.5">
                          {patient.fhir_resource_id}
                        </span>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <span className="text-xs font-mono text-gray-600">
                      {patient.document_id}
                    </span>
                  </TableCell>
                  <TableCell>
                    <span className="text-xs text-gray-500">
                      {patient.birth_date}
                    </span>
                  </TableCell>
                  <TableCell>
                    <span className="text-xs text-gray-500">
                      {patient.phone_number}
                    </span>
                  </TableCell>
                  <TableCell className="text-right">
                    <button
                      onClick={() => navigate(`/patients/${patient.fhir_resource_id}`)}
                      className="inline-flex items-center gap-1.5 text-[11px] font-semibold text-primary hover:text-primary/80 px-3 py-1.5 rounded-md hover:bg-primary/5 transition-all"
                    >
                      {t("patients.table.medicalRecord")}
                      <ArrowRight className="w-3 h-3 transition-transform group-hover:translate-x-0.5" />
                    </button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
      {totalPages > 1 && (
        <div className="flex items-center justify-between border-t border-border px-6 py-4 bg-gray-50/50">
          <span className="text-xs text-gray-500 font-semibold">
            {t("patients.pagination.pageInfo", { currentPage, totalPages, total: totalPatients })}
          </span>
          <div className="flex gap-2">
            <Button
              variantType="outline"
              className="px-3 py-1.5 text-xs font-bold"
              disabled={currentPage === 1}
              onClick={() => onPageChange(Math.max(currentPage - 1, 1))}
            >
              {t("patients.pagination.prev")}
            </Button>
            <Button
              variantType="outline"
              className="px-3 py-1.5 text-xs font-bold"
              disabled={currentPage === totalPages}
              onClick={() => onPageChange(Math.min(currentPage + 1, totalPages))}
            >
              {t("patients.pagination.next")}
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}
