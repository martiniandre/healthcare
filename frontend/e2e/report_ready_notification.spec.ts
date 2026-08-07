import { test, expect } from "@playwright/test"
import { loginAsPatient } from "./helpers"

test.describe("Report Ready Notification Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsPatient(page)
  })

  test("should land patient on the portal after login", async ({ page }) => {
    const patientHeaderName = page.locator("h1", { hasText: "Guilherme de Souza Araujo" })
    await expect(patientHeaderName).toBeVisible()
  })

  test("should show the report.ready notification in the patient notification bell", async ({ page }) => {
    const bellButton = page.getByTitle("Notificações")
    await expect(bellButton).toBeVisible()

    const unreadBadge = bellButton.locator("text=1")
    await expect(unreadBadge).toBeVisible()

    await bellButton.click()

    const readyNotificationTitle = page.locator("text=Laudo Pronto - Hemograma Completo")
    await expect(readyNotificationTitle).toBeVisible()
  })
})
