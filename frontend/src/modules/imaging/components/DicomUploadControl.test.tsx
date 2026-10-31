import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { DicomUploadControl } from "./DicomUploadControl"

interface UploadPayload {
  signal: AbortSignal
  onUploadProgress: (percentage: number) => void
  dicomBlob: File
}

const mutationState = vi.hoisted(() => ({
  mutateAsync: vi.fn(),
}))

vi.mock("../queries", () => ({
  useUploadImagingStudyMutation: () => ({ isPending: false, mutateAsync: mutationState.mutateAsync }),
}))

const selectFile = (file: File) => {
  fireEvent.change(screen.getByTestId("dicom-file-input"), { target: { files: [file] } })
}

describe("DicomUploadControl", () => {
  beforeEach(() => {
    mutationState.mutateAsync.mockReset()
  })

  it("should show the accepted format and size limit", () => {
    render(<DicomUploadControl patientFhirId="patient-1" />)

    expect(screen.getByText(/DICOM \(\.dcm\), up to 2 GB/)).toBeInTheDocument()
  })

  it("should reject files that are not dcm without uploading", () => {
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["x"], "photo.png"))

    expect(screen.getByRole("alert")).toHaveTextContent("is not a DICOM file")
    expect(mutationState.mutateAsync).not.toHaveBeenCalled()
  })

  it("should reject empty dcm files", () => {
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File([""], "empty.dcm"))

    expect(screen.getByRole("alert")).toHaveTextContent("is empty")
  })

  it("should upload a valid file and show the success state", async () => {
    mutationState.mutateAsync.mockResolvedValue({ id: "study-2" })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["dicom"], "mri.dcm"))

    await waitFor(() => expect(mutationState.mutateAsync).toHaveBeenCalledTimes(1))
    expect(mutationState.mutateAsync.mock.calls[0][0]).toMatchObject({ patientFhirId: "patient-1", title: "mri" })
    expect(await screen.findByText("DICOM sent. It is being processed.")).toBeInTheDocument()
  })

  it("should show a determinate progress bar and cancel the upload", async () => {
    mutationState.mutateAsync.mockImplementation(
      (payload: UploadPayload) =>
        new Promise((_resolve, reject) => {
          payload.onUploadProgress(40)
          payload.signal.addEventListener("abort", () => reject(new Error("canceled")))
        })
    )
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["dicom"], "mri.dcm"))

    const progressBar = await screen.findByRole("progressbar")
    await waitFor(() => expect(progressBar).toHaveAttribute("aria-valuenow", "40"))

    fireEvent.click(screen.getByRole("button", { name: "Cancel upload" }))

    expect(await screen.findByText("Upload cancelled.")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Try again" })).toBeInTheDocument()
  })

  it("should keep the file and retry after a failure", async () => {
    mutationState.mutateAsync.mockRejectedValueOnce(new Error("network")).mockResolvedValueOnce({ id: "study-3" })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["dicom"], "mri.dcm"))

    expect(await screen.findByText("The DICOM file could not be uploaded.")).toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: "Try again" }))

    await waitFor(() => expect(mutationState.mutateAsync).toHaveBeenCalledTimes(2))
    const retriedPayload = mutationState.mutateAsync.mock.calls[1][0] as UploadPayload
    expect(retriedPayload.dicomBlob.name).toBe("mri.dcm")
  })
})
