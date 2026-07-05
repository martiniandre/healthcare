import { useReactTable, getCoreRowModel, flexRender, type ColumnDef, type SortingState, getSortedRowModel } from "@tanstack/react-table"
import { useState, type ReactNode, type ReactElement } from "react"
import { ChevronUp, ChevronDown } from "lucide-react"
import { Card } from "../ui/Card"
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from "../ui/Table"

interface ClinicalTableProps<T> {
  title: string
  icon: ReactElement
  columns: ColumnDef<T>[]
  data: T[]
  isEmpty: boolean
  emptyIcon: ReactElement
  emptyText: string
  addButton?: ReactNode
  enableSorting?: boolean
}

export function ClinicalTable<T>({
  title,
  icon,
  columns,
  data,
  isEmpty,
  emptyIcon,
  emptyText,
  addButton,
  enableSorting = false,
}: ClinicalTableProps<T>) {
  const [sorting, setSorting] = useState<SortingState>([])

  const table = useReactTable({
    data,
    columns,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: enableSorting ? getSortedRowModel() : undefined,
  })

  return (
    <Card className="flex flex-col gap-5 min-h-[450px]">
      <div className="flex items-center justify-between border-b border-border pb-4">
        <h3 className="font-extrabold text-gray-900 text-md flex items-center gap-2">
          {icon}
          {title}
        </h3>
        {addButton}
      </div>

      {isEmpty ? (
        <div className="flex-1 flex flex-col items-center justify-center gap-2 py-16">
          {emptyIcon}
          <span className="text-xs text-muted">{emptyText}</span>
        </div>
      ) : (
        <div className="overflow-x-auto w-full">
          <Table className="w-full text-left border-collapse">
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id} className="border-b border-border bg-gray-50/80">
                  {headerGroup.headers.map((header) => (
                    <TableHead
                      key={header.id}
                      className={`py-3.5 px-4 text-xs font-black text-gray-400 uppercase tracking-wider${header.column.getCanSort() ? " cursor-pointer select-none" : ""}`}
                      onClick={header.column.getToggleSortingHandler()}
                    >
                      <div className="flex items-center gap-1">
                        {flexRender(header.column.columnDef.header, header.getContext())}
                        {header.column.getIsSorted() === "asc" && <ChevronUp className="w-3 h-3" />}
                        {header.column.getIsSorted() === "desc" && <ChevronDown className="w-3 h-3" />}
                      </div>
                    </TableHead>
                  ))}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {table.getRowModel().rows.map((row) => (
                <TableRow key={row.id} className="border-b border-border/60 hover:bg-gray-50 transition-colors duration-300">
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id} className="py-4 px-4 align-top">
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </Card>
  )
}
