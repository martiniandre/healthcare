import { useReactTable, getCoreRowModel, flexRender, type ColumnDef, type SortingState, getSortedRowModel } from "@tanstack/react-table"
import { useState, type ReactNode, type ReactElement } from "react"
import { ChevronUp, ChevronDown, Rows3, Rows4 } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Button } from "../ui/button"
import { useLayoutStore } from "../../store/layout_store"
import { Card, CardAction, CardContent, CardHeader, CardTitle } from "../ui/card"
import { cn } from "@/shared/utils/cn"
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from "../ui/table"

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
  const { t } = useTranslation("common")
  const tableDensity = useLayoutStore((state) => state.tableDensity)
  const toggleTableDensity = useLayoutStore((state) => state.toggleTableDensity)
  const isCompact = tableDensity === "compact"
  const densityLabel = isCompact ? t("tableDensity.comfortable") : t("tableDensity.compact")

  const table = useReactTable({
    data,
    columns,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: enableSorting ? getSortedRowModel() : undefined,
  })

  return (
    <Card className="min-h-[450px] gap-5">
      <CardHeader className="border-b">
        <CardTitle className="flex items-center gap-2 text-base font-extrabold">
          {icon}
          {title}
        </CardTitle>
        <CardAction className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={toggleTableDensity}
            aria-label={densityLabel}
            aria-pressed={isCompact}
            title={densityLabel}
          >
            {isCompact ? <Rows3 /> : <Rows4 />}
          </Button>
          {addButton}
        </CardAction>
      </CardHeader>

      {isEmpty ? (
        <CardContent className="flex flex-1 flex-col items-center justify-center gap-2 py-16">
          {emptyIcon}
          <span className="text-xs text-muted-foreground">{emptyText}</span>
        </CardContent>
      ) : (
        <CardContent className="w-full">
          <Table density={tableDensity} containerClassName="max-h-[560px] overflow-y-auto">
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id} className="hover:bg-transparent">
                  {headerGroup.headers.map((header) => (
                    <TableHead
                      key={header.id}
                      className={cn("sticky top-0 z-10 bg-card px-4 text-xs font-black uppercase tracking-wider text-muted-foreground", header.column.getCanSort() && "cursor-pointer select-none")}
                      onClick={header.column.getToggleSortingHandler()}
                      aria-sort={header.column.getIsSorted() === "asc" ? "ascending" : header.column.getIsSorted() === "desc" ? "descending" : undefined}
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
                <TableRow key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id} className={cn("px-4 align-top", isCompact ? "py-1.5" : "py-4")}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      )}
    </Card>
  )
}
