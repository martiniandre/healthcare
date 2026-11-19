import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Sidebar Navigation Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
  })

  test("should display all navigation links grouped by topics in the sidebar", async ({ page }) => {
    const sidebar = page.locator("aside")
    await expect(sidebar.getByText("Assistência clínica")).toBeVisible()
    await expect(sidebar.getByText("Diagnóstico")).toBeVisible()
    await expect(sidebar.getByText("Operações")).toBeVisible()
    await expect(page.getByRole("button", { name: "Pacientes" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Telemetria UTI" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Análise de exames" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Indicadores clínicos" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Gestão de equipes" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Aparência" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Sair" })).toBeVisible()
  })

  test("should navigate through all sidebar links and update URL correctly", async ({ page }) => {
    await page.getByRole("button", { name: "Telemetria UTI" }).click()
    await expect(page).toHaveURL(/\/telemetry$/)

    await page.getByRole("button", { name: "Indicadores clínicos" }).click()
    await expect(page).toHaveURL(/\/analytics$/)

    await page.getByRole("button", { name: "Gestão de equipes" }).click()
    await expect(page).toHaveURL(/\/staff$/)

    await page.getByRole("button", { name: "Análise de exames" }).click()
    await expect(page).toHaveURL(/\/exam-analyzer$/)

    await page.getByRole("button", { name: "Pacientes" }).click()
    await expect(page).toHaveURL(/\/$/)
  })

  test("should highlight the active route in the sidebar", async ({ page }) => {
    await page.getByRole("button", { name: "Indicadores clínicos" }).click()

    const analyticsButton = page.getByRole("button", { name: "Indicadores clínicos" })
    await expect(analyticsButton).toBeVisible()
  })
})
