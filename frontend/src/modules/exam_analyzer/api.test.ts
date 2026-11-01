import { examAnalyzerApi } from "./api"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("examAnalyzerApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should list analyses without filters", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await examAnalyzerApi.getAnalyses()

    expect(http.get).toHaveBeenCalledWith("/exam-analyses")
  })

  it("should list analyses filtered by patient", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await examAnalyzerApi.getAnalyses("patient-1")

    expect(http.get).toHaveBeenCalledWith("/exam-analyses?patientFhirId=patient-1")
  })

  it("should get a single analysis by id", async () => {
    vi.mocked(http.get).mockResolvedValue({ id: "analysis-1" })

    await examAnalyzerApi.getAnalysis("analysis-1")

    expect(http.get).toHaveBeenCalledWith("/exam-analyses/analysis-1")
  })

  it("should upload an exam file with consent, anonymize and patient flags", async () => {
    vi.mocked(http.post).mockResolvedValue({ id: "analysis-2" })
    const file = new File(["dicom-data"], "study.dcm", { type: "application/dicom" })

    await examAnalyzerApi.uploadExamFile({ file, consent: true, anonymize: false, patientFhirId: "patient-9" })

    expect(http.post).toHaveBeenCalledWith("/exam-analyses", expect.any(FormData), expect.any(Object))
    const formData = vi.mocked(http.post).mock.calls[0][1] as FormData
    expect(formData.get("file")).toBe(file)
    expect(formData.get("consent")).toBe("true")
    expect(formData.get("anonymize")).toBe("false")
    expect(formData.get("patientFhirId")).toBe("patient-9")
    expect(formData.get("acknowledgeExternalProcessing")).toBeNull()
  })

  it("should send the external processing acknowledgement only when anonymization is off", async () => {
    vi.mocked(http.post).mockResolvedValue({ id: "analysis-6" })
    const file = new File(["pdf"], "labs.pdf", { type: "application/pdf" })

    await examAnalyzerApi.uploadExamFile({ file, consent: true, anonymize: false, acknowledgeExternalProcessing: true })
    await examAnalyzerApi.uploadExamFile({ file, consent: true, anonymize: true, acknowledgeExternalProcessing: true })

    const withoutAnonymization = vi.mocked(http.post).mock.calls[0][1] as FormData
    const withAnonymization = vi.mocked(http.post).mock.calls[1][1] as FormData
    expect(withoutAnonymization.get("acknowledgeExternalProcessing")).toBe("true")
    expect(withAnonymization.get("acknowledgeExternalProcessing")).toBeNull()
  })

  it("should omit the patient filter from the upload when absent", async () => {
    vi.mocked(http.post).mockResolvedValue({ id: "analysis-3" })
    const file = new File(["dicom-data"], "study.dcm", { type: "application/dicom" })

    await examAnalyzerApi.uploadExamFile({ file, consent: false, anonymize: true })

    const formData = vi.mocked(http.post).mock.calls[0][1] as FormData
    expect(formData.get("patientFhirId")).toBeNull()
  })

  it("should report upload progress as a percentage", async () => {
    vi.mocked(http.post).mockResolvedValue({ id: "analysis-4" })
    const file = new File(["dicom-data"], "study.dcm", { type: "application/dicom" })
    const onUploadProgress = vi.fn()

    await examAnalyzerApi.uploadExamFile({ file, consent: true, anonymize: true, onUploadProgress })

    const requestConfig = vi.mocked(http.post).mock.calls[0][2] as { onUploadProgress: (event: { loaded: number; total: number }) => void }
    requestConfig.onUploadProgress({ loaded: 50, total: 200 })

    expect(onUploadProgress).toHaveBeenCalledWith(25)
  })

  it("should delete an analysis by id", async () => {
    vi.mocked(http.delete).mockResolvedValue({ success: "deleted" })

    await examAnalyzerApi.deleteAnalysis("analysis-5")

    expect(http.delete).toHaveBeenCalledWith("/exam-analyses/analysis-5")
  })
})
