# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: imaging.spec.ts >> Medical Imaging Module (PACS Console) >> should simulate DICOM file upload with progress bar
- Location: e2e\imaging.spec.ts:41:3

# Error details

```
Error: expect(received).toContain(expected) // indexOf

Expected substring: "DICOM carregado e processado com sucesso"
Received string:    ""
```

# Page snapshot

```yaml
- generic [ref=e2]:
  - generic [ref=e3]:
    - complementary [ref=e4]:
      - generic [ref=e6]:
        - img [ref=e8]
        - generic [ref=e10]:
          - heading "HealthCare" [level=1] [ref=e11]
          - text: Console Clínico v1.0
      - navigation [ref=e13]:
        - generic [ref=e14]: Menu Principal
        - button "Pacientes" [ref=e15]:
          - img [ref=e16]
          - text: Pacientes
        - button "Telemetria UTI" [ref=e21]:
          - img [ref=e22]
          - text: Telemetria UTI
        - button "PACS Viewer" [ref=e24]:
          - img [ref=e25]
          - text: PACS Viewer
        - button "Análise de Exames" [ref=e29]:
          - img [ref=e30]
          - text: Análise de Exames
        - button "Estatísticas" [ref=e33]:
          - img [ref=e34]
          - text: Estatísticas
        - button "Gestão de Equipes" [ref=e36]:
          - img [ref=e37]
          - text: Gestão de Equipes
        - button "Configurações Em breve" [disabled] [ref=e42]:
          - img [ref=e43]
          - text: Configurações
          - generic [ref=e46]: Em breve
      - button "Sair" [ref=e48]:
        - img [ref=e49]
        - text: Sair
      - generic [ref=e52]:
        - generic [ref=e55]: FHIR R4 · gRPC-Web
        - generic [ref=e56]: medico
    - generic [ref=e57]:
      - banner [ref=e58]:
        - button [ref=e59]:
          - img [ref=e60]
        - generic [ref=e65]:
          - generic [ref=e66]: M
          - generic [ref=e67]:
            - generic [ref=e68]: medico@clinica.com
            - generic [ref=e69]:
              - img [ref=e70]
              - generic [ref=e73]: Profissional
      - main [ref=e74]:
        - generic [ref=e75]:
          - generic [ref=e76]:
            - generic [ref=e77]:
              - button "Voltar Prontuário" [ref=e78] [cursor=pointer]:
                - img [ref=e79]
                - text: Voltar Prontuário
              - generic [ref=e81]:
                - heading "Console Cirúrgico PACS" [level=2] [ref=e82]
                - generic [ref=e83]: "Estudo: Tomografia Computadorizada de Tórax • UID: 1.2.840.10008.5.1.4.1.1.2.20260516.1"
            - button "Upload Novo .DCM" [ref=e85] [cursor=pointer]:
              - img [ref=e86]
              - text: Upload Novo .DCM
          - generic [ref=e89]:
            - generic [ref=e90]:
              - heading "Detalhes do Estudo" [level=3] [ref=e91]:
                - img [ref=e92]
                - text: Detalhes do Estudo
              - generic [ref=e95]:
                - generic [ref=e96]:
                  - generic [ref=e97]: ID do Estudo
                  - generic [ref=e98]: study-1
                - generic [ref=e99]:
                  - generic [ref=e100]: Modalidade Clínica
                  - generic [ref=e101]: CT
                - generic [ref=e102]:
                  - generic [ref=e103]: Status do Barramento
                  - generic [ref=e104]:
                    - img [ref=e105]
                    - text: completed
                - generic [ref=e108]:
                  - generic [ref=e109]: Série / Fatias
                  - generic [ref=e110]: "Slice #1 (Visualização Ativa)"
            - generic [ref=e111]:
              - img "Visualizador interativo de estudos médicos DICOM do prontuário eletrônico" [ref=e113]
              - generic [ref=e114]:
                - generic [ref=e115]:
                  - button "Zoom (Arrastar)" [ref=e116] [cursor=pointer]:
                    - img [ref=e117]
                    - text: Zoom (Arrastar)
                  - button "Luminosidade" [ref=e120] [cursor=pointer]:
                    - img [ref=e121]
                    - text: Luminosidade
                  - button "Régua (Medir)" [ref=e127] [cursor=pointer]:
                    - img [ref=e128]
                    - text: Régua (Medir)
                - generic [ref=e134]:
                  - button "Tecido Mole" [ref=e135] [cursor=pointer]
                  - button "Osso" [ref=e136] [cursor=pointer]
                  - button "Pulmão" [ref=e137] [cursor=pointer]
  - generic [ref=e138]:
    - img [ref=e140]
    - generic [ref=e143]:
      - generic [ref=e144]: Sucesso
      - paragraph [ref=e145]: DICOM carregado e processado com sucesso no barramento do PACS!
    - button [ref=e146] [cursor=pointer]:
      - img [ref=e147]
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
  47 |     const fileChooserPromise = page.waitForEvent("filechooser")
  48 |     await page.getByRole("button", { name: "Upload Novo .DCM" }).click()
  49 |     const fileChooser = await fileChooserPromise
  50 |     await fileChooser.setFiles({
  51 |       name: "test_exam.dcm",
  52 |       mimeType: "application/dicom",
  53 |       buffer: Buffer.from("mock_content")
  54 |     })
  55 |     const progressBarContainer = page.locator("text=Iniciando upload e validação de assinatura DICOM...")
  56 |     await expect(progressBarContainer).toBeVisible()
  57 |     await page.waitForTimeout(4000)
> 58 |     expect(receivedDialogAlertMessage).toContain("DICOM carregado e processado com sucesso")
     |                                        ^ Error: expect(received).toContain(expected) // indexOf
  59 |   })
  60 | })
  61 | 
```