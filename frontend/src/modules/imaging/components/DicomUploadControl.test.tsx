import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { DicomUploadControl } from "./DicomUploadControl"
import { DICOM_MAX_ARCHIVE_BYTES, DICOM_MAX_FILE_BYTES } from "../utils/dicom_file_validation"

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

const ZIP_HEADER_BYTES = [0x50, 0x4b, 0x03, 0x04, 0x14, 0x00]

const selectFile = (file: File) => {
  fireEvent.change(screen.getByTestId("dicom-file-input"), { target: { files: [file] } })
}

const buildZipFile = (fileName: string) => new File([new Uint8Array(ZIP_HEADER_BYTES)], fileName)

const buildFileWithSize = (fileName: string, sizeInBytes: number) => {
  const file = new File(["dicom"], fileName)
  Object.defineProperty(file, "size", { value: sizeInBytes })
  return file
}

const rejectWithCode = (code: string, error = "server text") => {
  mutationState.mutateAsync.mockRejectedValueOnce({ response: { data: { code, error } } })
}

describe("DicomUploadControl", () => {
  beforeEach(() => {
    mutationState.mutateAsync.mockReset()
  })

  it("should show the accepted formats, the one exam rule and the limits", () => {
    render(<DicomUploadControl patientFhirId="patient-1" />)

    expect(screen.getByText(/DICOM \(\.dcm\), up to 2 GB/)).toBeInTheDocument()
    expect(
      screen.getByText(/One exam: a single \.dcm file or a \.zip with all its files \(ZIP up to 2 GB and 2000 files\)/)
    ).toBeInTheDocument()
    expect(screen.getByTestId("dicom-file-input")).toHaveAttribute("accept", ".dcm,.zip")
  })

  it("should reject files that are neither dcm nor zip without uploading", async () => {
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["x"], "photo.png"))

    expect(await screen.findByRole("alert")).toHaveTextContent("is not a DICOM or ZIP file")
    expect(mutationState.mutateAsync).not.toHaveBeenCalled()
  })

  it("should reject a rar archive without uploading", async () => {
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["Rar!"], "exam.rar"))

    expect(await screen.findByRole("alert")).toHaveTextContent("\"exam.rar\" is not a DICOM or ZIP file")
    expect(mutationState.mutateAsync).not.toHaveBeenCalled()
  })

  it("should reject empty dcm files", async () => {
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File([""], "empty.dcm"))

    expect(await screen.findByRole("alert")).toHaveTextContent("is empty")
  })

  it("should reject a zip whose content is not a ZIP archive", async () => {
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["plain text"], "exam.zip"))

    expect(await screen.findByRole("alert")).toHaveTextContent("is not a valid ZIP file")
    expect(mutationState.mutateAsync).not.toHaveBeenCalled()
  })

  it("should reject a dcm file above the single file limit", async () => {
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(buildFileWithSize("huge.dcm", DICOM_MAX_FILE_BYTES + 1))

    expect(await screen.findByRole("alert")).toHaveTextContent("above the 2 GB limit")
    expect(mutationState.mutateAsync).not.toHaveBeenCalled()
  })

  it("should reject a zip above the archive limit", async () => {
    render(<DicomUploadControl patientFhirId="patient-1" />)

    const oversizedArchive = buildFileWithSize("exam.zip", DICOM_MAX_ARCHIVE_BYTES + 1)
    selectFile(oversizedArchive)

    expect(await screen.findByRole("alert")).toHaveTextContent("above the ZIP limit of 2 GB")
    expect(mutationState.mutateAsync).not.toHaveBeenCalled()
  })

  it("should upload a valid dcm file and show the new study outcome", async () => {
    mutationState.mutateAsync.mockResolvedValue({ id: "study-2", studyCreated: true, instancesAccepted: 1 })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(new File(["dicom"], "mri.dcm"))

    await waitFor(() => expect(mutationState.mutateAsync).toHaveBeenCalledTimes(1))
    const sentPayload = mutationState.mutateAsync.mock.calls[0][0] as Record<string, unknown>
    expect(sentPayload).toMatchObject({ patientFhirId: "patient-1" })
    expect(sentPayload.title).toBeUndefined()
    expect(sentPayload.modality).toBeUndefined()
    expect(await screen.findByText("New study created")).toBeInTheDocument()
    expect(screen.getByText("1 image accepted")).toBeInTheDocument()
  })

  it("should accept a zip archive and send it with its original name", async () => {
    mutationState.mutateAsync.mockResolvedValue({ id: "study-5", studyCreated: true, instancesAccepted: 3 })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(buildZipFile("exam-bundle.zip"))

    await waitFor(() => expect(mutationState.mutateAsync).toHaveBeenCalledTimes(1))
    const sentPayload = mutationState.mutateAsync.mock.calls[0][0] as UploadPayload
    expect(sentPayload.dicomBlob.name).toBe("exam-bundle.zip")
    expect(await screen.findByText("3 images accepted")).toBeInTheDocument()
  })

  it("should say files were added when the instances went into an existing study", async () => {
    mutationState.mutateAsync.mockResolvedValue({ id: "study-2", studyCreated: false, instancesAccepted: 2 })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(buildZipFile("second-batch.zip"))

    expect(await screen.findByText("Files added to the existing study")).toBeInTheDocument()
    expect(screen.queryByText("New study created")).not.toBeInTheDocument()
  })

  it("should list the rejected files with a translated reason and never the raw server text", async () => {
    mutationState.mutateAsync.mockResolvedValue({
      id: "study-2",
      studyCreated: false,
      instancesAccepted: 1,
      instancesRejected: [
        { fileName: "IMG-0002.dcm", code: "DICOM_DUPLICATE_INSTANCE", reason: "RAW SERVER TEXT ONE" },
        { fileName: "IMG-0003.dcm", code: "DICOM_UNSUPPORTED_TRANSFER_SYNTAX", reason: "RAW SERVER TEXT TWO" },
        { fileName: "IMG-0004.dcm", code: "DICOM_MISSING_IDENTIFIERS", reason: "RAW SERVER TEXT THREE" },
        { fileName: "IMG-0005.dcm", code: "BRAND_NEW_CODE", reason: "RAW SERVER TEXT FOUR" },
      ],
    })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(buildZipFile("partial.zip"))

    expect(await screen.findByText("4 files were not added")).toBeInTheDocument()
    expect(screen.getByText("IMG-0002.dcm")).toBeInTheDocument()
    expect(screen.getByText("Already stored in this study")).toBeInTheDocument()
    expect(screen.getByText("Compressed format not supported")).toBeInTheDocument()
    expect(screen.getByText("Missing patient or study identifiers")).toBeInTheDocument()
    expect(screen.getByText("The file could not be read")).toBeInTheDocument()
    expect(screen.queryByText(/RAW SERVER TEXT/)).not.toBeInTheDocument()
  })

  it("should not show the rejected list when every file was accepted", async () => {
    mutationState.mutateAsync.mockResolvedValue({ id: "study-2", studyCreated: true, instancesAccepted: 2, instancesRejected: [] })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(buildZipFile("clean.zip"))

    expect(await screen.findByText("New study created")).toBeInTheDocument()
    expect(screen.queryByText(/were not added|was not added/)).not.toBeInTheDocument()
  })

  it("should treat a success without accepted instances as an error", async () => {
    mutationState.mutateAsync.mockResolvedValue({ id: "study-2", studyCreated: false, instancesAccepted: 0 })
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(buildZipFile("empty-of-dicom.zip"))

    expect(
      await screen.findByText("No file in the ZIP is a valid DICOM of this study.")
    ).toBeInTheDocument()
    expect(screen.queryByText("New study created")).not.toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Try again" })).toBeInTheDocument()
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

  it.each([
    ["DICOM_PATIENT_MISMATCH", "This DICOM belongs to a different patient than the one selected.", false],
    ["INVALID_DICOM", "This is not a valid or supported DICOM file.", true],
    ["DICOM_DUPLICATE_INSTANCE", "This DICOM image was already uploaded.", true],
    ["ZIP_INVALID", "This file is not a valid ZIP archive.", true],
    ["ZIP_UNSAFE_ENTRY", "The ZIP contains a file with an unsafe name.", true],
    ["ZIP_LIMITS_EXCEEDED", "The ZIP is too large or has too many files.", true],
    ["DICOM_NO_VALID_INSTANCES", "No file in the ZIP is a valid DICOM of this study.", true],
    ["SOMETHING_UNKNOWN", "The DICOM file could not be uploaded.", true],
  ])(
    "should explain the %s error and show retry only when allowed",
    async (errorCode, expectedTitle, isRetryVisible) => {
      rejectWithCode(errorCode)
      render(<DicomUploadControl patientFhirId="patient-1" />)

      selectFile(new File(["dicom"], "mri.dcm"))

      expect(await screen.findByText(expectedTitle)).toBeInTheDocument()
      const retryButton = screen.queryByRole("button", { name: "Try again" })
      if (isRetryVisible) {
        expect(retryButton).toBeInTheDocument()
      } else {
        expect(retryButton).not.toBeInTheDocument()
      }
    }
  )

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

  it("should show the zip limits in the archive limit message", async () => {
    rejectWithCode("ZIP_LIMITS_EXCEEDED")
    render(<DicomUploadControl patientFhirId="patient-1" />)

    selectFile(buildZipFile("big.zip"))

    expect(await screen.findByText("Each ZIP can hold up to 2 GB and 2000 files.", { exact: false })).toBeInTheDocument()
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
