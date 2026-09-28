import { describe, it, expect, vi, afterEach } from "vitest"
import { render, screen, within, fireEvent, cleanup } from "@testing-library/react"
import { StatsEpidemiologyTable } from "./StatsEpidemiologyTable"
import type { PathologyData } from "../types"

afterEach(() => {
  cleanup()
})

const translations: Record<string, string> = {
  "epidemiology.title": "Epidemiology",
  "epidemiology.subtitle": "Active cases",
  "epidemiology.exportButton": "Export",
  "epidemiology.activeCasesTotalLabel": "Cases",
  "epidemiology.activeCasesUnit": "{{count}} active",
  "epidemiology.table.code": "Code",
  "epidemiology.table.description": "Description",
  "epidemiology.table.category": "Category",
  "epidemiology.table.activeCases": "Active",
  "epidemiology.table.trend": "Trend",
  "epidemiology.table.stable": "Stable",
  "epidemiology.table.proportion": "Proportion",
  "epidemiology.trend.severity.critical": "Critical",
  "epidemiology.trend.severity.elevated": "Under watch",
  "epidemiology.trend.severity.neutral": "Within expected range",
  "epidemiology.trend.severity.stable": "Declining",
  "empty.epidemiology": "No cases",
  "empty.epidemiologyDesc": "Nothing to show",
}

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) => {
      const template = translations[key]
      if (template === undefined) return key
      if (options === undefined) return template
      return template.replace(/{{(\w+)}}/g, (_match, token: string) => String(options[token]))
    },
  }),
}))

const basePathologies: PathologyData[] = [
  {
    code: "J45.9",
    description: "Unspecified asthma",
    category: "Respiratory",
    activeCases: 44,
    trend: "+5%",
  },
  {
    code: "I10",
    description: "Essential hypertension",
    category: "Cardiovascular",
    activeCases: 119,
    trend: "stable",
  },
  {
    code: "E11.9",
    description: "Type 2 diabetes",
    category: "Endocrine",
    activeCases: 85,
    trend: "-3%",
  },
]

const findTrendBadge = (code: string): HTMLElement => {
  const row = screen.getByRole("row", { name: new RegExp(code.replace(".", "\\.")) })
  const cells = within(row).getAllByRole("cell")
  const trendCell = cells.find((cell) => cell.querySelector("svg") !== null)
  return trendCell?.querySelector("span") as HTMLElement
}

describe("StatsEpidemiologyTable", () => {
  it("renders one row per pathology with its code and description", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)

    expect(screen.getByText("J45.9")).toBeInTheDocument()
    expect(screen.getByText("Unspecified asthma")).toBeInTheDocument()
    expect(screen.getByText("I10")).toBeInTheDocument()
    expect(screen.getByText("E11.9")).toBeInTheDocument()
  })

  it("colours a rising trend according to its magnitude", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)

    const trendBadge = findTrendBadge("J45.9")
    expect(trendBadge).toHaveTextContent("+5%")
    expect(trendBadge.className).toContain("text-warning")
  })

  it("escalates a large rise to the critical tone", () => {
    render(
      <StatsEpidemiologyTable
        pathologies={[{ ...basePathologies[0], trend: "+14%" }]}
      />
    )

    const trendBadge = findTrendBadge("J45.9")
    expect(trendBadge.className).toContain("text-danger")
  })

  it("keeps a marginal rise neutral instead of alarming", () => {
    render(
      <StatsEpidemiologyTable
        pathologies={[{ ...basePathologies[0], trend: "+1%" }]}
      />
    )

    const trendBadge = findTrendBadge("J45.9")
    expect(trendBadge.className).toContain("text-muted-foreground")
    expect(trendBadge.className).not.toContain("text-danger")
  })

  it("exposes the severity to assistive technology alongside the colour", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)

    const trendBadge = findTrendBadge("J45.9")
    expect(trendBadge).toHaveTextContent("Under watch")
  })

  it("colours a falling trend as stable", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)

    const trendBadge = findTrendBadge("E11.9")
    expect(trendBadge).toHaveTextContent("-3%")
    expect(trendBadge.className).toContain("text-success")
  })

  it("localizes the stable label and keeps it neutral", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)

    const trendBadge = findTrendBadge("I10")
    expect(trendBadge).toHaveTextContent("Stable")
    expect(trendBadge.className).toContain("text-muted-foreground")
  })

  it("does not colour by ICD code when the trend direction disagrees with the code", () => {
    render(
      <StatsEpidemiologyTable
        pathologies={[
          { ...basePathologies[0], code: "E11.9", trend: "-4%" },
          { ...basePathologies[2], code: "J45.9", trend: "+12%" },
        ]}
      />
    )

    const fallingDiabetesBadge = findTrendBadge("E11.9")
    expect(fallingDiabetesBadge.className).toContain("text-success")
    expect(fallingDiabetesBadge.className).not.toContain("text-danger")

    const risingAsthmaBadge = findTrendBadge("J45.9")
    expect(risingAsthmaBadge.className).toContain("text-danger")
  })

  it("sorts by descending active cases by default", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)

    const orderedCodes = screen
      .getAllByRole("row")
      .slice(1)
      .map((row) => within(row).getAllByRole("cell")[0].textContent)

    expect(orderedCodes).toEqual(["I10", "E11.9", "J45.9"])
  })

  it("reverses the active case order when the same header is clicked again", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)

    fireEvent.click(screen.getByRole("button", { name: /Active/ }))

    const orderedCodes = screen
      .getAllByRole("row")
      .slice(1)
      .map((row) => within(row).getAllByRole("cell")[0].textContent)

    expect(orderedCodes).toEqual(["J45.9", "E11.9", "I10"])
  })

  it("sorts by ICD code when the code header is clicked", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)

    fireEvent.click(screen.getByRole("button", { name: /Code/ }))

    const orderedCodes = screen
      .getAllByRole("row")
      .slice(1)
      .map((row) => within(row).getAllByRole("cell")[0].textContent)

    expect(orderedCodes).toEqual(["E11.9", "I10", "J45.9"])
  })

  it("sorts by trend severity so the worst trends surface first", () => {
    render(
      <StatsEpidemiologyTable
        pathologies={[
          { ...basePathologies[0], code: "A00.0", trend: "-2%" },
          { ...basePathologies[0], code: "B00.0", trend: "+1%" },
          { ...basePathologies[0], code: "C00.0", trend: "+15%" },
        ]}
      />
    )

    fireEvent.click(screen.getByRole("button", { name: /Trend/ }))

    const orderedCodes = screen
      .getAllByRole("row")
      .slice(1)
      .map((row) => within(row).getAllByRole("cell")[0].textContent)

    expect(orderedCodes).toEqual(["C00.0", "B00.0", "A00.0"])
  })

  it("marks the active sort column for assistive technology", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)

    expect(screen.getByRole("columnheader", { name: /Active/ })).toHaveAttribute(
      "aria-sort",
      "descending"
    )
    expect(screen.getByRole("columnheader", { name: /Code/ })).toHaveAttribute("aria-sort", "none")
  })

  it("renders a proportion bar relative to the largest case count", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)

    expect(screen.getByText("100%")).toBeInTheDocument()
    expect(screen.getByText("71%")).toBeInTheDocument()
    expect(screen.getByText("37%")).toBeInTheDocument()
  })

  it("shows the total active cases when there is data", () => {
    render(<StatsEpidemiologyTable pathologies={basePathologies} />)
    expect(screen.getByText("248 active")).toBeInTheDocument()
  })

  it("renders the empty state when there are no pathologies", () => {
    render(<StatsEpidemiologyTable pathologies={[]} />)

    expect(screen.getByText("No cases")).toBeInTheDocument()
    expect(screen.queryByText("248 active")).not.toBeInTheDocument()
  })
})
