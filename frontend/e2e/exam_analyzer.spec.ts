import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Exam Analyzer Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/exam-analyzer")
  })

  test("should render the initial history and page title", async ({ page }) => {
    const title = page.locator("h2", { hasText: "Medical Exam Analyzer" })
    await expect(title).toBeVisible()

    const historyItem = page.locator("text=rx_torax.png")
    await expect(historyItem).toBeVisible()
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
    const processingStatus = page.locator("text=Processando Análise Clínica...")
    await expect(processingStatus).toBeVisible()
    
    const newExamHistoryBlock = page.locator("div.group").filter({ hasText: "mock_uploaded_exam.jpg" }).first()
    await expect(newExamHistoryBlock.locator("text=Concluído")).toBeVisible({ timeout: 10000 })
    await expect(processingStatus).toBeHidden()
    const finding = page.locator("text=Nódulo pulmonar calcificado")
    await expect(finding).toBeVisible()
    const conclusion = page.locator("text=Achados benignos, sem necessidade de investigação adicional imediata.")
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
