import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Analytics Statistics Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/analytics")
  })

  test("should render the metric strip with values from the stats response", async ({ page }) => {
    await expect(page.getByText("Pacientes Ativos")).toBeVisible()
    await expect(page.getByText("340", { exact: true })).toBeVisible()

    await expect(page.getByText("Conformidade FHIR")).toBeVisible()
    await expect(page.getByText("99.4", { exact: true })).toBeVisible()

    await expect(page.getByText("T. Médio Consulta")).toBeVisible()
    await expect(page.getByText("14.5", { exact: true })).toBeVisible()

    await expect(page.getByText("Atendimentos Semanal")).toBeVisible()
    await expect(page.getByText("79", { exact: true })).toBeVisible()
  })

  test("should display exam modality distribution with counts and percentages", async ({ page }) => {
    await expect(page.getByText("Distribuição de Exames (PACS)")).toBeVisible()
    await expect(page.getByText("CT (Tomografia)", { exact: true })).toBeVisible()
    await expect(page.getByText("MR (Ressonância)", { exact: true })).toBeVisible()
    await expect(page.getByText("CR (Raio-X)", { exact: true })).toBeVisible()
    await expect(page.getByText("US (Ultrassom)", { exact: true })).toBeVisible()
    await expect(page.getByText("45%", { exact: true })).toBeVisible()
    await expect(page.getByText("35", { exact: true })).toBeVisible()
  })

  test("should display weekly consultations volume with localized day labels", async ({ page }) => {
    await expect(page.getByText("Volume de Atendimentos")).toBeVisible()
    await expect(page.getByText("Seg", { exact: true })).toBeVisible()
    await expect(page.getByText("Ter", { exact: true })).toBeVisible()
    await expect(page.getByText("Qua", { exact: true })).toBeVisible()
    await expect(page.getByText("Sex", { exact: true })).toBeVisible()
  })

  test("should summarize weekly consultations with min, average and peak day", async ({ page }) => {
    await expect(page.getByText("Menor", { exact: true })).toBeVisible()
    await expect(page.getByText("Média", { exact: true })).toBeVisible()
    await expect(page.getByText("Pico", { exact: true })).toBeVisible()
    await expect(page.getByText("15 · Sex", { exact: true })).toBeVisible()
  })

  test("should load epidemiology table with pathology classifications", async ({ page }) => {
    await expect(page.getByText("Epidemiologia e Diagnósticos (FHIR Core)")).toBeVisible()
    await expect(page.getByText("Asma não especificada")).toBeVisible()
    await expect(page.getByText("Hipertensão essencial primária")).toBeVisible()
    await expect(page.getByText("Diabetes mellitus tipo 2")).toBeVisible()
    await expect(page.getByText("J45.9", { exact: true })).toBeVisible()
    await expect(page.getByText("I10", { exact: true })).toBeVisible()
    await expect(page.getByText("E11.9", { exact: true })).toBeVisible()
    await expect(page.getByText("248 ativos", { exact: true })).toBeVisible()
  })

  test("should colour trends by direction and magnitude rather than by pathology code", async ({
    page,
  }) => {
    const asthmaRow = page.getByRole("row", { name: /J45\.9/ })
    await expect(asthmaRow.getByText("+5% — Em acompanhamento", { exact: true })).toBeVisible()

    const stableRow = page.getByRole("row", { name: /I10/ })
    await expect(stableRow.getByText("Estável — Dentro do esperado", { exact: true })).toBeVisible()

    const fallingRow = page.getByRole("row", { name: /E11\.9/ })
    await expect(fallingRow.getByText("-3% — Em queda", { exact: true })).toBeVisible()
  })

  test("should render a proportion bar scaled to the largest case count", async ({ page }) => {
    const peakRow = page.getByRole("row", { name: /I10/ })
    await expect(peakRow.getByText("100%", { exact: true })).toBeVisible()
  })
})
