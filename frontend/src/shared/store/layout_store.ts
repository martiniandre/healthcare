import { create } from "zustand"

type ThemeMode = "light" | "dark"
type TableDensity = "comfortable" | "compact"

const THEME_STORAGE_KEY = "healthcare.theme"
const DENSITY_STORAGE_KEY = "healthcare.tableDensity"

const readStoredDensity = (): TableDensity => {
  if (typeof window === "undefined") {
    return "comfortable"
  }
  return window.localStorage.getItem(DENSITY_STORAGE_KEY) === "compact" ? "compact" : "comfortable"
}

const readStoredTheme = (): ThemeMode | null => {
  if (typeof window === "undefined") {
    return null
  }
  const storedTheme = window.localStorage.getItem(THEME_STORAGE_KEY)
  return storedTheme === "light" || storedTheme === "dark" ? storedTheme : null
}

const readSystemTheme = (): ThemeMode => {
  if (typeof window !== "undefined" && window.matchMedia?.("(prefers-color-scheme: dark)").matches === true) {
    return "dark"
  }
  return "light"
}

const resolveInitialTheme = (): ThemeMode => readStoredTheme() ?? readSystemTheme()

const applyThemeToDocument = (theme: ThemeMode) => {
  if (typeof document === "undefined") {
    return
  }
  document.documentElement.classList.toggle("dark", theme === "dark")
  document.documentElement.style.colorScheme = theme
}

interface LayoutState {
  isMobileSidebarOpen: boolean
  theme: ThemeMode
  tableDensity: TableDensity
  toggleTableDensity: () => void
  toggleMobileSidebar: () => void
  closeMobileSidebar: () => void
  setTheme: (theme: ThemeMode) => void
  toggleTheme: () => void
}

const initialTheme = resolveInitialTheme()

applyThemeToDocument(initialTheme)

export const useLayoutStore = create<LayoutState>((set) => ({
  isMobileSidebarOpen: false,
  theme: initialTheme,
  tableDensity: readStoredDensity(),
  toggleTableDensity: () =>
    set((state) => {
      const nextDensity = state.tableDensity === "compact" ? "comfortable" : "compact"
      window.localStorage.setItem(DENSITY_STORAGE_KEY, nextDensity)
      return { tableDensity: nextDensity }
    }),
  toggleMobileSidebar: () => set((state) => ({ isMobileSidebarOpen: !state.isMobileSidebarOpen })),
  closeMobileSidebar: () => set({ isMobileSidebarOpen: false }),
  setTheme: (theme) => {
    window.localStorage.setItem(THEME_STORAGE_KEY, theme)
    applyThemeToDocument(theme)
    set({ theme })
  },
  toggleTheme: () =>
    set((state) => {
      const nextTheme = state.theme === "dark" ? "light" : "dark"
      window.localStorage.setItem(THEME_STORAGE_KEY, nextTheme)
      applyThemeToDocument(nextTheme)
      return { theme: nextTheme }
    }),
}))