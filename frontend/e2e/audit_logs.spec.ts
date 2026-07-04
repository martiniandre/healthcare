import { test, expect } from "@playwright/test"
import { loginAsAdmin } from "./helpers"

test.describe("Audit Logs Management Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page)
    await page.goto("/audit-logs")
  })

  test("should display audit logs page title and table with rows", async ({ page }) => {
    await expect(page.locator("text=Registros de Auditoria")).toBeVisible()
    await expect(page.locator("text=admin@hospital.com")).toBeVisible()
    await expect(page.locator("text=medico@clinica.com")).toBeVisible()
    await expect(page.locator("text=Sucesso")).toBeVisible()
    await expect(page.locator("text=Falha")).toBeVisible()
  })

  test("should show audit log details when expanding a row", async ({ page }) => {
    await page.locator("text=admin@hospital.com").click()
    await expect(page.locator("text=corr-001")).toBeVisible()
    await expect(page.locator("text=log-1")).toBeVisible()
  })

  test("should filter audit logs by status dropdown", async ({ page }) => {
    const statusSelect = page.locator("select").nth(1)
    await statusSelect.selectOption("FAILURE")
    await expect(page.locator("text=usuario.invalido@test.com")).toBeVisible()
  })
})
