import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Diagnostic Report Versioning Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/patients/fhir-pat-1/encounters/enc-2?tab=reports")
  })

  test("should open version history modal with snapshot entries", async ({ page }) => {
    const reportTitle = page.locator("text=Eletrocardiograma de Repouso")
    await expect(reportTitle).toBeVisible()

    await page.getByTitle("Histórico de Versões").first().click()

    const versionsModalTitle = page.locator(
      "text=Histórico de Versões do Laudo — Eletrocardiograma de Repouso"
    )
    await expect(versionsModalTitle).toBeVisible()

    const firstVersionLabel = page.locator("text=Versão 1")
    await expect(firstVersionLabel).toBeVisible()

    const secondVersionLabel = page.locator("text=Versão 2")
    await expect(secondVersionLabel).toBeVisible()

    const revisedConclusion = page.locator(
      "text=Conclusão revisada após reavaliação clínica do traçado."
    )
    await expect(revisedConclusion).toBeVisible()
  })

  test("should show empty state when report has no version history", async ({ page }) => {
    await page.route("**/api/v1/reports/rep-1/versions", async (networkRoute) => {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify([]),
      })
    })

    const reportTitle = page.locator("text=Eletrocardiograma de Repouso")
    await expect(reportTitle).toBeVisible()

    await page.getByTitle("Histórico de Versões").first().click()

    const emptyStateMessage = page.locator("text=Nenhuma versão registrada para este laudo.")
    await expect(emptyStateMessage).toBeVisible()
  })
})
