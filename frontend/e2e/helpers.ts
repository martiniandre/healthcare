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

export const loginAsDoctor = async (pageInstance: Page): Promise<void> => {
  await mockAuthAPI(pageInstance)
  await mockPatientsAPI(pageInstance)
  await pageInstance.goto("/#/login")
  await pageInstance.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
  await pageInstance.getByPlaceholder("••••••••").fill("senha123")
  await pageInstance.getByRole("button", { name: "Entrar no Console" }).click()
  await expect(pageInstance).toHaveURL(/.*#\/$/)
}
