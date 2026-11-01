import { test, expect } from "@playwright/test"
import { loginAsAdmin } from "./helpers"

test.describe("Patients Management Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page)
  })

  test("should render the initial list of patients", async ({ page }) => {
    const firstPatientName = page.locator("text=Guilherme de Souza Araujo")
    const secondPatientName = page.locator("text=Mariana Costa Silva")

    await expect(firstPatientName).toBeVisible()
    await expect(secondPatientName).toBeVisible()
  })

  test("should format patient birth dates according to the selected language", async ({ page }) => {
    const guilhermeRow = page.locator("tr", { hasText: "Guilherme de Souza Araujo" })
    await expect(guilhermeRow).toContainText("12/04/1988")

    await page.getByRole("button", { name: "Português", exact: true }).click()
    await page.getByRole("menuitemradio", { name: "English" }).click()

    await expect(page.getByRole("button", { name: "English", exact: true })).toBeVisible()
    await expect(guilhermeRow).toContainText("1988")
    await expect(guilhermeRow).not.toContainText("12/04/1988")
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

  test("should hide pagination when all results fit in a single page", async ({ page }) => {
    const searchField = page.getByPlaceholder("Buscar por nome, CPF ou telefone...")
    await searchField.fill("Nogueira")

    await expect(page.locator("text=Rafael Nogueira Alves")).toBeVisible()
    await expect(page.getByRole("button", { name: "Próxima" })).not.toBeVisible()
    await expect(page.getByRole("button", { name: "Anterior" })).not.toBeVisible()
  })

  test("should paginate through the patient list showing distinct rows per page", async ({ page }) => {
    await expect(page.locator("text=Guilherme de Souza Araujo")).toBeVisible()
    await expect(page.locator("text=Thiago Costa Matos")).not.toBeVisible()

    await page.getByRole("button", { name: "Próxima" }).click()

    await expect(page.locator("text=Thiago Costa Matos")).toBeVisible()
    await expect(page.locator("text=Guilherme de Souza Araujo")).not.toBeVisible()

    await page.getByRole("button", { name: "Anterior" }).click()

    await expect(page.locator("text=Guilherme de Souza Araujo")).toBeVisible()
    await expect(page.locator("text=Thiago Costa Matos")).not.toBeVisible()
  })

  test("should successfully register a new patient", async ({ page }) => {
    await page.getByRole("button", { name: "Novo Paciente" }).click()

    await page.getByPlaceholder("Nome Completo do Paciente").fill("Carlos Eduardo Rezende")
    await page.getByPlaceholder("DD/MM/AAAA").fill("15/05/1990")
    await page.getByPlaceholder("123.456.789-00").fill("111.222.333-44")
    await page.getByPlaceholder("(11) 98765-4321").fill("(11) 96543-2100")

    await page.getByRole("button", { name: "Confirmar Cadastro" }).click()

    const newPatientName = page.locator("text=Carlos Eduardo Rezende")
    await expect(newPatientName).toBeVisible()
  })

  test("should navigate to the patient details page when clicking the button", async ({ page }) => {
    const patientRow = page.locator("tr", { hasText: "Guilherme de Souza Araujo" })
    await patientRow.getByRole("button", { name: "Prontuário" }).click()

    await expect(page).toHaveURL(/\/patients\/fhir-pat-1$/)
    const recordHeading = page.locator("text=Recursos Clínicos")
    await expect(recordHeading).toBeVisible()
  })
})
