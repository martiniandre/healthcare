import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Patient Record and Clinical File Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/#/patients/fhir-pat-1")
  })

  test("should render the patient's initial clinical file", async ({ page }) => {
    const patientHeaderName = page.locator("h2", { hasText: "Guilherme de Souza Araujo" })
    await expect(patientHeaderName).toBeVisible()

    const initialEncounter = page.locator("text=Consulta de Rotina Geral")
    await expect(initialEncounter).toBeVisible()
  })

  test("should create a new medical encounter", async ({ page }) => {
    await page.getByRole("button", { name: "Nova Consulta" }).click()
    await page.getByPlaceholder("Ex: Consulta de Rotina, Dor Abdominal").fill("Checkup Anual Geral")
    await page.getByRole("button", { name: "Registrar Consulta" }).click()

    const newEncounterName = page.locator("text=Checkup Anual Geral")
    await expect(newEncounterName).toBeVisible()
  })

  test("should add a new vital sign observation", async ({ page }) => {
    await page.getByRole("button", { name: "Sinais Vitais" }).click()
    await page.getByRole("button", { name: "Adicionar Sinal" }).click()

    await page.locator("select").selectOption("8310-5")
    await page.getByPlaceholder("Insira o valor numérico").fill("38.5")
    await page.getByRole("button", { name: "Gravar Métrica" }).click()

    const temperatureText = page.locator("text=Temperatura Corporal")
    const temperatureValue = page.locator("text=38.5")
    await expect(temperatureText).toBeVisible()
    await expect(temperatureValue).toBeVisible()
  })

  test("should sign a new diagnostic report", async ({ page }) => {
    await page.getByRole("button", { name: "Laudos Clínicos" }).click()
    await page.getByRole("button", { name: "Novo Laudo" }).click()

    await page.getByPlaceholder("Ex: Hemograma Completo, ECG").fill("Hemograma Completo de Controle")
    await page.getByPlaceholder("Redija a conclusão clínica circunstanciada...").fill("Contagem de hemácias dentro do padrão de referência. Leucograma normal.")
    await page.getByRole("button", { name: "Assinar Laudo" }).click()

    const reportTitleText = page.locator("text=Hemograma Completo de Controle")
    const reportConclusionText = page.locator("text=Contagem de hemácias dentro do padrão de referência. Leucograma normal.")
    await expect(reportTitleText).toBeVisible()
    await expect(reportConclusionText).toBeVisible()
  })
})
