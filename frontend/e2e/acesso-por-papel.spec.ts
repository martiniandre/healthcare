import { test, expect } from "@playwright/test"
import { loginAsDoctor, loginAsPatient } from "./helpers"

test.describe("Controle de acesso por papel", () => {
  test("médico não vê o item Logs de Auditoria no menu", async ({ page }) => {
    await loginAsDoctor(page)

    await expect(page.getByRole("button", { name: "Pacientes" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Logs de Auditoria" })).toHaveCount(0)
  })

  test("médico é redirecionado ao acessar /audit-logs diretamente", async ({ page }) => {
    await loginAsDoctor(page)

    await page.goto("/audit-logs")
    await expect(page).not.toHaveURL(/\/audit-logs$/)
    await expect(page.getByRole("button", { name: "Pacientes" })).toBeVisible()
  })

  test("paciente é redirecionado para o portal ao acessar /staff", async ({ page }) => {
    await loginAsPatient(page)

    await page.goto("/staff")
    await expect(page).toHaveURL(/\/portal/)
  })
})
