import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import { findVitalSignDisplay, vitalSignDisplayMetadata } from "./vitalSignDisplay"
import { VitalSignValueDisplay } from "./VitalSignValueDisplay"

describe("vitalSignDisplayMetadata", () => {
  it("should cover every metric produced by the vital signs panel", () => {
    const coveredLoincCodes = vitalSignDisplayMetadata.map((metadata) => metadata.loincCode)
    expect(coveredLoincCodes).toEqual([
      "8867-4",
      "8310-5",
      "8480-6",
      "8462-4",
      "59408-5",
      "9279-1",
      "29463-7",
      "8302-2",
    ])
  })

  it("should provide an icon and a localized label key for each entry", () => {
    for (const metadata of vitalSignDisplayMetadata) {
      expect(metadata.IconComponent).toBeDefined()
      expect(metadata.iconClassName).toContain("bg-")
      expect(metadata.labelKey).toMatch(/^details\.vitalsCard\.metrics\./)
    }
  })

  it("should return no metadata for legacy panel observations so code_display is kept", () => {
    expect(findVitalSignDisplay("85354-9")).toBeUndefined()
  })

  it("should distinguish systolic and diastolic entries by color", () => {
    const systolic = findVitalSignDisplay("8480-6")
    const diastolic = findVitalSignDisplay("8462-4")
    expect(systolic?.iconClassName).not.toEqual(diastolic?.iconClassName)
  })
})

describe("VitalSignValueDisplay", () => {
  it("should render a muted N/A badge when the observation was not performed", () => {
    render(<VitalSignValueDisplay notPerformed valueQuantity={0} valueUnit="" />)
    expect(screen.getByText("N/A")).toBeDefined()
    expect(screen.queryByText("0")).toBeNull()
  })

  it("should render the measured value with its unit otherwise", () => {
    render(<VitalSignValueDisplay valueQuantity={72} valueUnit="bpm" />)
    expect(screen.getByText("72")).toBeDefined()
    expect(screen.getByText("bpm")).toBeDefined()
    expect(screen.queryByText("N/A")).toBeNull()
  })

  it("should apply the custom value class name when provided", () => {
    render(<VitalSignValueDisplay valueQuantity={36.5} valueUnit="°C" valueClassName="text-2xl font-bold" />)
    const valueElement = screen.getByText("36.5")
    expect(valueElement.className).toContain("text-2xl")
  })
})
