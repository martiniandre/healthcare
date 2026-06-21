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
})
