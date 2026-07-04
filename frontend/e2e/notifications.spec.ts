import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("In-App Notification Bell Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/")
  })

  test("should display unread notification count badge on the bell icon", async ({ page }) => {
    const bellButton = page.getByTitle("Notificações")
    await expect(bellButton).toBeVisible()

    const unreadBadge = page.locator("text=2")
    await expect(unreadBadge).toBeVisible()
  })

  test("should open dropdown with notification list when bell is clicked", async ({ page }) => {
    await page.getByTitle("Notificações").click()

    const dropdownTitle = page.locator("text=Notificações")
    await expect(dropdownTitle).toBeVisible()

    const firstNotification = page.locator("text=Alerta Crítico - Leito 01")
    await expect(firstNotification).toBeVisible()
  })

  test("should mark notification as read when clicked", async ({ page }) => {
    await page.getByTitle("Notificações").click()

    const criticalAlert = page.locator("text=Alerta Crítico - Leito 01")
    await criticalAlert.click()

    await expect(criticalAlert).not.toBeVisible()
  })
})
