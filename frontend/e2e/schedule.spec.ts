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
  await timeInputs.nth(1).fill("09:30")
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

    await expect(page.locator("text=Sem motivo informado")).toBeVisible()
    await expect(page.locator("text=Guilherme de Souza Araujo")).toBeVisible()
    await expect(page.locator("text=Agendado")).toBeVisible()
  })

  test("should cancel a scheduled appointment", async ({ page }) => {
    await page.goto("/schedule")
    await selectStaffOnPage(page)
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()
    await expect(page.locator("text=Sem motivo informado")).toBeVisible()

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
    await expect(page.locator("text=Sem motivo informado")).toBeVisible()

    await page.getByRole("button", { name: "Novo Agendamento" }).click()
    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()

    await expect(page.locator("text=Conflito de horário")).toBeVisible()
  })

  test("should show a validation error when the appointment date is in the past", async ({ page }) => {
    await page.goto("/schedule")
    await selectStaffOnPage(page)
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    await fillAppointmentForm(page)
    await page.getByRole("dialog").locator('input[type="date"]').fill("2020-01-01")
    await page.getByRole("button", { name: "Agendar" }).click()

    await expect(
      page.locator("text=A data do agendamento deve ser hoje ou uma data futura.")
    ).toBeVisible()
  })

  test("should show a validation error when the appointment reason exceeds the maximum length", async ({ page }) => {
    await page.goto("/schedule")
    await selectStaffOnPage(page)
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    await fillAppointmentForm(page)
    await page.getByPlaceholder("Motivo da consulta...").fill("a".repeat(501))
    await page.getByRole("button", { name: "Agendar" }).click()

    await expect(
      page.locator("text=O motivo deve ter no máximo 500 caracteres.")
    ).toBeVisible()
  })
})
