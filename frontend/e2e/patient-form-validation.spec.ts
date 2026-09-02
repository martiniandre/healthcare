import { test, expect } from "@playwright/test"
import { loginAsAdmin } from "./helpers"

test.describe("Patient Registration Validation", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page)
  })

  test("should display validation errors for invalid inputs", async ({ page }) => {
    await page.getByRole("button", { name: "Novo Paciente" }).click()

    await page.getByPlaceholder("Nome Completo do Paciente").fill("Jo")
    await page.getByPlaceholder("DD/MM/AAAA").fill("01/01/2050")
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

  test("should show a validation error when a field exceeds the maximum length", async ({ page }) => {
    await page.getByRole("button", { name: "Novo Paciente" }).click()

    await page.getByPlaceholder("Nome Completo do Paciente").fill("A".repeat(256))
    await page.getByPlaceholder("DD/MM/AAAA").fill("01/01/1990")
    await page.getByPlaceholder("123.456.789-00").fill("529.982.247-25")
    await page.getByPlaceholder("(11) 98765-4321").fill("(11) 98765-4321")

    await page.getByRole("button", { name: "Confirmar Cadastro" }).click()

    await expect(page.locator("text=Valor excede o limite máximo de caracteres")).toBeVisible()
  })
})
