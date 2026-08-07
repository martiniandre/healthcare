import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Medications / Prescriptions Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/patients/fhir-pat-1/encounters/enc-1?tab=medications")
  })

  test("should render the empty medications list initially", async ({ page }) => {
    const title = page.locator("h3", { hasText: "Prescrições (MedicationRequest - FHIR)" })
    await expect(title).toBeVisible()

    const emptyMessage = page.locator("text=Nenhuma medicação prescrita para este atendimento")
    await expect(emptyMessage).toBeVisible()
  })

  test("should create a new medication prescription", async ({ page }) => {
    await page.getByRole("button", { name: "Nova Prescrição" }).click()

    await page.getByPlaceholder("Ex: Amoxicilina 500mg, Losartana Potássica 50mg").fill("Dipirona 500mg")
    await page.getByPlaceholder("Ex: Tomar 1 comprimido via oral de 8 em 8 horas por 7 dias").fill("Tomar 1 comprimido de 6 em 6 horas se houver dor.")

    await page.getByRole("button", { name: "Prescrever" }).click()

    const medicationName = page.locator("text=Dipirona 500mg")
    const dosageInstruction = page.locator("text=Tomar 1 comprimido de 6 em 6 horas se houver dor.")

    await expect(medicationName).toBeVisible()
    await expect(dosageInstruction).toBeVisible()
  })
})
