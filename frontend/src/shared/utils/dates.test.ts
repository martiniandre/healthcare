import { describe, it, expect, beforeEach } from "vitest"
import i18n from "../i18n/i18n"
import { formatDate, formatLongDate, formatTime, formatDateTime, formatRelativeTime } from "./dates"

const localDate = new Date(2026, 7, 15, 14, 30, 0)
const isoDate = "2026-08-15T12:00:00Z"

describe("dates", () => {
  beforeEach(async () => {
    await i18n.changeLanguage("pt-BR")
  })

  describe("formatDate", () => {
    it("formats numeric date following the selected locale", async () => {
      await i18n.changeLanguage("pt-BR")
      expect(formatDate(localDate)).toBe("15/08/2026")
    })

    it("switches ordering when the locale changes", async () => {
      await i18n.changeLanguage("en-US")
      expect(formatDate(localDate)).toBe("8/15/2026")
    })

    it("accepts ISO strings and Date objects", () => {
      expect(formatDate(isoDate)).toContain("2026")
      expect(formatDate(localDate)).toContain("2026")
    })
  })

  describe("formatLongDate", () => {
    it("uses localized month names for en-US", async () => {
      await i18n.changeLanguage("en-US")
      expect(formatLongDate(localDate)).toMatch(/Aug/)
    })

    it("uses localized month names for pt-BR", async () => {
      await i18n.changeLanguage("pt-BR")
      expect(formatLongDate(localDate)).toMatch(/ago/)
    })
  })

  describe("formatTime", () => {
    it("renders only hour and minute in the selected locale", async () => {
      await i18n.changeLanguage("pt-BR")
      expect(formatTime(localDate)).toBe("14:30")
    })

    it("adds period markers for en-US", async () => {
      await i18n.changeLanguage("en-US")
      expect(formatTime(localDate)).toMatch(/PM/)
    })
  })

  describe("formatDateTime", () => {
    it("combines date and time segments", async () => {
      await i18n.changeLanguage("en-US")
      const formatted = formatDateTime(localDate)
      expect(formatted).toMatch(/Aug/)
      expect(formatted).toMatch(/PM/)
    })
  })

  describe("formatRelativeTime", () => {
    it("renders now for timestamps under one minute old", async () => {
      await i18n.changeLanguage("pt-BR")
      expect(formatRelativeTime(new Date())).toBe("agora")
      await i18n.changeLanguage("en-US")
      expect(formatRelativeTime(new Date())).toBe("just now")
    })

    it("renders minutes in the selected locale", async () => {
      await i18n.changeLanguage("en-US")
      expect(formatRelativeTime(new Date(Date.now() - 5 * 60000))).toBe("5 minutes ago")
      await i18n.changeLanguage("pt-BR")
      expect(formatRelativeTime(new Date(Date.now() - 5 * 60000))).toBe("há 5 minutos")
    })

    it("renders hours in the selected locale", async () => {
      await i18n.changeLanguage("es-ES")
      expect(formatRelativeTime(new Date(Date.now() - 3 * 3600000))).toBe("hace 3 horas")
    })

    it("renders days in the selected locale", async () => {
      await i18n.changeLanguage("en-US")
      expect(formatRelativeTime(new Date(Date.now() - 2 * 86400000))).toBe("2 days ago")
    })

    it("falls back to a formatted date beyond one week", async () => {
      await i18n.changeLanguage("pt-BR")
      const tenDaysAgo = new Date(Date.now() - 10 * 86400000)
      expect(formatRelativeTime(tenDaysAgo)).toContain(String(tenDaysAgo.getFullYear()))
    })
  })

  describe("invalid input handling", () => {
    it("returns the raw value when the date cannot be parsed", () => {
      expect(formatDate("not-a-date")).toBe("not-a-date")
      expect(formatDateTime("garbage")).toBe("garbage")
      expect(formatRelativeTime(undefined as unknown as string)).toBe("undefined")
    })
  })

  describe("explicit locale override", () => {
    it("formats with the given locale instead of the active language", async () => {
      await i18n.changeLanguage("en-US")
      expect(formatDate(localDate, "pt-BR")).toBe("15/08/2026")
    })
  })
})
