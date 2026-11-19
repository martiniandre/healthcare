import { useReactTable, getCoreRowModel, flexRender, type ColumnDef, type SortingState, getSortedRowModel } from "@tanstack/react-table"
import { useState, type ReactNode, type ReactElement } from "react"
import { ChevronUp, ChevronDown, Rows3, Rows4 } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Button } from "../ui/button"
import { useLayoutStore } from "../../store/layout_store"
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card"
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip"
import { cn } from "@/shared/utils/cn"
import { actionsCellClassName, actionsHeaderClassName, actionsHeaderContentClassName, isActionsColumn } from "./columnAlignment"
import { Table, TableHeader, TableBody, TableHead, TableRow, TableCell } from "../ui/table"

interface ClinicalTableProps<T> {
  title: string
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
      <CardHeader className="flex flex-col gap-3 border-b sm:flex-row sm:items-center sm:justify-between">
        <CardTitle className="flex min-w-0 items-center gap-2 text-base font-semibold [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:text-muted-foreground">
          <span className="min-w-0">{title}</span>
        </CardTitle>
        <div className="flex flex-wrap items-center gap-2 sm:shrink-0 sm:justify-end">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="icon-sm"
                onClick={toggleTableDensity}
                aria-label={densityLabel}
                aria-pressed={isCompact}
              >
                {isCompact ? <Rows3 /> : <Rows4 />}
              </Button>
            </TooltipTrigger>
            <TooltipContent>{densityLabel}</TooltipContent>
          </Tooltip>
          {addButton}
        </div>
      </CardHeader>

      {isEmpty ? (
        <CardContent className="flex flex-1 flex-col items-center justify-center gap-2 py-16">
          {emptyIcon}
          <span className="text-xs text-muted-foreground">{emptyText}</span>
        </CardContent>
      ) : (
        <CardContent className="w-full px-0">
          <Table density={tableDensity} containerClassName="max-h-[560px] overflow-y-auto">
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id} className="hover:bg-transparent">
                  {headerGroup.headers.map((header) => (
                    <TableHead
                      key={header.id}
                      className={cn("sticky top-0 z-10 bg-card px-4 first:pl-6 last:pr-6", header.column.getCanSort() && "cursor-pointer select-none hover:text-foreground", isActionsColumn(header.column.id) && actionsHeaderClassName)}
                      onClick={header.column.getToggleSortingHandler()}
                      aria-sort={header.column.getIsSorted() === "asc" ? "ascending" : header.column.getIsSorted() === "desc" ? "descending" : undefined}
                    >
                      <div className={cn("flex items-center gap-1", isActionsColumn(header.column.id) && actionsHeaderContentClassName)}>
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
                    <TableCell key={cell.id} className={cn("px-4 align-middle first:pl-6 last:pr-6", isCompact ? "py-1.5" : "py-4", isActionsColumn(cell.column.id) && actionsCellClassName)}>
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
