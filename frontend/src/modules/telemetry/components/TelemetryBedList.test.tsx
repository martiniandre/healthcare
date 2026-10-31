import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import { BedStatus, CardiacCondition } from "../../../shared/types"
import { TelemetryBedList } from "./TelemetryBedList"

const buildBed = (id: string, bedNumber: string, status: BedStatus) => ({
  id,
  bedNumber,
  patientName: `Patient ${id}`,
  status,
  condition: CardiacCondition.Normal,
  bpm: 80,
  spo2: 97,
})

describe("TelemetryBedList", () => {
  const beds = [
    buildBed("a", "Bed 1", BedStatus.Normal),
    buildBed("b", "Bed 2", BedStatus.Danger),
    buildBed("c", "Bed 3", BedStatus.Warning),
  ]

  it("should list critical beds first, then attention, then normal", () => {
    render(<TelemetryBedList beds={beds} selectedBedId={null} setSelectedBedId={() => {}} />)

    const bedButtons = screen.getAllByRole("button")
    expect(bedButtons[0]).toHaveTextContent("Bed 2")
    expect(bedButtons[1]).toHaveTextContent("Bed 3")
    expect(bedButtons[2]).toHaveTextContent("Bed 1")
  })

  it("should always show a text status and animate only the critical bed", () => {
    render(<TelemetryBedList beds={beds} selectedBedId={null} setSelectedBedId={() => {}} />)

    const criticalBadge = screen.getByText("Critical")
    expect(criticalBadge.className).toContain("animate-pulse-glow")
    expect(criticalBadge.className).toContain("motion-reduce:animate-none")
    expect(screen.getByText("Attention").className).not.toContain("animate-pulse-glow")
    expect(screen.getByText("Stable").className).not.toContain("animate-pulse-glow")
  })

  it("should render a trend sparkline when history is available", () => {
    render(
      <TelemetryBedList
        beds={beds}
        selectedBedId={null}
        setSelectedBedId={() => {}}
        history={{
          b: [
            { bpm: 70, spo2: 98 },
            { bpm: 90, spo2: 95 },
          ],
        }}
      />
    )

    expect(screen.getAllByRole("img")).toHaveLength(2)
  })
})
