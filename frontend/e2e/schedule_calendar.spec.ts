import { test, expect, type Page } from "@playwright/test"
import { loginAsAdmin, mockScheduleAPI } from "./helpers"

const isoDateOffset = (dayOffset: number): string => {
  const date = new Date()
  date.setDate(date.getDate() + dayOffset)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, "0")
  const day = String(date.getDate()).padStart(2, "0")
  return `${year}-${month}-${day}`
}

const seedAppointmentsInCurrentWeek = (): Record<string, unknown>[] => {
  const today = isoDateOffset(0)
  const tomorrow = isoDateOffset(1)
  return [
    {
      id: "seed-appt-1",
      patient_fhir_id: "fhir-pat-1",
      staff_id: "emp-1",
      starts_at: `${today}T10:00:00Z`,
      ends_at: `${today}T10:30:00Z`,
      status: "scheduled",
      reason: "Consulta cardiológica",
      version: 1,
      created_at: `${today}T08:00:00Z`,
    },
    {
      id: "seed-appt-2",
      patient_fhir_id: "fhir-pat-2",
      staff_id: "emp-1",
      starts_at: `${tomorrow}T14:00:00Z`,
      ends_at: `${tomorrow}T14:30:00Z`,
      status: "confirmed",
      reason: "Retorno clínico",
      version: 1,
      created_at: `${tomorrow}T08:00:00Z`,
    },
    {
      id: "seed-appt-cancelled",
      patient_fhir_id: "fhir-pat-2",
      staff_id: "emp-1",
      starts_at: `${today}T15:00:00Z`,
      ends_at: `${today}T15:30:00Z`,
      status: "cancelled",
      reason: "Cancelada pelo paciente",
      version: 2,
      created_at: `${today}T08:00:00Z`,
    },
  ]
}

test.describe("Schedule Calendar", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page)
  })

  const selectView = async (pageValue: Page, view: string) => {
    await pageValue.getByRole("combobox", { name: "Calendar view" }).selectOption(view)
  }

  test("should default to the week view with a view selector", async ({ page }) => {
    await mockScheduleAPI(page, [])
    await page.goto("/schedule")
    await expect(page.locator(".fc-timegrid-body")).toBeVisible()
    const viewSelector = page.getByRole("combobox", { name: "Calendar view" })
    await expect(viewSelector).toBeVisible()
    await expect(viewSelector).toHaveValue("week")
    await expect(viewSelector.locator("option")).toHaveText(["week", "day", "month", "year"])
  })

  test("should switch to the month view when the month option is selected", async ({ page }) => {
    await mockScheduleAPI(page, [])
    await page.goto("/schedule")
    await selectView(page, "month")
    await expect(page.locator(".fc-daygrid-body")).toBeVisible()
  })

  test("should switch to the year view when the year option is selected", async ({ page }) => {
    await mockScheduleAPI(page, [])
    await page.goto("/schedule")
    await selectView(page, "year")
    await expect(page.locator(".fc-multimonth")).toBeVisible()
    await expect(page.locator(".fc-multimonth-month")).toHaveCount(12)
  })

  test("should advance to the next week when the next button is clicked in week view", async ({ page }) => {
    await mockScheduleAPI(page, [])
    await page.goto("/schedule")
    await expect(page.locator(".fc-timegrid-body")).toBeVisible()
    const titleBefore = await page.locator(".fc-toolbar-title").innerText()
    await page.getByRole("button", { name: /next/i }).click()
    const titleAfter = await page.locator(".fc-toolbar-title").innerText()
    expect(titleAfter).not.toBe(titleBefore)
    await expect(page.locator(".fc-timegrid-body")).toBeVisible()
  })

  test("should render seeded appointments as colored chips with patient names", async ({ page }) => {
    await mockScheduleAPI(page, seedAppointmentsInCurrentWeek())
    await page.goto("/schedule")
    await expect(page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" })).toBeVisible()
    await expect(page.locator(".fc-event").filter({ hasText: "Mariana Costa Silva" })).toBeVisible()
  })

  test("should never render cancelled appointments on the grid", async ({ page }) => {
    await mockScheduleAPI(page, seedAppointmentsInCurrentWeek())
    await page.goto("/schedule")
    await expect(page.locator(".fc-timegrid-body")).toBeVisible()
    await expect(page.locator(".fc-event").filter({ hasText: "Cancelada pelo paciente" })).toHaveCount(0)
  })

  test("should hide a staff member's events when their checkbox is toggled off", async ({ page }) => {
    await mockScheduleAPI(page, seedAppointmentsInCurrentWeek())
    await page.goto("/schedule")
    const doctorEvents = page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" })
    await expect(doctorEvents.first()).toBeVisible()

    await page.getByRole("checkbox").first().uncheck()
    await expect(page.locator(".fc-event")).toHaveCount(0)

    await page.getByRole("checkbox").first().check()
    await expect(doctorEvents.first()).toBeVisible()
  })

  test("should collapse busy days behind a more link in the month view", async ({ page }) => {
    const today = isoDateOffset(0)
    const busyDayAppointments = Array.from({ length: 5 }, (_, appointmentIndex) => ({
      id: `busy-appt-${appointmentIndex}`,
      patient_fhir_id: "fhir-pat-1",
      staff_id: "emp-1",
      starts_at: `${today}T${String(9 + appointmentIndex).padStart(2, "0")}:00:00Z`,
      ends_at: `${today}T${String(9 + appointmentIndex).padStart(2, "0")}:30:00Z`,
      status: "scheduled",
      reason: `Consulta ${appointmentIndex + 1}`,
      version: 1,
      created_at: `${today}T08:00:00Z`,
    }))

    await mockScheduleAPI(page, busyDayAppointments)
    await page.goto("/schedule")
    await selectView(page, "month")

    await expect(page.locator(".fc-daygrid-more-link").first()).toBeVisible()
  })

  test("should open the create modal prefilled with the clicked time slot", async ({ page }) => {
    await mockScheduleAPI(page, seedAppointmentsInCurrentWeek())
    await page.goto("/schedule")
    await selectView(page, "day")

    const morningLane = page.locator('.fc-timegrid-slot-lane[data-time="09:00:00"]')
    await expect(morningLane).toBeVisible()
    await morningLane.scrollIntoViewIfNeeded()
    await morningLane.click({ position: { x: 80, y: 5 } })

    await expect(page.getByRole("dialog")).toBeVisible()
  })

  test("should reschedule an appointment when it is dragged to another slot", async ({ page }) => {
    await mockScheduleAPI(page, seedAppointmentsInCurrentWeek())
    await page.goto("/schedule")

    const appointmentEvent = page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" }).first()
    await expect(appointmentEvent).toBeVisible()
    const eventBox = (await appointmentEvent.boundingBox()) as { x: number; y: number; width: number; height: number }

    const dayColumn = page.locator(".fc-timegrid-col").nth(1)
    const columnBox = (await dayColumn.boundingBox()) as { x: number; y: number; width: number; height: number }
    const halfHourSlotHeight = columnBox.height / 48

    const rescheduleResponse = page.waitForResponse(
      (response) =>
        response.request().method() === "PUT" &&
        response.request().url().includes("/api/v1/appointments/") &&
        !response.request().url().includes("/cancel")
    )

    await page.mouse.move(eventBox.x + eventBox.width / 2, eventBox.y + eventBox.height / 2)
    await page.mouse.down()
    await page.mouse.move(eventBox.x + eventBox.width / 2, eventBox.y + eventBox.height / 2 + halfHourSlotHeight, {
      steps: 10,
    })
    await page.mouse.up()

    await expect(rescheduleResponse).resolves.toBeTruthy()
  })
})
