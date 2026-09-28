import { describe, it, expect, vi } from "vitest"
import { render, screen } from "@testing-library/react"
import { Activity, Users } from "lucide-react"
import { MetricStrip, type MetricStripItem } from "./MetricStrip"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}))

const buildItems = (): MetricStripItem[] => [
  {
    label: "Consultations today",
    value: "128",
    icon: Users,
    tone: "primary",
    meta: { text: "+4%", tone: "text-gray-700" },
  },
  {
    label: "Occupancy rate",
    value: "72.4",
    unit: "%",
    icon: Activity,
    tone: "elevated",
    meta: { text: "Elevated", tone: "text-warning" },
  },
]

describe("MetricStrip", () => {
  it("renders one card per item", () => {
    render(<MetricStrip items={buildItems()} />)

    expect(screen.getByText("Consultations today")).toBeInTheDocument()
    expect(screen.getByText("Occupancy rate")).toBeInTheDocument()
  })

  it("renders the value and the unit as distinct nodes", () => {
    render(<MetricStrip items={buildItems()} />)

    expect(screen.getByText("128")).toBeInTheDocument()
    expect(screen.getByText("72.4")).toBeInTheDocument()
    expect(screen.getByText("%")).toBeInTheDocument()
  })

  it("omits the unit node when no unit is provided", () => {
    render(<MetricStrip items={buildItems()} />)
    expect(screen.queryByText("%")).toBeInTheDocument()
    expect(screen.getAllByText("%")).toHaveLength(1)
  })

  it("renders the meta caption for every item", () => {
    render(<MetricStrip items={buildItems()} />)

    expect(screen.getByText("+4%")).toBeInTheDocument()
    expect(screen.getByText("Elevated")).toBeInTheDocument()
  })

  it("applies the tone bar class for each item tone", () => {
    const { container } = render(<MetricStrip items={buildItems()} />)

    const toneBars = container.querySelectorAll("span[aria-hidden='true'].bg-primary")
    expect(toneBars.length).toBe(1)
    expect(container.querySelectorAll("span[aria-hidden='true'].bg-warning").length).toBe(1)
  })

  it("renders nothing but the grid when the item list is empty", () => {
    render(<MetricStrip items={[]} />)
    expect(screen.queryByRole("heading")).not.toBeInTheDocument()
  })
})
