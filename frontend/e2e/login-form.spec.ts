import { test, expect } from "@playwright/test"
import { mockAuthAPI } from "./helpers"

test.describe("Login Form Interactions", () => {
  test.beforeEach(async ({ page }) => {
    await mockAuthAPI(page)
  })

  test("should toggle password visibility with the eye button", async ({ page }) => {
    await page.goto("/login")
    const passwordInput = page.getByPlaceholder("••••••••")
    await passwordInput.fill("senha123")

    await expect(passwordInput).toHaveAttribute("type", "password")

    await page.getByRole("button", { name: "Mostrar senha" }).click()
    await expect(passwordInput).toHaveAttribute("type", "text")
    await expect(passwordInput).toHaveValue("senha123")

    await page.getByRole("button", { name: "Ocultar senha" }).click()
    await expect(passwordInput).toHaveAttribute("type", "password")
  })

  test("should show i18n validation errors on empty submit instead of native browser validation", async ({ page }) => {
    await page.goto("/login")
    await page.getByRole("button", { name: "Entrar no Console" }).click()

    await expect(page.getByText("O e-mail é obrigatório")).toBeVisible()
    await expect(page.getByText("A senha deve ter no mínimo 8 caracteres")).toBeVisible()
  })
})

test.describe("Login Form Loading State", () => {
  test("should disable the submit button and show the loading state while authenticating", async ({ page }) => {
    await mockAuthAPI(page)

    let loginCompleted = false
    await page.route("**/api/v1/auth/login", async (networkRoute) => {
      await new Promise((resolve) => setTimeout(resolve, 1500))
      loginCompleted = true
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          token: "mock-jwt-token-doctor-123456",
          userId: "user-medico-123",
          role: "DOCTOR",
          email: "medico@clinica.com",
        }),
      })
    })

    await page.route("**/api/v1/auth/me", async (networkRoute) => {
      if (loginCompleted) {
        await networkRoute.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            token: "mock-jwt-token-doctor-123456",
            userId: "user-medico-123",
            role: "DOCTOR",
            email: "medico@clinica.com",
            fullName: "Dr. André Silva de Araujo",
            isActive: true,
          }),
        })
        return
      }
      await networkRoute.fulfill({
        status: 401,
        contentType: "application/json",
        body: JSON.stringify({ error: "Não autenticado." }),
      })
    })

    await page.goto("/login")
    await page.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
    await page.getByPlaceholder("••••••••").fill("senha123")
    await page.getByRole("button", { name: "Entrar no Console" }).click()

    const loadingSubmitButton = page.getByRole("button", { name: "Autenticando via gRPC..." })
    await expect(loadingSubmitButton).toBeVisible()
    await expect(loadingSubmitButton).toBeDisabled()
    await expect(loadingSubmitButton).toHaveAttribute("aria-busy", "true")

    await expect(page).toHaveURL(/\/$/)
  })
})
