import { create } from "zustand"

type ThemeMode = "light" | "dark"

const THEME_STORAGE_KEY = "healthcare.theme"

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