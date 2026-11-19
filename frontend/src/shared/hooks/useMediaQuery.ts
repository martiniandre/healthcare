import { useSyncExternalStore } from "react"

const readMatch = (query: string): boolean =>
  typeof window !== "undefined" && typeof window.matchMedia === "function" ? window.matchMedia(query).matches : false

export const useMediaQuery = (query: string): boolean =>
  useSyncExternalStore(
    (onChange) => {
      if (typeof window === "undefined" || typeof window.matchMedia !== "function") return () => undefined
      const mediaQueryList = window.matchMedia(query)
      mediaQueryList.addEventListener("change", onChange)
      return () => mediaQueryList.removeEventListener("change", onChange)
    },
    () => readMatch(query),
    () => false
  )

export const TABLET_MEDIA_QUERY = "(min-width: 768px) and (max-width: 1023px)"
