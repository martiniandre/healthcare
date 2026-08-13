import { test, expect } from "@playwright/test"
import {
  loginAsDoctor,
  mockAnalyticsAPI,
  mockAnalyzerAPI,
  mockAuthAPI,
  mockClinicalAPI,
  mockPatientsAPI,
  mockScheduleAPI,
  mockStaffAPI,
  mockTelemetryAPI,
} from "./helpers"

test.describe("In-App Notification Bell Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/")
  })

  test("should display unread notification count badge on the bell icon", async ({ page }) => {
    const bellButton = page.getByTitle("Notificações")
    await expect(bellButton).toBeVisible()

    const unreadBadge = bellButton.locator("text=2")
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

test.describe("Notification Bell Error States", () => {
  test("should show an error message when the notification list fails to load", async ({ page }) => {
    await mockAuthAPI(page)
    await mockPatientsAPI(page)
    await mockClinicalAPI(page)
    await mockAnalyzerAPI(page)
    await mockStaffAPI(page)
    await mockScheduleAPI(page)
    await mockTelemetryAPI(page)
    await mockAnalyticsAPI(page)

    await page.route("**/api/v1/notifications*", async (networkRoute) => {
      await networkRoute.fulfill({
        status: 500,
        contentType: "application/json",
        body: JSON.stringify({ error: "Internal Server Error" }),
      })
    })
    await page.route("**/api/v1/notifications/unread-count", async (networkRoute) => {
      await networkRoute.fulfill({
        status: 500,
        contentType: "application/json",
        body: JSON.stringify({ error: "Internal Server Error" }),
      })
    })
    await page.route("**/api/v1/notifications/*/read", async (networkRoute) => {
      await networkRoute.fulfill({
        status: 500,
        contentType: "application/json",
        body: JSON.stringify({ error: "Internal Server Error" }),
      })
    })
    await page.route("**/api/v1/notifications/stream", async (networkRoute) => {
      await networkRoute.fulfill({
        status: 200,
        contentType: "text/event-stream",
        headers: {
          "Cache-Control": "no-cache",
          Connection: "keep-alive",
        },
        body: "",
      })
    })

    await page.goto("/login")
    await page.getByPlaceholder("nome.sobrenome@hospital.com").fill("medico@clinica.com")
    await page.getByPlaceholder("••••••••").fill("senha123")
    await page.getByRole("button", { name: "Entrar no Console" }).click()
    await expect(page).toHaveURL(/\/$/)

    await page.getByTitle("Notificações").click()
    await expect(page.getByText("Não foi possível carregar as notificações.")).toBeVisible()
  })
})
