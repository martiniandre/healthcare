import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Analytics Dashboard Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/analytics")
  })

  test("should load and display all analytic metric cards with correct values", async ({ page }) => {
    await expect(page.locator("text=340")).toBeVisible()
    await expect(page.locator("text=99.4%")).toBeVisible()
    await expect(page.locator("text=14.5 min")).toBeVisible()
    await expect(page.locator("text=79")).toBeVisible()
  })

  test("should display exam modalities distribution chart with modality names", async ({ page }) => {
    await expect(page.locator("text=Distribuição de Exames")).toBeVisible()
    await expect(page.locator("text=CT (Tomografia)")).toBeVisible()
    await expect(page.locator("text=MR (Ressonância)")).toBeVisible()
    await expect(page.locator("text=CR (Raio-X)")).toBeVisible()
    await expect(page.locator("text=US (Ultrassom)")).toBeVisible()
  })

  test("should display weekly consultations volume chart with day labels", async ({ page }) => {
    await expect(page.locator("text=Volume de Atendimentos")).toBeVisible()
    await expect(page.getByText("Seg", { exact: true })).toBeVisible()
    await expect(page.getByText("Ter", { exact: true })).toBeVisible()
    await expect(page.getByText("Qua", { exact: true })).toBeVisible()
    await expect(page.getByText("Sex", { exact: true })).toBeVisible()
  })

  test("should load epidemiology table with pathology classifications", async ({ page }) => {
    await expect(page.locator("text=Epidemiologia e Diagnósticos")).toBeVisible()
    await expect(page.locator("text=Asma não especificada")).toBeVisible()
    await expect(page.locator("text=Hipertensão essencial primária")).toBeVisible()
    await expect(page.locator("text=Diabetes mellitus tipo 2")).toBeVisible()
    await expect(page.locator("text=J45.9")).toBeVisible()
    await expect(page.locator("text=I10")).toBeVisible()
    await expect(page.locator("text=E11.9")).toBeVisible()
  })
})
