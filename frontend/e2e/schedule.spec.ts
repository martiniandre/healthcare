import { test, expect } from "@playwright/test"
import { loginAsAdmin } from "./helpers"

const selectStaffOnPage = async (page: import("@playwright/test").Page) => {
  await page.locator('[role="combobox"]').first().click()
  await page.locator('[role="option"]', { hasText: "Dr. André Silva de Araujo" }).click()
}

const fillAppointmentForm = async (page: import("@playwright/test").Page) => {
  const modalDialog = page.getByRole("dialog")

  await modalDialog.locator('[role="combobox"]').first().click()
  await page.locator('[role="option"]', { hasText: "Guilherme de Souza Araujo" }).click()

  await modalDialog.locator('[role="combobox"]').nth(1).click()
  await page.locator('[role="option"]', { hasText: "Dr. André Silva de Araujo" }).click()

  const timeInputs = page.locator('input[type="time"]')
  await timeInputs.first().fill("09:00")
  await timeInputs.nth(1).fill("10:00")

  await page.getByPlaceholder("Motivo da consulta...").fill("Consulta de rotina")
}

test.describe("Appointment Scheduling Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page)
  })

  test("should display the schedule screen with staff selector", async ({ page }) => {
    await page.goto("/schedule")
    await expect(page.getByRole("heading", { name: "Agenda", exact: true })).toBeVisible()
    await expect(page.getByRole("button", { name: "Novo Agendamento" })).toBeVisible()
  })

  test("should book an appointment and show it on the day agenda", async ({ page }) => {
    await page.goto("/schedule")
    await selectStaffOnPage(page)
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()

    await expect(page.locator("text=Consulta de rotina")).toBeVisible()
    await expect(page.locator("text=Guilherme de Souza Araujo")).toBeVisible()
    await expect(page.locator("text=Agendado")).toBeVisible()
  })

  test("should cancel a scheduled appointment", async ({ page }) => {
    await page.goto("/schedule")
    await selectStaffOnPage(page)
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()
    await expect(page.locator("text=Consulta de rotina")).toBeVisible()

    await page.getByRole("button", { name: "Cancelar agendamento" }).click()
    await page.getByRole("button", { name: "Confirmar Cancelamento" }).click()

    await expect(page.locator("text=Cancelado")).toBeVisible()
  })

  test("should show conflict message when booking an overlapping slot", async ({ page }) => {
    await page.goto("/schedule")
    await selectStaffOnPage(page)
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()
    await expect(page.locator("text=Consulta de rotina")).toBeVisible()

    await page.getByRole("button", { name: "Novo Agendamento" }).click()
    await fillAppointmentForm(page)
    const timeInputs = page.locator('input[type="time"]')
    await timeInputs.first().fill("09:30")
    await timeInputs.nth(1).fill("10:30")
    await page.getByRole("button", { name: "Agendar" }).click()

    await expect(page.locator("text=Conflito de horário")).toBeVisible()
  })
})
