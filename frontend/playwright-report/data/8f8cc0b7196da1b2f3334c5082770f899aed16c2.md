# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: telemetry_stats_staff.spec.ts >> Clinical Statistics and Analytics Module >> should load diagnostic statistics and key KPIs
- Location: e2e\telemetry_stats_staff.spec.ts:58:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: locator('text=340')
Expected: visible
Timeout: 5000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for locator('text=340')

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
  - heading "Estatísticas Clínicas & Analytics" [level=2]
  - text: Visão epidemiológica agregada de prontuários FHIR e laudos radiológicos PACS Pacientes Ativos 2 +12% este mês Conformidade FHIR 99.4% R4 Core compliant T. Médio Consulta 14.5 min -1.2m vs mês anterior Atendimentos Semanal 7 5 leitos de UTI ativos
  - heading "Distribuição de Exames (PACS)" [level=3]
  - text: Modalidade dos estudos DICOM integrados
  - img
  - text: Total 16 Estudos CT (Tomografia) 7 exames 45% MR (Ressonância) 5 exames 30% CR (Raio-X) 2 exames 15% US (Ultrassom) 2 exames 10%
  - heading "Volume de Atendimentos" [level=3]
  - text: "Evolução diária de atendimentos médicos e triagem Seg Ter Qua Qui Sex Sáb Dom Menor: 1 (Dom) Média: 0 / Dia Pico: 1 (Sex)"
  - heading "Epidemiologia e Diagnósticos (FHIR Core)" [level=3]
  - text: Casos clínicos ativos mapeados na base de dados FHIR
  - button "Exportar CSV"
  - table:
    - rowgroup:
      - row "Código CID Descrição da Patologia Categoria FHIR Casos Ativos Tendência":
        - columnheader "Código CID"
        - columnheader "Descrição da Patologia"
        - columnheader "Categoria FHIR"
        - columnheader "Casos Ativos"
        - columnheader "Tendência"
    - rowgroup:
      - row "J45.9 Asma não especificada Respiratory 1 +5%":
        - cell "J45.9"
        - cell "Asma não especificada"
        - cell "Respiratory"
        - cell "1"
        - cell "+5%"
      - row "I10 Hipertensão essencial primária Cardiovascular 2 Estável":
        - cell "I10"
        - cell "Hipertensão essencial primária"
        - cell "Cardiovascular"
        - cell "2"
        - cell "Estável"
      - row "E11.9 Diabetes mellitus tipo 2 Endocrine 1 +12%":
        - cell "E11.9"
        - cell "Diabetes mellitus tipo 2"
        - cell "Endocrine"
        - cell "1"
        - cell "+12%"
```

# Test source

```ts
  1   | import { test, expect } from "@playwright/test"
  2   | import { loginAsDoctor } from "./helpers"
  3   | 
  4   | test.describe("Real-time Telemetry Dashboard Module", () => {
  5   |   test.beforeEach(async ({ page }) => {
  6   |     await loginAsDoctor(page)
  7   |     await page.goto("/#/telemetry")
  8   |   })
  9   | 
  10  |   test("should display telemetry rooms and default unlocked room status", async ({ page }) => {
  11  |     const greenRoomName = page.getByRole("heading", { name: "Sala Verde - Semi-Intensiva" })
  12  |     await expect(greenRoomName).toBeVisible()
  13  | 
  14  |     const redRoomName = page.getByRole("heading", { name: "Sala Vermelha - Choque & Emergência" })
  15  |     await expect(redRoomName).toBeVisible()
  16  | 
  17  |     const unlockedBadge = page.locator("text=Desbloqueada").first()
  18  |     await expect(unlockedBadge).toBeVisible()
  19  |   })
  20  | 
  21  |   test("should handle room unlocking workflow with passcode verification", async ({ page }) => {
  22  |     const redRoomCard = page.getByRole("heading", { name: "Sala Vermelha - Choque & Emergência" })
  23  |     await redRoomCard.click()
  24  | 
  25  |     const passwordInput = page.getByPlaceholder("Digite o código da sala...")
  26  |     const confirmButton = page.getByRole("button", { name: "Liberar Monitoramento" })
  27  | 
  28  |     await passwordInput.fill("1111")
  29  |     await confirmButton.click()
  30  | 
  31  |     const errorMessage = page.locator("text=Senha de Acesso incorreta. Verifique a escala do plantão.")
  32  |     await expect(errorMessage).toBeVisible()
  33  | 
  34  |     await passwordInput.fill("9999")
  35  |     await confirmButton.click()
  36  | 
  37  |     await expect(errorMessage).not.toBeVisible()
  38  |     const activeBedPatientName = page.getByRole("heading", { name: "Mariana Costa Silva", exact: true })
  39  |     await expect(activeBedPatientName).toBeVisible()
  40  |   })
  41  | 
  42  |   test("should toggle telemetry audible alarms status", async ({ page }) => {
  43  |     const alarmsButton = page.getByRole("button", { name: "Alarmes Silenciados" })
  44  |     await expect(alarmsButton).toBeVisible()
  45  | 
  46  |     await alarmsButton.click()
  47  |     const activeAlarmsButton = page.getByRole("button", { name: "Alarme Sonoro Ativo" })
  48  |     await expect(activeAlarmsButton).toBeVisible()
  49  |   })
  50  | })
  51  | 
  52  | test.describe("Clinical Statistics and Analytics Module", () => {
  53  |   test.beforeEach(async ({ page }) => {
  54  |     await loginAsDoctor(page)
  55  |     await page.goto("/#/stats")
  56  |   })
  57  | 
  58  |   test("should load diagnostic statistics and key KPIs", async ({ page }) => {
  59  |     const activePatientsCount = page.locator("text=340")
> 60  |     await expect(activePatientsCount).toBeVisible()
      |                                       ^ Error: expect(locator).toBeVisible() failed
  61  | 
  62  |     const fhirComplianceRate = page.locator("text=99.4%")
  63  |     await expect(fhirComplianceRate).toBeVisible()
  64  | 
  65  |     const averageDurationTime = page.locator("text=14.5 min")
  66  |     await expect(averageDurationTime).toBeVisible()
  67  |   })
  68  | 
  69  |   test("should load weekly consultations volume chart data", async ({ page }) => {
  70  |     const weeklyChartHeading = page.locator("text=Volume de Atendimentos")
  71  |     await expect(weeklyChartHeading).toBeVisible()
  72  | 
  73  |     const chartDaySeg = page.getByText("Seg", { exact: true })
  74  |     const chartDaySex = page.getByText("Sex", { exact: true })
  75  |     await expect(chartDaySeg).toBeVisible()
  76  |     await expect(chartDaySex).toBeVisible()
  77  |   })
  78  | 
  79  |   test("should load and list clinical pathlogies classification table", async ({ page }) => {
  80  |     const pathTitle = page.locator("text=Asma não especificada")
  81  |     const pathCode = page.locator("text=J45.9")
  82  |     await expect(pathTitle).toBeVisible()
  83  |     await expect(pathCode).toBeVisible()
  84  |   })
  85  | })
  86  | 
  87  | test.describe("Hospital Staff Management Module", () => {
  88  |   test.beforeEach(async ({ page }) => {
  89  |     await loginAsDoctor(page)
  90  |     await page.goto("/#/staff")
  91  |   })
  92  | 
  93  |   test("should render the staff members list page", async ({ page }) => {
  94  |     const doctorName = page.locator("text=Dr. André Silva de Araujo")
  95  |     await expect(doctorName).toBeVisible()
  96  | 
  97  |     const nurseName = page.locator("text=Enf. Roberta Santos Almeida")
  98  |     await expect(nurseName).toBeVisible()
  99  |   })
  100 | 
  101 |   test("should filter staff members list by role tab selection", async ({ page }) => {
  102 |     const searchField = page.getByPlaceholder("Buscar por nome, e-mail ou especialidade...")
  103 |     await searchField.fill("Roberta")
  104 | 
  105 |     const matchingNurse = page.locator("text=Enf. Roberta Santos Almeida")
  106 |     const nonMatchingDoctor = page.locator("text=Dr. André Silva de Araujo")
  107 | 
  108 |     await expect(matchingNurse).toBeVisible()
  109 |     await expect(nonMatchingDoctor).not.toBeVisible()
  110 |   })
  111 | 
  112 |   test("should successfully register a new hospital practitioner", async ({ page }) => {
  113 |     await page.getByRole("button", { name: "Cadastrar Profissional" }).click()
  114 | 
  115 |     await page.getByPlaceholder("Ex: Dr. André Silva de Araujo").fill("Dra. Paula Albuquerque")
  116 |     await page.locator("select").selectOption("Médico")
  117 |     await page.getByPlaceholder("Ex: CRM-SP 12345").fill("CRM-SP 777777")
  118 |     await page.getByPlaceholder("Ex: nome@hospital.com").fill("paula.albuquerque@hospital.com")
  119 |     await page.getByPlaceholder("Ex: Cardiologia").fill("Neurologia")
  120 | 
  121 |     await page.getByRole("button", { name: "Salvar Cadastro" }).click()
  122 | 
  123 |     const newPractitionerName = page.locator("text=Dra. Paula Albuquerque")
  124 |     await expect(newPractitionerName).toBeVisible()
  125 |   })
  126 | })
  127 | 
```