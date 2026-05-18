import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Patients Management Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
  })

  test("should render the initial list of patients", async ({ page }) => {
    const firstPatientName = page.locator("text=Guilherme de Souza Araujo")
    const secondPatientName = page.locator("text=Mariana Costa Silva")

    await expect(firstPatientName).toBeVisible()
    await expect(secondPatientName).toBeVisible()
  })

  test("should filter patients through the search field", async ({ page }) => {
    const searchField = page.getByPlaceholder("Buscar por nome, CPF ou telefone...")
    await searchField.fill("Guilherme")

    const matchingPatientName = page.locator("text=Guilherme de Souza Araujo")
    const nonMatchingPatientName = page.locator("text=Mariana Costa Silva")

    await expect(matchingPatientName).toBeVisible()
    await expect(nonMatchingPatientName).not.toBeVisible()

    await searchField.fill("")
    await expect(nonMatchingPatientName).toBeVisible()
  })

  test("should successfully register a new patient", async ({ page }) => {
    await page.getByRole("button", { name: "Novo Paciente" }).click()

    await page.getByPlaceholder("Nome Completo do Paciente").fill("Carlos Eduardo Rezende")
    await page.getByPlaceholder("AAAA-MM-DD").fill("1990-05-15")
    await page.getByPlaceholder("123.456.789-00").fill("111.222.333-44")
    await page.getByPlaceholder("(11) 98765-4321").fill("(11) 96543-2100")

    await page.getByRole("button", { name: "Confirmar Cadastro" }).click()

    const newPatientName = page.locator("text=Carlos Eduardo Rezende")
    await expect(newPatientName).toBeVisible()
  })

  test("should navigate to the patient details page when clicking the button", async ({ page }) => {
    const patientRow = page.locator("tr", { hasText: "Guilherme de Souza Araujo" })
    await patientRow.getByRole("button", { name: "Prontuário" }).click()

    await expect(page).toHaveURL(/.*#\/patients\/fhir-pat-1$/)
    const recordHeading = page.locator("text=Recursos Clínicos")
    await expect(recordHeading).toBeVisible()
  })
})
