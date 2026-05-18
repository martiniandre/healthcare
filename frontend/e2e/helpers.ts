import { type Page, expect } from "@playwright/test"

export const mockAuthAPI = async (pageInstance: Page): Promise<void> => {
  await pageInstance.route("**/api/auth/login", async (networkRoute) => {
    const httpRequest = networkRoute.request()
    const submittedJSON = httpRequest.postDataJSON()

    if (submittedJSON.email === "medico@clinica.com" && submittedJSON.password === "senha123") {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          token: "mock-jwt-token-123456",
          userId: "user-medico-123",
          role: "doctor",
          email: "medico@clinica.com",
        }),
      })
    } else {
      await networkRoute.fulfill({
        status: 401,
        contentType: "application/json",
        body: JSON.stringify({
          error: "Credenciais inválidas.",
        }),
      })
    }
  })

  await pageInstance.route("**/api/auth/logout", async (networkRoute) => {
    await networkRoute.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        success: true,
      }),
    })
  })
}

export const loginAsDoctor = async (pageInstance: Page): Promise<void> => {
  await mockAuthAPI(pageInstance)
  await pageInstance.goto("/#/login")
  await pageInstance.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
  await pageInstance.getByPlaceholder("••••••••").fill("senha123")
  await pageInstance.getByRole("button", { name: "Entrar no Console" }).click()
  await expect(pageInstance).toHaveURL(/.*#\/$/)
}
