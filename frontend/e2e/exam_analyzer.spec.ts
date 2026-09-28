import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Exam Analyzer Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/exam-analyzer")
  })

  test("should render the localized page title and the analysis history", async ({ page }) => {
    const title = page.locator("h2", { hasText: "Análise de Exames" })
    await expect(title).toBeVisible()

    const historyItem = page.locator("text=rx_torax.png")
    await expect(historyItem).toBeVisible()
  })

  test("should explain why the submit button is disabled before a file is chosen", async ({ page }) => {
    const submitButton = page.getByRole("button", { name: "Enviar para Análise" })
    await expect(submitButton).toBeDisabled()
    await expect(page.getByText("Selecione um arquivo para continuar.")).toBeVisible()
  })

  test("should demand patient consent before the submit button unlocks", async ({ page }) => {
    const fileChooserPromise = page.waitForEvent("filechooser")
    await page.locator("label", { hasText: "Selecione um Arquivo" }).click()
    const fileChooser = await fileChooserPromise
    await fileChooser.setFiles({
      name: "consent_check.jpg",
      mimeType: "image/jpeg",
      buffer: Buffer.from("mock_content")
    })

    const submitButton = page.getByRole("button", { name: "Enviar para Análise" })
    await expect(submitButton).toBeDisabled()
    await expect(
      page.getByText("Confirme o consentimento do paciente para habilitar o envio.")
    ).toBeVisible()

    await page.locator("input[type='checkbox']").first().check({ force: true })
    await expect(submitButton).toBeEnabled()
  })

  test("should surface the suggested course of action and the image quality verdict", async ({
    page,
  }) => {
    await page.getByRole("button", { name: "Ver análise de rx_torax.png" }).click()

    const suggestedConduct = page.getByText("Conduta Sugerida", { exact: true })
    await expect(suggestedConduct).toBeVisible()
    await expect(page.getByText("Urgência Médica", { exact: true })).toBeVisible()
    await expect(page.getByText("Agendar consulta com pneumologista ou clínico geral.")).toBeVisible()

    await expect(page.getByText("Qualidade da Imagem", { exact: true })).toBeVisible()
    await expect(page.getByText("Aprovada para leitura", { exact: true })).toBeVisible()
  })

  test("should rank each finding by combining significance and confidence", async ({ page }) => {
    await page.getByRole("button", { name: "Ver análise de rx_torax.png" }).click()

    const highSeverityFinding = page
      .getByText("Área de consolidação pulmonar no lobo inferior direito")
      .locator("xpath=ancestor::li")

    await expect(highSeverityFinding).toContainText("Relevância: Alta")
    await expect(highSeverityFinding).toContainText("Prioridade")
    await expect(highSeverityFinding).toContainText("88%")

    const lowSeverityFinding = page
      .getByText("Ausência de derrame pleural")
      .locator("xpath=ancestor::li")

    await expect(lowSeverityFinding).toContainText("Relevância: Baixa")
    await expect(lowSeverityFinding).toContainText("Informativo")
  })

  test("should keep the history selectable and deletable with the keyboard", async ({ page }) => {
    const selectButton = page.getByRole("button", { name: "Ver análise de rx_torax.png" })
    await expect(selectButton).toBeVisible()

    const deleteButton = page.getByRole("button", { name: "Excluir análise de rx_torax.png" })
    await expect(deleteButton).toBeVisible()
  })

  test("should upload a new exam and wait for processing", async ({ page }) => {
    const fileChooserPromise = page.waitForEvent("filechooser")
    await page.locator("label", { hasText: "Selecione um Arquivo" }).click()
    const fileChooser = await fileChooserPromise
    await fileChooser.setFiles({
      name: "mock_uploaded_exam.jpg",
      mimeType: "image/jpeg",
      buffer: Buffer.from("mock_content")
    })
    await page.locator("input[type='checkbox']").first().check({ force: true })
    await page.locator("input[type='checkbox']").nth(1).check({ force: true })
    await page.getByRole("button", { name: "Enviar para Análise" }).click()
    const processingStatus = page.getByText("Processando Análise Clínica...")
    await expect(processingStatus).toBeVisible()

    const newExamHistoryBlock = page.locator("div.group").filter({ hasText: "mock_uploaded_exam.jpg" }).first()
    await expect(newExamHistoryBlock.getByText("Concluído")).toBeVisible({ timeout: 10000 })
    await expect(processingStatus).toBeHidden()
    const finding = page.getByText("Nódulo pulmonar calcificado")
    await expect(finding).toBeVisible()
    const conclusion = page.getByText("Achados benignos, sem necessidade de investigação adicional imediata.")
    await expect(conclusion).toBeVisible()
  })

  test("should delete an analysis from history", async ({ page }) => {
    const historyItem = page.locator("text=rx_torax.png")
    await expect(historyItem).toBeVisible()
    const historyBlock = page.locator("div.group").filter({ hasText: "rx_torax.png" }).first()
    await historyBlock.hover()
    const deleteButton = historyBlock.locator("button.text-gray-400").first()
    await deleteButton.click({ force: true })
    await expect(historyItem).toBeHidden()
  })
})
