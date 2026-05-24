import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Exam Analyzer Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/#/exam-analyzer")
  })

  test("should render the initial history and page title", async ({ page }) => {
    const title = page.locator("h2", { hasText: "Medical Exam Analyzer" })
    await expect(title).toBeVisible()

    const historyItem = page.locator("text=rx_torax.png")
    await expect(historyItem).toBeVisible()
  })

  test("should upload a new exam and wait for processing", async ({ page }) => {
    // Fill the file uploader
    const fileChooserPromise = page.waitForEvent("filechooser")
    await page.locator("label", { hasText: "Selecione um Arquivo" }).click()
    const fileChooser = await fileChooserPromise
    
    // We will upload a dummy file
    await fileChooser.setFiles({
      name: "mock_uploaded_exam.jpg",
      mimeType: "image/jpeg",
      buffer: Buffer.from("mock_content")
    })

    // Check consent and anonymize
    await page.locator("input[type='checkbox']").first().check()
    await page.locator("input[type='checkbox']").nth(1).check()

    await page.getByRole("button", { name: "Enviar para Análise" }).click()

    // Wait for the processing status to appear
    const processingStatus = page.locator("text=Análise em andamento...")
    await expect(processingStatus).toBeVisible()

    // Wait for the simulated polling to complete (completed state should appear)
    const completedStatus = page.locator("text=Análise Concluída")
    await expect(completedStatus).toBeVisible({ timeout: 10000 })

    // Check if the mock findings are displayed
    const finding = page.locator("text=Nódulo pulmonar calcificado")
    await expect(finding).toBeVisible()
    
    const conclusion = page.locator("text=Achados benignos, sem necessidade de investigação adicional imediata.")
    await expect(conclusion).toBeVisible()
  })

  test("should delete an analysis from history", async ({ page }) => {
    const historyItem = page.locator("text=rx_torax.png")
    await expect(historyItem).toBeVisible()

    // Click the delete button inside the history item
    // Assuming the delete button has an aria-label or title. But let's find the specific button inside the item block
    const historyBlock = page.locator("div").filter({ hasText: "rx_torax.png" }).first()
    await historyBlock.hover()
    
    const deleteButton = historyBlock.locator("button.text-gray-400").first()
    await deleteButton.click({ force: true })

    // Since mock deletes it immediately, the item should be removed
    await expect(historyItem).toBeHidden()
  })
})
