import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Patient Registration Validation", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
  })

  test("should display validation errors for invalid inputs", async ({ page }) => {
    await page.getByRole("button", { name: "Novo Paciente" }).click()

    await page.getByPlaceholder("Nome Completo do Paciente").fill("Jo")
    await page.getByPlaceholder("AAAA-MM-DD").fill("2050-01-01")
    await page.getByPlaceholder("123.456.789-00").fill("111.111.111-11")
    await page.getByPlaceholder("(11) 98765-4321").fill("99999")

    await page.getByRole("button", { name: "Confirmar Cadastro" }).click()

    const nameError = page.locator("text=O nome deve ter no mínimo 3 caracteres")
    const birthDateError = page.locator("text=A data de nascimento deve ser no passado")
    const documentError = page.locator("text=CPF inválido")
    const phoneError = page.locator("text=Formato de telefone inválido. Ex: (11) 98765-4321")

    await expect(nameError).toBeVisible()
    await expect(birthDateError).toBeVisible()
    await expect(documentError).toBeVisible()
    await expect(phoneError).toBeVisible()
  })
})
