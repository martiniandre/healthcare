import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import { ImagingStudyDetails } from "./ImagingStudyDetails"
import { DicomModality, ImagingStudyStatus } from "../../../shared/types"
import type { ImagingStudy } from "../types"

const buildStudy = (overrides: Partial<ImagingStudy> = {}): ImagingStudy => ({
  id: "study-1",
  patient_fhir_id: "patient-1",
  title: "Chest CT",
  modality: DicomModality.CT,
  study_instance_uid: "1.2.3",
  status: ImagingStudyStatus.Available,
  created_at: "2026-01-01T00:00:00Z",
  ...overrides,
})

describe("ImagingStudyDetails", () => {
  it("should show the metadata read from the DICOM file when present", () => {
    render(
      <ImagingStudyDetails
        study={buildStudy({
          studyDate: "20260315",
          accessionNumber: "ACC-1",
          studyInstanceUid: "1.2.840.99",
          seriesInstanceUid: "1.2.840.99.1",
          sopInstanceUid: "1.2.840.99.1.1",
        })}
      />
    )

    expect(screen.getByText("Study Date")).toBeInTheDocument()
    expect(screen.getByText(/2026/)).toBeInTheDocument()
    expect(screen.getByText("ACC-1")).toBeInTheDocument()
    expect(screen.getByText("1.2.840.99")).toBeInTheDocument()
    expect(screen.getByText("1.2.840.99.1")).toBeInTheDocument()
    expect(screen.getByText("1.2.840.99.1.1")).toBeInTheDocument()
  })

  it("should fall back to the stored study UID and hide missing metadata", () => {
    render(<ImagingStudyDetails study={buildStudy()} />)

    expect(screen.getByText("Study UID")).toBeInTheDocument()
    expect(screen.getByText("1.2.3")).toBeInTheDocument()
    expect(screen.queryByText("Study Date")).not.toBeInTheDocument()
    expect(screen.queryByText("Accession Number")).not.toBeInTheDocument()
    expect(screen.queryByText("Series UID")).not.toBeInTheDocument()
  })
})
