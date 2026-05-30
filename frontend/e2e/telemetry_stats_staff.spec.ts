import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Real-time Telemetry Dashboard Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/#/telemetry")
  })

  test("should display telemetry rooms and default unlocked room status", async ({ page }) => {
    const greenRoomName = page.getByRole("heading", { name: "Sala Verde - Semi-Intensiva" })
    await expect(greenRoomName).toBeVisible()

    const redRoomName = page.getByRole("heading", { name: "Sala Vermelha - Choque & Emergência" })
    await expect(redRoomName).toBeVisible()

    const unlockedBadge = page.locator("text=Desbloqueada").first()
    await expect(unlockedBadge).toBeVisible()
  })

  test("should handle room unlocking workflow with passcode verification", async ({ page }) => {
    const redRoomCard = page.getByRole("heading", { name: "Sala Vermelha - Choque & Emergência" })
    await redRoomCard.click()

    const passwordInput = page.getByPlaceholder("Digite o código da sala...")
    const confirmButton = page.getByRole("button", { name: "Liberar Monitoramento" })

    await passwordInput.fill("1111")
    await confirmButton.click()

    const errorMessage = page.locator("text=Senha de Acesso incorreta. Verifique a escala do plantão.")
    await expect(errorMessage).toBeVisible()

    await passwordInput.fill("9999")
    await confirmButton.click()

    await expect(errorMessage).not.toBeVisible()
    await page.getByText("Leito 02").click()
    const activeBedPatientName = page.getByRole("heading", { name: "Mariana Costa Silva", exact: true })
    await expect(activeBedPatientName).toBeVisible()
  })

  test("should toggle telemetry audible alarms status", async ({ page }) => {
    const alarmsButton = page.getByRole("button", { name: "Alarmes Silenciados" })
    await expect(alarmsButton).toBeVisible()

    await alarmsButton.click()
    const activeAlarmsButton = page.getByRole("button", { name: "Alarme Sonoro Ativo" })
    await expect(activeAlarmsButton).toBeVisible()
  })
})

test.describe("Clinical Statistics and Analytics Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/#/stats")
  })

  test("should load diagnostic statistics and key KPIs", async ({ page }) => {
    const activePatientsCount = page.locator("text=340")
    await expect(activePatientsCount).toBeVisible()

    const fhirComplianceRate = page.locator("text=99.4%")
    await expect(fhirComplianceRate).toBeVisible()

    const averageDurationTime = page.locator("text=14.5 min")
    await expect(averageDurationTime).toBeVisible()
  })

  test("should load weekly consultations volume chart data", async ({ page }) => {
    const weeklyChartHeading = page.locator("text=Volume de Atendimentos")
    await expect(weeklyChartHeading).toBeVisible()

    const chartDaySeg = page.getByText("Seg", { exact: true })
    const chartDaySex = page.getByText("Sex", { exact: true })
    await expect(chartDaySeg).toBeVisible()
    await expect(chartDaySex).toBeVisible()
  })

  test("should load and list clinical pathlogies classification table", async ({ page }) => {
    const pathTitle = page.locator("text=Asma não especificada")
    const pathCode = page.locator("text=J45.9")
    await expect(pathTitle).toBeVisible()
    await expect(pathCode).toBeVisible()
  })
})

test.describe("Hospital Staff Management Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/#/staff")
  })

  test("should render the staff members list page", async ({ page }) => {
    const doctorName = page.locator("text=Dr. André Silva de Araujo")
    await expect(doctorName).toBeVisible()

    const nurseName = page.locator("text=Enf. Roberta Santos Almeida")
    await expect(nurseName).toBeVisible()
  })

  test("should filter staff members list by role tab selection", async ({ page }) => {
    const searchField = page.getByPlaceholder("Buscar por nome, e-mail ou especialidade...")
    await searchField.fill("Roberta")

    const matchingNurse = page.locator("text=Enf. Roberta Santos Almeida")
    const nonMatchingDoctor = page.locator("text=Dr. André Silva de Araujo")

    await expect(matchingNurse).toBeVisible()
    await expect(nonMatchingDoctor).not.toBeVisible()
  })

  test("should successfully register a new hospital practitioner", async ({ page }) => {
    await page.getByRole("button", { name: "Cadastrar Profissional" }).click()

    await page.getByPlaceholder("Ex: Dr. André Silva de Araujo").fill("Dra. Paula Albuquerque")
    await page.locator("select").selectOption("Médico")
    await page.getByPlaceholder("Ex: CRM-SP 12345").fill("CRM-SP 777777")
    await page.getByPlaceholder("Ex: nome@hospital.com").fill("paula.albuquerque@hospital.com")
    await page.getByPlaceholder("Ex: Cardiologia").fill("Neurologia")

    await page.getByRole("button", { name: "Salvar Cadastro" }).click()

    const newPractitionerName = page.locator("text=Dra. Paula Albuquerque")
    await expect(newPractitionerName).toBeVisible()
  })
})
