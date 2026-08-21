import { test, expect } from "@playwright/test"
import { loginAsAdmin, mockScheduleAPI } from "./helpers"

const currentMonthPrefix = (): string => {
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}`
}

const seedAppointmentsForCurrentMonth = (): Record<string, unknown>[] => {
  const monthPrefix = currentMonthPrefix()
  return [
    {
      id: "seed-appt-1",
      patient_fhir_id: "fhir-pat-1",
      staff_id: "fhir-emp-1",
      starts_at: `${monthPrefix}-10T09:00:00Z`,
      ends_at: `${monthPrefix}-10T09:30:00Z`,
      status: "scheduled",
      reason: "Consulta cardiológica",
      version: 1,
      created_at: `${monthPrefix}-01T08:00:00Z`,
    },
    {
      id: "seed-appt-2",
      patient_fhir_id: "fhir-pat-2",
      staff_id: "fhir-emp-1",
      starts_at: `${monthPrefix}-11T14:00:00Z`,
      ends_at: `${monthPrefix}-11T14:30:00Z`,
      status: "confirmed",
      reason: "Retorno clínico",
      version: 1,
      created_at: `${monthPrefix}-01T08:00:00Z`,
    },
    {
      id: "seed-appt-cancelled",
      patient_fhir_id: "fhir-pat-2",
      staff_id: "fhir-emp-1",
      starts_at: `${monthPrefix}-12T10:00:00Z`,
      ends_at: `${monthPrefix}-12T10:30:00Z`,
      status: "cancelled",
      reason: "Cancelada pelo paciente",
      version: 2,
      created_at: `${monthPrefix}-01T08:00:00Z`,
    },
  ]
}

test.describe("Schedule Month Calendar", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page)
  })

  test("should render seeded appointments as colored chips with patient names", async ({ page }) => {
    await mockScheduleAPI(page, seedAppointmentsForCurrentMonth())
    await page.goto("/schedule")
    await expect(page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" })).toBeVisible()
    await expect(page.locator(".fc-event").filter({ hasText: "Mariana Costa Silva" })).toBeVisible()
  })

  test("should never render cancelled appointments on the grid", async ({ page }) => {
    await mockScheduleAPI(page, seedAppointmentsForCurrentMonth())
    await page.goto("/schedule")
    await expect(page.locator(".fc-daygrid-body")).toBeVisible()
    await expect(page.locator(".fc-event").filter({ hasText: "Cancelada pelo paciente" })).toHaveCount(0)
  })

  test("should hide a staff member's events when their checkbox is toggled off", async ({ page }) => {
    await mockScheduleAPI(page, seedAppointmentsForCurrentMonth())
    await page.goto("/schedule")
    const doctorEvents = page.locator(".fc-event").filter({ hasText: "Guilherme de Souza Araujo" })
    await expect(doctorEvents.first()).toBeVisible()

    await page.getByRole("checkbox").first().uncheck()
    await expect(page.locator(".fc-event")).toHaveCount(0)

    await page.getByRole("checkbox").first().check()
    await expect(doctorEvents.first()).toBeVisible()
  })

  test("should collapse busy days behind a more link", async ({ page }) => {
    const monthPrefix = currentMonthPrefix()
    const busyDayAppointments = Array.from({ length: 5 }, (_, appointmentIndex) => ({
      id: `busy-appt-${appointmentIndex}`,
      patient_fhir_id: "fhir-pat-1",
      staff_id: "fhir-emp-1",
      starts_at: `${monthPrefix}-15T${String(9 + appointmentIndex).padStart(2, "0")}:00:00Z`,
      ends_at: `${monthPrefix}-15T${String(9 + appointmentIndex).padStart(2, "0")}:30:00Z`,
      status: "scheduled",
      reason: `Consulta ${appointmentIndex + 1}`,
      version: 1,
      created_at: `${monthPrefix}-01T08:00:00Z`,
    }))

    await mockScheduleAPI(page, busyDayAppointments)
    await page.goto("/schedule")

    await expect(page.locator(".fc-daygrid-more-link").first()).toBeVisible()
  })
})
