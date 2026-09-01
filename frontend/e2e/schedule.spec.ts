import { test, expect } from "@playwright/test"
import { loginAsAdmin } from "./helpers"

const getFutureAlignedSlot = (): { startTime: string; endTime: string } => {
  const nowDate = new Date()
  const slotDate = new Date(nowDate.getTime() + 2 * 60 * 60 * 1000)
  const alignedMinutes = Math.ceil(slotDate.getMinutes() / 15) * 15
  slotDate.setMinutes(alignedMinutes, 0, 0)
  const formatSlotTime = (dateValue: Date): string =>
    `${String(dateValue.getHours()).padStart(2, "0")}:${String(dateValue.getMinutes()).padStart(2, "0")}`

  const crossesMidnight =
    slotDate.getDate() !== nowDate.getDate() ||
    slotDate.getMonth() !== nowDate.getMonth() ||
    slotDate.getFullYear() !== nowDate.getFullYear()

  let startHour = slotDate.getHours()
  let startMinute = slotDate.getMinutes()
  if (crossesMidnight) {
    startHour = 21
    startMinute = 0
  }
  let endTotalMinutes = startHour * 60 + startMinute + 30
  const lastEndMinute = 23 * 60 + 45
  if (endTotalMinutes > lastEndMinute) {
    startHour = 22
    startMinute = 0
    endTotalMinutes = 22 * 60 + 30
  }
  const endHour = Math.floor(endTotalMinutes / 60)
  const endMinute = endTotalMinutes % 60
  return {
    startTime: formatSlotTime(new Date(2020, 0, 1, startHour, startMinute)),
    endTime: formatSlotTime(new Date(2020, 0, 1, endHour, endMinute)),
  }
}

const fillAppointmentForm = async (page: import("@playwright/test").Page) => {
  const modalDialog = page.getByRole("dialog")
  const appointmentSlot = getFutureAlignedSlot()

  await modalDialog.locator('[role="combobox"]').first().click()
  await page.locator('[role="option"]', { hasText: "Guilherme de Souza Araujo" }).click()

  await modalDialog.locator('[role="combobox"]').nth(1).click()
  await page.locator('[role="option"]', { hasText: "Dr. André Silva de Araujo" }).click()

  await modalDialog.locator('select[name="startTime"]').selectOption(appointmentSlot.startTime)
  await modalDialog.locator('select[name="endTime"]').selectOption(appointmentSlot.endTime)
}

test.describe("Appointment Scheduling Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page)
  })

  test("should display the calendar agenda with staff filters", async ({ page }) => {
    await page.goto("/schedule")
    await expect(page.getByRole("heading", { name: "Agenda", exact: true })).toBeVisible()
    await expect(page.getByRole("button", { name: "Novo Agendamento" })).toBeVisible()
    await expect(page.getByText("Profissionais")).toBeVisible()
    await expect(page.locator(".fc-timegrid-body")).toBeVisible()
  })

  test("should book an appointment and show it as a chip on the week calendar", async ({ page }) => {
    await page.goto("/schedule")
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()

    await expect(page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" }).first()).toBeVisible()
  })

  test("should show conflict message when booking an overlapping slot", async ({ page }) => {
    await page.goto("/schedule")
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()
    await expect(page.locator(".fc-event").first()).toBeVisible()

    await page.getByRole("button", { name: "Novo Agendamento" }).click()
    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()

    await expect(page.locator("text=Conflito de horário")).toBeVisible()
  })

  test("should show a validation error when the appointment date is in the past", async ({ page }) => {
    await page.goto("/schedule")
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    await fillAppointmentForm(page)
    await page.getByRole("dialog").locator('input[type="date"]').fill("2020-01-01")
    await page.getByRole("button", { name: "Agendar" }).click()

    await expect(
      page.locator("text=A data do agendamento deve ser hoje ou uma data futura.")
    ).toBeVisible()
  })

  test("should block past dates in the appointment date picker", async ({ page }) => {
    await page.goto("/schedule")
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    const now = new Date()
    const expectedMin = `${String(now.getFullYear()).padStart(4, "0")}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`
    await expect(
      page.getByRole("dialog").locator('input[type="date"]')
    ).toHaveAttribute("min", expectedMin)
  })

  test("should show a validation error when the appointment reason exceeds the maximum length", async ({ page }) => {
    await page.goto("/schedule")
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    await fillAppointmentForm(page)
    await page.getByPlaceholder("Motivo da consulta...").fill("a".repeat(501))
    await page.getByRole("button", { name: "Agendar" }).click()

    await expect(
      page.locator("text=O motivo deve ter no máximo 500 caracteres.")
    ).toBeVisible()
  })

  test("should display an opaque backdrop when the appointment modal opens", async ({ page }) => {
    await page.goto("/schedule")
    await page.getByRole("button", { name: "Novo Agendamento" }).click()

    const backdrop = page.locator('[data-state="open"].bg-black\\/60')
    await expect(backdrop).toBeVisible()
    await expect(backdrop).toHaveCSS("opacity", "1")
  })
})
