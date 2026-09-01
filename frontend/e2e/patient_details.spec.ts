import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Patient Record and Clinical File Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/patients/fhir-pat-1")
  })

  test("should render the patient's initial clinical file", async ({ page }) => {
    const patientHeaderName = page.locator("h2", { hasText: "Guilherme de Souza Araujo" })
    await expect(patientHeaderName).toBeVisible()

    const initialEncounter = page.locator("text=Consulta de Rotina Geral")
    await expect(initialEncounter).toBeVisible()
  })

  test("should format the header birth date according to the selected language", async ({ page }) => {
    const expectedBirthDate = new Intl.DateTimeFormat("pt-BR", {
      day: "numeric",
      month: "numeric",
      year: "numeric",
    }).format(new Date(1988, 3, 12))

    await expect(page.getByText(expectedBirthDate)).toBeVisible()
  })

  test("should create a new medical encounter", async ({ page }) => {
    await page.getByRole("button", { name: "Nova Consulta" }).click()
    await page.getByPlaceholder("Ex: Consulta de rotina cardiologia, Check-up anual").fill("Checkup Anual Geral")

    await page.getByRole("combobox").click()
    await page.getByRole("option", { name: "Dr. André Silva de Araujo" }).click()
    await page.getByRole("button", { name: "Confirmar Consulta" }).click()

    await expect(page).toHaveURL(/\/patients\/fhir-pat-1\/encounters\/enc-3/)
    await expect(page.getByText("Sinais Vitais (Observations)")).toBeVisible()
  })

  test("should add a new vital sign observation", async ({ page }) => {
    await page.goto("/patients/fhir-pat-1/encounters/enc-1")
    await page.getByRole("button", { name: "Adicionar Sinal" }).click()

    await page.getByRole("combobox").click()
    await page.getByRole("option", { name: "Temperatura Corporal (°C)" }).click()

    await page.getByPlaceholder("Ex: 72, 36.5, 120").fill("38.5")
    await page.getByRole("button", { name: "Registrar Sinal" }).click()

    const temperatureValue = page.locator("text=38.5")
    await expect(temperatureValue).toBeVisible()
    await expect(page.locator("text=Temperatura Corporal").first()).toBeVisible()
  })

  test("should sign a new diagnostic report", async ({ page }) => {
    await page.goto("/patients/fhir-pat-1/encounters/enc-1?tab=reports")
    await page.getByRole("button", { name: "Novo Laudo" }).click()

    await page.getByRole("combobox").click()
    await page.getByRole("option", { name: "Hemograma Completo (58410-2)" }).click()

    await page.getByPlaceholder("Ex: Hemograma Completo, Radiografia de Tórax PA").fill("Hemograma Completo de Controle")
    await page.getByPlaceholder("Descreva as conclusões e observações clínicas do exame...").fill("Contagem de hemácias dentro do padrão de referência. Leucograma normal.")
    await page.getByRole("button", { name: "Salvar Laudo" }).click()

    const reportTitleText = page.locator("text=Hemograma Completo de Controle")
    const reportConclusionText = page.locator("text=Contagem de hemácias dentro do padrão de referência. Leucograma normal.")
    await expect(reportTitleText).toBeVisible()
    await expect(reportConclusionText).toBeVisible()
  })

  test("should create a medical encounter with in-progress status", async ({ page }) => {
    await page.getByRole("button", { name: "Nova Consulta" }).click()
    await page.getByPlaceholder("Ex: Consulta de rotina cardiologia, Check-up anual").fill("Checkup Anual Geral")

    await page.getByRole("combobox").click()
    await page.getByRole("option", { name: "Dr. André Silva de Araujo" }).click()
    await page.getByRole("button", { name: "Confirmar Consulta" }).click()

    await expect(page).toHaveURL(/\/patients\/fhir-pat-1\/encounters\/enc-3/)

    await page.goto("/patients/fhir-pat-1")
    const newEncounterRow = page.locator("tr", { hasText: "Checkup Anual Geral" })
    await expect(newEncounterRow).toBeVisible()
    await expect(newEncounterRow).toContainText("in-progress")
  })

  test("should finish an in-progress encounter", async ({ page }) => {
    await page.getByRole("button", { name: "Nova Consulta" }).click()
    await page.getByPlaceholder("Ex: Consulta de rotina cardiologia, Check-up anual").fill("Checkup Anual Geral")

    await page.getByRole("combobox").click()
    await page.getByRole("option", { name: "Dr. André Silva de Araujo" }).click()
    await page.getByRole("button", { name: "Confirmar Consulta" }).click()

    await expect(page).toHaveURL(/\/patients\/fhir-pat-1\/encounters\/enc-3/)

    await page.goto("/patients/fhir-pat-1")
    const newEncounterRow = page.locator("tr", { hasText: "Checkup Anual Geral" })
    await expect(newEncounterRow).toContainText("in-progress")
    await newEncounterRow.getByRole("button", { name: "Finalizar" }).click()
    await expect(newEncounterRow).toContainText("finished")
  })
})

test.describe("Patient Record Error States", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
  })

  test("should show an error state instead of infinite loading when the patient record fails to load", async ({ page }) => {
    await page.route("**/api/v1/patients/*", async (networkRoute) => {
      await networkRoute.fulfill({
        status: 500,
        contentType: "application/json",
        body: JSON.stringify({ error: "Internal Server Error" }),
      })
    })

    await page.goto("/patients/fhir-pat-1")

    await expect(
      page.getByText("Não foi possível carregar a ficha clínica. Verifique se o paciente existe ou tente novamente.")
    ).toBeVisible()
  })
})
