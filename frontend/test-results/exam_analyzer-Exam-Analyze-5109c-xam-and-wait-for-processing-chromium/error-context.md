# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: exam_analyzer.spec.ts >> Exam Analyzer Module >> should upload a new exam and wait for processing
- Location: e2e\exam_analyzer.spec.ts:18:3

# Error details

```
Test timeout of 30000ms exceeded.
```

```
Error: page.waitForEvent: Test timeout of 30000ms exceeded.
=========================== logs ===========================
waiting for event "filechooser"
============================================================
```

# Page snapshot

```yaml
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
        - generic [ref=e77]:
          - generic [ref=e78]:
            - heading "Medical Exam Analyzer" [level=2] [ref=e79]
            - generic [ref=e80]:
              - img [ref=e81]
              - generic [ref=e84]: IA Multimodal
          - generic [ref=e85]: Análise clínica assistiva de exames médicos por processamento computacional
        - generic [ref=e86]:
          - generic [ref=e87]:
            - generic [ref=e88]:
              - heading "Histórico de Análises" [level=3] [ref=e89]:
                - img [ref=e90]
                - text: Histórico de Análises
              - generic [ref=e94]: 1 exames
            - generic [ref=e95]:
              - img [ref=e96]
              - textbox "Filtrar por nome ou tipo..." [ref=e99]
            - generic [ref=e101] [cursor=pointer]:
              - generic [ref=e102]:
                - generic [ref=e103]: rx_torax.png
                - generic [ref=e104]: Dados Insuficientes
                - generic [ref=e105]:
                  - generic [ref=e106]:
                    - img [ref=e107]
                    - text: 18/05, 07:00
                  - generic [ref=e109]: Concluído
              - button [ref=e110]:
                - img [ref=e111]
          - generic [ref=e114]:
            - generic [ref=e115]:
              - heading "Enviar Exame para Análise" [level=3] [ref=e116]
              - generic [ref=e117]: Arraste arquivos de exames radiológicos, fotos clínicas ou PDFs laboratoriais. A análise é assistiva e probabilística.
              - generic [ref=e118]:
                - generic [ref=e119] [cursor=pointer]:
                  - img [ref=e121]
                  - generic [ref=e124]:
                    - generic [ref=e125]: Selecione ou solte o arquivo aqui
                    - generic [ref=e126]: Imagens (PNG, JPG, DICOM) ou arquivos PDF de até 15MB
                - generic [ref=e127]:
                  - generic [ref=e128] [cursor=pointer]:
                    - img [ref=e130]
                    - generic [ref=e132]:
                      - generic [ref=e133]: Consentimento do Paciente
                      - generic [ref=e134]: Confirmo que possuo a autorização expressa do paciente para submeter seus dados e imagens para processamento clínico assistivo.
                  - generic [ref=e135] [cursor=pointer]:
                    - img [ref=e137]
                    - generic [ref=e139]:
                      - generic [ref=e140]: Anonimização de Segurança (Recomendado)
                      - generic [ref=e141]: Substituir o nome do arquivo enviado por um identificador UUID criptográfico seguro antes da gravação no armazenamento temporário.
                - button "Enviar e Analisar Exame" [disabled] [ref=e142] [cursor=pointer]
            - generic [ref=e143]:
              - img [ref=e145]
              - heading "Nenhum Exame Selecionado" [level=4] [ref=e148]
              - generic [ref=e149]: Selecione um exame no histórico lateral ou faça o upload de um novo exame para visualizar os achados preliminares.
```

# Test source

```ts
  1  | import { test, expect } from "@playwright/test"
  2  | import { loginAsDoctor } from "./helpers"
  3  | 
  4  | test.describe("Exam Analyzer Module", () => {
  5  |   test.beforeEach(async ({ page }) => {
  6  |     await loginAsDoctor(page)
  7  |     await page.goto("/#/exam-analyzer")
  8  |   })
  9  | 
  10 |   test("should render the initial history and page title", async ({ page }) => {
  11 |     const title = page.locator("h2", { hasText: "Medical Exam Analyzer" })
  12 |     await expect(title).toBeVisible()
  13 | 
  14 |     const historyItem = page.locator("text=rx_torax.png")
  15 |     await expect(historyItem).toBeVisible()
  16 |   })
  17 | 
  18 |   test("should upload a new exam and wait for processing", async ({ page }) => {
  19 |     // Fill the file uploader
> 20 |     const fileChooserPromise = page.waitForEvent("filechooser")
     |                                     ^ Error: page.waitForEvent: Test timeout of 30000ms exceeded.
  21 |     await page.locator("label", { hasText: "Selecione um Arquivo" }).click()
  22 |     const fileChooser = await fileChooserPromise
  23 |     
  24 |     // We will upload a dummy file
  25 |     await fileChooser.setFiles({
  26 |       name: "mock_uploaded_exam.jpg",
  27 |       mimeType: "image/jpeg",
  28 |       buffer: Buffer.from("mock_content")
  29 |     })
  30 | 
  31 |     // Check consent and anonymize
  32 |     await page.locator("input[type='checkbox']").first().check()
  33 |     await page.locator("input[type='checkbox']").nth(1).check()
  34 | 
  35 |     await page.getByRole("button", { name: "Enviar para Análise" }).click()
  36 | 
  37 |     // Wait for the processing status to appear
  38 |     const processingStatus = page.locator("text=Análise em andamento...")
  39 |     await expect(processingStatus).toBeVisible()
  40 | 
  41 |     // Wait for the simulated polling to complete (completed state should appear)
  42 |     const completedStatus = page.locator("text=Análise Concluída")
  43 |     await expect(completedStatus).toBeVisible({ timeout: 10000 })
  44 | 
  45 |     // Check if the mock findings are displayed
  46 |     const finding = page.locator("text=Nódulo pulmonar calcificado")
  47 |     await expect(finding).toBeVisible()
  48 |     
  49 |     const conclusion = page.locator("text=Achados benignos, sem necessidade de investigação adicional imediata.")
  50 |     await expect(conclusion).toBeVisible()
  51 |   })
  52 | 
  53 |   test("should delete an analysis from history", async ({ page }) => {
  54 |     const historyItem = page.locator("text=rx_torax.png")
  55 |     await expect(historyItem).toBeVisible()
  56 | 
  57 |     // Click the delete button inside the history item
  58 |     // Assuming the delete button has an aria-label or title. But let's find the specific button inside the item block
  59 |     const historyBlock = page.locator("div").filter({ hasText: "rx_torax.png" }).first()
  60 |     const deleteButton = historyBlock.locator("button.text-gray-400").first()
  61 |     
  62 |     await deleteButton.click()
  63 | 
  64 |     // Since mock deletes it immediately, the item should be removed
  65 |     await expect(historyItem).toBeHidden()
  66 |   })
  67 | })
  68 | 
```