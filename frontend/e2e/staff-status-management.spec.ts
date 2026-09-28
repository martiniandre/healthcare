import { test, expect } from "@playwright/test"
import { loginAsAdmin, loginAsDoctor } from "./helpers"

const firstEmployeeRow = (page: import("@playwright/test").Page, employeeName: string) =>
  page.getByRole("row").filter({ hasText: employeeName })

test.describe("Staff Status Management", () => {
  test("should render lowercase role labels instead of uppercase enum values", async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/staff")

    await expect(firstEmployeeRow(page, "Dr. André Silva de Araujo").getByText("médico")).toBeVisible()
    await expect(firstEmployeeRow(page, "Enf. Roberta Santos Almeida").getByText("enfermeiro")).toBeVisible()
    await expect(page.getByText("DOCTOR")).toHaveCount(0)
    await expect(page.getByText("NURSE")).toHaveCount(0)
  })

  test("should show the department of each staff member", async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/staff")

    await expect(firstEmployeeRow(page, "Dr. André Silva de Araujo").getByText("Cardiologia")).toBeVisible()
    await expect(firstEmployeeRow(page, "Enf. Roberta Santos Almeida").getByText("Pediatria")).toBeVisible()
  })

  test("should deactivate a staff member and reactivate it", async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/staff")

    const doctorRow = firstEmployeeRow(page, "Dr. André Silva de Araujo")
    await expect(doctorRow.getByText("Ativo")).toBeVisible()

    await doctorRow.getByRole("button", { name: "Desativar acesso" }).click()
    await expect(page.getByText("Acesso do profissional desativado.")).toBeVisible()
    await expect(firstEmployeeRow(page, "Dr. André Silva de Araujo").getByText("Inativo")).toBeVisible()

    const inactiveRow = firstEmployeeRow(page, "Dr. André Silva de Araujo")
    await inactiveRow.getByRole("button", { name: "Ativar acesso" }).click()
    await expect(page.getByText("Acesso do profissional ativado.")).toBeVisible()
    await expect(firstEmployeeRow(page, "Dr. André Silva de Araujo").getByText("Ativo")).toBeVisible()
  })

  test("should not expose the delete action to doctors", async ({ page }) => {
    await loginAsDoctor(page)
    await page.goto("/staff")

    await expect(firstEmployeeRow(page, "Dr. André Silva de Araujo").getByRole("button", { name: "Desativar acesso" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Excluir profissional" })).toHaveCount(0)
  })

  test("should soft delete a staff member after confirmation", async ({ page }) => {
    await loginAsAdmin(page)
    await page.goto("/staff")

    const nurseRow = firstEmployeeRow(page, "Enf. Roberta Santos Almeida")
    await nurseRow.getByRole("button", { name: "Excluir profissional" }).click()

    const confirmationDialog = page.getByRole("dialog")
    await expect(confirmationDialog.getByText("Excluir profissional")).toBeVisible()
    await confirmationDialog.getByRole("button", { name: "Excluir", exact: true }).click()

    await expect(page.getByText("Profissional excluído do quadro clínico.")).toBeVisible()
    await expect(firstEmployeeRow(page, "Enf. Roberta Santos Almeida").getByText("Inativo")).toBeVisible()
  })

  test("should keep the deleted staff member visible in the list for auditing", async ({ page }) => {
    await loginAsAdmin(page)
    await page.goto("/staff")

    await expect(firstEmployeeRow(page, "Enf. Roberta Santos Almeida")).toBeVisible()
    await expect(page.getByRole("row")).toHaveCount(3)
  })
})
