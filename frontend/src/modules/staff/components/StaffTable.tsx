import { useReactTable, getCoreRowModel, flexRender } from "@tanstack/react-table"
import { useTranslation } from "react-i18next"
import { TriangleAlert, UserSearch } from "lucide-react"
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "@/shared/components/ui/table"
import { Skeleton } from "@/shared/components/ui/skeleton"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { Button } from "@/shared/components/ui/button"
import { EmptyState } from "@/shared/components/ui/EmptyState"
import { useLayoutStore } from "@/shared/store/layout_store"
import { TableDensityToggle } from "@/shared/components/clinical/TableDensityToggle"
import { useStaffColumns } from "./useStaffColumns"
import type { StaffMember } from "../types"

interface StaffTableProps {
  isLoading: boolean
  isError: boolean
  onRetry: () => void
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
  isError,
  onRetry,
  filteredStaff,
  onToggleStatus,
  onDelete,
  canManageStatus,
  canDelete,
  isPendingEmployeeId,
}: StaffTableProps) => {
  const { t } = useTranslation("staff")
  const tableDensity = useLayoutStore((state) => state.tableDensity)
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

  if (isError) {
    return (
      <Alert variant="destructive">
        <TriangleAlert aria-hidden="true" />
        <AlertTitle>{t("error.title")}</AlertTitle>
        <AlertDescription>
          <p>{t("error.description")}</p>
          <Button type="button" variant="outline" size="sm" onClick={onRetry} className="mt-2">
            {t("error.retry")}
          </Button>
        </AlertDescription>
      </Alert>
    )
  }

  if (!isLoading && table.getRowModel().rows.length === 0) {
    return <EmptyState icon={UserSearch} title={t("emptyState")} className="py-12" />
  }

  return (
    <div className="flex flex-col gap-2">
      <div className="flex justify-end">
        <TableDensityToggle />
      </div>
      <div className="w-full overflow-hidden rounded-xl border border-border bg-card">
        <Table
          className="min-w-[860px] md:min-w-0"
          density={tableDensity}
          containerClassName="max-h-[560px] overflow-y-auto"
        >
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id} className="hover:bg-transparent">
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id} className="sticky top-0 z-10 bg-card">
                    {flexRender(header.column.columnDef.header, header.getContext())}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {isLoading
              ? skeletonRows.map((_, index) => (
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
              : table.getRowModel().rows.map((row) => (
                  <TableRow key={row.id} className="group">
                    {row.getVisibleCells().map((cell) => (
                      <TableCell key={cell.id}>
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </TableCell>
                    ))}
                  </TableRow>
                ))}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
