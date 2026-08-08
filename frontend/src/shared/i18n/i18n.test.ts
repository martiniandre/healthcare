import { describe, it, expect, beforeEach } from "vitest"
import i18n, { createModuleTranslator } from "./i18n"

describe("createModuleTranslator", () => {
  beforeEach(async () => {
    await i18n.changeLanguage("pt-BR")
  })

  it("should resolve a bare key inside the module namespace", () => {
    const translatePatients = createModuleTranslator("patients")
    const translatedValue = translatePatients("validation.reportCodeReq")
    expect(translatedValue).toBe(i18n.t("patients.validation.reportCodeReq"))
    expect(translatedValue).not.toBe("validation.reportCodeReq")
  })

  it("should resolve an already-prefixed key as-is", () => {
    const translateAnalytics = createModuleTranslator("analytics")
    expect(translateAnalytics("analytics.days.mon")).toBe(translateAnalytics("days.mon"))
    expect(translateAnalytics("analytics.days.mon")).toBe(i18n.t("analytics.days.mon"))
  })

  it("should render the default value when the key is missing", () => {
    const translatePatients = createModuleTranslator("patients")
    expect(translatePatients("validation.nonexistentKey", "Fallback label")).toBe("Fallback label")
  })

  it("should fall back to the fallback language when the active language is unknown", async () => {
    await i18n.changeLanguage("de-DE")
    const translatePatients = createModuleTranslator("patients")
    expect(translatePatients("validation.reportCodeReq")).toBe(
      i18n.t("patients.validation.reportCodeReq", { lng: "pt-BR" })
    )
  })
})
