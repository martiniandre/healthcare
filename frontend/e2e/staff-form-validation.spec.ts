import { test, expect } from "@playwright/test"
import { loginAsAdmin, loginAsDoctor } from "./helpers"

test.describe("Staff Registration Validation", () => {
  test("should display validation errors for invalid staff inputs", async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/staff")

    await page.getByRole("button", { name: "Cadastrar Profissional" }).click()

    await page.getByPlaceholder("Ex: Dr. André Silva de Araujo").fill("Dr")
    await page.getByPlaceholder("Ex: nome@hospital.com").fill("emailinvalido")
    await page.getByPlaceholder("Ex: CRM-SP 12345").fill("CRM-INVALID")

    await page.getByRole("button", { name: "Salvar Cadastro" }).click()

    await expect(page.locator("text=O nome deve ter no mínimo 3 caracteres")).toBeVisible()
    await expect(page.locator("text=E-mail inválido")).toBeVisible()
    await expect(page.locator("text=Formato inválido. Ex: CRM-SP 12345")).toBeVisible()
    await expect(page.locator("text=Selecione um departamento")).toBeVisible()
  })

  test("should successfully register a new staff member with valid inputs", async ({ page }) => {
    await loginAsAdmin(page)
    await page.goto("/staff")

    await page.getByRole("button", { name: "Cadastrar Profissional" }).click()

    await page.getByPlaceholder("Ex: Dr. André Silva de Araujo").fill("Dra. Roberta Santos")
    await page.locator("#staff-role").click()
    await page.getByRole("option", { name: "enfermeiro" }).click()
    await page.getByPlaceholder("Ex: nome@hospital.com").fill("roberta@hospital.com")
    await page.getByPlaceholder("Ex: CRM-SP 12345").fill("CRM-SP 54321")
    await page.locator("#staff-department").click()
    await page.getByRole("option", { name: "Dermatologia" }).click()

    await page.getByRole("button", { name: "Salvar Cadastro" }).click()

    await expect(page.getByText("Dra. Roberta Santos")).toBeVisible()
    await expect(page.getByText("Dermatologia")).toBeVisible()
  })

  test("should list the department catalog loaded from the staff endpoint", async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/staff")

    await page.getByRole("button", { name: "Cadastrar Profissional" }).click()
    await page.locator("#staff-department").click()

    await expect(page.getByRole("option", { name: "Clínica Geral" })).toBeVisible()
    await expect(page.getByRole("option", { name: "Cardiologia" })).toBeVisible()
    await expect(page.getByRole("option", { name: "Pediatria" })).toBeVisible()
  })
})
