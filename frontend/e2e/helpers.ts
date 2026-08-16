import { type Page, expect } from "@playwright/test"

interface MockAuthCredentials {
  email: string
  password: string
}

export const mockAuthAPI = async (
  pageInstance: Page,
  sessionRole: string = "DOCTOR",
  sessionCredentials: MockAuthCredentials = { email: "medico@clinica.com", password: "senha123" }
): Promise<void> => {
  const sessionToken = `mock-jwt-token-${sessionRole.toLowerCase()}-123456`
  const sessionUserId = `user-${sessionRole.toLowerCase()}-123`
  const sessionFullName = sessionRole === "PATIENT" ? "Guilherme de Souza Araujo" : "Dr. AndrǸ Silva de Araujo"
  let sessionLoggedIn = false

  await pageInstance.route("**/api/v1/auth/login", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const submittedJSON = httpRequest.postDataJSON()

    if (
      submittedJSON.email === sessionCredentials.email &&
      submittedJSON.password === sessionCredentials.password
    ) {
      sessionLoggedIn = true
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          token: sessionToken,
          userId: sessionUserId,
          role: sessionRole,
          email: sessionCredentials.email,
        }),
      })
    } else {
      await networkRoute.fulfill({
        status: 401,
        contentType: "application/json",
        body: JSON.stringify({
          error: "Credenciais inválidas.",
        }),
      })
    }
  })

  await pageInstance.route("**/api/v1/auth/logout", async (networkRoute) => {
    sessionLoggedIn = false
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        success: true,
      }),
    })
  })

  await pageInstance.route("**/api/v1/auth/me", async (networkRoute) => {
    if (sessionLoggedIn) {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          token: sessionToken,
          userId: sessionUserId,
          role: sessionRole,
          email: sessionCredentials.email,
          fullName: sessionFullName,
          isActive: true,
        }),
      })
      return
    }
    await networkRoute.fulfill({
      status: 401,
      contentType: "application/json",
      body: JSON.stringify({
        error: "Nǜo autenticado.",
      }),
    })
  })
}

export const mockPatientsAPI = async (pageInstance: Page): Promise<void> => {
  const currentPatientsList = [
    {
      patient_id: "pat-1",
      fhir_resource_id: "fhir-pat-1",
      full_name: "Guilherme de Souza Araujo",
      birth_date: "1988-04-12",
      document_id: "123.456.789-00",
      phone_number: "(11) 98765-4321",
    },
    {
      patient_id: "pat-2",
      fhir_resource_id: "fhir-pat-2",
      full_name: "Mariana Costa Silva",
      birth_date: "1995-11-23",
      document_id: "987.654.321-11",
      phone_number: "(21) 99999-8888",
    },
  ]

  await pageInstance.route("**/api/v1/patients*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    if (httpRequest.method() === "GET") {
      const requestURL = new URL(httpRequest.url())
      const searchTerm = (requestURL.searchParams.get("search") ?? "").toLowerCase()

      const filteredPatients = searchTerm
        ? currentPatientsList.filter((patient) =>
            patient.full_name.toLowerCase().includes(searchTerm) ||
            patient.document_id.toLowerCase().includes(searchTerm) ||
            patient.phone_number.toLowerCase().includes(searchTerm)
          )
        : currentPatientsList

      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filteredPatients),
      })
    } else if (httpRequest.method() === "POST") {
      const submittedJSON = httpRequest.postDataJSON()
      const newPatientId = `pat-${currentPatientsList.length + 1}`
      const newFhirResourceId = `fhir-pat-${currentPatientsList.length + 1}`
      const newPatient = {
        patient_id: newPatientId,
        fhir_resource_id: newFhirResourceId,
        full_name: submittedJSON.full_name,
        birth_date: submittedJSON.birth_date,
        document_id: submittedJSON.document_id,
        phone_number: submittedJSON.phone_number,
      }
      currentPatientsList.push(newPatient)
      await networkRoute.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify({
          patient_id: newPatientId,
          fhir_resource_id: newFhirResourceId,
        }),
      })
    }
  })

  await pageInstance.route("**/api/v1/patients/*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const targetResourceId = urlParts[urlParts.length - 1]

    const matchedPatient = currentPatientsList.find((patient) => patient.fhir_resource_id === targetResourceId)
    if (matchedPatient) {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(matchedPatient),
      })
    } else {
      await networkRoute.fulfill({
        status: 404,
        contentType: "application/json",
        body: JSON.stringify({ error: "Paciente não encontrado." }),
      })
    }
  })
}

export const mockClinicalAPI = async (pageInstance: Page): Promise<void> => {
  const currentEncountersList = [
    {
      fhir_id: "enc-1",
      patient_fhir_id: "fhir-pat-1",
      status: "finished",
      reason_display: "Consulta de Rotina Geral",
      created_at: "2026-05-10T10:00:00Z",
    },
    {
      fhir_id: "enc-2",
      patient_fhir_id: "fhir-pat-1",
      status: "finished",
      reason_display: "Retorno Cardiológico",
      created_at: "2026-05-15T14:30:00Z",
    },
  ]

  const currentObservationsList = [
    {
      fhir_id: "obs-1",
      encounter_fhir_id: "enc-1",
      patient_fhir_id: "fhir-pat-1",
      loinc_code: "8867-4",
      code_display: "Frequência Cardíaca",
      value_quantity: 72,
      value_unit: "bpm",
      created_at: "2026-05-10T10:05:00Z",
    },
    {
      fhir_id: "obs-2",
      encounter_fhir_id: "enc-1",
      patient_fhir_id: "fhir-pat-1",
      loinc_code: "85354-9",
      code_display: "Pressão Arterial Sistólica",
      value_quantity: 120,
      value_unit: "mmHg",
      created_at: "2026-05-10T10:05:00Z",
    },
    {
      fhir_id: "obs-3",
      encounter_fhir_id: "enc-1",
      patient_fhir_id: "fhir-pat-1",
      loinc_code: "8310-5",
      code_display: "Temperatura Corporal",
      value_quantity: 36.5,
      value_unit: "°C",
      created_at: "2026-05-10T10:05:00Z",
    },
    {
      fhir_id: "obs-4",
      encounter_fhir_id: "enc-2",
      patient_fhir_id: "fhir-pat-1",
      loinc_code: "8867-4",
      code_display: "Frequência Cardíaca",
      value_quantity: 85,
      value_unit: "bpm",
      created_at: "2026-05-15T14:35:00Z",
    },
    {
      fhir_id: "obs-5",
      encounter_fhir_id: "enc-2",
      patient_fhir_id: "fhir-pat-1",
      loinc_code: "85354-9",
      code_display: "Pressão Arterial Sistólica",
      value_quantity: 135,
      value_unit: "mmHg",
      created_at: "2026-05-15T14:35:00Z",
    },
  ]

  const currentConditionsList = [
    {
      fhir_id: "cond-1",
      patient_fhir_id: "fhir-pat-1",
      icd10_code: "I10",
      code_display: "Hipertensão Essencial Primária",
      clinical_status: "active",
      created_at: "2026-05-15T14:40:00Z",
    },
  ]

  const currentReportsList = [
    {
      fhir_id: "rep-1",
      encounter_fhir_id: "enc-2",
      patient_fhir_id: "fhir-pat-1",
      report_display: "Eletrocardiograma de Repouso",
      status: "final",
      conclusion: "Ritmo sinusal com leve taquicardia. Recomenda-se acompanhamento ambulatorial.",
      created_at: "2026-05-15T14:45:00Z",
    },
  ]

  const currentMedicationsList: unknown[] = []

  const currentStudiesList = [
    {
      id: "study-1",
      patient_fhir_id: "fhir-pat-1",
      title: "Tomografia Computadorizada de Tórax",
      modality: "CT",
      study_instance_uid: "1.2.840.10008.5.1.4.1.1.2.20260516.1",
      status: "completed",
      created_at: "2026-05-16T10:00:00Z",
    },
  ]

  await pageInstance.route("**/api/v1/patients/*/encounters", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const patientFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentEncountersList.filter((encounter) => encounter.patient_fhir_id === patientFhirId)
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filtered),
      })
    } else if (httpRequest.method() === "POST") {
      const submittedJSON = httpRequest.postDataJSON()
      const newEncounter = {
        fhir_id: `enc-${currentEncountersList.length + 1}`,
        patient_fhir_id: patientFhirId,
        status: "in-progress",
        reason_display: submittedJSON.reason_display,
        practitioner_id: submittedJSON.practitioner_id,
        created_at: new Date().toISOString(),
      }
      currentEncountersList.push(newEncounter)
      await networkRoute.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(newEncounter),
      })
    }
  })

  await pageInstance.route("**/api/v1/encounters/*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    if (httpRequest.method() !== "PUT") {
      await networkRoute.continue()
      return
    }
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const encounterFhirId = urlParts[urlParts.length - 1]
    const submittedJSON = httpRequest.postDataJSON()
    const matchedEncounter = currentEncountersList.find((encounter) => encounter.fhir_id === encounterFhirId)
    if (matchedEncounter) {
      matchedEncounter.status = submittedJSON.status
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(matchedEncounter),
      })
    } else {
      await networkRoute.fulfill({
        status: 404,
        contentType: "application/json",
        body: JSON.stringify({ error: "Encounter not found." }),
      })
    }
  })

  await pageInstance.route("**/api/v1/patients/*/observations", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const patientFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentObservationsList.filter((observation) => observation.patient_fhir_id === patientFhirId)
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filtered),
      })
    }
  })

  await pageInstance.route("**/api/v1/patients/*/conditions", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const patientFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentConditionsList.filter((condition) => condition.patient_fhir_id === patientFhirId)
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filtered),
      })
    } else if (httpRequest.method() === "POST") {
      const submittedJSON = httpRequest.postDataJSON()
      const newCondition = {
        fhir_id: `cond-${currentConditionsList.length + 1}`,
        patient_fhir_id: patientFhirId,
        icd10_code: submittedJSON.icd10_code,
        code_display: submittedJSON.code_display,
        clinical_status: "active",
        created_at: new Date().toISOString(),
      }
      currentConditionsList.push(newCondition)
      await networkRoute.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(newCondition),
      })
    }
  })

  await pageInstance.route("**/api/v1/encounters/*/observations", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const encounterFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentObservationsList.filter((observation) => observation.encounter_fhir_id === encounterFhirId)
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filtered),
      })
    } else if (httpRequest.method() === "POST") {
      const submittedJSON = httpRequest.postDataJSON()
      const newObservation = {
        fhir_id: `obs-${currentObservationsList.length + 1}`,
        encounter_fhir_id: encounterFhirId,
        patient_fhir_id: submittedJSON.patient_fhir_id,
        loinc_code: submittedJSON.loinc_code,
        code_display: submittedJSON.code_display,
        value_quantity: submittedJSON.value_quantity,
        value_unit: submittedJSON.value_unit,
        created_at: new Date().toISOString(),
      }
      currentObservationsList.push(newObservation)
      await networkRoute.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(newObservation),
      })
    }
  })

  await pageInstance.route("**/api/v1/encounters/*/reports", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const encounterFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentReportsList.filter((report) => report.encounter_fhir_id === encounterFhirId)
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filtered),
      })
    } else if (httpRequest.method() === "POST") {
      const submittedJSON = httpRequest.postDataJSON()
      const newReport = {
        fhir_id: `rep-${currentReportsList.length + 1}`,
        encounter_fhir_id: encounterFhirId,
        patient_fhir_id: submittedJSON.patient_fhir_id,
        report_display: submittedJSON.report_display,
        status: "final",
        conclusion: submittedJSON.conclusion,
        created_at: new Date().toISOString(),
      }
      currentReportsList.push(newReport)
      await networkRoute.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(newReport),
      })
    }
  })

  await pageInstance.route("**/api/v1/reports/*/versions", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const reportFhirId = urlParts[urlParts.length - 2]
    const reportEntry = currentReportsList.find((report) => report.fhir_id === reportFhirId)
    if (!reportEntry) {
      await networkRoute.fulfill({
        status: 404,
        contentType: "application/json",
        body: JSON.stringify({ error: "Report not found." }),
      })
      return
    }
    const versionsForReport = [
      {
        version: "1",
        snapshot: {
          fhir_resource_id: reportEntry.fhir_id,
          report_display: reportEntry.report_display,
          conclusion: reportEntry.conclusion,
        },
        changed_by: "Dr. André Silva de Araujo",
        changed_at: reportEntry.created_at,
      },
      {
        version: "2",
        snapshot: {
          fhir_resource_id: reportEntry.fhir_id,
          report_display: reportEntry.report_display,
          conclusion: "Conclusão revisada após reavaliação clínica do traçado.",
        },
        changed_by: "Dr. André Silva de Araujo",
        changed_at: new Date().toISOString(),
      },
    ]
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(versionsForReport),
    })
  })

  await pageInstance.route("**/api/v1/encounters/*/medications", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const encounterFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentMedicationsList.filter((medication) => medication.encounter_fhir_id === encounterFhirId)
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filtered),
      })
    } else if (httpRequest.method() === "POST") {
      const submittedJSON = httpRequest.postDataJSON()
      const newMedication = {
        fhir_id: `med-${currentMedicationsList.length + 1}`,
        encounter_fhir_id: encounterFhirId,
        patient_fhir_id: submittedJSON.patient_fhir_id,
        medication_name: submittedJSON.medication_name,
        dosage_instructions: submittedJSON.dosage_instructions,
        status: "active",
        created_at: new Date().toISOString(),
      }
      currentMedicationsList.push(newMedication)
      await networkRoute.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(newMedication),
      })
    }
  })

  await pageInstance.route("**/api/v1/patients/*/studies", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const patientFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentStudiesList.filter((study) => study.patient_fhir_id === patientFhirId)
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filtered),
      })
    } else if (httpRequest.method() === "POST") {
      const newStudy = {
        id: `study-${currentStudiesList.length + 1}`,
        patient_fhir_id: patientFhirId,
        title: "Nova Imagem (Simulada)",
        modality: "MR",
        study_instance_uid: `1.2.840.10008.5.1.4.1.1.2.20260516.${currentStudiesList.length + 1}`,
        status: "completed",
        created_at: new Date().toISOString(),
      }
      currentStudiesList.push(newStudy)
      await networkRoute.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(newStudy),
      })
    }
  })

  await pageInstance.route("**/api/v1/studies/*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const studyId = urlParts[urlParts.length - 1]

    if (httpRequest.method() === "GET") {
      const study = currentStudiesList.find((study) => study.id === studyId)
      if (study) {
        await networkRoute.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({ ...study, download_url: "mock_url" }),
        })
      } else {
        await networkRoute.fulfill({ status: 404 })
      }
    }
  })
}

export const mockAnalyzerAPI = async (pageInstance: Page): Promise<void> => {
  const currentAnalysesList = [
    {
      id: "ana-1",
      user_id: "user-medico-123",
      patient_fhir_id: "fhir-pat-1",
      exam_type: "Radiografia Digital de Tórax (PA)",
      file_name: "rx_torax.png",
      file_path: "tmp/exam_uploads/ana-1.png",
      status: "completed",
      analysis_response: {
        examType: "Radiografia Digital de Tórax (PA)",
        qualityAssessment: {
          score: 0.9,
          warnings: ["Inspiração adequada. Sem artefatos de movimento."]
        },
        detectedFindings: [
          {
            finding: "Área de consolidação pulmonar no lobo inferior direito",
            confidence: 0.88,
            severity: "high"
          },
          {
            finding: "Ausência de derrame pleural",
            confidence: 0.95,
            severity: "low"
          }
        ],
        possibleInterpretations: [
          "Sinais sugestivos de pneumonia lobar. Correlacionar com quadro clínico."
        ],
        recommendation: {
          urgency: "urgent",
          nextSteps: ["Agendar consulta com pneumologista ou clínico geral."]
        },
        limitations: ["Análise baseada em algoritmo assistivo."],
        disclaimer: "ESTE LAUDO É ASSISTIVO. OS RESULTADOS SÃO PRELIMINARES."
      },
      consent_given: true,
      anonymized: false,
      created_at: "2026-05-18T10:00:00Z",
      updated_at: "2026-05-18T10:00:00Z"
    }
  ]

  await pageInstance.route("**/api/v1/exam-analyses", async (networkRoute) => {
    const httpRequest = networkRoute.request()

    if (httpRequest.method() === "GET") {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(currentAnalysesList),
      })
    } else if (httpRequest.method() === "POST") {
      const newAnalysis: Record<string, unknown> = {
        id: `ana-${currentAnalysesList.length + 1}`,
        user_id: "user-medico-123",
        patient_fhir_id: undefined,
        exam_type: undefined,
        file_name: "mock_uploaded_exam.jpg",
        file_path: "tmp/exam_uploads/mock_uploaded_exam.jpg",
        status: "processing",
        analysis_response: {
          status: "pending"
        },
        consent_given: true,
        anonymized: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }
      currentAnalysesList.push(newAnalysis)
      
      setTimeout(() => {
        newAnalysis.status = "completed"
        newAnalysis.exam_type = "Radiografia Digital de Tórax (PA)"
        newAnalysis.analysis_response = {
          examType: "Radiografia Digital de Tórax (PA)",
          qualityAssessment: {
            score: 0.95,
            warnings: []
          },
          detectedFindings: [
            {
              finding: "Nódulo pulmonar calcificado",
              confidence: 0.92,
              severity: "low"
            },
            {
              finding: "Aorta normal",
              confidence: 0.98,
              severity: "low"
            }
          ],
          possibleInterpretations: [
            "Achados benignos, sem necessidade de investigação adicional imediata."
          ],
          recommendation: {
            urgency: "normal",
            nextSteps: ["Acompanhamento clínico periódico padrão."]
          },
          limitations: ["Radiografia simples possui limitações estruturais."],
          disclaimer: "ESTE LAUDO É ASSISTIVO. RECOMENDA-SE AVALIAÇÃO CLÍNICA COMPLETA."
        }
      }, 3000)

      await networkRoute.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(newAnalysis),
      })
    }
  })

  await pageInstance.route("**/api/v1/exam-analyses/*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const anaId = urlParts[urlParts.length - 1]

    if (httpRequest.method() === "GET") {
      const analysis = currentAnalysesList.find((analysis) => analysis.id === anaId)
      if (analysis) {
        await networkRoute.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(analysis),
        })
      } else {
        await networkRoute.fulfill({ status: 404 })
      }
    } else if (httpRequest.method() === "DELETE") {
      const index = currentAnalysesList.findIndex((analysis) => analysis.id === anaId)
      if (index !== -1) {
        currentAnalysesList.splice(index, 1)
        await networkRoute.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({ success: true }),
        })
      } else {
        await networkRoute.fulfill({ status: 404 })
      }
    }
  })
}

export const mockStaffAPI = async (pageInstance: Page): Promise<void> => {
  const currentEmployees = [
    {
      id: "emp-1",
      userId: "user-1",
      full_name: "Dr. André Silva de Araujo",
      email: "andre.silva@hospital.com",
      role: "doctor",
      crm_number: "CRM-SP 12345",
      is_active: true,
      department: "Cardiologia",
      fhir_resource_id: "fhir-emp-1",
    },
    {
      id: "emp-2",
      userId: "user-2",
      full_name: "Enf. Roberta Santos Almeida",
      email: "roberta.santos@hospital.com",
      role: "nurse",
      crm_number: "COREN-SP 54321",
      is_active: true,
      department: "Pediatria",
      fhir_resource_id: "fhir-emp-2",
    },
  ]

  await pageInstance.route("**/api/v1/staff/employees*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    if (httpRequest.method() === "GET") {
      const requestURL = new URL(httpRequest.url())
      const searchTerm = (requestURL.searchParams.get("search") ?? "").toLowerCase()
      const roleFilter = requestURL.searchParams.get("role")

      const filteredEmployees = currentEmployees.filter((employee) => {
        const matchesRole =
          !roleFilter ||
          roleFilter === "All" ||
          employee.role.toLowerCase() === roleFilter.toLowerCase()
        const matchesSearch =
          !searchTerm ||
          employee.full_name.toLowerCase().includes(searchTerm) ||
          employee.email.toLowerCase().includes(searchTerm)
        return matchesRole && matchesSearch
      })

      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filteredEmployees),
      })
    } else if (httpRequest.method() === "POST") {
      const submittedJSON = httpRequest.postDataJSON()
      const newEmployee = {
        id: `emp-${currentEmployees.length + 1}`,
        user_id: submittedJSON.user_id || `user-${currentEmployees.length + 1}`,
        full_name: submittedJSON.full_name,
        email: submittedJSON.email,
        role: submittedJSON.role,
        crm_number: submittedJSON.crm_number || "N/A",
        is_active: true,
        department: "Geral",
        fhir_resource_id: `fhir-emp-${currentEmployees.length + 1}`,
      }
      currentEmployees.push(newEmployee)
      await networkRoute.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify({ employeeId: newEmployee.id }),
      })
    }
  })
}

export const mockTelemetryAPI = async (pageInstance: Page): Promise<void> => {
  const roomsList = [
    {
      id: "room-1",
      name: "Sala Verde - Semi-Intensiva",
      description: "Monitoramento semi-intensivo",
    },
    {
      id: "room-2",
      name: "Sala Vermelha - Choque & Emergência",
      description: "Leitos críticos e trauma",
    },
  ]

  const bedsByRoom: Record<string, Record<string, string | number | boolean>[]> = {
    "room-1": [
      {
        id: "bed-1",
        roomId: "room-1",
        bedNumber: "Leito 01",
        patientName: "Guilherme de Souza Araujo",
        age: 38,
        gender: "male",
        bpm: 72,
        spo2: 98,
        temperature: 36.5,
        status: "normal",
        condition: "Normal",
      },
    ],
    "room-2": [
      {
        id: "bed-2",
        roomId: "room-2",
        bedNumber: "Leito 02",
        patientName: "Mariana Costa Silva",
        age: 30,
        gender: "female",
        bpm: 80,
        spo2: 95,
        temperature: 37.0,
        status: "normal",
        condition: "Normal",
      },
    ],
  }

  await pageInstance.route("**/api/v1/telemetry/rooms", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(roomsList),
    })
  })

  await pageInstance.route("**/api/v1/telemetry/rooms/*/unlock", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const submittedJSON = httpRequest.postDataJSON()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const roomId = urlParts[urlParts.length - 2]
    const matchedRoom = roomsList.find((roomItem) => roomItem.id === roomId)

    if (submittedJSON.passcode === "9999") {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          success: true,
          roomName: matchedRoom ? matchedRoom.name : "Sala Desbloqueada",
        }),
      })
    } else {
      await networkRoute.fulfill({
        status: 400,
        contentType: "application/json",
        body: JSON.stringify({ error: "Passcode incorreto" }),
      })
    }
  })

  await pageInstance.route("**/api/v1/telemetry/rooms/*/beds", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const roomId = urlParts[urlParts.length - 2]
    const beds = bedsByRoom[roomId] || []

    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(beds),
    })
  })

  await pageInstance.route("**/api/v1/telemetry/beds/*/condition", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const submittedJSON = httpRequest.postDataJSON()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const bedId = urlParts[urlParts.length - 2]

    for (const key of Object.keys(bedsByRoom)) {
      const bed = bedsByRoom[key].find((bedItem) => bedItem.id === bedId)
      if (bed) {
        bed.bpm = submittedJSON.bpm
        bed.spo2 = submittedJSON.spo2
        bed.temperature = submittedJSON.temperature
        bed.status = submittedJSON.status
        bed.condition = submittedJSON.condition
        break
      }
    }

    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ success: true }),
    })
  })
}

export const mockAnalyticsAPI = async (pageInstance: Page): Promise<void> => {
  await pageInstance.route("**/api/v1/analytics", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        totalRegisteredPatients: 340,
        fhirComplianceRate: 99.4,
        averageServiceDurationMinutes: 14.5,
        activeConsultationsTotal: 79,
        totalStudiesCount: 35,
        examModalitiesData: [
          { modality: "CT (Tomografia)", percentage: 45, count: 16, color: "#2563eb" },
          { modality: "MR (Ressonância)", percentage: 30, count: 11, color: "#0d9488" },
          { modality: "CR (Raio-X)", percentage: 15, count: 5, color: "#8b5cf6" },
          { modality: "US (Ultrassom)", percentage: 10, count: 3, color: "#f59e0b" }
        ],
        consultationsWeeklyData: [
          { dayName: "analytics.days.mon", count: 8 },
          { dayName: "analytics.days.tue", count: 12 },
          { dayName: "analytics.days.wed", count: 14 },
          { dayName: "analytics.days.thu", count: 11 },
          { dayName: "analytics.days.fri", count: 15 },
          { dayName: "analytics.days.sat", count: 5 },
          { dayName: "analytics.days.sun", count: 2 }
        ],
        pathologies: [
          { code: "J45.9", descriptionKey: "analytics.pathologies.asthma", categoryKey: "analytics.categories.respiratory", activeCases: 44, trend: "+5%" },
          { code: "I10", descriptionKey: "analytics.pathologies.hypertension", categoryKey: "analytics.categories.cardiovascular", activeCases: 119, trend: "stable" },
          { code: "E11.9", descriptionKey: "analytics.pathologies.diabetes", categoryKey: "analytics.categories.endocrine", activeCases: 85, trend: "+12%" }
        ]
      })
    })
  })
}

export const mockAuditLogsAPI = async (pageInstance: Page): Promise<void> => {
  const currentAuditLogs = [
    {
      id: "log-1",
      correlation_id: "corr-001",
      caller_user_id: "admin@hospital.com",
      caller_role: "ADMIN",
      method: "API_REQUEST",
      access_granted: true,
      resource_type: "patient",
      resource_id: "fhir-pat-1",
      action: "patient.created",
      payload_diff: { full_name: { from: null, to: "Guilherme de Souza Araujo" } },
      created_at: "2026-07-03T10:00:00Z",
    },
    {
      id: "log-2",
      correlation_id: "corr-002",
      caller_user_id: "medico@clinica.com",
      caller_role: "DOCTOR",
      method: "API_REQUEST",
      access_granted: true,
      resource_type: "diagnostic_report",
      resource_id: "fhir-rep-1",
      action: "report.updated",
      payload_diff: { meta: { versionId: { from: "1", to: "2" } } },
      created_at: "2026-07-03T09:30:00Z",
    },
    {
      id: "log-3",
      correlation_id: "corr-003",
      caller_user_id: "usuario.invalido@test.com",
      caller_role: "UNKNOWN",
      method: "LOGIN",
      access_granted: false,
      created_at: "2026-07-03T08:00:00Z",
    },
  ]

  await pageInstance.route("**/api/v1/audit-logs*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    if (httpRequest.method() === "GET") {
      const requestURL = new URL(httpRequest.url())
      const targetStatus = requestURL.searchParams.get("status")
      const targetEmail = requestURL.searchParams.get("email")
      const targetAction = requestURL.searchParams.get("action")

      const filteredLogs = currentAuditLogs.filter((logEntry) => {
        if (targetStatus && targetStatus !== "All") {
          const expectedAccess = targetStatus === "SUCCESS"
          if (logEntry.access_granted !== expectedAccess) {
            return false
          }
        }
        if (targetEmail && !logEntry.caller_user_id.toLowerCase().includes(targetEmail.toLowerCase())) {
          return false
        }
        if (targetAction && targetAction !== "All" && logEntry.method !== targetAction) {
          return false
        }
        return true
      })

      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ audit_logs: filteredLogs, total: filteredLogs.length }),
      })
    }
  })
}

export const mockScheduleAPI = async (pageInstance: Page): Promise<void> => {
  const currentAppointments: Record<string, unknown>[] = []
  const allowedSlotDurationsMinutes = [30, 45]
  const allowedStartMinutes = [0, 30, 45]

  const hasTimeOverlap = (
    firstStart: Date,
    firstEnd: Date,
    secondStart: Date,
    secondEnd: Date
  ): boolean => firstStart < secondEnd && firstEnd > secondStart

  await pageInstance.route("**/api/v1/appointments*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = new URL(httpRequest.url())

    if (httpRequest.method() === "GET") {
      const targetStaffId = requestURL.searchParams.get("staff_id")
      const targetDate = requestURL.searchParams.get("date")
      const targetPatientId = requestURL.searchParams.get("patient_fhir_id")
      const filtered = currentAppointments.filter((appointment) => {
        if (targetStaffId && appointment.staff_id !== targetStaffId) {
          return false
        }
        if (targetDate && !String(appointment.starts_at).startsWith(targetDate)) {
          return false
        }
        if (targetPatientId && appointment.patient_fhir_id !== targetPatientId) {
          return false
        }
        return true
      })
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filtered),
      })
      return
    }

    if (httpRequest.method() === "POST") {
      const submittedJSON = httpRequest.postDataJSON()
      const newStart = new Date(submittedJSON.starts_at)
      const newEnd = new Date(submittedJSON.ends_at)

      const slotDurationMinutes = (newEnd.getTime() - newStart.getTime()) / 60000
      const isAllowedSlotDuration = allowedSlotDurationsMinutes.includes(slotDurationMinutes)
      const isAlignedSlotStart =
        allowedStartMinutes.includes(newStart.getMinutes()) && newStart.getSeconds() === 0

      if (!isAllowedSlotDuration || !isAlignedSlotStart) {
        await networkRoute.fulfill({
          status: 400,
          contentType: "application/json",
          body: JSON.stringify({ error: "appointment slot must be 30 or 45 minutes and start aligned to a valid slot time" }),
        })
        return
      }

      const conflictingAppointment = currentAppointments.find(
        (appointment) =>
          appointment.staff_id === submittedJSON.staff_id &&
          appointment.status !== "cancelled" &&
          hasTimeOverlap(newStart, newEnd, new Date(appointment.starts_at as string), new Date(appointment.ends_at as string))
      )

      if (conflictingAppointment) {
        await networkRoute.fulfill({
          status: 409,
          contentType: "application/json",
          body: JSON.stringify({ error: "appointment time slot conflicts with an existing appointment" }),
        })
        return
      }

      const newAppointment = {
        id: `appt-${currentAppointments.length + 1}`,
        patient_fhir_id: submittedJSON.patient_fhir_id,
        staff_id: submittedJSON.staff_id,
        starts_at: submittedJSON.starts_at,
        ends_at: submittedJSON.ends_at,
        status: "scheduled",
        reason: submittedJSON.reason ?? "",
        version: 1,
        created_at: new Date().toISOString(),
      }
      currentAppointments.push(newAppointment)
      await networkRoute.fulfill({
        status: 201,
        contentType: "application/json",
        body: JSON.stringify(newAppointment),
      })
    }
  })

  await pageInstance.route("**/api/v1/appointments/my", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(currentAppointments),
    })
  })

  await pageInstance.route("**/api/v1/appointments/*/cancel", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = new URL(httpRequest.url())
    const urlParts = requestURL.pathname.split("/")
    const appointmentId = urlParts[urlParts.length - 2]

    const matchedAppointment = currentAppointments.find((appointment) => appointment.id === appointmentId)
    if (matchedAppointment) {
      matchedAppointment.status = "cancelled"
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(matchedAppointment),
      })
    } else {
      await networkRoute.fulfill({
        status: 404,
        contentType: "application/json",
        body: JSON.stringify({ error: "Appointment not found." }),
      })
    }
  })
}

export const mockPortalAPI = async (pageInstance: Page): Promise<void> => {
  const portalDashboardPayload = {
    patient_info: {
      fhir_resource_id: "fhir-pat-1",
      full_name: "Guilherme de Souza Araujo",
      birth_date: "1985-04-12",
      document_id: "123.456.789-00",
    },
    upcoming_encounters: [
      {
        fhir_resource_id: "enc-2",
        status: "finished",
        reason_display: "Retorno Cardiológico",
        started_at: "2026-05-15T14:30:00Z",
        ended_at: "2026-05-15T14:45:00Z",
      },
    ],
    recent_observations: [
      {
        fhir_resource_id: "obs-5",
        code_display: "Pressão Arterial Sistólica",
        loinc_code: "85354-9",
        value_quantity: 135,
        value_unit: "mmHg",
        observed_at: "2026-05-15T14:35:00Z",
      },
    ],
    active_conditions: [
      {
        fhir_resource_id: "cond-1",
        code_display: "Hipertensão Essencial Primária",
        icd10_code: "I10",
        clinical_status: "active",
        onset_at: "2026-05-15T14:40:00Z",
      },
    ],
    active_medications: [],
    recent_reports: [
      {
        fhir_resource_id: "rep-1",
        report_display: "Eletrocardiograma de Repouso",
        status: "final",
        conclusion: "Ritmo sinusal com leve taquicardia. Recomenda-se acompanhamento ambulatorial.",
        version: "2",
        issued_at: "2026-05-15T14:45:00Z",
      },
    ],
    recent_imaging: [
      {
        fhir_resource_id: "study-1",
        title: "Tomografia Computadorizada de Tórax",
        modality: "CT",
        status: "completed",
        created_at: "2026-05-16T10:00:00Z",
      },
    ],
  }

  await pageInstance.route("**/api/v1/portal/dashboard", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(portalDashboardPayload),
    })
  })

  await pageInstance.route("**/api/v1/portal/reports", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(portalDashboardPayload.recent_reports),
    })
  })

  await pageInstance.route("**/api/v1/portal/appointments", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify([]),
    })
  })
}

export const loginAsPatient = async (pageInstance: Page): Promise<void> => {
  await mockAuthAPI(pageInstance, "PATIENT", {
    email: "guilherme.paciente@hospital.com",
    password: "paciente123",
  })
  await mockPatientsAPI(pageInstance)
  await mockClinicalAPI(pageInstance)
  await mockNotificationsAPI(pageInstance, [
    {
      id: "notif-report-ready-1",
      type: "report_ready",
      priority: "high",
      title: "Laudo Pronto - Hemograma Completo",
      body: "Seu laudo do exame Hemograma Completo está disponível.",
      resource_type: "report",
      resource_id: "rep-1",
      is_read: false,
      created_at: new Date().toISOString(),
    },
  ])
  await mockPortalAPI(pageInstance)

  await pageInstance.goto("/login")
  await pageInstance.getByPlaceholder("nome.sobrenome@hospital.com").fill("guilherme.paciente@hospital.com")
  await pageInstance.getByPlaceholder("••••••••").fill("paciente123")
  await pageInstance.getByRole("button", { name: "Entrar no Console" }).click()
  await expect(pageInstance).toHaveURL(/\/portal/)
}

export const loginAsAdmin = async (pageInstance: Page): Promise<void> => {
  let adminSessionLoggedIn = false

  await pageInstance.route("**/api/v1/auth/login", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const submittedJSON = httpRequest.postDataJSON()

    if (submittedJSON.email === "admin@hospital.com" && submittedJSON.password === "admin123") {
      adminSessionLoggedIn = true
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          token: "mock-jwt-token-admin-789",
          userId: "user-admin-789",
          role: "ADMIN",
          email: "admin@hospital.com",
        }),
      })
    } else {
      await networkRoute.fulfill({
        status: 401,
        contentType: "application/json",
        body: JSON.stringify({
          error: "Credenciais inválidas.",
        }),
      })
    }
  })

  await pageInstance.route("**/api/v1/auth/logout", async (networkRoute) => {
    adminSessionLoggedIn = false
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        success: true,
      }),
    })
  })

  await pageInstance.route("**/api/v1/auth/me", async (networkRoute) => {
    if (adminSessionLoggedIn) {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          token: "mock-jwt-token-admin-789",
          userId: "user-admin-789",
          role: "ADMIN",
          email: "admin@hospital.com",
          fullName: "Administrador do Sistema",
          isActive: true,
        }),
      })
      return
    }
    await networkRoute.fulfill({
      status: 401,
      contentType: "application/json",
      body: JSON.stringify({
        error: "Não autenticado.",
      }),
    })
  })

  await mockPatientsAPI(pageInstance)
  await mockClinicalAPI(pageInstance)
  await mockAnalyzerAPI(pageInstance)
  await mockStaffAPI(pageInstance)
  await mockScheduleAPI(pageInstance)
  await mockTelemetryAPI(pageInstance)
  await mockAnalyticsAPI(pageInstance)
  await mockAuditLogsAPI(pageInstance)
  await mockNotificationsAPI(pageInstance)

  await pageInstance.goto("/login")
  await pageInstance.getByPlaceholder("nome.sobrenome@hospital.com").fill("admin@hospital.com")
  await pageInstance.getByPlaceholder("••••••••").fill("admin123")
  await pageInstance.getByRole("button", { name: "Entrar no Console" }).click()
  await expect(pageInstance).toHaveURL(/\/$/)
}

export const mockNotificationsAPI = async (
  pageInstance: Page,
  initialNotifications: Record<string, unknown>[] = []
): Promise<void> => {
  const currentNotifications: Record<string, unknown>[] = initialNotifications.length > 0
    ? initialNotifications
    : [
        {
          id: "notif-1",
          type: "system",
          priority: "low",
          title: "Bem-vindo ao sistema",
          body: "Seu cadastro foi realizado com sucesso.",
          resource_type: "",
          resource_id: "",
          is_read: false,
          created_at: new Date().toISOString(),
        },
        {
          id: "notif-2",
          type: "telemetry_alert",
          priority: "critical",
          title: "Alerta Crítico - Leito 01",
          body: "Paciente apresenta sinais vitais instáveis.",
          resource_type: "bed",
          resource_id: "bed-1",
          is_read: false,
          created_at: new Date().toISOString(),
        },
      ]

  await pageInstance.route("**/api/v1/notifications*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    if (httpRequest.method() === "GET") {
      const unreadOnlyNotifications = currentNotifications.filter((notification) => !notification.is_read)
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ notifications: unreadOnlyNotifications, total: unreadOnlyNotifications.length }),
      })
    }
  })

  await pageInstance.route("**/api/v1/notifications/unread-count", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    if (httpRequest.method() === "GET") {
      const unreadCount = currentNotifications.filter((notification) => notification.is_read === false).length
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ count: unreadCount }),
      })
    }
  })

  await pageInstance.route("**/api/v1/notifications/*/read", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const notificationId = urlParts[urlParts.length - 2]
    if (httpRequest.method() === "POST") {
      const notification = currentNotifications.find((notification) => notification.id === notificationId)
      if (notification) {
        notification.is_read = true
      }
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ success: true }),
      })
    }
  })

  await pageInstance.route("**/api/v1/notifications/stream", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "text/event-stream",
      headers: {
        "Cache-Control": "no-cache",
        Connection: "keep-alive",
      },
      body: "",
    })
  })
}

export const loginAsDoctor = async (pageInstance: Page): Promise<void> => {
  await mockAuthAPI(pageInstance)
  await mockPatientsAPI(pageInstance)
  await mockClinicalAPI(pageInstance)
  await mockAnalyzerAPI(pageInstance)
  await mockStaffAPI(pageInstance)
  await mockScheduleAPI(pageInstance)
  await mockTelemetryAPI(pageInstance)
  await mockAnalyticsAPI(pageInstance)
  await mockNotificationsAPI(pageInstance)
  await pageInstance.goto("/login")
  await pageInstance.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
  await pageInstance.getByPlaceholder("••••••••").fill("senha123")
  await pageInstance.getByRole("button", { name: "Entrar no Console" }).click()
  await expect(pageInstance).toHaveURL(/\/$/)
}

