import { test, expect } from "@playwright/test"
import { mockAuthAPI } from "./helpers"

test.describe("Internationalization (i18n) - Locale Switching", () => {
  test.beforeEach(async ({ page }) => {
    await mockAuthAPI(page)
  })

  test("should successfully switch between pt-BR, en-US, and es-ES", async ({ page }) => {
    await page.goto("/login")
    await page.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
    await page.getByPlaceholder("••••••••").fill("senha123")
    await page.getByRole("button", { name: "Entrar no Console" }).click()

    await expect(page).toHaveURL(/\/$/)

    const portuguesePatientsLabel = page.getByRole("button", { name: "Pacientes" })
    await expect(portuguesePatientsLabel).toBeVisible()

    const portugueseTelemetryLabel = page.getByRole("button", { name: "Telemetria UTI" })
    await expect(portugueseTelemetryLabel).toBeVisible()

    await page.getByRole("button", { name: "Português" }).click()
    await page.getByRole("button", { name: "English" }).click()

    const englishPatientsLabel = page.getByRole("button", { name: "Patients" })
    await expect(englishPatientsLabel).toBeVisible()

    const englishTelemetryLabel = page.getByRole("button", { name: "ICU Telemetry" })
    await expect(englishTelemetryLabel).toBeVisible()

    await page.getByRole("button", { name: "English" }).click()
    await page.getByRole("button", { name: "Español" }).click()

    const spanishPatientsLabel = page.getByRole("button", { name: "Pacientes" })
    await expect(spanishPatientsLabel).toBeVisible()

    const spanishTelemetryLabel = page.getByRole("button", { name: "Telemetría UCI" })
    await expect(spanishTelemetryLabel).toBeVisible()
  })

  test("should format dates according to the selected locale", async ({ page }) => {
    await page.route("**/api/v1/staff/employees*", async (networkRoute) => {
      await networkRoute.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify([]) })
    })
    await page.route("**/api/v1/appointments*", async (networkRoute) => {
      await networkRoute.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify([]) })
    })

    await page.goto("/login")
    await page.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
    await page.getByPlaceholder("••••••••").fill("senha123")
    await page.getByRole("button", { name: "Entrar no Console" }).click()
    await expect(page).toHaveURL(/\/$/)

    await page.goto("/schedule")

    const today = new Date()
    const expectedPtBrDate = new Intl.DateTimeFormat("pt-BR", {
      day: "numeric",
      month: "numeric",
      year: "numeric",
    }).format(today)
    await expect(page.getByRole("heading", { name: `Agenda de ${expectedPtBrDate}` })).toBeVisible()

    await page.getByRole("button", { name: "Português" }).click()
    await page.getByRole("button", { name: "English" }).click()

    const expectedEnUsDate = new Intl.DateTimeFormat("en-US", {
      day: "numeric",
      month: "numeric",
      year: "numeric",
    }).format(today)
    await expect(page.getByRole("heading", { name: `Schedule for ${expectedEnUsDate}` })).toBeVisible()
  })
})
