import { test, expect } from "@playwright/test"
import { loginAsPatient } from "./helpers"

test.describe("Report Ready Notification Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsPatient(page)
  })

  test("should land patient on the portal after login", async ({ page }) => {
    await expect(page).toHaveURL(/\/portal/)
    await expect(page.getByRole("heading", { level: 1, name: "Olá, Guilherme!" })).toBeVisible()
  })

  test("should show the report.ready notification in the patient notification bell", async ({ page }) => {
    const bellButton = page.getByRole("button", { name: "Notificações" })
    await expect(bellButton).toBeVisible()

    const unreadBadge = bellButton.locator("text=1")
    await expect(unreadBadge).toBeVisible()

    await bellButton.click()

    const readyNotificationTitle = page.locator("text=Laudo Pronto - Hemograma Completo")
    await expect(readyNotificationTitle).toBeVisible()
  })
})
