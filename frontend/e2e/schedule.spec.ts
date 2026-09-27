import { test, expect } from "@playwright/test"
import { loginAsAdmin } from "./helpers"

const formatSlotTime = (totalMinutes: number): string =>
  `${String(Math.floor(totalMinutes / 60)).padStart(2, "0")}:${String(totalMinutes % 60).padStart(2, "0")}`

const formatDateValue = (dateValue: Date): string =>
  `${String(dateValue.getFullYear()).padStart(4, "0")}-${String(dateValue.getMonth() + 1).padStart(2, "0")}-${String(dateValue.getDate()).padStart(2, "0")}`

const getFutureAlignedSlot = (): { startTime: string; endTime: string; appointmentDate: string } => {
  const nowDate = new Date()
  const nowTotalMinutes = nowDate.getHours() * 60 + nowDate.getMinutes()
  const earliestStartTotalMinutes = Math.ceil((nowTotalMinutes + 5) / 15) * 15
  const latestEndTotalMinutes = 23 * 60 + 45
  const shiftToNextDay = earliestStartTotalMinutes + 30 > latestEndTotalMinutes
  const scheduleDate = shiftToNextDay ? new Date(nowDate.getTime() + 24 * 60 * 60 * 1000) : nowDate
  const startTotalMinutes = shiftToNextDay ? 9 * 60 : earliestStartTotalMinutes
  const endTotalMinutes = startTotalMinutes + 30
  return {
    startTime: formatSlotTime(startTotalMinutes),
    endTime: formatSlotTime(endTotalMinutes),
    appointmentDate: formatDateValue(scheduleDate),
  }
}

const escapeRegExp = (value: string): string => value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")

const getEditSlot = (
  prefilledStartTime: string,
  prefilledDate: string
): { startTime: string; endTime: string; appointmentDate: string } => {
  const [startHour, startMinute] = prefilledStartTime.split(":").map(Number)
  const shiftedStartTotalMinutes = startHour * 60 + startMinute + 30
  const shiftedEndTotalMinutes = shiftedStartTotalMinutes + 30
  if (shiftedEndTotalMinutes <= 23 * 60 + 45) {
    return {
      startTime: formatSlotTime(shiftedStartTotalMinutes),
      endTime: formatSlotTime(shiftedEndTotalMinutes),
      appointmentDate: prefilledDate,
    }
  }
  const nextDayDate = new Date(`${prefilledDate}T00:00:00`)
  nextDayDate.setDate(nextDayDate.getDate() + 1)
  return {
    startTime: "09:00",
    endTime: "09:30",
    appointmentDate: formatDateValue(nextDayDate),
  }
}

const fillAppointmentForm = async (page: import("@playwright/test").Page) => {
  const modalDialog = page.getByRole("dialog")
  const appointmentSlot = getFutureAlignedSlot()

  await modalDialog.locator('input[type="date"]').fill(appointmentSlot.appointmentDate)
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

  test("should keep the displayed range when navigating between calendar weeks", async ({ page }) => {
    await page.goto("/schedule")
    const toolbarTitle = page.locator(".fc-toolbar-title")
    const initialTitle = (await toolbarTitle.innerText()).toLowerCase()
    const initialTitlePattern = new RegExp(`^${escapeRegExp(initialTitle)}$`, "i")

    await page.locator(".fc-next-button").click()
    await expect(toolbarTitle).not.toHaveText(initialTitlePattern)
    await page.waitForTimeout(1200)
    await expect(toolbarTitle).not.toHaveText(initialTitlePattern)

    await page.locator(".fc-prev-button").click()
    await expect(toolbarTitle).toHaveText(initialTitlePattern)
  })

  test("should open appointment details when clicking an existing event", async ({ page }) => {
    await page.goto("/schedule")
    await page.getByRole("button", { name: "Novo Agendamento" }).click()
    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()
    await expect(page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" }).first()).toBeVisible()

    await page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" }).first().click()

    await expect(page.getByRole("heading", { name: "Detalhes do Agendamento" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Editar" })).toBeVisible()
    await expect(page.getByRole("button", { name: "Cancelar" })).toBeVisible()
  })

  test("should edit an appointment from the details dialog", async ({ page }) => {
    await page.goto("/schedule")
    await page.getByRole("button", { name: "Novo Agendamento" }).click()
    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()
    await expect(page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" }).first()).toBeVisible()

    await page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" }).first().click()
    await page.getByRole("button", { name: "Editar" }).click()

    await expect(page.getByRole("heading", { name: "Editar Agendamento" })).toBeVisible()
    const prefilledAppointmentDate = await page.getByRole("dialog").locator('input[type="date"]').inputValue()
    const prefilledStartTime = await page.getByRole("dialog").locator('select[name="startTime"]').inputValue()
    const alternateSlot = getEditSlot(prefilledStartTime, prefilledAppointmentDate)
    if (alternateSlot.appointmentDate !== prefilledAppointmentDate) {
      await page.getByRole("dialog").locator('input[type="date"]').fill(alternateSlot.appointmentDate)
    }
    await page.getByRole("dialog").locator('select[name="startTime"]').selectOption(alternateSlot.startTime)
    await page.getByRole("dialog").locator('select[name="endTime"]').selectOption(alternateSlot.endTime)
    await page.getByRole("button", { name: "Salvar Alterações" }).click()

    await expect(page.locator("text=Agendamento atualizado com sucesso!")).toBeVisible()
  })

  test("should cancel an appointment from the details dialog", async ({ page }) => {
    await page.goto("/schedule")
    await page.getByRole("button", { name: "Novo Agendamento" }).click()
    await fillAppointmentForm(page)
    await page.getByRole("button", { name: "Agendar" }).click()
    await expect(page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" }).first()).toBeVisible()

    await page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" }).first().click()
    await page.getByRole("button", { name: "Cancelar" }).click()

    await expect(page.getByRole("heading", { name: "Cancelar Agendamento" })).toBeVisible()
    await page.getByRole("button", { name: "Confirmar Cancelamento" }).click()

    await expect(page.locator("text=Agendamento cancelado com sucesso!")).toBeVisible()
    await expect(page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" })).toHaveCount(0)
  })
})
