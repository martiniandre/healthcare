# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: imaging.spec.ts >> Medical Imaging Module (PACS Console) >> should simulate DICOM file upload with progress bar
- Location: e2e\imaging.spec.ts:41:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: locator('text=Iniciando upload e validação de assinatura DICOM...')
Expected: visible
Timeout: 5000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for locator('text=Iniciando upload e validação de assinatura DICOM...')

```

```yaml
- complementary:
  - heading "HealthCare" [level=1]
  - text: Console Clínico v1.0
  - navigation:
    - text: Menu Principal
    - button "Pacientes"
    - button "Telemetria UTI"
    - button "PACS Viewer"
    - button "Análise de Exames"
    - button "Estatísticas"
    - button "Gestão de Equipes"
    - button "Configurações Em breve" [disabled]
  - button "Sair"
  - text: FHIR R4 · gRPC-Web medico
- banner:
  - button
  - text: M medico@clinica.com Profissional
- main:
  - button "Voltar Prontuário"
  - heading "Console Cirúrgico PACS" [level=2]
  - text: "Estudo: Tomografia Computadorizada de Tórax • UID: 1.2.840.10008.5.1.4.1.1.2.20260516.1"
  - button "Upload Novo .DCM"
  - heading "Detalhes do Estudo" [level=3]
  - text: "ID do Estudo study-1 Modalidade Clínica CT Status do Barramento completed Série / Fatias Slice #1 (Visualização Ativa)"
  - img "Visualizador interativo de estudos médicos DICOM do prontuário eletrônico"
  - button "Zoom (Arrastar)"
  - button "Luminosidade"
  - button "Régua (Medir)"
  - button "Tecido Mole"
  - button "Osso"
  - button "Pulmão"
```

# Test source

```ts
  1  | import { test, expect } from "@playwright/test"
  2  | import { loginAsDoctor } from "./helpers"
  3  | 
  4  | test.describe("Medical Imaging Module (PACS Console)", () => {
  5  |   test.beforeEach(async ({ page }) => {
  6  |     await loginAsDoctor(page)
  7  |     await page.goto("/#/imaging/study-1")
  8  |   })
  9  | 
  10 |   test("should load the PACS surgical console with study details", async ({ page }) => {
  11 |     const surgicalConsoleHeading = page.locator("text=Console Cirúrgico PACS")
  12 |     await expect(surgicalConsoleHeading).toBeVisible()
  13 | 
  14 |     const studyTitleText = page.locator("text=Tomografia Computadorizada de Tórax")
  15 |     await expect(studyTitleText).toBeVisible()
  16 | 
  17 |     const imagingCanvas = page.locator("canvas")
  18 |     await expect(imagingCanvas).toBeVisible()
  19 |   })
  20 | 
  21 |   test("should select different tools from the control panel", async ({ page }) => {
  22 |     const zoomToolButton = page.getByRole("button", { name: "Zoom (Arrastar)" })
  23 |     const brightnessToolButton = page.getByRole("button", { name: "Luminosidade" })
  24 |     const rulerToolButton = page.getByRole("button", { name: "Régua (Medir)" })
  25 | 
  26 |     await brightnessToolButton.click()
  27 |     await expect(brightnessToolButton).toHaveClass(/.*bg-primary.*/)
  28 |     await expect(zoomToolButton).not.toHaveClass(/.*bg-primary.*/)
  29 | 
  30 |     await rulerToolButton.click()
  31 |     await expect(rulerToolButton).toHaveClass(/.*bg-primary.*/)
  32 |     await expect(brightnessToolButton).not.toHaveClass(/.*bg-primary.*/)
  33 |   })
  34 | 
  35 |   test("should allow triggering windowing presets", async ({ page }) => {
  36 |     await page.getByRole("button", { name: "Osso" }).click()
  37 |     await page.getByRole("button", { name: "Pulmão" }).click()
  38 |     await page.getByRole("button", { name: "Tecido Mole" }).click()
  39 |   })
  40 | 
  41 |   test("should simulate DICOM file upload with progress bar", async ({ page }) => {
  42 |     let receivedDialogAlertMessage = ""
  43 |     page.on("dialog", async (dialogWindow) => {
  44 |       receivedDialogAlertMessage = dialogWindow.message()
  45 |       await dialogWindow.accept()
  46 |     })
  47 | 
  48 |     await page.getByRole("button", { name: "Upload Novo .DCM" }).click()
  49 | 
  50 |     const progressBarContainer = page.locator("text=Iniciando upload e validação de assinatura DICOM...")
> 51 |     await expect(progressBarContainer).toBeVisible()
     |                                        ^ Error: expect(locator).toBeVisible() failed
  52 | 
  53 |     await page.waitForTimeout(4000)
  54 | 
  55 |     expect(receivedDialogAlertMessage).toContain("DICOM carregado e processado com sucesso")
  56 |   })
  57 | })
  58 | 
```