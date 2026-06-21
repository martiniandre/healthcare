import { test, expect } from "@playwright/test"
import { mockAuthAPI } from "./helpers"

test.describe("Medical Access Portal - Clinical Authentication", () => {
  test.beforeEach(async ({ page }) => {
    await mockAuthAPI(page)
  })

  test("should display error for invalid credentials", async ({ page }) => {
    await page.goto("/login")
    await page.getByPlaceholder("nome.sobrenome@hospital.com").fill("usuario.errado@clinica.com")
    await page.getByPlaceholder("••••••••").fill("senha_errada")
    await page.getByRole("button", { name: "Entrar no Console" }).click()

    const errorAlertElement = page.locator("text=Credenciais inválidas.")
    await expect(errorAlertElement).toBeVisible()
  })

  test("should successfully perform clinical login as doctor", async ({ page }) => {
    await page.goto("/login")
    await page.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
    await page.getByPlaceholder("••••••••").fill("senha123")
    await page.getByRole("button", { name: "Entrar no Console" }).click()

    await expect(page).toHaveURL(/\/$/)
    const headerTitleElement = page.locator("text=Gestão de prontuários e dados clínicos FHIR")
    await expect(headerTitleElement).toBeVisible()
  })

  test("should successfully perform clinical logout", async ({ page }) => {
    await page.goto("/login")
    await page.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
    await page.getByPlaceholder("••••••••").fill("senha123")
    await page.getByRole("button", { name: "Entrar no Console" }).click()

    await expect(page).toHaveURL(/\/$/)
    await page.getByRole("button", { name: "Sair" }).click()
    await expect(page).toHaveURL(/\/login$/)
  })
})
