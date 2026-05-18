import { test, expect } from "@playwright/test"

test.describe("Portal de Acesso Medico - Autenticacao Clinica", () => {
  test("deve exibir erro para credenciais invalidas", async ({ page }) => {
    await page.goto("/login")
    await page.getByPlaceholder("nome.sobrenome@hospital.com").fill("usuario.errado@clinica.com")
    await page.getByPlaceholder("••••••••").fill("senha_errada")
    await page.getByRole("button", { name: "Entrar no Console" }).click()
    const errorAlertElement = page.locator("text=Credenciais inválidas.")
    await expect(errorAlertElement).toBeVisible()
  })

  test("deve realizar login clinico com sucesso como medico", async ({ page }) => {
    await page.goto("/login")
    await page.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
    await page.getByPlaceholder("••••••••").fill("senha123")
    await page.getByRole("button", { name: "Entrar no Console" }).click()
    await expect(page).toHaveURL("/")
    const headerTitleElement = page.locator("text=Fila de Atendimento Clínico")
    await expect(headerTitleElement).toBeVisible()
  })

  test("deve realizar logout clinico com sucesso", async ({ page }) => {
    await page.goto("/login")
    await page.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
    await page.getByPlaceholder("••••••••").fill("senha123")
    await page.getByRole("button", { name: "Entrar no Console" }).click()
    await expect(page).toHaveURL("/")
    await page.getByRole("button", { name: "Sair" }).click()
    await expect(page).toHaveURL("/login")
  })
})
