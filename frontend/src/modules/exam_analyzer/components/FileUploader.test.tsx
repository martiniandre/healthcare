import { describe, it, expect, vi, beforeEach, afterEach } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { FileUploader, type UploadOutcome } from "./FileUploader"

class ResizeObserverStub {
  observe() {}
  unobserve() {}
  disconnect() {}
}

const selectFile = (file: File) => {
  fireEvent.change(screen.getByTestId("exam-file-input"), { target: { files: [file] } })
}

const buildPdf = (name = "labs.pdf") => new File(["pdf-content"], name, { type: "application/pdf" })

const renderUploader = (overrides: Partial<React.ComponentProps<typeof FileUploader>> = {}) => {
  const onUpload =
    vi.fn<
      (file: File, consent: boolean, anonymize: boolean, acknowledgeExternalProcessing: boolean) => Promise<UploadOutcome>
    >()
  onUpload.mockResolvedValue({ status: "success" })
  const onCancel = vi.fn()
  render(
    <FileUploader onUpload={onUpload} onCancel={onCancel} isPending={false} uploadProgress={null} {...overrides} />
  )
  return { onUpload, onCancel }
}

describe("FileUploader", () => {
  beforeEach(() => {
    vi.stubGlobal("ResizeObserver", ResizeObserverStub)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it("should show accepted formats and size limit as helper text", () => {
    renderUploader()

    expect(screen.getByText(/PNG, JPG or PDF. Up to 15 MB/)).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Choose file" })).toBeInTheDocument()
  })

  it("should show a specific message for unsupported types", () => {
    renderUploader()

    selectFile(new File(["x"], "study.dcm"))

    expect(screen.getByRole("alert")).toHaveTextContent("is not a supported format")
    expect(screen.getByRole("alert")).toHaveTextContent("study.dcm")
  })

  it("should show a specific message for empty files", () => {
    renderUploader()

    selectFile(new File([""], "blank.pdf"))

    expect(screen.getByRole("alert")).toHaveTextContent("blank.pdf")
    expect(screen.getByRole("alert")).toHaveTextContent("is empty")
  })

  it("should show a specific message for files above the limit", () => {
    renderUploader()
    const oversizedFile = buildPdf("huge.pdf")
    Object.defineProperty(oversizedFile, "size", { value: 16 * 1024 * 1024 })

    selectFile(oversizedFile)

    expect(screen.getByRole("alert")).toHaveTextContent("above the 15 MB limit")
  })

  it("should show the selected file and allow removing it", () => {
    renderUploader()

    selectFile(buildPdf())

    expect(screen.getByText("labs.pdf")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Replace file" })).toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: "Remove selected file" }))

    expect(screen.queryByText("labs.pdf")).not.toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Choose file" })).toBeInTheDocument()
  })

  it("should replace the selected file when another valid file is chosen", () => {
    renderUploader()

    selectFile(buildPdf("first.pdf"))
    selectFile(buildPdf("second.pdf"))

    expect(screen.queryByText("first.pdf")).not.toBeInTheDocument()
    expect(screen.getByText("second.pdf")).toBeInTheDocument()
  })

  it("should list what is missing and keep submit enabled", () => {
    const { onUpload } = renderUploader()

    const submitButton = screen.getByRole("button", { name: "Submit for Analysis" })
    expect(submitButton).toBeEnabled()
    expect(screen.getByText("Choose an exam file (step 1).")).toBeInTheDocument()
    expect(screen.getByText("Confirm the patient consent (step 2).")).toBeInTheDocument()

    fireEvent.click(submitButton)

    expect(screen.getByText("Cannot send yet. Please:")).toBeInTheDocument()
    expect(onUpload).not.toHaveBeenCalled()
  })

  it("should keep anonymization checked by default and explain what it does and does not cover", () => {
    renderUploader()

    expect(screen.getByRole("checkbox", { name: /Remove identifying file data/ })).toBeChecked()
    expect(screen.getByText(/Removes the file metadata and the original file name/)).toBeInTheDocument()
    expect(screen.getByText(/Text printed inside the image is NOT removed/)).toBeInTheDocument()
    expect(screen.queryByRole("checkbox", { name: /processed without anonymization/ })).not.toBeInTheDocument()
  })

  it("should mark consent as required and anonymization as optional", () => {
    renderUploader()

    expect(screen.getByText("Required")).toBeInTheDocument()
    expect(screen.getByText("Optional")).toBeInTheDocument()
    expect(screen.getByRole("checkbox", { name: /Patient Consent/ })).toHaveAttribute("aria-required", "true")
  })

  it("should upload with the chosen options and show the success next step", async () => {
    const { onUpload } = renderUploader()
    const pdfFile = buildPdf()

    selectFile(pdfFile)
    fireEvent.click(screen.getByRole("checkbox", { name: /Patient Consent/ }))
    fireEvent.click(screen.getByRole("button", { name: "Submit for Analysis" }))

    await waitFor(() => expect(onUpload).toHaveBeenCalledWith(pdfFile, true, true, false))
    expect(await screen.findByText("Exam sent. Analysis started.", { selector: "span" })).toBeInTheDocument()
    expect(screen.queryByText("labs.pdf")).not.toBeInTheDocument()
  })

  it("should keep the file and offer retry after a failed upload", async () => {
    const { onUpload } = renderUploader()
    onUpload.mockResolvedValueOnce({ status: "failed" })
    const pdfFile = buildPdf()

    selectFile(pdfFile)
    fireEvent.click(screen.getByRole("checkbox", { name: /Patient Consent/ }))
    fireEvent.click(screen.getByRole("button", { name: "Submit for Analysis" }))

    expect(await screen.findByText("The exam could not be sent.", { selector: "span" })).toBeInTheDocument()
    expect(screen.getByText("labs.pdf")).toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: "Try sending again" }))

    await waitFor(() => expect(onUpload).toHaveBeenCalledTimes(2))
    expect(onUpload).toHaveBeenLastCalledWith(pdfFile, true, true, false)
  })

  it("should require a separate acknowledgement when anonymization is unchecked", async () => {
    const { onUpload } = renderUploader()
    const pdfFile = buildPdf()

    selectFile(pdfFile)
    fireEvent.click(screen.getByRole("checkbox", { name: /Patient Consent/ }))
    fireEvent.click(screen.getByRole("checkbox", { name: /Remove identifying file data/ }))

    expect(screen.getByText("The exam will be sent with identifying data")).toBeInTheDocument()
    expect(screen.getByText(/will be processed by an external AI service/)).toBeInTheDocument()
    expect(
      screen.getByText("Confirm that you understand the exam will be processed without anonymization (step 2).")
    ).toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: "Submit for Analysis" }))
    expect(onUpload).not.toHaveBeenCalled()

    fireEvent.click(screen.getByRole("checkbox", { name: /processed without anonymization/ }))
    fireEvent.click(screen.getByRole("button", { name: "Submit for Analysis" }))

    await waitFor(() => expect(onUpload).toHaveBeenCalledWith(pdfFile, true, false, true))
  })

  it("should drop the acknowledgement when anonymization is checked again", async () => {
    const { onUpload } = renderUploader()
    const pdfFile = buildPdf()

    selectFile(pdfFile)
    fireEvent.click(screen.getByRole("checkbox", { name: /Patient Consent/ }))
    fireEvent.click(screen.getByRole("checkbox", { name: /Remove identifying file data/ }))
    fireEvent.click(screen.getByRole("checkbox", { name: /processed without anonymization/ }))
    fireEvent.click(screen.getByRole("checkbox", { name: /Remove identifying file data/ }))

    expect(screen.queryByText("The exam will be sent with identifying data")).not.toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: "Submit for Analysis" }))

    await waitFor(() => expect(onUpload).toHaveBeenCalledWith(pdfFile, true, true, false))
  })

  it.each([
    ["FILE_CONTENT_MISMATCH", "The file content does not match its extension."],
    ["FILE_CORRUPTED", "The file appears to be corrupted."],
    ["FILE_EMPTY", "The server received an empty file."],
    ["DEIDENTIFICATION_REQUIRED", "Sending without anonymization needs your confirmation."],
  ])("should show a specific message for the %s error and keep the file", async (errorCode, expectedTitle) => {
    const { onUpload } = renderUploader()
    onUpload.mockResolvedValueOnce({ status: "failed", errorCode })

    selectFile(buildPdf())
    fireEvent.click(screen.getByRole("checkbox", { name: /Patient Consent/ }))
    fireEvent.click(screen.getByRole("button", { name: "Submit for Analysis" }))

    expect(await screen.findByText(expectedTitle, { selector: "span" })).toBeInTheDocument()
    expect(screen.getByText("labs.pdf")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Replace file" })).toBeInTheDocument()
  })

  it("should show the generic failure for unknown error codes", async () => {
    const { onUpload } = renderUploader()
    onUpload.mockResolvedValueOnce({ status: "failed", errorCode: "SOMETHING_ELSE" })

    selectFile(buildPdf())
    fireEvent.click(screen.getByRole("checkbox", { name: /Patient Consent/ }))
    fireEvent.click(screen.getByRole("button", { name: "Submit for Analysis" }))

    expect(await screen.findByText("The exam could not be sent.", { selector: "span" })).toBeInTheDocument()
  })

  it("should expose a determinate progress bar and a cancel action while uploading", () => {
    const { onCancel } = renderUploader({ isPending: true, uploadProgress: 42 })

    expect(screen.getByRole("progressbar")).toHaveAttribute("aria-valuenow", "42")
    expect(screen.getByText("42%")).toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: "Cancel upload" }))

    expect(onCancel).toHaveBeenCalledTimes(1)
  })

  it("should keep the file after a cancelled upload", async () => {
    const { onUpload } = renderUploader()
    onUpload.mockResolvedValueOnce({ status: "cancelled" })

    selectFile(buildPdf())
    fireEvent.click(screen.getByRole("checkbox", { name: /Patient Consent/ }))
    fireEvent.click(screen.getByRole("button", { name: "Submit for Analysis" }))

    expect(await screen.findByText("Upload cancelled.", { selector: "span" })).toBeInTheDocument()
    expect(screen.getByText("labs.pdf")).toBeInTheDocument()
  })
})
