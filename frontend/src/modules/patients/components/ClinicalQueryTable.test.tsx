import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { ShieldAlert } from "lucide-react"
import type { ColumnDef } from "@tanstack/react-table"
import { ClinicalQueryTable } from "./ClinicalQueryTable"

interface Row {
  name: string
}

const columns: ColumnDef<Row>[] = [{ id: "name", header: "Name", accessorKey: "name" }]

const renderTable = (overrides: Partial<React.ComponentProps<typeof ClinicalQueryTable<Row>>> = {}) =>
  render(
    <ClinicalQueryTable<Row>
      title="Allergies"
      icon={<ShieldAlert />}
      columns={columns}
      data={[{ name: "Penicillin" }]}
      isLoading={false}
      isError={false}
      onRetry={vi.fn()}
      emptyText="Nothing here"
      emptyIcon={ShieldAlert}
      {...overrides}
    />
  )

describe("ClinicalQueryTable", () => {
  it("should render the rows when data is available", () => {
    renderTable()
    expect(screen.getByText("Penicillin")).toBeInTheDocument()
  })

  it("should render skeleton rows while loading", () => {
    renderTable({ isLoading: true, data: [] })
    expect(screen.getByRole("status")).toHaveAttribute("aria-busy", "true")
    expect(screen.queryByText("Nothing here")).toBeNull()
  })

  it("should render the empty state when there is no data", () => {
    renderTable({ data: [] })
    expect(screen.getByText("Nothing here")).toBeInTheDocument()
  })

  it("should render a destructive alert with a retry button on error", () => {
    const onRetry = vi.fn()
    renderTable({ isError: true, data: [], onRetry })
    expect(screen.getByRole("alert")).toBeInTheDocument()
    fireEvent.click(screen.getByRole("button", { name: "Try again" }))
    expect(onRetry).toHaveBeenCalledTimes(1)
  })
})
