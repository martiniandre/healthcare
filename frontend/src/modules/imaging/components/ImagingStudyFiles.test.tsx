import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { ImagingStudyFiles } from "./ImagingStudyFiles"
import type { ImagingStudyInstance } from "../types"

const instancesQueryState = vi.hoisted(() => ({
  current: {} as Record<string, unknown>,
  refetch: vi.fn(),
  calls: [] as string[],
}))

vi.mock("../queries", () => ({
  useImagingStudyInstancesQuery: (studyId: string) => {
    instancesQueryState.calls.push(studyId)
    return { ...instancesQueryState.current, refetch: instancesQueryState.refetch }
  },
}))

const buildInstance = (overrides: Partial<ImagingStudyInstance> = {}): ImagingStudyInstance => ({
  id: "instance-1",
  sopInstanceUid: "1.2.3.4.5",
  seriesInstanceUid: "1.2.3.4",
  instanceNumber: 1,
  fileName: "IMG-0001.dcm",
  downloadUrl: "https://storage.example.test/IMG-0001.dcm?signature=abc",
  expiresAt: "2026-10-31T12:00:00Z",
  ...overrides,
})

describe("ImagingStudyFiles", () => {
  beforeEach(() => {
    instancesQueryState.calls = []
    instancesQueryState.refetch.mockReset()
    instancesQueryState.current = { data: [], isPending: false, isError: false }
  })

  it("should not load the instances until the section is opened", () => {
    render(<ImagingStudyFiles studyId="study-1" />)

    expect(screen.getByRole("button", { name: "Show files" })).toHaveAttribute("aria-expanded", "false")
    expect(instancesQueryState.calls).toHaveLength(0)
  })

  it("should show the file count and a download button for each instance", () => {
    instancesQueryState.current = {
      data: [
        buildInstance(),
        buildInstance({
          id: "instance-2",
          fileName: "IMG-0002.dcm",
          downloadUrl: "https://storage.example.test/IMG-0002.dcm?signature=def",
        }),
      ],
      isPending: false,
      isError: false,
    }
    render(<ImagingStudyFiles studyId="study-1" />)

    fireEvent.click(screen.getByRole("button", { name: "Show files" }))

    expect(instancesQueryState.calls).toEqual(["study-1"])
    expect(screen.getByRole("button", { name: "Hide files" })).toHaveAttribute("aria-expanded", "true")
    expect(screen.getByText("2 files")).toBeInTheDocument()
    expect(screen.getByText("IMG-0001.dcm")).toBeInTheDocument()
    const firstDownload = screen.getByRole("link", { name: "Download IMG-0001.dcm" })
    expect(firstDownload).toHaveAttribute("href", "https://storage.example.test/IMG-0001.dcm?signature=abc")
    expect(firstDownload).toHaveAttribute("download", "IMG-0001.dcm")
    expect(screen.getByRole("link", { name: "Download IMG-0002.dcm" })).toHaveAttribute(
      "href",
      "https://storage.example.test/IMG-0002.dcm?signature=def"
    )
  })

  it("should show a single file with its own download", () => {
    instancesQueryState.current = { data: [buildInstance()], isPending: false, isError: false }
    render(<ImagingStudyFiles studyId="study-1" />)

    fireEvent.click(screen.getByRole("button", { name: "Show files" }))

    expect(screen.getByText("1 file")).toBeInTheDocument()
    expect(screen.getByRole("link", { name: "Download IMG-0001.dcm" })).toBeInTheDocument()
  })

  it("should show the loading state while the instances are fetched", () => {
    instancesQueryState.current = { data: undefined, isPending: true, isError: false }
    render(<ImagingStudyFiles studyId="study-1" />)

    fireEvent.click(screen.getByRole("button", { name: "Show files" }))

    expect(screen.getByText("Loading files...")).toBeInTheDocument()
  })

  it("should offer a retry when the instances cannot be loaded", () => {
    instancesQueryState.current = { data: undefined, isPending: false, isError: true }
    render(<ImagingStudyFiles studyId="study-1" />)

    fireEvent.click(screen.getByRole("button", { name: "Show files" }))
    fireEvent.click(screen.getByRole("button", { name: "Try again" }))

    expect(screen.getByText("The files of this study could not be loaded.")).toBeInTheDocument()
    expect(instancesQueryState.refetch).toHaveBeenCalledTimes(1)
  })

  it("should show an empty message when the study has no files", () => {
    render(<ImagingStudyFiles studyId="study-1" />)

    fireEvent.click(screen.getByRole("button", { name: "Show files" }))

    expect(screen.getByText("No files were found for this study.")).toBeInTheDocument()
  })
})
