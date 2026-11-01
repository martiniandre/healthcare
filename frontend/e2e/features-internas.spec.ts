import { test, expect, type Page } from "@playwright/test"
import { loginAsAdmin } from "./helpers"

const expectModuleLoaded = async (page: Page, routePattern: RegExp) => {
  await expect(page).toHaveURL(routePattern)
  await expect(page.getByRole("main")).toBeVisible()
  await expect(page.getByRole("heading").first()).toBeVisible()
}

test.describe("Internal features - authenticated admin", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page)
  })

  test("should open the patients list and navigate to a patient record", async ({ page }) => {
    await expectModuleLoaded(page, /\/$/)
    await page.getByRole("button", { name: "Prontuário" }).first().click()
    await expect(page).toHaveURL(/\/patients\/fhir-pat-1$/)
    await expect(page.getByRole("main")).toBeVisible()
  })

  test("should open the patient encounter detail", async ({ page }) => {
    await page.goto("/patients/fhir-pat-1/encounters/enc-2")
    await expect(page).toHaveURL(/\/patients\/fhir-pat-1\/encounters\/enc-2$/)
    await expect(page.getByRole("main")).toBeVisible()
  })

  test("should open the exam analyzer and list analyses", async ({ page }) => {
    await page.goto("/exam-analyzer")
    await expectModuleLoaded(page, /\/exam-analyzer$/)
    await expect(page.getByText("rx_torax.png").first()).toBeVisible()
  })

  test("should open telemetry and show the monitored beds", async ({ page }) => {
    await page.goto("/telemetry")
    await expectModuleLoaded(page, /\/telemetry$/)
    await expect(page.getByText("Sala Verde - Semi-Intensiva").first()).toBeVisible()
  })

  test("should open analytics", async ({ page }) => {
    await page.goto("/analytics")
    await expectModuleLoaded(page, /\/analytics$/)
  })

  test("should open staff management and list employees", async ({ page }) => {
    await page.goto("/staff")
    await expectModuleLoaded(page, /\/staff$/)
    await expect(page.getByText("Dr. André Silva de Araujo").first()).toBeVisible()
  })

  test("should open audit logs", async ({ page }) => {
    await page.goto("/audit-logs")
    await expectModuleLoaded(page, /\/audit-logs$/)
  })

  test("should show the schedule in the sidebar and open the agenda", async ({ page }) => {
    await page.goto("/")
    await page.getByRole("button", { name: /agenda|schedule/i }).first().click()
    await expect(page).toHaveURL(/\/schedule$/)
    await expect(page.getByText("Agenda", { exact: true }).first()).toBeVisible()
    await expect(page.getByRole("button", { name: "Novo Agendamento" })).toBeVisible()
  })

  test("should book and cancel an appointment from the agenda", async ({ page }) => {
    await page.goto("/schedule")
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    const modalDialog = page.getByRole("dialog")
    const appointmentDate = new Date()
    appointmentDate.setDate(appointmentDate.getDate() + 1)
    const isoDate = `${appointmentDate.getFullYear()}-${String(appointmentDate.getMonth() + 1).padStart(2, "0")}-${String(appointmentDate.getDate()).padStart(2, "0")}`

    await modalDialog.locator('input[type="date"]').fill(isoDate)
    await modalDialog.locator('[role="combobox"]').first().click()
    await page.locator('[role="option"]', { hasText: "Guilherme de Souza Araujo" }).click()
    await modalDialog.locator('[role="combobox"]').nth(1).click()
    await page.locator('[role="option"]', { hasText: "Dr. André Silva de Araujo" }).click()
    await modalDialog.getByText("Selecione o horário de início").click()
    await page.getByRole("option", { name: "09:00", exact: true }).click()
    await page.getByRole("button", { name: "Agendar" }).click()

    // Amanhã cai na semana seguinte quando é segunda-feira (a semana exibida começa na segunda).
    if (appointmentDate.getDay() === 1) {
      await page.locator(".fc-next-button").click()
    }

    const appointmentEvent = page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" }).first()
    await expect(appointmentEvent).toBeVisible()

    await appointmentEvent.click()
    await expect(page.getByRole("heading", { name: "Detalhes do Agendamento" })).toBeVisible()
    await page.getByRole("button", { name: "Cancelar" }).click()
    await expect(page.getByRole("heading", { name: "Cancelar Agendamento" })).toBeVisible()
    await page.getByRole("button", { name: "Confirmar Cancelamento" }).click()

    await expect(page.locator("text=Agendamento cancelado com sucesso!")).toBeVisible()
  })
})
