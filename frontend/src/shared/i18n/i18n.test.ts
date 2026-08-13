import { describe, it, expect, beforeEach } from "vitest"
import i18n, { createModuleTranslator } from "./i18n"

describe("createModuleTranslator", () => {
  beforeEach(async () => {
    await i18n.changeLanguage("pt-BR")
  })

  it("should resolve a bare key inside the module namespace", () => {
    const translatePatients = createModuleTranslator("patients")
    const translatedValue = translatePatients("validation.reportCodeReq")
    expect(translatedValue).toBe(i18n.t("validation.reportCodeReq", { ns: "patients" }))
    expect(translatedValue).not.toBe("patients.validation.reportCodeReq")
  })

  it("should resolve an already-prefixed dotted key as-is", () => {
    const translateAnalytics = createModuleTranslator("analytics")
    expect(translateAnalytics("analytics.days.mon")).toBe(translateAnalytics("days.mon"))
    expect(translateAnalytics("analytics.days.mon")).toBe(i18n.t("days.mon", { ns: "analytics" }))
  })

  it("should resolve an already-prefixed colon key as-is", () => {
    const translatePatients = createModuleTranslator("patients")
    expect(translatePatients("patients:validation.reportCodeReq")).toBe(
      translatePatients("validation.reportCodeReq")
    )
  })

  it("should render the default value when the key is missing", () => {
    const translatePatients = createModuleTranslator("patients")
    expect(translatePatients("validation.nonexistentKey", "Fallback label")).toBe("Fallback label")
  })

  it("should fall back to the fallback language when the active language is unknown", async () => {
    await i18n.changeLanguage("de-DE")
    const translatePatients = createModuleTranslator("patients")
    expect(translatePatients("validation.reportCodeReq")).toBe(
      i18n.t("validation.reportCodeReq", { lng: "pt-BR", ns: "patients" })
    )
  })
})
