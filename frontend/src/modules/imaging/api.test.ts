import { imagingApi } from "./api"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("imagingApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should list imaging studies for a patient", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await imagingApi.getImagingStudies("patient-1")

    expect(http.get).toHaveBeenCalledWith("/patients/patient-1/studies")
  })

  it("should get a single imaging study", async () => {
    vi.mocked(http.get).mockResolvedValue({ id: "study-1" })

    await imagingApi.getImagingStudy("study-1")

    expect(http.get).toHaveBeenCalledWith("/studies/study-1")
  })

  it("should upload an imaging study with title, modality and file", async () => {
    vi.mocked(http.post).mockResolvedValue({ id: "study-2" })
    const dicomBlob = new Blob(["dicom"], { type: "application/dicom" })

    await imagingApi.uploadImagingStudy({
      patientFhirId: "patient-1",
      title: "Raio-X tórax",
      modality: "CT",
      dicomBlob,
    })

    expect(http.post).toHaveBeenCalledWith("/patients/patient-1/studies", expect.any(FormData))
    const formData = vi.mocked(http.post).mock.calls[0][1] as FormData
    expect(formData.get("title")).toBe("Raio-X tórax")
    expect(formData.get("modality")).toBe("CT")
    const uploadedFile = formData.get("file") as Blob
    expect(uploadedFile).toBeInstanceOf(Blob)
    expect(uploadedFile.type).toBe("application/dicom")
    expect(uploadedFile.size).toBe(dicomBlob.size)
  })
})
