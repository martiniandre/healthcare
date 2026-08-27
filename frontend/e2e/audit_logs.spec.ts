import { test, expect } from "@playwright/test"
import { loginAsAdmin } from "./helpers"

test.describe("Audit Logs Management Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page)
    await page.goto("/audit-logs")
  })

  test("should display audit logs page title and table with rows", async ({ page }) => {
    await expect(page.locator("text=Registros de Auditoria")).toBeVisible()
    const logsTable = page.getByRole("table")
    await expect(logsTable.getByText("admin@hospital.com")).toBeVisible()
    await expect(logsTable.getByText("medico@clinica.com")).toBeVisible()
    await expect(logsTable.getByText("Sucesso").first()).toBeVisible()
    await expect(logsTable.getByText("Falha")).toBeVisible()
  })

  test("should show audit log details when expanding a row", async ({ page }) => {
    await page.getByRole("table").getByText("admin@hospital.com").click()
    await expect(page.getByText("corr-001")).toBeVisible()
    await expect(page.getByText("log-1")).toBeVisible()
  })

  test("should filter audit logs by status dropdown", async ({ page }) => {
    const statusSelect = page.getByRole("combobox").nth(1)
    await statusSelect.click()
    await page.getByRole("listbox").getByText("Falha", { exact: true }).click()
    await expect(page.getByRole("table").getByText("usuario.invalido@test.com")).toBeVisible()
  })
})
