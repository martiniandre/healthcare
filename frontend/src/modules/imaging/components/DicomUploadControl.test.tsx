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
    const sentPayload = mutationState.mutateAsync.mock.calls[0][0] as Record<string, unknown>
    expect(sentPayload).toMatchObject({ patientFhirId: "patient-1" })
    expect(sentPayload.title).toBeUndefined()
    expect(sentPayload.modality).toBeUndefined()
    expect(await screen.findByText("DICOM sent. It is being processed.")).toBeInTheDocument()
  })

  it("should summarize what the server read from the file", async () => {
    mutationState.mutateAsync.mockResolvedValue({
      id: "study-2",
      modality: "CT",
      studyDate: "20260315",
      accessionNumber: "ACC-778",
      studyInstanceUid: "1.2.840.113619.2.55.3.604688119.969.1234567890.123",
      seriesInstanceUid: "1.2.840.113619.2.55.3.604688119.969.1234567890.456",
      sopInstanceUid: "1.2.840.113619.2.55.3.604688119.969.1234567890.789",
      warnings: [],
    })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["dicom"], "ct.dcm"))

    const summary = await screen.findByTestId("dicom-upload-summary")
    expect(summary).toHaveTextContent("CT")
    expect(summary).toHaveTextContent("2026")
    expect(summary).toHaveTextContent("ACC-778")
    expect(summary).toHaveTextContent("1.2.840.11…567890.123")
    expect(screen.getByTitle("1.2.840.113619.2.55.3.604688119.969.1234567890.123")).toBeInTheDocument()
    expect(screen.getByText("Series UID")).toBeInTheDocument()
    expect(screen.getByText("Instance UID")).toBeInTheDocument()
  })

  it("should copy the full study UID", async () => {
    const writeText = vi.fn().mockResolvedValue(undefined)
    Object.defineProperty(navigator, "clipboard", { value: { writeText }, configurable: true })
    mutationState.mutateAsync.mockResolvedValue({
      id: "study-2",
      modality: "CT",
      studyInstanceUid: "1.2.840.113619.2.55.3.604688119.969.1234567890.123",
    })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["dicom"], "ct.dcm"))
    fireEvent.click(await screen.findByRole("button", { name: "Copy study UID" }))

    await waitFor(() =>
      expect(writeText).toHaveBeenCalledWith("1.2.840.113619.2.55.3.604688119.969.1234567890.123")
    )
  })

  it("should translate the server warnings and ignore their raw text", async () => {
    mutationState.mutateAsync.mockResolvedValue({
      id: "study-2",
      modality: "MR",
      warnings: [
        { code: "patient_identity_missing", message: "raw one" },
        { code: "anonymized_study", message: "raw two" },
        { code: "duplicate_instance", message: "raw three" },
        { code: "brand_new_code", message: "raw four" },
      ],
    })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["dicom"], "mr.dcm"))

    expect(await screen.findByText(/The file has no patient identification/)).toBeInTheDocument()
    expect(screen.getByText(/appears to be anonymized/)).toBeInTheDocument()
    expect(screen.getByText(/already uploaded before/)).toBeInTheDocument()
    expect(screen.getByText(/The server reported a notice about this file/)).toBeInTheDocument()
    expect(screen.queryByText(/raw one/)).not.toBeInTheDocument()
  })

  it("should explain a patient mismatch without exposing patient data and not offer retry", async () => {
    mutationState.mutateAsync.mockRejectedValueOnce({
      response: { data: { code: "DICOM_PATIENT_MISMATCH", error: "patient Jane Doe does not match" } },
    })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["dicom"], "mri.dcm"))

    expect(
      await screen.findByText("This DICOM belongs to a different patient than the one selected.")
    ).toBeInTheDocument()
    expect(screen.queryByText(/Jane Doe/)).not.toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "Try again" })).not.toBeInTheDocument()
  })

  it("should explain an invalid DICOM with the reason returned by the server", async () => {
    mutationState.mutateAsync.mockRejectedValueOnce({
      response: { data: { code: "INVALID_DICOM", error: "invalid", reason: "missing preamble" } },
    })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["dicom"], "mri.dcm"))

    expect(await screen.findByText("This is not a valid or supported DICOM file.")).toBeInTheDocument()
    expect(screen.getByText("Reason: missing preamble")).toBeInTheDocument()
  })

  it("should explain a duplicate DICOM instance", async () => {
    mutationState.mutateAsync.mockRejectedValueOnce({
      response: { data: { code: "DICOM_DUPLICATE_INSTANCE", error: "duplicate" } },
    })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["dicom"], "mri.dcm"))

    expect(await screen.findByText("This DICOM image was already uploaded.")).toBeInTheDocument()
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
