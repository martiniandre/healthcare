import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Medical Imaging Module (PACS Console)", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/#/imaging/study-1")
  })

  test("should load the PACS surgical console with study details", async ({ page }) => {
    const surgicalConsoleHeading = page.locator("text=Console Cirúrgico PACS")
    await expect(surgicalConsoleHeading).toBeVisible()

    const studyTitleText = page.locator("text=Tomografia Computadorizada de Tórax")
    await expect(studyTitleText).toBeVisible()

    const imagingCanvas = page.locator("canvas")
    await expect(imagingCanvas).toBeVisible()
  })

  test("should select different tools from the control panel", async ({ page }) => {
    const zoomToolButton = page.getByRole("button", { name: "Zoom (Arrastar)" })
    const brightnessToolButton = page.getByRole("button", { name: "Luminosidade" })
    const rulerToolButton = page.getByRole("button", { name: "Régua (Medir)" })

    await brightnessToolButton.click()
    await expect(brightnessToolButton).toHaveClass(/.*bg-primary.*/)
    await expect(zoomToolButton).not.toHaveClass(/.*bg-primary.*/)

    await rulerToolButton.click()
    await expect(rulerToolButton).toHaveClass(/.*bg-primary.*/)
    await expect(brightnessToolButton).not.toHaveClass(/.*bg-primary.*/)
  })

  test("should allow triggering windowing presets", async ({ page }) => {
    await page.getByRole("button", { name: "Osso" }).click()
    await page.getByRole("button", { name: "Pulmão" }).click()
    await page.getByRole("button", { name: "Tecido Mole" }).click()
  })

  test("should simulate DICOM file upload with progress bar", async ({ page }) => {
    let receivedDialogAlertMessage = ""
    page.on("dialog", async (dialogWindow) => {
      receivedDialogAlertMessage = dialogWindow.message()
      await dialogWindow.accept()
    })

    await page.getByRole("button", { name: "Upload Novo .DCM" }).click()

    const progressBarContainer = page.locator("text=Iniciando upload de chunks gRPC-Web...")
    await expect(progressBarContainer).toBeVisible()

    await page.waitForTimeout(4000)

    expect(receivedDialogAlertMessage).toContain("DICOM carregado e processado com sucesso")
  })
})
