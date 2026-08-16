import { patientsApi } from "./api"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("patientsApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should list patients without query params by default", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await patientsApi.getPatients()

    expect(http.get).toHaveBeenCalledWith("/patients")
  })

  it("should build the query string from pagination and sort filters", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await patientsApi.getPatients("Ana", "full_name", "desc", 2, 10)

    expect(http.get).toHaveBeenCalledWith("/patients?search=Ana&sortField=full_name&sortDirection=desc&page=2&limit=10")
  })

  it("should get a single patient by fhir id", async () => {
    vi.mocked(http.get).mockResolvedValue({ id: "p1" })

    await patientsApi.getPatient("p1")

    expect(http.get).toHaveBeenCalledWith("/patients/p1")
  })

  it("should combine the created patient ids with the submitted data", async () => {
    vi.mocked(http.post).mockResolvedValue({
      patient_id: "local-1",
      fhir_resource_id: "fhir-1",
    })

    const result = await patientsApi.createPatient({
      full_name: "Maria Silva",
      birth_date: "1990-01-01",
      document_id: "52998224725",
      phone_number: "(21) 99999-0000",
    })

    expect(http.post).toHaveBeenCalledWith("/patients", {
      full_name: "Maria Silva",
      birth_date: "1990-01-01",
      document_id: "52998224725",
      phone_number: "(21) 99999-0000",
    })
    expect(result).toMatchObject({
      patient_id: "local-1",
      fhir_resource_id: "fhir-1",
      full_name: "Maria Silva",
    })
  })

  it("should list encounters for a patient", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await patientsApi.getEncounters("patient-1")

    expect(http.get).toHaveBeenCalledWith("/patients/patient-1/encounters")
  })

  it("should create an encounter with the reason payload", async () => {
    vi.mocked(http.post).mockResolvedValue({ fhir_id: "enc-1" })

    await patientsApi.createEncounter({
      patient_fhir_id: "patient-1",
      practitioner_id: "practitioner-1",
      reason_display: "Dor abdominal",
    })

    expect(http.post).toHaveBeenCalledWith("/patients/patient-1/encounters", {
      reason_code: "",
      reason_display: "Dor abdominal",
      practitioner_id: "practitioner-1",
    })
  })

  it("should update an encounter status", async () => {
    vi.mocked(http.put).mockResolvedValue({ fhir_id: "enc-1", status: "finished" })

    await patientsApi.updateEncounter({ encounter_fhir_id: "enc-1", status: "finished" })

    expect(http.put).toHaveBeenCalledWith("/encounters/enc-1", { status: "finished" })
  })

  it("should list observations for an encounter", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await patientsApi.getObservations("enc-1")

    expect(http.get).toHaveBeenCalledWith("/encounters/enc-1/observations")
  })

  it("should list all observations for a patient", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await patientsApi.getAllPatientObservations("patient-1")

    expect(http.get).toHaveBeenCalledWith("/patients/patient-1/observations")
  })

  it("should create an observation with the clinical payload", async () => {
    vi.mocked(http.post).mockResolvedValue({ fhir_id: "obs-1" })

    await patientsApi.createObservation({
      encounter_fhir_id: "enc-1",
      patient_fhir_id: "patient-1",
      loinc_code: "8867-4",
      code_display: "Frequência cardíaca",
      value_quantity: 72,
      value_unit: "bpm",
    })

    expect(http.post).toHaveBeenCalledWith("/encounters/enc-1/observations", {
      patient_fhir_id: "patient-1",
      loinc_code: "8867-4",
      code_display: "Frequência cardíaca",
      value_quantity: 72,
      value_unit: "bpm",
    })
  })

  it("should list conditions for a patient", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await patientsApi.getConditions("patient-1")

    expect(http.get).toHaveBeenCalledWith("/patients/patient-1/conditions")
  })

  it("should create a condition for a patient", async () => {
    vi.mocked(http.post).mockResolvedValue({ fhir_id: "cond-1" })

    await patientsApi.createCondition({
      patient_fhir_id: "patient-1",
      icd10_code: "E11",
      code_display: "Diabetes",
    })

    expect(http.post).toHaveBeenCalledWith("/patients/patient-1/conditions", {
      icd10_code: "E11",
      code_display: "Diabetes",
    })
  })

  it("should list diagnostic reports for an encounter", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await patientsApi.getDiagnosticReports("enc-1")

    expect(http.get).toHaveBeenCalledWith("/encounters/enc-1/reports")
  })

  it("should list report versions", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await patientsApi.getDiagnosticReportVersions("report-1")

    expect(http.get).toHaveBeenCalledWith("/reports/report-1/versions")
  })

  it("should create a diagnostic report with the clinical payload", async () => {
    vi.mocked(http.post).mockResolvedValue({ fhir_id: "report-1" })

    await patientsApi.createDiagnosticReport({
      encounter_fhir_id: "enc-1",
      patient_fhir_id: "patient-1",
      report_code: "8867-4",
      report_display: "Frequência cardíaca",
      conclusion: "Normal",
    })

    expect(http.post).toHaveBeenCalledWith("/encounters/enc-1/reports", {
      patient_fhir_id: "patient-1",
      report_code: "8867-4",
      report_display: "Frequência cardíaca",
      conclusion: "Normal",
    })
  })

  it("should list allergies for a patient", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await patientsApi.getAllergies("patient-1")

    expect(http.get).toHaveBeenCalledWith("/patients/patient-1/allergies")
  })

  it("should create an allergy for a patient", async () => {
    vi.mocked(http.post).mockResolvedValue({ fhir_id: "allergy-1" })

    await patientsApi.createAllergy({
      patient_fhir_id: "patient-1",
      allergen_code: "J30",
      allergen_display: "Alergia a pólen",
      reaction: "Espirros",
    })

    expect(http.post).toHaveBeenCalledWith("/patients/patient-1/allergies", {
      allergen_code: "J30",
      allergen_display: "Alergia a pólen",
      reaction: "Espirros",
    })
  })

  it("should list medications for an encounter", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await patientsApi.getMedications("enc-1")

    expect(http.get).toHaveBeenCalledWith("/encounters/enc-1/medications")
  })

  it("should create a medication request", async () => {
    vi.mocked(http.post).mockResolvedValue({ fhir_id: "med-1" })

    await patientsApi.createMedication({
      encounter_fhir_id: "enc-1",
      patient_fhir_id: "patient-1",
      medication_name: "Metformina",
      dosage_instructions: "2x ao dia",
    })

    expect(http.post).toHaveBeenCalledWith("/encounters/enc-1/medications", {
      patient_fhir_id: "patient-1",
      medication_name: "Metformina",
      dosage_instructions: "2x ao dia",
    })
  })
})
