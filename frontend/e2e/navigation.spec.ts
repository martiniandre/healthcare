import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Sidebar Navigation Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
  })

  test("should display all navigation links in the sidebar", async ({ page }) => {
    await expect(page.getByRole("button", { name: "Pacientes" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Telemetria UTI" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Análise de Exames" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Analytics" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Gestão de Equipes" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Sair" })).toBeVisible()
  })

  test("should navigate through all sidebar links and update URL correctly", async ({ page }) => {
    await page.getByRole("button", { name: "Telemetria UTI" }).click()
    await expect(page).toHaveURL(/\/telemetry$/)

    await page.getByRole("button", { name: "Analytics" }).click()
    await expect(page).toHaveURL(/\/analytics$/)

    await page.getByRole("button", { name: "Gestão de Equipes" }).click()
    await expect(page).toHaveURL(/\/staff$/)

    await page.getByRole("button", { name: "Análise de Exames" }).click()
    await expect(page).toHaveURL(/\/exam-analyzer$/)

    await page.getByRole("button", { name: "Pacientes" }).click()
    await expect(page).toHaveURL(/\/$/)
  })

  test("should highlight the active route in the sidebar", async ({ page }) => {
    await page.getByRole("button", { name: "Analytics" }).click()

    const analyticsButton = page.getByRole("button", { name: "Analytics" })
    await expect(analyticsButton).toBeVisible()
  })
})
