import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Staff Registration Validation", () => {
  test.beforeEach(async ({ page }) => {
    const currentEmployees = [
      {
        id: "emp-1",
        user_id: "user-1",
        full_name: "Dr. André Silva de Araujo",
        email: "andre@hospital.com",
        role: "Médico",
        crm_number: "CRM-SP 12345",
        is_active: true,
        status: "Plantonista",
        department: "Cardiologia",
      },
    ]

    await page.route("**/staff/employees", async (route) => {
      const request = route.request()
      if (request.method() === "GET") {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(currentEmployees),
        })
      } else if (request.method() === "POST") {
        const data = request.postDataJSON()
        const newEmployee = {
          id: `emp-${currentEmployees.length + 1}`,
          user_id: data.user_id,
          full_name: data.full_name,
          email: data.email,
          role: data.role,
          crm_number: data.crm_number,
          is_active: true,
          status: "Plantonista",
          department: "Clínica Geral",
        }
        currentEmployees.push(newEmployee)
        await route.fulfill({
          status: 201,
          contentType: "application/json",
          body: JSON.stringify({ employee_id: newEmployee.id }),
        })
      }
    })

    await loginAsDoctor(page)
    await page.goto("/staff")
  })

  test("should display validation errors for invalid staff inputs", async ({ page }) => {
    await page.getByRole("button", { name: "Cadastrar Profissional" }).click()

    await page.getByPlaceholder("Ex: Dr. André Silva de Araujo").fill("Dr")
    await page.getByPlaceholder("Ex: nome@hospital.com").fill("emailinvalido")
    await page.getByPlaceholder("Ex: CRM-SP 12345").fill("CRM-INVALID")

    await page.getByRole("button", { name: "Salvar Cadastro" }).click()

    const nameError = page.locator("text=O nome deve ter no mínimo 3 caracteres")
    const emailError = page.locator("text=E-mail inválido")
    const licenseError = page.locator("text=Formato inválido. Ex: CRM-SP 12345")

    await expect(nameError).toBeVisible()
    await expect(emailError).toBeVisible()
    await expect(licenseError).toBeVisible()
  })

  test("should successfully register a new staff member with valid inputs", async ({ page }) => {
    await page.getByRole("button", { name: "Cadastrar Profissional" }).click()

    await page.getByPlaceholder("Ex: Dr. André Silva de Araujo").fill("Dra. Roberta Santos")
    await page.locator("select").selectOption("Médico")
    await page.getByPlaceholder("Ex: nome@hospital.com").fill("roberta@hospital.com")
    await page.getByPlaceholder("Ex: CRM-SP 12345").fill("CRM-SP 54321")

    await page.getByRole("button", { name: "Salvar Cadastro" }).click()

    const newName = page.locator("text=Dra. Roberta Santos")
    await expect(newName).toBeVisible()
  })
})
