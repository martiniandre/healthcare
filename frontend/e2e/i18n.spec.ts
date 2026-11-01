import { test, expect, type Page } from "@playwright/test"
import { mockAuthAPI } from "./helpers"

const switchLanguage = async (page: Page, currentLabel: string, targetLabel: string) => {
  await page.getByRole("button", { name: currentLabel, exact: true }).click()
  await page.getByRole("menuitemradio", { name: targetLabel }).click()
}

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

    await expect(page.getByRole("button", { name: "Pacientes" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Telemetria UTI" })).toBeVisible()

    await switchLanguage(page, "Português", "English")
    await expect(page.getByRole("button", { name: "Patients" })).toBeVisible()
    await expect(page.getByRole("button", { name: "ICU Telemetry" })).toBeVisible()

    await switchLanguage(page, "English", "Español")
    await expect(page.getByRole("button", { name: "Pacientes" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Telemetría UCI" })).toBeVisible()
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

    await expect(page.getByRole("button", { name: "Agenda", exact: true })).toBeVisible()
    await expect(page.locator(".fc-col-header").first()).toContainText(/seg/i)

    await switchLanguage(page, "Português", "English")

    await expect(page.getByRole("button", { name: "Schedule", exact: true })).toBeVisible()
  })
})
