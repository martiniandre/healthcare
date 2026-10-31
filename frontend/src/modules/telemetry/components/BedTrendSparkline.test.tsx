import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import { BedTrendSparkline } from "./BedTrendSparkline"
import { buildSparklinePoints } from "../utils/bed_helpers"

describe("BedTrendSparkline", () => {
  it("should expose an accessible label and a polyline", () => {
    const { container } = render(<BedTrendSparkline values={[60, 70, 65]} metricLabel="Heart rate" unit=" BPM" />)

    const sparkline = screen.getByRole("img")
    expect(sparkline.getAttribute("aria-label")).toContain("Heart rate")
    expect(container.querySelector("polyline")).not.toBeNull()
  })

  it("should render nothing with fewer than two samples", () => {
    const { container } = render(<BedTrendSparkline values={[60]} metricLabel="Heart rate" unit=" BPM" />)

    expect(container.querySelector("svg")).toBeNull()
  })

  it("should build one point per value", () => {
    expect(buildSparklinePoints([1, 2, 3]).split(" ")).toHaveLength(3)
    expect(buildSparklinePoints([])).toBe("")
  })
})
