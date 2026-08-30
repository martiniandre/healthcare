import { describe, it, expect, beforeEach, vi } from "vitest"
import { useLayoutStore } from "./layout_store"

const THEME_STORAGE_KEY = "healthcare.theme"

const mockSystemTheme = (matches: boolean) => {
  const previousMatchMedia = window.matchMedia
  Object.defineProperty(window, "matchMedia", {
    writable: true,
    configurable: true,
    value: vi.fn().mockImplementation((query: string) => ({
      matches,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    })),
  })
  return () => {
    if (previousMatchMedia) {
      Object.defineProperty(window, "matchMedia", {
        writable: true,
        configurable: true,
        value: previousMatchMedia,
      })
    } else {
      delete (window as { matchMedia?: unknown }).matchMedia
    }
  }
}

describe("useLayoutStore theme", () => {
  beforeEach(() => {
    window.localStorage.clear()
    useLayoutStore.setState({ theme: "light" })
    document.documentElement.classList.remove("dark")
  })

  it("should apply the theme class and color scheme to the document", () => {
    useLayoutStore.getState().setTheme("dark")

    expect(document.documentElement.classList.contains("dark")).toBe(true)
    expect(document.documentElement.style.colorScheme).toBe("dark")

    useLayoutStore.getState().setTheme("light")

    expect(document.documentElement.classList.contains("dark")).toBe(false)
    expect(document.documentElement.style.colorScheme).toBe("light")
  })

  it("should persist the selected theme in localStorage", () => {
    useLayoutStore.getState().toggleTheme()

    expect(useLayoutStore.getState().theme).toBe("dark")
    expect(window.localStorage.getItem(THEME_STORAGE_KEY)).toBe("dark")

    useLayoutStore.getState().toggleTheme()

    expect(useLayoutStore.getState().theme).toBe("light")
    expect(window.localStorage.getItem(THEME_STORAGE_KEY)).toBe("light")
  })

  it("should initialize the theme from the stored preference", async () => {
    window.localStorage.setItem(THEME_STORAGE_KEY, "dark")
    vi.resetModules()

    const reloadedStore = await import("./layout_store")

    expect(reloadedStore.useLayoutStore.getState().theme).toBe("dark")
    expect(document.documentElement.classList.contains("dark")).toBe(true)
  })

  it("should fall back to the system preference when nothing is stored", async () => {
    const restoreSystemTheme = mockSystemTheme(true)
    vi.resetModules()

    const reloadedStore = await import("./layout_store")

    expect(reloadedStore.useLayoutStore.getState().theme).toBe("dark")
    expect(document.documentElement.classList.contains("dark")).toBe(true)
    restoreSystemTheme()
  })
})