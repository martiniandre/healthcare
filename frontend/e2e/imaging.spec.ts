import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Medical Imaging Module (PACS Console)", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/patients/fhir-pat-1?tab=pacs")
    await page.getByRole("button", { name: "Abrir PACS" }).click()
  })

  test("should load the PACS surgical console with study details", async ({ page }) => {
    const surgicalConsoleHeading = page.locator("text=Console Clínico PACS")
    await expect(surgicalConsoleHeading).toBeVisible()

    const studyTitleText = page.locator("text=Tomografia Computadorizada de Tórax")
    await expect(studyTitleText).toBeVisible()

    const imagingCanvas = page.locator("canvas")
    await expect(imagingCanvas).toBeVisible()
  })

  test("should select different tools from the control panel", async ({ page }) => {
    const zoomToolButton = page.getByRole("button", { name: "Zoom", exact: true })
    const contrastToolButton = page.getByRole("button", { name: "Contraste", exact: true })
    const measureToolButton = page.getByRole("button", { name: "Medir", exact: true })

    await contrastToolButton.click()
    await expect(contrastToolButton).toHaveAttribute("aria-pressed", "true")
    await expect(zoomToolButton).toHaveAttribute("aria-pressed", "false")

    await measureToolButton.click()
    await expect(measureToolButton).toHaveAttribute("aria-pressed", "true")
    await expect(contrastToolButton).toHaveAttribute("aria-pressed", "false")
  })

  test("should allow triggering windowing presets", async ({ page }) => {
    await page.getByRole("button", { name: "Osso" }).click()
    await page.getByRole("button", { name: "Pulmão" }).click()
    await page.getByRole("button", { name: "Tecido Mole" }).click()
  })

  test("should simulate DICOM file upload with progress bar", async ({ page }) => {
    const fileChooserPromise = page.waitForEvent("filechooser")
    await page.getByRole("button", { name: "Upload Novo .DCM" }).click()
    const fileChooser = await fileChooserPromise
    await fileChooser.setFiles({
      name: "test_exam.dcm",
      mimeType: "application/dicom",
      buffer: Buffer.from("mock_content")
    })

    const successToast = page.getByText("DICOM carregado e processado com sucesso no barramento do PACS!")
    await expect(successToast).toBeVisible()
  })
})

test.describe("Imaging Workspace Error States", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
  })

  test("should show an error state instead of infinite loading when the imaging study fails to load", async ({ page }) => {
    await page.route("**/api/v1/studies/*", async (networkRoute) => {
      await networkRoute.fulfill({
        status: 500,
        contentType: "application/json",
        body: JSON.stringify({ error: "Internal Server Error" }),
      })
    })

    await page.goto("/patients/fhir-pat-1?tab=pacs")
    await page.getByRole("button", { name: "Abrir PACS" }).click()

    await expect(
      page.getByText("Não foi possível carregar o estudo de imagem. Verifique se ele existe ou tente novamente.")
    ).toBeVisible()
  })
})
