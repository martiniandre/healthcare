import { test, expect } from "@playwright/test"
import { loginAsDoctor, loginAsPatient } from "./helpers"

test.describe("Sidebar Navigation Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
  })

  test("should display all navigation links in the sidebar", async ({ page }) => {
    await expect(page.getByRole("button", { name: "Pacientes" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Telemetria UTI" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Análise de Exames" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Analytics" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Gestão de Equipes" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Sair" })).toBeVisible()
  })

  test("should organize navigation links into topic and feature groups", async ({ page }) => {
    await expect(page.getByText("Assistência ao Paciente")).toBeVisible()
    await expect(page.getByText("Diagnóstico & Insights")).toBeVisible()
    await expect(page.getByText("Administração")).toBeVisible()

    await expect(page.getByText("Assistência ao Paciente")).toHaveCount(1)
    await expect(page.getByText("Diagnóstico & Insights")).toHaveCount(1)
    await expect(page.getByText("Administração")).toHaveCount(1)
  })

  test("should hide restricted administration features from non-admin staff", async ({ page }) => {
    await expect(page.getByRole("button", { name: "Logs de Auditoria" })).toHaveCount(0)
    await expect(page.getByRole("button", { name: "Configurações" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Configurações" })).toBeDisabled()
  })

  test("should navigate through all sidebar links and update URL correctly", async ({ page }) => {
    await page.getByRole("button", { name: "Telemetria UTI" }).click()
    await expect(page).toHaveURL(/\/telemetry$/)

    await page.getByRole("button", { name: "Analytics" }).click()
    await expect(page).toHaveURL(/\/analytics$/)

    await page.getByRole("button", { name: "Gestão de Equipes" }).click()
    await expect(page).toHaveURL(/\/staff$/)

    await page.getByRole("button", { name: "Análise de Exames" }).click()
    await expect(page).toHaveURL(/\/exam-analyzer$/)

    await page.getByRole("button", { name: "Pacientes" }).click()
    await expect(page).toHaveURL(/\/$/)
  })

  test("should highlight the active route in the sidebar", async ({ page }) => {
    await page.getByRole("button", { name: "Analytics" }).click()

    const analyticsButton = page.getByRole("button", { name: "Analytics" })
    await expect(analyticsButton).toBeVisible()
    await expect(analyticsButton).toHaveAttribute("aria-current", "page")
  })
})

test.describe("Patient Sidebar Navigation Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsPatient(page)
  })

  test("should show only the patient access group", async ({ page }) => {
    await expect(page.getByText("Meu Acesso")).toBeVisible()
    await expect(page.getByRole("button", { name: "Meu Portal" })).toBeVisible()

    await expect(page.getByText("Assistência ao Paciente")).toHaveCount(0)
    await expect(page.getByText("Diagnóstico & Insights")).toHaveCount(0)
    await expect(page.getByText("Administração")).toHaveCount(0)
    await expect(page.getByRole("button", { name: "Pacientes" })).toHaveCount(0)
  })
})