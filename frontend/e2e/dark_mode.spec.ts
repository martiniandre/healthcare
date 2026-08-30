import { test, expect } from "@playwright/test"
import { loginAsDoctor } from "./helpers"

test.describe("Dark Mode Toggle", () => {
  test.beforeEach(async ({ page }) => {
    await page.emulateMedia({ colorScheme: "light" })
    await loginAsDoctor(page)
    await page.evaluate(() => localStorage.removeItem("healthcare.theme"))
    await page.reload()
    await expect(page.getByRole("button", { name: "Aparência" })).toBeVisible()
  })

  test("should toggle the dark theme class and persist the preference", async ({ page }) => {
    await expect(page.locator("html")).not.toHaveClass(/dark/)

    await page.getByRole("button", { name: "Aparência" }).click()

    await expect(page.locator("html")).toHaveClass(/dark/)
    expect(await page.evaluate(() => localStorage.getItem("healthcare.theme"))).toBe("dark")
    await expect(page.getByRole("button", { name: "Aparência" })).toHaveAttribute("aria-pressed", "true")

    await page.getByRole("button", { name: "Aparência" }).click()

    await expect(page.locator("html")).not.toHaveClass(/dark/)
    expect(await page.evaluate(() => localStorage.getItem("healthcare.theme"))).toBe("light")
    await expect(page.getByRole("button", { name: "Aparência" })).toHaveAttribute("aria-pressed", "false")
  })

  test("should keep the dark theme after a page reload", async ({ page }) => {
    await page.getByRole("button", { name: "Aparência" }).click()
    await expect(page.locator("html")).toHaveClass(/dark/)

    await page.reload()
    await expect(page.getByRole("button", { name: "Aparência" })).toBeVisible()
    await expect(page.locator("html")).toHaveClass(/dark/)
    expect(await page.evaluate(() => localStorage.getItem("healthcare.theme"))).toBe("dark")
  })
})