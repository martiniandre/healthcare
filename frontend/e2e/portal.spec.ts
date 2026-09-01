import { test, expect } from "@playwright/test"
import { loginAsPatient } from "./helpers"

test.describe("Patient Portal Module", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsPatient(page)
  })

  test("should show the patient's upcoming appointments in the portal", async ({ page }) => {
    await page.route("**/api/v1/appointments/my", async (networkRoute) => {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify([
          {
            id: "appt-1",
            patient_fhir_id: "fhir-pat-1",
            staff_id: "emp-1",
            starts_at: "2026-08-15T09:00:00Z",
            ends_at: "2026-08-15T09:30:00Z",
            status: "scheduled",
            reason: "Consulta de Rotina Geral",
            version: 1,
            created_at: "2026-08-01T10:00:00Z",
          },
          {
            id: "appt-2",
            patient_fhir_id: "fhir-pat-1",
            staff_id: "emp-1",
            starts_at: "2026-08-20T14:00:00Z",
            ends_at: "2026-08-20T14:30:00Z",
            status: "confirmed",
            reason: "Retorno Cardiológico",
            version: 1,
            created_at: "2026-08-02T10:00:00Z",
          },
        ]),
      })
    })

    await page.goto("/portal?tab=appointments")

    await expect(page.getByText("Consulta de Rotina Geral")).toBeVisible()
    await expect(page.getByText("Retorno Cardiológico")).toBeVisible()
    await expect(page.getByText("Agendado", { exact: true })).toBeVisible()
    await expect(page.getByText("Confirmado", { exact: true })).toBeVisible()
  })

  test("should show report version badge for patient reports", async ({ page }) => {
    await page.goto("/portal?tab=reports")

    await expect(page.getByText("Eletrocardiograma de Repouso")).toBeVisible()
    await expect(page.getByText("v2", { exact: true })).toBeVisible()
  })

  test("should show empty state when the patient has no appointments", async ({ page }) => {
    await page.route("**/api/v1/appointments/my", async (networkRoute) => {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify([]),
      })
    })

    await page.goto("/portal?tab=appointments")

    await expect(page.getByText("Nenhum agendamento encontrado.")).toBeVisible()
  })

  test("should render the encounter reason in the portal encounters list", async ({ page }) => {
    await page.route("**/api/v1/portal/encounters", async (networkRoute) => {
      await networkRoute.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify([
          {
            fhir_resource_id: "enc-1",
            status: "finished",
            reason_display: "Dor abdominal persistente",
            started_at: "2026-08-01T09:00:00Z",
            ended_at: "2026-08-01T09:30:00Z",
          },
          {
            fhir_resource_id: "enc-2",
            status: "finished",
            reason_display: "Retorno Cardiológico",
            started_at: "2026-08-05T14:00:00Z",
          },
        ]),
      })
    })

    await page.goto("/portal?tab=encounters")

    await expect(page.getByText("Dor abdominal persistente")).toBeVisible()
    await expect(page.getByText("Retorno Cardiológico")).toBeVisible()
  })
})
