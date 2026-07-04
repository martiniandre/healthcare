import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Logout Flow Module", () => {
  test("should logout and redirect to login page", async ({ page }) => {
    await loginAsDoctor(page)
    await page.getByRole("button", { name: "Sair" }).click()
    await expect(page).toHaveURL(/\/login$/)
  })

  test("should redirect to login when accessing a protected route after logout", async ({ page }) => {
    await loginAsDoctor(page)
    await page.getByRole("button", { name: "Sair" }).click()
    await expect(page).toHaveURL(/\/login$/)
    await page.goto("/")
    await expect(page).toHaveURL(/\/login$/)
  })
})
