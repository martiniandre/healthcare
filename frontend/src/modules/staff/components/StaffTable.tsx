import { useReactTable, getCoreRowModel, flexRender } from "@tanstack/react-table"
import { useTranslation } from "react-i18next"
import { UserSearch } from "lucide-react"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "../../../shared/components/ui/Table"
import { Skeleton } from "../../../shared/components/ui/Skeleton"
import { useStaffColumns } from "./useStaffColumns"
import type { StaffMember } from "../types"

interface StaffTableProps {
  isLoading: boolean
  filteredStaff: StaffMember[]
  onToggleStatus: (member: StaffMember) => void
  onDelete: (member: StaffMember) => void
  canManageStatus: boolean
  canDelete: boolean
  isPendingEmployeeId: string | null
}

const skeletonRows = Array.from({ length: 5 })

export const StaffTable = ({
  isLoading,
  filteredStaff,
  onToggleStatus,
  onDelete,
  canManageStatus,
  canDelete,
  isPendingEmployeeId,
}: StaffTableProps) => {
  const { t } = useTranslation("staff")
  const columns = useStaffColumns({
    onToggleStatus,
    onDelete,
    canManageStatus,
    canDelete,
    isPendingEmployeeId,
  })

  const table = useReactTable({
    data: filteredStaff,
    columns,
    getCoreRowModel: getCoreRowModel(),
  })

  return (
    <div className="overflow-x-auto border border-border rounded-xl w-full bg-card">
      <Table className="min-w-[860px] md:min-w-0">
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id} className="hover:bg-transparent">
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {flexRender(header.column.columnDef.header, header.getContext())}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {isLoading ? (
            skeletonRows.map((_, index) => (
              <TableRow key={`skeleton-${index}`}>
                <TableCell>
                  <div className="flex items-center gap-3">
                    <Skeleton className="w-8 h-8 rounded-lg" />
                    <div className="flex flex-col gap-1.5">
                      <Skeleton className="h-4 w-32" />
                      <Skeleton className="h-3 w-24" />
                    </div>
                  </div>
                </TableCell>
                <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                <TableCell><Skeleton className="h-8 w-16" /></TableCell>
              </TableRow>
            ))
          ) : table.getRowModel().rows.length === 0 ? (
            <TableRow className="hover:bg-transparent">
              <TableCell colSpan={columns.length}>
                <div className="flex flex-col items-center justify-center gap-2 py-12">
                  <UserSearch className="w-6 h-6 text-gray-300" />
                  <span className="text-xs text-muted">{t("emptyState")}</span>
                </div>
              </TableCell>
            </TableRow>
          ) : (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id} className="group">
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  )
}
