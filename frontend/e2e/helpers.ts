import { type Page, expect } from "@playwright/test"

export const mockAuthAPI = async (pageInstance: Page): Promise<void> => {
  await pageInstance.route("**/api/auth/login", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const submittedJSON = httpRequest.postDataJSON()

    if (submittedJSON.email === "medico@clinica.com" && submittedJSON.password === "senha123") {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          token: "mock-jwt-token-123456",
          userId: "user-medico-123",
          role: "doctor",
          email: "medico@clinica.com",
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

  await pageInstance.route("**/api/auth/logout", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        success: true,
      }),
    })
  })

  await pageInstance.route("**/api/auth/me", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        userId: "user-medico-123",
        role: "doctor",
        email: "medico@clinica.com",
        fullName: "Dr. André Silva de Araujo",
        isActive: true,
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

  await pageInstance.route("**/api/patients", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    if (httpRequest.method() === "GET") {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(currentPatientsList),
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

  await pageInstance.route("**/api/patients/*", async (networkRoute) => {
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

  await pageInstance.route("**/api/patients/*/encounters", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const patientFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentEncountersList.filter((e) => e.patient_fhir_id === patientFhirId)
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
        status: "finished",
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

  await pageInstance.route("**/api/patients/*/observations", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const patientFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentObservationsList.filter((o) => o.patient_fhir_id === patientFhirId)
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(filtered),
      })
    }
  })

  await pageInstance.route("**/api/patients/*/conditions", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const patientFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentConditionsList.filter((c) => c.patient_fhir_id === patientFhirId)
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

  await pageInstance.route("**/api/encounters/*/observations", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const encounterFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentObservationsList.filter((o) => o.encounter_fhir_id === encounterFhirId)
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

  await pageInstance.route("**/api/encounters/*/reports", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const encounterFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentReportsList.filter((r) => r.encounter_fhir_id === encounterFhirId)
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

  await pageInstance.route("**/api/encounters/*/medications", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const encounterFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentMedicationsList.filter((m) => m.encounter_fhir_id === encounterFhirId)
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
        medication_display: submittedJSON.medication_display,
        dosage_instruction: submittedJSON.dosage_instruction,
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

  await pageInstance.route("**/api/patients/*/studies", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const patientFhirId = urlParts[urlParts.length - 2]

    if (httpRequest.method() === "GET") {
      const filtered = currentStudiesList.filter((s) => s.patient_fhir_id === patientFhirId)
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

  await pageInstance.route("**/api/studies/*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const studyId = urlParts[urlParts.length - 1]

    if (httpRequest.method() === "GET") {
      const study = currentStudiesList.find((s) => s.id === studyId)
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

  await pageInstance.route("**/api/exam-analyses", async (networkRoute) => {
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

  await pageInstance.route("**/api/exam-analyses/*", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const requestURL = httpRequest.url()
    const urlParts = requestURL.split("/")
    const anaId = urlParts[urlParts.length - 1]

    if (httpRequest.method() === "GET") {
      const analysis = currentAnalysesList.find((a) => a.id === anaId)
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
      const index = currentAnalysesList.findIndex((a) => a.id === anaId)
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
      fullName: "Dr. André Silva de Araujo",
      email: "andre.silva@hospital.com",
      role: "doctor",
      crmNumber: "CRM-SP 12345",
      status: "active",
      department: "Cardiologia",
    },
    {
      id: "emp-2",
      userId: "user-2",
      fullName: "Enf. Roberta Santos Almeida",
      email: "roberta.santos@hospital.com",
      role: "nurse",
      crmNumber: "COREN-SP 54321",
      status: "active",
      department: "Pediatria",
    },
  ]

  await pageInstance.route("**/api/staff/employees", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    if (httpRequest.method() === "GET") {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(currentEmployees),
      })
    } else if (httpRequest.method() === "POST") {
      const submittedJSON = httpRequest.postDataJSON()
      const newEmployee = {
        id: `emp-${currentEmployees.length + 1}`,
        userId: submittedJSON.user_id || `user-${currentEmployees.length + 1}`,
        fullName: submittedJSON.full_name,
        email: submittedJSON.email,
        role: submittedJSON.role,
        crmNumber: submittedJSON.crm_number || "N/A",
        status: "active",
        department: "Geral",
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

  await pageInstance.route("**/api/telemetry/rooms", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify(roomsList),
    })
  })

  await pageInstance.route("**/api/telemetry/rooms/*/unlock", async (networkRoute) => {
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

  await pageInstance.route("**/api/telemetry/rooms/*/beds", async (networkRoute) => {
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

  await pageInstance.route("**/api/telemetry/beds/*/condition", async (networkRoute) => {
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

export const mockStatsAPI = async (pageInstance: Page): Promise<void> => {
  await pageInstance.route("**/api/stats", async (networkRoute) => {
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
          { dayName: "stats.days.mon", count: 8 },
          { dayName: "stats.days.tue", count: 12 },
          { dayName: "stats.days.wed", count: 14 },
          { dayName: "stats.days.thu", count: 11 },
          { dayName: "stats.days.fri", count: 15 },
          { dayName: "stats.days.sat", count: 5 },
          { dayName: "stats.days.sun", count: 2 }
        ],
        pathologies: [
          { code: "J45.9", descriptionKey: "stats.pathologies.asthma", categoryKey: "stats.categories.respiratory", activeCases: 44, trend: "+5%" },
          { code: "I10", descriptionKey: "stats.pathologies.hypertension", categoryKey: "stats.categories.cardiovascular", activeCases: 119, trend: "stable" },
          { code: "E11.9", descriptionKey: "stats.pathologies.diabetes", categoryKey: "stats.categories.endocrine", activeCases: 85, trend: "+12%" }
        ]
      })
    })
  })
}

export const loginAsDoctor = async (pageInstance: Page): Promise<void> => {
  await mockAuthAPI(pageInstance)
  await mockPatientsAPI(pageInstance)
  await mockClinicalAPI(pageInstance)
  await mockAnalyzerAPI(pageInstance)
  await mockStaffAPI(pageInstance)
  await mockTelemetryAPI(pageInstance)
  await mockStatsAPI(pageInstance)
  await pageInstance.goto("/login")
  await pageInstance.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
  await pageInstance.getByPlaceholder("••••••••").fill("senha123")
  await pageInstance.getByRole("button", { name: "Entrar no Console" }).click()
  await expect(pageInstance).toHaveURL(/\/$/)
}

